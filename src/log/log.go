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
	reset  = "\033[0m"
)

func Printf(color, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf(color + msg + reset)
}

func Println(color string, args ...interface{}) {
	msg := color + fmt.Sprint(args...) + reset
	log.Println(msg)
}

func Fatalf(format string, args ...interface{}) {
	msg := Red + fmt.Sprintf(format, args...) + reset
	log.Fatalf(msg)
}
