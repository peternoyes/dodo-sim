package dodosim

import (
	"fmt"
)

type Rom [0x8000]uint8

func (rom *Rom) Start() uint16 {
	return 0x8000
}

func (rom *Rom) Length() uint32 {
	return 0x8000
}

func (rom *Rom) Read(addr uint16) uint8 {
	v := rom[addr-0x8000]
	//fmt.Printf("Reading: %v from %v in ROM\n", v, addr)
	return v
}

func (rom *Rom) Write(addr uint16, val uint8) {
	fmt.Println("Address: ", addr)
	panic("ROM is readonly")
}
