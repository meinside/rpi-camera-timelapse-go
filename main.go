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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	storage "github.com/meinside/rpi-camera-timelapse-go/storage"
	camera "github.com/meinside/rpi-tools/hardware"
)

const (
	configFilename = "./config.json"

	imageExtension = "jpg"

	defaultShootIntervalMinutes = 1
	minImageWidth               = 400
	minImageHeight              = 300
)

type config struct {
	ShootWithinHours     string                 `json:"shoot_within_hours"` // e.g. 13-18
	ShootIntervalMinutes int                    `json:"shoot_interval_minutes"`
	ImageWidth           int                    `json:"image_width"`
	ImageHeight          int                    `json:"image_height"`
	CameraParams         map[string]interface{} `json:"camera_params"`
	StorageConfigs       []storage.Config       `json:"storages"`
	IsVerbose            bool                   `json:"is_verbose"`
}

// for making sure the camera is not used simultaneously
var cameraLock sync.Mutex

// ShootRequest type
type ShootRequest struct {
	ImageWidth   int
	ImageHeight  int
	CameraParams map[string]interface{}
}

// ShootWithinHoursConfig type
type ShootWithinHoursConfig struct {
	From int
	To   int
}

// ShouldCapture checks if it is a right time to capture.
func (hoursConfig ShootWithinHoursConfig) ShouldCapture(now time.Time) bool {
	return now.Hour() >= hoursConfig.From && now.Hour() <= hoursConfig.To
}

// variables
var captureChannel chan ShootRequest
var shootIntervalMinutes int
var shootWithinHours ShootWithinHoursConfig
var imageWidth, imageHeight int
var cameraParams map[string]interface{}
var storageInterfaces []storage.Interface
var isVerbose bool

// Read config
func getConfig() (config, error) {
	_, filename, _, _ := runtime.Caller(0) // = __FILE__

	file, err := ioutil.ReadFile(filepath.Join(path.Dir(filename), configFilename))

	if err == nil {
		var conf config
		err = json.Unmarshal(file, &conf)

		if err == nil {
			return conf, nil
		}
	}

	return config{}, err
}

func init() {
	// read config
	if conf, err := getConfig(); err != nil {
		panic(err)
	} else {
		// interval
		shootIntervalMinutes = conf.ShootIntervalMinutes
		if shootIntervalMinutes <= 0 {
			shootIntervalMinutes = defaultShootIntervalMinutes
		}

		// shoot within hours: from-to
		shootWithinHours = interpretWithinHours(conf.ShootWithinHours)

		// image width * height
		imageWidth = conf.ImageWidth
		if imageWidth < minImageWidth {
			imageWidth = minImageWidth
		}
		imageHeight = conf.ImageHeight
		if imageHeight < minImageHeight {
			imageHeight = minImageHeight
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
			case storage.TypeSMTP:
				loaded = storage.NewSMTPStorage(
					storageConf.SMTPEmail,
					storageConf.SMTPServer,
					storageConf.SMTPPasswd,
					storageConf.SMTPRecipients,
				)
			case storage.TypeS3:
				loaded = storage.NewS3Storage(storageConf.S3Bucket, storageConf.Path)
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

// parse 'shoot_within_hours' option
func interpretWithinHours(delimitedWithinHours string) ShootWithinHoursConfig {
	hours := strings.Split(delimitedWithinHours, "-")

	if len(hours) != 2 {
		return ShootWithinHoursConfig{0, 24}
	}

	from, err := strconv.Atoi(hours[0])
	if err != nil {
		from = 0
		log.Println("Invalid shoot_within_hours from, will use 0")
	}

	to, err := strconv.Atoi(hours[1])
	if err != nil {
		to = 24
		log.Println("Invalid shoot_within_hours to, will use 24")
	}

	return ShootWithinHoursConfig{from, to}
}

// capture
func capture(req ShootRequest) bool {
	// process result
	result := false

	cameraLock.Lock()
	defer cameraLock.Unlock()

	if !shootWithinHours.ShouldCapture(time.Now()) {
		if isVerbose {
			log.Println("Aborting capture as not within configured shooting hours")
		}
		return result
	}

	// capture image
	if bytes, err := camera.CaptureRaspiStill(req.ImageWidth, req.ImageHeight, req.CameraParams); err == nil {
		// generate a filename with current timestamp
		filename := fmt.Sprintf("%s.%s", time.Now().Format(time.RFC3339), imageExtension)

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
	log.Println("Starting up...")

	timer := time.NewTicker(time.Duration(shootIntervalMinutes) * time.Minute)
	quitter := make(chan struct{})

	// catch SIGINT and SIGTERM and terminate gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		quitter <- struct{}{}
	}()

	// capture a photo immediately before starting the infinite loop
	capture(ShootRequest{
		ImageWidth:   imageWidth,
		ImageHeight:  imageHeight,
		CameraParams: cameraParams,
	})

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
			log.Println("Shutting down...")
			os.Exit(1)
		}
	}
}
