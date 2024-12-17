package internal

import (
	"strings"
)

// returnSliceOrMapAndCount returns *RV.
func returnArr(idx, lineCount uint, inputRune []rune) ( returnedIdx, returnedLineCount uint, rs []any) {
	var (
		curIdx uint
		curToken rune
		returnedValue any
		arrTerminus, commanum = getArrTerminusAndCommaNum(idx, inputRune)// the number of commas, uses for allocating memory of "tempSlice"
	)
	// preallocate memory
	rs = make([]any, 0, commanum)
}