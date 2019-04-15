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

const json_parse_mult = ss => {
    ss = ss.split("\n").map(l => l.trim()).join("");
    let start = ss.indexOf("{");
    let open = 0;
    const res = [];
    for (let i = start; i < ss.length; i++) {
        if ((ss[i] === "{") && (i < 2 || ss.slice(i - 2, i) !== "\\\"")) {
            open++
        } else if ((ss[i] === "}") && (i < 2 || ss.slice(i - 2, i) !== "\\\"")) {
            open--;
            if (open === 0) {
                res.push(JSON.parse(ss.substring(start, i + 1)));
                start = i + 1
            }
        }
    }
    return res
};