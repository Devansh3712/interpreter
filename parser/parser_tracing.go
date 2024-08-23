package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

func indentLevel() string  { return strings.Repeat("\t", traceLevel-1) }
func tracePrint(fs string) { fmt.Printf("%s%s\n", indentLevel(), fs) }
func incrementIndent()     { traceLevel += 1 }
func decrementIndent()     { traceLevel -= 1 }

func trace(message string) string {
	incrementIndent()
	tracePrint("BEGIN " + message)
	return message
}

func untrace(message string) {
	tracePrint("END " + message)
	decrementIndent()
}
