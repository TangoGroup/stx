package stx

import "log"

// Logger provides functionality to conditionally output to stdout
type Logger struct {
	debug bool
}

// NewLogger returns a Log to be used for logging
func NewLogger(debug bool) *Logger {
	lgr := Logger{debug: debug}
	return &lgr
}

// Debug writes to stdout if debug==true
func (lgr *Logger) Debug(msg string) {
	if lgr.debug {
		log.Println(msg)
	}
}

// TODO implement
// func (lgr *Logger) Fatal(msg string) {}

// TODO implement
// func (lgr *Logger) Print(msg string) {}
