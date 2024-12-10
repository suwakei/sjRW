.PHONY: build test clean static_analysis lint vet fmt chkfmt


#go test個別実行の処理をかくgo test -run TestMapFromなど

all:
	go test -v .


# 実行した回数
# １回あたりの実行に掛かった時間(ns/op)
# １回あたりのアロケーションで確保した容量(B/op)
# 1回あたりのアロケーション回数(allocs/op) 

b:
	go test -bench=. -benchmem

bp:
	go test -bench=. -benchmem -cpuprofile cpu.prof

bm:
	go test -bench=. -benchmem -memprofile mem.prof

bpm:
	go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof

tm:
	go test -run TestReadAsMapFrom

ts:
	go test -run TestReadAsStrFrom

tb:
	go test -run TestReadAsByteFrom