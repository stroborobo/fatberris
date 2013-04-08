package main

import (
	"github.com/Knorkebrot/fatberris/fatberris-lib"
	"flag"
	"io/ioutil"
	"log"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-o out.m3u] [(chill | up | down | mix), ...]\n", os.Args[0])
}

func main() {
	var outFile string
	flag.StringVar(&outFile, "o", "", "Output file, defaults to STDOUT")

	flag.Usage = usage
	flag.Parse()

	moodArgs := flag.Args()
	if len(moodArgs) == 0 {
		usage()
		os.Exit(2)
	}

	out, err := fatberris.GetM3u(moodArgs)
	if err != nil {
		log.Fatal(err)
	}

	if outFile != "" {
		ioutil.WriteFile(outFile, []byte(out), 0644)
	} else {
		fmt.Print(out)
	}
}
