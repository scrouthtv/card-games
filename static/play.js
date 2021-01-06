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

    if (game.ruleset.state == statePreparing)
        document.getElementById("textinfo").innerHTML = "Game is still preparing";
    else if (game.ruleset.active == game.ruleset.me) 
        document.getElementById("textinfo").innerHTML = "It's your turn!";
    else
        document.getElementById("textinfo").innerHTML = `
            It's ${game.ruleset.active}'s turn,
            you are ${game.ruleset.me}!`;

    if (game.ruleset.state == statePlaying) {
        var hand = game.ruleset.hand.cards;
        var table = game.ruleset.table.cards;
			  var allowed = game.ruleset.allowedCards();
			  var elem;
        for (var i = 0; i < hand.length; i++) {
						elem = document.getElementById("hand" + (i + 1));
            elem.setCard(hand[i]);
            elem.classList.remove("hidden");
						if (allowed.includes(hand[i])) {
							elem.classList.add("allowed");
						} else {
							elem.classList.remove("allowed");
						}
        }
        for (; i < 12; i++)
            document.getElementById("hand" + (i + 1)).classList.add("hidden");
        for (i = 0; i < table.length; i++) {
            document.getElementById("table" + (i + 1)).setCard(table[i]);
            document.getElementById("table" + (i + 1)).classList.remove("hidden");
        }
        for (; i < 4; i++)
            document.getElementById("table" + (i + 1)).classList.add("hidden");
            
        //var storage = document.getElementById("storageme");
        var me = game.ruleset.me;
        for (i = 0; i < 4; i++) {
            if (i < me) drawStorage(game.ruleset.won[i], "storage" + i);
            else if (i == me) drawStorage(game.ruleset.won[i], "storageme");
            else drawStorage(game.ruleset.won[i], "storage" + (i - 1));
        }
    } else {
        console.log("Not painting this state");
    }
}

function drawStorage(amount, destination) {
    var storage = document.getElementById(destination);
    for (i = storage.children.length; i < amount; i++) {
        var card = new CardElement();
        card.classList.add("card");
        card.classList.add("small");
        storage.appendChild(card);
    }
}

function changeSize() {
    gameArea.width = window.innerWidth;
    gameArea.height = window.innerHeight;
    redraw();
}

window.onresize = changeSize;
