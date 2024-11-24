package main

import "fmt"

func main() {

	d := []rune("    \"add\": [")
	a := d
	fmt.Println(a, len(d))
	fmt.Println("\"", []byte("\""), "\\n", []byte("\n"), " ", []byte("		"))
	for _, i := range a {
		fmt.Println(string(i))
	}
}