import { Card, Deck, Ruleset, DokoGame, Game, ByteBuffer } from './serialize.mjs';
import { statePreparing, statePlaying, stateEnded } from './serialize-props.mjs';
export { join, play, pickup }

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
      var elem = document.getElementById("table" + (i + 1));
      elem.setCard(table[i]);
      elem.classList.remove("hidden");
      if (!game.ruleset.playable) {
        console.log("is not playable");
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
  } else {
    console.log("Not painting this state");
  }
}

function drawStorage(who, destination) {
  var storage = document.getElementById("storage" + destination);
  var amount = game.ruleset.won[who];
  var specials = game.ruleset.special[who];
  var specialAmt = 0;
  if (specials != undefined) specialAmt = specials.cards.length;

  for (var i = storage.children.length; i < amount - specialAmt; i++) {
    var card = new CardElement();
    card.classList.add("card");
    card.classList.add("small");
    storage.appendChild(card);
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
    card.classList.add("card");
    card.classList.add("small");
    card.setCard(specials.cards[i]);
    storage.appendChild(card);
  }
  return;
}

function changeSize() {
  gameArea.width = window.innerWidth;
  gameArea.height = window.innerHeight;
  redraw();
}

window.onresize = changeSize;
