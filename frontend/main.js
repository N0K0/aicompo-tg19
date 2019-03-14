
// function to generate random number
// noinspection JSUnusedGlobalSymbols
function random(min,max) {
  return  Math.floor(Math.random()*(max-min)) + min;
}

let admin_ws = new AdminConnection();
admin_ws.adminConnect();
let viewer = new Viewer();
viewer.viewConnect();
render_loop(viewer);