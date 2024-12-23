package internal

import (
	"strings"
)

//returnObjでinf recursionが起きてるっぽい.
// returnObj returns map[string]any
func returnObj(idx, lineCount uint, inputRune []rune) (returnedIdx, returnedLineCount uint, rm map[string]any) {
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		dc uint8 // Counter for number of ".
		keyMode bool = true // If true, mode which read json key. if false read json value.
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value any
	)

	// preallocation of memory
	rm = make(map[string]any, mapLength(idx, inputRune))

	keyBuf.Grow(20)

	valBuf.Grow(30)

	for ;; idx++ {
		curToken = inputRune[idx]

		if firstLoop {
			rm[string(curToken)] = ""
			firstLoop = false
			continue
		}

		if int(idx) + 1 <= len(inputRune) {
			peekToken = inputRune[idx + 1]
		}

		switch curToken {
		case SPACE, TAB:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case DOUBLEQUOTE:
			dc++
			if keyMode {
				keyBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
				}
			}
			
		case BACKSLASH:
			if keyMode {
				keyBuf.WriteRune(curToken)
				if peekToken == DOUBLEQUOTE {
					dc--
			}
		}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if peekToken == DOUBLEQUOTE {
					dc--
			}
		}

		case COLON:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}

				if dc == 0 {
					key = keyBuf.String()
					keyBuf.Reset()
					keyMode = false
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case LBRACE:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}

				if dc == 0 {
					rdx, rlc, rrs := returnObj(idx, lineCount, inputRune)
					idx += rdx
					lineCount += rlc
					value = rrs
				}
			}

		case RBRACE:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					ss := valBuf.String()
					valBuf.Reset()
					if ss != "" {
						value = determineType(ss)
					}
					rm[key] = value
					rm[string(curToken)] = ""
					returnedIdx = idx
					returnedLineCount = lineCount
					return returnedIdx, returnedLineCount, rm
				}
			}

		case LBRACKET:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					rdx, rlc, rrs := returnArr(idx, lineCount, inputRune)
					idx += rdx
					lineCount += rlc
					value = rrs
					valBuf.Reset()
				}
			}


		case COMMA:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					ss := valBuf.String()
					valBuf.Reset()
					if ss != "" {
						value = determineType(ss)
						rm[key] = value
					} else {
						rm[key] = ss
					}

					keyMode = true
				}
			}

		case lrTOKEN:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}

				if dc == 0 {
					if peekToken == lnTOKEN {
						continue
					}
					lineCount++
			}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					if peekToken == lnTOKEN {
						continue
					}
					lineCount++
				}
			}
		
		case lnTOKEN:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}

				if dc == 0 {
					lineCount++
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					lineCount++
				}
			}
		
		default:
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		}
	}

}

func mapLength(idx uint, inputRune []rune) uint {
	var (
		curToken rune
		peekToken rune
		firstLoop bool = true
		dc uint8
		mapLength uint
		lb uint8
		rb uint8
	)

	for ;; idx++ {
		curToken = inputRune[idx]

		if int(idx + 1) <= len(inputRune) {
			peekToken = inputRune[idx + 1]
		}

		if firstLoop {
			idx++
			firstLoop =  false
			continue
		}

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
				mapLength++
			}

		case LBRACE:
			if dc == 0 {
				lb++
			}

		case RBRACE:
			if dc == 0 {
				rb++
			}
			if lb == rb {
				return mapLength + 1
			}
		}
	}
}