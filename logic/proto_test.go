package logic

import "testing"

func TestSimplePacket(t *testing.T) {
	var p *Packet = NewPacket("play the card")

	if p.Action() != "play" {
		t.Errorf("Packet action should be play, it is %s", p.Action())
	}

	var sbargs []string = []string{"the", "card"}
	var isargs []string = p.Args()

	if len(sbargs) != len(isargs) {
		t.Errorf("Wrong amount of args returned, got %d, should be %d",
			len(isargs), len(sbargs))
	}

	var i int
	var s string
	for i, s = range isargs {
		if s != sbargs[i] {
			t.Errorf("Wrong arg @ %d, should be %s, is %s",
				i, s, isargs[i])
		}
	}
}

func TestNilPacket(t *testing.T) {
	var p *Packet = nil

	var cmd string = p.Action()
	var args []string = p.Args()

	if cmd != "" {
		t.Errorf("Wrong action returned: \"%s\", should be \"\"", cmd)
	}

	if len(args) != 0 {
		t.Errorf("Wrong args returned: %v", args)
	}
}
