
// function to generate random number
// noinspection JSUnusedGlobalSymbols
function random(min,max) {
  return  Math.floor(Math.random()*(max-min)) + min;
}

adminConnect();

let viewer = new Viewer();
viewer.viewConnect();
render_loop(viewer);