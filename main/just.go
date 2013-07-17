package main

import (
	"flag"
	"fmt"
	"just"
	// "log"
	"os"
)

const (
	VER    = "1.0"
	AUTHOR = "肖立群(xlqstar@gmail.com)"
	USAGE  = "Usage:\n\tje [-config arg] [command]"
)

var (
	configFile = flag.String("config", "config", "配置文件路径")
	args       []string
	arg        string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "JustExpress")
		flag.PrintDefaults()
		// fmt.Fprintln(os.Stderr, "  -h   : show help usage")
	}
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	fmt.Println("Version\t:" + VER)
	fmt.Println("Author\t:" + AUTHOR)
	flag.Parse()
	// fmt.Println(flag.Value)
	args = flag.Args()
	if len(args) > 1 {
		fmt.Println(USAGE)
		os.Exit(1)
	} else if len(args) == 1 {
		arg = args[0]
	}
}

func main() {
	switch arg {
	default:
		just.Complie(*configFile)
	case "complie":
		just.Complie(*configFile)
	}
}
