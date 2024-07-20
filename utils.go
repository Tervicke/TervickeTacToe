package main

import (
	"encoding/json"
	"log"

)

func sendJSONMessage(client_id string , eventType string , data map[string]string ){
	s := sessions[client_id];
	data["EVENT"] = eventType;
	message,err := json.Marshal(data);
	if err != nil {
		log.Printf("%v",err);
	}
	if err := s.Write(message); err != nil{
		log.Printf("%v",err);
	}
}

