package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/huin/goserial"
)

type Packet struct {
	header     string
	coord      string
	data_r     string
	data_g     string
	terminator string
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

	fmt.Println(serialPort)
}
