import { Card, DokoGame } from '../static/serialize.mjs'

const doko = new DokoGame();

var c = new Card(2, 10); // hearts 10

console.log(doko.trumpValue(c));
console.log(doko.color(c));
