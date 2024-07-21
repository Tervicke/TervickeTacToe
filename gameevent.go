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
	if(room.gameover){
		return;
	}
	if ( checkValidMove(pos,room)){
		if( playersMove(room,client) ){
			executeMove(move,pos,room,client);
		}
	}

	win,winner := checkForWin(room.gameBoard);
	if( win ){
		declareWinner(room , winner);
	}

	if (checkForDraw(room.gameBoard) ){
		data := map[string]string{
			"RESULT":"DRAW",
		}
		sendJSONMessage(room.player1.Id,"GAMEOVER",data);
		sendJSONMessage(room.player2.Id,"GAMEOVER",data);
	}

}

func declareWinner(room *Room , winner string){
	room.gameover = true;
	Wdata := map[string]string{
		"RESULT":"WIN",
	}
	Ldata := map[string]string{
		"RESULT":"LOSE",
	}
	
	if room.player1.symbol == winner{
		sendJSONMessage(room.player1.Id , "GAMEOVER" , Wdata);
		sendJSONMessage(room.player2.Id , "GAMEOVER" , Ldata);
	}else{
		sendJSONMessage(room.player2.Id , "GAMEOVER" , Wdata);
		sendJSONMessage(room.player1.Id , "GAMEOVER" , Ldata);
	}

}

func checkForDraw(gameboard [9]string) (bool) {
	for _,cell := range gameboard{
		if cell == ""{
			return false;
		}
	}
	return true;
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
func HandleReplay( s *melody.Session ){
	room,err := GetRoomByClientId(s.Request.RemoteAddr);
	if err != nil {
		log.Printf("%v",err);
	}
	if room.player1.Id == s.Request.RemoteAddr{
		room.player1replay = true;
	}

	if room.player2.Id == s.Request.RemoteAddr{
		room.player2replay = true;
	}

	if room.player1replay && room.player2replay {
		restartGame(room);
	}

}
func restartGame(room *Room){
	room.current = getRandomSymbol();

	//clear the array
	clearStringArray(&room.gameBoard);

	room.player1replay = false;
	room.player2replay = false;
	
	room.gameover = false;

	data := map[string]string{};
	
	sendJSONMessage(room.player1.Id,"RESTART",data);
	sendJSONMessage(room.player2.Id,"RESTART",data);
}


func clearStringArray(arr *[9]string) {
    for i := range arr {
        arr[i] = ""
    }
}
