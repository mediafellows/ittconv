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

- Conversion of .itt to TTML and WebVTT.
- Precise timecode conversions using rational numbers.
- Efficient XML parsing with SAX.
- Configurable frame rates, precision, and TTML profiles.
- Structured logging.
- Comprehensive unit, property, mutation, and integration tests.
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

The `ittconv` CLI tool allows you to convert .itt files with various options.

**Basic Conversion to TTML:**

```bash
./ittconv --input <input.itt> --output <output.ttml> --framerate <frame_rate>
```

Example:

```bash
./ittconv --input input.itt --output output.ttml --framerate 24
```

**Conversion to WebVTT:**

Use the `--vtt` flag to convert to WebVTT format.

```bash
./ittconv --input <input.itt> --vtt --output <output.vtt> --framerate <frame_rate>
```

Example:

```bash
./ittconv --input input.itt --vtt --output output.vtt --framerate 23.976
```

**Optional Flags:**

- `--profile <profile>`: Specify the TTML profile (e.g., `imsc1`).
- `--precision <decimal_places>`: Set the decimal places for time precision.
- `--log-level <level>`: Configure the logging level (debug, info, warn, error).
- `--version`: Display the application version.

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
	"ittconv/internal/ttml"
	"ittconv/internal/vtt"
)

func main() {
	// Example for TTML conversion
	ittSource := `<itt>...</itt>` // Your .itt XML content
	ttmlOutput, err := ttml.ConvertToTTML(ittSource, "24")
	if err != nil {
		log.Fatalf("Error converting to TTML: %v", err)
	}
	fmt.Println("TTML Output:\n", ttmlOutput)

	// Example for WebVTT conversion
	vttOutput, err := vtt.ConvertToVTT(ittSource, "23.976")
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

To check test coverage:

```bash
go test ./... -cover
```

Unit, property, mutation, and integration tests are included to ensure robustness and correctness.

## Project Structure

The project is organized into the following main directories:

- `cmd/ittconv`: Contains the main CLI application.
- `internal/parser`: Handles .itt XML parsing.
- `internal/timecode`: Manages timecode conversions.
- `internal/ttml`: Implements .itt to TTML conversion logic.
- `internal/vtt`: Provides TTML to WebVTT conversion functionality.
- `docs`: Documentation files, including conversion guides and checklists.
- `testdata`: Sample .itt, TTML, and WebVTT files for testing, along with golden files for expected outputs.

## Contributing

Contributions are welcome! Please refer to the `docs/GUIDE.md` and `docs/CHECKLIST.md` for guidelines and acceptance criteria.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: `LICENSE` file is not created yet, but will be added.) 