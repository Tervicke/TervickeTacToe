package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/olahol/melody"
)


var (
	Users []Client;
	Rooms []Room;
	m *melody.Session;
	sessions map[string]*melody.Session;
)

type Client struct{
	Id string;
	symbol string;
}

type Room struct{
	id string;
	player1 *Client;
	player2 *Client;
	gameBoard [9]string; 
	current string;
	gameover bool;
	player1replay bool;
	player2replay bool;
}


func main(){

	m := melody.New();
	sessions = make(map[string]*melody.Session)
	http.HandleFunc("/ping",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,"Pong");
	});

	http.HandleFunc("/ws",func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w,r);
	});

	//handle the connect
	m.HandleConnect(HandleConnect);

	//handle the disconnect	
	m.HandleDisconnect(HandleDisconnect)

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//HandleMessage(s,msg[] byte)	
		var data map[string]interface{}

		err := json.Unmarshal(msg,&data)

		if err != nil{
			fmt.Println("Error decoding")
			return
		}
		if event,ok := data["EVENT"]; ok{
			client,err:= GetClientById(s.Request.RemoteAddr);

			if err != nil {
				log.Printf("%v",err);
			}

			if(event == "CREATE-ROOM"){

				handleCreateEvent(s,client);

			}

			if(event == "JOIN-ROOM"){

				handleJoinEvent(data["ROOM_ID"].(string) , s );

			}
		
			if event == "GAMEMOVE"{

				handleGameMoveEvent(data["MOVE"].(string),s);

			}

			if event == "REPLAY"{
				HandleReplay(s);
			}

		}
	})
	port := os.Getenv("PORT");
	if port == ""{
		port = "5000"
	}
	fmt.Printf("Server running on %s",port);
	err := http.ListenAndServe(":"+port,nil);
	if err != nil{
		log.Fatal(err);
	}
}




/*
{
	EVENT : CREATE-ROOM
	CONTENT: NULL

	EVENT : JOIN-ROOM
	CONTENT: {
		ROOMID: id:
	}

	EVENT: GAME 
	CONTENT: {
		MOVE: A1,A2,A3 / B1 , B2 , B3 / C1 , C2 , C3
	}
}
*/
