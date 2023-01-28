/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package db

import (
	"encoding/json"
	"fmt"
	"os"
	"samedup/internal/checksummer"
	"samedup/internal/logger"
	"samedup/internal/multmap"
	"samedup/internal/node"
	"time"

	"github.com/pterm/pterm"
)

var (
	// SortByName sort duplicates by name
	SortByName = "name"
	// SortByModTime sort duplicates by modification time
	SortByModTime = "modtime"
)

// Database samedup main database
type Database struct {
	version     string
	bySize      multmap.MultMap[uint64]
	newFileChan chan (*node.Node)
	resultChan  chan (*Duplicates)
	workerPool  *checksummer.Pool
	cnt         int64
	spinner     *pterm.SpinnerPrinter
	hashMethod  int
	sortMode    string
	flogger     *logger.FileLogger
}

// Duplicates found duplicates
type Duplicates struct {
	Duplicates  []*Duplicate
	TotalWasted uint64
}

// Duplicate list of nodes sharing the same key (checksum, name, size, etc)
type Duplicate struct {
	Nodes  []*node.Node
	Key    string
	Wasted uint64
}

type header struct {
	Version string       `json:"version"`
	Created int64        `json:"created"`
	Entries []*node.Node `json:"files"`
}

func (d *Database) addBySize(n *node.Node) {
	entries := d.bySize.Add(n.Size, n)
	length := len(entries)
	if length == 2 {
		// checksum first and this one
		for _, entry := range entries {
			subn := entry.(*node.Node)
			d.workerPool.AddJob(subn, d.hashMethod, -1)
		}
	} else if length > 2 {
		// checksum this one
		d.workerPool.AddJob(n, d.hashMethod, -1)
	}
}

func (d *Database) findDuplicates() multmap.MultMap[string] {
	if len(d.bySize) < 2 {
		// no(t enough) files
		return nil
	}

	// aggregate all with the same fast checksum
	if d.spinner != nil {
		d.spinner.UpdateText(fmt.Sprintf("aggregate by size (%d)...", len(d.bySize)))
	}

	potentials := multmap.New[string]()
	for _, lst := range d.bySize {
		if len(lst) < 2 {
			// skip we need at least two nodes with same size
			continue
		}
		first := lst[0].(*node.Node)
		potentials.Add(first.Checksum, lst...)
	}

	// wait for all checksums to be done
	d.workerPool.Wait()

	// aggregate all with the same checksum
	all := potentials.GetAllValues()
	if d.spinner != nil {
		d.spinner.UpdateText(fmt.Sprintf("aggregate by checksum (%d)...", len(all)))
	}
	actuals := multmap.New[string]()
	for _, v := range all {
		n := v.(*node.Node)
		if len(n.Checksum) < 1 {
			continue
		}
		actuals.Add(n.Checksum, n)
	}

	logger.Debugf("done building duplicates")
	if d.spinner != nil {
		d.spinner.UpdateText("done")
	}
	return actuals
}

func (d *Database) index() {
	for n := range d.newFileChan {
		logger.Debugf("adding node: %#v", n)
		d.addBySize(n)
		d.cnt++
	}

	// we get here once the chan is closed
	// that is when all files have been indexed
	d.resultChan <- d.getDuplicatesByContent()
	close(d.resultChan)
}

func (d *Database) sort(lst []*node.Node) []*node.Node {
	switch d.sortMode {
	case SortByName:
		return node.SortNodesByBaseName(lst)
	case SortByModTime:
		return node.SortNodesByModTime(lst)
	}
	logger.Warnf("no such sorting mode: %s", d.sortMode)
	return node.SortNodesByBaseName(lst)
}

func (d *Database) getDuplicatesByContent() *Duplicates {
	potentials := d.findDuplicates()

	// construct the duplicates list
	if d.spinner != nil {
		d.spinner.UpdateText("compile duplicates list...")
	}
	var dups Duplicates
	for sum, lst := range potentials {
		if len(lst) < 2 {
			continue
		}

		logger.Debugf("%d with sum %s", len(lst), sum)

		var sublst []*node.Node
		var size uint64
		for _, entry := range lst {
			n := entry.(*node.Node)
			size = n.Size
			sublst = append(sublst, n)
		}
		dup := &Duplicate{
			Key:    sum,
			Nodes:  d.sort(sublst),
			Wasted: uint64(len(lst)-1) * size,
		}
		dups.Duplicates = append(dups.Duplicates, dup)
		dups.TotalWasted += dup.Wasted
	}
	return &dups
}

// Serialize serialize database data
func (d *Database) Serialize() ([]byte, error) {
	var nodes []*node.Node
	for _, entry := range d.bySize.GetAllValues() {
		n := entry.(*node.Node)
		nodes = append(nodes, n)
	}

	hdr := header{
		Version: d.version,
		Created: time.Now().Unix(),
		Entries: nodes,
	}

	content, err := json.Marshal(hdr)
	if err != nil {
		return nil, err
	}

	return content, err
}

// Save save database to file
func (d *Database) Save(path string) error {
	content, err := d.Serialize()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Start start the database for indexing and processing nodes
func (d *Database) Start() {
	d.workerPool.Start()
	go d.index()
}

// NewDatabase create a new database
func NewDatabase(version string, newFileChan chan (*node.Node), resultChan chan (*Duplicates), hashMethod int, spinner *pterm.SpinnerPrinter, sortMode string, flogger *logger.FileLogger) *Database {
	db := &Database{
		version:     version,
		newFileChan: newFileChan,
		resultChan:  resultChan,
		bySize:      multmap.New[uint64](),
		workerPool:  checksummer.NewPool(flogger),
		spinner:     spinner,
		hashMethod:  hashMethod,
		sortMode:    sortMode,
		flogger:     flogger,
	}
	return db
}
