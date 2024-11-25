package sjrw

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	//"strconv"
	"strings"
)
// mapのキーの型はstringとは限らないからcomparable型にする
aa

type SjReader struct {
}


//ReadAsStr returns json content as string
func (sj *SjReader) ReadAsStrFrom(readFilePath string) (contentAsStr string, err error) {
	var jsonByte []byte = make([]byte, 0)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return "", errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0666)

	if err != nil {
		log.Fatalf("could not open file \"%s\"", readFilePath)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not close file \"%s\"", readFilePath)
		}
	}()

	reader := bufio.NewReaderSize(f, 24 * 1024)

	for {
		readByte, _, err := reader.ReadLine()
		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("could not read content")
		}

	}
	contentAsStr = strings.TrimSpace(string(jsonByte))
	return contentAsStr, err
}


func (sj *SjReader) ReadAsByteFrom(readFilePath string) (contentAsByte []byte, err error) {
	var jsonByte []byte = make([]byte, 0)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return []byte(""), errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0666)

	if err != nil {
		log.Fatalf("could not open file \"%s\"", readFilePath)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not close file \"%s\"", readFilePath)
		}
	}()

	reader := bufio.NewReaderSize(f, 24 * 1024)
	for {
		readByte, _, err := reader.ReadLine()

		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("could not read content")
		}
	}
	contentAsByte = jsonByte
	return contentAsByte, err
}


// 本当に一行ずつ取得とかではなくて配列またはオブジェクトがキーの中にあったらそれをすべて参照できるようにする
func (sj *SjReader) ReadAsMapFrom(readFilePath string) (contentAsMap map[int]map[string]any, err error) {
	var (
		jsonByte []byte = make([]byte, 0)
	)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return nil, errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0666)

	if err != nil {
		log.Fatalf("could not open file \"%s\"", readFilePath)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not close file \"%s\"", readFilePath)
		}
	}()

	reader := bufio.NewReaderSize(f, 24 * 1024)

	for {
		readByte, _, err := reader.ReadLine()

		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("could not read content")
		}
	}

	trimedString := strings.TrimSpace(string(jsonByte))

	contentAsMap = assembleMap(trimedString)

	return contentAsMap, err
}


func assembleMap(str string) (assembledMap map[int]map[string]any) {
	var initMap map[int]map[string]any = make(map[int]map[string]any, 0)

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

	var runifiedStr []rune = make([]rune, strLength)
	runifiedStr = r


	for idx := range runifiedStr {
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
				initMap[lineCount] = make(map[string]any, 0)
			}
			initMap[lineCount][string(curToken)] = ""
			keyBuf.Reset()
			valBuf.Reset()
			break
		}

		if sliceMode {
			if idx <= idx + sliceModeCount {
				continue
			}
			sliceMode = false
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
				fmt.Println("!K", idx)
				if doubleQuoteCnt == 0 {
					sliceMode = true
					sliceModeCount = 0
					var (
						tempSlice[]any = make([]any, strLength / 8)
						dc int = 0
						lineCountBuf int = lineCount
						num int
						tempRune rune
					)
					
					for i := idx + 1; i < strLength; i++ {
						fmt.Println("s", i)
						sliceModeCount += 1
						tempRune = rune(runifiedStr[i])
						switch tempRune {
						case DOUBLEQUOTE:
							dc += 1
							sliceBuf.WriteRune(tempRune)

						case COMMA:
							s := sliceBuf.String()
							if dc < 2 {
								num, _ = strconv.Atoi(s)
								tempSlice = append(tempSlice, num)
								dc = 0
								continue
							}
							dc = 0
							tempSlice = append(tempSlice, s)

						case lnTOKEN:
							s := sliceBuf.String()
							if dc < 2 {
								num, _ = strconv.Atoi(s)
								tempSlice = append(tempSlice, num)
								dc = 0
								continue
							}
							dc = 0
							tempSlice = append(tempSlice, s)
							lineCount += 1

						case RBRACKET:
							if dc < 0 {
								initMap[lineCountBuf][key] = tempSlice
								break
							}
							sliceBuf.WriteRune(tempRune)

						default:
							sliceBuf.WriteRune(tempRune)
				}
				}
				fmt.Println(tempSlice)
			}
		}

		case lnTOKEN:
			lineCount += 1

			valBuf.WriteRune(curToken)
			value = valBuf.String()

			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[string]any, 0)
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
