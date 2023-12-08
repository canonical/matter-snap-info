package logger

import (
	"fmt"
	"log"
)

const (
	// ANSI color codes
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	White  = "\033[97m"
	Reset  = "\033[0m"
)

func Printf(color, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf(color + msg + Reset)
}

func Println(color string, args ...interface{}) {
	msg := color + fmt.Sprint(args...) + Reset
	log.Println(msg)
}

func Fatalf(format string, args ...interface{}) {
	msg := Red + fmt.Sprintf(format, args...) + Reset
	log.Fatalf(msg)
}
