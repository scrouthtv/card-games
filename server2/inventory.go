package main

import (
	"fmt"
	"strings"
)

// Inventory is a collection of int-identified slots which can hold multiple cards
type Inventory map[int][]*Card

const idelim string = "-"
const odelim string = "#"

// NewInventory creates a new inventory by placing
// all provided items in the first (0th) slot
func NewInventory(items []*Card) *Inventory {
	var inv Inventory = make(Inventory)

	inv[0] = items

	return &inv
}

// RemoveItem removes a single item from this inventory
// while keeping everything in order
// It searches through the first slot first, afterwards the second, and so on
// The first occurence is removed
// Returns whether the item was removed ( = whether it was found)
func (inv *Inventory) RemoveItem(itm *Card) bool {
	var idx, jdx int
	var oldslot []*Card
	var item *Card
	for idx, oldslot = range *inv {
		for jdx, item = range oldslot {
			if *item == *itm {
				// remove this card:
				(*inv)[idx] = append(oldslot[:jdx], oldslot[jdx+1:]...)
				return true
			}
		}
	}
	return false
}

// Get returns the item stack at the specified slot
func (inv *Inventory) Get(slot int) []*Card {
	return (*inv)[slot]
}

// AddToSlot adds the specified item(s) to this inventory at the
// given slot. If there are no cards at this slot yet, it is created.
func (inv *Inventory) AddToSlot(slot int, items ...*Card) {
	var oldslot []*Card = (*inv)[slot]
	(*inv)[slot] = append(oldslot, items...)
}

// Send converts the inventory into a sendable string
func (inv *Inventory) Send() string {
	var out strings.Builder
	var slot int
	var items []*Card

	for slot, items = range *inv {
		if out.Len() > 0 {
			out.WriteString(odelim)
		}
		out.WriteString(fmt.Sprintf("%d:", slot))
		var i int
		var item *Card
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
