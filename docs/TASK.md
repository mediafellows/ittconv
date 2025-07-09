Coding Agent Task: Develop a Go 1.24 Module for iTT to TTML and VTT Converter
Task Overview
You are a coding agent tasked with developing a Go 1.24 module that converts .itt (iTunes Timed Text) subtitle files into standard TTML (Timed Text Markup Language) and WebVTT formats. The module must be robust, well-structured, and adhere to Go conventions, leveraging Go 1.24+ features. It includes a CLI application for user interaction and comprehensive testing to ensure correctness. The implementation must follow provided requirements, reference docs/*.md files for conversion rules and acceptance criteria, and pass all checks in docs/CHECKLIST.md and then in docs/CHECKLIST2.md. 
All docs/*.md files must be preserved in the final repository.
Your goal is to produce clean, maintainable code with thorough documentation and testing, addressing all specified requirements and additional considerations to ensure a production-ready module.

Task Objectives

Implement a Go 1.24 module with functions to convert .itt to TTML and WebVTT.
Use precise rational number calculations for timecode conversions.
Parse .itt XML efficiently using a SAX parser.
Build a user-friendly CLI application.
Include comprehensive unit, property, mutation, and integration tests.
Ensure code follows Go best practices and passes all acceptance criteria.
Retain all docs/*.md files and use them as development guides.


Requirements
General Requirements

Go Version: Use Go 1.24+, utilizing features like math/big.Rat, enhanced slices/maps packages, and modern error handling.
Module Structure: Organize code into modular packages (e.g., parsing, conversion, CLI) under internal/ to enforce encapsulation.
Code Quality: Adhere to Go conventions (gofmt, go vet, Effective Go), with consistent naming and clear comments.
Documentation:
Add GoDoc comments for all exported types, functions, and methods.
Create a README.md with setup, usage, and testing instructions.


Dependencies:
github.com/orisano/gosax for SAX-based XML parsing.
github.com/asticode/go-astisub for TTML to WebVTT conversion.
github.com/alecthomas/kong for CLI argument parsing.
math/big for rational number calculations.
go.uber.org/zap for structured logging.
github.com/leanovate/gopter for property-based testing (optional).
github.com/zimmski/go-mutesting for mutation testing (optional).
Minimize additional dependencies.


Logging: Use structured logging with configurable levels (debug, info, warn, error).
Configuration: Support CLI flags and optional config files (JSON/YAML) for frame rate, time precision, and TTML profile.
Docs Preservation: Keep all docs/*.md files (e.g., docs/GUIDE.md, docs/CHECKLIST.md, docs/CHECKLIST2.md) in the repository.

Functional Requirements
Entry Point Functions

TTML Conversion:
func ConvertToTTML(ittSource, framerate string) (string, error)


Input:
ittSource: .itt XML content as a string.
framerate: Frame rate (e.g., 24, 23.976, 29.97).


Output: TTML XML content as a string or an error.
Behavior: Convert .itt to TTML per docs/GUIDE.md, meeting docs/CHECKLIST.md and docs/CHECKLIST2.md criteria.


WebVTT Conversion:
func ConvertToVTT(ittSource, framerate string) (string, error)


Input: Same as ConvertToTTML.
Output: WebVTT content as a string or an error.
Behavior: Convert .itt to TTML, then use github.com/asticode/go-astisub to produce WebVTT, ensuring valid output.



Conversion Logic

Parsing:

Use github.com/orisano/gosax for streaming XML parsing.
Handle malformed XML with clear error messages.
Support .itt structures as defined in docs/GUIDE.md.


Timecode Conversion:

Use math/big.Rat for precise frame-to-millisecond conversions.
Convert SMPTE timecodes (HH:MM:SS:FF) to clock time (HH:MM:SS.sss).
Support non-standard frame rates (e.g., 23.976 fps with frameRateMultiplier).
Validate timecode consistency (begin < end).


TTML Generation:

Remove Apple-specific features (e.g., forced_subtitles, markerMode, dropMode).
Adjust namespaces, styling, and regions per TTML standards.
Preserve text, <span> styling, and <br/> line breaks.
Support configurable TTML profiles (e.g., IMSC1).


WebVTT Conversion:

Use github.com/asticode/go-astisub to convert TTML to WebVTT.
Ensure proper WebVTT cues, timestamps, and styling.
Address WebVTT limitations (e.g., limited styling support).


Corner Cases:

Handle cases from docs/CHECKLIST2.md (e.g., missing styles, overlapping timecodes, non-UTF-8 characters).
Implement fallbacks (e.g., default region, style) for undefined elements.



CLI Application

Framework: Use github.com/alecthomas/kong for argument parsing.
Features:
Required arguments: input .itt file path, output file path, frame rate.
Optional flags: --profile (e.g., imsc1), --precision (decimal places), --log-level, --vtt (for WebVTT output).
Support batch processing of multiple files.
Include --version flag.
Allow stdout output for piping.


Example Commands:
ittconv --input input.itt --output output.ttml --framerate 24
ittconv --input input.itt --vtt --output output.vtt --framerate 23.976


Error Handling: Show user-friendly error messages for invalid inputs.

Testing Requirements

Unit Tests:

Test components (e.g., timecode conversion, style mapping).
Cover edge cases (e.g., zero-duration subtitles, invalid frame rates).
Use table-driven tests with testing package.


Property Tests:

Use github.com/leanovate/gopter (optional) for properties like timecode monotonicity and XML well-formedness.
Ensure robustness across input ranges.


Mutation Tests:

Use github.com/zimmski/go-mutesting (optional) for critical functions (e.g., ConvertToTTML).
Verify test suite catches mutations.


Integration Tests:

Test end-to-end conversion with sample .itt files (multilingual, complex).
Validate output with W3C TTML schemas (e.g., via jing or ttml2-validator).
Check compatibility with players (e.g., VLC, browsers).


Test Coverage:

Achieve 90%+ coverage (go test -cover).
Include coverage reports in CI.


Test Data:

Create testdata/ with sample .itt, TTML, and WebVTT files.
Include golden files for expected outputs.



Non-Functional Requirements

Performance:

Convert a 2-hour .itt file (1000 subtitles) in < 5 seconds (2.5 GHz CPU, 8 GB RAM).
Process 100 files in < 1 minute.
Use < 100 MB peak memory for a 10 MB .itt file via streaming parsing.


Error Handling:

Provide detailed errors with context (e.g., line number).
Support configurable error tolerance (skip vs. fail).
Log warnings for non-critical issues.


Maintainability:

Use packages like internal/parser, internal/ttml, internal/vtt.
Follow single-responsibility principle.
Comment complex logic.


Portability:

Run on Linux, macOS, Windows.
Use standard library where possible.



Acceptance Criteria

Pass all criteria in docs/CHECKLIST.md:
Correct namespaces, timecodes, styling, regions, content.
Remove Apple-specific features.
Validate against TTML schemas.
Ensure player compatibility.


CLI supports all use cases intuitively.
All tests pass with 90%+ coverage.
Handle corner cases and mistakes from docs/CHECKLIST2.md.


Deliverables

Go Module:

Git repository with go.mod for Go 1.24.
Packages: internal/parser, internal/ttml, internal/vtt, internal/timecode.
GoDoc for exported APIs.


CLI Application:

In cmd/ittconv/ with main.go.
Supports all CLI features and configs.


Tests:

Unit, property, mutation, integration tests.
testdata/ with sample files and golden outputs.
Coverage reports.


Documentation:

README.md with setup, usage, testing.
Preserve docs/*.md files.
Inline comments for complex logic.


CI/CD:

GitHub Actions for linting (gofmt, go vet, golint), testing, coverage.
Multi-platform testing (Linux, macOS, Windows).


Build Artifacts:

CLI binary (ittconv) for major platforms.
Build/install instructions.




Instructions for Coding Agent

Setup:

Initialize a Go module (go mod init github.com/username/ittconv).
Add dependencies to go.mod.
Create directory structure: cmd/ittconv/, internal/, testdata/, docs/.
Copy docs/*.md files into docs/.


Implementation:

Parser: Use github.com/orisano/gosax to parse .itt XML into a structured model.
Timecode Conversion: Implement precise conversions with math/big.Rat.
TTML Conversion: Build TTML output per docs/itt_to_ttml_conversion_guide.md.
WebVTT Conversion: Chain TTML to WebVTT using github.com/asticode/go-astisub.
CLI: Use github.com/alecthomas/kong for argument parsing and config.
Logging: Integrate go.uber.org/zap for structured logs.


Testing:

Write unit tests for each package, covering edge cases.
Add property tests for timecode and XML properties.
Implement mutation tests for critical functions.
Create integration tests with testdata/ files.
Ensure 90%+ coverage.


Validation:

Validate TTML output with a schema validator.
Test WebVTT output with players.
Verify all checklist criteria programmatically.


Documentation:

Write GoDoc and README.md.
Comment complex logic (e.g., timecode math).
Update docs/*.md if new insights arise (but preserve originals).


CI/CD:

Configure GitHub Actions for linting, testing, coverage.
Test on Linux, macOS, Windows.




Additional Considerations

Input Validation:

Check ittSource for valid XML and framerate for numeric format.
Support common frame rate formats (e.g., 24, 23.976).


Output Formatting:

Pretty-print TTML XML (indented) with configurable compact option.
Ensure WebVTT follows standard formatting.


Internationalization:

Support UTF-8 multilingual text (e.g., Arabic, Chinese).
Handle xml:lang appropriately.


Extensibility:

Use interfaces for parser and converter to allow future formats.
Design for easy addition of new features.


Security:

Sanitize XML input to prevent injection.
Limit memory usage for large inputs.


Monitoring:

Log conversion metrics (e.g., time, subtitle count).
Support optional Prometheus integration.




Success Criteria

Module converts .itt to valid TTML and WebVTT.
CLI is intuitive and supports all features.
Tests achieve 90%+ coverage with no failures.
All criteria in docs/itt_to_ttml_converter_checklist.md are met.
docs/*.md files are preserved.
Code is clean, documented, and follows Go 1.24 conventions.


References

Use docs/itt_to_ttml_conversion_guide.md for conversion rules.
Follow docs/itt_to_ttml_converter_checklist.md for acceptance criteria.
Address corner cases and mistakes in docs/itt_to_ttml_conversion_checklist_with_corner_cases.md.
Refer to W3C TTML and WebVTT specifications for standards compliance.


Agent Instructions

Start: Initialize the repository and set up dependencies.
Iterate: Implement parser, converter, CLI, and tests incrementally, validating each component.
Validate: Continuously check against docs/itt_to_ttml_converter_checklist.md.
Test: Run tests frequently and address coverage gaps.
Document: Update documentation as you code.
Finalize: Ensure CI passes, build binaries, and verify checklist compliance.

Proceed with development, ensuring all requirements and considerations are addressed. Report progress and any blockers promptly.
