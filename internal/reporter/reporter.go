/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"samedup/internal/db"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

var (
	// FormatTree tree output
	FormatTree = "tree"
	// FormatCSV CSV output
	FormatCSV = "csv"
	// FormatStairs stairs output
	FormatStairs = "stairs"
	// FormatLine one line per duplicate
	FormatLine = "oneline"
	// FormatNLines duplicates separated by a crlf like fdupes
	FormatNLines = "nlines"
	// FormatScript outputs a script to handle duplicates
	FormatScript = "script"
)

// SizeToHuman convert size to human readable string
func SizeToHuman(bytes uint64) string {
	str := humanize.Bytes(bytes)
	return strings.ReplaceAll(str, " ", "")
}

// IsValidFormat return true if format is among the supported ones
func IsValidFormat(format string) bool {
	switch format {
	case FormatTree:
		return true
	case FormatCSV:
		return true
	case FormatStairs:
		return true
	case FormatLine:
		return true
	case FormatNLines:
		return true
	case FormatScript:
		return true
	}
	return false
}

// GetFormats return available output formats
func GetFormats() []string {
	var out []string
	out = append(out, FormatTree)
	out = append(out, FormatCSV)
	out = append(out, FormatStairs)
	out = append(out, FormatLine)
	out = append(out, FormatNLines)
	out = append(out, FormatScript)
	return out
}

// Print print duplicates
func Print(format string, dups *db.Duplicates) {
	switch format {
	case FormatTree:
		printTree(dups)
	case FormatCSV:
		printCSV(dups)
	case FormatStairs:
		printStairs(dups)
	case FormatLine:
		printOneLine(dups)
	case FormatNLines:
		printNLines(dups)
	case FormatScript:
		printScript(dups)
	default:
		printTree(dups)
	}
}
