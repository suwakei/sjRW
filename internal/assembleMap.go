package internal

import (
	"strings"
	"strconv"
)

// AssembleMap returns map created by input []rune
func AssembleMap(input []rune) (assembledMap map[int]map[string]any) {
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
		BACKSLASH = '\\'
	)

	var (
		curToken rune // The token of target
		peekToken rune // The token for confirmation of next character

		r []rune = input // "Input" stored rune slice
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

	var valBufMemoryNumber float32 = float32(strLength) * 0.5
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
				if doubleQuoteCnt > 0 {
				keyBuf.WriteRune(curToken) // o
				continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
				valBuf.WriteRune(curToken) // o
				continue
				}
			}

		case COLON: // ":"
			if keyMode {
				if doubleQuoteCnt == 2 { // o
					doubleQuoteCnt = 0
					keyMode = false
					continue
				}
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if doubleQuoteCnt > 0 { // o
					valBuf.WriteRune(curToken)
					continue
				}

				if doubleQuoteCnt == 0 { // o
					continue
				}
			}

		case DOUBLEQUOTE: // "
			if keyMode {
				doubleQuoteCnt += 1
				keyBuf.WriteRune(curToken) // o
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					continue
				}
				continue
			}

			if !keyMode {
				doubleQuoteCnt += 1
				valBuf.WriteRune(curToken) // o
				if doubleQuoteCnt == 2 {
					doubleQuoteCnt = 0
					continue
				}
				continue
			}

		case BACKSLASH: // "\"
			if keyMode {
				keyBuf.WriteRune(curToken)
				if doubleQuoteCnt > 0 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt -= 1 // o
					continue
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if doubleQuoteCnt > 0 && peekToken == DOUBLEQUOTE {
					doubleQuoteCnt -= 1 // o
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
				if doubleQuoteCnt == 0 { // o
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
					keyBuf.WriteRune(curToken) // o
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					sliceMode = true
					internalLineCount, sliceModeIdx, returnedSlice := returnSliceOrMapAndCount(idx, runifiedStr)

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
					}

					initMap[lineCount][key] = returnedSlice
					keyBuf.Reset()
					valBuf.Reset()
					lineCount += internalLineCount // o
					tempCount = idx + sliceModeIdx
					continue
				}

				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken)
				}
			}

		case RBRACE, RBRACKET: // "}" "]"
			if keyMode {
				if doubleQuoteCnt > 0 {
				keyBuf.WriteRune(curToken) // o
				continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken) // o
					continue
				}
			}

		case lnTOKEN: // "\n"
			if keyMode {
				if doubleQuoteCnt > 0 {
					keyBuf.WriteRune(curToken) // o
					continue
				}
			}

			if !keyMode {
				if doubleQuoteCnt > 0 {
					valBuf.WriteRune(curToken) // o
					continue
				}

				if doubleQuoteCnt == 0 {
					lineCount += 1
					value = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
					}

					initMap[lineCount][key] = value // o
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
// "tempSlice" is slice created inside this func, this slice is used to be assigned to initMap
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
				BACKSLASH = '\\'
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
			if dc == 2{
				dc = 0
				continue
			}

		case BACKSLASH:
			sliceBuf.WriteRune(tempRune)
			if peekTempRune = rune(runifiedStr[i + 1]); dc > 0 && peekTempRune == DOUBLEQUOTE{
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
  				// whether 
				if num, err := strconv.Atoi(ss); err == nil {
					tempSlice = append(tempSlice, num)
					continue
				}

				if tr := strings.TrimSpace(ss); tr == "true"|| tr == "false" {
					b, _ := strconv.ParseBool(tr)
					tempSlice = append(tempSlice, b)
					continue
}
				tempSlice = append(tempSlice, ss)
				continue
			}

		case lnTOKEN:
			if dc > 0 {
			sliceBuf.WriteRune(tempRune)
			continue
			}

			if dc == 0 {
			interLineCount += 1
			continue
			}

		case RBRACKET:
			// When the token is last
			if dc == 0 {
				interLineCount += 1
				ss= sliceBuf.String()
				sliceBuf.Reset()
  				// whether 
				if num, err := strconv.Atoi(ss); err == nil {
					tempSlice = append(tempSlice, num)
					continue
				}

				if tr := strings.TrimSpace(ss); tr == "true" || tr == "false" {
					b, _ := strconv.ParseBool(tr)
					tempSlice = append(tempSlice, b)
					continue
				}
				tempSlice = append(tempSlice, ss)
				return interLineCount, sliceModeIdx, tempSlice
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
	return 0, 0, nil
}