package main

import (
	"log"

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
	sendJSONMessage(s,"ROOM CONNECTED",data)
}


