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
		}
		if( data["RESULT"] == "WIN" ){
			Display("**U WIN**");
		}
		if( data["RESULT"] == "LOSE" ){
			Display("**U LOSE**");
		}
		showReplayButton();
		showGameText();
	}

	if(data.EVENT == "EXIT"){
		showHomeScreen();
		hideGameScreen();
		document.getElementById("NotifContainer").style.display = "block";
		setTimeout(() => {
		document.getElementById("NotifContainer").style.display = "None";
		}, 1500); 

	}
	if(data.EVENT == "RESTART"){
		resetGame();
	}

}
