package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func fillHost(path []byte, host string) string {
	if strings.HasPrefix(string(path), "http://") {
		return string(path)
	}

}

func getURI(target string) (uris []string, err error) {
	res, err := http.Get(target)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	reg, err := regexp.Compile(`href=["|']([^"']+)["|']`)
	if err != nil {
		return
	}
	matches := reg.FindAllSubmatch(body, -1)
	for _, match := range matches {
		for index, it := range match {
			if index == 0 {
				continue
			}
			uris = append(uris, fillHost(it))
		}
	}

	return
}
