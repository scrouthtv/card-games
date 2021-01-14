import { Game } from "./serialize.mjs";
import { statePreparing, statePlaying, stateEnded } from "./serialize-props.mjs";
import { CardElement } from "./cardElement.mjs";
export { DokoAreaElement };

class StorageElement extends HTMLElement {
	
	constructor(logic, who) {
		super();

		this.logic = logic;
		this.player = who;

		this.root = this.attachShadow({ mode: "closed" });

		var elem = document.createElement("span");
		elem.id = "player" + who;

		var message = document.createElement("span");
		message.id = "player" + who + "message";
		message.innerHTML = "Hier ist ein Stich:";
		elem.appendChild(message);
		this.message = message;

		elem.appendChild(document.createElement("br"));
		elem.appendChild(document.createElement("br"));

		var storage = document.createElement("span");
		storage.id = "storage" + storage;
		elem.appendChild(storage);
		this.storage = storage;

		var specials = document.createElement("span");
		specials.id = "special" + storage;
		elem.appendChild(specials);
		this.specials = specials;
		
		this.root.appendChild(elem);
	}

	update() {
		console.log(this.logic);
	}

}
customElements.define("doko-storage", StorageElement);

class DokoAreaElement extends HTMLElement {

	constructor(conn) {
		super();

		console.log("here we have a new game");

		this.conn = conn;

		this.logic = null;

		this.root = this.attachShadow({ mode: "closed" });

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
		for (i = 0; i < 4; i++)
			this.storage[i] = new StorageElement(i, this.ruleset);
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
				elem = this.table.children.items(i);
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

			for (var storage of this.storage) {
				storage.update();
			}
		} else if (this.logic.ruleset.state == stateEnded) {
			console.log(this.logic);
		} else {
			console.log("Not painting this state (" + this.logic.state + ")");
			console.log(this.logic);
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
		storage = this.storage[destination].children.items(3);

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
