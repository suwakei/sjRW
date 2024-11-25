package main

import (
	"github.com/suwakei/sjrw"
)

func main() {

sj sjrw.SjReader

// dlvでステップ実行するためエントリポイントで実行
m, _ := sj.ReadAsMapFrom("../testdata/readtest.json")
	fmt.Println(m)

}