package main

import (
	"fmt"
	"strings"
)

// Item contains all methods needed for storing and sending items (cards / ...) in the inventories
type Item interface {
	// Short() returns a short alphanumeric identifier for this item,
	// that is also recognizable by the client
	Short() string
}

// Inventory is a collection of int-identified slots which can hold multiple items
type Inventory map[int][]Item

const idelim string = "-"
const odelim string = "#"

// Send converts the inventory into a sendable string
func (i *Inventory) Send() string {
	var out strings.Builder
	var slot int
	var items []Item

	for slot, items = range *i {
		if out.Len() > 0 {
			out.WriteString(odelim)
		}
		out.WriteString(fmt.Sprintf("%d:", slot))
		var i int
		var item Item
		for i, item = range items {
			if i > 0 {
				out.WriteString(idelim)
			}
			out.WriteString(item.Short())
		}
		out.WriteString(odelim)
	}

	return out.String()
}
