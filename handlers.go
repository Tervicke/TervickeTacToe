package main

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/olahol/melody"
)

func HandleConnect(s *melody.Session){
	
		sessions[s.Request.RemoteAddr] = s	
		Users = append(Users,Client{Id:s.Request.RemoteAddr} );	
		fmt.Println(s.Request.RemoteAddr, " Client connected")

}
func HandleDisconnect(s *melody.Session){
		for i,c := range Users{
			if c.Id == s.Request.RemoteAddr{
				Users = append(Users[:i],Users...);
				fmt.Println("Client Disconnected");
				delete(sessions,s.Request.RemoteAddr);
				break;
			}
		}
}

func HandleMessage(s *melody.Session , msg []byte ){
	
		var data map[string]interface{};

		err := json.Unmarshal(msg,&data)

		if err != nil{
			log.Printf("%v", err);
			return;
		}
		
		if event,ok := data["EVENT"]; ok{
			client,err := GetClientById(s.Request.RemoteAddr);

			if err != nil {
				log.Printf("%v", err);
			}

			if event == "CREATE-ROOM"{
				handleCreateEvent(s,client);
			}
				
		}
}
