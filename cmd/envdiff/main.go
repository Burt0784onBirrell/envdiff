package main

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func main() {
	var (
		maskSecrets = flag.Bool("mask", true, "mask secret values in output")
		outputFile  = flag.String("output", "", "write report to file instead of stdout")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <file1> <file2>\n\n")
		fmt.Fprintf(os.Stderr, "Compare two .env files and report differences.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	leftPath, rightPath := args[0], args[1]

	left, err := parser.ParseFile(leftPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", leftPath, err)
		os.Exit(1)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", rightPath, err)
		os.Exit(1)
	}

	result := diff.Compare(left, right)

	out := os.Stdout
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	opts := diff.ReportOptions{
		LeftLabel:   leftPath,
		RightLabel:  rightPath,
		MaskSecrets: *maskSecrets,
	}

	if err := diff.WriteReport(out, result, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if result.HasDiff() {
		os.Exit(2)
	}
}
