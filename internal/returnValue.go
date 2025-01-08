package internal

import (
	"strings"
)

func (a *Assemble) returnValue(inputRune []rune) (value any) {
	var (
		dc        uint8
		ss        string
		curToken  rune
		peekToken rune
		valBuf    strings.Builder
	)
	// preallocate memory
	valBuf.Grow(40)

	for ; ; a.idx++ {
		curToken = inputRune[a.idx]

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if dc == 0 {
				a.ignoreSpaceTab(inputRune)
			}

		case DOUBLEQUOTE:
			valBuf.WriteRune(curToken)
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			valBuf.WriteRune(curToken)
			if dc > 0 {
				if peekToken = inputRune[a.idx+1]; dc > 0 && peekToken == DOUBLEQUOTE {
					dc--
				}
			}

		case SLASH:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if peekToken = inputRune[a.idx+1]; dc == 0 && peekToken == SLASH  {
				a.ignoreComments(inputRune)
			}

		case COMMA:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
				valBuf.Reset()
				return value
			}

		default:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}
		}
	}
}
