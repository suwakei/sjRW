package internal

import (
	"fmt"
	"strconv"
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
	}

// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {

	var (
		curToken rune // The target token
		peekToken rune // The token for confirmation of next character

		runeLength uint = uint(len(inputRune)) // The length of input rune slice

		doubleQuoteCnt uint8 = 0 // Counter for number of ".
		lineCount uint = 0 // Counter for current number of line.
		tempCount uint = 0 // Match with idx later when exiting sliceMode.
		keyMode bool = true // If true, mode which read json key. if false read json value.
		sliceMode bool = false // If true, which read array of json value.
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for storing key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for storing value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value SA // The variable is for concatenated tokens stored in "valBuf".
		rv RV
		ilc uint // The variable is for storing "returnedValue.internalLineCount"
	)

	// initalize
	value.valArrAny = nil

	// preallocation of memory
	var initMap map[uint]map[string]any = make(map[uint]map[string]any, lnNum(inputRune))

	var keyBufMemoryNumber float32 = float32(runeLength) * 0.2
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(runeLength) * 0.5
	valBuf.Grow(int(valBufMemoryNumber))


	for idx := range inputRune {
		if sliceMode {
			if uint(idx) <= tempCount {
				continue
			}
			sliceMode = false
		}

		curToken = inputRune[idx]

		// For preventation of "index out of range" error
		if uint(idx) + 1 < runeLength {
		peekToken = inputRune[idx + 1]
		}

		/* Delete when debug finised */
		a := string(curToken)
		b := string(peekToken)
		fmt.Println(a)
		fmt.Println(b)

		if firstLoop {
			keyBuf.WriteRune(curToken)
			key = keyBuf.String()
			firstLoop = false
			keyMode = false
			continue
		}

		// When the token is last
		if (uint(idx) + 1 == runeLength) && (curToken == RBRACE || curToken == RBRACKET){
			lineCount += 1
			strCurToken := string(curToken)
			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[string]any, len(strCurToken))
			}
			initMap[lineCount][strCurToken] = ""
			keyBuf.Reset()
			valBuf.Reset()
			break
		}

// ポインタも意識してかく
		switch curToken {

		case SPACE, TAB: // "space" "\t"
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}
			}

		case COLON: // ":"
			if keyMode {
				if doubleQuoteCnt == 0 {
					key = keyBuf.String()
					keyMode = false
					continue
				}
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}

				if doubleQuoteCnt == 0 {
					continue
				}
			}

		case DOUBLEQUOTE: // "
			if keyMode {
				doubleQuoteCnt += 1
				keyBuf.WriteRune(curToken)
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					continue
				}
				continue
			}

			if !keyMode {
				doubleQuoteCnt += 1
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
				}
				continue
			}

		case BACKSLASH: // "\"
			if keyMode {
				keyBuf.WriteRune(curToken)
				if doubleQuoteCnt == 1 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt -= 1
					continue
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt == 1 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt -= 1
					continue
			}
		}

		case COMMA: // ","
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue // o
				}
				continue
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					continue
				}
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}
		}

		case LBRACKET: // "["
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					sliceMode = true
					returnedValue := rv.returnSliceOrMapAndCount(uint(idx), inputRune)
					value.valArrAny = returnedValue.rs
					ilc = returnedValue.internalLineCount
					keyBuf.Reset()
					valBuf.Reset()
					tempCount = uint(idx) + returnedValue.sliceModeIdx
					continue
				}

				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case RBRACE, RBRACKET: // "}" "]"
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}
			}

		case lrTOKEN: // "\r"
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}

				if doubleQuoteCnt == 0 {
					if peekToken == lnTOKEN {
						continue
					}

					if peekToken != lnTOKEN {
						lineCount += 1
						value.valStr = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
						if value.valArrAny != nil {
							initMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valArrAny = nil

						} else if value.valArrAny == nil {
							initMap[lineCount][key] = value.valStr
						}
					}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
					continue
					}
				}
			}

		case lnTOKEN: // "\n"
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
					continue
				}

				if doubleQuoteCnt == 0 {
					lineCount += 1
					value.valStr = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
						if value.valArrAny != nil {
							initMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valArrAny = nil

						} else if value.valArrAny == nil {
							initMap[lineCount][key] = value.valStr
						}
					}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
					continue
			}
		}

		default: // undefined token
			if keyMode {
				keyBuf.WriteRune(curToken)
				continue
			}
			if !keyMode {
				valBuf.WriteRune(curToken)
				continue
			}
		}
	}
	assembledMap = initMap
	return assembledMap
	}


