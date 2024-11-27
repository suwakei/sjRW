package internal

import (
	"strings"
	"strconv"

)
// キーやバリューに対応した空白を格納しておく場所を作ってもいい


func AssembleMap(str string) (assembledMap map[int]map[any]any) {
	var initMap map[int]map[any]any = make(map[int]map[any]any, 0)

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
		curToken rune
		peekToken rune
        r []rune = []rune(str)
		strLength int = len(r)
		doubleQuoteCnt int = 0
		sliceModeCount int = 0
		lineCount int = 0
		keyMode bool = true
		sliceMode bool = false
		firstLoop bool = true
		key string
		value string

		keyBuf strings.Builder
		valBuf strings.Builder
		sliceBuf strings.Builder
	)

	// preallocation of memory
	var keyBufMemoryNumber float32 = float32(strLength) * 0.2
	keyBuf.Grow(int(keyBufMemoryNumber))

	var valBufMemoryNumber float32 = float32(strLength) * 0.7
	valBuf.Grow(int(valBufMemoryNumber))

	var sliceBufMemoryNumber float32 = float32(strLength) * 0.1
	sliceBuf.Grow(int(sliceBufMemoryNumber))

	var runifiedStr []rune = make([]rune, 0, strLength)
	runifiedStr = r


	for idx := range runifiedStr {
		L:
		if sliceMode {
			if idx <= idx + sliceModeCount {
				continue
			}
			sliceMode = false
		}

		curToken = runifiedStr[idx]

		// index out of lengthを防ぐため
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

		// 最後のトークンの時
		if (idx + 1 == strLength) && (curToken == RBRACE || curToken == RBRACKET){
			lineCount += 1
			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[any]any, 0)
			}
			initMap[lineCount][string(curToken)] = ""
			keyBuf.Reset()
			valBuf.Reset()
			break
		}


		switch curToken {

		case SPACE, TAB:
			if keyMode {
				keyBuf.WriteRune(curToken)
			}
			if !keyMode {
				valBuf.WriteRune(curToken)
			}

		case DOUBLEQUOTE:
			if keyMode {
				keyBuf.WriteRune(curToken)
				doubleQuoteCnt += 1
			}
			if !keyMode {
				valBuf.WriteRune(curToken)
			}
			
			if doubleQuoteCnt == 2 && peekToken == COLON {
				key = keyBuf.String()
				keyMode = false
				continue
			}

		case LBRACKET:
			if keyMode {
				keyBuf.WriteRune(curToken)
			}

			if !keyMode {
				if doubleQuoteCnt == 0 {
					sliceMode = true
					sliceModeCount = 0
					var (
						tempSlice[]any = make([]any, 0, strLength / 10)
						dc int = 0
						lineCountBuf int = lineCount
						tempRune rune
						peekTempRune rune
						s string // sliceBufのstringを格納する
					)
					
					for i := idx + 1; i < strLength; i++ {
						sliceModeCount += 1
						tempRune = rune(runifiedStr[i])
						switch tempRune {
						case DOUBLEQUOTE:
							dc += 1
							// "が二個以上カウントされないようにする
							if peekTempRune = rune(runifiedStr[i + 1]);peekTempRune != COMMA && dc == 2{
								dc -= 1
							}
							sliceBuf.WriteRune(tempRune)

						case COMMA:
							s = sliceBuf.String()
							if num, err := strconv.Atoi(s); dc < 2 && err == nil {
								dc = 0
								sliceBuf.Reset()
								tempSlice = append(tempSlice, num)
								continue
							}
							dc = 0
							sliceBuf.Reset()
							tempSlice = append(tempSlice, s)

						case lnTOKEN:
							sliceBuf.WriteRune(tempRune)
							s = sliceBuf.String()
							// slicebufが数値かどうか判定
							if num, err := strconv.Atoi(s); dc < 2 && err == nil {
								dc = 0
								sliceBuf.Reset()
								tempSlice = append(tempSlice, num)
								lineCount += 1
								continue
							}
							dc = 0
							tempSlice = append(tempSlice, s)
							lineCount += 1

						case RBRACKET:
							if dc == 0 {
								initMap[lineCountBuf][key] = tempSlice
								lineCount += 1
								goto L
							}
							sliceBuf.WriteRune(tempRune)

						case SPACE, TAB:
							if dc != 1 {
								continue
							}
							sliceBuf.WriteRune(tempRune)

						default:
							sliceBuf.WriteRune(tempRune)
				}
				}
			}
		}

		case lnTOKEN:
			lineCount += 1

			valBuf.WriteRune(curToken)
			value = valBuf.String()

			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[any]any, 0)
			}

			initMap[lineCount][key] = value
			keyBuf.Reset()
			valBuf.Reset()
			keyMode = true
		
		case COLON:
			if doubleQuoteCnt == 2 {
				doubleQuoteCnt = 0
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
	assembledMap = initMap
	return assembledMap
}