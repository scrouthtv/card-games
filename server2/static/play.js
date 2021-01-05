function sendSomething() {

    if (!conn) {
        console.log("Not yet connected");
        return false;
    }
    conn.send("hello");
    return false;
};

function join(id) {
    if (!conn) {
        console.log("Failed to connect");
    }
    conn.send(`join ${id}`);
}

var conn;
var id;

const url = window.location.search;
const params = new URLSearchParams(url);
id = [...params.keys()][0];
console.log(id);

if (id == undefined) {
    window.location.href = "/";
}

if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.binaryType = "arraybuffer"
    conn.onclose = function (evt) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
    };
    conn.onmessage = function (evt) {
        console.log("Server sent" + evt.data + ":");
        console.log(new Uint8Array(evt.data))
    };
    conn.onopen = function (evt) {
        console.log("Connection is open");
        join(id);
    };
} else {
    alert("<b>Your browser does not support WebSockets.</b>");
}