package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

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

func (worker *Worker) RunTask() task.DockerResult {
	t := worker.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}

	taskQueued := t.(task.Task)

	taskPersisted := worker.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		worker.Db[taskQueued.ID] = &taskQueued
	}

	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = worker.StartTask(taskQueued)
		case task.Completed:
			result = worker.StopTask(taskQueued)
		default:
			result.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}
	return result
}

func (worker *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	docker := task.NewDocker(config)
	result := docker.Run()
	if result.Error != nil {
		log.Printf("Err running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		worker.Db[t.ID] = &t
		return result
	}

	t.ContainerId = result.ContainerId
	t.State = task.Running
	worker.Db[t.ID] = &t

	return result
}

func (worker *Worker) StopTask(t task.Task) task.DockerResult {
	config := t.NewConfig(&t)
	docker := t.NewDocker(config)

	result := docker.Stop(t.ContainerId)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", t.ContainerId, result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	worker.Db[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v\n", t.ContainerId, t.ID)

	return result
}

func (worker *Worker) AddTask(t task.Task) {
	worker.Queue.Enqueue(t)
}
