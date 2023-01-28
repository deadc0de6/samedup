/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"fmt"
	"strings"

	"github.com/deadc0de6/samedup/internal/db"

	"github.com/TwiN/go-color"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type treeWriter struct{}

func (w treeWriter) Write(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, nil
	}
	out := string(data)
	out = strings.TrimSuffix(out, "\n")
	out = strings.TrimSuffix(out, "\n")
	fmt.Println(string(out))
	return len(out), nil
}

func printTree(dups *db.Duplicates) {
	for _, dup := range dups.Duplicates {
		var tree pterm.LeveledList
		item := pterm.LeveledListItem{
			Level: 0,
			Text:  color.InGreen(fmt.Sprintf("%s (%d)", dup.Key, len(dup.Nodes))),
		}
		tree = append(tree, item)
		for _, entry := range dup.Nodes {
			txt := fmt.Sprintf("%s (size:%s)", entry.AbsPath, SizeToHuman(entry.Size))
			item := pterm.LeveledListItem{
				Level: 1,
				Text:  txt,
			}
			tree = append(tree, item)
		}
		root := putils.TreeFromLeveledList(tree)
		pterm.DefaultTree.WithRoot(root).WithWriter(treeWriter{}).Render()
	}
}
