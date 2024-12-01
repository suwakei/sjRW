package main

import (
	"fmt"
	"github.com/suwakei/sjrw"
)

func main() {

	var sj sjrw.SjReader

	// カレントディレクトリを変える
	m, _ := sj.ReadAsMapFrom("../testdata/readtest5.json")
	fmt.Println(m)
}