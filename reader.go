package sjrw

import (
	"bufio"
	"bytes"
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
	if readFile != nil {

		reader := bufio.NewReaderSize(readFile, 24 * 1024)

		for {
			readByte, err := reader.ReadByte()
			jsonByte = append(jsonByte, readByte)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

		}
	}
	trimedByte = bytes.TrimSpace(jsonByte)
    contentAsStr = unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
	return contentAsStr, err
}


func (SjReader) ReadAsByteFrom(readFile io.Reader) (contentAsByte []byte, err error) {
	var jsonByte []byte = make([]byte, 200)
	if readFile != nil {
		reader := bufio.NewReaderSize(readFile, 24 * 1024)
		for {
			readByte, err := reader.ReadByte()
			jsonByte = append(jsonByte, readByte)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
	contentAsByte = bytes.TrimSpace(jsonByte)
	return contentAsByte, err
}


func (SjReader) ReadAsMapFrom(readFile io.Reader) (contentAsMap map[uint]map[string]any, err error) {
	var jsonRune []rune = make([]rune, 0, 800)

	if readFile != nil {
		reader := bufio.NewReaderSize(readFile, 24 * 1024)
		runeSlice := make([]rune, 0, 400)

		for {
			readRune, _, err := reader.ReadRune()

			jsonRune = append(jsonRune, append(runeSlice, readRune)...)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
	contentAsMap = internal.AssembleMap(jsonRune)
	return contentAsMap, err
}
