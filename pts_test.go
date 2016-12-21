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

import "testing"

func TestPTSIsAfterWithoutRollover(t *testing.T) {
	p := PTS(2)
	other := PTS(1)
	if !p.After(other) {
		t.Errorf("PTS=%v, not After other=%v", p, other)
	}
}

func TestPTSIsAfterWithRollover(t *testing.T) {
	p := PTS(1)
	other := PTS(8589934591) // maxPts
	if !p.After(other) {
		t.Errorf("PTS=%v, not After other=%v", p, other)
	}
	if other.After(p) {
		t.Errorf("PTS=%v, After other=%v", other, p)
	}

}

func TestPTSIsBeforePositiveInfinity(t *testing.T) {
	p := PTS(1)
	other := PtsPositiveInfinity
	if p.After(other) {
		t.Errorf("PTS=%v, not before other=%v", p, other)
	}
}

func TestPTSIsAfterNegativeInfinity(t *testing.T) {
	p := PTS(0)
	other := PtsNegativeInfinity
	if !p.After(other) {
		t.Errorf("PTS=%v, not before other=%v", p, other)
	}
}

func TestPTSIsNotAfterWithoutRollover(t *testing.T) {
	p := PTS(2)
	other := PTS(3)
	if p.After(other) {
		t.Errorf("PTS=%v, not After other=%v", p, other)
	}
}

func TestPTSIsNotAfterWithRolloverOverThreshold(t *testing.T) {
	p := PTS(162000001)
	other := PTS(8589934591)
	if p.After(other) {
		t.Errorf("PTS=%v, not After other=%v", p, other)
	}
}

func TestPTSRolledOver(t *testing.T) {
	p := PTS(1)
	other := PTS(8589934591) // maxPts
	if !p.RolledOver(other) {
		t.Errorf("PTS=%v, not After other=%v", p, other)
	}
}

func TestPTSDurationFrom(t *testing.T) {
	if 5 != PTS(10).DurationFrom(PTS(5)) {
		t.Error("Expected duration of 5")
	}

	if 16 != PTS(5).DurationFrom(PTS(MaxPts-10)) {
		t.Error("Expected duration of 16")
	}
}

func TestPTSGreaterOrEqual(t *testing.T) {
	if PTS(8589904323).GreaterOrEqual(PTS(146909)) {
		t.Error("Greater or equal failed rollover")
	}

	if !PTS(146909).GreaterOrEqual(PTS(8589904323)) {
		t.Error("Greater or equal failed rollover")
	}

	if !PTS(8589904323).GreaterOrEqual(PTS(8589904323)) {
		t.Error("Greater or equal failed rollover")
	}
}

func TestAdd(t *testing.T) {
	if PTS(6674924900) != PTS(7594224546).Add(PTS(7670634945)) {
		t.Error("PTS addition 1 test failed")
	}
	if PTS(2000) != PTS(1500).Add(PTS(500)) {
		t.Error("PTS addition 2 test failed")
	}
}

func TestInsertPTS(t *testing.T) {
	var pts uint64 = 0x1DEADBEEF
	b := make([]byte, 5)
	InsertPTS(b, pts)
	if ExtractTime(b) != 0x1DEADBEEF {
		t.Error("Insert PTS test 1 failed")
	}
}

func TestMaxPtsConstants(t *testing.T) {
	if MaxPts != PTS_MAX {
		t.Error("PTS_MAX does not equal MaxPts")
	}
}
