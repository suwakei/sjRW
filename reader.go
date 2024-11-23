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
	var jsonByte []byte
	path := readFilePath

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("this path is not exist %s", path)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	reader := bufio.NewReaderSize(f, 24 * 1024)

	for {
		readByte, _, err := reader.ReadLine()
		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

	}
	contentAsStr = strings.TrimSpace(string(jsonByte))
	return contentAsStr, err
}


func (sj *SjReader) ReadAsByteFrom(readFilePath string) (contentAsByte []byte, err error) {
	var jsonByte []byte
	path := readFilePath

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("this path is not exist %s", path)
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	reader := bufio.NewReaderSize(f, 24 * 1024)
	for {
		readByte, _, err := reader.ReadLine()

		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}
	contentAsByte = jsonByte

	return contentAsByte, err
}


func (sj *SjReader) ReadAsMapFrom(readFilePath string) (contentAsMap map[string]interface{}, err error) {
	var jsonByte []byte
	path := readFilePath

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("this path is not exist %s", path)
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	reader := bufio.NewReaderSize(f, 24 * 1024)
	for {
		readByte, _, err := reader.ReadLine()

		jsonByte = append(jsonByte, append(readByte, []byte("\n")...)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}

	// 行の順番を持ったマップとjsonのキーをキーとしたマップ
}