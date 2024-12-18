package worker

import (
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

func (w *Worker) StartTask() {
	fmt.Println("I start a task")
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

func (w *Worker) RunTask() {
	fmt.Println("I start or stop a task")
}
