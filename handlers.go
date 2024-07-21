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

	id := s.Request.RemoteAddr;
	for i,room := range Rooms{
		if room.player1.Id == id || room.player2.Id == id{
			sendExitMessage(&Rooms[i]);
			Rooms = append(Rooms[:i], Rooms[i+1:]...)
		}
	}
		for i,c := range Users{
			if c.Id == id{
				Users = append(Users[:i],Users...);
				fmt.Println("Client Disconnected");
				delete(sessions,s.Request.RemoteAddr);
				break;
			}
		}
}
func sendExitMessage(room *Room){ 
	fmt.Println("---------------------------");
	fmt.Println(room.player1);
	fmt.Println(room.player1);
	sendJSONMessage(room.player1.Id,"EXIT",map[string]string{});
	sendJSONMessage(room.player2.Id,"EXIT",map[string]string{});
	fmt.Println("---------------------------");
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

			if event ==  "JOIN-ROOM"{
				handleJoinEvent(data["ROOM_ID"].(string) , s );
			}
			
			if event == "GAMEMOVE"{
				handleGameMoveEvent(data["MOVE"].(string) , s);
			}

			if event == "REPLAY"{
				HandleReplay(s);
			}

				
		}
}

