package crs

//参考文献
//http://homepage.cs.uiowa.edu/~sriram/21/fall07/code/CRS.java

import "fmt"

// Crs is crs container
type Crs struct {
	//非要素リスト
	//疎行列が0, 1しかないので1の長さがわかればok
	Val    []uint32
	ColIdx []uint32
	RowPtr []uint32
	MSize  uint32
}

// NewCrs is Constractor
func NewCrs(matrix [16][16]uint32) *Crs {
	var i uint32
	var j uint32

	var totalNonZeros uint32
	var index uint32

	mSize := uint32(len(matrix))

	totalNonZeros = 0

	for i = 0; i < mSize; i++ {
		for j = 0; j < mSize; j++ {
			totalNonZeros++
		}
	}

	val := make([]uint32, totalNonZeros)
	colIdx := make([]uint32, totalNonZeros)

	rowPtr := make([]uint32, mSize+1)
	rowPtr[0] = 0

	index = 0

	for i = 0; i < mSize; i++ {
		for j = 0; j < mSize; j++ {
			if matrix[i][j] != 0 {
				val[index] = matrix[i][j]
				colIdx[index] = uint32(j)
				index++
			}
		}
		rowPtr[i+1] = index
	}

	return &Crs{Val: val, ColIdx: colIdx, RowPtr: rowPtr, MSize: mSize}
}

// PrintCRS print CRS
func PrintCRS(crs *Crs) {
	var i int

	fmt.Println("print vectors used in CRS")
	fmt.Println("The vector val")

	for i = 0; i < len(crs.Val); i++ {
		fmt.Print(crs.Val[i], ", ")
	}
	fmt.Println()
	fmt.Println("The vector colIdx")
	for i = 0; i < len(crs.ColIdx); i++ {
		fmt.Print(crs.ColIdx[i], ", ")
	}
	fmt.Println()
	fmt.Println("The vector row_ptr")
	for i = 0; i < len(crs.RowPtr); i++ {
		fmt.Print(crs.RowPtr[i], ", ")
	}
	fmt.Println()
}
