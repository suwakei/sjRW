package internal

import (
	"strings"
)

func returnKey(idx uint, inputRune []rune) (returnedIdx uint, key string) {
	var (
		dc        uint8
		curToken  rune
		peekToken rune
		keyBuf    strings.Builder
	)
	// preallocate memory
	keyBuf.Grow(20)

	for ; ; idx++ {
		curToken = inputRune[idx]

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreSpaceTab(idx, inputRune)
				continue
			}

		case DOUBLEQUOTE:
			keyBuf.WriteRune(curToken)
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			keyBuf.WriteRune(curToken)
			if peekToken = inputRune[idx+1]; dc > 0 && peekToken == DOUBLEQUOTE {
				dc--
			}

		case SLASH:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreComments(idx, inputRune)
			}

		case COLON:
			if dc > 0 {
				keyBuf.WriteRune(curToken)
			} else if dc == 0 {
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
