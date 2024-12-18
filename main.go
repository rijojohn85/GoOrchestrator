package main

import (
	"fmt"
	"os"
	"rijojohn85/cube/task"
	"time"

	"github.com/docker/docker/client"
)

func main() {
	d, result := createContainer()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		os.Exit(1)
	}
	fmt.Printf("Container is running\n")
	time.Sleep(time.Second * 5)
	result = stopContainer(d, result.ContainerId)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		os.Exit(1)
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
