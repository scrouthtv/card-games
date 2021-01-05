package main

// Inventory is a collection of int-identified slots which can hold multiple cards
type Inventory map[int]*Deck

// NewInventory creates a new inventory by placing
// all provided items in the first (0th) slot
func NewInventory(items []*Card) *Inventory {
	var inv Inventory = make(Inventory)

	inv.Get(0).AddAll(items...)

	return &inv
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

// Length returns how many defined slots this inventory has,
// they could also be empty
func (inv *Inventory) Length() int {
	return len(*inv)
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
