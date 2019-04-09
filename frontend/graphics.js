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

    this.mapsizeX = game_status.GameStatus.GameMap.SizeX -1;
    this.mapsizeY = game_status.GameStatus.GameMap.SizeY -1;

    // Render snakes
    let snakes = game_status.Players;
    for (let snake in snakes){
        let s = snakes[snake];
        console.log(s);
        let col_str = s.Color;

        for(let b in s.PosX){
            let x = s.PosX[b];
            let y = s.PosY[b];
            // TODO: Respect colors
            this.paint_cell(x,y,col_str)
        }
    }
    // Render food
    let foods = game_status.GameStatus.GameMap.Foods;
    for (let food in foods) {
        let f = foods[food];
        console.log(f);
        let x = f.X;
        let y = f.Y;
        this.paint_cell(x,y,"green");
    }

    // Render walls

    let walls = game_status.GameStatus.GameMap.Walls;
    for (let wall in walls){
        let w = walls[wall];
        let x = w.X;
        let y = w.Y;

        this.paint_cell(x,y,"white");
    }

    // Render UI
    switch (game_status.GameStatus.Status) {
        case "pregame":
            this.ui_pregame();
            break;
        case "running":
            this.ui_running();
            break;
        case "done":
            this.ui_done();
            break;
        default:
            break;
    }

};

Viewer.prototype.paint_cell = function(x, y, color) {
    //console.log("Painting cell", x, y, color);
    let h = Math.floor(ctx.height / (this.mapsizeY +1));
    let w = Math.floor(ctx.width / (this.mapsizeX+1));

    ctx.fillStyle = color;
    ctx.fillRect(x*w, y*h, w, h);
};

Viewer.prototype.ui_pregame = function (game_status) {

};

Viewer.prototype.ui_running = function (game_status) {

};

Viewer.prototype.ui_done = function (game_status) {

};

function render(){
    requestAnimationFrame(render);
}

function render_loop() {
    requestAnimationFrame(render);
}
