/*
MIT License

Copyright 2016 Comcast Cable Communications Management, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gots

import "math"

// PTS constants
const (
	PTS_DTS_INDICATOR_BOTH     = 3 // 11
	PTS_DTS_INDICATOR_ONLY_PTS = 2 // 10
	PTS_DTS_INDICATOR_NONE     = 0 // 00

	// MaxPtsValue is the highest value the PTS can hold before it rolls over, since its a 33 bit timestamp.
	MaxPtsValue = (1 << 33) - 1 // 2^33 - 1 = 8589934591 = 0x1FFFFFFFF

	// MaxPtsTicks is the length of the complete PTS timeline.
	MaxPtsTicks = 1 << 33 // 2^33 = 8589934592 = 0x200000000

	// Used as a sentinel values for algorithms working against PTS
	PtsNegativeInfinity = PTS(math.MaxUint64 - 1) //18446744073709551614
	PtsPositiveInfinity = PTS(math.MaxUint64)     //18446744073709551615
	PtsClockRate        = 90000

	// UpperPtsRolloverThreshold is the threshold for a rollover on the upper end, MaxPtsValue - 30 min
	UpperPtsRolloverThreshold = 8427934591
	// LowerPtsRolloverThreshold is the threshold for a rollover on the lower end, 30 min
	LowerPtsRolloverThreshold = 162000000
)

// PTS represents PTS time
type PTS uint64

// After checks if this PTS is after the other PTS
func (p PTS) After(other PTS) bool {
	switch {
	case other == PtsPositiveInfinity:
		return false
	case other == PtsNegativeInfinity:
		return true
	case p.RolledOver(other):
		return true
	case other.RolledOver(p):
		return false
	default:
		return p > other
	}
}

// GreaterOrEqual returns true if the method reciever is >= the provided PTS
func (p PTS) GreaterOrEqual(other PTS) bool {
	if p == other {
		return true
	}

	return p.After(other)
}

// RolledOver checks if this PTS just rollover compared to the other PTS
func (p PTS) RolledOver(other PTS) bool {
	if other == PtsNegativeInfinity || other == PtsPositiveInfinity {
		return false
	}

	if p < LowerPtsRolloverThreshold && other > UpperPtsRolloverThreshold {
		return true
	}
	return false
}

// DurationFrom returns the difference between the two pts times.
// One possible distance has rollover and wraps around the number line,
// the other one does not wrap around the number line and doesn't have rollover.
// DurationFrom will assume the smallest difference is the correct one.
// This number is always positive.
// This function will only be accurate if the difference between the two
// PTS values is less than MaxPtsTicks/2 (2^32 = 4294967296 = 0x100000000)
func (p PTS) DurationFrom(from PTS) PTS {
	var difference PTS
	if p == PtsPositiveInfinity || p == PtsNegativeInfinity || from == PtsPositiveInfinity || from == PtsNegativeInfinity {
		return PtsPositiveInfinity
	}
	if p > from {
		difference = p - from
	} else {
		difference = from - p
	}
	if difference >= MaxPtsTicks/2 {
		difference = MaxPtsTicks - difference
	}
	return difference
}

// Add adds the two PTS times together and returns a new PTS
func (p PTS) Add(x PTS) PTS {
	return (p + x) & MaxPtsValue
}

// ExtractTime extracts a PTS time
func ExtractTime(bytes []byte) uint64 {
	var a, b, c, d, e uint64
	a = uint64((bytes[0] >> 1) & 0x07)
	b = uint64(bytes[1])
	c = uint64((bytes[2] >> 1) & 0x7f)
	d = uint64(bytes[3])
	e = uint64((bytes[4] >> 1) & 0x7f)
	return (a << 30) | (b << 22) | (c << 15) | (d << 7) | e
}

// InsertPTS insterts a given pts time into a byte slice and sets the
// marker bits.  len(b) >= 5
func InsertPTS(b []byte, pts uint64) {
	b[0] = byte(pts >> 29 & 0x0f) // PTS[32..30]
	b[1] = byte(pts >> 22 & 0xff) // PTS[29..22]
	b[2] = byte(pts >> 14 & 0xff) // PTS[21..15]
	b[3] = byte(pts >> 7 & 0xff)  // PTS[14..8]
	b[4] = byte(pts&0xff) << 1    // PTS[7..0]

	// Set the marker bits as appropriate
	b[0] |= 0x21
	b[2] |= 0x01
	b[4] |= 0x01
}
