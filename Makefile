.PHONY: build test clean static_analysis lint vet fmt chkfmt


default:
	all

all:
	test fmt b


# 実行した回数
# １回あたりの実行に掛かった時間(ns/op)
# １回あたりのアロケーションで確保した容量(B/op)
# 1回あたりのアロケーション回数(allocs/op) 

b:
	@go test -bench=. -benchmem

bp:
	@go test -bench=. -benchmem -cpuprofile cpu.prof

bm:
	@go test -bench=. -benchmem -memprofile mem.prof

bpm:
	@go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof

tm:
	@go test -run TestReadAsMapFrom

ts:
	@go test -run TestReadAsStrFrom

tb:
	@go test -run TestReadAsByteFrom


test:
	@go test -v .

fmt:
	@echo "[info ***********************formatting****************************]"
	@go fmt ./...