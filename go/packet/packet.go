package packet

import "github.com/RobotClubKut/MatrixLEDGUI/go/matrix"

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

func printLCD(str matrix.MatrixString, shift int) *LcdMatrix {
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
