package packet

import "github.com/RobotClubKut/MatrixLEDTOOLS/go/matrix"

type LcdMatrix struct {
	DataR [3][16]uint32
	DataG [3][16]uint32
}

type Packet struct {
	Header     string
	Coord      string
	DataR      []byte
	DataG      []byte
	Terminator string
}

func CreateTestPacket() *Packet {
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

func printMatrix(str matrix.MatrixString, shift int) *LcdMatrix {
	var ret LcdMatrix

	for y := 0; y < 16; y++ {
		var bufR []uint8
		var bufG []uint8
		for cn := 0; cn < len(str.Char); cn++ {
			for x := 0; x < 16; x++ {
				//bit := str.c[cn].bitmap[y]
				bit := matrix.ReadMatrixChar(str.Char[cn]).Bitmap[y]
				bit = bit >> (15 - uint(x))
				//mask 0b0000000000000001
				bit = bit & 0x0001
				//赤があるとき
				//buf[coordX] = uint8(bit)
				if str.Char[cn].Color&0xff0000 == 0xff0000 {
					bufR = append(bufR, uint8(bit))
				} else {
					bufR = append(bufR, uint8(0))
				}
				if str.Char[cn].Color&0x00ff00 == 0x00ff00 {
					bufG = append(bufG, uint8(bit))
				} else {
					bufG = append(bufG, uint8(0))
				}
			}
		}
		i := 0
		counter := 0

		for x := shift; x < 96+shift; x++ {
			if counter == 32 {
				i++
				if i == 3 {
					break
				}
				counter = 0
			}
			ret.DataR[i][y] = ret.DataR[i][y] << 1
			ret.DataG[i][y] = ret.DataG[i][y] << 1
			if int(str.Coord)+x < len(bufR) && x+int(str.Coord) >= 0 {
				ret.DataR[i][y] |= uint32(bufR[x+int(str.Coord)])
			} else {
				ret.DataR[i][y] |= 0
			}
			if int(str.Coord)+x < len(bufG) && int(str.Coord)+x >= 0 {
				ret.DataG[i][y] |= uint32(bufG[x+int(str.Coord)])
			} else {
				ret.DataG[i][y] |= 0
			}
			counter++
		}
	}

	return &ret
}

func CreatePacket(str matrix.MatrixString, shift int) *Packet {
	data := printMatrix(str, shift)
	//str.coord += shift

	var packet Packet
	packet.Header = "pcmat\r"
	packet.Coord = "000\r5ff\r"

	var bufr []byte
	var bufg []byte
	for y := 0; y < 16; y++ {

		for x := 0; x < 3; x++ {
			r := data.DataR[x][y]
			rbyte := (r & 0xff000000) >> 24
			bufr = append(bufr, byte(rbyte))
			rbyte = (r & 0x00ff0000) >> 16
			bufr = append(bufr, byte(rbyte))
			rbyte = (r & 0x0000ff00) >> 8
			bufr = append(bufr, byte(rbyte))
			rbyte = (r & 0x000000ff) >> 0
			bufr = append(bufr, byte(rbyte))

			g := data.DataG[x][y]
			//bing := make([]byte, 4)
			//binary.LittleEndian.PutUint32(bing, g)
			gbyte := (g & 0xff000000) >> 24
			bufg = append(bufg, byte(gbyte))
			gbyte = (g & 0x00ff0000) >> 16
			bufg = append(bufg, byte(gbyte))
			gbyte = (g & 0x0000ff00) >> 8
			bufg = append(bufg, byte(gbyte))
			gbyte = (g & 0x000000ff) >> 0
			bufg = append(bufg, byte(gbyte))
		}

	}
	//packet.dataR = []byte(string(bufr) + "\r")
	//packet.dataG = []byte(string(bufg) + "\r")
	fin := make(chan bool)
	go func() {
		packet.DataR = append(packet.DataR, bufr...)
		packet.DataR = append(packet.DataR, []byte("\r")...)
		fin <- true
	}()
	go func() {
		packet.DataG = append(packet.DataG, bufg...)
		packet.DataG = append(packet.DataG, []byte("\r")...)
		fin <- true
	}()
	<-fin
	<-fin

	packet.Terminator = "end\r"

	return &packet
}
