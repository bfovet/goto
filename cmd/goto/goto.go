package main

import (
	"fmt"
	"time"

	"github.com/bfovet/goto/internal/app/goto/task"
	"github.com/bfovet/goto/internal/app/goto/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)
	worker := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("Starting task")
	worker.AddTask(t)
	result := worker.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerId = result.ContainerId
	fmt.Printf("Task %s is running in container %s\n", t.ID, t.ContainerId)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 30)

	fmt.Printf("Stopping task %s\n", t.ID)
	t.State = task.Completed
	worker.AddTask(t)
	result = worker.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
