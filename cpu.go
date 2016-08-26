package dodosim

import (
//"fmt"
)

type Cpu struct {
	SP     uint8
	A      uint8
	X      uint8
	Y      uint8
	Status Status
	PC     uint16
}

type Resolve struct {
	Cpu     *Cpu
	Space   Space
	Mode    AddressMode
	Address uint16
	Penalty bool
	Opcode  uint8
}

type Status uint8

const (
	Carry     Status = 0x01
	Zero      Status = 0x02
	Interrupt Status = 0x04
	Decimal   Status = 0x08
	Break     Status = 0x10
	Constant  Status = 0x20
	Overflow  Status = 0x40
	Sign      Status = 0x80
)

type AddressMode int

const (
	Imp AddressMode = iota
	Acc
	Imm
	Zp
	Zpx
	Zpy
	Rel
	Abso
	Absx
	Absy
	Ind
	Indx
	Indy
	Indzp
)

func (cpu *Cpu) Reset(space Space) {
	cpu.PC = uint16(space.Read(0xFFFC)) | (uint16(space.Read(0xFFFD)) << 8)
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.SP = 0xFD
	cpu.Status |= Constant
}

func (cpu *Cpu) SetCarry() {
	cpu.Status |= Carry
}

func (cpu *Cpu) ClearCarry() {
	cpu.Status &^= Carry
}

func (cpu *Cpu) SetZero() {
	cpu.Status |= Zero
}

func (cpu *Cpu) ClearZero() {
	cpu.Status &^= Zero
}

func (cpu *Cpu) SetInterrupt() {
	cpu.Status |= Interrupt
}

func (cpu *Cpu) ClearInterrupt() {
	cpu.Status &^= Interrupt
}

func (cpu *Cpu) SetDecimal() {
	cpu.Status |= Decimal
}

func (cpu *Cpu) ClearDecimal() {
	cpu.Status &^= Decimal
}

func (cpu *Cpu) SetOverflow() {
	cpu.Status |= Overflow
}

func (cpu *Cpu) ClearOverflow() {
	cpu.Status &^= Overflow
}

func (cpu *Cpu) SetSign() {
	cpu.Status |= Sign
}

func (cpu *Cpu) ClearSign() {
	cpu.Status &^= Sign
}

func (cpu *Cpu) ZeroCalc8(val uint8) {
	if val&0xFF != 0 {
		cpu.ClearZero()
	} else {
		cpu.SetZero()
	}
}

func (cpu *Cpu) ZeroCalc(val uint16) {
	if val&0x00FF != 0 {
		cpu.ClearZero()
	} else {
		cpu.SetZero()
	}
}

func (cpu *Cpu) SignCalc8(val uint8) {
	if val&0x80 != 0 {
		cpu.SetSign()
	} else {
		cpu.ClearSign()
	}
}

func (cpu *Cpu) SignCalc(val uint16) {
	if val&0x0080 != 0 {
		cpu.SetSign()
	} else {
		cpu.ClearSign()
	}
}

func (cpu *Cpu) CarryCalc(val uint16) {
	if val&0xFF00 != 0 {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}
}

func (cpu *Cpu) OverflowCalc(val uint16, m uint8, o uint16) {
	if ((val ^ uint16(m)) & (val ^ o) & 0x0080) != 0 {
		cpu.SetOverflow()
	} else {
		cpu.ClearOverflow()
	}
}

func (cpu *Cpu) SaveAccum(val uint16) {
	cpu.A = uint8(val & 0x00FF)
}

