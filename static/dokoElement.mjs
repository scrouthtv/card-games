import { Game } from "./serialize.mjs";
import { statePreparing, statePlaying, stateEnded } from "./serialize-props.mjs";
import { CardElement } from "./cardElement.mjs";
export { DokoAreaElement };

class StorageElement extends HTMLElement {
	
	constructor(screen, who) {
		super();

		this.screen = screen;
		this.who = who;

		this.root = this.attachShadow({ mode: "closed" });
		this.root.innerHTML = "<link rel=\"stylesheet\" href=\"doko-storage.css\" />";

		var message = document.createElement("span");
		message.id = "player" + who + "message";
		message.innerHTML = "Hier ist ein Stich:";
		message.classList.add("message");
		this.root.appendChild(message);
		this.message = message;

		this.root.appendChild(document.createElement("br"));
		this.root.appendChild(document.createElement("br"));

		var storage = document.createElement("span");
		storage.id = "storage" + who;
		storage.classList.add("storage");
		this.root.appendChild(storage);
		this.storage = storage;

		var specials = document.createElement("span");
		specials.id = "special" + who;
		specials.classList.add("special");
		this.root.appendChild(specials);
		this.specials = specials;
	}

	/**
	 * @param {playerID} the id of the player to draw
	 */
	update() {
		var rs = this.screen.logic.ruleset;

		var amount = rs.won[this.who];
		var specials = rs.special[this.who];
		var specialAmt = 0;
		if (specials != undefined) specialAmt = specials.cards.length;
		console.log(specialAmt);

		for (var i = this.storage.children.length; i < amount - specialAmt; i++) {
			var elem = new CardElement();
			elem.classList.add("card");
			elem.classList.add("small");
			this.storage.appendChild(elem);
		}

		if (specials == undefined) return;

		// Remove special cards that have become irrelevant since last update:
		while (this.specials.children.length > specialAmt) {
			this.specials.removeChild(this.specials.lastChild);
		}

		for (i = 0; i < this.specials.children.length; i++) {
			if (this.specials.children[i].getCard() != specials.cards[i]) {
				this.specials.children[i].setCard(specials.cards[i]);
			}
		}
		for (; i < specials.cards.length; i++) {
			document.getElementById("adding a special card");
			var card = new CardElement();
			card.classList.add("card");
			card.classList.add("small");
			card.setCard(specials.cards[i]);
			this.specials.appendChild(card);
		}
	}
}

customElements.define("doko-storage", StorageElement);

class DokoAreaElement extends HTMLElement {

	constructor(conn) {
		super();

		this.conn = conn;

		this.logic = null;

		this.root = this.attachShadow({ mode: "open" });
		this.root.innerHTML = "<link rel=\"stylesheet\" href=\"doko.css\" />";

		this.initScreen();
	}

	play(evt) {
		if (this.logic.ruleset.active != this.logic.ruleset.me) return;
		var card = evt.target.getCard();
		this.conn.send(`card ${card.toString().toLowerCase()}`);
	}

	pickup() {
		if (this.logic.ruleset.active != this.logic.ruleset.me) return;
		this.conn.send("pickup");
	}

	/**
	 * @type {buf} ByteBuffer
	 */
	msg(buf) {
		if (buf.dataView.byteLength > 1) {
			this.logic = Game.fromBinary(buf);
		}
		this.redraw();
	}

	initScreen() {
		var elem = document.createElement("div");
		elem.id = "textinfo";
		this.root.appendChild(elem);
		this.textinfo = elem;

		// Create the table
		elem = document.createElement("span");
		elem.id = "table";
		var card;
		for (var i = 1; i <= 4; i++) {
			card = new CardElement();
			card.classList.add("card", "hidden");
			card.id = "table" + i;
			card.onclick = () => this.pickup();
			elem.appendChild(card);
		}
		this.root.appendChild(elem);
		this.table = elem;

		// Create my hand
		elem = document.createElement("span");
		elem.id = "hand";
		for (i = 1; i <= 12; i++) {
			card = new CardElement();
			card.classList.add("card");
			card.id = "hand" + i;
			card.onclick = (evt) => this.play(evt);
			elem.appendChild(card);
		}
		this.root.appendChild(elem);
		this.hand = elem;

		this.storage = [];
		for (i = 0; i < 4; i++) {
			this.storage[i] = new StorageElement(this, i);
			this.root.appendChild(this.storage[i]);
			this.storage[i].id = "player" + i;
		}
	}

