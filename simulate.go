package dodosim

import (
	"strings"
	"time"
)

type Simulator struct {
	Renderer         Renderer
	Input            chan string
	Ticker           *time.Ticker
	IntervalCallback func() bool
	Complete         func(*Cpu)
	CyclesPerFrame   func(cycles uint64)
}

func Simulate(s *Simulator, firmware, game []byte) {
	bus := new(Bus)
	bus.New()

	ram := new(Ram)
	bus.Add(ram)

	rom := new(Rom)

	ssd1305 := new(Ssd1305)
	ssd1305.New(ram, s.Renderer)

	bus.Add(ssd1305)

	gamepad := new(Gamepad)
	gamepad.New()

	fram := new(Fram)
	fram.New(game)

	via := new(Via)
	via.New(gamepad, fram)
	bus.Add(via)

	acia := new(Acia)
	bus.Add(acia)

	for i, b := range firmware {
		rom[i] = b
	}
	bus.Add(rom)

	bus.BuildMap()

	cpu := new(Cpu)
	cpu.Reset(bus)

	BuildTable()

	var cycles uint64 = 0

	syncer := make(chan int)
	go func(syncer chan int) {
		for {
			time.Sleep(50 * time.Millisecond) // 50 ms
			syncer <- 50
		}
	}(syncer)

	var lastOp uint8 = 0
	var waitTester int = 0
	var waitingForInterrupt bool = false
	var missedFrames int = 0

L:
	for {
		opcode := bus.Read(cpu.PC)

		cpu.PC++
		cpu.Status |= Constant
		o := GetOperation(opcode)
		c := o.Execute(cpu, bus, opcode)
		cycles += uint64(c)

		if (lastOp == 0xA5 && opcode == 0xF0) || (lastOp == 0xF0 && opcode == 0xA5) {
			waitTester++
		} else {
			waitTester = 0
		}

		// If Lda, Beq sequences happens 5 times in a row then assume we are waiting for interrupt
		if waitTester == 10 { // Getting in here tells us that a complete game cycle was performed
			s.CyclesPerFrame(uint64(missedFrames*50000) + cycles)
			missedFrames = 0
			cycles = 0
			waitTester = 0
			waitingForInterrupt = true
		} else if cycles >= 50000 { // If we hit 50000 cycles then that means we need to pause for a whole additional interrupt cycle
			cycles = 0
			waitingForInterrupt = true
			missedFrames++
		}

		if waitingForInterrupt {
			<-syncer
			cpu.Irq(bus)
			waitingForInterrupt = false
		}

		lastOp = opcode

		select {
		case v, ok := <-s.Input:
			if !ok {
				break
			} else {
				gamepad.A = strings.Contains(v, "A")
				gamepad.B = strings.Contains(v, "B")
				gamepad.U = strings.Contains(v, "U")
				gamepad.D = strings.Contains(v, "D")
				gamepad.L = strings.Contains(v, "L")
				gamepad.R = strings.Contains(v, "R")

				if strings.Contains(v, "X") {
					s.Complete(cpu)
					return
				}
			}
		case <-s.Ticker.C:
			if !s.IntervalCallback() {
				break L
			}
		default:
		}
	}
}
