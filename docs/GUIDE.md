# Comprehensive iTT to TTML Conversion Guide

This guide provides an in-depth explanation for converting `.itt` (iTunes Timed Text) subtitle files into standard TTML (Timed Text Markup Language) format. It includes detailed rules, comprehensive examples, additional insights from common practices, and solutions to frequent issues encountered during conversion.

---

## Table of Contents
1. [Overview](#overview)
2. [File Header and Namespaces](#file-header-and-namespaces)
3. [Time Format Conversion](#time-format-conversion)
4. [Styling](#styling)
5. [Layout and Regions](#layout-and-regions)
6. [Subtitle Content (Body)](#subtitle-content-body)
7. [Inline Styling with `<span>`](#inline-styling-with-span)
8. [Line Breaks](#line-breaks)
9. [Unsupported or Apple-Specific Features](#unsupported-or-apple-specific-features)
10. [Validation and Testing](#validation-and-testing)
11. [Common Issues and Solutions](#common-issues-and-solutions)
12. [Tools and Resources](#tools-and-resources)

---

## Overview

`.itt` is a proprietary subtitle format used by Apple for iTunes and other platforms, based on TTML but with Apple-specific extensions. Converting `.itt` to standard TTML ensures compatibility with a broader range of players and platforms, such as web browsers, streaming services, and broadcast systems. The conversion process involves adjusting namespaces, time formats, styling, and removing Apple-specific features while preserving the subtitle content and presentation.

This guide expands on standard practices and incorporates real-world insights to address complexities and common pitfalls.

---

## File Header and Namespaces

### `.itt` Example
```xml
<tt xmlns="http://www.w3.org/ns/ttml"
    xmlns:tt="http://www.w3.org/ns/ttml"
    xmlns:tts="http://www.w3.org/ns/ttml#styling"
    xmlns:ttp="http://www.w3.org/ns/ttml#parameter"
    xml:lang="en-US"
    ttp:timeBase="smpte"
    ttp:frameRate="24"
    ttp:frameRateMultiplier="999 1000"
    ttp:dropMode="nonDrop">
```

### TTML Equivalent
```xml
<tt xmlns="http://www.w3.org/ns/ttml"
    xmlns:tts="http://www.w3.org/ns/ttml#styling"
    xmlns:ttp="http://www.w3.org/ns/ttml#parameter"
    xml:lang="en"
    ttp:timeBase="media">
```

### Key Changes
- **Remove redundant namespaces**: The `xmlns:tt` namespace is unnecessary in standard TTML and can be omitted.
- **Simplify `xml:lang`**: Use a general language code (e.g., `en`) unless a specific dialect (e.g., `en-US`) is required by the target platform.
- **Change `ttp:timeBase`**: Convert `smpte` to `media` for clock-based timing, which is more universally supported.
- **Remove SMPTE-specific attributes**: Attributes like `ttp:frameRate`, `ttp:frameRateMultiplier`, and `ttp:dropMode` are often irrelevant in standard TTML and should be removed unless explicitly required by the target system.
- **Additional consideration**: Some platforms require `ttp:profile` (e.g., `http://www.w3.org/ns/ttml/profile/imsc1/text`) for IMSC1 compliance. Check the target platform's requirements.

---

## Time Format Conversion

### `.itt` Example
```xml
<p begin="01:00:06:09" end="01:00:11:13">Hello World</p>
```
(At 24 fps, frame 09 = 375ms, frame 13 = 542ms)

### TTML Equivalent
```xml
<p begin="01:00:06.375" end="01:00:11.542">Hello World</p>
```

### Conversion Process
- **SMPTE to Clock Time**:
  - SMPTE timecode in `.itt` uses `HH:MM:SS:FF` (hours, minutes, seconds, frames).
  - Convert frames to milliseconds: `milliseconds = (frame_number / frame_rate) * 1000`.
  - Example: At 24 fps, frame 09 = `(9 / 24) * 1000 = 375ms`.
- **Precision**: Ensure millisecond values are rounded to three decimal places for consistency.
- **Tooling**: Use libraries like `pysrt` or custom scripts to automate frame-to-millisecond conversion.

### Example Calculation
For `01:00:06:09` at 24 fps:
- Base time: `01:00:06` = 3606 seconds.
- Frame 09: `(9 / 24) * 1000 = 375ms`.
- Result: `01:00:06.375`.

---

## Styling

### `.itt` Example
```xml
<head>
  <styling>
    <style xml:id="bold"
      tts:fontFamily="sansSerif"
      tts:fontWeight="bold"
      tts:fontStyle="normal"
      tts:color="white"
      tts:fontSize="100%"/>
    <style xml:id="italic"
      tts:fontStyle="italic"/>
  </styling>
</head>
```

### TTML Equivalent
```xml
<head>
  <styling>
    <style xml:id="bold"
      tts:fontFamily="sans-serif"
      tts:fontWeight="bold"
      tts:fontStyle="normal"
      tts:color="#FFFFFF"
      tts:fontSize="1em"/>
    <style xml:id="italic"
      tts:fontStyle="italic"/>
  </styling>
</head>
```

### Key Changes
- **Font size**: Replace `100%` with `1em` for standard TTML compliance.
- **Font family**: Use `sans-serif` (with hyphen) for broader compatibility, as some players may not recognize `sansSerif`.
- **Color**: Convert named colors (e.g., `white`) to hexadecimal (e.g., `#FFFFFF`) for consistency, though named colors are generally valid.
- **Additional styling**: Include `tts:backgroundColor` or `tts:padding` if used in `.itt` to maintain visual fidelity.

---

## Layout and Regions

### `.itt` Example
```xml
<head>
  <layout>
    <region xml:id="bottom"
      tts:origin="0% 85%"
      tts:extent="100% 15%"
      tts:textAlign="center"
      tts:displayAlign="after"/>
  </layout>
</head>
```

### TTML Equivalent
```xml
<head>
  <layout>
    <region xml:id="bottom"
      tts:origin="0% 85%"
      tts:extent="100% 15%"
      tts:textAlign="center"
      tts:displayAlign="after"/>
  </layout>
</head>
```

### Notes
- **No changes needed**: Most `.itt` region definitions are TTML-compliant.
- **Validation**: Ensure `tts:origin` and `tts:extent` use percentage or pixel values that are valid for the target player.
- **Multiple regions**: If `.itt` uses multiple regions (e.g., `top` and `bottom`), verify that all referenced regions are defined in `<layout>`.

---

## Subtitle Content (Body)

### `.itt` Example
```xml
<body region="bottom" style="normal">
  <div>
    <p begin="01:00:06:09" end="01:00:11:13">
      <span style="bold">iPhone 4</span> is here.
    </p>
    <p begin="01:00:12:00" end="01:00:15:00">
      Revolutionary design.
    </p>
  </div>
</body>
```

### TTML Equivalent
```xml
<body region="bottom" style="normal">
  <div>
    <p begin="01:00:06.375" end="01:00:11.542">
      <span style="bold">iPhone 4</span> is here.
    </p>
    <p begin="01:00:12.000" end="01:00:15.000">
      Revolutionary design.
    </p>
  </div>
</body>
```

### Key Changes
- **Time conversion**: Convert SMPTE timecodes to clock time as described in [Time Format Conversion](#time-format-conversion).
- **Style references**: Ensure all `style` attributes (e.g., `normal`, `bold`) are defined in the `<styling>` section.
- **Region references**: Verify that the `region` attribute (e.g., `bottom`) matches a defined region in `<layout>`.

---

## Inline Styling with `<span>`

### `.itt` Example
```xml
<p>Behind it is <span style="italic"><span style="yellow">intense</span> technology.</span></p>
```

### TTML Equivalent
```xml
<p>Behind it is <span style="italic"><span style="yellow">intense</span> technology.</span></p>
```

### Notes
- **Style definitions**: Ensure `italic` and `yellow` are defined in `<styling>`. Example:
  ```xml
  <style xml:id="yellow" tts:color="#FFFF00"/>
  <style xml:id="italic" tts:fontStyle="italic"/>
  ```
- **Nesting**: TTML supports nested `<span>` elements, but avoid excessive nesting to maintain compatibility with simpler players.
- **Common issue**: Missing style definitions cause validation errors. Always cross-check `<span style>` references.

---

## Line Breaks

### `.itt` Example
```xml
<p>First line<br/>Second line</p>
```

### TTML Equivalent
```xml
<p>First line<br/>Second line</p>
```

### Notes
- **No changes needed**: The `<br/>` element is standard in TTML.
- **Alternative**: Some platforms prefer `<p>` elements for each line to apply separate timings or styles. Example:
  ```xml
  <p begin="00:00:01.000" end="00:00:03.000">First line</p>
  <p begin="00:00:03.000" end="00:00:05.000">Second line</p>
  ```

---

## Unsupported or Apple-Specific Features

| Feature                | Action                                                                 |
|------------------------|----------------------------------------------------------------------|
| `forced_subtitles`     | Remove; not supported in standard TTML.                              |
| `markerMode`           | Remove; Apple-specific and irrelevant for most platforms.             |
| `metadata.xml` link    | Ignore; not applicable to TTML.                                      |
| `dropMode`             | Remove; use `media` time base instead.                               |
| Time Base `smpte`      | Convert to `media` and use clock time.                               |
| `ittp:aspectRatio`     | Remove unless target platform explicitly supports it (rare).          |
| Custom namespaces      | Remove any non-standard namespaces (e.g., `xmlns:itt`).               |

### Additional Notes
- **Forced subtitles**: If forced subtitles are needed, create a separate TTML file with only the forced content.
- **Metadata**: Transfer relevant metadata (e.g., title, description) to TTML `<metadata>` if supported by the target platform.

---

## Validation and Testing

### Steps
1. **Schema Validation**:
   - Use W3C TTML schema (e.g., TTML1 or TTML2) to validate the output.
   - Tools: `jing`, `ttv`, `ttml2-validator`, or XML linters like `xmllint`.
2. **Player Testing**:
   - Test with target players (e.g., VLC, web browsers, streaming platforms).
   - Verify timing accuracy, styling, and region placement.
3. **Manual Review**:
   - Check for missing subtitles, incorrect timings, or misaligned text.
   - Ensure all style and region references are valid.

### Common Validation Errors
- Undefined styles or regions.
- Incorrect time formats (e.g., SMPTE timecodes in TTML).
- Missing or mismatched namespaces.

---

## Common Issues and Solutions

Based on community feedback and real-world experiences, here are the top issues faced during `.itt` to TTML conversion and their solutions:

1. **Issue**: Incorrect timecode conversion leading to desynchronized subtitles.
   - **Solution**: Use a reliable frame-to-millisecond calculator or library (e.g., `pysrt`, `ffmpeg`). Double-check the frame rate (e.g., 24, 25, or 29.97 fps).
   - **Example**: For 29.97 fps, frame 15 = `(15 / 29.97) * 1000 ≈ 500ms`.

2. **Issue**: Styles not rendering due to undefined or incompatible style IDs.
   - **Solution**: Ensure all `style` attributes in `<p>` or `<span>` match IDs in `<styling>`. Use standard style properties (e.g., `tts:color="#FFFFFF"` instead of named colors).
   - **Example**:
     ```xml
     <style xml:id="yellow" tts:color="#FFFF00"/>
     <p><span style="yellow">Text</span></p>
     ```

3. **Issue**: Regions not displaying correctly on some players.
   - **Solution**: Verify `tts:origin` and `tts:extent` values are in percentages or pixels. Test with multiple players to ensure compatibility.
   - **Example**:
     ```xml
     <region xml:id="bottom" tts:origin="10% 80%" tts:extent="80% 10%"/>
     ```

4. **Issue**: Apple-specific attributes causing validation errors.
   - **Solution**: Remove attributes like `forced_subtitles`, `markerMode`, and `dropMode`. Use a TTML validator to catch these issues early.
   - **Example**: Remove `<tt ttp:dropMode="nonDrop">`.

5. **Issue**: Inconsistent font rendering across platforms.
   - **Solution**: Use generic font families (e.g., `sans-serif`, `monospace`) and avoid Apple-specific fonts like `HelveticaNeue`.
   - **Example**:
     ```xml
     <style xml:id="default" tts:fontFamily="sans-serif"/>
     ```

---

## Tools and Resources

### Conversion Tools
- **Python Libraries**: `pysrt`, `subtitle-converter`, or custom scripts for timecode conversion.
- **FFmpeg**: Can convert `.itt` to other formats, though manual TTML adjustment may be needed.
- **Subtitle Edit**: GUI tool for editing and converting subtitle files.
- **ttconv**: Command-line tool for TTML conversions (`pip install ttconv`).

### Validation Tools
- **jing**: Java-based validator for TTML schemas.
- **ttv**: W3C’s Timed Text Validator.
- **ttml2-validator**: Online or CLI tool for TTML validation.
- **xmllint**: General XML linter for syntax checking.

### Resources
- [W3C TTML Specification](https://www.w3.org/TR/ttml2/)
- [IMSC1 Specification](https://www.w3.org/TR/ttml-imsc1.0.1/) for streaming platforms.
- [Apple iTT Documentation](https://developer.apple.com/documentation/itunes_timed_text_itt) (for `.itt` specifics).
- Community forums like Stack Overflow or GitHub issues for `ttconv`.

---

## Comprehensive Example

Below is a complete `.itt` file and its TTML equivalent, incorporating all conversion rules.

### `.itt` Input
```xml
<?xml version="1.0" encoding="UTF-8"?>
<tt xmlns="http://www.w3.org/ns/ttml"
    xmlns:tts="http://www.w3.org/ns/ttml#styling"
    xmlns:ttp="http://www.w3.org/ns/ttml#parameter"
    xml:lang="en-US"
    ttp:timeBase="smpte"
    ttp:frameRate="24"
    ttp:dropMode="nonDrop">
  <head>
    <styling>
      <style xml:id="bold" tts:fontWeight="bold" tts:fontFamily="sansSerif" tts:fontSize="100%"/>
      <style xml:id="italic" tts:fontStyle="italic"/>
    </styling>
    <layout>
      <region xml:id="bottom" tts:origin="0% 85%" tts:extent="100% 15%" tts:textAlign="center"/>
    </layout>
  </head>
  <body region="bottom">
    <div>
      <p begin="00:00:01:00" end="00:00:03:00">
        <span style="bold">Welcome</span> to the future.
      </p>
      <p begin="00:00:04:00" end="00:00:06:00">
        It’s <span style="italic">here</span>.<br/>Now.
      </p>
    </div>
  </body>
</tt>
```

### TTML Output
```xml
<?xml version="1.0" encoding="UTF-8"?>
<tt xmlns="http://www.w3.org/ns/ttml"
    xmlns:tts="http://www.w3.org/ns/ttml#styling"
    xmlns:ttp="http://www.w3.org/ns/ttml#parameter"
    xml:lang="en"
    ttp:timeBase="media">
  <head>
    <styling>
      <style xml:id="bold" tts:fontWeight="bold" tts:fontFamily="sans-serif" tts:fontSize="1em"/>
      <style xml:id="italic" tts:fontStyle="italic"/>
    </styling>
    <layout>
      <region xml:id="bottom" tts:origin="0% 85%" tts:extent="100% 15%" tts:textAlign="center"/>
    </layout>
  </head>
  <body region="bottom">
    <div>
      <p begin="00:00:01.000" end="00:00:03.000">
        <span style="bold">Welcome</span> to the future.
      </p>
      <p begin="00:00:04.000" end="00:00:06.000">
        It’s <span style="italic">here</span>.<br/>Now.
      </p>
    </div>
  </body>
</tt>
```

### Conversion Notes
- Timecodes converted from `00:00:01:00` (frame 00) to `00:00:01.000` and `00:00:03:00` to `00:00:03.000`.
- `sansSerif` changed to `sans-serif`.
- `100%` font size changed to `1em`.
- Removed `ttp:frameRate` and `ttp:dropMode`.
- Changed `ttp:timeBase` to `media`.

