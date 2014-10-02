package main

import (
	"flag"
	"os"
	"tf/command"
)

func main() {
	s := flag.Bool("s", false, "transfer secure")
	t := flag.String("t", "", "tag marker")
	h := flag.Bool("h", false, "show help")
	var operation string
	flag.Parse()

	arguments := parseArguments(os.Args[1:])
	if len(arguments) < 1 || *h {
		operation = "help"
	} else {
		operation = arguments[0]
	}

	switch operation {
	case "upload":
		command.Upload(arguments[1:], *t, *s)
	case "download":
		command.Download(arguments[1:], *s)
	case "find":
		command.Find(*t)
	default:
		command.Help()
	}
}

func parseArguments(args []string) (parsed []string) {
	var r []rune
	flagValue := false
	for _, arg := range args {
		r = []rune(arg)
		if string(r[0]) == "-" {
			flagValue = true
			continue
		} else if flagValue {
			flagValue = false
			continue
		}
		parsed = append(parsed, arg)
	}

	return
}
