package internal

import (
	"strconv"
	"strings"
)

//returnObjでinf recursionが起きてるっぽい.
// returnObj returns map[string]any
func returnObj(commonIdx, commonLineCount *uint, inputRune []rune) map[string]any {
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		curRuneLength uint = uint(len(inputRune[*commonIdx:])) // The length of input rune slice.

		dc uint8 // Counter for number of ".
		internalIdx uint
		internalLineCount uint
		keyMode bool = true // If true, mode which read json key. if false read json value.
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value *SA = new(SA) // The variable is for concatenated tokens stored in "valBuf".
		rs []any
		rm map[string]any
	)
// Store *commonIdx and *commonLineCount to avoid calling type *uint again and again
	internalIdx = *commonIdx
	internalLineCount = *commonLineCount

	// preallocation of memory
	rm = make(map[string]any, commaNum(internalIdx, inputRune))

	var keyBufMemoryNumber float32 = float32(commaNum(internalIdx, inputRune)) * 0.3
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(commaNum(internalIdx, inputRune)) * 0.3
	valBuf.Grow(int(valBufMemoryNumber))


	for ; internalIdx < curRuneLength; internalIdx++ {
		curToken = inputRune[internalIdx]

		if firstLoop {
			keyBuf.WriteRune(curToken)
			key = keyBuf.String()
			firstLoop = false
			keyMode = false
			continue
		}

		// This if expression for preventation of "index out of range" error
		if internalIdx + 1 < curRuneLength {
		peekToken = inputRune[internalIdx + 1]
		}


		// When the token is last
		if (internalIdx + 1 == curRuneLength) && (curToken == RBRACE || curToken == RBRACKET){
			internalLineCount++
			strCurToken := string(curToken)
			if _, ok := rm[internalLineCount]; !ok {
				rm[internalLineCount] = make(map[string]any, 1)
			}
			rm[internalLineCount][strCurToken] = ""
			keyBuf.Reset()
			valBuf.Reset()
			*commonIdx = internalIdx
			*commonLineCount = internalLineCount
			return rm

		switch curToken {

		case SPACE, TAB, COMMA: // "space" "\t"
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

		case COLON: // ":"
			if keyMode {
				if dc == 0 {
					key = keyBuf.String()
					keyMode = false
					continue
				}
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				}

				if dc == 0 {
					continue
				}
			}

		case DOUBLEQUOTE: // "
			if keyMode {
				dc++
				keyBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
				}
			}

			if !keyMode {
				dc++
				valBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
				}
			}

		case BACKSLASH: // "\"
			if keyMode {
				keyBuf.WriteRune(curToken)
				if dc == 1 && peekToken == DOUBLEQUOTE {
					dc--
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if dc == 1 && peekToken == DOUBLEQUOTE {
					dc--
			}
		}

		case LBRACKET: // "["
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				}
			}

			if !keyMode {
				if dc == 0 {
					rs = returnArr(&internalIdx, &internalLineCount, inputRune)
					value.valArrAny = rs
					keyBuf.Reset()
					valBuf.Reset()
				}

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
				if dc == 0 {
					rm = returnObj(&internalIdx, &internalLineCount, inputRune)
					value.valMap = rm
					keyBuf.Reset()
					valBuf.Reset()
				}
			}

		case RBRACE, RBRACKET: // "}" "]"
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

		case lrTOKEN: // "\r"
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
					if peekToken == lnTOKEN {
						continue
					}

					if peekToken != lnTOKEN {
						internalLineCount++
						value.valStr = valBuf.String()

						if _, ok := rm[internalLineCount]; !ok {
							rm[internalLineCount] = make(map[string]any, 1)
							if value.valArrAny == nil && value.valMap == nil {

								// determine whether "value.valStr" is int or not 
								if num, err := strconv.Atoi(value.valStr); err == nil {
									rm[internalLineCount][key] = num
									continue
								}

								// determine whether "value.valStr" is bool or not
								if tr := strings.TrimSpace(value.valStr); tr == "true"|| tr == "false" {
									b, _ := strconv.ParseBool(tr)
									rm[internalLineCount][key] = b
									continue
								}

								rm[internalLineCount][key] = value.valStr
							}

							if value.valArrAny != nil {
								rm[internalLineCount][key] = value.valArrAny
								value.valArrAny = nil
							}

							if value.valMap != nil {
								rm[internalLineCount][key] = value.valMap
								value.valMap = nil
							}
						}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
					}
				}
			}

		case lnTOKEN: // "\n"
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
					internalLineCount++
					value.valStr = valBuf.String()

					if _, ok := rm[internalLineCount]; !ok {
						rm[internalLineCount] = make(map[string]any, 1)

						if value.valArrAny == nil && value.valMap == nil{

							// determine whether "value.valStr" is int or not 
							if num, err := strconv.Atoi(value.valStr); err == nil {
								rm[internalLineCount][key] = num
								continue
							}

							// determine whether "value.valStr" is bool or not
							if tr := strings.TrimSpace(value.valStr); tr == "true"|| tr == "false" {
								b, _ := strconv.ParseBool(tr)
								rm[internalLineCount][key] = b
								continue
							}

							rm[internalLineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							rm[internalLineCount][key] = value.valArrAny
							value.valArrAny = nil
						}

						if value.valMap != nil {
							rm[internalLineCount][key] = value.valArrAny
							value.valMap = nil
						}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
			}
		}
	}

		default: // undefined token
			if keyMode {
				keyBuf.WriteRune(curToken)
			}
			if !keyMode {
				valBuf.WriteRune(curToken)
			}
		}
	}
	return rm
}
}

func searchBlockTerminal() {
	
}