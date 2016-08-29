package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

const (
	WEBHANDLER = "webhandler"
)

func fillHost(path []byte, host string) string {
	spath := strings.TrimLeft(string(path), "/")
	if strings.HasPrefix(spath, "http://") {
		return spath
	}

	return host + "/" + spath
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

	_uri, err := url.Parse(target)
	if err != nil {
		return
	}

	host := _uri.Scheme + "://" + _uri.Host
	uniqs := make(map[string]bool, 0)

	matches := reg.FindAllSubmatch(body, -1)
	for _, match := range matches {
		for index, it := range match {
			if index == 0 {
				continue
			}
			fullpath := fillHost(it, host)
			if _, ok := uniqs[fullpath]; ok {
				continue
			}
			uris = append(uris, fullpath)
			uniqs[fullpath] = true
		}
	}

	return
}

func responseJson(w http.ResponseWriter, status bool, message string, data ...interface{}) {
	result := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func callWebHandler(target, ss_dir string) ([]byte, error) {
	cmd := exec.Command(WEBHANDLER, target, ss_dir)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
