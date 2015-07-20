package packet

type lcdMatrix struct {
	dataR [3][16]uint32
	dataG [3][16]uint32
}

type Packet struct {
	Header     string
	Coord      string
	DataR      []byte
	DataG      []byte
	Terminator string
}

func createTestPacket() *Packet {
	var p Packet
	p.Header = "pcmat\r"
	p.Coord = "000\r5ff\r"
	s1 := ""
	s2 := ""
	for i := 0; i < 192; i++ {
		s1 += "p"
		s2 += "a"
	}
	s1 += "\r"
	s2 += "\r"
	p.DataR = []byte(s1)
	p.DataG = []byte(s2)
	p.Terminator = "end\r"

	return &p
}
