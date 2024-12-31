package internal

import (
	"strconv"
	"strings"
)

const (
	SPACE       rune = ' '
	TAB         rune = '\t'
	lnTOKEN     rune = '\n'
	lrTOKEN     rune = '\r'
	DOUBLEQUOTE rune = '"'
	COLON       rune = ':'
	LBRACE      rune = '{'
	RBRACE      rune = '}'
	LBRACKET    rune = '['
	RBRACKET    rune = ']'
	COMMA       rune = ','
	BACKSLASH   rune = '\\'
	SLASH       rune = '/'
)

// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken rune // The target token.

		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		idx               uint
		returnedIdx       uint
		lineCount         uint = 1 // Counter for current number of line.
		returnedLineCount uint

		firstLoop bool = true // First loop flag.
		keyMode   bool = true //  If true, read jsonKey.

		returnedSlice []any
		returnedMap   map[string]any
		returnedKey   string
		returnedValue any
	)

	// preallocation of memory.
	assembledMap = make(map[uint]map[string]any, lnNum(inputRune))

	for ; ; idx++ {
		curToken = inputRune[idx]

		if firstLoop {
			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
			}
			assembledMap[lineCount][string(curToken)] = ""
			firstLoop = false
			continue
		}

		// last loop.
		if idx+1 == runeLength-1 {
			lineCount++
			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
			}

			assembledMap[lineCount][string(curToken)] = ""
			break
		}

		if curToken == SLASH {
			idx = ignoreComments(idx, inputRune)
		}

		if curToken == lrTOKEN {
			if inputRune[idx+1] == lnTOKEN {
				continue
			}
			lineCount++
			continue
		}

		if curToken == lnTOKEN {
			lineCount++
			continue
		}

		if keyMode && isIgnores(curToken) {
			continue
		}

		if keyMode {
			returnedIdx, returnedKey = returnKey(idx, inputRune)
			idx = returnedIdx
			keyMode = false
			continue
		}

		if !keyMode && curToken == LBRACKET {
			returnedIdx, returnedLineCount, returnedSlice = returnArr(idx, lineCount, inputRune)
			idx = returnedIdx

			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
				assembledMap[lineCount][returnedKey] = returnedSlice
			}
			lineCount = returnedLineCount
			keyMode = true
			continue

		} else if !keyMode && curToken == LBRACE {
			returnedIdx, returnedLineCount, returnedMap = returnObj(idx, lineCount, inputRune)
			idx = returnedIdx

			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
				assembledMap[lineCount][returnedKey] = returnedMap
			}
			lineCount = returnedLineCount
			keyMode = true
			continue

		} else if !keyMode && !isIgnores(curToken) {
			returnedIdx, returnedValue = returnValue(idx, inputRune)
			idx = returnedIdx

			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
				assembledMap[lineCount][returnedKey] = returnedValue
			}
			keyMode = true
			continue

		} else if !keyMode && isIgnores(curToken) {
			continue
		}
	}
	return assembledMap
}

func determineType(ss string) any {
	if num, err := strconv.Atoi(ss); err == nil {
		return num

	} else if tr := strings.TrimSpace(ss); tr == "true" || tr == "false" {
		b, _ := strconv.ParseBool(tr)
		return b

	} else {
		return ss
	}
}

// "lnNum" returns the number of "\n" or "\r" from "r".
// this return value used for initializing memory of "initMap"
func lnNum(r []rune) uint {
	var lnCount uint = 0
	var dc uint8 = 0

	for i := 0; i < len(r); i++ {
		if dc > 0 && r[i] == BACKSLASH && r[i+1] == DOUBLEQUOTE {
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
	return lnCount + 1
}

func isIgnores(curToken rune) bool {
	if curToken == SPACE {
		return true
	}
	if curToken == TAB {
		return true
	}
	if curToken == COMMA {
		return true
	}
	if curToken == COLON {
		return true
	}
	return false
}

func ignoreComments(idx uint, inputRune []rune) uint {
	var (
		curToken  rune
		peekToken rune
	)

	for ; ; idx++ {
		curToken = inputRune[idx]
		switch curToken {
		case lrTOKEN:
			if peekToken = inputRune[idx+1]; peekToken == lnTOKEN {
				continue
			} else if peekToken = inputRune[idx+1]; peekToken != lnTOKEN {
				return idx
			}
		case lnTOKEN:
			return idx

		default:
			continue
		}
	}
}

func ignoreSpaceTab(idx uint, inputRune []rune) uint {
	var (
		curToken rune
	)

	for ; ; idx++ {
		curToken = inputRune[idx]
		switch curToken {
		case SPACE, TAB:
			continue

		default:
			return idx
		}
	}
}
