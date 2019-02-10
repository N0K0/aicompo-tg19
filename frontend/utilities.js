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