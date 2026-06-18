# iTT to TTML and WebVTT Converter

This Go module provides functionality to convert iTunes Timed Text (.itt) subtitle files into standard TTML (Timed Text Markup Language) and WebVTT formats. It also includes a command-line interface (CLI) for easy interaction.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [CLI Application](#cli-application)
  - [Go Module](#go-module)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- Conversion of .itt to TTML and reference-compatible WebVTT.
- Precise SMPTE timecode conversion using rational numbers.
- Efficient XML parsing with SAX.
- Preservation of line breaks, cue placement, color classes, and italics in WebVTT.
- Structured logging.
- Unit tests and golden subtitle fixtures.
- User-friendly CLI.

## Installation

To get started with the `ittconv` module, ensure you have Go 1.24+ installed.

1. Clone the repository:

```bash
git clone https://github.com/your-username/ittconv.git
cd ittconv
```

2. Build the CLI application:

```bash
go build -o ittconv ./cmd/ittconv
```

This will create an executable named `ittconv` in your current directory.

## Usage

### CLI Application

The `ittconv` CLI reads one `.itt` file and writes WebVTT by default.

```bash
./ittconv input.itt
./ittconv input.itt -o output.vtt
./ittconv input.itt -f ttml -o output.ttml
```

Use `go run ./cmd/ittconv --help` to inspect the current flags.

**Batch Processing:**

Currently, batch processing is not directly supported via a single command. You can use shell scripting to process multiple files.

### Go Module

You can also use `ittconv` as a Go module in your projects:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mediafellows/ittconv"
)

func main() {
	// Example for TTML conversion
	ittSource := `<itt>...</itt>` // Your .itt XML content
	ttmlOutput, err := ittconv.ToTTML(ittSource)
	if err != nil {
		log.Fatalf("Error converting to TTML: %v", err)
	}
	fmt.Println("TTML Output:\n", ttmlOutput)

	// Example for WebVTT conversion
	vttOutput, err := ittconv.ToVTT(ittSource)
	if err != nil {
		log.Fatalf("Error converting to WebVTT: %v", err)
	}
	fmt.Println("WebVTT Output:\n", vttOutput)
}
```

## Testing

To run the tests for the module:

```bash
go test ./...
```

The same command is available through `make test` or `just test`.

To check test coverage:

```bash
go test ./... -cover
```

The `integration-test` automation target is a stub until integration tests are explicitly requested.

## Project Structure

The project is organized into the following main directories:

- `cmd/ittconv`: Contains the main CLI application.
- `internal/parser`: Handles .itt XML parsing.
- `internal/timecode`: Manages timecode conversions.
- `internal/ttml`: Implements .itt to TTML conversion logic.
- `internal/vtt`: Provides direct ITT model to WebVTT conversion functionality.
- `docs`: Documentation files, including conversion guides and checklists.
- `testdata`: Sample .itt, TTML, and WebVTT files for testing, along with golden files for expected outputs.

## Contributing

Contributions are welcome! Please refer to the `docs/GUIDE.md` and `docs/CHECKLIST.md` for guidelines and acceptance criteria.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: `LICENSE` file is not created yet, but will be added.) 
