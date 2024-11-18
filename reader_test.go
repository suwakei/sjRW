package sjrw

import (
	"log"
	"testing"
)


func BenchmarkReadAsStr(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	sj := Sj{ReadFilePath: jsonPath}
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsStr()
	if err != nil {
		log.Fatal(err)
	}
}
}


func BenchmarkReadAsBytes(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	sj := Sj{ReadFilePath: jsonPath}
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsBytes()
	if err != nil {
		log.Fatal(err)
	}
}
}