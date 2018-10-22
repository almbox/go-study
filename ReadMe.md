
# Watchdog

## Install Go

[https://golang.org/doc/install](https://golang.org/doc/install)

## Get Go Package for sending email using gmail and deamonize watchdog

```golang
go get gopkg.in/gomail.v2
go get github.com/sevlyar/go-daemon
```

## Configure variables

```golang
// challenge doc variables
const command string = "systemctl"
const serviceName string = "docker"
const checkInterval int = 2 // Seconds
const restartCount int = 2
const logFilePath string = "log/watchdog.log"

// email conf
const sender string = "example@gmail.com"
const password string = "password"
const recipient string = "example@gmail.com"
const smtpServer string = "smtp.gmail.com"

// set to true if you want to simulate service fail
const debug = false
```

## Run Package

```golang
go run watchdog.go
```

## Kill daemon

```bash
kill -9 <pid>
```

## Warning!

This script had been developed and published for study purposes only