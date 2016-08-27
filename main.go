package main

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/golang/glog"
)

func main() {
	m := martini.Classic()

	glog.Fatal(http.ListenAndServe(":5555", m))
}
