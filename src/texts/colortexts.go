package texts

import (
	"github.com/fatih/color"
)

var Prompt func(format string, a ...interface{})
var Ok func(format string, a ...interface{})
var Warning func(format string, a ...interface{})
var Error func(format string, a ...interface{})

func init() {
	Prompt = color.New(color.Bold, color.FgBlue).PrintfFunc()
	Ok = color.New(color.Bold, color.FgGreen).PrintfFunc()
	Warning = color.New(color.Bold, color.FgYellow).PrintfFunc()
	Error = color.New(color.Bold, color.FgRed).PrintfFunc()
}