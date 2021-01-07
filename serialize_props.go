package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/scrouthtv/card-games/doko"
	"github.com/scrouthtv/card-games/logic"
)

func serveProps(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/serialize-props.js" {
		return
	}
	w.Write([]byte(fmt.Sprintf("const cardMaxSuit = %d;\n", logic.CardMaxSuit)))
	w.Write([]byte(fmt.Sprintf("const valueOrder = %s;\n",
		intSliceToString(logic.ValueOrder))))
	w.Write([]byte(fmt.Sprintf("const dokoGameUUID = %d;\n", doko.DokoGameUUID)))
	w.Write([]byte(fmt.Sprintf("const statePreparing = %d;\n",
		logic.StatePreparing)))
	w.Write([]byte(fmt.Sprintf("const statePlaying = %d;\n", logic.StatePlaying)))
	w.Write([]byte(fmt.Sprintf("const stateEnded = %d;\n", logic.StateEnded)))
	w.Write([]byte(fmt.Sprintf("const dokoTrumpOrder = %s;\n",
		cardSliceToString(doko.DokoTrumpOrder))))
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

func cardSliceToString(slice []logic.Card) string {
	var out strings.Builder
	out.WriteString("[")

	var k int
	var c logic.Card
	for k, c = range slice {
		if k > 0 {
			out.WriteString(", ")
		}
		out.WriteString(fmt.Sprintf("new Card(%d, %d)", c.Suit(), c.Value()))
	}
	out.WriteString("]")

	return out.String()
}