type RV struct {
    internalLineCount uint // "internalLineCount" is lineCount that counted inside this method.
    sliceModeIdx uint // "sliceModeIdx" is idx that counted inside this method.
    rs []any // "rs" stands for returnSlice is slice created inside this method, this slice is used to be assigned to initMap
}

// returnSliceOrMapAndCount returns three return value.
func (rv RV) returnSliceOrMapAndCount(curIdx uint, inputRune []rune) RV {

	var (
		commanum = commaNum(inputRune, curIdx)// the number of commas, uses for allocating memory of "tempSlice"
		runeLength uint = uint(len(inputRune))// The length of input rune slice
		dc uint8 = 0 // This variable is same as variable "doubleQuoteCount" in the function "AssembleMap"
		tempRune rune // This variable works just like the variable "curToken" in the function "AssembleMap"
		peekTempRune rune // This variable works just like the variable "peekToken" in the function "AssembleMap"
		sliceBuf strings.Builder // When in "sliceMode" is true, buf for storing slice token.
		ss string // The variable for concatnated tokens stored in "sliceBuf".
	)

	rv.sliceModeIdx = 0

	// preallocate memory
	var sliceBufMemoryNumber float32 = float32(runeLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))
	rv.rs = make([]any, 0, commanum)

	for i := uint(curIdx) + 1; i < runeLength; i++ {
		rv.sliceModeIdx += 1
		tempRune = rune(inputRune[i])
		switch tempRune {

		case DOUBLEQUOTE:
			sliceBuf.WriteRune(tempRune)
			dc += 1
			if dc == 2{
				dc = 0
				continue
			}

		case BACKSLASH:
			sliceBuf.WriteRune(tempRune)
			if peekTempRune = rune(inputRune[i + 1]); dc == 1 && peekTempRune == DOUBLEQUOTE{
				dc -= 1
				continue
			}

		case COMMA:
			if dc > 0 {
				sliceBuf.WriteRune(tempRune)
				continue
			}

			if dc == 0 {
				ss = sliceBuf.String()
				sliceBuf.Reset()
  				// determine whether "ss" is int or not 
				if num, err := strconv.Atoi(ss); err == nil {
					rv.rs = append(rv.rs, num)
					continue
				}

				// determine whether "ss" is bool or not
				if tr := strings.TrimSpace(ss); tr == "true"|| tr == "false" {
					b, _ := strconv.ParseBool(tr)
					rv.rs = append(rv.rs, b)
					continue
}
				rv.rs = append(rv.rs, ss)
				continue
			}

		case lrTOKEN:
			if dc > 0 {
				sliceBuf.WriteRune(tempRune)
				continue
				}

			if dc == 0 {
				if peekTempRune = rune(inputRune[i + 1]); peekTempRune == lnTOKEN {
					continue
				}

				if peekTempRune = rune(inputRune[i + 1]); peekTempRune != lnTOKEN {
					rv.internalLineCount += 1
					continue
				}
			}

		case lnTOKEN:
			if dc > 0 {
			sliceBuf.WriteRune(tempRune)
			continue
			}

			if dc == 0 {
			rv.internalLineCount += 1
			continue
			}

		case RBRACKET:
			// When the token is last
			if dc == 0 {
				rv.internalLineCount += 1
				ss= sliceBuf.String()
				sliceBuf.Reset()

				if num, err := strconv.Atoi(ss); err == nil {
					rv.rs = append(rv.rs, num)
					continue
				}

				if tr := strings.TrimSpace(ss); tr == "true" || tr == "false" {
					b, _ := strconv.ParseBool(tr)
					rv.rs = append(rv.rs, b)
					continue
				}
				rv.rs = append(rv.rs, ss)
				return RV{
					internalLineCount: rv.internalLineCount,
					sliceModeIdx: rv.sliceModeIdx,
					rs: rv.rs,
				}
			}

			if dc > 0 {
				sliceBuf.WriteRune(tempRune)
				continue
			}

		case SPACE, TAB:
			if dc < 1 {
				continue
			}
			sliceBuf.WriteRune(tempRune)

		default:
			sliceBuf.WriteRune(tempRune)
		}
	}
	return RV{}
}

func commaNum(r []rune, curIdx uint) uint {
	var (
		commaCount uint = 0 // Counter for number of commas
		lBracketCount uint = 0 // Counter for number of left brackets
		rBracketCount uint = 0 // Counter for number of right brackets
	)

	for i := curIdx;; i++ {
		if r[i] == '[' {
			lBracketCount += 1
		}
		if r[i] == ',' {
			commaCount += 1
		}
		if r[i] == ']' {
			rBracketCount += 1
			if lBracketCount == rBracketCount {
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
	for _, n := range r {
		if n == lnTOKEN || n == lrTOKEN {
			lnCount += 1
		}
	}
	return lnCount
}