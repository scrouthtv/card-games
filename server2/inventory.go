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

// NewInventory creates a new inventory by placing
// all provided items in the first (0th) slot
func NewInventory(items []Item) *Inventory {
	var inv Inventory = make(Inventory)

	inv[0] = items

	return &inv
}

// ItemIndex searches this inventory for the specified item
// and returns the pair of numbers (slot x position in slot)
// this item first appears at.
// Returns (-1, -1) if the item is not in this inventory
func (inv *Inventory) ItemIndex(item Item) (int, int) {
	var idx, jdx int
	var items []Item
	var i Item
	for idx, items = range *inv {
		for jdx, i = range items {
			if i == item {
				return idx, jdx
			}
		}
	}

	return -1, -1
}

// Send converts the inventory into a sendable string
func (inv *Inventory) Send() string {
	var out strings.Builder
	var slot int
	var items []Item

	for slot, items = range *inv {
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
