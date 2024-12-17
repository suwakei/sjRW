package internal

import (
	"strings"
)


func returnKey(idx uint, inputRune []rune, keyBuf strings.Builder) (returnedIdx uint, key string){
	var (
		dc uint8
		curToken rune
		peekToken rune
		keyTerminus int = int(searchKeyTerminus(idx, inputRune))
	)

	// preallocation of memory.
	keyBuf.Grow(keyTerminus)

	for ; idx <= uint(keyTerminus); idx++ {
		curToken = inputRune[idx]
		peekToken = inputRune[idx + 1]

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
 keyBuf.Reset()
	returnedIdx = idx
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