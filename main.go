package main

import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

func prepare(r *http.Request) {
	r.ParseForm()
}

var (
	handlers = make(map[string]martini.Handler, 0)
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("./template"))
	m.Use(martini.Static("./static"))
	m.Use(prepare)

	log.Fatal(http.ListenAndServe(":5555", m))
}
