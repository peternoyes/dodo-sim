package dodosim

const (
	PORTB uint16 = 0x6000
	PORTA uint16 = 0x6001
	DDRB  uint16 = 0x6002
	DDRA  uint16 = 0x6003
)

type Via struct {
	PortA Parallel
	PortB Parallel
	DirA  uint8
	DirB  uint8
}

type Parallel interface {
	ReadBit(bit int) bool
	WriteBit(bit int, val bool)
}

func (v *Via) New(portA, portB Parallel) {
	v.PortA = portA
	v.PortB = portB
	v.DirA = 0xFF
	v.DirB = 0xFF
}

func (v *Via) Start() uint16 {
	return 0x6000
}

func (v *Via) Length() uint32 {
	return 0x10
}

func (v *Via) Read(addr uint16) uint8 {
	val := uint8(0xFF)

	switch addr {
	case PORTB:
		val = readParallel(v.PortB, v.DirB)
	case PORTA:
		val = readParallel(v.PortA, v.DirA)
	case DDRB:
		val = v.DirB
	case DDRA:
		val = v.DirA
	}

	return val
}

func (v *Via) Write(addr uint16, val uint8) {
	switch addr {
	case PORTB:
		writeParallel(v.PortB, v.DirB, val)
	case PORTA:
		writeParallel(v.PortA, v.DirA, val)
	case DDRB:
		v.DirB = val
	case DDRA:
		v.DirA = val
	}
}

func readParallel(p Parallel, d uint8) uint8 {
	var v uint8 = 0x00
	if p != nil {
		for i := 7; i >= 0; i-- {
			v = v << 1
			if (d&0x80 == 0x0) && p.ReadBit(i) {
				v |= 0x1

			}

			d = d << 1
		}
	}
	return v
}

func writeParallel(p Parallel, d, v uint8) {
	if p != nil {
		for i := 0; i < 8; i++ {
			if d&0x1 == 0x1 {
				p.WriteBit(i, v&0x1 == 0x1)
			}

			v = v >> 1
			d = d >> 1
		}
	}
}
