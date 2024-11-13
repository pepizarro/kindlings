package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Paths with dynamic user directory
	source := flag.String("s", filepath.Join("/run/media", usr.Username, "Kindle", "My Clippings.txt"), "Path for the clippings file (My Clippings.txt)")
	target := flag.String("t", filepath.Join(usr.HomeDir, "books", "kindle"), "Target directory")

	flag.Usage = usage
	flag.Parse()

	if *source == "" || *target == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	parser := NewParser(*source, *target)
	if err := parser.Parse(); err != nil {
		fmt.Println("Error parsing clippings:", err)
		return
	}
	// parser.Print()

	if err := Write(*target, parser.Clippings); err != nil {
		fmt.Println("Error in Write:", err)
		return
	}

}

func usage() {
	fmt.Println("Usage: kindle-clippings [options]")
	flag.PrintDefaults()
}
