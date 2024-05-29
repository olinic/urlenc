package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
)

type Parser struct {
	doubleEncode bool
	decode       bool
	args         []string
}

func (p Parser) Parse(text string) string {
	if p.decode {
		t, err := url.QueryUnescape(text)
		if err != nil {
			log.Fatal(err)
		}
		return t
	} else {
		encoded := url.QueryEscape(text)
		if p.doubleEncode {
			encoded = url.QueryEscape(encoded)
		}
		return encoded
	}
}

func (p Parser) Run() {
	if p.isInputFromPipe() {
		p.parsePipe(os.Stdin, os.Stdout)
	} else {
		p.parseCmd()
	}
}

func (p Parser) isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func (p Parser) parsePipe(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		_, e := fmt.Fprintln(
			w, p.Parse(scanner.Text()))
		if e != nil {
			return e
		}
	}
	return nil
}

func (p Parser) parseCmd() {
	var args = flag.Args()
	if len(args) == 1 {
		fmt.Println(p.Parse(args[0]))
	} else {
		fmt.Println("Expected one input string.")
	}
}

func main() {
	doubleEncode := flag.Bool("D", false, "Double encode")
	urlDecode := flag.Bool("d", false, "URL decode")
	flag.Parse()
	parser := Parser{*doubleEncode, *urlDecode, flag.Args()}
	parser.Run()
}
