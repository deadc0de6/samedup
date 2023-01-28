/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/deadc0de6/samedup/internal/checksummer"
	"github.com/deadc0de6/samedup/internal/db"
	"github.com/deadc0de6/samedup/internal/logger"
	"github.com/deadc0de6/samedup/internal/node"
	"github.com/deadc0de6/samedup/internal/reporter"
	"github.com/deadc0de6/samedup/internal/walker"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	myName  = "samedup"

	rootCmd = &cobra.Command{
		Use:     "samedup <dir> [<dir>...]",
		Short:   "samedup - duplicate finder",
		Long:    `Duplicate file finder`,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		PreRun:  preRunDebug,
		RunE:    samedup,
	}

	rootOptDebugMode   bool
	rootOptIgnoreEmpty = true
	rootOptIgnore      []string
	rootOptOutFormat   string
	rootOptQuiet       bool
	rootOptHashMethod  int
	rootOptFilter      []string
	rootOptSortMode    string
)

func init() {
	// env variables
	viper.SetEnvPrefix(strings.ToUpper(myName))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().BoolVarP(&rootOptDebugMode, "debug", "d", viper.GetBool("DEBUG"), "enable debug mode")
	rootCmd.PersistentFlags().BoolVarP(&rootOptIgnoreEmpty, "noempty", "z", false, "ignore empty files")
	rootCmd.PersistentFlags().StringSliceVarP(&rootOptIgnore, "ignore", "i", []string{}, "patterns to ignore")
	hlp := fmt.Sprintf("output format (%s)", strings.Join(reporter.GetFormats(), ","))
	rootCmd.PersistentFlags().StringVarP(&rootOptOutFormat, "output", "o", "tree", hlp)
	rootCmd.PersistentFlags().BoolVarP(&rootOptQuiet, "quiet", "q", false, "disable stats output")
	hlp = fmt.Sprintf("hash method to use (%s)", checksummer.GetHashMethods())
	rootCmd.PersistentFlags().IntVarP(&rootOptHashMethod, "hash", "H", checksummer.UseSHA1, hlp)
	rootCmd.PersistentFlags().StringSliceVarP(&rootOptFilter, "filter", "f", []string{}, "patterns to filter")
	hlp = fmt.Sprintf("sort duplicates by \"%s\" or \"%s\"", db.SortByName, db.SortByModTime)
	rootCmd.PersistentFlags().StringVarP(&rootOptSortMode, "sort", "s", db.SortByName, hlp)
}

func preRunDebug(ccmd *cobra.Command, args []string) {
	if rootOptDebugMode {
		logger.DebugMode = true
	}
}

func samedup(ccmd *cobra.Command, args []string) error {
	if !reporter.IsValidFormat(rootOptOutFormat) {
		logger.Fatalf("no such format: %s", rootOptOutFormat)
	}
	if rootOptOutFormat == reporter.FormatCSV || rootOptOutFormat == reporter.FormatScript {
		// disable stats for csv
		rootOptQuiet = true
	}

	// spinner
	var spinner *pterm.SpinnerPrinter
	if !rootOptQuiet {
		spinner = pterm.DefaultSpinner.WithRemoveWhenDone(true)
		spinner.Sequence = []string{` ⠋ `, ` ⠙ `, ` ⠹ `, ` ⠸ `, ` ⠼ `, ` ⠴ `, ` ⠦ `, ` ⠧ `, ` ⠇ `, ` ⠏ `}
		spinner.ShowTimer = true
	}

	// create the channels
	newFileChan := make(chan *node.Node)
	resultChan := make(chan *db.Duplicates)

	// file logger for warnings and errors
	flogger := logger.NewFileLogger()
	defer flogger.Close()

	// create the db
	database := db.NewDatabase(version, newFileChan, resultChan, rootOptHashMethod, spinner, rootOptSortMode, flogger)

	// build ignore pattern
	var ignPatterns []*regexp.Regexp
	for _, ign := range rootOptIgnore {
		re, err := regexp.Compile(ign)
		if err != nil {
			logger.Fatal(err)
		}
		ignPatterns = append(ignPatterns, re)
	}

	// build filter patterns
	var filterPatterns []*regexp.Regexp
	for _, filt := range rootOptFilter {
		re, err := regexp.Compile(filt)
		if err != nil {
			logger.Fatal(err)
		}
		filterPatterns = append(filterPatterns, re)
	}

	t0 := time.Now()

	// index all
	if spinner != nil {
		var err error
		spinner, err = spinner.Start("indexing...")
		if err != nil {
			logger.Error(err)
		}
	}
	database.Start()

	// process each argument
	var cnt int64
	var stats []string
	for _, arg := range args {
		path, err := filepath.Abs(arg)
		if err != nil {
			logger.Error(err)
			continue
		}

		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			logger.Error(err)
			continue
		}

		info, err := os.Lstat(path)
		if err != nil {
			logger.Error(err)
			continue
		}

		if !info.IsDir() {
			logger.Errorf("ignoring non directory %s", path)
			continue
		}

		// walk
		logger.Debugf("processing %s (dir:%v)", path, info.IsDir())
		if spinner != nil {
			spinner.UpdateText(fmt.Sprintf("indexing \"%s\"", path))
		}
		logPath, subcnt, skipcnt, errcnt, err := walker.Walk(path, newFileChan, filterPatterns, ignPatterns, rootOptIgnoreEmpty, flogger)
		if err != nil {
			logger.Errorf("walking %s: %v", path, err)
		}

		// stats
		substat := fmt.Sprintf("%d file(s) indexed from \"%s\"", subcnt, path)
		if errcnt > 0 || skipcnt > 0 {
			substat += fmt.Sprintf(" (errors:%d, skipped:%d, log:%s)", errcnt, skipcnt, logPath)
		}
		stats = append(stats, substat)

		cnt += subcnt
		if spinner != nil {
			spinner.UpdateText(fmt.Sprintf("done indexing %s (found %d file(s))", path, subcnt))
		}
	}

	close(newFileChan)

	if spinner != nil {
		spinner.UpdateText(fmt.Sprintf("done walking, compile duplicates for %d...", cnt))
	}
	duplicates := <-resultChan
	if spinner != nil {
		spinner.Stop()
	}

	reporter.Print(rootOptOutFormat, duplicates)

	if !rootOptQuiet {
		for _, stat := range stats {
			logger.Info(stat)
		}
		logger.Infof("%d duplicates found, %d file(s) processed, duration: %v", len(duplicates.Duplicates), cnt, time.Since(t0))
		if duplicates.TotalWasted > 0 {
			logger.Infof("%s freed if removing all duplicates", reporter.SizeToHuman(duplicates.TotalWasted))
		}
	}

	return nil
}

// Execute entry point
func Execute() error {
	return rootCmd.Execute()
}
