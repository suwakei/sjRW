package main

import (
	"log"
	"os"

	"github.com/suwakei/sjrw"
)

func main() {
	var sj sjrw.SjReader

	// カレントディレクトリを変える
	jsonPath5 := "../testdata/readtest5.json"
	f5, oerr := os.OpenFile(jsonPath5, os.O_RDONLY, 0o666)

	if oerr != nil {
		log.Fatal(oerr)
	}


	defer func() {
		err := f5.Close()
		if err != nil {
			log.Fatalf("could not close file \"%s\"", jsonPath5)
		}
	}()

	m, _ := sj.ReadAsMapFrom(f5)
	_ = m
}