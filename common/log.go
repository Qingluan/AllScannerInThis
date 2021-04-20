package common

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Green       = color.New(color.FgGreen, color.Bold).SprintFunc()
	Blue        = color.New(color.FgBlue, color.Bold).SprintfFunc()
	Yellow      = color.New(color.FgYellow).SprintfFunc()
	LabelYellow = color.New(color.BgYellow, color.FgWhite, color.Bold).SprintfFunc()

	White = color.New(color.FgWhite, color.Bold).SprintFunc()
)

func Info(args ...interface{}) {
	fmt.Printf("%s : %s\r\n", Blue("[*]"), White(args...))
}

func Infor(args ...interface{}) {
	fmt.Printf("                                                                              \r")
	fmt.Printf("%s : %s\r", LabelYellow("[.V.]"), White(args...))
}

// func Infor(args ...interface{}) {
// 	fmt.Printf("                                                                              \r")
// 	fmt.Printf("%s : %s\r", Blue("[*]"), Green(args...))
// }
