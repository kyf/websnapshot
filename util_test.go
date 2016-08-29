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

	ss_dir := "/home/kyf/demo/snapshot/"
	for _, target := range uris {

		output, err := callWebHandler(target, ss_dir)
		if err != nil {
			t.Fatal(err)
		}
		log.Print(string(output))

	}
}

/*
func Test_callWebHandler(t *testing.T) {
	target, ss_dir := "http://im2.6renyou.com:6060/login", "/home/kyf/demo/snapshot/"

	output, err := callWebHandler(target, ss_dir)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(string(output))
}
*/
