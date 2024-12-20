package internal

import (
	"strings"
)


func returnKey(idx uint, inputRune []rune) (returnedIdx uint, key string){
	var (
		dc uint8
		curToken rune
		peekToken rune
		keyBuf strings.Builder
	)
	// preallocate memory
	keyBuf.Grow(20)

	for ;; idx++ {
		curToken = inputRune[idx]
		peekToken = inputRune[idx + 1]

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			}
			continue

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
				key = keyBuf.String()
				keyBuf.Reset()
				returnedIdx = idx
				return returnedIdx, key
			}

		default:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			}
		}
	}
}