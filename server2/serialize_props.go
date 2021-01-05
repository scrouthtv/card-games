package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func serveProps(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/serialize-props.js" {
		return
	}
	log.Print("Serving props.js")
	w.Write([]byte(fmt.Sprintf("const cardMaxSuit = %d;\n", cardMaxSuit)))
	w.Write([]byte(fmt.Sprintf("const valueOrder = %s;\n", intSliceToString(valueOrder))))
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
