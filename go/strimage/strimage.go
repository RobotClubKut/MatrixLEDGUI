package strimage

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/disintegration/imaging"

	"code.google.com/p/freetype-go/freetype"
)

func ConvertString2image(s string, fontfile string) (*image.RGBA, error) {
	dpi := float64(72.0)
	//fontfile := "../font/MS Gothic.ttf"

	//fontfile := "../font/VL.ttf"
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

func CancellationAntiAliasing(img *image.RGBA) *image.NRGBA {
	gray := imaging.Grayscale(img)
	//imaging.Save(gray, "./grayscaled.png")
	w := gray.Rect.Max.X
	h := gray.Rect.Max.Y

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			if r > 38000 {
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
