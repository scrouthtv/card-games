package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/scrouthtv/card-games/doko"
	"github.com/scrouthtv/card-games/logic"
)

func writeProps(w io.Writer) {
	writeJSintconst(w, "cardMaxSuit", logic.CardMaxSuit)
	writeJSconst(w, "valueOrder", intSliceToString(logic.ValueOrder))
	writeJSintconst(w, "dokoGameUUID", doko.DokoGameUUID)
	writeJSintconst(w, "statePreparing", logic.StatePreparing)
	writeJSintconst(w, "statePlaying", logic.StatePlaying)
	writeJSintconst(w, "stateEnded", logic.StateEnded)
	writeJSconst(w, "dokoTrumpOrder", cardSliceToString(doko.DokoTrumpOrder))
	writeJSintconst(w, "reasonWon", doko.ReasonWon)
	writeJSintconst(w, "reasonAgainstTheElders", doko.ReasonAgainstTheElders)
	writeJSintconst(w, "reasonNo90", doko.ReasonNo90)
	writeJSintconst(w, "reasonNo60", doko.ReasonNo60)
	writeJSintconst(w, "reasonNo30", doko.ReasonNo30)
	writeJSintconst(w, "reasonBlack", doko.ReasonBlack)
	w.Write([]byte("\nexport { cardMaxSuit, valueOrder, dokoGameUUID, statePreparing, statePlaying, stateEnded, dokoTrumpOrder };\n"))
}

func writeJSintconst(w io.Writer, name string, value int) {
	writeJSconst(w, name, strconv.Itoa(value))
}

func writeJSconst(w io.Writer, name string, value string) {
	w.Write([]byte(fmt.Sprintf("const %s = %s;\n", name, value)))
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
	out.WriteString("[ ")

	var k int
	var c logic.Card
	for k, c = range slice {
		if k > 0 {
			out.WriteString(", ")
		}
		out.WriteString(fmt.Sprintf("\"%s\"", c.Short()))
	}
	out.WriteString(" ]")

	return out.String()
}
