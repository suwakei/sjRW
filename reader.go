package sjrw

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)


type Sj struct {
}

//ReadAsStr returns json content as string
func (sj *Sj) ReadAsStrFrom(readFilePath string) (string, error) {
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
	return strings.TrimSpace(string(jsonByte)), err
}


func (sj *Sj) ReadAsBytesFrom(readFilePath string) ([]byte, error) {
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
	return jsonByte, err
}