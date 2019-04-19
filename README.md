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

## 2. How can I configure it?

You need to create your own config file.

Sample file is included, so feel free to copy it and change values as you want.

```bash
$ cp config.json.sample config.json
$ vi config.json
```

You can configure it to save files locally, send via SMTP, upload to Dropbox or S3 like this:

```json
"storages": [
	{
		"type": "local",
		"path": "/home/meinside/photos/timelapse"
	},
	{
		"type": "smtp",
		"smtp_recipients": "recipient-email-address1@outlook.com,recipient-email-address2@yahoo.com",
		"smtp_email": "sender-email-address@email.com",
		"smtp_passwd": "sender-email-password",
		"smtp_server": "sender.smtp-server.com:587"
	},
	{
		"type": "dropbox",
		"path": "/timelapse",
		"dropbox_token": "Tttttttt_oOOOOOOO-kkkkkkkk-eeeeeee_NNNNNNNN"
	},
	{
		"type": "s3",
		"path": "/optional/directory/in/bucket",
		"s3_bucket": "my-s3-bucket-name"
	}
]
```

When not needed, just remove the unwanted one from __storages__.

## 3. How can I build it?

### A. Manually

```bash
$ go get -d github.com/meinside/rpi-camera-timelapse-go
$ cd $GOPATH/src/github.com/meinside/rpi-camera-timelapse-go
$ go build
```

### B. With docker-compose

```bash
$ docker-compose build
```

## 4. How can I run it?

### A. Manually

Just execute the compiled binary:

```bash
$ ./rpi-camera-timelapse-go
```

If nothing goes wrong, images will be captured and stored periodically as you configured.

### B. With docker-compose

```bash
$ docker-compose run app
```

## 4. How can I run it as a service?

### A. With systemd

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

### B. With docker-compose

```bash
$ docker-compose up -d
```

## 5. How do I merge captured images to a timelapse video?

Use ffmpeg:

```bash
$ ffmpeg -framerate 30 -pattern_type glob -i '*.jpg' -c:v libx264 timelapse.mp4
```

## 998. Any trouble?

Please open an issue.

## 999. License?

MIT

