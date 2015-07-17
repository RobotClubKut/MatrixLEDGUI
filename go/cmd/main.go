package main

import (
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
	c         []lcdChar
	drowCoord int
}

type lcdMatrix struct {
	data [96][16]int
}

func createTestPacket() *packet {
	var p packet
	p.header = "pcmat\r"
	p.coord = "000\r5ff\r"
	s := ""
	for i := 0; i < 192; i++ {
		s += "p"
	}
	s += "\r"
	p.dataR = s
	p.dataG = s
	p.terminator = "end\r"

	return &p
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

func main() {
	/*
		ttyPort, err := viewTtySelecterUI()
		if err != nil {
			log.Fatalln(err)
		}

		serialConfigure := &goserial.Config{Name: ttyPort, Baud: 9600}
		serialPort, _ := goserial.OpenPort(serialConfigure)

		packet := createTestPacket()
		writeLCDMatrix(packet, serialPort)
	*/
	//gray := imaging.Grayscale(convertString2image("A"))
	img, _ := convertString2image("０")
	gray := cancellationAntiAliasing(img)

	imaging.Save(gray, "./grayscaled.png")
}
