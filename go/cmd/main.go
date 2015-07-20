package main

import (
	"log"

	"github.com/RobotClubKut/MatrixLEDTOOLS/go/career"
	"github.com/RobotClubKut/MatrixLEDTOOLS/go/font"
	"github.com/RobotClubKut/MatrixLEDTOOLS/go/matrix"
	"github.com/huin/goserial"
)

//"github.com/RobotClubKut/MatrixLEDTOOLS/go/ledMatrix/"

//import "github.com/RobotClubKut/MatrixLEDTOOLS/go/crs"

func main() {
	font, err := font.SelectFont()
	if err != nil {
		log.Fatalln(err)
	}

	str := matrix.NewMatrixString("にゃんぱすといろはす似てね？", 0xffff00, font)
	ttyPort, err := career.ViewTtySelecterUI()
	if err != nil {
		log.Fatalln(err)
	}

	serialConfigure := &goserial.Config{Name: ttyPort, Baud: 9600}
	//fmt.Println(serialConfigure)

	serialFin := make(chan bool)
	go career.SendMatrixString(serialConfigure, str, serialFin)
	<-serialFin

}
