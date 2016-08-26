package dodosim

type Gamepad struct {
	U   bool
	D   bool
	L   bool
	R   bool
	A   bool
	B   bool
	LED bool
}

func (g *Gamepad) New() {
	g.U = false
	g.D = false
	g.L = false
	g.R = false
	g.A = false
	g.B = false
	g.LED = false
}

func (g *Gamepad) ReadBit(bit int) bool {
	switch bit {
	case 0:
		return !g.U
	case 1:
		return !g.D
	case 2:
		return !g.L
	case 3:
		return !g.R
	case 4:
		return !g.A
	case 5:
		return !g.B
	default:
		return true
	}

}

func (g *Gamepad) WriteBit(bit int, val bool) {
	switch bit {
	case 6:
		g.LED = val
	}
}
