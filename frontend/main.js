// setup canvas





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

let viewer = new Viewer();
viewer.viewConnect();
render_loop(viewer);