// This class is a simple implementation of a player to test some things quickly

let local_players = {};


function create_player() {
    new_player = new Player();
    local_players[new_player.name] = new_player;
    add_player_gui(new_player);

    new_player.playerConnect();
}


class Player{
    constructor() {
        this.name = makeid();
        this.connection = null;
    }
}


Player.prototype.playerConnect = function() {
    if(this.connection != null){
        toLog("Client1 already connected!");
    }
    toLog("Connecting as: " + this.name);
    player = this;
    this.connection = new WebSocket("ws://localhost:8080/ws");
    // noinspection JSUnusedLocalSymbols
    this.connection.onopen = function(evt) {
        console.log("OPEN");
        toLog("Websocket Open");
        player.updateUsername()
    };
    // noinspection JSUnusedLocalSymbols
    this.connection.onclose = function(evt) {
        toLog("Websocket Closed");
        player.connection = null;
    };
    this.connection.onmessage = function(evt) {
        toLog(evt.data);
        console.log(evt);
    };
    this.connection.onerror = function(evt) {
        toLog("ERROR: " + evt.data);
    };
};

Player.prototype.disconnect = function() {
    if(this.connection != null){
        this.connection.close(1000, "Goodbye!");
    }
};

Player.prototype.updateUsername = function() {
    toLog("Sending username: " + this.name);
    let payload = {
        type: "username",
        command: this.name,
    };
    this.connection.send(JSON.stringify(payload));
};

function add_player_gui(new_player) {
    p_name = new_player.name;

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
    toLog("Player ("+  player.name + "): " + command);
    let payload = {
        type:  type,
        value: command,
    };
    player.connection.send(JSON.stringify(payload));
}

function makeid() {
    let text = "";
    let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  
    for (let i = 0; i < 5; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));
  
    return text;
  }