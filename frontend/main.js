// setup canvas

let canvas = document.querySelector('canvas');
let parent = document.getElementById("parent");
let ctx = canvas.getContext('2d');

canvas.width = width = parent.offsetWidth;
canvas.height = height = parent.offsetHeight;



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
  ctx.fillStyle = 'rgba(0,0,0,0.25)';
  ctx.fillRect(0,0,width,height);

  
  requestAnimationFrame(loop);
}


loop();
