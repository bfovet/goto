package main

import (
	"fmt"
	"time"

	"github.com/bfovet/goto/internal/app/goto/manager"
	"github.com/bfovet/goto/internal/app/goto/node"
	"github.com/bfovet/goto/internal/app/goto/task"
	"github.com/bfovet/goto/internal/app/goto/worker"
	"github.com/moby/moby/client"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	config := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	client, _ := client.NewClientWithOpts(client.FromEnv)
	docker := task.Docker{
		Client: client,
		Config: config,
	}

	result := docker.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, config)

	return &docker, &result
}

func stopContainer(docker *task.Docker, id string) *task.DockerResult {
	result := docker.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}

func main() {
	task1 := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	taskEvent := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      task1,
	}

	fmt.Printf("task: %v\n", task1)
	fmt.Printf("task event: %v\n", taskEvent)

	worker := worker.Worker{
		Name: "worker-1", Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker: %v\n", worker)
	worker.CollectStats()
	worker.RunTask()
	worker.StartTask()
	worker.StopTask()

	manager := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{worker.Name},
	}
	fmt.Printf("manager: %v\n", manager)
	manager.SelectWorker()
	manager.UpdateTasks()
	manager.SendWork()
	node := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", node)

	fmt.Printf("create a test container\n")
	dockerTask, createResult := createContainer()

	time.Sleep(time.Second * 5)

	fmt.Printf("stopping container %s\n", createResult.ContainerId)
	_ = stopContainer(dockerTask, createResult.ContainerId)
}
