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
	msg := green + fmt.Sprint(args...) + reset
	log.Println(msg)
}

func Infof(format string, args ...interface{}) {
	msg := white + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Infoln(args ...interface{}) {
	msg := white + fmt.Sprint(args...) + reset
	log.Println(msg)
}

func Warnf(format string, args ...interface{}) {
	msg := yellow + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Warnln(args ...interface{}) {
	msg := yellow + fmt.Sprint(args...) + reset
	log.Println(msg)
}

func Errorf(format string, args ...interface{}) {
	msg := red + fmt.Sprintf(format, args...) + reset
	log.Printf(msg)
}

func Errorln(args ...interface{}) {
	msg := red + fmt.Sprint(args...) + reset
	log.Println(msg)
}

func Fatalf(format string, args ...interface{}) {
	msg := red + fmt.Sprintf(format, args...) + reset
	log.Fatalf(msg)
}

func Fatalln(args ...interface{}) {
	msg := red + fmt.Sprint(args...) + reset
	log.Fatalln(msg)
}
