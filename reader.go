package sjrw

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/suwakei/sjrw/internal"
)

type SjReader struct {
}


//ReadAsStr returns json content as string
func (SjReader) ReadAsStrFrom(readFile io.Reader) (contentAsStr string, err error) {
	var jsonByte []byte = make([]byte, 200)

	// if _, err := os.Stat(readFilePath); err != nil {
	// 	log.Fatalf("this path is not exist \"%s\"", readFilePath)
	// }

	// if filepath.Ext(readFilePath) != ".json" {
	// 	return "", errors.New("read file is not json file")
	// }

	// f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

	// if err != nil {
	// 	log.Fatalf("could not open file \"%s\"", readFilePath)
	// }

	// defer func() {
	// 	err := f.Close()
	// 	if err != nil {
	// 		log.Fatalf("could not close file \"%s\"", readFilePath)
	// 	}
	// }()

	reader := bufio.NewReaderSize(readFile, 24 * 1024)

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


func (SjReader) ReadAsByteFrom(readFile io.Reader) (contentAsByte []byte, err error) {
	var jsonByte []byte = make([]byte, 200)

	// if _, err := os.Stat(readFilePath); err != nil {
	// 	log.Fatalf("this path is not exist \"%s\"", readFilePath)
	// }

	// if filepath.Ext(readFilePath) != ".json" {
	// 	return []byte(""), errors.New("read file is not json file")
	// }

	// f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

	// if err != nil {
	// 	log.Fatalf("could not open file \"%s\"", readFilePath)
	// }

	// defer func() {
	// 	err := f.Close()
	// 	if err != nil {
	// 		log.Fatalf("could not close file \"%s\"", readFilePath)
	// 	}
	// }()

	reader := bufio.NewReaderSize(readFile, 24 * 1024)
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
func (SjReader) ReadAsMapFrom(readFile io.Reader) (contentAsMap map[int]map[string]any, err error) {
	var (
		jsonByte []byte = make([]byte, 200)
	)

	// if _, err := os.Stat(readFilePath); err != nil {
	// 	log.Fatalf("this path is not exist \"%s\"", readFilePath)
	// }

	// if filepath.Ext(readFilePath) != ".json" {
	// 	return nil, errors.New("read file is not json file")
	// }

	// f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

	// if err != nil {
	// 	log.Fatalf("could not open file \"%s\"", readFilePath)
	// }

	// defer func() {
	// 	err := f.Close()
	// 	if err != nil {
	// 		log.Fatalf("could not close file \"%s\"", readFilePath)
	// 	}
	// }()

	reader := bufio.NewReaderSize(readFile, 24 * 1024)

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
	trimedByte := bytes.TrimSpace(jsonByte)
	st := string(trimedByte)
	fmt.Println(st)
	rs := []rune(st)

	contentAsMap = internal.AssembleMap(rs)

	return contentAsMap, err
}
