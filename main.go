package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-martini/martini"
)

func prepare(r *http.Request) {
	r.ParseForm()
}

var (
	handlers = make(map[string]martini.Handler, 0)
	SS_DIR   string
)

func main() {
	flag.StringVar(&SS_DIR, "snapshot_dir", "", "store snapshot image directory")
	flag.Parse()

	if len(SS_DIR) == 0 {
		log.Fatal(color.RedString("snapshot_dir is not be initial"))
	}

	m := martini.Classic()
	m.Use(martini.Static("./template"))
	m.Use(martini.Static("./static"))
	m.Use(martini.Static(SS_DIR))
	m.Use(prepare)

	for k, v := range handlers {
		m.Get(k, v)
		m.Post(k, v)
	}

	taskch := make(chan Task, 1)

	tm := NewTaskManager(SS_DIR)
	err := tm.Run(taskch)
	if err != nil {
		log.Fatal(err)
	}
	defer tm.Stop()
	m.Map(taskch)
	m.Map(tm)

	log.Fatal(http.ListenAndServe(":5555", m))
}
