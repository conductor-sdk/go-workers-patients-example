package main

import (
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	taskRunner := worker.NewTaskRunner(
		settings.NewAuthenticationSettings(
			os.Getenv("KEY"),
			os.Getenv("SECRET"),
		),
		settings.NewHttpSettings(
			os.Getenv("CONDUCTOR_SERVER_URL"),
		),
	)
	log.Info("Starting workers with conductor-go version 1.3.7")
	log.Info("Will connect to ", os.Getenv("CONDUCTOR_SERVER_URL"))

	taskRunner.StartWorker(
		"find_patient",
		FindPatientWorker,
		20,
		100,
	)
	taskRunner.StartWorker(
		"update_patient",
		UpdatePatientWorker,
		20,
		100,
	)
	taskRunner.WaitWorkers()
}
