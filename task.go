package main

import (
	"errors"
	"time"
)

type State struct {
	state bool
}

func NewState(state bool) State {
	return State{state}
}

func (this *State) String() string {
	if this.state {
		return "running"
	} else {
		return "stop"
	}
}

type Task struct {
	id     int64
	target string
}

func NewTask(target string) Task {
	id := time.Now().UnixNano()
	return Task{id: id, target: target}
}

type TaskManager struct {
	list  []Task
	state State
	exit  chan int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{make([]Task, 0), NewState(false), make(chan int, 1)}
}

func (this *TaskManager) Run(taskCh <-chan Task) error {
	if this.state.state {
		return errors.New("task manager has started!")
	}
	this.state.state = true
	go func() {
		for it := range taskCh {
			this.push(it)
		}
	}()

	outCh := make(chan Task, 1)

	go func(out chan<- Task) {
		for {
			select {
			case <-this.exit:
				close(outCh)
				goto Exit
			default:
				if it, ok := this.pop(); ok {
					outCh <- it
				} else {
					time.Sleep(time.Second * 5)
				}
			}
		}
	Exit:
	}(outCh)

	for it := range outCh {
		callWebHandler()
	}

	return nil
}

func (this *TaskManager) State() State {
	return this.state
}

func (this *TaskManager) Stop() error {
	if !this.state.state {
		return errors.New("task manager has not been start!")
	}
	this.state.state = false
	this.exit <- 1
	return nil
}
