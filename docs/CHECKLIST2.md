Comprehensive iTT to TTML Conversion Checklist with Corner Cases and Common Mistakes
This checklist provides a thorough guide for converting .itt (iTunes Timed Text) subtitle files to standard TTML (Timed Text Markup Language) format. It includes acceptance criteria, corner cases that may cause issues, and solutions to common mistakes encountered during conversion. The goal is to ensure accurate, compatible, and valid TTML output suitable for various platforms, addressing both typical scenarios and edge cases.

Table of Contents

General Requirements
File Header and Namespaces
Time Format Conversion
Styling
Layout and Regions
Subtitle Content
Inline Styling and Line Breaks
Apple-Specific Features
Validation and Output Quality
Error Handling
Corner Cases
Common Mistakes and Solutions


General Requirements

 Converter processes .itt files and produces TTML files compliant with W3C TTML1 or TTML2 specifications.
 Output TTML file uses UTF-8 encoding with a valid XML declaration (<?xml version="1.0" encoding="UTF-8"?>).
 Preserves all subtitle text content without loss or unintended modifications.
 Supports multiple frame rates (e.g., 24, 25, 29.97, 30 fps) specified in ttp:frameRate.
 Handles single file and batch processing efficiently.

Corner Case: Input .itt file lacks an XML declaration or uses an unsupported encoding (e.g., UTF-16).

Solution: Default to UTF-8 and add XML declaration if missing; log a warning for non-UTF-8 encodings.

Common Mistake: Assuming all .itt files have consistent structure.

Solution: Validate input file structure (e.g., presence of <tt>, <head>, <body>) before processing.


File Header and Namespaces

 Removes redundant namespaces (e.g., xmlns:tt) while retaining required ones (xmlns, xmlns:tts, xmlns:ttp).
 Simplifies xml:lang (e.g., en from en-US) unless a specific dialect is required.
 Changes ttp:timeBase from smpte to media for clock-based timing.
 Removes SMPTE-specific attributes (ttp:frameRate, ttp:frameRateMultiplier, ttp:dropMode) unless needed.
 Adds ttp:profile (e.g., IMSC1 profile) if required by the target platform.
 Ensures <tt> root element includes all necessary namespaces.

Corner Case: .itt file includes custom namespaces (e.g., xmlns:itt="http://itunes.apple.com/itt") that cause validation errors.

Solution: Strip all non-standard namespaces and log their removal.

Common Mistake: Forgetting to update xml:lang for target platform compatibility.

Solution: Map xml:lang to a standardized list (e.g., ISO 639-1 codes) and allow user configuration.


Time Format Conversion

 Converts SMPTE timecodes (HH:MM:SS:FF) to clock time (HH:MM:SS.sss) with millisecond precision.
 Calculates milliseconds accurately based on frame rate (e.g., frame 09 at 24 fps = 375ms).
 Rounds millisecond values to three decimal places.
 Ensures timing accuracy for all <p> elements (no desynchronization).
 Supports non-standard frame rates (e.g., 23.976 fps with frameRateMultiplier="999 1000").
 Validates begin and end times (begin < end).

Corner Case: Frame rate is missing or ambiguous, leading to incorrect time conversions.
# Comprehensive iTT to TTML Conversion Checklist with Corner Cases and Common Mistakes

This checklist provides a thorough guide for converting `.itt` (iTunes Timed Text) subtitle files to standard TTML (Timed Text Markup Language) format. It includes acceptance criteria, corner cases that may cause issues, and solutions to common mistakes encountered during conversion. The goal is to ensure accurate, compatible, and valid TTML output suitable for various platforms, addressing both typical scenarios and edge cases.

---

