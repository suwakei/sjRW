package internal

import (
	"strings"
)

// returnObj returns map[string]any
func returnObj(idx, lineCount uint, inputRune []rune) (returnedIdx, returnedLineCount uint, rm map[string]any) {
	var (
		curToken  rune // The target token.
		peekToken rune // The token for confirmation of next character.

		dc        uint8        // Counter for number of ".
		keyMode   bool  = true // If true, mode which read json key. if false read json value.
		firstLoop bool  = true // First loop flag.

		keyBuf strings.Builder // When in "keyMode" is true, buf for accumulating key token.
		valBuf strings.Builder // When in "keyMode" is false, buf for accumulating value token.
		key    string          // The variable is for concatenated tokens stored in "keyBuf".
		value  any             = nil
	)

	// preallocation of memory
	rm = make(map[string]any, mapLength(idx, inputRune))
	keyBuf.Grow(20)
	valBuf.Grow(30)

	for ; ; idx++ {
		curToken = inputRune[idx]

		if firstLoop {
			rm[string(curToken)] = ""
			firstLoop = false
			continue
		}

		switch curToken {
		case SPACE, TAB:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
				} else if dc == 0 {
					idx = ignoreSpaceTab(idx, inputRune)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
				} else if dc == 0 {
					idx = ignoreSpaceTab(idx, inputRune)
				}
			}

		case DOUBLEQUOTE:
			dc++
			if keyMode {
				keyBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
					continue
				}
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if dc == 2 {
					dc = 0
					continue
				}
			}

		case BACKSLASH:
			if keyMode {
				keyBuf.WriteRune(curToken)
				if peekToken = inputRune[idx+1]; peekToken == DOUBLEQUOTE {
					dc--
					continue
				}
				continue
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
				if peekToken = inputRune[idx+1]; peekToken == DOUBLEQUOTE {
					dc--
					continue
				}
				continue
			}

		case SLASH:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					idx = ignoreComments(idx, inputRune)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					idx = ignoreComments(idx, inputRune)
					continue
				}
			}

		case COLON:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					key = keyBuf.String()
					keyBuf.Reset()
					keyMode = false
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				}
			}

		case LBRACE:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					rdx, rlc, rrs := returnObj(idx, lineCount, inputRune)
					idx = rdx
					lineCount = rlc
					value = rrs
					continue
				}
			}

		case RBRACE:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					if value != nil {
						rm[key] = value

					} else {
						ss := valBuf.String()
						valBuf.Reset()
						value = nil

						if ss != "" {
							value = determineType(ss)
							rm[key] = value

						} else if ss == "" {
							rm[key] = ss
						}
					}
					rm[string(curToken)] = ""
					returnedIdx = idx
					returnedLineCount = lineCount
					return returnedIdx, returnedLineCount, rm
				}
			}

		case LBRACKET:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					rdx, rlc, rrs := returnArr(idx, lineCount, inputRune)
					idx = rdx
					lineCount = rlc
					value = rrs
					valBuf.Reset()
					continue
				}
			}

		case RBRACKET:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				}
			}

		case COMMA:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					if value != nil {
						rm[key] = value
						value = nil

					} else {
						ss := valBuf.String()
						valBuf.Reset()

						if ss != "" {
							value = determineType(ss)
							rm[key] = value
							value = nil

						} else if ss == "" {
							rm[key] = ss
							value = nil
						}
					}

					keyMode = true
					continue
				}
			}

		case lrTOKEN:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					if peekToken = inputRune[idx+1]; peekToken == lnTOKEN {
						continue
					}
					lineCount++
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					if peekToken = inputRune[idx+1]; peekToken == lnTOKEN {
						continue
					}
					lineCount++
					continue
				}
			}

		case lnTOKEN:
			if keyMode {
				if dc > 0 {
					keyBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					lineCount++
					continue
				}
			}

			if !keyMode {
				if dc > 0 {
					valBuf.WriteRune(curToken)
					continue
				} else if dc == 0 {
					lineCount++
					continue
				}
			}

		default:
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		}
	}

}

func mapLength(idx uint, inputRune []rune) uint {
	var (
		curToken  rune
		peekToken rune
		dc        uint8
		mapLength uint
		lb        uint8
		rb        uint8
	)

	for ; lb != rb; idx++ {
		curToken = inputRune[idx]

		switch curToken {
		case DOUBLEQUOTE:
			dc++
			if dc == 2 {
				dc = 0
			}

		case BACKSLASH:
			if peekToken = inputRune[idx+1]; peekToken == DOUBLEQUOTE {
				dc--
			}

		case COLON:
			if dc == 0 {
				mapLength++
			}

		case LBRACE:
			if dc == 0 {
				lb++
			}

		case RBRACE:
			if dc == 0 {
				rb++
			}
		}
	}
	return mapLength + 2
}
