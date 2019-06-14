package main

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func BackupJob() {

	log().Info("Starting backup task")

	instances, loadErr := loadContainers()
	if loadErr != nil {
		log().Warn("Error", loadErr.Error())
		return
	}

	for _, i := range instances {
		performBackup(&i)
	}
	return

}

func performBackup(instance *Instance) {

	var backupCommand []string

	t := time.Now()
	timestamp := t.Format("2006-01-02-1504")
	outputFileName := "${dbname}_" + timestamp

	mysqlConnect := "-h " + instance.IPAddress + " -u root -p" + instance.Password
	backupCommand = append(backupCommand, "mysql", mysqlConnect, "-N", "-e", "'show databases;'")
	backupCommand = append(backupCommand, "|")
	backupCommand = append(backupCommand, "while", "read", "dbname;", "do")
	backupCommand = append(backupCommand, "mysqldump", mysqlConnect, "--skip-lock-tables", "\"$dbname\" > " + outputFileName + ".sql")
	backupCommand = append(backupCommand, "&&", "tar", "-czvf", outputFileName + ".tar.gz", outputFileName + ".sql;", "done")

	//cmdExec := exec.Command("bash", "-c", strings.Join(backupCommand, " "))

	cmd := exec.Command("ash", "-c", strings.Join(backupCommand, " "))
	cmd.Env = os.Environ()
	cmd.Dir = "/tmp"
	out, err := cmd.CombinedOutput()

	if err != nil {
		log().Error("Fatal error running mysql dump", "error", err.Error())
	} else {
		log().Info("Mysql backup output", "output", string(out))
	}

}