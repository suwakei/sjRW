package sjrw

import (
	"strconv"
	"strings"
)

const (
	SPACE       rune = ' '  // U+0020 32
	TAB         rune = '\t' // U+0009 9 
	lnTOKEN     rune = '\n' // U+000A 10
	lrTOKEN     rune = '\r' // U+000D 13
	DOUBLEQUOTE rune = '"'  // U+0022 34
	COLON       rune = ':'  // U+003A 58
	LBRACE      rune = '{'  // U+007B 123
	RBRACE      rune = '}'  // U+007D 125
	LBRACKET    rune = '['  // U+005B 91
	RBRACKET    rune = ']'  // U+005D 93
	COMMA       rune = ','  // U+002C 44
	BACKSLASH   rune = '\\' // U+005C 92
	SLASH       rune = '/'  // U+002F 47
	ASTERISK    rune = '*'  // U+002A 42
)

type assemble struct {
	idx       uint
	lineCount uint // Counter for current number of line.
}

// AssembleMap returns map created by input []rune
func (a *assemble) assembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken   rune // The target token.
		peekToken  rune
		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		tempLineCount uint

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

	for ; ; a.idx++ {
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
			a.ignoreComments(inputRune)
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
			returnedKey = a.returnKey(inputRune)
			keyMode = false
			continue
		}

		if !keyMode && curToken == LBRACKET {
			tempLineCount = a.lineCount
			returnedSlice = a.returnArr(inputRune)

			if _, ok := assembledMap[tempLineCount]; !ok {
				assembledMap[tempLineCount] = make(map[string]any, 1)
				assembledMap[tempLineCount][returnedKey] = returnedSlice
			}
			keyMode = true
			continue

		} else if !keyMode && curToken == LBRACE {
			tempLineCount = a.lineCount
			returnedMap = a.returnObj(inputRune)

			if _, ok := assembledMap[tempLineCount]; !ok {
				assembledMap[tempLineCount] = make(map[string]any, 1)
				assembledMap[tempLineCount][returnedKey] = returnedMap
			}
			keyMode = true
			continue

		} else if !keyMode && !isIgnores(curToken) {
			returnedValue = a.returnValue(inputRune)

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

func (a *assemble) ignoreComments(inputRune []rune) {
	var (
		curToken rune = inputRune[a.idx]
		peekToken rune = inputRune[a.idx+1]
	)

	if curToken == SLASH && peekToken == SLASH {
		for ;; a.idx++ {
			peekToken = inputRune[a.idx+1]
			if peekToken == lrTOKEN || peekToken == lnTOKEN {
				return
			}
		}
	} else if curToken == SLASH && peekToken == ASTERISK {
		for ;; a.idx++ {
			curToken = inputRune[a.idx]
			peekToken = inputRune[a.idx+1]
			if curToken == ASTERISK && peekToken == SLASH {
				a.idx += 1
				return
			}
		}
	}
}

func (a *assemble) ignoreSpaceTab(inputRune []rune) {
	var (
		curToken  rune
		peekToken rune
	)

	for ; ; a.idx++ {
		curToken = inputRune[a.idx]
		switch curToken {
		case SPACE, TAB:
			if peekToken = inputRune[a.idx+1]; peekToken != SPACE && peekToken != TAB {
				return
			}
			continue
		}
	}
}
