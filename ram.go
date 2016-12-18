package dodosim

type Ram [0x7F00]uint8

func (ram *Ram) Start() uint16 {
	return 0x0
}

func (ram *Ram) Length() uint32 {
	return 0x7F00
}

func (ram *Ram) Read(addr uint16) uint8 {
	return ram[addr]
}

func (ram *Ram) Write(addr uint16, val uint8) {
	ram[addr] = val
}
