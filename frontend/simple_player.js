// This class is a simple implementation of a player to test some things quickly

local_players = {}


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
        this.tick = 0;
    }
}


Player.prototype.playerConnect = function() {
    if(this.connection != null){
        toLog("Client1 already connected!");
    }
    toLog("Connecting as: " + this.name);
    this.connection = new WebSocket("ws://localhost:8080/ws");
    this.connection.onopen = function(evt) {
        console.log("OPEN");
        toLog("Websocket Open");
    };
    this.connection.onclose = function(evt) {
        console.log("CLOSE");
        toLog("Websocket Closed");
        this.connection = null;
    };
    this.connection.onmessage = function(evt) {
        parseServerMessage(evt);
    };
    this.connection.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    };
};

function add_player_gui(new_player) {
    p_name = new_player.name;

    var log = document.getElementById("debug");
    var item = document.createElement("div");

    var content = `
        ${p_name}   
        <div>
            <button onclick="pmove('${p_name}','left')"> ← </button>
            <button onclick="pmove('${p_name}','up')"> ↑ </button>
            <button onclick="pmove('${p_name}','right')"> → </button>
            <button onclick="pmove('${p_name}','down')"> ↓ </button>
            <button> Disconnect </button>
            <button> Invalid </button>
        </div>
    `;

    item.innerHTML = content;

    log.appendChild(item);
}

function pmove(p_name, command) {
    player = local_players[p_name];
    console.log(p_name,command);
}


function makeid() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  
    for (var i = 0; i < 5; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));
  
    return text;
  }