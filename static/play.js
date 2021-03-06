import { ByteBuffer } from "./serialize.mjs";
import { dokoGameUUID } from "./serialize-props.mjs";
import { DokoAreaElement } from "./dokoElement.mjs";

function join(id) {
  if (!conn) {
    console.log("Failed to connect");
  }
  conn.send(`join ${id}`);
}

var conn;

var gamearea;

const url = window.location.search;
const params = new URLSearchParams(url);
var id = [...params.keys()][0];

if (id == undefined) {
  window.location.href = "/";
}

if (window["WebSocket"]) {
  conn = new WebSocket("ws://" + document.location.host + "/ws");
  conn.binaryType = "arraybuffer";
  conn.onclose = function () {
    console.log("Connection closed.");
  };
  conn.onmessage = function (evt) {
		const buf = new ByteBuffer(evt.data);
		if (gamearea != undefined) {
			gamearea.msg(buf);
			return;
		}

		const gameID = buf.getInt8();
		switch (gameID) {
			case dokoGameUUID:
				gamearea = new DokoAreaElement(conn);
				document.getElementById("gamescreen").appendChild(gamearea);
				gamearea.msg(buf);
				break;
			default:
				console.log("I don't know this game");
				break;
		}
  };
  conn.onopen = function () {
    console.log("Connection is open");
    join(id);
  };
} else {
  alert("<b>Your browser does not support WebSockets.</b>");
}

/*function changeSize() {
  gameArea.width = window.innerWidth;
  gameArea.height = window.innerHeight;
  redraw();
}

window.onresize = changeSize; */
