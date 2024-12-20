package internal

import (
	"strings"
	"strconv"
)


func returnValue(idx uint, inputRune []rune, valBuf strings.Builder) (returnedIdx uint, value any) {
    var (
		dc uint8
		ss string
		curToken rune
		peekToken rune
	)


	for ;; idx++ {
		curToken = inputRune[idx]
		peekToken = inputRune[idx + 1]

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
				valBuf.Reset()
				returnedIdx = idx
				return returnedIdx, value
			}

		case RBRACKET:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
				valBuf.Reset()
				returnedIdx = idx
				return returnedIdx, value
			}

		case RBRACE:
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}

			if dc == 0 {
				ss = valBuf.String()
				value = determineType(ss)
				valBuf.Reset()
				returnedIdx = idx
				return returnedIdx, value
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
			if dc > 0 {
				valBuf.WriteRune(curToken)
			}
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