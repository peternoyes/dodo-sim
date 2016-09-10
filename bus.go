package dodosim

type Bus struct {
	Devices  []Space
	SpaceMap [0x10000]Space
}

func (bus *Bus) New() {
	bus.Devices = make([]Space, 0, 0)
}

func (bus *Bus) Add(device Space) {
	bus.Devices = append(bus.Devices, device)
}

func (bus *Bus) Start() uint16 {
	return 0
}

func (bus *Bus) Length() uint32 {
	return 0x10000
}

// Big optimization for GopherJS (only 64k total address space so map each and every address to the correct device)
func (bus *Bus) BuildMap() {
	for _, d := range bus.Devices {
		s := d.Start()
		l := d.Length()
		e := uint32(s) + l

		for i := uint32(s); i < e; i++ {
			bus.SpaceMap[i] = d
		}
	}
}

func (bus *Bus) Read(addr uint16) uint8 {
	return bus.SpaceMap[addr].Read(addr)
}

func (bus *Bus) Write(addr uint16, val uint8) {
	bus.SpaceMap[addr].Write(addr, val)
}
