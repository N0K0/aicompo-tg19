// setup canvas

let canvas = document.querySelector('canvas');
let parent = document.getElementById("parent");
let ctx = canvas.getContext('2d');

canvas.width = width = parent.offsetWidth;
canvas.height = height = parent.offsetHeight;

let w = canvas.width;
let h = canvas.height;

let viewer = new Viewer();

function paint () {
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, w, h);
    ctx.strokeStyle = "black";
    ctx.strokeRect(0, 0, w, h);
}

function paint_cell(x, y, color) {
    ctx.fillStyle = color;
    ctx.fillRect(x*cw, y*cw, cw, cw);
    ctx.strokeStyle = "white";
    ctx.strokeRect(x*cw, y*cw, cw, cw);
}


// Function to add lines to the log
function toLog(text) {
    let log = document.getElementById("logArea");

    let item = document.createElement("div");
    item.innerText = text;

    let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

// function to generate random number
// noinspection JSUnusedGlobalSymbols
function random(min,max) {
  return  Math.floor(Math.random()*(max-min)) + min;
}


function loop() {
  paint();

  requestAnimationFrame(loop);
}

viewer.viewConnect();
loop();
