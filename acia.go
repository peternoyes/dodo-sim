package dodosim

type Acia struct {
}

func (a *Acia) Start() uint16 {
	return 0x7F10
}

func (a *Acia) Length() uint32 {
	return 0x10
}

func (a *Acia) Read(addr uint16) uint8 {
	panic("Reading from Acia")
}

func (a *Acia) Write(addr uint16, val uint8) {

}
