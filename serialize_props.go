package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func serveProps(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/serialize-props.js" {
		return
	}
	w.Write([]byte(fmt.Sprintf("const cardMaxSuit = %d;\n", cardMaxSuit)))
	w.Write([]byte(fmt.Sprintf("const valueOrder = %s;\n",
		intSliceToString(valueOrder))))
	w.Write([]byte(fmt.Sprintf("const dokoGameUUID = %d;\n", dokoGameUUID)))
	w.Write([]byte(fmt.Sprintf("const statePreparing = %d;\n",
		StatePreparing)))
	w.Write([]byte(fmt.Sprintf("const statePlaying = %d;\n", StatePlaying)))
	w.Write([]byte(fmt.Sprintf("const stateEnded = %d;\n", StateEnded)))
	w.Write([]byte(fmt.Sprintf("const dokoTrumpOrder = %s;\n",
		cardSliceToString(dokoTrumpOrder))))
}

func intSliceToString(slice []int) string {
	var out strings.Builder
	out.WriteString("[")

	var k, v int
	for k, v = range slice {
		if k > 0 {
			out.WriteString(", ")
		}
		out.WriteString(strconv.Itoa(v))
	}

	out.WriteString("]")

	return out.String()
}

func cardSliceToString(slice []Card) string {
	var out strings.Builder
	out.WriteString("[")

	var k int
	var c Card
	for k, c = range slice {
		if k > 0 {
			out.WriteString(", ")
		}
		out.WriteString(fmt.Sprintf("new Card(%d, %d)", c.suit, c.value))
	}
	out.WriteString("]")

	return out.String()
}
