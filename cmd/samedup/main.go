/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package main

import (
	"os"

	"github.com/deadc0de6/samedup/internal/commands"
	"github.com/deadc0de6/samedup/internal/logger"
)

func main() {
	err := commands.Execute()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
