import { Card, Deck, Ruleset, DokoGame, Game, ByteBuffer } from '../static/serialize.mjs';

var card = Card.fromBinary(5);
console.log(card);

if (card.suit != 1 || card.value != 1) {
	console.log("Card should be clubs ace");
	process.exit(1);
}

card = Card.fromString("h10");
console.log(card);

if (card.suit != 2 || card.value != 10) {
	console.log("Card should be hearts 10");
	process.exit(1);
}

console.log("Success");
