package dodosim

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
)

type Bus struct {
	Devices []Space
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

func (bus *Bus) Read(addr uint16) uint8 {
	for _, d := range bus.Devices {
		s := d.Start()
		l := d.Length()
		e := uint32(s) + l

		if addr >= s && uint32(addr) < e {
			return d.Read(addr)
		}
	}

	cmd := exec.Command("/bin/stty", "-raw", "echo")
	cmd.Stdin = os.Stdin
	cmd.Run()

	fmt.Println("Attempting to read from: ", addr)
	debug.PrintStack()
	panic("Unmapped Address Space")
	return 0
}

func (bus *Bus) Write(addr uint16, val uint8) {
	for _, d := range bus.Devices {
		s := d.Start()
		l := d.Length()
		e := uint32(s) + l

		if addr >= s && uint32(addr) < e {
			d.Write(addr, val)
			return
		}
	}

	fmt.Println("Attempting to write to: ", addr)
	panic("Unmapped Address Space")
}
