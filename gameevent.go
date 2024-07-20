package main

import (
	"log"
	"strconv"

	"github.com/olahol/melody"
)

//CREATE-ROOM
func handleCreateEvent(s *melody.Session ,client *Client){
	updateClientSymbol(client.Id ,"o");
	room := createRoom(client);
	log.Printf("Room created id: %s",room.id);
	data := map[string]string{
		"ROOM_ID":room.id,
	}
	sendJSONMessage(s.Request.RemoteAddr,"ROOM CONNECTED",data)
}

//JOIN-ROOM
func handleJoinEvent(room_id string , s *melody.Session){

	room,err := addToRoom(room_id,s) ;
	updateClientSymbol(s.Request.RemoteAddr,"x");

	if err != nil {
		sendJSONMessage(s.Request.RemoteAddr,"NO ROOM FOUND",map[string]string{})
		log.Printf("%v",err);
		return;
	}

	//send to the player1 that player2 has been connected	
	data := map[string]string{}
	sendJSONMessage(room.player1.Id , "PLAYER 2 CONNECTED" , data);

	//send to the second player that he has joined the room 
	data = map[string]string{}
	sendJSONMessage(room.player2.Id , "ROOM CONNECTED" , data);

}

func handleGameMoveEvent(move string , s *melody.Session) {
	pos , err := strconv.Atoi(move)
	if err != nil {
		log.Printf("%v",err)
	}

	room,err := GetRoomByClientId(s.Request.RemoteAddr);
	client,_:= GetClientById(s.Request.RemoteAddr);
	if err != nil{
		log.Printf("%v",err);
	}
	if ( checkValidMove(pos,room) ){
		if( playersMove(room,client) ){
			executeMove(move,pos,room,client);
		}
	}

}

func playersMove(room *Room, client *Client) (bool){
	if room.current == client.symbol{
		return true;
	}
	return false;
}

func checkValidMove(pos int , room *Room) (bool) {
	if pos < 10 && pos > 0{
		if room.gameBoard[pos-1] == ""{
			return true;
		}
		return false;
	}
	return false
}

func executeMove(move string ,pos int, room *Room , client *Client){
	room.gameBoard[pos-1] = client.symbol;
	toogleSymbol(room);

	data := map[string]string{
		"POSITION":move,
		"SYMBOL":client.symbol,
	}

	sendJSONMessage(room.player1.Id,"GAMEMOVE",data);
	sendJSONMessage(room.player2.Id,"GAMEMOVE",data);
}
func toogleSymbol(room *Room){
	if room.current == "x"{
		room.current = "o";
	}else{
		room.current = "x";
	}
}
