package font

import (
	"fmt"
	"io/ioutil"
	"log"
)

func SelectFont() (string, error) {
	fontDir := "../fonts/"
	list, err := ioutil.ReadDir(fontDir)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for i, f := range list {
		fmt.Printf("%d: ", i)
		fmt.Println(f.Name())
	}

	n := 0
	fmt.Println()
	fmt.Print("Select font file: ")
	fmt.Scan(&n)
	fmt.Println()
	return "../font/" + list[n].Name(), nil
}
