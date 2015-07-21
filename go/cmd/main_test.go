package main

import "testing"

func BenchmarkXxxx(b *testing.B) {
	str := convertLCDString("高知工科大学　", 0xff0000)
	shiftCoord := len(str.c) * 16
	for i := 0; i < shiftCoord+96+1; i++ {
		createPacket(*str, i-96)
	}
}
