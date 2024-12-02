package sjrw

import (
	"bufio"
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


func (SjReader) ReadAsMapFrom(readFile io.Reader) (contentAsMap map[int]map[string]any, err error) {
	var (
		jsonRune []rune = make([]rune, 0, 800)
	)


	reader := bufio.NewReaderSize(readFile, 24 * 1024)
	runeSlice := make([]rune, 0, 400)

	for {
		readRune, _, err := reader.ReadRune()

		jsonRune = append(jsonRune, append(runeSlice, readRune)...)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("could not read content")
		}
	}
	fmt.Println(len(jsonRune))

	contentAsMap = internal.AssembleMap(jsonRune)

	return contentAsMap, err
}
