class Viewer{
    constructor() {
        this.connection = null;
        this.canvas = document.querySelector('canvas');
        this.ctx = this.canvas.getContext('2d');

        this.ctx.height = 800;
        this.ctx.width = 800;

        this.mapsizeX = 0;
        this.mapsizeY = 0;

    }
}

Viewer.prototype.viewConnect = function() {
    if(this.connection != null){
        toLog("Client1 already connected!");
    }
    console.log("Connecting view");
    let viewer = this;
    this.connection = new WebSocket("ws://localhost:8080/view");
    // noinspection JSUnusedLocalSymbols
    this.connection.onopen = function(evt) {
        console.log("Viewer Open");
    };
    // noinspection JSUnusedLocalSymbols
    this.connection.onclose = function(evt) {
        console.log("Viewer Closed");
        viewer.connection = null;
    };
    this.connection.onmessage = function(evt) {
        console.log("Viewer update");
        viewer.render_scene(JSON.parse(evt.data))
    };
    this.connection.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    };
};


Viewer.prototype.paint_canvas = function () {
    let w = ctx.width;
    let h = ctx.height;

    this.ctx.fillStyle = "black";
    this.ctx.fillRect(0, 0, w, h);
    this.ctx.strokeStyle = "black";
    this.ctx.strokeRect(0, 0, w, h);
};

Viewer.prototype.render_scene = function (game_status) {
    this.paint_canvas();
    console.log(game_status);


    this.current_round = game_status["GameStatus"]["RoundNumber"];
    this.round_total = game_status["GameStatus"]["TotalRounds"];

    console.log(this.current_round);
    console.log(this.round_total);

    if(this.current_round > this.round_total) {
        showScreen(5);
        return
    }

    this.mapsizeX = game_status["GameStatus"]["GameMap"]["SizeX"] - 1;
    this.mapsizeY = game_status["GameStatus"]["GameMap"]["SizeY"] - 1;


    // Render food
    let foods = game_status["GameStatus"]["GameMap"]["Foods"];
    for (let food in foods) {
        let f = foods[food];
        //console.log(f);
        let x = f.X;
        let y = f.Y;
        this.paint_cell(x, y, "green");
    }

    // Render walls
    let walls = game_status["GameStatus"]["GameMap"]["Walls"];
    for (let wall in walls) {
        let w = walls[wall];
        let x = w.X;
        let y = w.Y;
        this.paint_cell(x, y, "white");
    }

    // Render snakes
    let snakes = game_status["Players"];
    for (let snake in snakes) {
        let s = snakes[snake];
        let col_str = s["Color"];

        for (let b in s["PosX"]) {
            let x = s["PosX"][b];
            let y = s["PosY"][b];
            // TODO: Respect colors
            this.paint_cell(x, y, col_str)
        }
    }

    // Render head

    // Render tail

    // Render UI

    let div_tick = document.getElementById("ticks");
    let div_scoreboard = document.getElementById("scoreboard");
    let div_round = document.getElementById("round");

    div_tick.innerText = "Tick: " + game_status["GameStatus"]["CurrentTick"];
    div_round.innerText = "Round: " + this.current_round + " / " + this.round_total;

    let div_score_clone = div_scoreboard.cloneNode(false);

    let div_score_clone_text = document.createElement("div");
    div_score_clone.innerText = "Score:";

    div_score_clone.appendChild(div_score_clone_text);

    let players = Object.entries(game_status.Players);
    players = players.sort(compareScore);
    for (let player in players) {
        let p = players[player][1];
        let tmp_div = document.createElement("div");
        //console.log(p);
        tmp_div.innerText = p.username + ": " + p["RoundScore"] + "  (" + p["TotalScore"]+ ")";

        div_score_clone.appendChild(tmp_div)
    }

    div_scoreboard.replaceWith(div_score_clone)
};

// a and b are object elements of your array
function compareScore(a,b) {
    return a["RoundScore"] - b["RoundScore"];
}

Viewer.prototype.paint_cell = function(x, y, color) {
    //console.log("Painting cell", x, y, color);
    let h = Math.floor(ctx.height / (this.mapsizeY +1));
    let w = Math.floor(ctx.width / (this.mapsizeX+1));

    ctx.fillStyle = color;
    ctx.fillRect(x*w, y*h, w, h);
};



function render(){
    requestAnimationFrame(render);
}

function render_loop() {
    requestAnimationFrame(render);
}
