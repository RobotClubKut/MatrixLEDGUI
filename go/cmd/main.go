package main

import "github.com/RobotClubKut/MatrixLEDGUI/go/crs"

//"github.com/RobotClubKut/MatrixLEDGUI/go/ledMatrix/"

//import "github.com/RobotClubKut/MatrixLEDGUI/go/crs"

func main() {
	var matrix [16][16]uint32
	c := crs.NewCrs(matrix)

	crs.PrintCRS(c)
}
