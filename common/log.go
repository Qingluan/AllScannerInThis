package common

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Green       = color.New(color.FgGreen, color.Bold).SprintFunc()
	Blue        = color.New(color.FgBlue, color.Bold).SprintFunc()
	Yellow      = color.New(color.FgYellow).SprintFunc()
	Red         = color.New(color.FgRed).SprintFunc()
	Cyan        = color.New(color.FgCyan).SprintFunc()
	LabelYellow = color.New(color.BgYellow, color.FgWhite, color.Bold).SprintFunc()

	White = color.New(color.FgWhite, color.Bold).SprintFunc()
)

func Info(args ...interface{}) {
	fmt.Printf("%s : %s        \r\n", Blue("[*]"), White(args...))
}
func InfoOk(args ...interface{}) {
	fmt.Printf("%s : %s        \r\n", Green("[•]"), White(args...))
}

func Infor(args ...interface{}) {
	fmt.Printf("                                                                              \r")
	fmt.Printf("%s : %s\r", Yellow("[*]"), White(args...))
}

func InfoErr(err error, args ...interface{}) {
	fmt.Printf("%s : %s        \r\n", Red("[*]"), White(args...), Yellow(err))
}

// func Infor(args ...interface{}) {
// 	fmt.Printf("                                                                              \r")
// 	fmt.Printf("%s : %s\r", Blue("[*]"), Green(args...))
// }
