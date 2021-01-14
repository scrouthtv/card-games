import {
	dokoGameUUID, statePreparing, statePlaying, stateEnded,
	cardMaxSuit, dokoTrumpOrder }
from "./serialize-props.mjs";

class Card {
  /**
   * @param {number} suit
   * @param {number} value
   */
  constructor(suit, value) {
    /** @type {number} */
    this.suit = suit;
    /** @type {number} */
    this.value = value;
  }

  static fromBinary(byte) {
    return new Card(byte % cardMaxSuit, Math.floor(byte / cardMaxSuit));
  }

  toString() {
    var out = "";
    switch (this.suit) {
      case 0:
        out = "C";
        break;
      case 1:
        out = "D";
        break;
      case 2:
        out = "H";
        break;
      case 3:
        out = "S";
        break;
      default:
        out = "X";
        break;
    }
    switch (this.value) {
      case 1:
        out += "a";
        break;
      case 11:
        out += "j";
        break;
      case 12:
        out += "q";
        break;
      case 13:
        out += "k";
        break;
      default:
        out += this.value.toString();
        break;
    }
    return out;
  }

	static fromString(str) {
		str = str.toLowerCase();
		var card = new Card();
		if (str.length < 2) {
			return undefined; 
		}
		switch (str.substring(0, 1)) {
			case "c":
				card.suit = 0;
				break;
			case "d":
				card.suit = 1;
				break;
			case "h":
				card.suit = 2;
				break;
			case "s":
				card.suit = 3;
				break;
			default:
				return undefined;
		}
		switch (str.substring(1)) {
			case "a":
				card.value = 1;
				break;
			case "j":
				card.value = 11;
				break;
			case "q":
				card.value = 12;
				break;
			case "k":
				card.value = 13;
				break;
			default:
				card.value = parseInt(str.substring(1));
		}
		return card;
	}

  /**
   * @param {Card} other
   */
  equal(other) {
    if (other == undefined) return false;
    return this.suit == other.suit && this.value == other.value;
  }
}

/**
 * @class
 */
class Deck {
  constructor() {
    /** @type {Card} */
    this.cards = [];
  }

  /**
   * @param {Card} card
   */
  addCard(card) {
    this.cards.push(card);
  }

  toString() {
    return this.cards.toString();
  }

  length() {
    return this.cards.length;
  }

  get(idx) {
    return this.cards[idx];
  }

  /**
   * @param {ByteBuffer} buf
   * @returns {Deck}
   */
  static fromBinary(buf) {
    var deck = new Deck();
    var length = buf.getUint8();
    for (var i = 1; i <= length; i++) {
      deck.addCard(Card.fromBinary(buf.getUint8()));
    }
    return deck;
  }
}

class DokoScore {

	constructor() {
		this.scores = [0, 0, 0, 0];
		this.rereasons = [];
		this.contrareasons = [];
	}

	/**
	 * @param {ByteBuffer} buf
	 */
	static fromBinary(buf) {
		var ds = new DokoScore();
		var slen = buf.getUint8();
		for (var i = 0; i < slen; i++) {
			ds.scores.push(buf.getUint8());
		}
		ds.rereasons = arrayFromBuf(buf);
		ds.contrareasons = arrayFromBuf(buf);
		return ds;
	}

}

/**
 * @param {ByteBuffer} buf
 */
function arrayFromBuf(buf) {
	var arr = [];
	var len = buf.getUint8();
	for (var i = 0; i < len; i++) {
		arr.push(buf.getUint8());
	}
	return arr;
}

/**
 * @interface
 */
class Ruleset {}

/**
 * @implements {Ruleset}
 */
class DokoGame {

	constructor() {
		this.state = statePreparing;
		this.active = 0;
		this.me = 0;
		this.playable = 0;
		this.hand = new Deck();
		this.table = new Deck();
		this.won = [];
		this.special = [];
	}

  /**
   * @param {ByteBuffer} buf
   */
  static fromBinary(buf) {
    /** @type{DokoGame} */
    var dg = new DokoGame();
    var stateInfo = buf.getInt8();
    var state = stateInfo & 0b11;
    switch (state) {
      case statePreparing: {
        dg.state = statePreparing;
        break;
			}
      case statePlaying: {
        dg.state = statePlaying;
        dg.active = 	(stateInfo & 0b00001100) >> 2;
        dg.me = (stateInfo & 0b00110000) >> 4;
				dg.playable = (stateInfo & 0b01000000) >> 6;
        dg.hand = Deck.fromBinary(buf);
        dg.table = Deck.fromBinary(buf);
        dg.won = [];
        let i;
        for (i = 0; i < 4; i++) dg.won[i] = buf.getUint8();

        dg.special = [];
        for (let player = 0; player < 4; player++) {
          player = buf.getUint8();
          dg.special[player] = Deck.fromBinary(buf);
        }
        break;
			}
      case stateEnded: {
        dg.state = stateEnded;

				dg.scores = DokoScore.fromBinary(buf);

				dg.won = [];
				for (let player = 0; player < 4; player++) {
					player = buf.getUint8();
					dg.won[player] = Deck.fromBinary(buf);
				}
        break;
			}
      default:
        dg.state = -1;
        break;
    }

    return dg;
  }

  allowedCards() {
    if (this.table.length() == 0) {
      return this.hand.cards;
    }

    const show = this.table.get(0);
    var allowed = [];
    var has = this.hand;

    for (var i = 0; i < has.length(); i++) {
      var ownedCard = has.get(i);
      if (this.color(ownedCard) == this.color(show)) {
        allowed.push(ownedCard);
      }
    }

    if (allowed.length == 0) return this.hand.cards;
    return allowed;
  }

  /**
   * @param {Card} card
   */
  color(card) {
    if (this.trumpValue(card) == -1) {
      return card.suit;
    }
    return -1;
  }

  /**
   * @param {Card} card
   */
  trumpValue(card) {
    /**
     * @type {number}
     */
    var value;
    /**
     * @type {Card}
     */
    var trump;
    for (value = 0; value < dokoTrumpOrder.length; value++) {
      trump = dokoTrumpOrder[value];
      if (trump == card.toString().toLowerCase()) {
        return dokoTrumpOrder.length - value;
      }
    }
    return -1;
  }
}

class Game {
  /**
   * @param {ByteBuffer} buf
   */
  static fromBinary(buf) {
    /** @type {Game} */
    var g = new Game();
    switch (buf.getInt8()) {
      case dokoGameUUID:
				console.log("This is a doko game");
        g.ruleset = DokoGame.fromBinary(buf);
        break;
      default:
        console.log("Unknown game");
        break;
    }

    return g;
  }
}

class ByteBuffer {
  /**
   * @param {ArrayBuffer} buf
   */
  constructor(buf) {
    this.dataView = new DataView(buf);
    this.offset = 0;
  }

  getInt8() {
    return this.dataView.getInt8(this.offset++);
  }

  getUint8() {
    return this.dataView.getUint8(this.offset++);
  }

  hasNext() {
    return this.offset < this.dataView.byteLength;
  }
}

export { Card, Deck, Ruleset, DokoGame, Game, ByteBuffer };
