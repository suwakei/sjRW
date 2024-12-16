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

	str := getArrBlock(idx, arrTerminus, inputRune)
}


func getArrTerminusAndCommaNum(internalIdx uint, inputRune []rune) (uint, uint) {
	var (
		dc uint8 = 0 // dc stands for doubleQuoteCount
		commaCount uint = 0 // Counter for the number of commas
		terminalIdx uint = internalIdx
		lBracketCount uint = 0 // Counter for the number of left brackets
		rBracketCount uint = 0 // Counter for the number of right brackets
		curToken rune
		peekToken rune
	)

	for {
		curToken = inputRune[internalIdx]
		peekToken = inputRune[internalIdx + 1]
		switch curToken {
		case DOUBLEQUOTE:
			dc++
			terminalIdx++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			terminalIdx++
			if peekToken == DOUBLEQUOTE {
				dc--
			}

		case COMMA:
			terminalIdx++
			if dc == 0 {
				commaCount++
			}

		case LBRACKET:
			terminalIdx++
			if dc == 0 {
				lBracketCount++
			}

		case RBRACKET:
			terminalIdx++
			if dc == 0 {
				rBracketCount++
			}

			if lBracketCount == rBracketCount {
				return terminalIdx, commaCount + 1
			}
		}
	internalIdx++
	}
}

func getArrBlock(idx uint, arrTerminus uint, inputRune []rune) string {
	var (
		curToken rune = inputRune[idx]
		strBuf strings.Builder
	)

	for {
		if idx != arrTerminus {
			strBuf.WriteRune(curToken)
			break
		}
	}
	return strBuf.String()
}