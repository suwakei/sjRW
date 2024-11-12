package main

import (
	"fmt"
	"github.com/suwakei/sjrw"
)

func main() {
	sj1 := sjrw.SjReader{JsonPath: "../testdata/test.json"}
	sj2 := sjrw.SjReader{JsonPath: "../testdata/test2.json"}
	b1, _ := sj1.ReadAsBytes()
	b2, _ := sj2.ReadAsBytes()

	// bb := Diff("test1", b1, "test2", b2)

	fmt.Println(string(b1), string(b2))
}