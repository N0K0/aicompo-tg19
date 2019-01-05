c1ws = null;

function clientConnect() {
    if(c1ws != null){
        toLog("Client1 already connected");
    }
    toLog("Trying to connect to the backend as client");
    c1ws = new WebSocket("ws://localhost:8080/ws");
    c1ws.onopen = function(evt) {
        console.log("OPEN");
        toLog("Websocket Open");
    };
    c1ws.onclose = function(evt) {
        console.log("CLOSE");
        toLog("Websocket Closed");
        c1ws = null;
    };
    c1ws.onmessage = function(evt) {
        parseServerMessage(evt);
    };
    c1ws.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    };
}

function clientDisconnect() {
    if(c1ws != null) {
        toLog("Closing client");
        c1ws.close();
        c1ws = null;
    } else {
        toLog("Socket not open");
    }
}

// Parse messages from the server
function parseServerMessage(evt){
    toLog("New message from server!");
    console.log(evt);
}
