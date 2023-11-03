package simplediff

import (
	"strings"
)

// Separator is the separator of lines.
//
//nolint:gochecknoglobals
var Separator = "\n"

type (
	// DiffOperation represents a single operation in a diff.
	DiffOperation struct {
		Op   string // "+" for add, "-" for delete, " " for no change
		Text string
	}
	// DiffOperations represents a collection of DiffOperation.
	DiffOperations []DiffOperation
)

// diff returns a slice of DiffOp representing the diff of two slices.
//
//nolint:cyclop
func Diff(before, after string) DiffOperations {
	return diff(strings.Split(before, Separator), strings.Split(after, Separator))
}

//nolint:cyclop
func diff(a, b []string) []DiffOperation {
	m := len(a)
	n := len(b)
	diffs := []DiffOperation{}

	// Create a 2D slice to store the edit distance between slices
	edits := make([][]int, m+1)
	for i := range edits {
		edits[i] = make([]int, n+1)
	}

	// Fill the table
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			switch {
			case i == 0:
				edits[i][j] = j
			case j == 0:
				edits[i][j] = i
			case a[i-1] == b[j-1]:
				edits[i][j] = edits[i-1][j-1]
			default:
				edits[i][j] = min(edits[i-1][j]+1, edits[i][j-1]+1)
			}
		}
	}

	// Backtrack to find the diff
	for i, j := m, n; i > 0 || j > 0; {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			diffs = append([]DiffOperation{{" ", a[i-1]}}, diffs...)
			i--
			j--
		case j > 0 && (i == 0 || edits[i][j-1] <= edits[i-1][j]):
			diffs = append([]DiffOperation{{"+", b[j-1]}}, diffs...)
			j--
		case i > 0 && (j == 0 || edits[i][j-1] > edits[i-1][j]):
			diffs = append([]DiffOperation{{"-", a[i-1]}}, diffs...)
			i--
		}
	}

	return diffs
}

func (diffOps DiffOperations) String() string {
	var result strings.Builder
	for _, diff := range diffOps {
		result.WriteString(diff.Op + diff.Text + Separator)
	}
	return result.String()
}
