package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleCreate(w http.ResponseWriter, r *http.Request, logger *log.Logger, taskch chan Task) {
	target := r.Form.Get("target")

	if strings.EqualFold("", target) {
		responseJson(w, false, "target is empty")
		return
	}

	task := NewTask(target)

	taskch <- task
	responseJson(w, true, "", map[string]int64{"taskid": task.id})
}

func HandleProcess(w http.ResponseWriter, r *http.Request, logger *log.Logger, tm *TaskManager) {
	taskid := r.Form.Get("taskid")

	if strings.EqualFold("", taskid) {
		responseJson(w, false, "taskid is empty")
		return
	}

	_taskid, err := strconv.ParseInt(taskid, 10, 0)
	if err != nil {
		responseJson(w, false, "taskid is invalid")
		return
	}

	data := tm.response[_taskid]
	responseJson(w, true, "", data)
}

func init() {
	handlers["/create"] = HandleCreate
}
