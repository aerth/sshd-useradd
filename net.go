package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func newuser(username, pubkey, host string) string {

	body := strings.NewReader(`{"user":"` + username + `","key":"` + pubkey + `","host":"` + host + `"}`)
	req, err := http.NewRequest("POST", "https://hashbang.sh/user/create", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "sshd-adduser/"+version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(respbody)
}
