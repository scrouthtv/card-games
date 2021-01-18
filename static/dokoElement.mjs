import { Card, Game } from "./serialize.mjs";
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
		message.id = "message";
		message.innerHTML = "Hier ist ein Stich:";
		message.classList.add("message");
		this.root.appendChild(message);
		this.message = message;

		this.root.appendChild(document.createElement("br"));
		this.root.appendChild(document.createElement("br"));

		var hand = document.createElement("div");
		hand.id = "hand";
		hand.classList.add("hand");
		this.root.appendChild(hand);
		this.hand = hand;
		this.handVisible = false;

		for (let i = 0; i < 12; i++) {
			let card = new CardElement();
			card.classList.add("card", "small");
			this.hand.appendChild(card);
		}

		var storage = document.createElement("span");
		storage.id = "storage";
		storage.classList.add("storage");
		this.root.appendChild(storage);
		this.storage = storage;

		var specials = document.createElement("span");
		specials.id = "special";
		specials.classList.add("special");
		this.root.appendChild(specials);
		this.specials = specials;

		this.setFinished(false);
	}

	/**
	 * If the game has finished, the storage will draw the won cards, 
	 * not only the back side.
	 */
	setFinished(isFinished) {
		this.finished = isFinished;
		if (isFinished) {
			this.classList.add("finished");
			this.specials.classList.add("finished");
			this.classList.remove("playing");
			this.specials.classList.remove("playing");
		} else {
			this.classList.remove("finished");
			this.specials.classList.remove("finished");
			this.classList.add("playing");
			this.specials.classList.add("playing");
		}
	}

	updateHand() {
		var rs = this.screen.logic.ruleset;

		if (this.who == rs.me) {
			const hand = rs.hand.cards;
			const allowed = rs.allowedCards();
			for (let i = 0; i < hand.length; i++) {
				if (!this.handVisible) {
					this.hand.children[i].classList.remove("small");
					this.hand.children[i].setCard(hand[i]);
					this.hand.children[i].onclick = (evt) => this.screen.play(evt);
					// this.hand.children[i].onclick =
					//	(evt) => this.screen.animateHandToTable(this.who, evt.target);
				}

				if (rs.playingState == 1) { // is playing phase?
					if (allowed.includes(hand[i])) { // is card allowed?
						this.hand.children[i].classList.add("allowed");
					} else {
						this.hand.children[i].classList.remove("allowed");
					}
				} else {
					this.hand.children[i].classList.remove("allowed");
				}

				if (rs.playingState == 1 && rs.active == rs.me) {
					this.hand.children[i].classList.add("active");
				} else {
					this.hand.children[i].classList.remove("active");
				}
			}
			if (!this.handVisible) this.handVisible = true;
		}
	}

	/**
	 * @param {playerID} the id of the player to draw
	 */
	update() {
		var rs = this.screen.logic.ruleset;

		this.updateHand();

		if (this.screen.logic.players[this.who] == undefined) {
			this.message.innerHTML = "noch niemand";
		} else {
			this.message.innerHTML = this.screen.logic.players[this.who] + ":";
		}

		var amount = rs.won[this.who];
		if (this.finished && amount != undefined)
			amount = rs.won[this.who].cards.length;
		var specials = rs.special[this.who];
		var specialAmt = 0;
		if (specials != undefined) specialAmt = specials.cards.length;

		for (let i = this.storage.children.length; i < amount - specialAmt; i++) {
			var elem = new CardElement();
			elem.classList.add("card");
			elem.classList.add("small");
			this.storage.appendChild(elem);
		}

		if (this.finished) {
			for (let i = 0; i < amount - specialAmt; i++) {
				this.storage.children.item(i).setCard(rs.won[this.who].cards[i]);
			}
		}

		if (specials == undefined) return;

		// Remove special cards that have become irrelevant since last update:
		while (this.specials.children.length > specialAmt) {
			this.specials.removeChild(this.specials.lastChild);
		}

		for (var i = 0; i < this.specials.children.length; i++) {
			if (this.specials.children[i].getCard() != specials.cards[i]) {
				this.specials.children[i].setCard(specials.cards[i]);
			}
		}
		for (; i < specials.cards.length; i++) {
			var card = new CardElement();
			card.classList.add("card");
			card.classList.add("small");
			card.setCard(specials.cards[i]);
			this.specials.appendChild(card);
		}

		if (this.finished) {
			var score = rs.scores.scores[this.who];
			if (score == 1)
				this.message.innerHTML = "1 Punkt:";
			else
				this.message.innerHTML = score + " Punkte:";
		}
	}
}

customElements.define("doko-storage", StorageElement);

class DokoCallerElement extends HTMLElement {

	constructor(screen, conn) {
		super();

		this.screen = screen;
		this.conn = conn;

		this.root = this.attachShadow({ mode: "open" });
		this.root.innerHTML = "<link rel=\"stylesheet\" href=\"doko-call.css\">";

		this.contents = document.createElement("div");
		this.contents.id = "contents";

		this.root.appendChild(this.contents);
	}

