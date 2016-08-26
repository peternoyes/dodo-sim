package dodosim

type Space interface {
	Start() uint16
	Length() uint32
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}
