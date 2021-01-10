package main

import (
	"bufio"
	"os"
	"testing"
)

func TestExportProps(t *testing.T) {
	var f *os.File
	var err error
	f, err = os.Create("./static/serialize-props.mjs")
	if err != nil {
		t.Error("Error creating file:", err)
		t.FailNow()
	}
	var w *bufio.Writer = bufio.NewWriter(f)
	writeProps(w)
	w.Flush()
	f.Close()
	t.Log(f.Name())
}

func TestSlicePrinter(t *testing.T) {
	var str string = intSliceToString([]int{2, 1, 3, 5})
	if str != "[2, 1, 3, 5]" {
		t.Errorf("Wrong string representation: %s", str)
	}
}
