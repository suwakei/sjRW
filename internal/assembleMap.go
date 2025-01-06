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

type Assemble struct {
	idx uint
	lineCount uint // Counter for current number of line.
}

// AssembleMap returns map created by input []rune
func (a *Assemble) AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken rune // The target token.
		peekToken rune
		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		returnedIdx       uint
		lineCount         uint = 1 
		returnedLineCount uint

		firstLoop bool = true // First loop flag.
		keyMode   bool = true //  If true, read jsonKey.

		returnedKey   string
		returnedValue any
		returnedSlice []any
		returnedMap   map[string]any
	)

	a.lineCount = 1

	// preallocation of memory.
	assembledMap = make(map[uint]map[string]any, lnNum(inputRune))

	for ;; a.idx++ {
		curToken = inputRune[a.idx]

		if firstLoop {
			if _, ok := assembledMap[a.lineCount]; !ok {
				assembledMap[a.lineCount] = make(map[string]any, 1)
			}
			assembledMap[a.lineCount][string(curToken)] = ""
			firstLoop = false
			continue
		}

		// last loop.
		if a.idx+1 == runeLength {
			a.lineCount++
			if _, ok := assembledMap[a.lineCount]; !ok {
				assembledMap[a.lineCount] = make(map[string]any, 1)
			}

			assembledMap[a.lineCount][string(curToken)] = ""
			break
		}

		if peekToken = inputRune[a.idx+1]; curToken == SLASH && peekToken == SLASH {
			a.idx = a.ignoreComments(a.idx, inputRune)
		}

		if curToken == lrTOKEN {
			if inputRune[a.idx+1] == lnTOKEN {
				continue
			}
			a.lineCount++
			continue
		}

		if curToken == lnTOKEN {
			a.lineCount++
			continue
		}

		if keyMode && isIgnores(curToken) {
			continue
		}

		if keyMode {
			returnedIdx, returnedKey = returnKey(a.idx, inputRune)
			a.idx = returnedIdx
			keyMode = false
			continue
		}

		if !keyMode && curToken == LBRACKET {
			returnedIdx, returnedLineCount, returnedSlice = returnArr(a.idx, lineCount, inputRune)
			a.idx = returnedIdx

			if _, ok := assembledMap[a.lineCount]; !ok {
				assembledMap[a.lineCount] = make(map[string]any, 1)
				assembledMap[a.lineCount][returnedKey] = returnedSlice
			}
			a.lineCount = returnedLineCount
			keyMode = true
			continue

		} else if !keyMode && curToken == LBRACE {
			returnedIdx, returnedLineCount, returnedMap = returnObj(a.idx, lineCount, inputRune)
			a.idx = returnedIdx

			if _, ok := assembledMap[a.lineCount]; !ok {
				assembledMap[a.lineCount] = make(map[string]any, 1)
				assembledMap[a.lineCount][returnedKey] = returnedMap
			}
			a.lineCount = returnedLineCount
			keyMode = true
			continue

		} else if !keyMode && !isIgnores(curToken) {
			returnedIdx, returnedValue = returnValue(a.idx, inputRune)
			a.idx = returnedIdx

			if _, ok := assembledMap[a.lineCount]; !ok {
				assembledMap[a.lineCount] = make(map[string]any, 1)
				assembledMap[a.lineCount][returnedKey] = returnedValue
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

func (a *Assemble )ignoreComments(idx uint, inputRune []rune) uint {
	var (
		peekToken rune
	)

	for ;; a.idx++ {
		peekToken = inputRune[a.idx+1]
		if peekToken == lrTOKEN || peekToken == lnTOKEN {
			return a.idx
		}
	}
}

func (a *Assemble)ignoreSpaceTab(idx uint, inputRune []rune) uint {
	var (
		curToken rune
		peekToken rune
	)

	for ;; a.idx++ {
		curToken = inputRune[a.idx]
		switch curToken {
		case SPACE, TAB:
			if peekToken = inputRune[a.idx+1]; peekToken != SPACE && peekToken != TAB {
				return a.idx
			}
			continue
		}
	}
}
