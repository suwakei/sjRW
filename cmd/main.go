package main

import (

	"github.com/suwakei/sjrw"
)

func main() {

	var sj sjrw.SjReader

	// dlvでステップ実行するためエントリポイントで実行
	m, _ := sj.ReadAsMapFrom("../testdata/readtest5.json")
	_ = m

}