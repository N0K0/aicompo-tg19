
let canvas = document.getElementById("snake");
let ctx = canvas.getContext("2d");

// Screens
let screen_snake = document.getElementById("snake");
let screen_menu = document.getElementById("menu");
let screen_gameover = document.getElementById("gameover");
let screen_setting = document.getElementById("setting");

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

showScreen("menu");

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
            break;

        case 1:
            screen_snake.style.display = "none";
            screen_menu.style.display = "block";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            break;

        case 2:
            fetch_settings();
            screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "block";
            screen_gameover.style.display = "none";
            break;

        case 3:
            screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "block";
            break;
    }
}

function fetch_settings(){
    // Todo make front fetch settings from backend

    let envelope = {
        "type": "config_get",
        "message": ""
    };

    admin_ws.send(JSON.stringify(envelope));
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

    admin_ws.send(JSON.stringify(envelope));
    showScreen(1);
}

function import_settings(payload){
    console.log("Importing settings");
    let settings = JSON.parse(payload);
    console.log(settings);
    for(let k in settings){
        console.log(k);

        value = settings[k];

        switch(k){
            case "minTurnUpdate":
                check_setting(min_time_setting, value.toString());
                break;
            case "maxTurnUpdate":
                check_setting(max_time_setting, value.toString());
                break;
            case "outerWalls":
                check_setting(wall_setting, value.toString());
                break;
            case "mapSize":
                check_setting(map_size_setting, value);
                break;
        }
    }
}

function check_setting(elements, value) {
    console.log("Checking value", value);
    console.log(elements);

    for (let i = 0; i < elements.length; ++i) {
        console.log(elements[i]);
        elements[i].checked = elements[i].value === value;
    }

}
/////////////////////////////////////////////////////////////

function setMapSize(value){
    map_size_str = value;
}

/////////////////////////////////////////////////////////////

function setMinTime(value){
    time_min_turn = value;
}

/////////////////////////////////////////////////////////////

function setMaxTime(value){
    time_max_turn = value;
}



/////////////////////////////////////////////////////////////
function setWall(wall_value){
    wall = wall_value;
    if(wall === 0){screen_snake.style.borderColor = "#606060";}
    if(wall === 1){screen_snake.style.borderColor = "#FFFFFF";}
}