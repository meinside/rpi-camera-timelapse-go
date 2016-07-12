# Raspberry Pi: Timelapse Camera Daemon

## 0. What is it for?

It is for capturing images with some interval, using Raspberry Pi Camera Module.

These captured images can be used as each frame of a timelapse video.

## 1. What do I need before running it?

You need:

* Raspberry Pi
* Raspberry Pi Camera Module enabled, and its cable correctly connected
* [golang installed on Raspberry Pi](https://github.com/meinside/rpi-configs/blob/master/bin/prep_go.sh)
* and this README.md.

## 2. How can I build it?

```bash
$ go get -u github.com/meinside/rpi-camera-timelapse-go
$ cd $GOPATH/src/github.com/meinside/rpi-camera-timelapse-go
$ go build
```

or

```bash
$ go get -u github.com/stacktic/dropbox
$ git clone https://github.com/meinside/rpi-camera-timelapse-go.git
$ cd rpi-camera-timelapse-go
$ go build
```

## 3. How can I run it?

You need to create your own config file.

Sample file is included, so feel free to copy it and change values as you want.

```bash
$ cp config.json.sample config.json
$ vi config.json
```

You can configure it to save files locally or on Dropbox like this:

```json
"storages": [
	{
		"type": "local",
		"path": "/home/meinside/photos/timelapse"
	},
	{
		"type": "dropbox",
		"path": "/timelapse",
		"key": "0a1b2c3d4e5f6g",
		"secret": "0987654321jihgfedcba",
		"token": "Tttttttt_oOOOOOOO-kkkkkkkk-eeeeeee_NNNNNNNN"
	}
]
```

When not needed, just remove it from the config file.

After the configuration is finished, just execute the binary:

```bash
$ ./rpi-cameera-timelapse-go
```

If nothing goes wrong, images will be captured and stored as you configured.

## 4. How can I run it as a service?

### systemd

```bash
$ sudo cp systemd/rpi-camera-timelapse-go.service /lib/systemd/system/
$ sudo vi /lib/systemd/system/rpi-camera-timelapse-go.service
```

and edit **User**, **Group**, **WorkingDirectory**, and **ExecStart** values.

You can simply start/stop it with:

```
$ sudo systemctl start rpi-camera-timelapse-go.service
$ sudo systemctl stop rpi-camera-timelapse-go.service
```

If you want to launch it automatically on boot:

```bash
$ sudo systemctl enable rpi-camera-timelapse-go.service
```

## 998. Any trouble?

Please open an issue.

## 999. License?

MIT
