package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	//"sync"

	"github.com/fatih/color"
)

func HandleCreate(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	target := r.Form.Get("target")

	if strings.EqualFold("", target) {
		responseJson(w, false, "target is empty")
		return
	}

	list := make([]string, 0)
	outputch := make(chan string, 1)

	//var wg sync.WaitGroup
	uris, err := getURI(target)
	if err != nil {
		logger.Print(color.RedString("getURI error: %v", err))
		responseJson(w, false, fmt.Sprintf("%v", err))
		return
	}

	uris = append(uris, target)

	go func() {
		for {
			select {
			case output := <-outputch:
				list = append(list, output)
			}
		}
	}()

	uris = uris[0:2]

	for _, uri := range uris {
		//wg.Add(1)
		//go func() {
		output, err := callWebHandler(uri, SS_DIR)
		if err != nil {
			logger.Print(color.RedString("[%s] error:%v", uri, err))
			return
		}
		outputch <- string(output)
		logger.Printf("[%v] success!", uri)
		logger.Print(string(output))
		//wg.Done()
		//}()
	}

	//wg.Wait()
	close(outputch)
	data := map[string]interface{}{
		"prefix": SS_DIR,
		"list":   list,
	}
	responseJson(w, true, "", data)
}

func init() {
	handlers["/create"] = HandleCreate
}
