package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/huin/goserial"

	"code.google.com/p/freetype-go/freetype"
)

type packet struct {
	header     string
	coord      string
	dataR      string
	dataG      string
	terminator string
}

type lcdChar struct {
	c      string
	bitmap [16]uint16
	color  int
}

type lcdString struct {
	c     []lcdChar
	coord int
}

type lcdMatrix struct {
	dataR [3][16]uint32
	dataG [3][16]uint32
}

func createTestPacket() *packet {
	var p packet
	p.header = "pcmat\r"
	p.coord = "000\r5ff\r"
	s1 := ""
	s2 := ""
	for i := 0; i < 192; i++ {
		s1 += "p"
		s2 += "a"
	}
	s1 += "\r"
	s2 += "\r"
	p.dataR = s1
	p.dataG = s2
	p.terminator = "end\r"

	return &p
}

func createPacket(str lcdString) *packet {
	data, _ := printLCD(str, 0)

	packet := createTestPacket()
	var bufr []byte
	var bufg []byte
	for y := 0; y < 16; y++ {

		for x := 0; x < 3; x++ {
			r := data.dataR[x][y]
			binr := make([]byte, 4)
			binary.LittleEndian.PutUint32(binr, r)
			bufr = append(bufr, binr...)

			g := data.dataR[x][y]
			bing := make([]byte, 4)
			binary.LittleEndian.PutUint32(bing, g)
			bufg = append(bufg, bing...)
		}

	}
	packet.dataR = string(bufr) + "\r"
	packet.dataG = string(bufg) + "\r"

	return packet
}

func getUsbttyList() []string {
	contents, _ := ioutil.ReadDir("/dev")
	var ret []string

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usb") ||
			strings.Contains(f.Name(), "ttyUSB") {
			//return "/dev/" + f.Name()
			ret = append(ret, "/dev/"+f.Name())
		}
	}

	return ret
}

func writeLCDMatrix(p *packet, s io.ReadWriteCloser) {
	s.Write([]byte(p.header))
	s.Write([]byte(p.coord))
	s.Write([]byte(p.dataR))
	s.Write([]byte(p.dataG))
	s.Write([]byte(p.terminator))
}

func ttySelecter(ttys []string) (string, error) {
	for i, s := range ttys {
		fmt.Println(strconv.Itoa(i) + ": " + s)
	}
	fmt.Println()
	fmt.Print("Select tty prot: ")

	n := 0
	fmt.Scan(&n)

	if n > len(ttys)-1 {
		return "", errors.New("tty port select error. " + strconv.Itoa(n) + " is not exist.")
	}

	return ttys[n], nil
}

func viewTtySelecterUI() (string, error) {
	ttys := getUsbttyList()
	tty, err := ttySelecter(ttys)

	fmt.Println(tty)
	return tty, err
}

func convertString2image(s string) (*image.RGBA, error) {
	dpi := float64(72.0)
	fontfile := "../font/MS Gothic.ttf"
	hinting := "none"
	size := float64(17)
	spacing := float64(0)

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 16, 16))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(freetype.NoHinting)
	case "full":
		c.SetHinting(freetype.FullHinting)
	}

	// Draw the text.
	pt := freetype.Pt(-1, -3+int(c.PointToFix32(size)>>8))

	_, err = c.DrawString(s, pt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pt.Y += c.PointToFix32(size * spacing)

	return rgba, nil
}

func cancellationAntiAliasing(img *image.RGBA) *image.NRGBA {
	gray := imaging.Grayscale(img)
	//imaging.Save(gray, "./grayscaled.png")
	w := gray.Rect.Max.X
	h := gray.Rect.Max.Y

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			if r > 27000 {
				c := color.RGBA{0xff, 0xff, 0xff, 0xff}
				gray.Set(x, y, c)
			} else {
				c := color.RGBA{0x00, 0x00, 0x00, 0xff}
				gray.Set(x, y, c)
			}
		}
	}
	return gray
}

func convertLCDChar(c string, color int) *lcdChar {
	image, _ := convertString2image(c)
	img := cancellationAntiAliasing(image)

	var ret lcdChar
	ret.color = color
	ret.c = c

	for y := 0; y < 16; y++ {
		ret.bitmap[y] = 0
		for x := 0; x < 16; x++ {
			ret.bitmap[y] = ret.bitmap[y] << 1
			r, _, _, _ := img.At(x, y).RGBA()
			if r == 0 {
				ret.bitmap[y] |= 1
			} else {
				ret.bitmap[y] |= 0
			}
		}
	}
	return &ret
}

func convertLCDString(str string, color int) *lcdString {
	var ret lcdString

	for _, c := range str {
		ret.c = append(ret.c, *convertLCDChar(string(c), color))
	}
	ret.coord = 0
	return &ret
}

func connectLCDStr(str0 *lcdString, str1 *lcdString) *lcdString {
	var ret lcdString
	ret.c = append(ret.c, str0.c...)
	ret.c = append(ret.c, str1.c...)
	ret.coord = 0
	return &ret
}

func printLCD(str lcdString, shift int) (*lcdMatrix, *lcdString) {
	var ret lcdMatrix

	for y := 0; y < 16; y++ {
		var bufR []uint8
		var bufG []uint8
		for cn := 0; cn < len(str.c); cn++ {
			for x := 0; x < 16; x++ {
				bit := str.c[cn].bitmap[y]
				bit = bit >> (15 - uint(x))
				//mask 0b0000000000000001
				bit = bit & 0x0001
				//赤があるとき
				//buf[coordX] = uint8(bit)
				if str.c[cn].color == 0xff0000 {
					bufR = append(bufR, uint8(bit))
				}
				if str.c[cn].color == 0x00ff00 {
					bufG = append(bufG, uint8(bit))
				}
			}
		}
		i := 0
		counter := 0

		for x := 0; x < 96; x++ {
			if counter == 32 {
				i++
				if i == 3 {
					break
				}
				counter = 0
			}
			ret.dataR[i][y] = ret.dataR[i][y] << 1
			if str.coord+x < len(bufR) {
				ret.dataR[i][y] |= uint32(bufR[str.coord+x])
			} else {
				ret.dataR[i][y] |= 0
			}
			ret.dataG[i][y] = ret.dataG[i][y] << 1
			if str.coord+x < len(bufG) {
				ret.dataG[i][y] |= uint32(bufG[str.coord+x])
			} else {
				ret.dataG[i][y] |= 0
			}
			counter++
		}
	}
	str.coord += shift
	fmt.Println(ret)
	return &ret, &str

}

func main() {

	//str := convertLCDString("a", 0x00ff00)
	//packet := createPacket(*str)

	ttyPort, err := viewTtySelecterUI()
	if err != nil {
		log.Fatalln(err)
	}

	serialConfigure := &goserial.Config{Name: ttyPort, Baud: 9600}
	serialPort, _ := goserial.OpenPort(serialConfigure)

	packet := createTestPacket()
	for i := 0; i < 100; i++ {
		writeLCDMatrix(packet, serialPort)
	}
	//gray := imaging.Grayscale(convertString2image("A"))
	//imaging.Save(gray, "./grayscaled.png")
}
