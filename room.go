package main
import (
	"math/rand"

	"github.com/olahol/melody"
)
func createRoom(player1 *Client) (*Room){
	room := Room{
		id:generateRoomId(),
		player1: player1,
		player2: nil,
		current: getRandomSymbol(),
	}
	Rooms = append(Rooms, room)
	return &room;
}

func generateRoomId() string {
	const charset = "0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func getRandomSymbol() string {
	const charset = "xo"
	b := make([]byte, 1)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

//check for valid room id
func checkRoomId(id string) (bool) {
	for _,room := range Rooms{
		if room.id == id{
			return true;
		}
	}
	return false;
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
	return nil,ErrRoomNotFound
}

func GetRoomByClientId(id string) (*Room,error){
	for i,r := range Rooms{
		if r.player1.Id == id || r.player2.Id == id{
			return &Rooms[i],nil;
		}
	}
	return nil,ErrRoomNotFound;
}
