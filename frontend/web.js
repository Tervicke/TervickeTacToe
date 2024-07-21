let socket;
let reconnectAttempts = 0;
const maxReconnectAttempts = 4;
const reconnectInterval = 3000;

function connectWebSocket(){
	socket = new WebSocket("wss://tervicketactoe.onrender.com/ws");
	//socket = new WebSocket("ws://localhost:5000/ws")
	
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

