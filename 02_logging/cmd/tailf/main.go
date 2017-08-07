package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dustin/go-follow"
)

func main() {
	f := flag.String("f", "", "File to follow")
	flag.Parse()

	if *f == "" {
		log.Fatal("Please specify a file to follow")
	}

	file, err := os.Open(*f)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	tail := follow.New(file)
	defer tail.Close()

	scanner := bufio.NewScanner(tail)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading:", err)
	}
}
