package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suwakei/sjrw"
)

func main() {

	var sj sjrw.SjReader

	// カレントディレクトリを変える
	jsonPath5 := "../testdata/readtest5.json"
	f5, _ := os.OpenFile(jsonPath5, os.O_RDONLY, 0o666)


	defer func() {
		err := f5.Close()
		if err != nil {
			log.Fatalf("could not close file \"%s\"", jsonPath5)
		}
	}()

	m, _ := sj.ReadAsMapFrom(f5)
	_ = m
}