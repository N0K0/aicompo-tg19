class AdminConnection {
    constructor() {
        this.conn = null;
    }
}

AdminConnection.prototype.adminConnect = function() {
    console.log("Starting admin connect");

    if(this.conn != null) {
        console.log("Socket is already connected!");
        return
    }

    let admin = this;

    this.conn = new WebSocket("ws://localhost:8080/admin");
    // noinspection JSUnusedLocalSymbols
    this.conn.onopen = function(evt) {
        console.log("Admin socket started");
    };
    // noinspection JSUnusedLocalSymbols
    this.conn.onclose = function(evt) {
        console.log("Admin socket closed!");
        admin.adminDisconnect()
    };
    this.conn.onmessage = function(evt) {
        console.log("A RESPONSE: " + evt.data);
        admin.parseAdminEvent(evt);
    };
    this.conn.onerror = function(evt) {
        console.log("A ERROR: " + evt.data);
    };
};

AdminConnection.prototype.adminDisconnect = function() {
    if(this.conn != null) {
        console.log("Closing socket");
        this.conn.close();
        this.conn = null;
    } else {
        console.log("Socket not open");
    }
};


AdminConnection.prototype.parseAdminEvent = function(evt) {
    let json = JSON.parse(evt.data);

    switch (json.type) {
        case "config_push":
            import_settings(json.message);
            break;
        default:
            console.log("Unable to parse admin message");
            break;
    }
};

AdminConnection.prototype.startSystem = function() {
    if(admin_ws == null) {
        console.log("No socket found... rejecting this call and trying to reconnect");
        this.adminConnect();
        return false;
    }
    console.log("Admin: Starting system");
    this.sendCommand("start");
};

AdminConnection.prototype.pauseSystem = function() {
    if(admin_ws == null) {
        console.log("No socket found... rejecting this call and trying to reconnect");
        this.adminConnect();
        return false;
    }
    console.log("Admin: pauseSystem");
    this.sendCommand("pause");
};

AdminConnection.prototype.restartSystem = function() {
    if(admin_ws == null) {
        console.log("No socket found... rejecting this call and trying to reconnect");
        this.adminConnect();
        return false;
    }
    console.log("Admin: restartSystem");
    this.sendCommand("restart");
};


// Function for testing via sending garbage
AdminConnection.prototype.sendInvalid = function() {
    console.log("Sending invalid data");

    let payload = {
        invalid: "test"
    };
    admin_ws.send(JSON.stringify(payload))
};

AdminConnection.prototype.sendCommand = function(command) {
    console.log("Sending admin command: " + command);
    let payload = {
        command: command,
    };
    admin_ws.send(JSON.stringify(payload));
};


AdminConnection.prototype.create_simple_player = function(){
    create_player();
    console.log("Players: "  +Object.keys(local_players).length)
};
