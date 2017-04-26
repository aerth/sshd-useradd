package main

import (
	"encoding/json"
)

type ServerResponse struct {
	Message string `json:"message"`
}

type ServerStatus struct {
	Address string `json:"ip"`
	Location string `json:"location"`
	CurrentUsers int `json:"currentUsers"`
	MaxUsers int `json:"maxUsers"`
	Coordinates GPS `json:"coordinates"`
}

type GPS struct {
	Latitude float32  `json:"lat"`
	Longitude float32 `json:"lon"`
}

func decode(b []byte) string {
	var m ServerResponse
	err := json.Unmarshal(b, &m)
	if err != nil {
		println(err.Error())
		return "{}"
	}
	return m.Message

}                  

func decodestatus(b []byte) map[string]ServerStatus {
	var	status = make(map[string]ServerStatus)
	err := json.Unmarshal(b, status)
	if err != nil {
		println(err.Error())
	}
	return status
}                                        
