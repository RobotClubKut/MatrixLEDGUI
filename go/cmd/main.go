package main

import (
	"fmt"
	"log"

	"github.com/RobotClubKut/MatrixLEDGUI/go/career"
	"github.com/RobotClubKut/MatrixLEDGUI/go/font"
	"github.com/RobotClubKut/MatrixLEDGUI/go/matrix"
	"github.com/huin/goserial"
)

//"github.com/RobotClubKut/MatrixLEDGUI/go/ledMatrix/"

//import "github.com/RobotClubKut/MatrixLEDGUI/go/crs"

func main() {
	font, err := font.SelectFont()
	if err != nil {
		log.Fatalln(err)
	}

	str := matrix.NewMatrixString("にゃんぱす", 0xffff00, font)
	ttyPort, err := career.ViewTtySelecterUI()
	if err != nil {
		log.Fatalln(err)
	}

	serialConfigure := &goserial.Config{Name: ttyPort, Baud: 9600}
	fmt.Println(serialConfigure)

	serialFin := make(chan bool)
	go career.SendMatrixString(serialConfigure, str, serialFin)
	<-serialFin

}
