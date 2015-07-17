package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/huin/goserial"
)

type packet struct {
	header     string
	coord      string
	dataR      string
	dataG      string
	terminator string
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

func main() {
	ttyPort, err := viewTtySelecterUI()
	if err != nil {
		log.Fatalln(err)
	}

	serialConfigure := &goserial.Config{Name: ttyPort, Baud: 9600}
	serialPort, _ := goserial.OpenPort(serialConfigure)

	packet := createTestPacket()
	for {
		writeLCDMatrix(packet, serialPort)
	}

}
