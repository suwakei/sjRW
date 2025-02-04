package sjrw

import (
	"strings"
)

func (a *assemble) handlePremitive(inputRune []rune, assembledMap map[uint]map[string]any, key string, ) {
	var (
		dc        uint8
		inQuote bool = false
		ss        string
		curToken  rune
		peekToken rune
		valBuf    strings.Builder
	)
	// preallocate memory
	valBuf.Grow(40)

	for ; ; a.idx++ {
		curToken = inputRune[a.idx]
		peekToken = inputRune[a.idx+1]

		switch curToken {
		case SPACE, TAB:
			if inQuote {
				valBuf.WriteRune(curToken)
			} else if !inQuote {
				a.ignoreSpaceTab(inputRune)
			}

		case DOUBLEQUOTE:
			valBuf.WriteRune(curToken)
			dc++
			if dc == 1 {
				inQuote = true
			}
			if dc == 2 {
				dc = 0
				inQuote = false
			}

		case BACKSLASH:
			valBuf.WriteRune(curToken)
			if inQuote{
				if peekToken == DOUBLEQUOTE {
					dc--
				}
			}

		case SLASH:
			if inQuote {
				valBuf.WriteRune(curToken)
			} else if !inQuote && peekToken == SLASH {
				a.ignoreComments(inputRune)
			}

		case COMMA:
			if inQuote {
				valBuf.WriteRune(curToken)
			} else if !inQuote {
				ss = valBuf.String()
				value = determineType(ss)
				valBuf.Reset()
				if _, ok := assembledMap[lineCount]; !ok {
					assembledMap[lineCount] = make(map[string]any, 1)
				}
				assembledMap[lineCount][key] = value
			}

		default:
			if inQuote {
				valBuf.WriteRune(curToken)
			}
		}
	}
}
