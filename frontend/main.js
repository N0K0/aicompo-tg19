// setup canvas

var canvas = document.querySelector('canvas');
var parent = document.getElementById("parent");
var ctx = canvas.getContext('2d');

var parent = document.getElementById("parent");
canvas.width = width = parent.offsetWidth;
canvas.height = height = parent.offsetHeight;



// Function to add lines to the log
function toLog(text) {
    var log = document.getElementById("logArea");

    var item = document.createElement("div");
    item.innerText = text;

    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

// function to generate random number
function random(min,max) {
  var num = Math.floor(Math.random()*(max-min)) + min;
  return num;
}


function loop() {
  ctx.fillStyle = 'rgba(0,0,0,0.25)';
  ctx.fillRect(0,0,width,height);

  
  requestAnimationFrame(loop);
}


loop();
