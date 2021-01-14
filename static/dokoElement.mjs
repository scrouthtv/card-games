import { Game } from "./serialize.mjs";
import { statePreparing, statePlaying, stateEnded } from "./serialize-props.mjs";
import { CardElement } from "./cardElement.mjs";
export { DokoAreaElement };

class DokoAreaElement extends HTMLElement {

	constructor(container, conn) {
		super();

		console.log("here we have a new game");

		this.container = container;

		this.conn = conn;

		this.logic = null;

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
			card.onclick = () => this.pickup();
			elem.appendChild(card);
		}
		screen.appendChild(elem);

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
		screen.appendChild(elem);

		// Create the storages
		const storages = ["me", "0", "1", "2"];
		var ielem;
		for (var storage of storages) {
			elem = document.createElement("span");
			elem.id = "player" + storage;

			ielem = document.createElement("span");
			ielem.id = "player" + storage + "message";
			ielem.innerHTML = "Hier ist ein Stich:";
			elem.appendChild(ielem);
			elem.appendChild(document.createElement("br"));
			elem.appendChild(document.createElement("br"));

			ielem = document.createElement("span");
			ielem.id = "storage" + storage;
			elem.appendChild(ielem);

			ielem = document.createElement("span");
			ielem.id = "special" + storage;
			elem.appendChild(ielem);
			
			screen.appendChild(elem);
		}
	}

	redraw() {
		if (this.logic == undefined) {
			return;
		}

		if (this.logic.ruleset.state == statePreparing)
			document.getElementById("textinfo").innerHTML = "Game is still preparing";
		else if (this.logic.ruleset.active == this.logic.ruleset.me)
			document.getElementById("textinfo").innerHTML = "It's your turn!";
		else
			document.getElementById("textinfo").innerHTML = `
							It's ${this.logic.ruleset.active}'s turn,
							you are ${this.logic.ruleset.me}!`;

		if (this.logic.ruleset.state == statePlaying) {
			var hand = this.logic.ruleset.hand.cards;
			var table = this.logic.ruleset.table.cards;
			var allowed = this.logic.ruleset.allowedCards();
			var elem;
			for (var i = 0; i < hand.length; i++) {
				elem = document.getElementById("hand" + (i + 1));
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
			for (; i < 12; i++)
				document.getElementById("hand" + (i + 1)).classList.add("hidden");
			for (i = 0; i < table.length; i++) {
				elem = document.getElementById("table" + (i + 1));
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
			for (; i < 4; i++)
				document.getElementById("table" + (i + 1)).classList.add("hidden");

			var me = this.logic.ruleset.me;
			for (i = 0; i < 4; i++) {
				if (i < me) {
					this.drawStorage(i, i);
					document.getElementById("player" + i + "message").innerText =
						"Hier ist " + i + "'s Stich";
				} else if (i == me) this.drawStorage(i, "me");
				else {
					this.drawStorage(i, i - 1);
					document.getElementById("player" + (i - 1) + "message").innerText =
						"Hier ist " + i + "'s Stich";
				}
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
		var storage = document.getElementById("storage" + destination);
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

}
