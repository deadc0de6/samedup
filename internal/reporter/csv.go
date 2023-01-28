/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"fmt"
	"strings"

	"github.com/deadc0de6/samedup/internal/db"
)

var (
	// CSVSeparator CSV separator
	CSVSeparator = ","
)

// format: <path>, <checksum>, <mode>, <size>
func printCSV(dups *db.Duplicates) {
	for _, dup := range dups.Duplicates {
		for _, entry := range dup.Nodes {
			var out []string
			out = append(out, entry.AbsPath)
			out = append(out, entry.Checksum)
			out = append(out, entry.Mode)
			sz := SizeToHuman(entry.Size)
			out = append(out, sz)
			fmt.Println(strings.Join(out, CSVSeparator))
		}
	}
}
