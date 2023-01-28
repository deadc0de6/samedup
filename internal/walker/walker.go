/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package walker

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"samedup/internal/logger"
	"samedup/internal/node"
)

func isMatch(path string, ignores []*regexp.Regexp) bool {
	for _, patt := range ignores {
		matched := patt.MatchString(path)
		if matched {
			return matched
		}
	}
	return false
}

// Walk walks the filesystem and pushes found files to the channel
// it returns the logfile, total number of found files, nb skipped, nb error
func Walk(rootPath string, newFileChan chan (*node.Node), filters []*regexp.Regexp, ignores []*regexp.Regexp, ignoreZero bool, flogger *logger.FileLogger) (string, int64, int64, int64, error) {
	var cntPushed int64
	var cntSkipped int64
	var cntErr int64
	logger.Debugf("walking %s", rootPath)

	logPath := flogger.GetPath()
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			flogger.Error(fmt.Errorf("reading %s: %v", path, err))
			cntErr++
			return nil
		}
		if info == nil {
			flogger.Errorf("no file info for %s", path)
			cntErr++
			return nil
		}
		baseName := filepath.Base(path)
		if isMatch(baseName, ignores) {
			flogger.Warnf("skip ignored: %s", path)
			cntSkipped++
			return nil
		}
		if len(filters) > 0 {
			if !isMatch(path, filters) {
				flogger.Warnf("skip filtered: %s", path)
				cntSkipped++
				return nil
			}
		}

		if path == rootPath {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !info.Mode().IsRegular() {
			flogger.Warnf("skipping non regular file %s", path)
			cntSkipped++
			return nil
		}

		// process the file
		size := info.Size()
		if size < 1 && ignoreZero {
			logger.Debugf("skip zero size %s", path)
			cntSkipped++
			return nil
		}
		logger.Debugf("found new file: %s", path)
		n, err := node.NewNode(path, uint64(size), info.Mode().String(), info.ModTime())
		if err != nil {
			flogger.Error(fmt.Errorf("creating node for %s: %v", path, err))
			cntErr++
			return nil
		}
		newFileChan <- n
		cntPushed++
		return nil
	})

	return logPath, cntPushed, cntSkipped, cntErr, err
}
