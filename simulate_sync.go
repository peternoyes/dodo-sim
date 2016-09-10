package dodosim

import (
	"strings"
)

type SimulatorSync struct {
	Renderer       Renderer
	CyclesPerFrame func(cycles uint64)

	Bus     *Bus
	Cpu     *Cpu
	Gamepad *Gamepad
	Resolve *Resolve
	Fram    *Fram

	Cycles              uint64
	LastOp              uint8
	WaitTester          int
	WaitingForInterrupt bool
	MissedFrames        int
}

func (s *SimulatorSync) PumpClock(input string) {
	s.Gamepad.A = strings.Contains(input, "A")
	s.Gamepad.B = strings.Contains(input, "B")
	s.Gamepad.U = strings.Contains(input, "U")
	s.Gamepad.D = strings.Contains(input, "D")
	s.Gamepad.L = strings.Contains(input, "L")
	s.Gamepad.R = strings.Contains(input, "R")

	for {
		opcode := s.Bus.Read(s.Cpu.PC)

		s.Cpu.PC++
		s.Cpu.Status |= Constant
		o := GetOperation(opcode)

		c := o.Cycles

		s.Resolve.Opcode = opcode
		s.Resolve.Mode = o.Mode

		s.Resolve.Resolve()

		pop, rc := o.Handler(s.Resolve)
		c += rc
		if s.Resolve.Penalty && pop {
			c += 1
		}

		s.Cycles += uint64(c)

		if (s.LastOp == 0xA5 && opcode == 0xF0) || (s.LastOp == 0xF0 && opcode == 0xA5) {
			s.WaitTester++
		} else {
			s.WaitTester = 0
		}

		s.LastOp = opcode

		// If Lda, Beq sequences happens 5 times in a row then assume we are waiting for interrupt
		if s.WaitTester == 10 { // Getting in here tells us that a complete game cycle was performed
			s.CyclesPerFrame(uint64(s.MissedFrames*50000) + s.Cycles)
			s.MissedFrames = 0
			s.Cycles = 0
			s.WaitTester = 0
			s.Cpu.Irq(s.Bus)

			return
		} else if s.Cycles >= 50000 { // If we hit 50000 cycles then that means we need to pause for a whole additional interrupt cycle
			s.Cycles = 0
			s.MissedFrames++

			s.Cpu.Irq(s.Bus)

			return
		}
	}
}

func (s *SimulatorSync) SwitchFram(game []byte) {
	s.Fram.New(game)
	s.Cpu.Reset(s.Bus)
}

func (s *SimulatorSync) SimulateSyncInit(firmware, game []byte) {
	s.Bus = new(Bus)
	s.Bus.New()

	s.Resolve = new(Resolve)

	ram := new(Ram)
	s.Bus.Add(ram)

	rom := new(Rom)

	ssd1305 := new(Ssd1305)
	ssd1305.New(ram, s.Renderer)

	s.Bus.Add(ssd1305)

	s.Gamepad = new(Gamepad)
	s.Gamepad.New()

	s.Fram = new(Fram)
	s.Fram.New(game)

	via := new(Via)
	via.New(s.Gamepad, s.Fram)
	s.Bus.Add(via)

	acia := new(Acia)
	s.Bus.Add(acia)

	for i, b := range firmware {
		rom[i] = b
	}

	s.Bus.Add(rom)

	s.Bus.BuildMap()

	s.Cpu = new(Cpu)
	s.Cpu.Reset(s.Bus)

	s.Resolve.Cpu = s.Cpu
	s.Resolve.Space = s.Bus

	BuildTable()

	s.Cycles = 0
	s.LastOp = 0
	s.WaitTester = 0
	s.WaitingForInterrupt = false
	s.MissedFrames = 0
}
