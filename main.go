package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lucasepe/uri/templates"
)

const (
	banner = `         _            _   
o  / / _|     URI      |_ 
o / /   |_ Templates  _|  `
)

var (
	flagInput   string
	flagNewLine bool
	flagVersion bool

	gitCommit string
)

func main() {
	configureFlags()

	if flagVersion {
		fmt.Printf("%s version: %s\n", appName(), gitCommit)
		os.Exit(0)
	}

	uri := flag.CommandLine.Arg(0)
	if len(uri) == 0 {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	var values map[string]interface{}
	var err error

	if len(flagInput) > 0 {
		fp, err := os.Open(flagInput)
		exitOnErr(err)
		defer fp.Close()
		values, err = decodeValues(fp)
	} else {
		values, err = decodeValues(os.Stdin)
	}
	exitOnErr(err)

	tpl, err := templates.Parse(uri)
	exitOnErr(err)

	res, err := tpl.Expand(values)
	exitOnErr(err)

	if flagNewLine {
		fmt.Println(res)
	} else {
		fmt.Print(res)
	}
}

func decodeValues(in io.Reader) (map[string]interface{}, error) {
	const limit = 512 * 1024
	data, err := ioutil.ReadAll(io.LimitReader(in, limit))
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		fmt.Print(banner, "\n\n")
		//fmt.Printf("URI builder through variable expansion as specified in RFC 6570.\n\n")
		fmt.Println("URI Template is a way to specify a URL that includes parameters that must be")
		fmt.Printf("substituted before the URL is resolved; e.g.: http://example.org/{userid}.\n\n")
		fmt.Printf("This tool supports up to Level 4 template expressions. See RFC6570 for more information.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <your/uri/template>\n\n", name)

		fmt.Print("EXAMPLE(s):\n\n")
		fmt.Printf("  %s -i data.json http://example.org/~{username}/{term:1}/{term}{?q*,lang}\n", name)
		fmt.Printf("  cat data.json | %s http://example.org/~{username}/{term:1}/{term}{?q*,lang}\n\n", name)
		fmt.Print("  where data.json content is: ")
		fmt.Println(`{"username":"scarlett","term":"black widow","q":{"a":"mars","b":"jupiter"},"lang":"en"}`)
		fmt.Println()
		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()

		fmt.Println("Crafted with passion by Luca Sepe - https://github.com/lucasepe/uri")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.BoolVar(&flagVersion, "v", false, "print the current version and exit")
	flag.CommandLine.StringVar(&flagInput, "i", "", "json input file containing the values ​​of the variables")
	flag.CommandLine.BoolVar(&flagNewLine, "n", false, "append new line character '\\n' to the output")

	flag.CommandLine.Parse(os.Args[1:])
}

func appName() string {
	return filepath.Base(os.Args[0])
}

// exitOnErr check for an error and eventually exit
func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
