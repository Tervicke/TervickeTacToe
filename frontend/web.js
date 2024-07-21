let socket;

function connectWebSocket(){
	socket = new WebSocket("wss://tervicketactoe.onrender.com/ws");
	//socket = new WebSocket("ws://localhost:5000/ws")
	
	socket.onopen = function(){
		console.log("Socket is now open");
	}

	socket.onopen = function(){
		socket.addEventListener('message', (event) => {
			HandleMessages(event)
		})
	}

	socket.onclose = function(){
		console.log("socket closed")
	}

}


