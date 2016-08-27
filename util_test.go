package main

import (
	"log"
	"testing"
)

func Test_getURI(t *testing.T) {
	uris, err := getURI("http://www.6renyou.com")
	if err != nil {
		t.Fatal(err)
	}
	for _, uri := range uris {
		log.Print(uri)
	}
}
