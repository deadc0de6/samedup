/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package node

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Node a file
type Node struct {
	BaseName string `json:"basename"`
	AbsPath  string `json:"abspath"`
	Checksum string `json:"checksum"`
	Size     uint64 `json:"size"`
	Mode     string `json:"mode"`
	ModTime  int64  `json:"modification"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// SortNodesByBaseName sort nodes by basename
func SortNodesByBaseName(nodes []*Node) []*Node {
	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]
		return left.BaseName < right.BaseName
	})
	return nodes
}

// SortNodesByFullName sort nodes by fullname
func SortNodesByFullName(nodes []*Node) []*Node {
	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]
		return left.AbsPath < right.AbsPath
	})
	return nodes
}

// SortNodesByModTime sort nodes by modification time
func SortNodesByModTime(nodes []*Node) []*Node {
	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]
		return left.ModTime < right.ModTime
	})
	return nodes
}

// NewNode create a new entry node
func NewNode(path string, size uint64, mode string, modTime time.Time) (*Node, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if !fileExists(path) {
		return nil, fmt.Errorf("no such file: %s", path)
	}
	n := &Node{
		BaseName: filepath.Base(path),
		AbsPath:  path,
		Size:     size,
		Mode:     mode,
		ModTime:  modTime.Unix(),
	}
	return n, nil
}
