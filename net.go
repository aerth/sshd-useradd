package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func newuser(username, pubkey, host string) string {
	jsoncode := `{"user":"` + username + `","key":"` + pubkey + `","host":"` + host + `"}`
	body := strings.NewReader(jsoncode)
	println("sending request:")
	println(body)
	req, err := http.NewRequest("POST", "https://hashbang.sh/user/create", body)
	if err != nil {
		println(err.Error())
		return "error"
	}
	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("User-Agent", "sshd-adduser/"+version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return "error"
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(respbody)
}

func getter(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
		return "error"
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return "error"
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(respbody)

}
