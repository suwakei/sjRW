package sjrw

import (
	"strings"
	"fmt"
)

func (a *assemble) returnArr(inputRune []rune) (rs []any) {
	var (
		curToken   rune
		peekToken  rune
		dc         uint8
		ss         string
		arrVal     any
		tempArrBuf strings.Builder
		firstLoop  bool = true
	)
	// preallocate memory
	rs = make([]any, 0, arrLength(a.idx, inputRune))
	tempArrBuf.Grow(15)

	for ; ; a.idx++ {
		curToken = inputRune[a.idx]
		peekToken = inputRune[a.idx+1]

		// FIXME for debug
		fmt.Println("idx", a.idx, "lineCount", a.lineCount, "curToken", string(curToken))
		fmt.Println("idx", a.idx+1, "lineCount", a.lineCount, "peekToken", string(peekToken))

		if firstLoop {
			firstLoop = false
			continue
		}

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				a.ignoreSpaceTab(inputRune)
			}

		case DOUBLEQUOTE:
			dc++
			tempArrBuf.WriteRune(curToken)
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			tempArrBuf.WriteRune(curToken)
			if peekToken == DOUBLEQUOTE {
				dc--
			}

		case SLASH:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if peekToken == SLASH && dc == 0 {
				a.ignoreComments(inputRune)
			}

		case LBRACKET:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				rrs := a.returnArr(inputRune)
				rs = append(rs, rrs)
				tempArrBuf.Reset()
			}

		case RBRACKET:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss = tempArrBuf.String()
				arrVal = determineType(ss)
				rs = append(rs, arrVal)
				tempArrBuf.Reset()
				return rs
			}

		case COMMA:// TODO commaとlbracket が続いたときいらない処理をしてしまう
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss = tempArrBuf.String()
				arrVal = determineType(ss)
				rs = append(rs, arrVal)
				tempArrBuf.Reset()
				return rs
			}

		case LBRACE:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				rrs := a.returnObj(inputRune)
				rs = append(rs, rrs)
			}

		case lrTOKEN:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				if peekToken == lnTOKEN {
					continue
				}
				a.lineCount++
			}

		case lnTOKEN:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				a.lineCount++
			}

		default:
			tempArrBuf.WriteRune(curToken)
		}
	}
}

func arrLength(idx uint, inputRune []rune) uint {
	var (
		curToken  rune
		peekToken rune
		dc        uint8 = 0
		arrLength uint  = 0
		lb        uint8 = 0
		rb        uint8 = 0
	)

	for ; ; idx++ {
		curToken = inputRune[idx]
		peekToken = inputRune[idx+1]
		switch curToken {
		case DOUBLEQUOTE:
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			if peekToken == DOUBLEQUOTE {
				dc--
			}

		case COMMA:
			if dc == 0 {
				arrLength++
			}

		case LBRACKET:
			if dc == 0 {
				lb++
			}

		case RBRACKET:
			if dc == 0 {
				rb++
			}
			if lb == rb {
				return arrLength + 1
			}
		}
	}
}
