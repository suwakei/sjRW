package sjrw

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
 "unicode/utf8"
 "bytes"
	"path/filepath"
	"strings"
	"github.com/suwakei/sjrw/internal"
)

type SjReader struct {
}


//ReadAsStr returns json content as string
func (SjReader) ReadAsStrFrom(readFilePath string) (contentAsStr string, err error) {
	var jsonByte []byte = make([]byte, 0)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return "", errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

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


func (SjReader) ReadAsByteFrom(readFilePath string) (contentAsByte []byte, err error) {
	var jsonByte []byte = make([]byte, 0)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return []byte(""), errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

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
func (SjReader) ReadAsMapFrom(readFilePath string) (contentAsMap map[int]map[string]any, err error) {
	var (
		jsonByte []byte = make([]byte, 0)
	)

	if _, err := os.Stat(readFilePath); err != nil {
		log.Fatalf("this path is not exist \"%s\"", readFilePath)
	}

	if filepath.Ext(readFilePath) != ".json" {
		return nil, errors.New("read file is not json file")
	}

	f, err := os.OpenFile(readFilePath, os.O_RDONLY, 0o666)

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

// stringにしないでbyte型からutf8.EncodeRune使ってruneにする
//trimspaceはbytesパッケージ使ってstring経由しないでやる
	trimedString := strings.TrimSpace(string(jsonByte))

	contentAsMap = internal.AssembleMap(trimedString)

	return contentAsMap, err
}
