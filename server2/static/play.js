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

const cardsprite = new Image();
cardsprite.src = "card-deck-161536.svg";
cardsprite.onload = redraw;
const canvas = document.getElementById("gamescreen");
const ctx = canvas.getContext("2d");
changeSize();

function redraw() {
    if (game == undefined) {
        return;
    }

    console.log("redrawing");

    if (game.ruleset.state == statePlaying) {
        console.log(game);
        var hand = game.ruleset.hand.cards;
        for (var i = 0; i < hand.length; i++) {
            let pos = cardPosition(hand[i]);
            console.log("drawing image to ");
            //console.log(pos.x);
            //console.log(pos);
            ctx.drawImage(cardsprite,
                pos.x, pos.y, 79, 123,
                i * 20, 90, 79, 123);
        }
    } else {
        console.log("Not painting this state");
    }

    var pos = cardPosition(new Card(3, 8));
    ctx.drawImage(cardsprite, 
        pos.x, pos.y, 79, 123,
        30, 30, 188, 306);
}

/**
 * @returns {{x: number, y: number}}
 */
function cardPosition(card) {
    return {
        x: card.value * 79,
        y: card.suit * 123
    }
}

function changeSize() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    redraw();
}

window.onresize = changeSize;