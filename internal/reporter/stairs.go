/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"fmt"
	"samedup/internal/db"
)

func printStairs(dups *db.Duplicates) {
	for _, dup := range dups.Duplicates {
		fmt.Println(dup.Key)
		for _, entry := range dup.Nodes {
			fmt.Printf("  %s\n", entry.AbsPath)
		}
	}
}
