package internal

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type pair struct{ x, y int }

// Diff returns result of comparing "from" and "to". two dimention map "mapFromDiff" of "rm" line list and "add" line list
// and diff log as []byte.
// "mapLog" used for debug.
func Diff(fromName string, from []byte, toName string, to []byte) (mapFromDiff map[string]map[int]string, mapLog []byte) {
	if bytes.Equal(from, to) {
		fmt.Printf("%s and %s are the same value", fromName, toName)
		return nil, nil
	}
	fmt.Println(fromName, toName)

	x := lines(from)
	y := lines(to)

	// Print diff header.
	var out bytes.Buffer
	fmt.Fprintf(&out, "diff '%s' '%s'\n", fromName, toName)

	var (
		done  pair     // printed up to x[:done.x] and y[:done.y]
		chunk pair     // start lines of current chunk
		count pair     // number of lines from each side in current chunk
		ctext []string // lines for current chunk
	)

	const spl string = "------------------------------------------------------\n"

	for _, m := range tgs(x, y) {
		if m.x < done.x {
			// Already handled scanning forward from earlier match.
			continue
		}

		// Expand matching lines as far as possible,
		// establishing that x[start.x:end.x] == y[start.y:end.y].
		// Note that on the first (or last) iteration we may (or definitely do)
		// have an empty match: start.x==end.x and start.y==end.y.
		start := m
		for start.x > done.x && start.y > done.y && x[start.x-1] == y[start.y-1] {
			start.x--
			start.y--
		}

		end := m

		for end.x < len(x) && end.y < len(y) && x[end.x] == y[end.y] {
			end.x++
			end.y++
		}

		// Emit the mismatched lines before start into this chunk.
		// (No effect on first sentinel iteration, when start = {0,0}.)
		for _, s := range x[done.x:start.x] {
			ctext = append(ctext, "-"+s)
			count.x++
		}

		for _, s := range y[done.y:start.y] {
			ctext = append(ctext, "+"+s)
			count.y++
		}

		// If we're not at EOF and have too few common lines,
		// the chunk includes all the common lines and continues.
		const C = 3 // number of context lines

		if (end.x < len(x) || end.y < len(y)) &&
			(end.x-start.x < C || (len(ctext) > 0 && end.x-start.x < 2*C)) {

			for _, s := range x[start.x:end.x] {
				ctext = append(ctext, " "+s)
				count.x++
				count.y++
			}

			done = end
			continue
		}

		// End chunk with common lines for context.
		if len(ctext) > 0 {
			n := end.x - start.x

			if n > C {
				n = C
			}

			for _, s := range x[start.x : start.x+n] {
				ctext = append(ctext, " "+s)
				count.x++
				count.y++
			}

			done = pair{start.x + n, start.y + n}

			// Format and emit chunk.
			// Convert line numbers to 1-indexed.
			// Special case: empty file shows up as 0,0 not 1,0.
			if count.x > 0 {
				chunk.x++
			}

			if count.y > 0 {
				chunk.y++
			}

			fmt.Fprintf(&out, "@@ '%s' %d,%d | '%s' %d,%d @@\n", fromName, chunk.x+2, count.x, toName, chunk.y+3, count.y)

			out.WriteString(fmt.Sprintf("Lines being removed or added from [%s]\n", fromName))
			out.WriteString(spl)

			for _, s := range ctext {
				out.WriteString(s)
			}

			count.x = 0
			count.y = 0
			ctext = ctext[:0]
		}

		// If we reached EOF, we're done.
		if end.x >= len(x) && end.y >= len(y) {
			break
		}

		// Otherwise start a new chunk.
		chunk = pair{end.x - C, end.y - C}

		for _, s := range x[chunk.x:end.x] {
			ctext = append(ctext, " "+s)
			count.x++
			count.y++
		}

		done = end
	}

	var editMap map[string]map[int]string = make(map[string]map[int]string)
	// prepare two dimenton map for return value
	if _, ok := editMap["rm"]; !ok {
		editMap["rm"] = make(map[int]string)
	}

	if _, ok := editMap["add"]; !ok {
		editMap["add"] = make(map[int]string)
	}

	str := out.String()
	trim := strings.TrimSpace(str)
	ar := strings.Split(trim, spl)[1:]

	var strarr []string
	for i := range len(ar) {
		strarr = strings.Split(ar[i], "\n")
	}

	for idx, s := range strarr {
		flag := strings.Split(s, "")[0:1][0]
		if flag == "+" {
			editMap["rm"][idx+1] = strings.Replace(s, "+", "", 1)
		}

		if flag == "-" {
			editMap["add"][idx+1] = strings.Replace(s, "-", "", 1)
		}
	}

	mapFromDiff = editMap

	mapLog = out.Bytes()
	return mapFromDiff, mapLog
}

// lines returns the lines in the file x, including newlines.
// If the file does not end in a newline, one is supplied
// along with a warning about the missing newline.
func lines(x []byte) []string {
	l := strings.SplitAfter(string(x), "\n")

	if l[len(l)-1] == "" {
		l = l[:len(l)-1]
	} else {
		l[len(l)-1] += "\n\\ No newline at end of file\n"
	}
	return l
}

func tgs(x, y []string) []pair {
	// Count the number of times each string appears in a and b.
	// We only care about 0, 1, many, counted as 0, -1, -2
	// for the x side and 0, -4, -8 for the y side.
	// Using negative numbers now lets us distinguish positive line numbers later.
	m := make(map[string]int)
	for _, s := range x {
		if c := m[s]; c > -2 {
			m[s] = c - 1
		}
	}
	for _, s := range y {
		if c := m[s]; c > -8 {
			m[s] = c - 4
		}
	}

	// Now unique strings can be identified by m[s] = -1+-4.
	//
	// Gather the indexes of those strings in x and y, building:
	//	xi[i] = increasing indexes of unique strings in x.
	//	yi[i] = increasing indexes of unique strings in y.
	//	inv[i] = index j such that x[xi[i]] = y[yi[j]].
	var xi, yi, inv []int
	for i, s := range y {
		if m[s] == -1+-4 {
			m[s] = len(yi)
			yi = append(yi, i)
		}
	}
	for i, s := range x {
		if j, ok := m[s]; ok && j >= 0 {
			xi = append(xi, i)
			inv = append(inv, j)
		}
	}

	// In those terms, A = J = inv and B = (0, n).
	// We add sentinel pairs {0,0}, and {len(x),len(y)}
	// to the returned sequence, to help the processing loop.
	J := inv
	n := len(xi)
	T := make([]int, n)
	L := make([]int, n)

	for i := range T {
		T[i] = n + 1
	}

	for i := 0; i < n; i++ {
		k := sort.Search(n, func(k int) bool {
			return T[k] >= J[i]
		})

		T[k] = J[i]
		L[i] = k + 1
	}

	k := 0

	for _, v := range L {
		if k < v {
			k = v
		}
	}

	seq := make([]pair, 2+k)
	seq[1+k] = pair{len(x), len(y)} // sentinel at end
	lastj := n

	for i := n - 1; i >= 0; i-- {
		if L[i] == k && J[i] < lastj {
			seq[k] = pair{xi[i], yi[J[i]]}
			k--
		}
	}

	seq[0] = pair{0, 0} // sentinel at start
	return seq
}
