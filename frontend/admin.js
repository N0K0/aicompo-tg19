
function adminConnect() {

    console.log("Starting admin connect");
    let admin_ws = new WebSocket("ws://localhost:8080/admin");
    // noinspection JSUnusedLocalSymbols
    admin_ws.onopen = function(evt) {
        console.log("Admin socket started");
    };
    // noinspection JSUnusedLocalSymbols
    admin_ws.onclose = function(evt) {
        console.log("Admin socket closed!");
        admin_ws = null;
    };
    admin_ws.onmessage = function(evt) {
        console.log("A RESPONSE: " + evt.data);
        parseAdminEvent(evt)
    };
    admin_ws.onerror = function(evt) {
        console.log("A ERROR: " + evt.data);
    };

    return admin_ws
}

function adminDisconnect() {
    if(admin_ws != null) {
        console.log("Closing socket");
        admin_ws.close();
        admin_ws = null;
    } else {
        console.log("Socket not open");
    }
}

function parseAdminEvent(evt) {
    let json = JSON.parse(evt.data);

    switch (json.type) {
        case "config_push":
            import_settings(json.message);
            break;
        default:
            console.log("Unable to parse admin message");
            break;

    }
}

function startSystem() {
    if(admin_ws == null) {
        return false;
    }
    console.log("Admin: Starting system");
    sendCommand("start");
}

function pauseSystem() {
    if(admin_ws == null) {
        return false;
    }
    console.log("Admin: pauseSystem");
    sendCommand("pause");
}
function restartSystem() {
    if(admin_ws == null) {
        return false;
    }
    console.log("Admin: restartSystem");
    sendCommand("restart");
}

function sendInvalid() {
    console.log("Sending invalid data");

    let payload = {
        invalid: "test"
    };
    admin_ws.send(JSON.stringify(payload))
}

function sendCommand(command) {
    console.log("Sending admin command: " + command);
    let payload = {
        command: command,
    };
    admin_ws.send(JSON.stringify(payload));
}


function create_simple_player(){
    create_player();
    console.log("Players: "  +Object.keys(local_players).length)
}