	addButtons() {
		let calls = this.screen.logic.ruleset.availableCalls();

		this.contents.innerHTML = "";

		var btn;
		for (let call of calls) {
			this.contents.appendChild(document.createElement("br"));
			btn = document.createElement("button");
			btn.innerHTML = call.name;
			btn.onclick = () => call.callback(this.conn);
			this.contents.appendChild(btn);
		}
	}

	activate() {
		this.contents.classList.add("active");
		var nodes = this.contents.children;
		for (let node of nodes) {
			if (node.nodeName == "BUTTON") {
				node.disabled = false;
			}
		}
	}

	deactivate() {
		this.contents.classList.remove("active");
		var nodes = this.contents.children;
		for (let node of nodes) {
			if (node.nodeName == "BUTTON") {
				node.disabled = true;
			}
		}
	}

	update() {
		this.addButtons();
		if (this.screen.logic.ruleset.active == this.screen.logic.ruleset.me) this.activate();
		else this.deactivate();
	}

}

customElements.define("doko-caller", DokoCallerElement);

class DokoAreaElement extends HTMLElement {

	constructor(conn) {
		super();

		this.conn = conn;

		this.logic = null;

		this.hasInit = false;

		this.root = this.attachShadow({ mode: "open" });
		this.root.innerHTML = "<link rel=\"stylesheet\" href=\"doko.css\" />";

		this.tablecards = [];
	}

	play(evt) {
		if (this.logic.ruleset.active != this.logic.ruleset.me) return;
		var card = evt.target.getCard();
		this.conn.send(`card ${card.toString().toLowerCase()}`);
		this.animateHandToTable(0, evt.target);
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
		if (this.hasInit)
			this.redraw();
		else {
			this.initScreen();
			this.hasInit = true;
			this.redraw();
		}
	}

	initScreen() {
		console.log(this.logic);
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

		this.storage = [];
		for (i = 0; i < 4; i++) {
			this.storage[i] = new StorageElement(this, (i + this.logic.ruleset.me) % 4);
			this.root.appendChild(this.storage[i]);
			this.storage[i].id = "player" + i;
		}

		this.caller = new DokoCallerElement(this, this.conn);
		this.caller.classList.add("hidden");
		this.root.appendChild(this.caller);
	}

	NOPRODUCTIONexampleCards() {
		for (let i = 0; i < this.hand.children.length; i++)
			this.hand.children[i].setCard(Card.fromBinary((i + 2) * 3));
		for (let i = 0; i < 2; i++) {
			this.table.children[i].classList.remove("hidden");
			this.table.children[i].setCard(Card.fromBinary(41 - 5 * i));
			this.tablecards.push(this.table.children[i]);
		}
	}

	animateHandToTable(player, element) {
		var target = this.table.children[this.tablecards.length];
		if (target == undefined) return;
		target.classList.remove("hidden");
		this.tablecards.push(element);

		// coords of the target:
		const tstyle = window.getComputedStyle(target);
		const tt = tstyle.transform;
		const to = tstyle.transformOrigin;

		element.style.zIndex = this.tablecards.length;

		// wait for the properties to apply
		window.setTimeout(function() {
			element.style.left = "105px"; // these numbers are magic,
			element.style.top = "-327px"; // don't touch them
			element.style.transform = tt;
			element.style.transformOrigin = to;
		}, 1);

		target.classList.add("hidden");
	}

	updateTable() {
		var table = this.logic.ruleset.table.cards;
		var elem;
		for (var i = 0; i < table.length; i++) {
			elem = this.table.children.item(i);
			elem.setCard(table[i]);
			elem.classList.remove("hidden");
			if (this.logic.ruleset.playingState == 2) {
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

		if (this.logic.ruleset.state == statePreparing) {
			for (let i = 0; i < 4; i++) {
				this.storage[i].update();
			}
		} else if (this.logic.ruleset.state == statePlaying) {
			if (this.logic.ruleset.playingState == 0) {
				//this.updateHand();
				this.caller.classList.remove("hidden");
				this.caller.update();
			} else {
				//this.updateHand();
				this.updateTable();
				this.caller.classList.add("hidden");

				for (let i = 0; i < 4; i++) {
					this.storage[i].update();
				} 
			}
		} else if (this.logic.ruleset.state == stateEnded) {
			console.log("this is the end");
			this.drawEnd();
		} else {
			console.log("Not painting this state (" + this.logic.ruleset.state + ")");
			console.log(this.logic);
			console.log(statePreparing);
			console.log(statePlaying);
			console.log(stateEnded);
		}
	}

	drawEnd() {
		this.hand.classList.add("hidden");
		this.table.classList.add("hidden");

		var re = 0, contra = 0;
		for (let i = 0; i < 4; i++) {
			if (this.logic.ruleset.reteam.includes(i))
				this.storage[i].id = "re" + re++;
			else
				this.storage[i].id = "contra" + contra++;
			this.storage[i].setFinished(true);
			this.storage[i].update();
		}

		if (this.badge == undefined) {
			this.badge = document.createElement("h1");
			this.root.appendChild(this.badge);
			this.badge.innerHTML = "Someone won!";
			this.badge.id = "badge";
			if (this.logic.ruleset.didIWin()) {
				this.badge.classList.add("winner");
				this.badge.innerHTML = "You win!";
			} else {
				this.badge.classList.add("looser");
				this.badge.innerHTML = "You loose!";
			}
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
