package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sevlyar/go-daemon"
	"gopkg.in/gomail.v2"
)

// ========================================================
// Variables

// challenge doc variables
const command string = "systemctl"
const serviceName string = "docker"
const checkInterval int = 2 // Seconds
const restartCount int = 2
const logFilePath string = "watchdog.log"

// email conf
const sender string = "example@gmail.com"
const password string = "password"
const recipient string = "example@gmail.com"
const smtpServer string = "smtp.gmail.com"

// set to true if you want to simulate service fail
const debug = false

// ========================================================

func timeNow() string {
	return (time.Now()).Format(time.RFC3339)
}

// To terminate the daemon use:
//  kill `cat pid`
func main() {
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-watchdog]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	watchdog()
}

// log file
var (
	outfile, _ = os.Create(logFilePath) // update path for your needs
	l          = log.New(outfile, "", 0)
)

func watchdog() {
LOOP:
	for {
		// check interval
		time.Sleep(time.Second * time.Duration(checkInterval))

		// log file handling
		defer outfile.Close()

		// check service status
		// result: true = active; false = inactive
		status := checkStatus()
		if status != true {
			activated := false
			for i := 1; i <= restartCount; i++ {
				l.Printf("%s - ERROR: Restart attempt nr.: %d\n", timeNow(), i)
				l.Printf("%s - ERROR: Waiting %d seconds before restart ...\n", timeNow(), checkInterval)
				time.Sleep(time.Second * time.Duration(checkInterval))
				// start the service
				reborn()
				if debug {
					failForDebug()
				}
				// check service status after reborn
				if checkStatus() {
					activated = true
					l.Printf("%s - INFO: Service succesfully restart by attempt nr.: %d\n", timeNow(), i)
					break
				}
			}
			// send email with status
			if activated == true {
				SendMail(true)
			} else {
				SendMail(false)
				l.Printf("%s - ERROR: Terminating Daemon!", timeNow())
				break LOOP
			}

		} else {
			l.Printf("%s - INFO: Service active", timeNow())
		}
		if debug {
			failForDebug()
		}
	}
}

// get service status
func checkStatus() bool {
	cmd := exec.Command(command, "check", serviceName)
	out, _ := cmd.CombinedOutput()
	status := strings.TrimSpace(string(out))
	if strings.EqualFold(string(status), "active") {
		l.Printf("%s - INFO: Service status: %s\n", timeNow(), status)
		return true
	}
	l.Printf("%s - ERROR: Service status: %s\n", timeNow(), status)
	return false
}

// start service
func reborn() {
	reborn := exec.Command(command, "start", serviceName)
	_, rebornErr := reborn.CombinedOutput()
	if rebornErr != nil {
		l.Printf("%s - ERROR: exec.Command('"+command+" start "+serviceName+"') failed with %s\n", timeNow(), rebornErr)
		log := exec.Command(command, "status", serviceName+".service")
		out, _ := log.CombinedOutput()
		l.Printf("%s - ERROR: Error message: \n%s\n", timeNow(), string(out))
	}
}

// debug function to simulate failure (stop the service)
func failForDebug() {
	cmd := exec.Command(command, "stop", serviceName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		l.Printf("exec.Command('"+command+" stop "+serviceName+"') failed with %s\n", err)
	}
}

// SendMail function send email with attached log
func SendMail(status bool) {
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Watchdog Report")
	if status {
		m.SetBody("text/html", "INFO: Service is up after restart")
	} else {
		m.SetBody("text/html", "ERROR: Can not start service '"+serviceName+"'")
	}
	m.Attach(logFilePath)

	d := gomail.NewDialer(smtpServer, 587, sender, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		l.Printf("ERROR: Failed to send email %s\n", err)
	}
}
