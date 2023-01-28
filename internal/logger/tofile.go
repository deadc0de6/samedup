/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package logger

import (
	"fmt"
	"log"
	"os"
)

// FileLogger a file logger
type FileLogger struct {
	fd *os.File
}

// Error print error to stderr
func (fl *FileLogger) Error(err error) {
	txt := errorPre + err.Error() + eol
	if fl.fd != nil {
		fl.fd.WriteString(txt)
	} else {
		Error(fmt.Errorf(txt))
	}
}

// Errorf print error to stderr
func (fl *FileLogger) Errorf(format string, a ...interface{}) {
	out := errorPre + fmt.Sprintf(format, a...) + eol
	txt := clear + out
	if fl.fd != nil {
		fl.fd.WriteString(txt)
	} else {
		Error(fmt.Errorf(txt))
	}
}

// Warn print warning to stderr
func (fl *FileLogger) Warn(text string) {
	out := warnPre + text + eol
	txt := clear + out
	if fl.fd != nil {
		fl.fd.WriteString(txt)
	} else {
		Warn(txt)
	}
}

// Warnf print warning to stderr
func (fl *FileLogger) Warnf(format string, a ...interface{}) {
	out := warnPre + fmt.Sprintf(format, a...) + eol
	txt := clear + out
	if fl.fd != nil {
		fl.fd.WriteString(txt)
	} else {
		Error(fmt.Errorf(txt))
	}
}

// Close close the file
func (fl *FileLogger) Close() {
	if fl.fd != nil {
		fl.fd.Close()
	}
}

// GetPath return the path where logs are pushed
func (fl *FileLogger) GetPath() string {
	if fl.fd != nil {
		return fl.fd.Name()
	}
	return ""
}

// NewFileLogger a file logger
func NewFileLogger() *FileLogger {
	fd, err := os.CreateTemp("", "samedup.*.log")
	if err != nil {
		log.Fatal(err)
	}

	fl := &FileLogger{
		fd: fd,
	}
	return fl
}
