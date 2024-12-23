package internal

import (
	"unsafe"
	"log"
)

type Buffer struct {
	address *Buffer
	buf []byte
}

func (b *Buffer) bufToString() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
	
}

// func (b *Buffer) grow(n int) {
// 	if n < 0 {
// 		log.Fatal("neganive number is invalid")
// 	}
// 	if cap(b.buf)-len(b.buf) < n {
// 		var tempBuf []byte = bytealg.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
// 		copy(tempBuf, b.buf)
// 		b.buf = tempBuf
// 	}
// }

func (b *Buffer) bufLen() int {
	return len(b.buf)
}

func (b *Buffer) bufCap() int  {
	return cap(b.buf)
}

func (b *Buffer) accumuRune(r rune) error {
	b.buf = append(b.buf, byte(r))
	return nil
}

func (b *Buffer) bufReset() {
	b.address = nil
	b.buf = nil
}

// resetを使うとGrowした容量が消えて再度Builderを使うときまた再割り当てが発生するため
// bufの容量を残したままbufの内容を消す処理をする
// b.buf[:0]で容量のみ遺す
func (b *Buffer) LeaveCap() {
	b.address = nil
	b.buf = b.buf[:0]
}