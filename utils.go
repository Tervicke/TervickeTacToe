package main

import (
	"encoding/json"
	"log"

	"github.com/olahol/melody"
)

func sendJSONMessage(s *melody.Session , eventType string , data map[string]string ){
	data["EVENT"] = eventType;
	message,err := json.Marshal(data);
	if err != nil {
		log.Printf("%v",err);
	}
	if err := s.Write(message); err != nil{
		log.Printf("%v",err);
	}
}
