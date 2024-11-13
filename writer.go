package sjrw

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type SjWriter struct{
	JsonPath string
}


func WriteAsStr() {

}


func WriteAsBytes() {

}


// jsonファイルの差分があった時すべて変更するのではなく
// 変更した部分のみ書き込む処理
// 行の複数選択も可能
func (s *SjWriter) editLine(editMapFromDiff map[string]map[int]string)  {
	var jsonByte []byte
	path := s.JsonPath

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("this path is not exist %s", path)
	}

	f, err := os.OpenFile(s.JsonPath, os.O_RDWR, 0666)

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

	content := strings.TrimSpace(string(jsonByte))
	contentLines := strings.Split(content, "\n")

	indexes := make(map[int]string)
	for i, line := range contentLines {
		indexes[i + 1] = line
	}

	// とりあえずのテスト用二次元マップ
	if _, ok := editMapFromDiff["rm"]; !ok {
		editMapFromDiff["rm"] = make(map[int]string)
	}

	editMapFromDiff["rm"][1] = "git add ."
	editMapFromDiff["rm"][2] = "git commit -m "

	if _, ok := editMapFromDiff["add"]; !ok {
		editMapFromDiff["add"] = make(map[int]string)
	}

	editMapFromDiff["add"][1] = "git push origin main"
	editMapFromDiff["add"][2] = "git status"


	
	if len(indexes) <  len(editMapFromDiff["rm"]) + len(editMapFromDiff["add"]){
		log.Fatal("invalid line number")
	}
	
	// lines = append(lines[:editNumber], lines[editNumber + 1:]...)
	// output := strings.Join(lines, "\n")
	// fmt.Println(output)
	// return os.WriteFile(jsonPath, []byte(output), 0644)
}


