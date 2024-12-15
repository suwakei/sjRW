package internal

import (
	"fmt"
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

// AssembleMap returns map created by input []rune
func AssembleMap(inputRune []rune) (assembledMap map[uint]map[string]any) {
	var (
		curToken rune // The target token.
		peekToken rune // The token for confirmation of next character.

		runeLength uint = uint(len(inputRune)) // The length of input rune slice.

		idx uint
		returnedIdx uint
		lineCount uint // Counter for current number of line.
		returnedLineCount uint

		firstLoop bool = true // First loop flag.
		keyMode bool = true //  If true, read jsonKey.

		returnedSlice []any
		returnedMap map[string]any
		returnedKey string
		returnedValue any
	)

	// preallocation of memory.
	assembledMap = make(map[uint]map[string]any, lnNum(inputRune))


	for ;; idx++ {
		curToken = inputRune[idx]

		if firstLoop {
			if _, ok := assembledMap[idx]; !ok {
				assembledMap[idx] = make(map[string]any, 1)
			}
			assembledMap[idx][string(curToken)] = ""
			firstLoop = false
			continue
		}

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
			lineCount++
			if _, ok := assembledMap[lineCount]; !ok {
				assembledMap[lineCount] = make(map[string]any, 1)
			}

			assembledMap[lineCount][string(curToken)] = ""
			break
		}

		if keyMode {
			returnedIdx, returnedKey = returnKey(idx, inputRune)
			idx += returnedIdx
			keyMode = false
		}

		if !keyMode && curToken == LBRACKET {
			keyMode = true

		} else if !keyMode && curToken == LBRACE {
			returnedIdx, returnedLineCount, returnedSlice = returnArr(idx, lineCount, inputRune)
			idx += returnedIdx
			lineCount += returnedLineCount
			keyMode = true

		} else if !keyMode && !isIgnores(curToken) {
			returnedIdx, returnedValue = returnValue(idx, inputRune)
			idx += returnedIdx
			keyMode = true

		}else if !keyMode && isIgnores(curToken) {
			continue
		}
	}
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
	return false
}