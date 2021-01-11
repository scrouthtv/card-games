package doko

import "testing"

func TestUUID(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	if doko.TypeID() != DokoGameUUID {
		t.Errorf("Type id should be %d, is %d", DokoGameUUID, doko.TypeID())
	}
}
