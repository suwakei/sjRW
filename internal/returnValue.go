package internal

import (
	"strings"
	"strconv"
)

func returnValue(commonIdx *uint, inputRune []rune) (value any) {
    var (
		dc uint8
		ss string
		curToken rune
		peekToken rune
		valBuf strings.Builder
		internalIdx uint = *commonIdx
		valueTerminus int = int(searchValueTerminus(internalIdx, inputRune))
	)

	// preallocation of memory.
	valBuf.Grow(valueTerminus)

	for ; internalIdx <= uint(valueTerminus); internalIdx++ {
		curToken = inputRune[internalIdx]
		peekToken = inputRune[internalIdx + 1]

		switch curToken {
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

		case lrTOKEN:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				if peekToken == lnTOKEN {
					continue
				}

				ss = valBuf.String()
				value = determineType(ss)
			}

		case lnTOKEN:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
			}

		default:
			valBuf.WriteRune(curToken)

			if dc == 0 {
				if curToken == RBRACE || curToken == RBRACKET {
					ss = valBuf.String()
					value = determineType(ss)
				}
			}
		}
	}
	*commonIdx = internalIdx
	return value
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

		case lrTOKEN:
			if dc > 0 {
				terminalIdx++
			}

			if dc == 0 {
				if peekToken == lnTOKEN {
					continue
				}
				return terminalIdx
			}

		case lnTOKEN:
			if dc > 0 {
				terminalIdx++
			}

			if dc == 0 {
				return terminalIdx
			}

		default:
			terminalIdx++
		}
	}
}


func determineType(ss string) any {
	var value any

	if num, err := strconv.Atoi(ss); err == nil {
		value = num
		return value

	} else if tr := strings.TrimSpace(ss); tr == "true" || tr == "false" {
		b, _ := strconv.ParseBool(tr)
		value = b
		return value

	} else {
		return ss
	}
}