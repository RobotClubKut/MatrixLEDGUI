package matrix

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"

	"github.com/RobotClubKut/MatrixLEDGUI/go/strimage"
	"github.com/pwaller/go-clz4"
)

type MatrixCharData struct {
	Bitmap [16][]byte
	Color  uint32
}

type MatrixChar struct {
	Bitmap [16]uint32
	Color  uint32
}

type MatrixString struct {
	Char  []MatrixCharData
	Coord uint32
}

func NewMatrixString(str string, color uint32, font string) *MatrixString {
	var ret MatrixString

	for _, char := range str {
		//圧縮した文字を格納する
		ret.Char = append(ret.Char, *compressMatrixChar(string(char), font, color))
	}
	return &ret
}

func ConnectMatrixString(s0 *MatrixString, s1 *MatrixString) *MatrixString {
	var ret MatrixString

	ret.Char = append(s0.Char, s1.Char...)
	ret.Coord = 0
	return &ret
}

func compressMatrixChar(c string, font string, color uint32) *MatrixCharData {
	image, err := strimage.ConvertString2image(c, font)
	if err != nil {
		log.Fatalln(err)
	}
	img := strimage.CancellationAntiAliasing(image)

	var ret MatrixCharData
	ret.Color = color

	for y := 0; y < 16; y++ {
		//ret.bitmap[y] = 0
		var bitmap uint32
		bitmap = 0
		for x := 0; x < 16; x++ {
			//ret.bitmap[y] = ret.bitmap[y] << 1
			bitmap = bitmap << 1
			r, _, _, _ := img.At(x, y).RGBA()
			if r == 0 {
				bitmap |= 1
			} else {
				bitmap |= 0
			}
		}
		var buf []byte
		b := (bitmap & 0xff000000) >> 24
		buf = append(buf, byte(b))
		b = (bitmap & 0x00ff0000) >> 16
		buf = append(buf, byte(b))
		b = (bitmap & 0x0000ff00) >> 8
		buf = append(buf, byte(b))
		b = (bitmap & 0x000000ff) >> 0
		buf = append(buf, byte(b))
		//圧縮して押し込む
		ch := make(chan []byte)
		go compressChar("clz4", buf, ch)
		ret.Bitmap[y] = <-ch
	}
	return &ret
}

func ReadMatrixChar(cm MatrixCharData) *MatrixChar {
	var ret MatrixChar
	ret.Color = cm.Color
	for y := 0; y < 16; y++ {
		buf := uncompressChar("clz4", cm.Bitmap[y])
		ret.Bitmap[y] = 0
		fin := make(chan bool)
		go func() {
			ret.Bitmap[y] |= ((uint32(buf[0]) << 24) & 0xff000000)
			fin <- true
		}()
		go func() {
			ret.Bitmap[y] |= ((uint32(buf[1]) << 16) & 0x00ff0000)
			fin <- true
		}()
		go func() {
			ret.Bitmap[y] |= ((uint32(buf[2]) << 8) & 0x0000ff00)
			fin <- true
		}()
		go func() {
			ret.Bitmap[y] |= ((uint32(buf[3]) << 0) & 0x000000ff)
			fin <- true
		}()
		<-fin
		<-fin
		<-fin
		<-fin
	}
	return &ret
}

func compressChar(mode string, srcBytes []byte, ch chan []byte) {
	if mode == "zlib" {
		var srcBuffer bytes.Buffer
		zlibWriter := zlib.NewWriter(&srcBuffer)
		zlibWriter.Write(srcBytes)
		zlibWriter.Close()
		ch <- srcBuffer.Bytes()
		//return srcBuffer.Bytes()
	} else if mode == "clz4" {
		out := []byte{}
		err := clz4.Compress(srcBytes, &out)
		if err != nil {
			log.Fatalln(err)
		}
		if len(out) == 0 {
			log.Fatalln("empty")
		}
		ch <- out
	}
}
func uncompressChar(mode string, srcBytes []byte) []byte {
	if mode == "zlib" {
		var srcBuffer bytes.Buffer
		var distBuf bytes.Buffer
		srcBuffer.Write(srcBytes)
		zlibReader, err := zlib.NewReader(&srcBuffer)
		if err != nil {
			fmt.Println("Can't reading data")
		}
		io.Copy(&distBuf, zlibReader)
		zlibReader.Close()
		return distBuf.Bytes()
	} else if mode == "clz4" {
		dec := make([]byte, len(srcBytes))
		err := clz4.Uncompress(srcBytes, &dec)
		if err != nil {
			log.Fatalln(err)
		}
		return dec
	}
	return nil
}

/*
func compressChar(srcBytes []byte, ch chan []byte) {
	out := []byte{}
	err := clz4.Compress(srcBytes, &out)
	if err != nil {
		log.Fatalln(err)
	}
	if len(out) == 0 {
		log.Fatalln("empty")
	}
	ch <- out
}
func uncompressChar(srcBytes []byte) []byte {
	dec := make([]byte, len(srcBytes))
	err := clz4.Uncompress(srcBytes, &dec)
	if err != nil {
		log.Fatalln(err)
	}
	return dec
}
*/
