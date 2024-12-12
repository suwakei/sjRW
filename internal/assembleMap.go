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
	valMap map[string]any
}

type RV struct {
    internalLineCount uint // "internalLineCount" is lineCount that counted inside this method.
    modeIdx uint // "modeIdx" is idx that counted inside retuen~ method.
    rs []any // "rs" stands for returnSlice is created inside return method, this slice is used to be assigned to initMap.
	rm map[string]any // "rm" stands for returnMap is created inside return~ method, this slice is used to be assigned to initMap.
}

// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		doubleQuoteCnt uint8 = 0 // Counter for the number of ".
		lineCount uint = 0 // Counter for current number of line.
		tempCount uint = 0 // Match with idx later when exiting sliceMode.
		
		sliceMode bool = false // If true, read array of json value.
		mapMode bool = false // If true, read map of json value.
		firstLoop bool = true // First loop flag.
		keyMode bool = true //  If true, read jsonKey.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value SA // The variable is for concatenated tokens stored in "valBuf".
		rv RV // The struct for return value.
		ilc uint // The variable is for storing "returnedValue.internalLineCount".
	)

	// initalize.
	value.valArrAny = nil
	value.valMap = nil

	// preallocation of memory.
	var keyBufMemoryNumber float32 = float32(runeLength) * 0.1
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(runeLength) * 0.1
	valBuf.Grow(int(valBufMemoryNumber))

	assembledMap = make(map[uint]map[string]any, lnNum(inputRune))


	for idx := range inputRune {
		if sliceMode {
			if uint(idx) <= tempCount {
				continue
			}
			sliceMode = false
		}

		if mapMode {
			if uint(idx) <= tempCount {
				continue
			}
			mapMode = false
		}

		curToken = inputRune[idx]

		// For preventation of "index out of range" error.
		if uint(idx) + 1 < runeLength {
		peekToken = inputRune[idx + 1]
		}

		/* Delete when debug finished */
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

		// When the token is last.
		if (uint(idx) + 1 == runeLength) && (curToken == RBRACE || curToken == RBRACKET){
			lineCount++
			strCurToken := string(curToken)
			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
			}
			assembledMap[lineCount][strCurToken] = ""
			keyBuf.Reset()
			valBuf.Reset()
			break
		}


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
				doubleQuoteCnt++
				keyBuf.WriteRune(curToken)
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					continue
				}
				continue
			}

			if !keyMode {
				doubleQuoteCnt++
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
					doubleQuoteCnt--
					continue
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt == 1 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt--
					continue
			}
		}

		case COMMA: // ","
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
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
					rv = returnSliceAndCount(uint(idx), inputRune)
					value.valArrAny = rv.rs
					ilc = rv.internalLineCount
					keyBuf.Reset()
					valBuf.Reset()
					tempCount = uint(idx) + rv.modeIdx
					continue
				}

				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case LBRACE:
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					mapMode = true
					rv = returnMapAndCount(uint(idx), inputRune)
					value.valMap = rv.rm
					ilc = rv.internalLineCount
					keyBuf.Reset()
					valBuf.Reset()
					tempCount = uint(idx) + rv.modeIdx
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
						lineCount++
						value.valStr = valBuf.String()

					if _, ok := assembledMap[lineCount]; !ok {
						assembledMap[lineCount] = make(map[string]any, 1)
						if value.valArrAny == nil && value.valMap == nil {
							assembledMap[lineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							assembledMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valArrAny = nil
						}

						if value.valMap != nil {
							assembledMap[lineCount][key] = value.valMap
							lineCount += ilc
							value.valMap = nil
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
					lineCount++
					value.valStr = valBuf.String()

					if _, ok := assembledMap[lineCount]; !ok {
						assembledMap[lineCount] = make(map[string]any, 1)
						if value.valArrAny == nil && value.valMap == nil {
							assembledMap[lineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							assembledMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valArrAny = nil
						}

						if value.valMap != nil {
							assembledMap[lineCount][key] = value.valMap
							lineCount += ilc
							value.valMap = nil
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
	return assembledMap
	}


// returnSliceOrMapAndCount returns three return value.
func returnSliceAndCount(curIdx uint, inputRune []rune) RV {
	var (
		commanum = commaNum(inputRune, curIdx)// the number of commas, uses for allocating memory of "tempSlice"
		curRuneLength uint = uint(len(inputRune[curIdx:]))// The length of input rune slice
		dc uint8 = 0 // This variable is same as variable "doubleQuoteCount" in the function "AssembleMap"
		tempRune rune // This variable works just like the variable "curToken" in the function "AssembleMap"
		peekTempRune rune // This variable works just like the variable "peekToken" in the function "AssembleMap"
		sliceBuf strings.Builder // When in "sliceMode" is true, buf for storing slice token.
		ss string // The variable for concatnated tokens stored in "sliceBuf".
		rv RV
	)

	rv.modeIdx = 0

	// preallocate memory
	var sliceBufMemoryNumber float32 = float32(curRuneLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))
	rv.rs = make([]any, 0, commanum)

	for i := uint(curIdx) + 1; i < curRuneLength; i++ {
		rv.modeIdx++
		tempRune = rune(inputRune[i])
		switch tempRune {

		case DOUBLEQUOTE:
			sliceBuf.WriteRune(tempRune)
			dc++
			if dc == 2{
				dc = 0
				continue
			}

		case BACKSLASH:
			sliceBuf.WriteRune(tempRune)
			if peekTempRune = rune(inputRune[i + 1]); dc == 1 && peekTempRune == DOUBLEQUOTE{
				dc--
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
					rv.internalLineCount++
					continue
				}
			}

		case lnTOKEN:
			if dc > 0 {
			sliceBuf.WriteRune(tempRune)
			continue
			}

			if dc == 0 {
			rv.internalLineCount++
			continue
			}

		case RBRACKET:
			// When the token is last
			if dc == 0 {
				rv.internalLineCount++
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
					modeIdx: rv.modeIdx,
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
returnMapAndCountでinf recursionが起きてるっぽい
func returnMapAndCount(curIdx uint, inputRune []rune) RV{
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		curRuneLength uint = uint(len(inputRune[curIdx:])) // The length of input rune slice.

		doubleQuoteCnt uint8 = 0 // Counter for number of ".
		lineCount uint = 0 // Counter for current number of line.
		tempCount uint = 0 // Match with idx later when exiting sliceMode.
		keyMode bool = true // If true, mode which read json key. if false read json value.
		sliceMode bool = false // If true, which read array of json value.
		mapMode bool = false
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key string // The variable is for concatenated tokens stored in "keyBuf". 
		value SA // The variable is for concatenated tokens stored in "valBuf".
		rv RV // The struct for return value.
		ilc uint // The variable is for storing "returnedValue.internalLineCount"
	)

	// initalize
	value.valArrAny = nil
	value.valMap = nil

	// preallocation of memory
	var initMap map[uint]map[string]any = make(map[uint]map[string]any, commaNum(inputRune, curIdx))

	var keyBufMemoryNumber float32 = float32(commaNum(inputRune, curIdx)) * 0.3
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(commaNum(inputRune, curIdx)) * 0.3
	valBuf.Grow(int(valBufMemoryNumber))


	for i := uint(curIdx) + 1; i < curRuneLength; i++ {
		if sliceMode {
			if uint(curIdx) <= tempCount {
				continue
			}
			sliceMode = false
		}

		if mapMode {
			if uint(curIdx) <= tempCount {
				continue
			}
			mapMode = false
		}

		curToken = inputRune[curIdx]

		// For preventation of "index out of range" error
		if uint(curIdx) + 1 < curRuneLength {
		peekToken = inputRune[curIdx + 1]
		}

		if firstLoop {
			keyBuf.WriteRune(curToken)
			key = keyBuf.String()
			firstLoop = false
			keyMode = false
			continue
		}

		// When the token is last
		if (uint(curIdx) + 1 == curRuneLength) && (curToken == RBRACE || curToken == RBRACKET){
			lineCount++
			strCurToken := string(curToken)
			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[string]any, 1)
			}
			initMap[lineCount][strCurToken] = ""
			keyBuf.Reset()
			valBuf.Reset()
			return RV{
				internalLineCount: rv.internalLineCount,
				modeIdx: rv.modeIdx,
				rm: rv.rm,
			}
		}

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
				doubleQuoteCnt++
				keyBuf.WriteRune(curToken)
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					continue
				}
				continue
			}

			if !keyMode {
				doubleQuoteCnt++
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
					doubleQuoteCnt--
					continue
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt == 1 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt--
					continue
			}
		}

		case COMMA: // ","
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
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
					rv = returnSliceAndCount(uint(curIdx), inputRune)
					value.valArrAny = rv.rs
					ilc = rv.internalLineCount
					keyBuf.Reset()
					valBuf.Reset()
					tempCount = uint(curIdx) + rv.modeIdx
					continue
				}

				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case LBRACE:
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					mapMode = true
					rv = returnMapAndCount(uint(curIdx), inputRune)
					value.valMap = rv.rm
					ilc = rv.internalLineCount
					keyBuf.Reset()
					valBuf.Reset()
					tempCount = uint(curIdx) + rv.modeIdx
					continue
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
						lineCount++
						value.valStr = valBuf.String()

						if _, ok := initMap[lineCount]; !ok {
							initMap[lineCount] = make(map[string]any, 1)
							if value.valArrAny == nil && value.valMap == nil {
								initMap[lineCount][key] = value.valStr
							}

							if value.valArrAny != nil {
								initMap[lineCount][key] = value.valArrAny
								lineCount += ilc
								value.valArrAny = nil
							}

							if value.valMap != nil {
								initMap[lineCount][key] = value.valMap
								lineCount += ilc
								value.valMap = nil
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
				if doubleQuoteCnt == 0 {
					lineCount++
					value.valStr = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)

						if value.valArrAny == nil && value.valMap == nil{
							initMap[lineCount][key] = value.valStr
						}

						if value.valArrAny != nil {
							initMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valArrAny = nil
						}

						if value.valMap != nil {
							initMap[lineCount][key] = value.valArrAny
							lineCount += ilc
							value.valMap = nil
						}
					keyBuf.Reset()
					valBuf.Reset()
					keyMode = true
					continue
			}
		}

		if doubleQuoteCnt > 0 {
			valBuf.WriteRune(curToken)
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
	return RV{}
	}



func commaNum(r []rune, curIdx uint) uint {
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
	for _, n := range r {
		if n == DOUBLEQUOTE {
			dc++
			if dc == 2 {
				dc = 0
			}
		}

		if n == lrTOKEN && dc == 0{
			lnCount++
		}
	}
	return lnCount
}