package main

import (
	"github.com/hashicorp/go-hclog"
	"os"
	"os/exec"
)

type Instance struct {
	IPAddress	string
	Password	string
}

func main() {
	if configured() {
		BackupJob()
	} else {
		log().Warn("Missing required environmental variables, please rebuild me with the correct settings.")
	}
}

func configured() bool {

	// Verify we have the correct params
	isReady := true

	log().Info("Starting MySQL Dump Tool")

	// Just basic sanity checking
	if string(os.Getenv("METADATA_AUTH")) == "" {
		isReady = false
	}
	if string(os.Getenv("METADATA_URL")) == "" {
		isReady = false
	}

	// Ensure mysql and mysqldump are available
	_, err := exec.LookPath("mysql")
	if err != nil {
		log().Error("Missing executable mysql!")
		isReady = false
	}

	_, err = exec.LookPath("mysqldump")
	if err != nil {
		log().Error("Missing executable mysqldump!")
		isReady = false
	}

	return isReady
}

func log() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:  "mysql-dump-tool",
		Level: hclog.LevelFromString("INFO"),
		TimeFormat: "2006/01/02 15:04:05",
	})
}