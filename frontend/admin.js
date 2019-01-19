let admin_ws = null;

function adminConnect() {
    if(admin_ws != null){
        toLog("Admin already connected..");
        return false;
    }
    toLog("Starting admin connect");
    admin_ws = new WebSocket("ws://localhost:8080/admin");
    // noinspection JSUnusedLocalSymbols
    admin_ws.onopen = function(evt) {
        toLog("Admin socket started");
    };
    // noinspection JSUnusedLocalSymbols
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


function startSystem() {
    if(admin_ws == null) {
        return false;
    }
    toLog("Admin: Starting system");
    sendCommand("start");
}

function pauseSystem() {
    if(admin_ws == null) {
        return false;
    }
    toLog("Admin: pauseSystem");
    sendCommand("pause");
}
function restartSystem() {
    if(admin_ws == null) {
        return false;
    }
    toLog("Admin: restartSystem");
    sendCommand("restart");
}

function sendInvalid() {
    toLog("Sending invalid data");

    let payload = {
        invalid: "test"
    };
    admin_ws.send(JSON.stringify(payload))
}

function sendCommand(command) {
    toLog("Sending admin command: " + command);
    let payload = {
        command: command,
    };
    admin_ws.send(JSON.stringify(payload));
}


function create_simple_player(){
    create_player()
}
