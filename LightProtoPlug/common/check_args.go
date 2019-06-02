package common

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var Opts map[string]bool

func init() {
	Opts = make(map[string]bool)
	Opts["cli"] = true
	Opts["server"] = true
	Opts["all"] = true
}

func CheckArgs() (string, string, string) {
	// 1. pb file
	// 2. opt
	file := flag.String("file", "", "pb file path")
	codeOpt := flag.String("opt", "all", "[cli,server,all] code for pb.go")

	flag.Parse()

	if _, ok := Opts[*codeOpt]; !ok {
		fmt.Printf("unknown opts:%s, support opts:%v", *codeOpt, Opts)
		panic("unknown opts.")
	}

	if *file == "" {
		fmt.Println("please print .pb file path.")
		os.Exit(0)
		//LOG.Fatal("please print .pb file path.")
	}
	if !Exists(*file) {
		fmt.Printf("pb file:%s cant find.", *file)
		panic("file not found error.")
	}

	contents := strings.Split(*file, ".")
	if contents[1] != "proto" {
		fmt.Printf("file:%s is not proto file", *file)
		panic("file type error.")
	}

	return *file, contents[0], *codeOpt
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
