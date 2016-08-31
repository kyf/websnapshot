package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"
)

type State struct {
	state bool
}

func NewState(state bool) State {
	return State{state}
}

func (this *State) isRun() bool {
	return this.state
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

type HandlerResponse struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
	Snapshot    string `json:"snapshot"`
}

type TaskManager struct {
	list     []Task
	state    State
	exit     chan int
	ssDir    string
	response map[int64][]HandlerResponse
	sync.Mutex
}

func NewTaskManager(ssDir string) *TaskManager {
	return &TaskManager{list: make([]Task, 0), state: NewState(false), exit: make(chan int, 1), ssDir: ssDir, response: make(map[int64][]HandlerResponse, 0)}
}

func (this *TaskManager) push(it Task) {
	this.Lock()
	defer this.Unlock()

	this.list = append(this.list, it)
}

func (this *TaskManager) pop() (it Task, ok bool) {
	if len(this.list) == 0 {
		ok = false
		return
	}

	this.Lock()
	defer this.Unlock()
	length := len(this.list)
	it = this.list[length-1]
	this.list = this.list[:length-1]
	ok = true
	return
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
		uris, err := getURI(it.target)
		if err != nil {
			log.Printf("getURI error:%v", err)
			continue
		}
		uris = append(uris, it.target)

		uris = []string{it.target}

		reslist := make([]HandlerResponse, 0)
		for index, uri := range uris {
			output, err := callWebHandler(uri, this.ssDir)

			if err != nil {
				log.Printf("callWebHandler error:%v", err)
				continue
			}
			var res HandlerResponse
			err = json.Unmarshal(output, &res)
			if err != nil {
				log.Printf("json.Unmarshal output error:%v, origin is [%s]", err, string(output))
				continue
			}
			res.Snapshot = strings.Replace(res.Snapshot, this.ssDir, "", -1)
			reslist = append(reslist, res)
			log.Printf("[%v-%s] has finished!", index, uri)
		}
		this.response[it.id] = reslist
		log.Printf("task[%v] has complete!", it.id)
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
