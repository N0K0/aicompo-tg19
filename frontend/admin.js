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

        case "players":
            update_players(json.message);
            break;
        default:
            console.log("Unable to parse admin message");
            break;
    }
};

AdminConnection.prototype.startSystem = function () {
    console.log("Admin: Starting system");

    showScreen(0);
    let payload = {
        type: "start",
        message: ""
    };

    this.conn.send(JSON.stringify(payload));

};

AdminConnection.prototype.pauseSystem = function() {
    console.log("Admin: pauseSystem");

    let payload = {
        type: "pause",
        message: ""
    };

    this.conn.send(JSON.stringify(payload));
};

AdminConnection.prototype.restartSystem = function() {
    console.log("Admin: restartSystem");

    let payload = {
        type: "restart",
        message: ""
    };

    this.conn.send(JSON.stringify(payload));

};


// Function for testing via sending garbage
AdminConnection.prototype.sendInvalid = function() {
    console.log("Sending invalid data");

    let payload = {
        invalid: "test"
    };
    this.conn.send(JSON.stringify(payload))
};

AdminConnection.prototype.sendCommand = function(command) {
    console.log("Sending admin command: " + command);
    let payload = {
        command: command,
    };
    this.conn.send(JSON.stringify(payload));
};


AdminConnection.prototype.kickPlayer = function(player) {
    console.log("Kicking player: " + player);

    let payload = {
        type: "kick",
        message: player
    };

    this.conn.send(JSON.stringify(payload))

};


AdminConnection.prototype.create_simple_player = function(){
    create_player();
    console.log("Players: "  +Object.keys(local_players).length)
};
