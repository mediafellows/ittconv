package timecode

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// FrameRate represents a video frame rate as a rational number for precision.
type FrameRate struct {
	*big.Rat
}

// NewFrameRate creates a new FrameRate from a string (e.g., "24", "23.976", "29.97").
func NewFrameRate(s string) (*FrameRate, error) {
	r := new(big.Rat)
	if strings.Contains(s, ".") {
		// Handle decimal frame rates
		parts := strings.Split(s, ".")
		integerPart := parts[0]
		decimalPart := parts[1]

		numStr := integerPart + decimalPart
		denStr := "1" + strings.Repeat("0", len(decimalPart))

		num, ok := new(big.Int).SetString(numStr, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number format in framerate: %s", numStr)
		}
		den, ok := new(big.Int).SetString(denStr, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number format in framerate: %s", denStr)
		}
		r.SetFrac(num, den)

	} else {
		// Handle integer frame rates
		_, ok := r.SetString(s)
		if !ok {
			return nil, fmt.Errorf("invalid framerate format: %s", s)
		}
	}
	return &FrameRate{r}, nil
}

// SMPTETimecode represents a timecode in HH:MM:SS:FF format.
type SMPTETimecode struct {
	Sign    int
	Hours   int
	Minutes int
	Seconds int
	Frames  int
}

// ParseSMPTETimecode parses a string like "HH:MM:SS:FF" into a SMPTETimecode.
func ParseSMPTETimecode(s string) (*SMPTETimecode, error) {
	// Drop-frame timecodes often use ';' between seconds and frames.
	// Normalize them to ':' so they can be parsed alongside standard SMPTE.
	normalized := strings.ReplaceAll(s, ";", ":")
	sign := 1
	if strings.HasPrefix(normalized, "-") {
		sign = -1
		normalized = strings.TrimPrefix(normalized, "-")
	} else {
		normalized = strings.TrimPrefix(normalized, "+")
	}

	parts := strings.Split(normalized, ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid SMPTE timecode format: %s, expected HH:MM:SS:FF", s)
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid hours in timecode %s: %w", s, err)
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minutes in timecode %s: %w", s, err)
	}
	s2, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid seconds in timecode %s: %w", s, err)
	}
	f, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid frames in timecode %s: %w", s, err)
	}

	if m < 0 || m >= 60 || s2 < 0 || s2 >= 60 || f < 0 {
		return nil, fmt.Errorf("timecode out of range: %s", s)
	}

	return &SMPTETimecode{
		Sign:    sign,
		Hours:   h,
		Minutes: m,
		Seconds: s2,
		Frames:  f,
	}, nil
}

// ToMilliseconds converts SMPTETimecode to milliseconds using the given FrameRate.
func (t *SMPTETimecode) ToMilliseconds(fr *FrameRate) (*big.Rat, error) {
	if fr == nil || fr.Num().Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("invalid framerate: cannot be nil or zero")
	}

	// Total seconds from HH:MM:SS
	totalSeconds := big.NewRat(int64(t.Hours*3600+t.Minutes*60+t.Seconds), 1)
	framesAsRat := big.NewRat(int64(t.Frames), 1)

	// Convert frames to seconds: frames / framerate
	secondsFromFrames := new(big.Rat).Quo(framesAsRat, fr.Rat)

	// Total seconds
	totalPreciseSeconds := new(big.Rat).Add(totalSeconds, secondsFromFrames)

	// Convert total seconds to milliseconds: totalPreciseSeconds * 1000
	milliseconds := new(big.Rat).Mul(totalPreciseSeconds, big.NewRat(1000, 1))

	sign := t.Sign
	if sign == 0 {
		sign = 1
	}

	return new(big.Rat).Mul(milliseconds, big.NewRat(int64(sign), 1)), nil
}

// MillisecondsToSMPTETimecode converts milliseconds to SMPTETimecode using the given FrameRate.
func MillisecondsToSMPTETimecode(ms *big.Rat, fr *FrameRate) (*SMPTETimecode, error) {
	if ms == nil || fr == nil || fr.Num().Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("invalid input: milliseconds or framerate cannot be nil or zero")
	}

	sign := 1
	if ms.Sign() < 0 {
		sign = -1
		ms = new(big.Rat).Abs(ms)
	}

	// Convert milliseconds to total seconds as a rational number
	totalSecondsRat := new(big.Rat).Quo(ms, big.NewRat(1000, 1))

	// Get integer part of total seconds
	totalSecondsInt := new(big.Int).Quo(totalSecondsRat.Num(), totalSecondsRat.Denom())

	// Calculate hours, minutes, seconds from integer total seconds
	hours := new(big.Int).Quo(totalSecondsInt, big.NewInt(3600))
	remainingSeconds := new(big.Int).Mod(totalSecondsInt, big.NewInt(3600))
	minutes := new(big.Int).Quo(remainingSeconds, big.NewInt(60))
	seconds := new(big.Int).Mod(remainingSeconds, big.NewInt(60))

	// Calculate fractional part of seconds as frames
	fractionalSecondsRat := new(big.Rat).Sub(totalSecondsRat, new(big.Rat).SetInt(totalSecondsInt))
	framesRat := new(big.Rat).Mul(fractionalSecondsRat, fr.Rat)

	// Round frames to nearest integer (add 0.5 and truncate)
	roundedFramesRat := new(big.Rat).Add(framesRat, big.NewRat(1, 2))
	frames := new(big.Int).Quo(roundedFramesRat.Num(), roundedFramesRat.Denom())

	return &SMPTETimecode{
		Sign:    sign,
		Hours:   int(hours.Int64()),
		Minutes: int(minutes.Int64()),
		Seconds: int(seconds.Int64()),
		Frames:  int(frames.Int64()),
	}, nil
}

// ToClockTime converts SMPTETimecode to HH:MM:SS.sss clock time string.
func (t *SMPTETimecode) ToClockTime(fr *FrameRate, precision int) (string, error) {
	ms, err := t.ToMilliseconds(fr)
	if err != nil {
		return "", err
	}

	sign := ""
	if ms.Sign() < 0 {
		sign = "-"
		ms = new(big.Rat).Abs(ms)
	}

	// Convert milliseconds to total seconds
	totalSeconds := new(big.Rat).Quo(ms, big.NewRat(1000, 1))

	// Get integer part of seconds
	secInt := new(big.Int).Quo(totalSeconds.Num(), totalSeconds.Denom())

	// Get fractional part of seconds
	fractionalSec := new(big.Rat).Sub(totalSeconds, new(big.Rat).SetInt(secInt))

	// Convert fractional seconds to a string with desired precision
	// We'll multiply by 10^precision, convert to int, then format.
	multiplier := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))
	fractionalMillisInt := new(big.Rat).Mul(fractionalSec, multiplier)

	// Round to nearest integer for precision (add 0.5 and truncate).
	roundedFractionalMillisRat := new(big.Rat).Add(fractionalMillisInt, big.NewRat(1, 2))
	roundedFractionalMillis := new(big.Int).Quo(roundedFractionalMillisRat.Num(), roundedFractionalMillisRat.Denom())

	// Format fractional part with leading zeros
	fracFormat := fmt.Sprintf("%%0%dd", precision)
	fracStr := fmt.Sprintf(fracFormat, roundedFractionalMillis.Int64())

	return fmt.Sprintf("%s%02d:%02d:%02d.%s",
		sign,
		int(secInt.Int64()/3600),
		int((secInt.Int64()%3600)/60),
		int(secInt.Int64()%60),
		fracStr,
	), nil
}
