package sjrw

import (
	"fmt"
	"strings"
	"os"
)



func WriteAsStr() {

}


func WriteAsBytes() {

}


func removeLine(filename string, lineToRemove int) error {
	content, err := os.ReadFile(filename)
	
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	if len(lines) < lineToRemove {
		return fmt.Errorf("invalid line number")
	}
	
	lines = append(lines[:lineToRemove], lines[lineToRemove + 1:]...)
	output := strings.Join(lines, "\n")
	
	return os.WriteFile(filename, []byte(output), 0644)
}