/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package reporter

import (
	"fmt"

	"github.com/deadc0de6/samedup/internal/db"
)

var (
	header      = "#!/usr/bin/env bash\n#\n"
	statsHeader = "# total %d duplicates, total wasted: %s\n#\n"
	dupHeader   = "# %d duplicates for \"%s\" - would free %s\n"
	dupFirst    = "#rm -fv '%s'\n"
	dupDup      = "rm -fv '%s'\n"
	crlf        = "\n"
)

func printScript(dups *db.Duplicates) {
	fmt.Print(header)
	stats := fmt.Sprintf(statsHeader, len(dups.Duplicates), SizeToHuman(dups.TotalWasted))
	fmt.Print(stats)
	for _, dup := range dups.Duplicates {
		wasted := SizeToHuman(dup.Wasted)
		l := fmt.Sprintf(dupHeader, len(dup.Nodes), dup.Key, wasted)
		fmt.Print(l)
		for idx, node := range dup.Nodes {
			tmpl := dupDup
			if idx == 0 {
				tmpl = dupFirst
			}
			l := fmt.Sprintf(tmpl, node.AbsPath)
			fmt.Print(l)
		}
		fmt.Print(crlf)
	}
}
