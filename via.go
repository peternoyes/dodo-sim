package dodosim

const (
	PORTB uint16 = 0x6000
	PORTA uint16 = 0x6001
	DDRB  uint16 = 0x6002
	DDRA  uint16 = 0x6003
	T1CL  uint16 = 0x6004
	T1CH  uint16 = 0x6005
	T1LL  uint16 = 0x6006
	T1LH  uint16 = 0x6007
	T2CL  uint16 = 0x6008
	T2CH  uint16 = 0x6009
	SR    uint16 = 0x600A
	ACR   uint16 = 0x600B
	PCR   uint16 = 0x600C
	IFR   uint16 = 0x600D
	IER   uint16 = 0x600E
	ORAX  uint16 = 0x600F
)

const (
	SR_DISABLED    uint8 = 0x0
	SR_IN_T2       uint8 = 0x1
	SR_IN_PHI2     uint8 = 0x2
	SR_IN_EXT      uint8 = 0x3
	SR_OUT_T2_FREE uint8 = 0x4
	SR_OUT_T2      uint8 = 0x5
	SR_OUT_PHI2    uint8 = 0x6
	SR_OUT_EXT     uint8 = 0x7
)

type Via struct {
	PortA   Parallel
	PortB   Parallel
	Speaker Speaker

	DirA uint8
	DirB uint8
	ACR  uint8
	T2CL uint8
	SR   uint8
}

type Parallel interface {
	ReadBit(bit int) bool
	WriteBit(bit int, val bool)
}

type Speaker interface {
	SetFrequency(freq int)
}

func (v *Via) New(portA, portB Parallel, speaker Speaker) {
	v.PortA = portA
	v.PortB = portB
	v.Speaker = speaker

	v.DirA = 0xFF
	v.DirB = 0xFF
	v.ACR = 0x0
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
	case ACR:
		val = v.ACR
	case T2CL:
		val = v.T2CL
	case SR:
		val = v.SR
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
	case ACR:
		v.ACR = val
	case T2CL:
		v.T2CL = val
		v.processAudio()
	case SR:
		v.SR = val
	}
}

func (v *Via) processAudio() {
	if v.Speaker != nil && v.getSRMode() == SR_OUT_T2_FREE && v.SR == 15 {
		freq := 0
		if v.T2CL != 0 {
			freq = 1000000 / ((int(v.T2CL) + 2) * 8)
		}

		v.Speaker.SetFrequency(freq)
	}
}

func (v *Via) getSRMode() uint8 {
	return (v.ACR >> 2) & 0x7
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
