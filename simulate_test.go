package dodosim

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type Ramtest [0x10000]uint8

func (ram *Ramtest) Start() uint16 {
	return 0x0
}

func (ram *Ramtest) Length() uint32 {
	return 0x10000
}

func (ram *Ramtest) Read(addr uint16) uint8 {
	return ram[addr]
}

func (ram *Ramtest) Write(addr uint16, val uint8) {
	ram[addr] = val
}

func TestSimulator(t *testing.T) {
	bus := new(dodosim.Bus)
	bus.New()

	ram := new(Ramtest)
	bus.Add(ram)

	dat, err := ioutil.ReadFile("6502_functional_test.bin")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, b := range dat {
		ram[i] = b
	}

	cpu := new(dodosim.Cpu)
	cpu.Reset(bus)

	cpu.PC = 0x400

	dodosim.BuildTable()

	for {
		before := cpu.PC
		opcode := bus.Read(cpu.PC)

		cpu.PC++
		cpu.Status |= dodosim.Constant
		o := dodosim.GetOperation(opcode)
		o.Execute(cpu, bus, opcode)

		//fmt.Println(opcode)

		if before == cpu.PC {
			if cpu.PC != 13209 {
				t.Error("Failure. Trap at ", cpu.PC)
			}
			return
		}
	}

	t.Error("End Condition")
}
