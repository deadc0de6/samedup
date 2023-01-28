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

func printOneLine(dups *db.Duplicates) {
	for _, dup := range dups.Duplicates {
		var out []string
		for _, entry := range dup.Nodes {
			out = append(out, entry.AbsPath)
		}
		fmt.Println(strings.Join(out, " "))
	}
}
