package internal

import (
	"strings"
	"strconv"
)

func returnValue(idx uint, inputRune []rune) (returnedIdx uint, value any) {
    var (
		dc uint8
		ss string
		curToken rune
		peekToken rune
		valBuf strings.Builder
		internalIdx uint = idx
		valueTerminus int = int(searchValueTerminus(internalIdx, inputRune))
	)

	// preallocation of memory.
	valBuf.Grow(valueTerminus)

	for ; internalIdx <= uint(valueTerminus); internalIdx++ {
		curToken = inputRune[internalIdx]
		peekToken = inputRune[internalIdx + 1]

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				valBuf.WriteRune(curToken)
				continue
			}
			continue
		
		case DOUBLEQUOTE:
			valBuf.WriteRune(curToken)
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			valBuf.WriteRune(curToken)
			if dc > 0 && peekToken == DOUBLEQUOTE {
				dc--
			}

		case COMMA:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
			}

		case RBRACKET:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
			}

		case RBRACE:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
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
	returnedIdx = internalIdx
	return returnedIdx, value
	}


func searchValueTerminus(internalIdx uint, inputRune []rune) uint {
	var (
		dc uint8
		curToken rune
		peekToken rune
		terminalIdx uint = internalIdx
	)

	for ;; internalIdx++{
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

		case COMMA:
			if dc > 0 {
				terminalIdx++
			}

			if dc == 0 {
				return terminalIdx
			}

		case RBRACKET:
			if dc > 0 {
				terminalIdx++
			}

			if dc == 0 {
				return terminalIdx
			}

		case RBRACE:
			if dc > 0 {
				terminalIdx++
			}

			if dc == 0 {
				return terminalIdx
			}

		case lrTOKEN:
			if dc > 0 {
				terminalIdx++
			}

		case lnTOKEN:
			if dc > 0 {
				terminalIdx++
			}

		default:
			terminalIdx++
		}
	}
}


func determineType(ss string) any {
	if num, err := strconv.Atoi(ss); err == nil {
		return num

	} else if tr := strings.TrimSpace(ss); tr == "true" || tr == "false" {
		b, _ := strconv.ParseBool(tr)
		return b

	} else {
		return ss
	}
}