package internal

import (
	"strings"
)

interkalIdxをidxにする
func returnKey(idx uint, inputRune []rune) (returnedIdx uint, key string){
	var (
		dc uint8
		curToken rune
		peekToken rune
		keyBuf strings.Builder
		internalIdx uint = idx
		keyTerminus int = int(searchKeyTerminus(internalIdx, inputRune))
	)

	// preallocation of memory.
	keyBuf.Grow(keyTerminus)

	for ; internalIdx <= uint(keyTerminus); internalIdx++ {
		curToken = inputRune[internalIdx]
		peekToken = inputRune[internalIdx + 1]

		switch curToken {
		case DOUBLEQUOTE:
			keyBuf.WriteRune(curToken)
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			keyBuf.WriteRune(curToken)
			if dc > 0 && peekToken == DOUBLEQUOTE {
				dc--
			}

		case COLON:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			}

			if dc == 0 {
				break
			}
		}
	}
	key = keyBuf.String()
	returnedIdx = internalIdx
	return returnedIdx, key
}


func searchKeyTerminus(internalIdx uint, inputRune []rune) uint {
	var (
		dc uint8
		curToken rune
		peekToken rune
		terminalIdx uint = internalIdx
	)

	for ;; internalIdx++ {
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
			if dc > 0 && peekToken == DOUBLEQUOTE {
				dc--
				terminalIdx++
			}

		case COLON:
			if dc == 0 {
				return terminalIdx
			}

		default:
			terminalIdx++
		}
	}
}