
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
let button_back_setting = document.getElementById("setting_back");
button_back_setting.onclick = function() { showScreen(1)};

let button_save_setting = document.getElementById("setting_save");
button_save_setting.onclick = save_settings();

// etc
let ele_score = document.getElementById("score_value");
let speed_setting = document.getElementsByName("speed");
let wall_setting = document.getElementsByName("wall");

// --------------------


setSnakeSpeed(150);
setWall(1);

showScreen("menu");

// --------------------
// Settings

// speed
for(var i = 0; i < speed_setting.length; i++){
    speed_setting[i].addEventListener("click", function(){
        for(var i = 0; i < speed_setting.length; i++){
            if(speed_setting[i].checked){
                setSnakeSpeed(speed_setting[i].value);
            }
        }
    });
}

// wall
for(var i = 0; i < wall_setting.length; i++){
    wall_setting[i].addEventListener("click", function(){
        for(var i = 0; i < wall_setting.length; i++){
            if(wall_setting[i].checked){
                setWall(wall_setting[i].value);
            }
        }
    });
}

document.onkeydown = function(evt){
    if(screen_gameover.style.display == "block"){
        evt = evt || window.event;
        if(evt.keyCode == 32){
            newGame();
        }
    }
};

// 0 for the game
// 1 for the main menu
// 2 for the settings screen
// 3 for the game over screen

function showScreen(screen_opt){
    switch(screen_opt){

        case 0:  screen_snake.style.display = "block";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            break;

        case 1:  screen_snake.style.display = "none";
            screen_menu.style.display = "block";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "none";
            break;

        case 2:  screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "block";
            screen_gameover.style.display = "none";
            break;

        case 3: screen_snake.style.display = "none";
            screen_menu.style.display = "none";
            screen_setting.style.display = "none";
            screen_gameover.style.display = "block";
            break;
    }
}

// Settings logic
function save_settings(){

    showScreen(1)
}


/////////////////////////////////////////////////////////////

// Change the snake speed...
// 150 = slow
// 100 = normal
// 50 = fast
function setSnakeSpeed(speed_value){
    snake_speed = speed_value;
}

/////////////////////////////////////////////////////////////
function setWall(wall_value){
    wall = wall_value;
    if(wall == 0){screen_snake.style.borderColor = "#606060";}
    if(wall == 1){screen_snake.style.borderColor = "#FFFFFF";}
}