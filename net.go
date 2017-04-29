package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func newuser(username, pubkey, host string) string {
	jsoncode := `{"user":"` + username + `","key":"` + pubkey + `","host":"` + host + `"}`
	body := strings.NewReader(jsoncode)
	println("sending request...")
	req, err := http.NewRequest("POST", "https://hashbang.sh/user/create", body)
	if err != nil {
		println(err.Error())
		return "error"
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return "error"
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return decode(respbody)
}

func getter(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
		return "error"
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return "error"
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return "error"
	}
	resp.Body.Close()
	return string(respbody)

}

func getstatus(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
		return []byte("error")
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return []byte("error")
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return []byte("error")
	}
	resp.Body.Close()
	return respbody

}
