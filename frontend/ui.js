
let canvas = document.getElementById("snake");
let ctx = canvas.getContext("2d");

// Screens
let screen_snake = document.getElementById("snake");
let screen_menu = document.getElementById("menu");
let screen_gameover = document.getElementById("gameover");
let screen_setting = document.getElementById("setting");
let screen_lobby = document.getElementById("lobby");

// Main menu
let button_newgame_menu = document.getElementById("newgame_menu");
button_newgame_menu.onclick = function(){newGame();};

// Main menu
let button_setting_menu = document.getElementById("setting_menu");
button_setting_menu.onclick = function(){showScreen(2);};

// Game over
let button_newgame_gameover = document.getElementById("newgame_gameover");
button_newgame_gameover.onclick = function(){newGame();};

let button_setting_gameover = document.getElementById("setting_gameover");
button_setting_gameover.onclick = function(){showScreen(2)};

// Settings
let button_save_setting = document.getElementById("setting_save");
button_save_setting.onclick = function () {save_settings()};

let button_main_menu = document.getElementById("main_menu");
button_main_menu.onclick = function () {showScreen(1)};

let button_start_game = document.getElementById("start_game");
button_start_game.onclick = function () {admin_ws.startSystem()};




// etc
let ele_score = document.getElementById("score_value");

// Settings levels
let wall_setting = document.getElementsByName("wall");
let map_size_setting = document.getElementsByName("map_size");
let min_time_setting = document.getElementsByName("min_time");
let max_time_setting = document.getElementsByName("max_time");

let map_size_str;
let time_min_turn;
let time_max_turn;
let wall;

// --------------------

setMapSize("0x0");
setMinTime("200");
setMaxTime("800");
setWall("1");

showScreen(4);

// --------------------
// Settings


// wall
for(let i = 0; i < wall_setting.length; i++){
    wall_setting[i].addEventListener("click", function(){
        for(let i = 0; i < wall_setting.length; i++){
            if(wall_setting[i].checked){
                setWall(wall_setting[i].value);
            }
        }
    });
}

// map
for(let i = 0; i < map_size_setting.length; i++){
    map_size_setting[i].addEventListener("click", function(){
        for(let i = 0; i < map_size_setting.length; i++){
            if(map_size_setting[i].checked){
                setMapSize(map_size_setting[i].value);
            }
        }
    });
}

// min time
for(let i = 0; i < min_time_setting.length; i++){
    min_time_setting[i].addEventListener("click", function(){
        for(let i = 0; i < min_time_setting.length; i++){
            if(min_time_setting[i].checked){
                setMinTime(min_time_setting[i].value);
            }
        }
    });
}

// max time
for(let i = 0; i < max_time_setting.length; i++){
    max_time_setting[i].addEventListener("click", function(){
        for(let i = 0; i < max_time_setting.length; i++){
            if(max_time_setting[i].checked){
                setMaxTime(max_time_setting[i].value);
            }
        }
    });
}


function newGame() {
    showScreen(4);
}


// 0 for the game
// 1 for the main menu
// 2 for the settings screen
// 3 for the game over screen

function showScreen(screen_opt){
    switch(screen_opt){

        case 0:
            screen_snake.style.display = "block";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            screen_lobby.style.display = "none";
            break;

        case 1:
            screen_snake.style.display = "none";
            screen_menu.style.display = "block";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            screen_lobby.style.display = "none";
            break;

        case 2:
            fetch_settings();
            screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "block";
            screen_gameover.style.display = "none";
            screen_lobby.style.display = "none";
            break;

        case 3:
            screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "block";
            screen_lobby.style.display = "none";
            break;

        case 4:
            screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            screen_lobby.style.display = "block";
            break;
    }
}

function fetch_settings(){
    let envelope = {
        "type": "config_get",
        "message": ""
    };
    admin_ws.conn.send(JSON.stringify(envelope));
}

