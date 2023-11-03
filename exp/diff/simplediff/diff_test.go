package simplediff

import (
	"testing"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	const before = `package main

import (
	"fmt"
	"os"
)

func main() {
	world := "world"
	// hello
	fmt.Printf("Hello, %s!", world)
	os.Exit(0)
}
`
	const after = `package main

import (
	"fmt"
	"os"
)

func main() {
	name := "world"
	// hello
	fmt.Printf("Hello, %s!", name)
	os.Exit(0)
}
`

	const golden = ` package main
 
 import (
 	"fmt"
 	"os"
 )
 
 func main() {
-	world := "world"
+	name := "world"
 	// hello
-	fmt.Printf("Hello, %s!", world)
+	fmt.Printf("Hello, %s!", name)
 	os.Exit(0)
 }
 
`

	diffString := Diff(before, after).String()
	if expected, actual := golden, diffString; expected != actual {
		t.Errorf("❌: expected != actual:\n---EXPECTED\n+++ACTUAL\n%s", Diff(expected, actual))
	}
}
