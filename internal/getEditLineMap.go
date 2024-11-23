package internal

import (
	"sort"
	"strings"
)

// getEditLineMap returns a map with only the difference changed from the original string
func GetEditLineMap[targetType []byte | string] (str targetType, editMapFromDiff map[string]map[int]string) map[int]string {
	content := strings.TrimSpace(string(str))
	contentLines := strings.Split(content, "\n")

	indexes := make(map[int]string, 0)

	for i, line := range contentLines {
		indexes[i + 1] = line
	}

	rm := editMapFromDiff["rm"]
	RMNUMS := GetKey(rm)
	sort.Ints(RMNUMS)

	add := editMapFromDiff["add"]
	ADDNUMS := GetKey(add)
	sort.Ints(ADDNUMS)

	// remove
	for _, rmn := range RMNUMS {
		delete(indexes, rmn)
	}

	// add
	for _, addn := range ADDNUMS {
		indexes[addn] = editMapFromDiff["add"][addn]
	}

	return indexes
}

// GetKey returns keySlice of "m"
func GetKey(m map[int]string) (keySlice []int) {
	keySlice = make([]int, len(m))
	for k := range m {
		keySlice = append(keySlice, k)
	}
	return keySlice
}