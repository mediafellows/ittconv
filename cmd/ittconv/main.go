package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mediafellows/ittconv"

	"github.com/alecthomas/kong"
)

var CLI struct {
	InputFile  string `kong:"arg,required,help='Input .itt file path.',type='existingfile'"`
	OutputFile string `kong:"short='o',help='Output file path. If not provided, output is written to stdout.'"`
	Format     string `kong:"short='f',help='Output format (vtt or ttml). Defaults to vtt.',default='vtt'"`
}

func main() {
	ctx := kong.Parse(&CLI)

	// Read input file
	inputData, err := ioutil.ReadFile(CLI.InputFile)
	if err != nil {
		ctx.Fatalf("Failed to read input file: %v", err)
	}

	// Convert to the target format
	var output string
	switch CLI.Format {
	case "vtt":
		output, err = ittconv.ToVTT(string(inputData))
	case "ttml":
		output, err = ittconv.ToTTML(string(inputData))
	default:
		ctx.Fatalf("Unsupported format: %s. Please use 'vtt' or 'ttml'.", CLI.Format)
	}

	if err != nil {
		ctx.Fatalf("Failed to convert to %s: %v", CLI.Format, err)
	}

	// Write output
	if CLI.OutputFile != "" {
		err = ioutil.WriteFile(CLI.OutputFile, []byte(output), 0644)
		if err != nil {
			ctx.Fatalf("Failed to write to output file: %v", err)
		}
		fmt.Printf("Successfully converted %s to %s (%s).\n", filepath.Base(CLI.InputFile), filepath.Base(CLI.OutputFile), CLI.Format)
	} else {
		fmt.Print(output)
	}
}
