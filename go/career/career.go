package career

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/RobotClubKut/MatrixLEDGUI/go/matrix"
	"github.com/RobotClubKut/MatrixLEDGUI/go/packet"
	"github.com/huin/goserial"
)

func getUsbttyList() []string {
	contents, _ := ioutil.ReadDir("/dev")
	var ret []string

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usb") ||
			strings.Contains(f.Name(), "ttyUSB") ||
			strings.Contains(f.Name(), "ttyACM") {
			//return "/dev/" + f.Name()
			ret = append(ret, "/dev/"+f.Name())
		}
	}

	return ret
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

// ViewTtySelecterUI UI Viewer
func ViewTtySelecterUI() (string, error) {
	ttys := getUsbttyList()
	tty, err := ttySelecter(ttys)

	fmt.Println(tty)
	return tty, err
}

func writeLCDMatrix(p *packet.Packet, s io.ReadWriteCloser) {
	s.Write([]byte(p.Header))
	s.Write([]byte(p.Coord))
	s.Write([]byte(p.DataR))
	s.Write([]byte(p.DataG))
	s.Write([]byte(p.Terminator))
}

func SendMatrixString(serialConfigure *goserial.Config, str *matrix.MatrixString, fin chan bool) {
	for {
		shiftCoord := len(str.Char) * 16
		for i := 0; i < shiftCoord+96+1; i++ {
			serialPort, _ := goserial.OpenPort(serialConfigure)
			packet := packet.CreatePacket(*str, i-96)
			writeLCDMatrix(packet, serialPort)
			//time.Sleep(1 * time.Millisecond)
			serialPort.Close()
		}
	}
	fin <- true
}
