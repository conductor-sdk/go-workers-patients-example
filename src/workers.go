package main

import (
	"database/sql"
	"errors"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	_ "github.com/lib/pq"
	"regexp"
)

func FindPatientWorker(task *model.Task) (result interface{}, err error) {

	taskResult := model.NewTaskResultFromTask(task)
	taskResult.Status = model.FailedTask
	connStr, ok := task.InputData["DBConnectionString"]
	if !ok {
		err := errors.New("DBConnectionString is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findFirstName, ok := task.InputData["first_name"]
	if !ok {
		err := errors.New("first_name is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findLastName, ok := task.InputData["last_name"]
	if !ok {
		err := errors.New("last_name is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findDob, ok := task.InputData["dob"]
	if !ok {
		err := errors.New("dob is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	tableName, ok := task.InputData["table"].(string)
	if !ok {
		err := errors.New("table is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(tableName)
	if !isAlphanumeric {
		err = errors.New("table name is invalid")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}

	db, err := sql.Open("postgres", connStr.(string))
	if db != nil {
		defer db.Close()
	}

	if err != nil {
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	queryString := " SELECT first_name, last_name, dob::text, family_doctor_assigned FROM  " + tableName +
		" WHERE first_name = $1 AND last_name = $2 AND dob = $3::date "
	rows, err := db.Query(queryString,
		findFirstName, findLastName, findDob)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	var firstName, lastName, dob string
	var doctorAssigned bool
	if rows.Next() {
		err := rows.Scan(&firstName, &lastName, &dob, &doctorAssigned)
		if err != nil {
			taskResult.ReasonForIncompletion = err.Error()
			return taskResult, err
		}
	} else {
		err = errors.New("Patient not found")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	taskResult.OutputData = map[string]interface{}{"first_name": firstName, "last_name": lastName,
		"dob": dob, "family_doctor_assigned": doctorAssigned}
	taskResult.Status = model.CompletedTask
	return taskResult, nil
}

func UpdatePatientWorker(task *model.Task) (result interface{}, err error) {

	taskResult := model.NewTaskResultFromTask(task)
	taskResult.Status = model.FailedTask
	connStr, ok := task.InputData["DBConnectionString"]
	if !ok {
		err := errors.New("DBConnectionString is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findFirstName, ok := task.InputData["first_name"]
	if !ok {
		err := errors.New("first_name is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findLastName, ok := task.InputData["last_name"]
	if !ok {
		err := errors.New("last_name is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	findDob, ok := task.InputData["dob"]
	if !ok {
		err := errors.New("dob is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	tableName, ok := task.InputData["table"].(string)
	if !ok {
		err := errors.New("table is missing")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(tableName)
	if !isAlphanumeric {
		err = errors.New("table name is invalid")
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	db, err := sql.Open("postgres", connStr.(string))
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	_, err = db.Exec(" UPDATE "+tableName+" SET family_doctor_assigned = true "+
		"WHERE first_name = $1 AND last_name = $2 AND dob = $3::date ",
		findFirstName, findLastName, findDob)
	if err != nil {
		taskResult.ReasonForIncompletion = err.Error()
		return taskResult, err
	}
	taskResult.OutputData = map[string]interface{}{"first_name": findFirstName, "last_name": findLastName,
		"dob": findDob, "family_doctor_assigned": true}
	taskResult.Status = model.CompletedTask
	return taskResult, nil

}
