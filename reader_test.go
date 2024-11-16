package sjrw

import (
	"fmt"
	"log"
	"testing"

	"github.com/suwakei/sjrw/internal"
)


func BenchmarkReadAsStr(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	sj := Sj{ReadJsonPath: jsonPath}
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsStr()
	if err != nil {
		log.Fatal(err)
	}
}
}


func BenchmarkReadAsBytes(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	sj := Sj{ReadJsonPath: jsonPath}
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsBytes()
	if err != nil {
		log.Fatal(err)
	}
}
}

func TestDiff(t *testing.T) {
	var jsonPath1 string = "./testdata/difftest.json"
	var jsonPath2 string = "./testdata/difftest2.json"
	sj1 := Sj{ReadJsonPath: jsonPath1}
	sj2 := Sj{ReadJsonPath: jsonPath2}
	s1, _ := sj1.ReadAsBytes()
	s2, _ := sj2.ReadAsBytes()
	i := internal.Diff("test1", s1, "test2", s2)
	fmt.Println(string(i))
}