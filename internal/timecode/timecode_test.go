package timecode

import (
	"math/big"
	"testing"
)

func tc(h, m, s, f int) *SMPTETimecode {
	return &SMPTETimecode{Sign: 1, Hours: h, Minutes: m, Seconds: s, Frames: f}
}

func tcSigned(sign, h, m, s, f int) *SMPTETimecode {
	return &SMPTETimecode{Sign: sign, Hours: h, Minutes: m, Seconds: s, Frames: f}
}

func TestNewFrameRate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string // Expected rational number string
		expectErr bool
	}{
		{name: "Integer Frame Rate", input: "24", expected: "24/1", expectErr: false},
		{name: "Decimal Frame Rate", input: "23.976", expected: "2997/125", expectErr: false},
		{name: "Another Decimal Frame Rate", input: "29.97", expected: "2997/100", expectErr: false},
		{name: "Invalid String", input: "abc", expected: "", expectErr: true},
		{name: "Empty String", input: "", expected: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr, err := NewFrameRate(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected an error for input %s, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for input %s, but got: %v", tt.input, err)
				}
				if fr.String() != tt.expected {
					t.Errorf("For input %s, expected %s, but got %s", tt.input, tt.expected, fr.String())
				}
			}
		})
	}
}

func TestParseSMPTETimecode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *SMPTETimecode
		expectErr bool
	}{
		{name: "Valid Timecode", input: "01:02:03:04", expected: tc(1, 2, 3, 4), expectErr: false},
		{name: "Zero Timecode", input: "00:00:00:00", expected: tc(0, 0, 0, 0), expectErr: false},
		{name: "Max Values", input: "23:59:59:29", expected: tc(23, 59, 59, 29), expectErr: false},
		{name: "Negative Timecode", input: "-00:00:01:00", expected: tcSigned(-1, 0, 0, 1, 0), expectErr: false},
		{name: "Invalid Format - too few parts", input: "01:02:03", expected: nil, expectErr: true},
		{name: "Invalid Format - too many parts", input: "01:02:03:04:05", expected: nil, expectErr: true},
		{name: "Invalid Hours", input: "XX:02:03:04", expected: nil, expectErr: true},
		{name: "Invalid Minutes", input: "01:XX:03:04", expected: nil, expectErr: true},
		{name: "Invalid Seconds", input: "01:02:XX:04", expected: nil, expectErr: true},
		{name: "Invalid Frames", input: "01:02:03:XX", expected: nil, expectErr: true},
		{name: "Minutes Out of Range", input: "01:60:03:04", expected: nil, expectErr: true},
		{name: "Seconds Out of Range", input: "01:02:60:04", expected: nil, expectErr: true},
		{name: "Negative Frames", input: "01:02:03:-1", expected: nil, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := ParseSMPTETimecode(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for input %s, but got: %v", tt.input, err)
				}
				if tc == nil || *tc != *tt.expected {
					t.Errorf("For input %s, expected %+v, but got %+v", tt.input, tt.expected, tc)
				}
			}
		})
	}
}

func TestToMilliseconds(t *testing.T) {
	fr24, _ := NewFrameRate("24")
	fr2997, _ := NewFrameRate("29.97")
	fr25, _ := NewFrameRate("25")
	fr30, _ := NewFrameRate("30")

	tests := []struct {
		name      string
		timecode  *SMPTETimecode
		framerate *FrameRate
		expected  *big.Rat
		expectErr bool
	}{
		{name: "00:00:00:00 @ 24fps", timecode: tc(0, 0, 0, 0), framerate: fr24, expected: big.NewRat(0, 1), expectErr: false},
		{name: "00:00:01:00 @ 24fps", timecode: tc(0, 0, 1, 0), framerate: fr24, expected: big.NewRat(1000, 1), expectErr: false},
		{name: "00:00:00:12 @ 24fps", timecode: tc(0, 0, 0, 12), framerate: fr24, expected: big.NewRat(500, 1), expectErr: false},
		{name: "00:00:01:12 @ 24fps", timecode: tc(0, 0, 1, 12), framerate: fr24, expected: big.NewRat(1500, 1), expectErr: false},
		{name: "00:00:00:29 @ 29.97fps", timecode: tc(0, 0, 0, 29), framerate: fr2997, expected: big.NewRat(2900000, 2997), expectErr: false},
		{name: "00:00:01:00 @ 25fps", timecode: tc(0, 0, 1, 0), framerate: fr25, expected: big.NewRat(1000, 1), expectErr: false},
		{name: "00:00:00:15 @ 30fps", timecode: tc(0, 0, 0, 15), framerate: fr30, expected: big.NewRat(500, 1), expectErr: false},
		{name: "Negative timecode", timecode: tcSigned(-1, 0, 0, 1, 0), framerate: fr24, expected: big.NewRat(-1000, 1), expectErr: false},
		{name: "Nil FrameRate", timecode: tc(0, 0, 0, 0), framerate: nil, expected: nil, expectErr: true},
		{name: "Zero FrameRate", timecode: tc(0, 0, 0, 0), framerate: &FrameRate{big.NewRat(0, 1)}, expected: nil, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms, err := tt.timecode.ToMilliseconds(tt.framerate)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected an error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %s, but got: %v", tt.name, err)
				}
				if ms.Cmp(tt.expected) != 0 {
					t.Errorf("For %s, expected %s, but got %s", tt.name, tt.expected.String(), ms.String())
				}
			}
		})
	}
}

