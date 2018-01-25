package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	storage "github.com/meinside/rpi-camera-timelapse-go/storage"
	camera "github.com/meinside/rpi-tools/hardware"
)

const (
	ConfigFilename = "./config.json"

	ImageExtension = "jpg"

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
	ImageWidth   int
	ImageHeight  int
	CameraParams map[string]interface{}
}

// variables
var captureChannel chan ShootRequest
var shootIntervalMinutes int
var imageWidth, imageHeight int
var cameraParams map[string]interface{}
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
		var loaded storage.Interface
		storageInterfaces = []storage.Interface{}
		for _, storageConf := range conf.StorageConfigs {
			switch storageConf.Type {
			case storage.TypeLocal:
				loaded = storage.NewLocalStorage(storageConf.Path)
			case storage.TypeDropbox:
				loaded = storage.NewDropboxStorage(
					storageConf.DropboxToken,
					storageConf.Path)
			case storage.TypeSmtp:
				loaded = storage.NewSmtpStorage(
					storageConf.SmtpEmail,
					storageConf.SmtpServer,
					storageConf.SmtpPasswd,
					storageConf.SmtpRecipients)
			default:
				log.Printf("*** Unknown storage type: %s\n", storageConf.Type)
				continue
			}

			log.Printf("Storage config loaded: %s\n", storageConf.Type)

			storageInterfaces = append(storageInterfaces, loaded)
		}
		if len(storageInterfaces) <= 0 {
			panic("No storages were configured.")
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
	if bytes, err := camera.CaptureRaspiStill(req.ImageWidth, req.ImageHeight, req.CameraParams); err == nil {
		// generate a filename with current timestamp
		filename := fmt.Sprintf("%.4f.%s", float64(time.Now().UnixNano())/float64(time.Millisecond), ImageExtension)

		// store captured image
		for _, storage := range storageInterfaces {
			if err := storage.Save(filename, bytes); err == nil {
				if isVerbose {
					log.Printf("Saved %d bytes to storage: %+v\n", len(bytes), storage)
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
