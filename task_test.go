package main

import (
	"log"
	"testing"
	"time"
)

func TestTaskManager(t *testing.T) {
	ss_dir := "/home/kyf/demo/snapshot/"
	tm := NewTaskManager(ss_dir)

	taskCh := make(chan Task, 1)

	go tm.Run(taskCh)

	defer tm.Stop()

	//task := NewTask("http://www.6renyou.com")
	task := NewTask("https://www.baidu.com")

	var taskid int64 = task.id

	select {
	case <-time.After(time.Second * 1):
		taskCh <- task
		log.Print("start new task ...")
	}

	select {
	case <-time.After(time.Second * 180):
		res := tm.response[taskid]
		for _, it := range res {
			log.Print(it)
		}
		log.Print("will be finish!")
	}
}
