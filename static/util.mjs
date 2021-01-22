function parentOffset(elem) {
	var offset = { "top": 0, "left": 0 };
	var paren = elem.offsetParent;

	while (paren != null) {
		offset.top += paren.offsetTop;
		offset.left += paren.offsetLeft;
		paren = paren.offsetParent;
	}

	return offset;
}

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

export { parentOffset, getRandomInt };
