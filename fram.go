package dodosim

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	CLK  int = 0
	CS   int = 1
	MOSI int = 6
	MISO int = 7
)

const (
	WREN  uint8 = 0x06
	WRDI  uint8 = 0x04
	RDSR  uint8 = 0x05
	WRSR  uint8 = 0x01
	READ  uint8 = 0x03
	WRITE uint8 = 0x02
	RDID  uint8 = 0x9F
)

const (
	None int = iota
	ReadAddr1
	ReadAddr2
	WriteAddr1
	WriteAddr2
	Reading
	Writing
)

type Fram struct {
	Data         [0x2000]uint8
	Clock        bool
	Off          bool
	WaitingMosi  bool
	ByteIn       uint8
	ByteOut      uint8
	Bit          uint8
	BufferOut    []uint8
	BufferOutPos int
	Address      uint16
	State        int
	WriteEnable  bool
}

func (f *Fram) New() {
	f.Clock = false
	f.Off = true
	f.WaitingMosi = false
	f.ByteIn = 0
	f.ByteOut = 0
	f.Bit = 0
	f.BufferOut = nil
	f.BufferOutPos = 0
	f.Address = 0
	f.State = None
	f.WriteEnable = false

	if _, err := os.Stat("fram.bin"); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile("fram.bin")
		if err != nil {
			panic(err)
		}

		for i, v := range b {
			f.Data[i] = v
		}
	}
}

func (f *Fram) Flush() {
	ioutil.WriteFile("fram.bin", f.Data[:], 0644)
}

func (f *Fram) ReadBit(bit int) bool {
	if bit == MISO {
		return (f.ByteOut<<f.Bit)&0x80 == 0x80
	}
	return true
}

func (f *Fram) WriteBit(bit int, val bool) {
	switch bit {
	case CLK:
		if !f.Off {
			if val != f.Clock {
				f.Clock = val
				if f.Clock { // Rising edge of clock
					var b uint8 = 0
					if f.WaitingMosi {
						b = 0x80
					}
					b = b >> f.Bit
					f.ByteIn |= b
				} else { // Falling Edge
					f.Bit++
					if f.Bit == 8 {
						f.Bit = 0
						t := f.ByteIn
						f.ByteIn = 0
						f.ByteOut = 0
						f.processByte(t)
					}
				}
			}
		}
		break
	case CS:
		if val != f.Off {
			f.Off = val
			if f.Off {
				f.State = None
				f.Address = 0
			}
		}
	case MOSI:
		if !f.Off {
			f.WaitingMosi = val
		}
	}
}

func (f *Fram) processByte(val uint8) {
	switch f.State {
	case ReadAddr1:
		f.Address = uint16(val) << 8
		f.State = ReadAddr2
	case ReadAddr2:
		f.Address |= uint16(val) & 0xFF
		f.State = Reading
		f.ByteOut = f.Data[f.Address]
		f.incAddress()
	case WriteAddr1:
		f.Address = uint16(val) << 8
		f.State = WriteAddr2
	case WriteAddr2:
		f.Address |= uint16(val) & 0xFF
		f.State = Writing
	case Reading:
		f.ByteOut = f.Data[f.Address]
		f.incAddress()
	case Writing:
		if f.WriteEnable {
			f.Data[f.Address] = val
		}
		f.incAddress()
	case None:
		if f.BufferOut != nil {
			f.ByteOut = f.BufferOut[f.BufferOutPos]
			f.BufferOutPos++
			if f.BufferOutPos >= len(f.BufferOut) {
				f.BufferOut = nil
				f.BufferOutPos = 0
			}

			return
		}
		switch val {
		case 0:
			break
		case WREN:
			f.WriteEnable = true
		case WRDI:
			f.WriteEnable = false
		case READ:
			f.State = ReadAddr1
		case WRITE:
			f.State = WriteAddr1
		case RDID:
			f.ByteOut = 0x04
			f.BufferOut = make([]uint8, 3)
			f.BufferOutPos = 0
			f.BufferOut[0] = 0x7F
			f.BufferOut[1] = 0x03
			f.BufferOut[2] = 0x02
		default:
			fmt.Println("Val: ", val)
			panic("Yikes")
		}
	}
}

func (f *Fram) incAddress() {
	f.Address++
	if f.Address == 0x2000 {
		f.Address = 0
	}
}
