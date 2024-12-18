package worker

import (
	"errors"
	"fmt"
	"log"
	"rijojohn85/cube/task"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (w *Worker) GetTasks() []task.Task {
	var tasks []task.Task
	for _, v := range w.Db {
		tasks = append(tasks, *v)
	}
	return tasks
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	config := t.NewConfig()
	d := config.NewDocker()
	result := d.Run()
	if result.Error != nil {
		log.Printf("Error running task %v: %v", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t
	t.StartTime = time.Now().UTC()
	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := t.NewConfig()
	d := config.NewDocker()
	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", t.ContainerID, result.Error)
		return result
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Stopped and removed container for %v for task %v", t.ContainerID, t.ID)
	return result
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in queue.")
		return task.DockerResult{Error: nil}
	}
	taskQueued := t.(task.Task)
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}
	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("we should not be here, contact help")
		}
	} else {
		err := fmt.Errorf("invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}
	return result
}
