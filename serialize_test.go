package main

import (
	"testing"
)

func TestSlicePrinter(t *testing.T) {
	var str string = intSliceToString([]int{2, 1, 3, 5})
	if str != "[2, 1, 3, 5]" {
		t.Errorf("Wrong string representation: %s", str)
	}
}
