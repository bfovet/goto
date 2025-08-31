package worker

import (
	"fmt"

	"github.com/bfovet/goto/internal/app/goto/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (worker *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (worker *Worker) RunTask() {
	fmt.Println("I will start or stop a task")
}

func (worker *Worker) StartTask() {
	fmt.Println("I will start a task")
}

func (worker *Worker) StopTask() {
	fmt.Println("I will stop a task")
}
