# iTT to TTML Converter Acceptance Criteria Checklist

This checklist outlines comprehensive acceptance criteria for a converter that transforms `.itt` (iTunes Timed Text) subtitle files into standard TTML (Timed Text Markup Language) format. It ensures the converter produces valid, compatible, and accurate TTML files suitable for a wide range of platforms. The criteria cover functionality, validation, error handling, performance, and usability.

---

## Table of Contents
1. [General Requirements](#general-requirements)
2. [File Header and Namespaces](#file-header-and-namespaces)
3. [Time Format Conversion](#time-format-conversion)
4. [Styling](#styling)
5. [Layout and Regions](#layout-and-regions)
6. [Subtitle Content](#subtitle-content)
7. [Inline Styling and Line Breaks](#inline-styling-and-line-breaks)
8. [Handling Apple-Specific Features](#handling-apple-specific-features)
9. [Validation and Output Quality](#validation-and-output-quality)
10. [Error Handling](#error-handling)
11. [Performance](#performance)
12. [Usability](#usability)
13. [Testing Requirements](#testing-requirements)

---

## General Requirements

- [x] The converter processes `.itt` files and produces TTML files compliant with W3C TTML1 or TTML2 specifications.
- [x] The output TTML file is encoded in UTF-8 with a proper XML declaration (`<?xml version="1.0" encoding="UTF-8"?>`).
- [x] The converter preserves all subtitle text content without loss or alteration.
- [x] The converter supports `.itt` files with varying frame rates (e.g., 24, 25, 29.97, 30 fps).
- [x] The converter handles both single and batch file processing.

---

## File Header and Namespaces

- [x] Removes redundant namespaces (e.g., `xmlns:tt`) while retaining required ones (`xmlns`, `xmlns:tts`, `xmlns:ttp`).
- [x] Converts `xml:lang` to a general language code (e.g., `en` from `en-US`) unless specified otherwise.
- [x] Changes `ttp:timeBase` from `smpte` to `media` for clock-based timing.
- [x] Removes SMPTE-specific attributes (`ttp:frameRate`, `ttp:frameRateMultiplier`, `ttp:dropMode`) unless required by the target platform.
- [x] Adds `ttp:profile` (e.g., `http://www.w3.org/ns/ttml/profile/imsc1/text`) if specified for IMSC1 compliance.
- [x] Ensures the `<tt>` root element includes all necessary namespaces for TTML compliance.

---

## Time Format Conversion

- [x] Converts SMPTE timecodes (`HH:MM:SS:FF`) to clock time (`HH:MM:SS.sss`) with millisecond precision.
- [x] Correctly calculates milliseconds based on the input frame rate (e.g., frame 09 at 24 fps = `375ms`).
- [x] Rounds millisecond values to three decimal places (e.g., `01:00:06.375`).
- [x] Preserves timing accuracy for all `<p>` elements (no drift or desynchronization).
- [x] Handles non-standard frame rates (e.g., 23.976 fps with `frameRateMultiplier="999 1000"`) accurately.
- [x] Validates that `begin` and `end` times are logically consistent (`begin` < `end`).

---

## Styling

- [x] Converts `tts:fontSize="100%"` to `tts:fontSize="1em"`.
- [x] Replaces `sansSerif` with `sans-serif` for `tts:fontFamily` to ensure compatibility.
- [x] Converts named colors (e.g., `white`) to hexadecimal (e.g., `#FFFFFF`) for consistency, while supporting both formats.
- [x] Preserves all styling attributes (`tts:fontWeight`, `tts:fontStyle`, `tts:color`, etc.) defined in `<styling>`.
- [x] Ensures all `style` references in `<p>` or `<span>` match defined `<style>` IDs in `<styling>`.
- [x] Supports additional styling attributes (e.g., `tts:backgroundColor`, `tts:padding`) if present.

---

## Layout and Regions

- [x] Retains valid `<region>` definitions, including `tts:origin`, `tts:extent`, `tts:textAlign`, and `tts:displayAlign`.
- [x] Validates that `tts:origin` and `tts:extent` use percentages or pixels compatible with TTML.
- [x] Ensures all `region` attributes in `<body>` or `<div>` reference defined `<region>` IDs in `<layout>`.
- [x] Preserves multiple regions (e.g., `top`, `bottom`) if defined in the `.itt` file.
- [x] Handles missing or undefined regions by applying a default region (e.g., centered bottom) or raising a warning.

---

## Subtitle Content

- [x] Preserves all text content within `<p>` elements without modification.
- [x] Maintains the structure of `<body>`, `<div>`, and `<p>` elements as in the input file.
- [x] Ensures all `<p>` elements include valid `begin` and `end` attributes after time conversion.
- [x] Retains `style` and `region` attributes on `<body>`, `<div>`, or `<p>` elements if defined.
- [x] Handles empty `<p>` elements by either removing them or logging a warning, based on configuration.

---

## Inline Styling and Line Breaks

- [x] Preserves `<span>` elements with valid `style` attributes, ensuring referenced styles are defined.
- [x] Supports nested `<span>` elements while maintaining compatibility with simpler TTML players.
- [x] Retains `<br/>` elements for line breaks without modification.
- [x] Validates that all `style` attributes in `<span>` reference existing `<style>` IDs in `<styling>`.
- [x] Handles invalid or undefined `style` references by logging an error or applying a default style.

---

## Handling Apple-Specific Features

- [x] Removes `forced_subtitles` attributes or elements, as they are not supported in standard TTML.
- [x] Removes `markerMode` attributes, as they are Apple-specific.
- [x] Ignores `metadata.xml` links, as they are not applicable to TTML.
- [x] Removes `ttp:dropMode` and related SMPTE attributes.
- [x] Removes custom namespaces (e.g., `xmlns:itt`) or attributes (e.g., `ittp:aspectRatio`) unless explicitly supported.
- [x] Optionally extracts relevant metadata to TTML `<metadata>` if supported by the target platform.

---

## Validation and Output Quality

- [ ] Produces TTML output that validates against W3C TTML1 or TTML2 schemas using tools like `jing` or `ttml2-validator`.
- [x] Ensures no syntax errors in the output XML (e.g., unclosed tags, invalid characters).
- [x] Verifies that all referenced styles, regions, and times are defined and valid.
- [ ] Produces output compatible with common TTML players (e.g., VLC, web browsers, streaming platforms).
- [x] Maintains visual fidelity (e.g., text positioning, styling) as close as possible to the original `.itt` file.
- [x] Includes warnings for non-critical issues (e.g., unsupported attributes) in logs or reports.

---

## Error Handling

- [x] Gracefully handles malformed `.itt` files (e.g., missing `<tt>` root, invalid XML) by logging an error and exiting.
- [x] Detects and reports missing or undefined styles, regions, or timecodes.
- [x] Handles invalid timecodes (e.g., negative times, `end` < `begin`) by logging an error or skipping the affected `<p>`.
- [x] Provides clear error messages with line numbers or element details for debugging.
- [x] Supports configurable error tolerance (e.g., skip invalid elements vs. fail conversion).

---

## Performance

- [ ] Processes a 2-hour `.itt` file with 1000 subtitles in under 5 seconds on standard hardware (e.g., 2.5 GHz CPU, 8 GB RAM).
- [ ] Scales efficiently for batch processing of multiple files (e.g., 100 files in under 1 minute).
- [ ] Minimizes memory usage for large files (e.g., < 100 MB peak memory for a 10 MB `.itt` file).
- [ ] Avoids infinite loops or excessive recursion during parsing or conversion.

---

## Usability

- [x] Provides a command-line interface (CLI) with clear usage instructions (e.g., `converter --input file.itt --output file.ttml`).
- [x] Supports configuration options for frame rate, time precision, and target TTML profile.
- [ ] Optionally provides a graphical user interface (GUI) for non-technical users.
- [x] Generates detailed logs or reports for conversion results, including warnings and errors.
- [ ] Includes documentation with examples, installation steps, and troubleshooting tips.
- [x] Supports input validation to warn users about potential issues before conversion.

---

## Testing Requirements

- [ ] Includes unit tests for timecode conversion accuracy across common frame rates (24, 25, 29.97, 30 fps).
- [ ] Includes integration tests for end-to-end conversion of sample `.itt` files.
- [ ] Tests edge cases, such as empty files, missing styles/regions, or invalid timecodes.
- [ ] Validates output with multiple TTML validators (e.g., `jing`, `ttv`, `ttml2-validator`).
- [ ] Tests compatibility with at least three TTML players (e.g., VLC, browser-based player, streaming platform).
- [ ] Includes performance tests for large files and batch processing.
- [ ] Tests error handling for malformed inputs and unsupported features.

