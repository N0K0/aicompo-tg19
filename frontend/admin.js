var admin_ws = null

function adminConnect() {
    if(admin_ws != null){
        toLog("Admin already connected..");
        return false;
    }
    toLog("Starting admin connect");
    admin_ws = new WebSocket("ws://localhost:8080/admin");
    admin_ws.onopen = function(evt) {
        toLog("Admin socket started");
    };
    admin_ws.onclose = function(evt) {
        toLog("Admin socket closed!");
        admin_ws = null;
    };
    admin_ws.onmessage = function(evt) {
        toLog("A RESPONSE: " + evt.data);
    };
    admin_ws.onerror = function(evt) {
        toLog("A ERROR: " + evt.data);
    };
}

function adminDisconnect() {
    if(admin_ws != null) {
        toLog("Closing socket");
        admin_ws.close();
        admin_ws = null;
    } else {
        toLog("Socket not open");
    }
}


var startSystem = function() {
    if(admin_ws == null) {
        return False;
    }
    sendCommand("Start");
}

var sendCommand = function(command) {
    toLog("Sending admin command: " + command);
    payload = {
        command: command,
    }
    admin_ws.send(JSON.stringify(payload));
}