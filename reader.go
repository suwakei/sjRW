package sjrw

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)


type SjReader struct {
}


//ReadAsStr returns json content as string
func (sj *SjReader) ReadAsStrFrom(readFilePath string) (contentAsStr string, err error) {
	var jsonByte []byte = make([]byte, 0)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
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


//[]map[string]any
func (sj *SjReader) ReadAsMapFrom(readFilePath string) (contentAsMap map[int]map[string]any, err error) {
	var (
		jsonByte []byte = make([]byte, 0)
	)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
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
	// 行の順番を持ったマップとjsonのキーをキーとしたマップ
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

	r := []rune(str)
	var runifiedStr []rune = make([]rune, len(r))
	runifiedStr = r

	var (
		curToken rune
		peekToken rune
		strLength int = len(runifiedStr)
		doubleQuoteCnt int = 0
		lineCount int = 0
		key string
		value string

		strBuf strings.Builder
	)

	for idx := range runifiedStr {
		curToken = runifiedStr[idx]

		// index out of lengthを防ぐため
		if idx + 1 < strLength {
		peekToken = runifiedStr[idx + 1]
		}

		// 最後のトークンの時
		if idx + 1 == strLength {
			lineCount += 1
			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[string]any, 0)
			}
			initMap[lineCount][string(curToken)] = ""
			break
		}

		if curToken == SPACE || curToken == TAB {
			strBuf.WriteRune(curToken)
			continue
		}

		if curToken == DOUBLEQUOTE {
			strBuf.WriteRune(curToken)
			if doubleQuoteCnt == 2 && peekToken == COLON {
				key = strBuf.String()
				strBuf.Reset()
				doubleQuoteCnt = 0
				continue
			}
			doubleQuoteCnt += 1
			continue
		}

		if curToken == lnTOKEN {
			lineCount += 1
			strBuf.WriteRune(curToken)
			value = strBuf.String()

			if _, ok := initMap[lineCount]; !ok {
				initMap[lineCount] = make(map[string]any, 0)
			}
			initMap[lineCount][key] = value
			strBuf.Reset()
			continue
			}
		
		if curToken == COMMA && doubleQuoteCnt < 2 {
			continue
		}

		// どのトークンにも当てはまらない場合
		strBuf.WriteRune(curToken)
	}
	assembledMap = initMap
	return assembledMap
}