// Settings logic
function save_settings(){

    let payload = {
        "configs": [
            {"name": "minTurnUpdate",
             "value": time_min_turn},
            {"name": "maxTurnUpdate",
                "value": time_max_turn},
            {"name": "mapSize",
                "value": map_size_str},
            {"name": "outerWalls",
                "value": wall},
        ]
    };

    let envelope = {
        "type": "config",
        "message": payload
    };

    console.log(envelope);

    admin_ws.conn.send(JSON.stringify(envelope));
    showScreen(1);
}

function import_settings(payload){
    console.log("Importing settings");
    let settings = JSON.parse(payload);
    console.log(settings);
    for(let k in settings){
        console.log(k);

        let value = settings[k];

        switch(k){
            case "minTurnUpdate":
                check_setting(min_time_setting, value.toString());
                setMinTime(value.toString());
                break;
            case "maxTurnUpdate":
                check_setting(max_time_setting, value.toString());
                setMaxTime(value.toString());
                break;
            case "outerWalls":
                check_setting(wall_setting, value.toString());
                setWall(value.toString());
                break;
            case "mapSize":
                check_setting(map_size_setting, value);
                setMapSize(value);
                break;
        }
    }
}

function update_players(message) {
    console.log("Updating the players state");
    let players = JSON.parse(message);

    update_player_ui(players);

}

function update_player_ui(players_json){
    // Lobby
    let lobby_player_div = document.getElementById("players");
    let lobby_no_player_div = document.getElementById("no_players");
    let lobby_player_number = document.getElementById("num_players");

    let cloned_player_div = lobby_player_div.cloneNode(false);


    if(Object.getOwnPropertyNames(players_json).length  === 0 ){
        lobby_no_player_div.style.display = "block";
        lobby_player_div.style.display = "none";
        return
    }

    let label_num_player = document.createElement("label");
    label_num_player.id = "num_players";
    label_num_player.innerText = "Players: " + Object.keys(players_json).length.toString();

    cloned_player_div.appendChild(label_num_player);

    for( let p in players_json) {
        let v = players_json[p];


        let tmp_div_player = document.createElement("div");
        tmp_div_player.id = "player";

        let tmp_player_name = document.createElement("label");
        tmp_player_name.innerText = v.username;
        tmp_player_name.id = "player_name";

        let tmp_player_kick = document.createElement("label");
        tmp_player_kick.innerText = "[X]";
        tmp_player_kick.id = "player_kick";


        tmp_div_player.appendChild(tmp_player_name);
        tmp_div_player.appendChild(tmp_player_kick);

        tmp_player_kick.onclick = function() {
            admin_ws.kickPlayer( v.username)
        };

        cloned_player_div.appendChild(tmp_div_player);
    }

    lobby_player_div.replaceWith(cloned_player_div);

    cloned_player_div.style.display = "block";
    lobby_no_player_div.style.display = "none";
    lobby_player_div.style.display = "block";


}

function check_setting(elements, value) {

    for (let i = 0; i < elements.length; ++i) {
        if (elements[i].value === value) {
            elements[i].checked = true;

        } else {
            elements[i].checked = false;

        }
    }

}
/////////////////////////////////////////////////////////////

function setMapSize(value){
    console.log("Setting value to " + value + " old value: " + map_size_str);
    map_size_str = value;
}

/////////////////////////////////////////////////////////////

function setMinTime(value){
    console.log("Setting value to " + value + " old value: " + time_min_turn);
    time_min_turn = value;
}

/////////////////////////////////////////////////////////////

function setMaxTime(value){
    console.log("Setting value to " + value + " old value: " + time_max_turn);

    time_max_turn = value;
}



/////////////////////////////////////////////////////////////
function setWall(value){
    console.log("Setting value to " + value + " old value: " + wall);

    wall = value;
    if(wall === 0){screen_snake.style.borderColor = "#606060";}
    if(wall === 1){screen_snake.style.borderColor = "#FFFFFF";}
}