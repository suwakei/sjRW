package internal

import (
	"strconv"
	"strings"
)

// returnSliceOrMapAndCount returns *RV.
func returnArr(commonIdx, commonLineCount *uint, inputRune []rune) *RV {
	var (
		commanum = commaNum(*commonIdx, inputRune)// the number of commas, uses for allocating memory of "tempSlice"
		curRuneLength uint = uint(len(inputRune[*commonIdx:]))// The length of input rune slice
		dc uint8 = 0 // This variable is same as variable "doubleQuoteCount" in the function "AssembleMap"
		sliceRune rune // This variable works just like the variable "curToken" in the function "AssembleMap"
		peekSliceRune rune // This variable works just like the variable "peekToken" in the function "AssembleMap"
		sliceBuf strings.Builder // When in "sliceMode" is true, buf for storing slice token.
		internalIdx uint
		internalLineCount uint
		ss string // The variable for concatnated tokens stored in "sliceBuf".
		rv *RV = new(RV)
	)
	// Store *commonIdx and *commonLineCount to avoid calling type *uint again and again
	internalIdx = *commonIdx
	internalLineCount = *commonLineCount

	// preallocate memory
	var sliceBufMemoryNumber float32 = float32(curRuneLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))
	rv.rs = make([]any, 0, commanum)

	for ; internalIdx < curRuneLength; internalIdx++ {
		sliceRune = rune(inputRune[internalIdx])
		switch sliceRune {

		case SPACE, TAB:
			if dc == 0 {
				continue
			}
			sliceBuf.WriteRune(sliceRune)

		case DOUBLEQUOTE:
			sliceBuf.WriteRune(sliceRune)
			dc++
			if dc == 2{
				dc = 0
			}

		case BACKSLASH:
			sliceBuf.WriteRune(sliceRune)
			if peekSliceRune = rune(inputRune[internalIdx + 1]); dc == 1 && peekSliceRune == DOUBLEQUOTE{
				dc--
			}

		case COMMA:
			if dc > 0 {
				sliceBuf.WriteRune(sliceRune)
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
			}

		case lrTOKEN:
			if dc > 0 {
				sliceBuf.WriteRune(sliceRune)
				}

			if dc == 0 {
				if peekSliceRune = rune(inputRune[internalIdx + 1]); peekSliceRune == lnTOKEN {
					continue
				}

				if peekSliceRune = rune(inputRune[internalIdx + 1]); peekSliceRune != lnTOKEN {
					internalLineCount++
				}
			}

		case lnTOKEN:
			if dc > 0 {
			sliceBuf.WriteRune(sliceRune)
			}

			if dc == 0 {
			internalLineCount++
			}

		case RBRACKET:
			// When the token is last
			if dc == 0 {
				internalLineCount++
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

				*commonIdx = internalIdx
				*commonLineCount = internalLineCount
				return &RV{
					rs: rv.rs,
				}
			}

			if dc > 0 {
				sliceBuf.WriteRune(sliceRune)
				continue
			}

		default:
			sliceBuf.WriteRune(sliceRune)
		}
	}
	return &RV{}
}