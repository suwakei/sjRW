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


// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[int]map[string]any) {

	var (
		curToken rune // The target token
		peekToken rune // The token for confirmation of next character

		runeLength int = len(inputRune) // The length of input rune slice

		doubleQuoteCnt int = 0 // Counter for number of ".
		lineCount int = 0 // Counter for current number of line.
		tempCount int = 0 // Match with idx later when exiting sliceMode.
		keyMode bool = true // If true, mode which read json key. if false read json value.
		sliceMode bool = false // If true, which read array of json value.
		firstLoop bool = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for storing key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for storing value token.
		key string // The variable for concatenated tokens stored in "keyBuf". 
		value string// The variable for concatenated tokens stored in "valBuf".
	)

	// preallocation of memory
	var initMap map[int]map[string]any = make(map[int]map[string]any, lnNum(inputRune))

	var keyBufMemoryNumber float32 = float32(runeLength) * 0.2
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(runeLength) * 0.5
	valBuf.Grow(int(valBufMemoryNumber))


	for idx := range inputRune {
		if sliceMode {
			if idx <= tempCount {
				continue
			}
			sliceMode = false
		}

		curToken = inputRune[idx]

		// For preventation of "index out of range" error
		if idx + 1 < runeLength {
		peekToken = inputRune[idx + 1]
		}

		/* Delete when finised debug */
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
		if (idx + 1 == runeLength) && (curToken == RBRACE || curToken == RBRACKET){
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
					internalLineCount, sliceModeIdx, returnedSlice := returnSliceOrMapAndCount(idx, inputRune)

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount + 1] = make(map[string]any, 1)
					}

					initMap[lineCount + 1][key] = returnedSlice
					keyBuf.Reset()
					valBuf.Reset()
					lineCount += internalLineCount
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
						value = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
					}

					initMap[lineCount][key] = value
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
					value = valBuf.String()

					if _, ok := initMap[lineCount]; !ok {
						initMap[lineCount] = make(map[string]any, 1)
					}

					initMap[lineCount][key] = value
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


// returnSliceOrMapAndCount returns three return value.
// "internalLineCount" is lineCount which counted inside this func.
// "sliceModeIdx" is idx which counted inside this func.
// "tempSlice" is slice created inside this func, this slice is used to be assigned to initMap
func returnSliceOrMapAndCount(curIdx int, inputRune []rune) (
		interLineCount,
		sliceModeIdx int,
		tempSlice []any,
		) {

	var (
		commanum = commaNum(inputRune, curIdx)// the number of commas, uses for allocating memory of "tempSlice"
		runeLength = len(inputRune)// The length of input rune slice
		dc int = 0 // This variable is same as variable "doubleQuoteCount" in the function "AssembleMap"
		tempRune rune // This variable works just like the variable "curToken" in the function "AssembleMap"
		peekTempRune rune // This variable works just like the variable "peekToken" in the function "AssembleMap"
		sliceBuf strings.Builder // When in "sliceMode" is true, buf for storing slice token.
		ss string // The variable for concatnated tokens stored in "sliceBuf".
	)

	sliceModeIdx = 0

	// preallocate memory
	var sliceBufMemoryNumber float32 = float32(runeLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))
	tempSlice = make([]any, 0, commanum)

	for i := curIdx + 1; i < runeLength; i++ {
		sliceModeIdx += 1
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
					tempSlice = append(tempSlice, num)
					continue
				}

				// determine whether "ss" is bool or not
				if tr := strings.TrimSpace(ss); tr == "true"|| tr == "false" {
					b, _ := strconv.ParseBool(tr)
					tempSlice = append(tempSlice, b)
					continue
}
				tempSlice = append(tempSlice, ss)
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
					interLineCount += 1
					continue
				}
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

// "lnNum" returns the number of "\n" or "\r" from "r".
// this return value used for initializing memory of "initMap" 
func lnNum(r []rune) int {
	var lnCount int = 0
	for _, n := range r {
		if n == lnTOKEN || n == lrTOKEN {
			lnCount += 1
		}
	}
	return lnCount
}