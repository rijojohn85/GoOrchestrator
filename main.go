package main

import (
	"fmt"
	"log"
	"os"
	"rijojohn85/cube/task"
	"rijojohn85/cube/worker"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w)
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env:   []string{"POSTGRES_USER=cube", "POSTGRES_PASSWORD=secret"},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}
	fmt.Printf(
		"Container %s is running with config %v",
		result.ContainerId,
		c,
	)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
	}
	fmt.Printf(
		"Container %s has been stopped.",
		id,
	)
	return &result
}