	redraw() {
		if (this.logic == undefined) {
			return;
		}

		if (this.logic.ruleset.state == statePreparing)
			this.textinfo.innerHTML = "Game is still preparing";
		else if (this.logic.ruleset.active == this.logic.ruleset.me)
			this.textinfo.innerHTML = "It's your turn!";
		else
			this.textinfo.innerHTML = `
							It's ${this.logic.ruleset.active}'s turn,
							you are ${this.logic.ruleset.me}!`;

		if (this.logic.ruleset.state == statePlaying) {
			var hand = this.logic.ruleset.hand.cards;
			var table = this.logic.ruleset.table.cards;
			var allowed = this.logic.ruleset.allowedCards();
			var elem;
			for (var i = 0; i < hand.length; i++) {
				elem = this.hand.children.item(i);
				elem.setCard(hand[i]);
				elem.classList.remove("hidden");
				if (this.logic.ruleset.playable) {
					if (allowed.includes(hand[i])) {
						elem.classList.add("allowed");
					} else {
						elem.classList.remove("allowed");
					}
				} else {
					elem.classList.remove("allowed");
				}

				if (this.logic.ruleset.playable && this.logic.ruleset.active == this.logic.ruleset.me) {
					elem.classList.add("active");
				} else {
					elem.classList.remove("active");
				}
			}

			// hide the other cards:
			for (; i < 12; i++)
				this.hand.children.item(i).classList.add("hidden");

			// set the table:
			for (i = 0; i < table.length; i++) {
				elem = this.table.children.item(i);
				elem.setCard(table[i]);
				elem.classList.remove("hidden");
				if (!this.logic.ruleset.playable) {
					elem.classList.add("allowed");
					if (this.logic.ruleset.active == this.logic.ruleset.me) {
						elem.classList.add("active");
					} else {
						elem.classList.remove("active");
					}
				} else {
					elem.classList.remove("allowed");
				}
			}

			// hide the rest of the table:
			for (; i < 4; i++)
				this.table.children.item(i).classList.add("hidden");

			for (i = 0; i < 4; i++) {
				if (i < this.logic.ruleset.me) {
					this.storage[i].update();
					this.storage[i].id = "player" + i;
				} else if (i == this.logic.ruleset.me) {
					this.storage[i].update();
					this.storage[i].id = "playerme";
				} else {
					this.storage[i].update();
					this.storage[i].id = "player" + (i - 1);
				}
			} 
		} else if (this.logic.ruleset.state == stateEnded) {
			console.log(this.logic.ruleset);
		} else {
			console.log("Not painting this state (" + this.logic.ruleset.state + ")");
			console.log(this.logic);
			console.log(statePreparing);
			console.log(statePlaying);
			console.log(stateEnded);
		}
	}

	drawStorage(who, destination) {
		var storage = this.storage[destination];
		var amount = this.logic.ruleset.won[who];
		var specials = this.logic.ruleset.special[who];
		var specialAmt = 0;
		if (specials != undefined) specialAmt = specials.cards.length;

		for (var i = storage.children.length; i < amount - specialAmt; i++) {
			var elem = new CardElement();
			elem.classList.add("card");
			elem.classList.add("small");
			storage.appendChild(elem);
		}

		if (specials == undefined) return;
		console.log(this.storage[destination]);
		storage = this.storage[destination].children.item(3);

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
			var card = new CardElement();
			card.classList.add("card", "small");
			card.setCard(specials.cards[i]);
			storage.appendChild(card);
		}
	}

}
