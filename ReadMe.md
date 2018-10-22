
# Watchdog

## Install Go

Min version:  1.11.1

[https://golang.org/doc/install](https://golang.org/doc/install)

Please don't forgot to set GOPATH:

[https://github.com/golang/go/wiki/SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

## Configuration

### Get go package for sending email using gmail and deamonize watchdog

```golang
go get gopkg.in/gomail.v2
go get github.com/sevlyar/go-daemon
```

### Install watchdog

```bash
git clone https://github.com/almbox/go-study.git

```

Navigate in "go-study" folder

### Configure variables in watchdog.go

```golang
// challenge doc variables
const command string = "systemctl"
const serviceName string = "docker"
const checkInterval int = 2 // Seconds
const restartCount int = 2
const logFilePath string = "./watchdog.log"

// email conf
const sender string = "example@gmail.com"
const password string = "password"
const recipient string = "example@gmail.com"
const smtpServer string = "smtp.gmail.com"

// set to true if you want to simulate service fail
const debug = false
```

## Run package

sudo priviledges required to be able to start service

```golang
sudo /usr/local/go/bin/go run watchdog.go
```

### Kill daemon

Check the containing go-study folder for running watchdog "pid" file.

```bash
kill -9 <pid>
```

# Issue

google authentification level should be adjusted to "allow less secure app access"
https://myaccount.google.com/lesssecureapps?pli=1

# Warning

Logs are not rotated. You better to stop the deamon after test otherwise you disk will be running out of space.

This script had been developed and published for study purposes only!
