class Viewer{
    constructor() {
        this.connection = null;
    }
}

Viewer.prototype.viewConnect = function() {
    if(this.connection != null){
        toLog("Client1 already connected!");
    }
    toLog("Connecting view");
    let viewer = this;
    this.connection = new WebSocket("ws://localhost:8080/view");
    // noinspection JSUnusedLocalSymbols
    this.connection.onopen = function(evt) {
        toLog("Viewer Open");
    };
    // noinspection JSUnusedLocalSymbols
    this.connection.onclose = function(evt) {
        toLog("Websocket Closed");
        viewer.connection = null;
    };
    this.connection.onmessage = function(evt) {
        toLog(evt.data);
        console.log(evt);
        console.log(evt.data);
    };
    this.connection.onerror = function(evt) {
        toLog("ERROR: " + evt.data);
    };
};