## Table of Contents
1. [General Requirements](#general-requirements)
2. [File Header and Namespaces](#file-header-and-namespaces)
3. [Time Format Conversion](#time-format-conversion)
4. [Styling](#styling)
5. [Layout and Regions](#layout-and-regions)
6. [Subtitle Content](#subtitle-content)
7. [Inline Styling and Line Breaks](#inline-styling-and-line-breaks)
8. [Apple-Specific Features](#apple-specific-features)
9. [Validation and Output Quality](#validation-and-output-quality)
10. [Error Handling](#error-handling)
11. [Corner Cases](#corner-cases)
12. [Common Mistakes and Solutions](#common-mistakes-and-solutions)

---

## General Requirements

- [ ] Converter processes `.itt` files and produces TTML files compliant with W3C TTML1 or TTML2 specifications.
- [ ] Output TTML file uses UTF-8 encoding with a valid XML declaration (`<?xml version="1.0" encoding="UTF-8"?>`).
- [ ] Preserves all subtitle text content without loss or unintended modifications.
- [ ] Supports multiple frame rates (e.g., 24, 25, 29.97, 30 fps) specified in `ttp:frameRate`.
- [ ] Handles single file and batch processing efficiently.

**Corner Case**: Input `.itt` file lacks an XML declaration or uses an unsupported encoding (e.g., UTF-16).
- **Solution**: Default to UTF-8 and add XML declaration if missing; log a warning for non-UTF-8 encodings.

**Common Mistake**: Assuming all `.itt` files have consistent structure.
- **Solution**: Validate input file structure (e.g., presence of `<tt>`, `<head>`, `<body>`) before processing.

---

## File Header and Namespaces

- [ ] Removes redundant namespaces (e.g., `xmlns:tt`) while retaining required ones (`xmlns`, `xmlns:tts`, `xmlns:ttp`).
- [ ] Simplifies `xml:lang` (e.g., `en` from `en-US`) unless a specific dialect is required.
- [ ] Changes `ttp:timeBase` from `smpte` to `media` for clock-based timing.
- [ ] Removes SMPTE-specific attributes (`ttp:frameRate`, `ttp:frameRateMultiplier`, `ttp:dropMode`) unless needed.
- [ ] Adds `ttp:profile` (e.g., IMSC1 profile) if required by the target platform.
- [ ] Ensures `<tt>` root element includes all necessary namespaces.

**Corner Case**: `.itt` file includes custom namespaces (e.g., `xmlns:itt="http://itunes.apple.com/itt"`) that cause validation errors.
- **Solution**: Strip all non-standard namespaces and log their removal.

**Common Mistake**: Forgetting to update `xml:lang` for target platform compatibility.
- **Solution**: Map `xml:lang` to a standardized list (e.g., ISO 639-1 codes) and allow user configuration.

---

## Time Format Conversion

- [ ] Converts SMPTE timecodes (`HH:MM:SS:FF`) to clock time (`HH:MM:SS.sss`) with millisecond precision.
- [ ] Calculates milliseconds accurately based on frame rate (e.g., frame 09 at 24 fps = `375ms`).
- [ ] Rounds millisecond values to three decimal places.
- [ ] Ensures timing accuracy for all `<p>` elements (no desynchronization).
- [ ] Supports non-standard frame rates (e.g., 23.976 fps with `frameRateMultiplier="999 1000"`).
- [ ] Validates `begin` and `end` times (`begin` < `end`).

**Corner Case**: Frame rate is missing or ambiguous, leading to incorrect time conversions.
- **Solution**: Prompt user for frame rate or default to a common value (e.g., 24 fps) with a logged warning.

**Common Mistake**: Incorrect frame-to-millisecond conversion due to ignoring `frameRateMultiplier`.
- **Solution**: Apply multiplier in calculations (e.g., for 23.976 fps, use `frameRate * (999/1000)`).

---

## Styling

- [ ] Converts `tts:fontSize="100%"` to `tts:fontSize="1em"`.
- [ ] Replaces `sansSerif` with `sans-serif` for `tts:fontFamily`.
- [ ] Converts named colors (e.g., `white`) to hexadecimal (e.g., `#FFFFFF`) while supporting both.
- [ ] Preserves all styling attributes (`tts:fontWeight`, `tts:fontStyle`, `tts:color`, etc.).
- [ ] Ensures `style` references in `<p>` or `<span>` match `<style>` IDs in `<styling>`.
- [ ] Supports additional attributes (e.g., `tts:backgroundColor`) if present.

**Corner Case**: `.itt` file uses Apple-specific fonts (e.g., `HelveticaNeue`) unsupported by TTML players.
- **Solution**: Map to generic fonts (e.g., `sans-serif`) and log the substitution.

**Common Mistake**: Missing style definitions cause rendering issues.
- **Solution**: Validate all `style` references against `<styling>` and apply a default style if undefined.

---

## Layout and Regions

- [ ] Retains valid `<region>` definitions (`tts:origin`, `tts:extent`, `tts:textAlign`, `tts:displayAlign`).
- [ ] Validates `tts:origin` and `tts:extent` use percentages or pixels.
- [ ] Ensures `region` attributes reference defined `<region>` IDs.
- [ ] Preserves multiple regions (e.g., `top`, `bottom`).
- [ ] Handles undefined regions by applying a default (e.g., bottom-centered).

**Corner Case**: Region coordinates exceed display bounds (e.g., `tts:origin="150% 150%"`).
- **Solution**: Clamp values to valid ranges (0%–100%) or log an error.

**Common Mistake**: Incorrect `tts:displayAlign` (e.g., `after` vs. `bottom`) causes misaligned subtitles.
- **Solution**: Map `.itt` alignment to TTML standards and test with multiple players.

---

## Subtitle Content

- [ ] Preserves all text within `<p>` elements.
- [ ] Maintains `<body>`, `<div>`, and `<p>` structure.
- [ ] Ensures `<p>` elements have valid `begin` and `end` times.
- [ ] Retains `style` and `region` attributes on `<body>`, `<div>`, or `<p>`.
- [ ] Handles empty `<p>` elements by removing or logging a warning.

**Corner Case**: Overlapping timecodes (e.g., two `<p>` elements with identical `begin` and `end` times).
- **Solution**: Merge overlapping subtitles or prioritize based on order, with a logged warning.

**Common Mistake**: Losing text due to improper XML parsing.
- **Solution**: Use a robust XML parser (e.g., `lxml` in Python) and validate text content post-conversion.

---

## Inline Styling and Line Breaks

- [ ] Preserves `<span>` elements with valid `style` attributes.
- [ ] Supports nested `<span>` elements for complex styling.
- [ ] Retains `<br/>` elements for line breaks.
- [ ] Validates `style` attributes in `<span>` against `<styling>`.
- [ ] Handles invalid `style` references with a default style or error log.

**Corner Case**: Excessive `<span>` nesting causes rendering issues in simple players.
- **Solution**: Flatten nested styles where possible (e.g., combine `italic` and `bold` into one style).

**Common Mistake**: Ignoring `<br/>` elements, resulting in merged lines.
- **Solution**: Explicitly check for `<br/>` during parsing and preserve in output.

---

## Apple-Specific Features

- [ ] Removes `forced_subtitles` attributes or elements.
- [ ] Removes `markerMode` attributes.
- [ ] Ignores `metadata.xml` links.
- [ ] Removes `ttp:dropMode` and SMPTE attributes.
- [ ] Removes custom namespaces or attributes (e.g., `ittp:aspectRatio`).
- [ ] Optionally maps metadata to TTML `<metadata>` if supported.

**Corner Case**: `.itt` file includes forced subtitles critical for content (e.g., translations).
- **Solution**: Generate a separate TTML file for forced subtitles or prompt user for handling.

**Common Mistake**: Retaining Apple-specific attributes causes validation errors.
- **Solution**: Use a whitelist of TTML-compliant attributes and strip others.

---

## Validation and Output Quality

- [ ] Produces TTML output that validates against W3C TTML1/TTML2 schemas.
- [ ] Ensures no XML syntax errors (e.g., unclosed tags).
- [ ] Verifies all styles, regions, and times are defined.
- [ ] Ensures compatibility with common players (e.g., VLC, browsers, streaming platforms).
- [ ] Maintains visual fidelity (positioning, styling) as close as possible to `.itt`.
- [ ] Logs warnings for non-critical issues (e.g., unsupported attributes).

**Corner Case**: TTML output validates but fails in specific players due to strict requirements.
- **Solution**: Test with target players and adjust output (e.g., simplify styles) based on compatibility.

**Common Mistake**: Skipping validation, leading to runtime errors.
- **Solution**: Integrate validators (e.g., `jing`, `ttml2-validator`) into the conversion pipeline.

---

## Error Handling

- [ ] Handles malformed `.itt` files (e.g., missing `<tt>`) by logging and exiting.
- [ ] Detects and reports undefined styles, regions, or timecodes.
- [ ] Manages invalid timecodes (e.g., `end` < `begin`) by skipping or logging.
- [ ] Provides detailed error messages with line numbers.
- [ ] Supports configurable error tolerance (skip vs. fail).

**Corner Case**: Corrupted `.itt` file with invalid characters crashes the converter.
- **Solution**: Sanitize input during parsing and handle encoding errors gracefully.

**Common Mistake**: Silent failures leave users unaware of issues.
- **Solution**: Log all errors and warnings to a file or console with actionable details.

---

## Corner Cases

1. **Missing `<head>` or `<styling>` Section**:
   - **Issue**: `.itt` file lacks style definitions but uses `style` attributes.
   - **Solution**: Create a default style (e.g., `sans-serif`, normal weight) and apply it.

2. **Zero-Duration Subtitles** (`begin` = `end`):
   - **Issue**: Subtitles with no duration may be ignored by players.
   - **Solution**: Extend `end` time by a small amount (e.g., 100ms) or skip the subtitle.

3. **Non-Sequential Timecodes**:
   - **Issue**: Subtitles appear out of order (e.g., `begin="00:00:05.000"` before `00:00:01.000`).
   - **Solution**: Sort `<p>` elements by `begin` time or log a warning.

4. **Large Files with Thousands of Subtitles**:
   - **Issue**: High memory usage or slow processing.
   - **Solution**: Use streaming XML parsing (e.g., `xml.sax`) and optimize time conversions.

5. **Mixed Time Formats**:
   - **Issue**: `.itt` file uses both SMPTE and clock time.
   - **Solution**: Detect and normalize all times to clock format, logging inconsistencies.

6. **Language-Specific Characters**:
   - **Issue**: Non-UTF-8 characters (e.g., Arabic, Chinese) cause encoding errors.
   - **Solution**: Ensure UTF-8 encoding and test with multilingual `.itt` files.

---

## Common Mistakes and Solutions

1. **Mistake**: Incorrect frame rate assumption (e.g., assuming 24 fps when file uses 29.97 fps).
   - **Solution**: Extract `ttp:frameRate` and `ttp:frameRateMultiplier` from `.itt` and validate with user input if missing.

2. **Mistake**: Ignoring undefined regions, causing subtitles to disappear.
   - **Solution**: Define a fallback region (e.g., `tts:origin="0% 85%"`) and log undefined references.

3. **Mistake**: Overwriting existing styles during conversion.
   - **Solution**: Check for duplicate `xml:id` values and merge or rename conflicting styles.

4. **Mistake**: Failing to handle nested `<span>` elements properly.
   - **Solution**: Parse `<span>` recursively and validate all style references at each level.

5. **Mistake**: Not testing with real players, leading to compatibility issues.
   - **Solution**: Test output with VLC, browser-based players, and streaming platforms (e.g., Netflix, YouTube).

6. **Mistake**: Hardcoding conversion rules for specific `.itt` files.
   - **Solution**: Make the converter configurable (e.g., frame rate, time precision, profile) to handle diverse inputs.

---

## Example Checklist Application

### Sample `.itt` Input
```xml
<tt xmlns="http://www.w3.org/ns/ttml" ttp:timeBase="smpte" ttp:frameRate="24" xml:lang="en-US">
  <head>
    <styling>
      <style xml:id="bold" tts:fontWeight="bold" tts:fontSize="100%"/>
    </styling>
    <layout>
      <region xml:id="bottom" tts:origin="0% 85%" tts:extent="100% 15%"/>
    </layout>
  </head>
  <body region="bottom">
    <div>
      <p begin="00:00:01:00" end="00:00:03:00">
        <span style="bold">Test</span><br/>Line 2
      </p>
    </div>
  </body>
</tt>
```

### Expected TTML Output
```xml
<?xml version="1.0" encoding="UTF-8"?>
<tt xmlns="http://www.w3.org/ns/ttml" xmlns:tts="http://www.w3.org/ns/ttml#styling" xmlns:ttp="http://www.w3.org/ns/ttml#parameter" ttp:timeBase="media" xml:lang="en">
  <head>
    <styling>
      <style xml:id="bold" tts:fontWeight="bold" tts:fontSize="1em"/>
    </styling>
    <layout>
      <region xml:id="bottom" tts:origin="0% 85%" tts:extent="100% 15%"/>
    </layout>
  </head>
  <body region="bottom">
    <div>
      <p begin="00:00:01.000" end="00:00:03.000">
        <span style="bold">Test</span><br/>Line 2
      </p>
    </div>
  </body>
</tt>
```

### Verification Steps
- [ ] Validate XML syntax with `xmllint`.
- [ ] Check time conversion (frame 00 at 24 fps = `0ms`).
- [ ] Confirm `style` and `region` references are defined.
- [ ] Test output in VLC and a browser-based player.
- [ ] Verify `<br/>` and `<span>` are preserved.

