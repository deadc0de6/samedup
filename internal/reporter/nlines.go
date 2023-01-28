/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"fmt"

	"github.com/deadc0de6/samedup/internal/db"
)

func printNLines(dups *db.Duplicates) {
	for _, dup := range dups.Duplicates {
		for _, entry := range dup.Nodes {
			fmt.Printf("%s\n", entry.AbsPath)
		}
		fmt.Println("")
	}
}
