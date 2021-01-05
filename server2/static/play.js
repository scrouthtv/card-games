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
        alert("<b>Connection closed.</b>");
    };
    conn.onmessage = function (evt) {
        game = Game.fromBinary(new ByteBuffer(evt.data));
        redraw();
    };
    conn.onopen = function (evt) {
        console.log("Connection is open");
        join(id);
    };
} else {
    alert("<b>Your browser does not support WebSockets.</b>");
}

/** @type {Game} */
var game;

const gameArea = document.getElementById("gamescreen");
changeSize();

function redraw() {
    if (game == undefined) {
        return;
    }

    console.log("redrawing");

    if (game.ruleset.state == statePlaying) {
        console.log(game);
        var hand = game.ruleset.hand.cards;
        console.log(hand);
        for (var i = 0; i < hand.length; i++) {
            document.getElementById("hand" + (i + 1)).src = "deck/card-deck-" + hand[i].toString().toLowerCase() + ".png";
        }
    } else {
        console.log("Not painting this state");
    }
}

function changeSize() {
    gameArea.width = window.innerWidth;
    gameArea.height = window.innerHeight;
    redraw();
}

window.onresize = changeSize;