package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Devansh3712/interpreter/repl"
)

func main() {
	timestamp := time.Now()
	fmt.Printf("dev 0.0.1 (%s)\n", timestamp.Format("Jan 2 2006, 15:04:05"))
	repl.Start(os.Stdin, os.Stdout)
}
