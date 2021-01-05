

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
            case 0: out = "C"; break;
            case 1: out = "D"; break;
            case 2: out = "H"; break;
            case 3: out = "S"; break;
            default: out = "X"; break;
        }
        switch (this.value) {
            case 1: out += "a"; break;
            case 11: out += "j"; break;
            case 12: out += "q"; break;
            case 13: out += "k"; break;
            default: out += this.value.toString(); break;
        }
        return out;
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
        return this.cards.toString()
    }

    /**
     * @param {ByteBuffer} buf
     * @returns {Deck}
     */
    static fromBinary(buf) {
        var deck = new Deck();
        dv.byteOffset = 1;
        var length = buf.getUint8();
        for (var i = 1; i <= length; i++) {
            deck.addCard(Card.fromBinary(buf.getUint8()));
        }
    }
}

/**
 * @interface
 */
class Ruleset {

}

/**
 * @implements {Ruleset}
 */
class DokoGame {
    /**
     * @param {ByteBuffer} buf 
     */
    static fromBinary(buf) {
        /** @type{DokoGame} */
        var dg = new DokoGame();
        var state = buf.getInt8() & 0b11
        switch (state) {
            case statePreparing:
                dg.state = statePreparing;
                console.log("We are preparing");
                break;
            case statePlaying:
                dg.state = statePlaying;
                console.log("We are playing");
                break;
            case stateEnded:
                dg.state = stateEnded;
                console.log("We have ended");
                break;
            default:
                dg.state = -1
                console.log("Unknown state!");
                break;
        }

        return dg;
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
                g.ruleset = DokoGame.fromBinary(buf)
                break;
            default:
                console.log("Unknown game")
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
}