func TestMillisecondsToSMPTETimecode(t *testing.T) {
	fr24, _ := NewFrameRate("24")
	fr2997, _ := NewFrameRate("29.97")

	tests := []struct {
		name      string
		ms        *big.Rat
		framerate *FrameRate
		expected  *SMPTETimecode
		expectErr bool
	}{
		{name: "0ms @ 24fps", ms: big.NewRat(0, 1), framerate: fr24, expected: tc(0, 0, 0, 0), expectErr: false},
		{name: "1000ms @ 24fps", ms: big.NewRat(1000, 1), framerate: fr24, expected: tc(0, 0, 1, 0), expectErr: false},
		{name: "500ms @ 24fps", ms: big.NewRat(500, 1), framerate: fr24, expected: tc(0, 0, 0, 12), expectErr: false},
		{name: "1500ms @ 24fps", ms: big.NewRat(1500, 1), framerate: fr24, expected: tc(0, 0, 1, 12), expectErr: false},
		{name: "2900000/2997ms @ 29.97fps", ms: big.NewRat(2900000, 2997), framerate: fr2997, expected: tc(0, 0, 0, 29), expectErr: false},
		{name: "Rounding Up", ms: big.NewRat(499, 1), framerate: fr24, expected: tc(0, 0, 0, 12), expectErr: false},   // Should round up to 12 frames
		{name: "Rounding Down", ms: big.NewRat(450, 1), framerate: fr24, expected: tc(0, 0, 0, 11), expectErr: false}, // Should round down to 11 frames
		{name: "Negative milliseconds", ms: big.NewRat(-1000, 1), framerate: fr24, expected: tcSigned(-1, 0, 0, 1, 0), expectErr: false},
		{name: "Nil Milliseconds", ms: nil, framerate: fr24, expected: nil, expectErr: true},
		{name: "Nil FrameRate", ms: big.NewRat(0, 1), framerate: nil, expected: nil, expectErr: true},
		{name: "Zero FrameRate", ms: big.NewRat(0, 1), framerate: &FrameRate{big.NewRat(0, 1)}, expected: nil, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := MillisecondsToSMPTETimecode(tt.ms, tt.framerate)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected an error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %s, but got: %v", tt.name, err)
				}
				if tc == nil || *tc != *tt.expected {
					t.Errorf("For %s, expected %+v, but got %+v", tt.name, tt.expected, tc)
				}
			}
		})
	}
}

func TestToClockTime(t *testing.T) {
	fr24, _ := NewFrameRate("24")
	fr2997, _ := NewFrameRate("29.97")
	fr25, _ := NewFrameRate("25")
	fr30, _ := NewFrameRate("30")

	tests := []struct {
		name      string
		timecode  *SMPTETimecode
		framerate *FrameRate
		precision int
		expected  string
		expectErr bool
	}{
		{name: "00:00:00:00 @ 24fps, prec 3", timecode: tc(0, 0, 0, 0), framerate: fr24, precision: 3, expected: "00:00:00.000", expectErr: false},
		{name: "00:00:01:00 @ 24fps, prec 3", timecode: tc(0, 0, 1, 0), framerate: fr24, precision: 3, expected: "00:00:01.000", expectErr: false},
		{name: "00:00:00:12 @ 24fps, prec 3", timecode: tc(0, 0, 0, 12), framerate: fr24, precision: 3, expected: "00:00:00.500", expectErr: false},
		{name: "00:00:01:12 @ 24fps, prec 3", timecode: tc(0, 0, 1, 12), framerate: fr24, precision: 3, expected: "00:00:01.500", expectErr: false},
		{name: "00:00:00:29 @ 29.97fps, prec 3", timecode: tc(0, 0, 0, 29), framerate: fr2997, precision: 3, expected: "00:00:00.968", expectErr: false}, // 29 frames / 29.97 fps = 0.9676343 seconds -> 968ms
		{name: "00:00:00:29 @ 29.97fps, prec 2", timecode: tc(0, 0, 0, 29), framerate: fr2997, precision: 2, expected: "00:00:00.97", expectErr: false},  // 968ms -> 970ms with prec 2
		{name: "00:00:00:12 @ 25fps, prec 3", timecode: tc(0, 0, 0, 12), framerate: fr25, precision: 3, expected: "00:00:00.480", expectErr: false},
		{name: "00:00:00:15 @ 30fps, prec 3", timecode: tc(0, 0, 0, 15), framerate: fr30, precision: 3, expected: "00:00:00.500", expectErr: false},
		{name: "Negative timecode", timecode: tcSigned(-1, 0, 0, 1, 0), framerate: fr24, precision: 3, expected: "-00:00:01.000", expectErr: false},
		{name: "Nil FrameRate", timecode: tc(0, 0, 0, 0), framerate: nil, precision: 3, expected: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clockTime, err := tt.timecode.ToClockTime(tt.framerate, tt.precision)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected an error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %s, but got: %v", tt.name, err)
				}
				if clockTime != tt.expected {
					t.Errorf("For %s, expected %s, but got %s", tt.name, tt.expected, clockTime)
				}
			}
		})
	}
}
