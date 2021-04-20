package common

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Green = color.New(color.FgGreen, color.Bold).SprintFunc()
	Blue  = color.New(color.FgBlue, color.Bold).SprintfFunc()
)

func Info(args ...interface{}) {
	fmt.Printf("%s : %s\r\n", Blue("[*]"), Green(args...))
}

func Infor(args ...interface{}) {
	fmt.Printf("                                                                              \r")
	fmt.Printf("%s : %s\r", Blue("[*]"), Green(args...))
}
