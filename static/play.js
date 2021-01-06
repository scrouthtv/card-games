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

function play(cint) {
    var card = document.getElementById("hand" + cint).getCard();
    console.log(card);
    conn.send(`card ${card.toString().toLowerCase()}`);
}

/** @type {Game} */
var game;

const gameArea = document.getElementById("gamescreen");
changeSize();

function redraw() {
    if (game == undefined) {
        return;
    }

    if (game.ruleset.active == game.ruleset.me) 
        document.getElementById("textinfo").innerHTML = "It's your turn!";
    else
        document.getElementById("textinfo").innerHTML = "It's " + game.ruleset.active + "'s turn!";

    if (game.ruleset.state == statePlaying) {
        //console.log(game);
        var hand = game.ruleset.hand.cards;
        //console.log(hand);
        for (var i = 0; i < hand.length; i++) {
            document.getElementById("hand" + (i + 1)).setCard(hand[i])
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