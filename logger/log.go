package logger

import (
	"fmt"
	"log"
)

const (
	// ANSI color codes
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	white  = "\033[97m"
	reset  = "\033[0m"
)

func Successf(format string, args ...interface{}) {
	msg := green + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Successln(args ...interface{}) {
	Successf(fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	msg := white + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Infoln(args ...interface{}) {
	Infof(fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	msg := yellow + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Warnln(args ...interface{}) {
	Warnf(fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	msg := red + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Errorln(args ...interface{}) {
	Errorf(fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	msg := red + fmt.Sprintf(format, args...) + reset
	log.Fatalf(msg)
}

func Fatalln(args ...interface{}) {
	Fatalf(fmt.Sprint(args...))
}
