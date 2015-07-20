package web

import (
	"fmt"
	"html"
	"net/http"
)

func top(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func update(w http.ResponseWriter, r *http.Request) (string, string) {

	str := r.FormValue("str")
	col := r.FormValue("col")

	fmt.Fprintf(w, "<html><body>Input String: %s, %s</body></html>",
		html.EscapeString(str), html.EscapeString(col))
	return str, col
}

func webServer(fin chan bool) {
	http.HandleFunc("/", top)
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan bool)
		go func() {
			str, col := update(w, r)
			if str != "" {
				c := 0
				if col == "red" {
					c = 0xff0000
				} else if col == "green" {
					c = 0x00ff00
				} else if col == "orange" {
					c = 0xffff00
				} else {
					c = 0xffff00
				}
				//lcdStringBuffer := ConvertLCDString(str, c)
				fmt.Println(c)
			}
			ch <- true

		}()
		<-ch
	})
	http.ListenAndServe(":8080", nil)

	fin <- true
}
