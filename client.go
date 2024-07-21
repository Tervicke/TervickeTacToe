package main

//find the client by id
func GetClientById(Id string) (*Client,error) {
	for _,user := range Users {
		if user.Id == Id{
			return &user,nil;
		}
	}
	return nil,ErrClientNotFound
}

//update the client symbol or set it 
func updateClientSymbol(id string , symbol string){
	for i,user := range Users{
		if user.Id == id{
			user.symbol = symbol
			Users[i] = user;
		}		
	}
}
