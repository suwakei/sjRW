package internal

import (
	"strings"
	"strconv"
)

// AssembleMap returns map created by input str
func AssembleMap(str string) (assembledMap map[int]map[string]any) {
	const (
		SPACE = ' '
		TAB = '\t'
		lnTOKEN = '\n'
		DOUBLEQUOTE = '"'
		COLON = ':'

		LBRACE = '{'
		RBRACE = '}'
		LBRACKET = '['
		RBRACKET = ']'
		COMMA = ','
	)

	var (
		curToken rune // The token of target
		peekToken rune // The token for confirmation of next character

		r []rune = []rune(str) // Input str transrated into rune slice
		strLength int = len(r) // The length of input rune slice

		doubleQuoteCnt int = 0 // Counter for number of ".
		lineCount int = 0 // Counter for current number of line.
		tempCount int = 0 // Match with idx later when exiting sliceMode.
		keyMode bool = true // If true, mode which read json key. if false read json value.
		sliceMode bool = false // If true, which read array of json value.
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for storing key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for storing value token.
		key string // The variable for concatnated tokens stored in "keyBuf". 
		value string// The variable for concatnated tokens stored in "valBuf".
	)

	// preallocation of memory
	var initMap map[int]map[string]any = make(map[int]map[string]any, strLength / 20)

	var keyBufMemoryNumber float32 = float32(strLength) * 0.2
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(strLength) * 0.7
	valBuf.Grow(int(valBufMemoryNumber))

	var runifiedStr []rune = make([]rune, 0, strLength)
	runifiedStr = r


	for idx := range runifiedStr {
		if sliceMode {
			if idx <= tempCount {
				continue
			}
			sliceMode = false
		}

		curToken = runifiedStr[idx]

		// For preventatin of "index out of range" error
		if idx + 1 < strLength {
		peekToken = runifiedStr[idx + 1]
		}

		if firstLoop {
			keyBuf.WriteRune(curToken)
			key = keyBuf.String()
			firstLoop = false
			keyMode = false
			continue
		}

		// When the token is last
		if (idx + 1 == strLength) && (curToken == RBRACE || curToken == RBRACKET){
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
				continue
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		case COLON: // ":"
			if keyMode {
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					keyMode = false
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		case DOUBLEQUOTE: // "
			if keyMode {
				keyBuf.WriteRune(curToken)
				doubleQuoteCnt += 1
				if doubleQuoteCnt == 2 && peekToken == COLON {
					doubleQuoteCnt = 0
					key = strings.TrimSpace(keyBuf.String())
					keyMode = false
					continue
				}
			}

			if !keyMode {
				doubleQuoteCnt += 1
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt == 2 && peekToken != COMMA {
					doubleQuoteCnt -= 1
					continue
				}
				if doubleQuoteCnt == 2 && peekToken == COMMA {
					doubleQuoteCnt = 0
				}
			}

		case COMMA: // ","
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if doubleQuoteCnt == 0 && peekToken == lnTOKEN{
					continue
				}
				if doubleQuoteCnt != 0 {
					valBuf.WriteRune(curToken)
				}
		}

		case LBRACKET: // "["
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					sliceMode = true
					internalLineCount, sliceModeIdx, returnedSlice := returnSliceOrMapAndCount(idx, runifiedStr)
					initMap[lineCount][key] = returnedSlice
					keyBuf.Reset()
					valBuf.Reset()
					lineCount += internalLineCount
					tempCount = idx + sliceModeIdx
				}
			}

		case RBRACE, RBRACKET: // "}" "]"
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		case lnTOKEN: // "\n"
			if keyMode {
				keyBuf.WriteRune(curToken)
			}
			if !keyMode {
				valBuf.WriteRune(curToken)

				if doubleQuoteCnt == 0 {
				lineCount += 1
				value = strings.TrimSpace(valBuf.String())

				if _, ok := initMap[lineCount]; !ok {
					initMap[lineCount] = make(map[string]any, len(value))
				}

				initMap[lineCount][key] = value
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
	assembledMap = initMap
	return assembledMap
	}



func commaNum(r []rune, curIdx int) int {
	var (
		commaCount int = 0 // Counter for number of commas
		lBracketCount int = 0 // Counter for number of left brackets
		rBracketCount int = 0 // Counter for number of right brackets
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

// returnSliceOrMapAndCount returns three return value.
// "internalLineCount" is lineCount which counted inside this func.
// "sliceModeIdx" is idx which counted inside this func.
// "tempSlice" is slice created inside this func, this slice used for assigning initMap
func returnSliceOrMapAndCount(curIdx int, runifiedStr []rune) (
		interLineCount,
		sliceModeIdx int,
		tempSlice []any,
		) {

			const (
				SPACE = ' '
				TAB = '\t'
				lnTOKEN = '\n'
				DOUBLEQUOTE = '"'
				COLON = ':'
		
				LBRACE = '{'
				RBRACE = '}'
				LBRACKET = '['
				RBRACKET = ']'
				COMMA = ','
			)

	var (
		commanum = commaNum(runifiedStr, curIdx)// number of commas, uses for allocating memory of "tempSlice"
		strLength = len(runifiedStr)// The length of input rune slice
		dc int = 0 // This variable is same as variable "doubleQuoteCount" in the function "AssembleMap"
		tempRune rune // This variable works just like the variable "curToken" in the function "AssembleMap"
		peekTempRune rune // This variable works just like the variable "peekToken" in the function "AssembleMap"
		sliceBuf strings.Builder // When in "sliceMode" is true, buf for storing slice token.
		ss string // The variable for concatnated tokens stored in "sliceBuf".
	)

	sliceModeIdx = 0

	// preallocate memory
	var sliceBufMemoryNumber float32 = float32(strLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))
	tempSlice = make([]any, 0, commanum)

	for i := curIdx + 1; i < strLength; i++ {
		sliceModeIdx += 1
		tempRune = rune(runifiedStr[i])
		switch tempRune {
		case DOUBLEQUOTE:
			sliceBuf.WriteRune(tempRune)
			dc += 1

			if peekTempRune = rune(runifiedStr[i + 1]); peekTempRune != COMMA && dc == 2{
				dc -= 1
			}

		case COMMA:
			ss = sliceBuf.String()

			if num, err := strconv.Atoi(ss); dc < 2 && err == nil {
				dc = 0
				sliceBuf.Reset()
				tempSlice = append(tempSlice, num)
				continue
			}
			dc = 0
			sliceBuf.Reset()
			tempSlice = append(tempSlice, ss)

		case lnTOKEN:
			sliceBuf.WriteRune(tempRune)

			if num, err := strconv.Atoi(ss); dc < 2 && err == nil {
				dc = 0
				sliceBuf.Reset()
				tempSlice = append(tempSlice, num)
				interLineCount += 1
				continue
			}
			dc = 0
			interLineCount += 1

		case RBRACKET:
			// When the token is last
			if peekTempRune = rune(runifiedStr[i + 1]); peekTempRune == COMMA && dc == 0{
				return interLineCount, sliceModeIdx, tempSlice
			}
			if dc != 0 {
				sliceBuf.WriteRune(tempRune)
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
	return 0, 0, nil
}