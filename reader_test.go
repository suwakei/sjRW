package sjrw

import (
	"log"
	"testing"
)


func BenchmarkReadAsStr(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	var sj Sj
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsStrFrom(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
}
}


func BenchmarkReadAsBytes(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	var sj Sj
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsBytesFrom(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
}
}