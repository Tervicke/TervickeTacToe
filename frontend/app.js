connectWebSocket();
showHomeScreen()

function resetGame(){
	var buttons = document.getElementById("ButtonContainer").children;
	Array.from(buttons).forEach(function(button) {
    button.innerHTML = "";
		button.classList.remove("NeonRed");
		button.classList.remove("NeonBlue");
	});
	removeDisplay();
	removeReplayButton();
}

function showReplayButton(){
	var button = document.getElementById("ReplayButton")
	button.style.display = "block";
}
function removeRoomDetails(){
	var container = document.getElementById("room_details")
	while(container.firstChild){
		container.removeChild(container.firstChild);
	}
}
function Display(textstring){
	var text = document.getElementById("result");
	text.innerHTML = textstring;
}
function showGameText(){
	document.getElementById("GameText").style.display = "flex";
}
function removeDisplay(){
	document.getElementById("GameText").style.display = "None";
}
function removeReplayButton(){
	document.getElementById("ReplayButton").style.display = "None";
}

function hideHomeScreen(){
	ButtonContainer = document.getElementById("GameMenu");
	ButtonContainer.style.display = "None";
}

function showHomeScreen(){
	document.getElementById("GameMenu").display = "flex";
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


function showGameScreen(){
	const ButtonContainer = document.getElementById('ButtonContainer')
	for(var i = 1 ; i <= 9 ; i++){
		var GameButton = document.createElement('button')
		GameButton.classList.add('GameButton')
		GameButton.id=i;
		GameButton.addEventListener('click', ButtonClick)
		ButtonContainer.appendChild(GameButton);
	}
}

function hideGameScreen(){
	document.getElementById("GameText").remove();
	document.getElementById("GameContainer").remove();
}
function ButtonClick(event){
	if(event.target.innerHTML == ''){
		data = {
			"EVENT":"GAMEMOVE",
			"MOVE":event.target.id
		}
		socket.send(JSON.stringify(data));
	}
}

function HandleReplay(){
		data = {
			"EVENT":"REPLAY",
		}
		socket.send(JSON.stringify(data));
}
