package main

// Inventory is a collection of int-identified slots which can hold multiple cards
type Inventory map[int]*Deck

const idelim string = "-"
const odelim string = "#"

// NewInventory creates a new inventory by placing
// all provided items in the first (0th) slot
func NewInventory(items []*Card) *Inventory {
	var inv Inventory = make(Inventory)

	inv.Get(0).AddAll(items...)

	return &inv
}

// RemoveItem removes an item from this inventory
// while keeping everything in order
// It searches through the first slot first, afterwards the second, and so on
// The first n occurences is removed
// Returns how many items were removed
func (inv *Inventory) RemoveItem(card Card, n int) int {
	var deleted int = 0
	var deck *Deck
	for _, deck = range *inv {
		deleted += deck.Remove(card, n-deleted)
		if deleted == n {
			return deleted
		}
	}
	return deleted
}

// Get returns the item stack at the specified slot
func (inv *Inventory) Get(slot int) *Deck {
	if (*inv)[slot] == nil {
		var deck Deck = make(Deck, 0)
		return &deck
	}
	return (*inv)[slot]
}

// AddToSlot adds the specified item(s) to this inventory at the
// given slot. If there are no cards at this slot yet, it is created.
func (inv *Inventory) AddToSlot(slot int, items ...*Card) {
	inv.Get(slot).AddAll(items...)
}

// Serialize converts the inventory into a sendable string
func (inv *Inventory) Serialize() string {
	panic("not impl")
}

// Clear clears the specified slot, e. g. discards any cards in it
func (inv *Inventory) Clear(slot int) {
	(*inv)[slot] = nil
}

// ClearAll clears all slots by calling Clear() on them
func (inv *Inventory) ClearAll() {
	var i int
	for i = 0; i < len(*inv); i++ {
		inv.Clear(i)
	}
}
