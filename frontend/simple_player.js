// This class is a simple implementation of a player to test some things quickly

let local_players = {};


function create_player() {
    let new_player = new Player();
    local_players[new_player.name] = new_player;
    add_player_gui(new_player);

    new_player.playerConnect();
}


class Player{
    constructor() {
        this.name = makeid();
        this.color = getRandomRgb();
        this.connection = null;
    }
}


Player.prototype.playerConnect = function() {
    if(this.connection != null){
        console.log("Client1 already connected!");
    }
    console.log("Connecting as: " + this.name);
    let player = this;
    this.connection = new WebSocket("ws://localhost:8080/ws");
    // noinspection JSUnusedLocalSymbols
    this.connection.onopen = function(evt) {
        console.log("OPEN");
        console.log("Websocket Open");
        player.updateUsername();
        player.updateColor();
    };
    // noinspection JSUnusedLocalSymbols
    this.connection.onclose = function(evt) {
        console.log("Websocket Closed");
        player.connection = null;
    };
    this.connection.onmessage = function(evt) {
        player.parseEvent(evt.data);
    };
    this.connection.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    };
};

Player.prototype.disconnect = function() {
    if(this.connection != null){
        this.connection.close(1000, "Goodbye!");
    }
};

Player.prototype.updateUsername = function() {
    console.log("Sending username: " + this.name);
    let payload = {
        type: "username",
        value: this.name,
    };
    this.connection.send(JSON.stringify(payload));
};

Player.prototype.updateColor = function() {
    console.log("Sending color: " + this.color);

    let payload = {
        type: "color",
        value: this.color,
    };
    this.connection.send(JSON.stringify(payload));
};

Player.prototype.move = function(command) {
    console.log("Moving:" + command);
    this.sendCommand("move", command);
};

Player.prototype.sendCommand = function(type, command) {
    let payload = {
        type:  type,
        value: command,
    };
    this.connection.send(JSON.stringify(payload));
};

Player.prototype.moveToFood = function(){
    let head = this.data.Self.Head;
    let food = this.data.GameStatus.GameMap.Foods;

    let closest;
    let closest_dist = Number.MAX_SAFE_INTEGER;

    for (let i = 0; i < food.length; ++i) {
        let f = food[i];
        let x_len = f["X"] - head["X"];
        let y_len = f["Y"] - head["Y"];
        let len = Math.abs(x_len) + Math.abs(y_len);

        if (len < closest_dist) {
            closest_dist = len;
            closest = f;
        }
    }

    let x_offset = closest.X - head.X;
    let y_offset = closest.Y - head.Y;

    if(y_offset > 0) {
        this.move("up")
    } else if(y_offset < 0) {
        this.move("down")
    }else if(x_offset > 0) {
        this.move("right")
    } else if(x_offset < 0) {
        this.move("left")
    } else {
        console.log("No moves left");
        this.move("left")
    }
};

Player.prototype.findClosestTile = function(type) {

};

Player.prototype.parseEvent = function(evt_data) {
    this.data = JSON.parse(evt_data);
    this.moveToFood(this.data)
};

function getRandomRgb() {
    let num = Math.round(0xffffff * Math.random());
    let r = num >> 16;
    let g = num >> 8 & 255;
    let b = num & 255;
    return 'rgb(' + r + ', ' + g + ', ' + b + ')';
}


function add_player_gui(new_player) {
    let p_name = new_player.name;

    let log = document.getElementById("debug");
    let item = document.createElement("div");
    item.setAttribute("id", p_name);

    item.innerHTML = `
        ${p_name}   
        <div>
            <button onclick="pmove('${p_name}','left')"> ← </button>
            <button onclick="pmove('${p_name}','up')"> ↑ </button>
            <button onclick="pmove('${p_name}','right')"> → </button>
            <button onclick="pmove('${p_name}','down')"> ↓ </button>
            <button onclick="disconnect_player('${p_name}')"> Disconnect </button>
            <button onclick="send_invalid('${p_name}')"> Send invalid </button>
        </div>
    `;

    log.appendChild(item);
}

function disconnect_player(p_name) {
    let player = local_players[p_name];
    player.disconnect();

    let element = document.getElementById(p_name);
    element.parentNode.removeChild(element);

}


function pmove(p_name, command) {
    let player = local_players[p_name];
    console.log(p_name,command);
    sendCommand(player,"move", command);
}

function send_invalid(playerID) {
    let player = local_players[playerID];
    sendCommand(player, makeid(), "invalid");
}

function sendCommand(player,type, command) {
    console.log("Player ("+  player.name + "): " + command);
    let payload = {
        type:  type,
        value: command,
    };
    player.connection.send(JSON.stringify(payload));
}

function create_simple_player(){
    create_player();
    console.log("Players: "  +Object.keys(local_players).length)
}

function makeid() {
    let text = "";
    // noinspection SpellCheckingInspection
    let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  
    for (let i = 0; i < 5; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));
  
    return text;
  }