func (cpu *Cpu) Irq(space Space) {
	if cpu.Status&Interrupt != 0 {
		return
	}

	space.Write(0x100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
	space.Write(0x100+((uint16(cpu.SP)-1)&0x00FF), uint8(cpu.PC&0x00FF))
	cpu.SP -= 2

	space.Write(0x100+uint16(cpu.SP), uint8(cpu.Status))
	cpu.SP -= 1

	cpu.Status |= Interrupt
	cpu.PC = uint16(space.Read(0xFFFE)) | (uint16(space.Read(0xFFFF)) << 8)
}

func (r Resolve) Push16(val uint16) {
	cpu := r.Cpu
	r.Space.Write(0x100+uint16(cpu.SP), uint8((val>>8)&0x00FF))
	r.Space.Write(0x100+((uint16(cpu.SP)-1)&0x00FF), uint8(val&0x00FF))
	cpu.SP -= 2
}

func (r Resolve) Push8(val uint8) {
	cpu := r.Cpu
	r.Space.Write(0x100+uint16(cpu.SP), val)
	cpu.SP -= 1
}

func (r Resolve) Pull16() uint16 {
	cpu := r.Cpu
	var t uint16
	t = uint16(r.Space.Read(0x100+((uint16(cpu.SP)+1)&0x00FF))) | (uint16(r.Space.Read(0x100+((uint16(cpu.SP)+2)&0x00FF))) << 8)
	cpu.SP += 2
	return t
}

func (r Resolve) Pull8() uint8 {
	cpu := r.Cpu
	cpu.SP += 1
	return r.Space.Read(0x100 + uint16(cpu.SP)&0x00FF)
}

func (r Resolve) Write(val uint16) {
	if r.Mode == Acc {
		r.Cpu.A = uint8(val & 0x00FF)
	} else {
		r.Space.Write(r.Address, uint8(val&0x00FF))
	}
}

func (r Resolve) Read() uint16 {
	if r.Mode == Acc {
		return uint16(r.Cpu.A)
	} else {
		return uint16(r.Space.Read(r.Address))
	}
}

func (a AddressMode) Resolve(cpu *Cpu, space Space, opcode uint8) Resolve {
	pc := cpu.PC
	var r uint16 = 0
	penalty := false
	switch a {
	case Imp:
		return Resolve{cpu, space, a, r, penalty, opcode}
	case Acc:
		return Resolve{cpu, space, a, r, penalty, opcode}
	case Imm:
		r = pc
		cpu.PC = pc + 1
	case Zp:
		r = uint16(space.Read(pc))
		cpu.PC = pc + 1
	case Zpx:
		r = (uint16(space.Read(pc)) + uint16(cpu.X)) & 0xFF
		cpu.PC = pc + 1
	case Zpy:
		r = (uint16(space.Read(pc)) + uint16(cpu.Y)) & 0xFF
		cpu.PC = pc + 1
	case Rel:
		r = uint16(space.Read(cpu.PC))
		cpu.PC = pc + 1
		if (r & 0x80) != 0 {
			r |= 0xFF00
		}
	case Abso:
		r = uint16(space.Read(pc)) | (uint16(space.Read(pc+1)) << 8)
		cpu.PC = pc + 2
	case Absx:
		var startPage uint16
		r = uint16(space.Read(pc)) | (uint16(space.Read(pc+1)) << 8)
		startPage = r & 0xFF00
		r += uint16(cpu.X)
		if startPage != (r & 0xFF00) {
			penalty = true // 1 cycle CPU penalty
		}
		cpu.PC = pc + 2
	case Absy:
		var startPage uint16
		r = uint16(space.Read(pc)) | (uint16(space.Read(pc+1)) << 8)
		startPage = r & 0xFF00
		r += uint16(cpu.Y)
		if startPage != (r & 0xFF00) {
			penalty = true
		}
		cpu.PC = pc + 2
	case Ind:
		var h1, h2 uint16
		h1 = uint16(space.Read(pc)) | (uint16(space.Read(pc+1)) << 8)
		h2 = (h1 & 0xFF00) | ((h1 + 1) & 0x00FF) // Replicating 6502 bug
		r = uint16(space.Read(h1)) | (uint16(space.Read(h2)) << 8)
		cpu.PC = pc + 2
	case Indx:
		var h uint16
		h = (uint16(space.Read(pc)) + uint16(cpu.X)) & 0xFF
		r = uint16(space.Read(h&0x00FF)) | (uint16(space.Read((h+1)&0x00FF)) << 8)
		cpu.PC = pc + 1
	case Indy:
		var h1, h2, startPage uint16
		h1 = uint16(space.Read(pc))
		h2 = (h1 & 0xFF00) | ((h1 + 1) & 0x00FF)
		r = uint16(space.Read(h1)) | (uint16(space.Read(h2)) << 8)
		startPage = r & 0xFF00
		r += uint16(cpu.Y)
		if startPage != (r & 0xFF00) {
			penalty = true
		}
		cpu.PC = pc + 1
	case Indzp:
		var h1, h2 uint16
		h1 = uint16(space.Read(pc))
		h2 = (h1 & 0xFF00) | ((h1 + 1) & 0x00FF)
		r = uint16(space.Read(h1)) | (uint16(space.Read(h2)) << 8)
		cpu.PC = pc + 1
	}
	return Resolve{cpu, space, a, r, penalty, opcode}
}
