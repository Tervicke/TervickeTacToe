package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		
		var data map[string]interface{}

		err := json.Unmarshal(msg,&data)

		if err != nil{
			fmt.Println("Error decoding")
			return
		}
		if event,ok := data["EVENT"]; ok{
			client,err:= GetClientById(s.Request.RemoteAddr);

			if err != nil {
				log.Fatal("%v",err);
			}

			if(event == "CREATE-ROOM"){

				handleCreateEvent(s,client);

			}

			if(event == "JOIN-ROOM"){

				if room_id,ok := data["ROOM_ID"]; ok{
					fmt.Println(room_id);
					room,err := addToRoom(room_id.(string),s) ;
					updateClientSymbol(s.Request.RemoteAddr,"x");
					if err != nil{
						s.Write([]byte("NO ROOM FOUND"));
						return;
					}
					fmt.Printf("2nd User connected to the room");
					fmt.Println(room.id);
					fmt.Println("The two players now are");
					fmt.Println(room.player1.Id);
					fmt.Println(room.player2.Id);
					data := map[string]string{
						"MESSAGE":"PLAYER 2 CONNECTED",
					};
					JSONData,_ := json.Marshal(data);
					sendMessageToClient(room.player1.Id,JSONData);
					data = map[string]string{
						"MESSAGE":"ROOM CONNECTED",
					};
					JSONData,_ = json.Marshal(data);
					sendMessageToClient(room.player2.Id,JSONData);
				}
			}
		
			if event == "GAME"{
				fmt.Println("event recieved");

				if move,ok := data["MOVE"].(string); ok{ //get the position
					fmt.Println("Move recieved");

					pos,err := strconv.Atoi(move);
					fmt.Println(pos);
					if err != nil {
						fmt.Println("NOT A VALID MOVE")
						return
					}	

					if pos < 1 || pos > 9{
						fmt.Println("NOT A VALID MOVE")
						return;
					}
					

					room,err := GetRoom(s.Request.RemoteAddr);

					if err != nil{
						fmt.Print("error is not nil");
						return;
					}

					if room == nil{
						fmt.Println("Room is nil");
						return;
					}

					if room.gameBoard[pos-1] != ""{
						fmt.Println("NOT A VALID MOVE")
						return;
					}

					if room.gameover {
						fmt.Println("THE GAME IS OVER")
						return;
					}

					client,err = GetClientById(s.Request.RemoteAddr);
					if err != nil{
						fmt.Println("client error is not nil");
						return;
					}
					if client == nil{
						fmt.Println("Client is nil");
						return;
					}

					if client.symbol != room.current {
						fmt.Println("Not your turn");
						return;
					}

					//update the board

					fmt.Println("Client symbol is");
					fmt.Println(client.symbol);

					printAllClients();
					room.gameBoard[pos-1] = client.symbol;
					fmt.Println(room.gameBoard)
					
					if room.current == "x"{
						room.current = "o";
					}else{
						room.current = "x"
					}

					data := map[string]string{
						"MESSAGE":"GAMEMOVE",
						"POSITION":move,
						"SYMBOL":client.symbol,
					}

					JSONData,_ := json.Marshal(data);
					sendMessageToClient(room.player1.Id , JSONData);
					sendMessageToClient(room.player2.Id , JSONData);
					//check for win 
					win,winner := checkForWin(room.gameBoard);
					if( win ){
						fmt.Println("Winner spotted the Winner is");
						fmt.Println(winner);
						
						room.gameover = true;

						data = map[string]string{
							"MESSAGE":"GAMEOVER",
							"RESULT":"WIN",
							"UWIN":"",
						}

						player1,_ := GetClientById(room.player1.Id);

						WinnerData := map[string]string{
							"MESSAGE":"GAMEOVER",
							"RESULT":"WIN",
						}

						LoserData := map[string]string{
							"MESSAGE":"GAMEOVER",
							"RESULT":"LOSE",
						}

						WJSONData,_ := json.Marshal(WinnerData);
						LJSONData,_ := json.Marshal(LoserData);
						if winner == player1.symbol {
							sendMessageToClient(room.player1.Id,WJSONData);
							sendMessageToClient(room.player2.Id,LJSONData);
						}else{
							sendMessageToClient(room.player2.Id,WJSONData);
							sendMessageToClient(room.player1.Id,LJSONData);
						}
						return;
					}

					//check for draw
					if ( checkForDraw(room.gameBoard) ){
						fmt.Printf("Draw occured in game with id %s ",room.id);
						data = map[string]string{
							"MESSAGE":"GAMEOVER",
							"RESULT":"DRAW",
						}

						JSONData,_ := json.Marshal(data);
						sendMessageToClient(room.player1.Id , JSONData);
						sendMessageToClient(room.player2.Id , JSONData);

					}

				}

			}

		}
	})
	port := os.Getenv("PORT");
	if port == ""{
		port = "5000"
	}
	err := http.ListenAndServe(":"+port,nil);
	if err != nil{
		log.Fatal(err);
	}
	fmt.Printf("Server running on %s",port);
}
func checkForWin(gameBoard [9]string) (bool,string) {
	winningLines := [][]int{
        {0, 1, 2}, // Row 1
        {3, 4, 5}, // Row 2
        {6, 7, 8}, // Row 3
        {0, 3, 6}, // Column 1
        {1, 4, 7}, // Column 2
        {2, 5, 8}, // Column 3
        {0, 4, 8}, // Main diagonal
        {2, 4, 6}, // Anti-diagonal
    };
		for _, line := range winningLines {
			if gameBoard[line[0]] != "" &&
			 gameBoard[line[0]] == gameBoard[line[1]] &&
			 gameBoard[line[1]] == gameBoard[line[2]] {
				return true, gameBoard[line[0]];
			}
    }
		return false,"";
}

func checkForDraw(gameboard [9]string) (bool) {
	for _,cell := range gameboard{
		if cell == ""{
			return false;
		}
	}
	return true;
}

func printAllClients(){

	for i,_ := range Users{
		fmt.Println(Users[i]);
	}

}
func GetRoom(id string) (*Room,error){
	for i,r := range Rooms{
		if r.player1.Id == id || r.player2.Id == id{
			return &Rooms[i],nil;
		}
	}
	return nil,errors.New("ROOM NOT FOUND !");
}

func sendMessageToClient(remoteAddr string, message []byte) {
	if session, ok := sessions[remoteAddr]; ok {
		session.Write(message)
	} else {
		fmt.Printf("Session not found for %s\n", remoteAddr)
	}
}

func addToRoom(id string , s *melody.Session) (*Room,error){
	client, err := GetClientById(s.Request.RemoteAddr)
    if err != nil {
        return nil, err // Return an error if the client was not found
  }
	for i,room := range Rooms{
		if room.id == id && room.player2 == nil{

			room.player2 = client;
			Rooms[i] = room;

			return &Rooms[i],nil;

		}
	}
	return nil,errors.New("ROOM NOT FOUND");
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
