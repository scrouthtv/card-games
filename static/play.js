import { Game, ByteBuffer } from "./serialize.mjs";
import { statePreparing, statePlaying, stateEnded, dokoGameUUID } from "./serialize-props.mjs";
import { CardElement } from "./cardElement.mjs";
export { join, play, pickup };

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

if (id == undefined) {
  window.location.href = "/";
}

if (window["WebSocket"]) {
  conn = new WebSocket("ws://" + document.location.host + "/ws");
  conn.binaryType = "arraybuffer";
  conn.onclose = function () {
    console.log("Connection closed.");
  };
  conn.onmessage = function (evt) {
		const initBuf = new ByteBuffer(evt.data);
		const gameID = initBuf.getInt8();
		switch (gameID) {
			case dokoGameUUID:
				setupGame();
				break;
			default:
				console.log("I don't know this game");
				break;
		}
  };
  conn.onopen = function () {
    console.log("Connection is open");
    join(id);
  };
} else {
  alert("<b>Your browser does not support WebSockets.</b>");
}

function play(cint) {
  var card = document.getElementById("hand" + cint).getCard();
  conn.send(`card ${card.toString().toLowerCase()}`);
}

function pickup() {
  if (game.ruleset.active != game.ruleset.me) return;
  conn.send("pickup");
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
      if (game.ruleset.playable) {
        if (allowed.includes(hand[i])) {
          elem.classList.add("allowed");
        } else {
          elem.classList.remove("allowed");
        }
      } else {
        elem.classList.remove("allowed");
      }

      if (game.ruleset.playable && game.ruleset.active == game.ruleset.me) {
        elem.classList.add("active");
      } else {
        elem.classList.remove("active");
      }
    }
    for (; i < 12; i++)
      document.getElementById("hand" + (i + 1)).classList.add("hidden");
    for (i = 0; i < table.length; i++) {
      elem = document.getElementById("table" + (i + 1));
      elem.setCard(table[i]);
      elem.classList.remove("hidden");
      if (!game.ruleset.playable) {
        elem.classList.add("allowed");
        if (game.ruleset.active == game.ruleset.me) {
          elem.classList.add("active");
        } else {
          elem.classList.remove("active");
        }
      } else {
        elem.classList.remove("allowed");
      }
    }
    for (; i < 4; i++)
      document.getElementById("table" + (i + 1)).classList.add("hidden");

    var me = game.ruleset.me;
    for (i = 0; i < 4; i++) {
      if (i < me) {
        drawStorage(i, i);
        document.getElementById("player" + i + "message").innerText =
          "Hier ist " + i + "'s Stich";
      } else if (i == me) drawStorage(i, "me");
      else {
        drawStorage(i, i - 1);
        document.getElementById("player" + (i - 1) + "message").innerText =
          "Hier ist " + i + "'s Stich";
      }
    }
	} else if (game.ruleset.state == stateEnded) {
		console.log(game);
	} else {
    console.log("Not painting this state (" + game.state + ")");
		console.log(game);
		console.log(stateEnded);
	}
}

function drawStorage(who, destination) {
  var storage = document.getElementById("storage" + destination);
  var amount = game.ruleset.won[who];
  var specials = game.ruleset.special[who];
  var specialAmt = 0;
  if (specials != undefined) specialAmt = specials.cards.length;

  for (var i = storage.children.length; i < amount - specialAmt; i++) {
    var elem = new CardElement();
    elem.classList.add("card");
    elem.classList.add("small");
    storage.appendChild(elem);
  }

  if (specials == undefined) return;
  storage = document.getElementById("special" + destination);

  // Remove special cards that have become irrelevant since last update:
  while (storage.children.length > specialAmt) {
    storage.removeChild(storage.lastChild);
  }

  for (i = 0; i < storage.children.length; i++) {
    if (storage.children[i].getCard() != specials.cards[i]) {
      storage.children[i].setCard(specials.cards[i]);
    }
  }
  for (; i < specials.cards.length; i++) {
    document.getElementById("adding a special card");
    var card = new CardElement();
    card.classList.add("card", "small");
    card.setCard(specials.cards[i]);
    storage.appendChild(card);
  }
  
}

function changeSize() {
  gameArea.width = window.innerWidth;
  gameArea.height = window.innerHeight;
  redraw();
}

window.onresize = changeSize;

function setupGame() {
	const screen = document.getElementById("gamescreen");

	var elem = document.createElement("div");
	elem.id = "textinfo";
	screen.appendChild(elem);

	// Create the table
	elem = document.createElement("span");
	elem.id = "table";
	var card;
	for (var i = 1; i <= 4; i++) {
		card = new CardElement();
		card.classList.add("card", "hidden");
		card.id = "table" + i;
		card.setAttribute("onclick", "pickup()");
		elem.appendChild(card);
	}
	console.log(elem);
	screen.appendChild(elem);

	// Create my hand
	elem = document.createElement("span");
	elem.id = "hand";
	for (i = 1; i <= 12; i++) {
		card = new CardElement();
		card.classList.add("card");
		card.id = "hand" + i;
		card.setAttribute("onclick", "play('" + i + "')");
		elem.appendChild(card);
	}
	screen.appendChild(elem);

	// Create the storages
	const storages = ["me", "0", "1", "2"];
	var ielem;
	for (var storage of storages) {
		elem = document.createElement("span");
		elem.id = "player" + storage;

		ielem = document.createElement("span");
		ielem.id = "player" + storage + "message";
		elem.appendChild(ielem);

		ielem = document.createElement("span");
		ielem.id = "storage" + storage;
		elem.appendChild(ielem);

		ielem = document.createElement("span");
		ielem.id = "special" + storage;
		elem.appendChild(ielem);
		
		screen.appendChild(elem);
	}

	conn.onmessage = function(evt) {
		game = Game.fromBinary(new ByteBuffer(evt.data));
		redraw();
	};
}
