package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rijojohn85/cube/task"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, req *http.Request) {
	var te task.TaskEvent
	err := json.NewDecoder(req.Body).Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error decoding json: %v", err)
		log.Println(msg)
		http.Error(w, msg, 400)
		return
	}
	a.Worker.AddTask(te.Task)
	log.Printf("Task %s added", te.Task.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, req *http.Request) {
	params := chi.URLParam(req, "taskID")
	if params == "" {
		msg := fmt.Errorf("invalid request, taskID not present")
		http.Error(w, msg.Error(), http.StatusBadRequest)
		return
	}
	taskID, err := uuid.Parse(params)
	if err != nil {
		msg := fmt.Errorf("invalid taskID")
		http.Error(w, msg.Error(), http.StatusBadRequest)
		return
	}
	taskToStop, ok := a.Worker.Db[taskID]
	if !ok {
		msg := fmt.Errorf("invalid taskID")
		http.Error(w, msg.Error(), http.StatusNotFound)
		return
	}
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Printf("Added task %v to stop container %v", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(http.StatusAccepted)
}
