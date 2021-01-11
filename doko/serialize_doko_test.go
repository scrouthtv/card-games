package doko

import "testing"
import "bytes"
import "github.com/scrouthtv/card-games/logic"

func TestSerializePreparingGame(t *testing.T) {
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var d *Doko = NewDoko(gs)
	var buf bytes.Buffer
	d.Reset()
	d.WriteBinary(0, &buf)
	testBinary(t, &buf, []byte{logic.StatePreparing})

	gs.state = 5
	buf.Reset()
	d.WriteBinary(0, &buf)
	testBinary(t, &buf, []byte{255})
}

func TestIntArrayBinary(t *testing.T) {
	var arr []int = []int{5, 0, 13, 255, 0, 0, 255}
	var buf bytes.Buffer
	writeIntArray(arr, &buf)
	testBinary(t, &buf, []byte{7, 5, 0, 13, 255, 0, 0, 255})
}

func testBinary(t *testing.T, buf *bytes.Buffer, exp []byte) {
	if t != nil {
		t.Helper()
	}

	if buf.Len() < len(exp) {
		t.Errorf("Wrong binary, is %d long, should be %d long", buf.Len(), len(exp))
	}

	var is []byte = buf.Bytes()

	var i int
	var b1, b2 byte
	for i, b1 = range exp {
		b2 = is[i]
		if b1 != b2 {
			t.Errorf("Wrong byte at %d, should be %d, is %d", i, b1, b2)
		}
	}
}

func TestUUID(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	if doko.TypeID() != DokoGameUUID {
		t.Errorf("Type id should be %d, is %d", DokoGameUUID, doko.TypeID())
	}
}
