var admin_connect = function() {
    toLog("Starting admin connect")
    ws = new WebSocket("ws://localhost:8080/admin");
    ws.onopen = function(evt) {
        toLog("Admin socket started")
    }
    ws.onclose = function(evt) {
        toLog("Admin socket closed!")
        ws = null;
    }
    ws.onmessage = function(evt) {
        toLog("A RESPONSE: " + evt.data);
    }
    ws.onerror = function(evt) {
        toLog("A ERROR: " + evt.data);
    }
}
