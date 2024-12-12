package internal

import (
	"fmt"
	"strings"
)


const (
	SPACE = ' '
	TAB = '\t'
	lnTOKEN = '\n'
	lrTOKEN = '\r'
	DOUBLEQUOTE = '"'
	COLON = ':'
	LBRACE = '{'
	RBRACE = '}'
	LBRACKET = '['
	RBRACKET = ']'
	COMMA = ','
	BACKSLASH = '\\'
	)


type SA struct {
	valStr string
	valArrAny []any
	valMap map[string]any
}


type RV struct {
    rs []any // "rs" stands for returnSlice is created inside return method, this slice is used to be assigned to initMap.
	rm map[string]any // "rm" stands for returnMap is created inside return~ method, this slice is used to be assigned to initMap.
}

// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		dc uint8 = 0 // Counter for the number of ".
		idx uint = 0
		commonIdx *uint = new(uint)
		commonLineCount *uint = new(uint) // Counter for current number of line.

		firstLoop bool = true // First loop flag.
		keyMode bool = true //  If true, read jsonKey.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value *SA = new(SA) // The variable is for concatenated tokens stored in "valBuf".
		rv *RV = new(RV) // The struct for return value.
	)

	// preallocation of memory.
	var keyBufMemoryNumber float32 = float32(runeLength) * 0.1
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(runeLength) * 0.1
	valBuf.Grow(int(valBufMemoryNumber))

	assembledMap = make(map[uint]map[string]any, lnNum(inputRune))


	for {
		curToken = inputRune[idx]

		if firstLoop {
			keyBuf.WriteRune(curToken)
			key = keyBuf.String()
			value.valStr = ""
			firstLoop = false
			keyMode = false
			idx++
			continue
		}

		idx++

		// This if expression for preventation of "index out of range" error.
		if idx + 1 < runeLength {
		peekToken = inputRune[idx + 1]
		}

		/* Delete when debug finished */
		a := string(curToken)
		b := string(peekToken)
		fmt.Println(a)
		fmt.Println(b)


		// last loop.
		if (idx + 1 == runeLength) && (curToken == RBRACE || curToken == RBRACKET){
			*commonLineCount++
			strCurToken := string(curToken)
			if _, ok := assembledMap[*commonLineCount]; !ok {
				assembledMap[*commonLineCount] = make(map[string]any, 1)
			}

			assembledMap[*commonLineCount][strCurToken] = ""
			keyBuf.Reset()
			valBuf.Reset()
			break
		}


		switch curToken {

		case SPACE, TAB, COMMA: // "space" "\t" ",".
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

		case COLON: // ":".
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
					*commonIdx = idx
					rv = returnArr(commonIdx, commonLineCount, inputRune)
					value.valArrAny = rv.rs
					keyBuf.Reset()
					valBuf.Reset()
					idx = *commonIdx
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
					*commonIdx = idx
					rv = returnObj(commonIdx, commonLineCount, inputRune)
					value.valMap = rv.rm
					keyBuf.Reset()
					valBuf.Reset()
					idx = *commonIdx
				}

				if dc > 0 {
					valBuf.WriteRune(curToken)
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
						*commonLineCount++
						value.valStr = valBuf.String()

					if _, ok := assembledMap[*commonLineCount]; !ok {
						assembledMap[*commonLineCount] = make(map[string]any, 1)
						if value.valArrAny == nil && value.valMap == nil {
							assembledMap[*commonLineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							assembledMap[*commonLineCount][key] = value.valArrAny
							value.valArrAny = nil
						}

						if value.valMap != nil {
							assembledMap[*commonLineCount][key] = value.valMap
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
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				}

				if dc == 0 {
					*commonLineCount++
					value.valStr = valBuf.String()

					if _, ok := assembledMap[*commonLineCount]; !ok {
						assembledMap[*commonLineCount] = make(map[string]any, 1)
						if value.valArrAny == nil && value.valMap == nil {
							assembledMap[*commonLineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							assembledMap[*commonLineCount][key] = value.valArrAny
							value.valArrAny = nil
						}

						if value.valMap != nil {
							assembledMap[*commonLineCount][key] = value.valMap
							value.valMap = nil
						}
					}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
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
	return assembledMap
	}


func commaNum(curIdx uint, r []rune) uint {
	var (
		dc uint8 = 0 // dc stands for doubleQuoteCount
		commaCount uint = 0 // Counter for the number of commas
		lBracketCount uint = 0 // Counter for the number of left brackets
		rBracketCount uint = 0 // Counter for the number of right brackets
		lBraceCount uint = 0 // Counter for the number of left braces
		rBraceCount uint = 0 // Counter for the number of right braces
	)

	for i := curIdx;; i++ {
		if r[i] == BACKSLASH && r[i + 1] == DOUBLEQUOTE {
			dc--
		}

		if r[i] == DOUBLEQUOTE {
			dc++
			if dc == 2 {
				dc = 0
			}
		}

		if dc == 0 && r[i] == COMMA {
			commaCount++
		}

		if dc == 0 && r[i] == LBRACKET {
			lBracketCount++
		}

		if dc == 0 && r[i] == RBRACKET {
			rBracketCount++
			if lBracketCount == rBracketCount {
				break
			}
		}

		if dc == 0 && r[i] == LBRACE {
			lBraceCount++
		}

		if dc == 0 && r[i] == RBRACE {
			rBraceCount++
			if lBraceCount == rBraceCount {
				break
			}
		}
	}
	return commaCount + 1
}

// "lnNum" returns the number of "\n" or "\r" from "r".
// this return value used for initializing memory of "initMap" 
func lnNum(r []rune) uint {
	var lnCount uint = 0
	var dc uint8 = 0

	for i := 0 ; i < len(r); i++ {
		if dc > 0 && r[i] == BACKSLASH && r[i + 1] == DOUBLEQUOTE {
			dc--
		}

		if r[i] == DOUBLEQUOTE {
			dc++
			if dc == 2 {
				dc = 0
			}
		}

		if dc == 0 && r[i] == lrTOKEN {
			lnCount++
		}
	}
	return lnCount
}