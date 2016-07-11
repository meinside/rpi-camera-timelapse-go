package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	storage "github.com/meinside/rpi-camera-timelapse-go/storage"
)

const (
	// constants for config
	ConfigFilename = "./config.json"

	// absolute path of raspistill
	RaspiStillBin = "/usr/bin/raspistill"

	// temp directory
	RecommendedTempDir = "/var/tmp" // 'tmpfs /var/tmp tmpfs nodev,nosuid,size=10M 0 0' in /etc/fstab
	DefaultTempDir     = "/tmp"

	NumQueue                    = 4
	DefaultShootIntervalMinutes = 1
	MinImageWidth               = 400
	MinImageHeight              = 300
)

type config struct {
	ShootIntervalMinutes int                    `json:"shoot_interval_minutes"`
	ImageWidth           int                    `json:"image_width"`
	ImageHeight          int                    `json:"image_height"`
	CameraParams         map[string]interface{} `json:"camera_params"`
	StorageConfigs       []storage.Config       `json:"storages"`
	IsVerbose            bool                   `json:"is_verbose"`
}

// for making sure the camera is not used simultaneously
var cameraLock sync.Mutex

type ShootRequest struct {
	Directory    string
	ImageWidth   int
	ImageHeight  int
	CameraParams map[string]interface{}
}

// variables
var captureChannel chan ShootRequest
var shootIntervalMinutes int
var imageWidth, imageHeight int
var cameraParams map[string]interface{}
var tmpDir string
var storageInterfaces []storage.Interface
var isVerbose bool

// Read config
func getConfig() (config, error) {
	_, filename, _, _ := runtime.Caller(0) // = __FILE__

	if file, err := ioutil.ReadFile(filepath.Join(path.Dir(filename), ConfigFilename)); err == nil {
		var conf config
		if err := json.Unmarshal(file, &conf); err == nil {
			return conf, nil
		} else {
			return config{}, err
		}
	} else {
		return config{}, err
	}
}

func init() {
	// read config
	if conf, err := getConfig(); err != nil {
		panic(err)
	} else {
		// temporary directory
		tmpDir = RecommendedTempDir
		if _, err := os.Stat(tmpDir); err != nil {
			if os.IsNotExist(err) { // file does not exist
				tmpDir = DefaultTempDir
			}
		}
		log.Printf("Using temporary directory: %s\n", tmpDir)

		// interval
		shootIntervalMinutes = conf.ShootIntervalMinutes
		if shootIntervalMinutes <= 0 {
			shootIntervalMinutes = DefaultShootIntervalMinutes
		}

		// image width * height
		imageWidth = conf.ImageWidth
		if imageWidth < MinImageWidth {
			imageWidth = MinImageWidth
		}
		imageHeight = conf.ImageHeight
		if imageHeight < MinImageHeight {
			imageHeight = MinImageHeight
		}

		// other camera params
		cameraParams = conf.CameraParams

		// storage configurations
		storageInterfaces = []storage.Interface{}
		for _, storageConf := range conf.StorageConfigs {
			switch storageConf.Type {
			case storage.TypeLocal:
				storageInterfaces = append(storageInterfaces, storage.NewLocalStorage(storageConf.Path))
			case storage.TypeDropbox:
				storageInterfaces = append(storageInterfaces, storage.NewDropboxStorage(
					storageConf.Key,
					storageConf.Secret,
					storageConf.Token,
					storageConf.Path))
			}

			log.Printf("Read storage config: %s\n", storageConf.Type)
		}

		// show verbose messages or not
		isVerbose = conf.IsVerbose
	}
}

// capture
func capture(req ShootRequest) bool {
	// process result
	result := false

	cameraLock.Lock()
	defer cameraLock.Unlock()

	// capture image
	if filepath, err := captureImage(req.Directory, req.ImageWidth, req.ImageHeight, req.CameraParams); err == nil {
		defer removeImage(filepath)

		// store captured image
		for _, storage := range storageInterfaces {
			if err := storage.Save(&filepath); err == nil {
				if isVerbose {
					log.Printf("Saved %s to storage: %+v\n", filepath, storage)
				}
				result = true
			} else {
				log.Printf("*** Failed to store image: %s\n", err)
			}
		}
	} else {
		log.Printf("*** Image capture failed: %s\n", err)
	}

	return result
}

// capture an image with given width, height, and other parameters
// return the captured image's filepath (for deleting it after use)
func captureImage(directory string, width, height int, cameraParams map[string]interface{}) (filepath string, err error) {
	// filepath
	filepath = fmt.Sprintf("%s/captured_%d.jpg", directory, time.Now().UnixNano()/int64(time.Millisecond))

	// command line arguments
	args := []string{
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		"-o", filepath,
	}
	for k, v := range cameraParams {
		args = append(args, k)
		if v != nil {
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	// execute command
	if bytes, err := exec.Command(RaspiStillBin, args...).CombinedOutput(); err != nil {
		log.Printf("*** Error running %s: %s\n", RaspiStillBin, string(bytes))
		return "", err
	} else {
		if isVerbose {
			log.Printf("Captured image: %s\n", filepath)
		}
		return filepath, nil
	}
}

// remove temporary image
func removeImage(filepath string) {
	if err := os.Remove(filepath); err != nil {
		log.Printf("*** Failed to delete temp file: %s\n", err)
	}
}

func main() {
	log.Printf("Starting up...")

	timer := time.NewTicker(time.Duration(shootIntervalMinutes) * time.Minute)
	quitter := make(chan struct{})

	// catch SIGINT and SIGTERM and terminate gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		quitter <- struct{}{}
	}()

	// infinite loop
	for {
		select {
		case <-timer.C:
			capture(ShootRequest{
				Directory:    tmpDir,
				ImageWidth:   imageWidth,
				ImageHeight:  imageHeight,
				CameraParams: cameraParams,
			})
		case <-quitter:
			log.Printf("Shutting down...")
			os.Exit(1)
		}
	}
}
