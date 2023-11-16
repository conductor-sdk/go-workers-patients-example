package main

import (
	"encoding/json"
	client2 "github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func runSyncWorkflow(w http.ResponseWriter, r *http.Request) {

	client := client2.NewAPIClient(settings.NewAuthenticationSettings(
		os.Getenv("KEY"),
		os.Getenv("SECRET"),
	),
		settings.NewHttpSettings(
			os.Getenv("CONDUCTOR_SERVER_URL"),
		),
	)
	log.Info("Request received")
	workflowExecutor := executor.NewWorkflowExecutor(client)

	startRequest := model.NewStartWorkflowRequest(
		"PatientWorkflow",
		1,
		"",
		make(map[string]interface{}))

	workflow, err := workflowExecutor.ExecuteWorkflow(startRequest, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Executed workflow ", workflow.WorkflowId)
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	response, err := json.Marshal(workflow)
	_, err = w.Write(response)
	if err != nil {
		log.Error("Failed to write response body for executed workflow, ", err.Error())
		return
	}
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
		100*time.Millisecond,
	)
	taskRunner.StartWorker(
		"update_patient",
		UpdatePatientWorker,
		20,
		100*time.Millisecond,
	)
	http.HandleFunc("/", runSyncWorkflow)
	go func() {
		err := http.ListenAndServe(":8083", nil)
		if err != nil {
			log.Error("aasd")
		}
	}()

	taskRunner.WaitWorkers()
}
