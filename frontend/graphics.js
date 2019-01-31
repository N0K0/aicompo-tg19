class Viewer{
    constructor() {
        this.connection = null;
        this.canvas = document.querySelector('canvas');
        this.parent = document.getElementById("parent");
        this.ctx = this.canvas.getContext('2d');

        this.canvas.width = this.width = this.parent.offsetWidth;
        this.canvas.height = this.height = this.parent.offsetHeight;
    }
}




Viewer.prototype.viewConnect = function() {
    if(this.connection != null){
        toLog("Client1 already connected!");
    }
    toLog("Connecting view");
    console.log("Connecting view");
    let viewer = this;
    this.connection = new WebSocket("ws://localhost:8080/view");
    // noinspection JSUnusedLocalSymbols
    this.connection.onopen = function(evt) {
        toLog("Viewer Open");
        console.log("Viewer Open");
    };
    // noinspection JSUnusedLocalSymbols
    this.connection.onclose = function(evt) {
        toLog("Viewer Closed");
        console.log("Viewer Closed");
        viewer.connection = null;
    };
    this.connection.onmessage = function(evt) {
        console.log("Viewer update");
        viewer.render_scene(JSON.parse(evt.data))
    };
    this.connection.onerror = function(evt) {
        toLog("ERROR: " + evt.data);
    };
};


Viewer.prototype.paint_canvas = function () {
    let w = this.width;
    let h = this.height;

    this.ctx.fillStyle = "black";
    this.ctx.fillRect(0, 0, w, h);
    this.ctx.strokeStyle = "black";
    this.ctx.strokeRect(0, 0, w, h);
};

Viewer.prototype.render_scene = function (game_status) {
    this.paint_canvas(game_status);

    console.log(game_status);
    console.log(game_status.GameStatus.Status);

    // Render snakes

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
    let h = this.height;
    let w = this.width;

    this.ctx.fillStyle = color;
    this.ctx.fillRect(x*w, y*h, w, h);
    this.ctx.strokeStyle = "white";
    this.ctx.strokeRect(x*w, y*h, w, h);
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
