package main

import (
	"fmt"
	"log"
	"os"

	"github.com/christiangda/flyover-reverse-engineering/pkg/fly/c3m"
	"github.com/christiangda/flyover-reverse-engineering/pkg/oth"
)

var l = log.New(os.Stderr, "", 0)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [c3m_file]\n", os.Args[0])
		os.Exit(1)
	}

	file := os.Args[1]
	data, err := os.ReadFile(file)
	oth.CheckPanic(err)
	l.Printf("File size: %d bytes\n", len(data))
	c3m, err := c3m.Parse(data)
	oth.CheckPanic(err)
	_ = c3m
}
