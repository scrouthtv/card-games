import { Card, Deck, DokoGame } from '../static/serialize.mjs';

var cards = [ "dk", "ca", "c10", "ca", "da", "h9", "d9", "s10", "cq", "cj", "dj", "dk" ];
var deck = new Deck();
cards.forEach(e => deck.addCard(Card.fromString(e)));

const doko = new DokoGame();
doko.hand = deck;

var allowed = doko.allowedCards();

if (allowed != deck.cards) {
	console.log("All cards should be allowed right now, they aren't");
	process.exit(1);
}

doko.table.addCard(new Card(2, 13));
allowed = doko.allowedCards();

if (allowed.length != 1) {
	console.log("Wrong amount of cards allowed, should be 2, is " + allowed.length);
	process.exit(1);
}

if (!allowed[0].equal(new Card(2, 9))) {
	console.log("Wrong card allowed, should be H9, is " + allowed[0].toString());
	process.exit(1);
}
