//const socket = new WebSocket("ws://localhost:5000/ws");
let socket;
let reconnectAttempts = 0;
const maxReconnectAttempts = 10;
const reconnectInterval = 3000;

function connectWebSocket(){
	//socket = new WebSocket("wss://tervicketactoe.onrender.com/ws");
	socket = new WebSocket("ws://localhost:5000/ws")
	
	socket.onopen = function(){
		console.log("Socket is now open");
		reconnectAttempts = 0;
	}

	socket.onopen = function(){
		socket.addEventListener('message', (event) => {
			HandleMessages(event)
		})
	}

	socket.onclose = function(){
		console.log("socket was being closed restarted")
		attemptReconnect();
	}

}
function attemptReconnect(){
	if (reconnectAttempts < maxReconnectAttempts) {
        setTimeout(() => {
            reconnectAttempts++;
            console.log(`Reconnecting attempt ${reconnectAttempts}`);
            connectWebSocket();
        }, reconnectInterval);
    } else {
        console.log("Max reconnect attempts reached. Please check the server.");
  }
}

connectWebSocket();

function HandleMessages(event){
	data = JSON.parse(event.data);
	if(data.EVENT == "ROOM CONNECTED"){
		room_id_text = document.createElement("h3");
		if("ROOM_ID" in data ){
			room_id_text.innerHTML = "Room id is " + data.ROOM_ID; 
			waiting_for_player = document.createElement("h3");
			waiting_for_player.innerHTML = "waiting for the other player to join";
			room_details = document.createElement('div');
			room_details.id = "room_details";
			room_details.appendChild(room_id_text);
			room_details.appendChild(waiting_for_player);
			document.body.appendChild(room_details);
		}else{
			showGameScreen()
		}
		hideHomeScreen();
	}

	if(data.EVENT == "PLAYER 2 CONNECTED"){
		removeRoomDetails();
		showGameScreen();
	}

	if(data.EVENT == "GAMEMOVE"){
		pos = data["POSITION"]
		symbol = data["SYMBOL"]
		var button = document.getElementById(pos);
		if(data["SYMBOL"] == "x"){
			button.classList.add('NeonRed');
		}else{
			button.classList.add('NeonBlue');
		}
		button.innerHTML = symbol;
	}

	if(data.EVENT == "GAMEOVER"){
		if(data["RESULT"] == "DRAW"){
			Display("Draw !!")
		}else{
			if( data["RESULT"] == "WIN" ){
				Display("**U WIN**");
			}else{
				Display("**U LOSE**");
			}
		}
	}

}

function removeRoomDetails(){
	var container = document.getElementById("room_details")
	while(container.firstChild){
		container.removeChild(container.firstChild);
	}
}
function Display(textstring){
	var container = document.getElementById("GameText");
	var text = document.createElement("h3");
	text.innerHTML = textstring
	container.appendChild(text);
}

function hideHomeScreen(){
	ButtonContainer = document.getElementById("GameMenu");
	ButtonContainer.remove();
}

function showHomeScreen(){
	const GameMenu = document.createElement('div')
	GameMenu.id = 'GameMenu'
	const JoinButton = document.createElement('button')
	JoinButton.classList.add("HomeButton")
	JoinButton.addEventListener('click', JoinEvent)
	JoinButton.textContent = "Join";
	const CreateButton = document.createElement('button')
	CreateButton.textContent = "Create"
	CreateButton.addEventListener('click', CreateEvent)
	CreateButton.classList.add("HomeButton")
	GameMenu.append(CreateButton)
	GameMenu.append(JoinButton)
	document.body.append(GameMenu)
}
function JoinEvent(){
	let id = prompt("ENTER THE Room id")
	data = {
		"EVENT":"JOIN-ROOM",
		"ROOM_ID":id,
	}	
	socket.send(JSON.stringify(data));
}

function CreateEvent() {
	data = {
		"EVENT":"CREATE-ROOM",
		"CONTENT":null
	}	
	socket.send(JSON.stringify(data));
}

showHomeScreen()

function showGameScreen(){
	const ButtonContainer = document.createElement('div')
	ButtonContainer.id='ButtonContainer'
	for(var i = 1 ; i <= 9 ; i++){
		var GameButton = document.createElement('button')
		GameButton.classList.add('GameButton')
		GameButton.id=i;
		GameButton.addEventListener('click', ButtonClick)
		ButtonContainer.appendChild(GameButton);
	}
	GameContainer = document.createElement('div');
	GameContainer.appendChild(ButtonContainer)
	GameContainer.id="GameContainer";
	GameText = document.createElement('div');
	GameText.id = "GameText";
	document.body.appendChild(GameContainer);
	document.body.appendChild(GameText);
}

function ButtonClick(event){
	if(event.target.innerHTML == ''){
		data = {
			"EVENT":"GAMEMOVE",
			"MOVE":event.target.id
		}
		socket.send(JSON.stringify(data));
		console.log("data sent");
	}
}
