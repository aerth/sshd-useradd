package main

import (
	"encoding/json"
	"fmt"
)

type ServerResponse struct {
	Message string `json:"message"`
}

type ServerStatus struct {
	Address      string `json:"ip"`
	Location     string `json:"location"`
	CurrentUsers int    `json:"currentUsers, int"` // api returns int for current users
	MaxUsers     string `json:"maxUsers, string"`  // api returns string for max users
	Coordinates  GPS    `json:"coordinates, string"`
	hostname     string `json:"-"`
}

type GPS struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

func decode(b []byte) string {
	var m ServerResponse
	err := json.Unmarshal(b, &m)
	if err != nil {
		println(err.Error())
		return "error decoding json: " + string(b)
	}
	return m.Message

}

func decodestatus(b []byte) map[string]ServerStatus {
	var status = make(map[string]ServerStatus)
	err := json.Unmarshal(b, &status)
	if err != nil {
		println(err.Error())
	}
	for k, v := range status {
		v.hostname = k
	}
	return status
}

func (s ServerStatus) String() string {
	var str string
	if s.hostname != "" {
		str += fmt.Sprintf("Host: %s\n", s.hostname)
	}
	str += fmt.Sprintf("Current Users: %v\nMax Users %s\n", s.CurrentUsers, s.MaxUsers)
	str += fmt.Sprintf("Location: %s (%v, %v)\n", s.Location, s.Coordinates.Latitude, s.Coordinates.Longitude)
	return str
}
