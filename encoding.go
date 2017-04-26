package main

import (
	"encoding/json"
)

type ServerResponse struct {
	Message string `json:"message"`
	
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
