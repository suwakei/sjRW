package internal

import (
	"strings"
)

func returnArr(idx, lineCount uint, inputRune []rune) (returnedIdx, returnedLineCount uint, rs []any) {
	var (
		curToken   rune
		peekToken  rune
		firstLoop  bool = true
		dc         uint8
		ss         string
		arrVal     any
		tempArrBuf strings.Builder
	)
	// preallocate memory
	rs = make([]any, 0, arrLength(idx, inputRune))
	tempArrBuf.Grow(15)

	for ; ; idx++ {
		curToken = inputRune[idx]

		if firstLoop {
			firstLoop = false
			continue
		}

		switch curToken {
		case SPACE, TAB:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreSpaceTab(idx, inputRune)
			}

		case DOUBLEQUOTE:
			dc++
			tempArrBuf.WriteRune(curToken)
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			tempArrBuf.WriteRune(curToken)
			if peekToken = inputRune[idx+1]; peekToken == DOUBLEQUOTE {
				dc--
			}

		case SLASH:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				idx = ignoreComments(idx, inputRune)
			}

		case LBRACKET:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				rdx, rlc, rrs := returnArr(idx, lineCount, inputRune)
				idx = rdx
				lineCount = rlc
				rs = append(rs, rrs)
				tempArrBuf.Reset()
			}

		case RBRACKET:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss = tempArrBuf.String()
				if ss != "" {
					arrVal = determineType(ss)
				} else if ss == "" {
					arrVal = ss
				}
				rs = append(rs, arrVal)
				returnedIdx = idx
				returnedLineCount = lineCount
				tempArrBuf.Reset()
				return returnedIdx, returnedLineCount, rs
			}

		case COMMA:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				ss := tempArrBuf.String()
				arrVal = determineType(ss)
				rs = append(rs, arrVal)
				tempArrBuf.Reset()
			}

		case LBRACE:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				rdx, rlc, rrs := returnObj(idx, lineCount, inputRune)
				idx = rdx
				lineCount = rlc
				rs = append(rs, rrs)
			}

		case lrTOKEN:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				if peekToken = inputRune[idx+1]; peekToken == lnTOKEN {
					continue
				}
				lineCount++
			}

		case lnTOKEN:
			if dc > 0 {
				tempArrBuf.WriteRune(curToken)
			} else if dc == 0 {
				lineCount++
			}

		default:
			tempArrBuf.WriteRune(curToken)
		}
	}
}

func arrLength(idx uint, inputRune []rune) uint {
	var (
		curToken  rune
		dc        uint8
		arrLength uint
		lb        uint8
		rb        uint8
	)

	for ; ; idx++ {
		curToken = inputRune[idx]
		switch curToken {
		case DOUBLEQUOTE:
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			if peekToken := inputRune[idx+1]; peekToken == DOUBLEQUOTE {
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
