package internal

import (
	"strings"
)

func returnValue(idx uint, inputRune []rune) (returnedIdx uint, value any) {
	var (
		dc        uint8
		ss        string
		curToken  rune
		peekToken rune
		valBuf    strings.Builder
	)
	// preallocate memory
	valBuf.Grow(40)

	for ;; idx++ {
		curToken = inputRune[idx]

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreSpaceTab(idx, inputRune)
			}

		case DOUBLEQUOTE:
			valBuf.WriteRune(curToken)
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			valBuf.WriteRune(curToken)
			if peekToken = inputRune[idx+1]; dc > 0 && peekToken == DOUBLEQUOTE {
				dc--
			}

			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreComments(idx, inputRune)
			}

		case COMMA:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss = valBuf.String()
				if ss != "" {
					value = determineType(ss)
				} else if ss == "" {
					value = ss
				}
				valBuf.Reset()
				returnedIdx = idx
				return returnedIdx, value
			}

		case RBRACKET:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

		case RBRACE:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

		case lrTOKEN:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

		case lnTOKEN:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

		default:
			valBuf.WriteRune(curToken)
		}
	}
}
