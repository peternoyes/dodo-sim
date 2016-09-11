package dodosim

import (
	//"fmt"
	"math/rand"
	"time"
)

type Renderer interface {
	Render(data [1024]byte)
}

type Ssd1305 struct {
	Buffer        [1056]uint8
	On            bool
	Column        uint8
	Page          uint8
	Mode          uint8
	CmdInProcess  uint8
	Args          []uint8
	ArgsRemaining int
	Ram           *Ram
	Output        Renderer
}

func (s *Ssd1305) New(ram *Ram, output Renderer) {
	s.On = false
	s.Column = 0
	s.Mode = 2
	s.Ram = ram
	s.Output = output

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1056; i++ {
		s.Buffer[i] = uint8(r.Intn(256))
	}
}

func (s *Ssd1305) Start() uint16 {
	return 0x4800
}

func (s *Ssd1305) Length() uint32 {
	return 0x10
}

func (s *Ssd1305) Read(addr uint16) uint8 {
	panic("Reading from Display")
}

func (s *Ssd1305) Write(addr uint16, val uint8) {
	//fmt.Println("Addr, Val: ", addr, val)
	if addr&1 == 0 { // Command
		if s.ArgsRemaining > 0 {
			s.ArgsRemaining--
			s.Args = append(s.Args, val)

			if s.ArgsRemaining == 0 {
				if s.CmdInProcess == 0x20 {
					//fmt.Println("Address Mode: ", s.Args[0])
				} else if s.CmdInProcess == 0x21 {
					//fmt.Println("Column Address: ", s.Args[0], s.Args[1])
				} else if s.CmdInProcess == 0x22 {
					//fmt.Println("Page Address: ", s.Args[0]&3, s.Args[1]&3)
				} else if s.CmdInProcess == 0x81 {
					//fmt.Println("Contrast: ", s.Args[0])
				} else if s.CmdInProcess == 0x82 {
					//fmt.Println("Brightness: ", s.Args[0])
				} else if s.CmdInProcess == 0x91 {
					//fmt.Println("LUT: ", s.Args[0], s.Args[1], s.Args[2], s.Args[3])
				} else if s.CmdInProcess == 0x92 {

				} else if s.CmdInProcess == 0x93 {

				} else if s.CmdInProcess == 0xA8 {
					//fmt.Println("Multiplex: ", s.Args[0])
				} else if s.CmdInProcess == 0xAB {
					//fmt.Println("DIM: ", s.Args[0], s.Args[1], s.Args[2])
				} else if s.CmdInProcess == 0xAD {

				} else if s.CmdInProcess == 0xD3 {
					//fmt.Println("Display Offset: ", s.Args[0])
				} else if s.CmdInProcess == 0xD5 {
					//fmt.Println("Clock Divide: ", s.Args[0])
				}
			}

			return
		}

		s.ArgsRemaining = 0
		s.CmdInProcess = 0
		s.Args = make([]uint8, 0, 0)
		if val&0xF0 == 0 {
			low := val & 0x0F
			s.Column = (s.Column & 0xF0) | low
			//fmt.Println("Set Lower Column: ", low)
		} else if val&0xF0 == 0x10 {
			high := val & 0x0F
			s.Column = (s.Column & 0x0F) | ((high << 4) & 0xF0)
			//fmt.Println("Set Higher Column: ", high)
		} else if val == 0x20 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0x21 {
			s.CmdInProcess = val
			s.ArgsRemaining = 2
		} else if val == 0x22 {
			s.CmdInProcess = val
			s.ArgsRemaining = 2
		} else if val&0xC0 == 0x40 {
			//start := val & 0x3F
			//fmt.Println("Start Line: ", start)
		} else if val == 0x81 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0x82 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0x91 {
			s.CmdInProcess = val
			s.ArgsRemaining = 4
		} else if val == 0x92 {
			s.CmdInProcess = val
			s.ArgsRemaining = 4
		} else if val == 0x93 {
			s.CmdInProcess = val
			s.ArgsRemaining = 4
		} else if val&0xFE == 0xA0 {
			//fmt.Println("Remap: ", val&0x01)
		} else if val&0xFE == 0xA4 {
			//fmt.Println("Entire ON: ", val&0x01)
		} else if val&0xFE == 0xA6 {
			//fmt.Println("Normal: ", val&0x01)
		} else if val == 0xA8 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xAB {
			s.CmdInProcess = val
			s.ArgsRemaining = 3
		} else if val == 0xAD {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xAC {
			//fmt.Println("ON Dim")
		} else if val == 0xAE {
			//fmt.Println("Off")
		} else if val == 0xAF {
			//fmt.Println("On")
			s.Render()
		} else if val&0xF8 == 0xB0 {
			//fmt.Println("Page Start: ", val&0x07)
			s.Page = val & 0x07
		} else if val&0xF7 == 0xC0 {
			//fmt.Println("Com Dir: ", val&0x08)
		} else if val == 0xD3 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xD5 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xD8 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xD9 {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xDA {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xDB {
			s.CmdInProcess = val
			s.ArgsRemaining = 1
		} else if val == 0xE0 {
			//fmt.Println("Read Modify Write")
		} else if val == 0xE3 {
			//fmt.Println("NOP")
		} else if val == 0xEE {
			//fmt.Println("Exit Read Modify Write")
		} else if val&0xFE == 0x26 {
			s.CmdInProcess = val
			s.ArgsRemaining = 4
		} else if val&0xFC == 0x29 {
			s.CmdInProcess = val
			s.ArgsRemaining = 5
		} else if val == 0x2E {
			//fmt.Println("Deactivate Scroll")
		} else if val == 0x2F {
			//fmt.Println("Activate Scroll")
		} else if val == 0xA3 {
			s.CmdInProcess = val
			s.ArgsRemaining = 2
		}
	} else { // Data
		if s.Mode == 2 {
			i := int(s.Page)*132 + int(s.Column)
			s.Buffer[i] = val
			s.Column++

			if s.Page == 7 && s.Column == 128 {
				s.Render()
			}

		} else {
			panic("Unsupported Mode")
		}
	}
}

func (s *Ssd1305) Render() {
	var data [1024]byte
	var x, y int
	for y = 0; y < 8; y++ {
		for x = 0; x < 128; x++ {
			b := s.Buffer[(x+2)+y*132]
			data[x+y*128] = b
		}
	}
	s.Output.Render(data)
}
