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
package psi

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

var testData = []struct {
	patBytes          string
	wantStreams       int
	wantProgramNumber uint16
	wantProgramMapPID uint16
	wantProgramMap    map[uint16]uint16
}{
	{"0000b00d0000c100000001e064dee0f320", 1, 1, 100, map[uint16]uint16{1: 100}},
	{"0000b0150000c100000001e0640002e0c80003e12ce8f16345", 3, 0, 8192, map[uint16]uint16{1: 100, 2: 200, 3: 300}},
}

func TestNumPrograms(t *testing.T) {
	for _, test := range testData {
		pat_bytes, _ := hex.DecodeString(test.patBytes)

		pat, err := NewPAT(pat_bytes)
		if err != nil {
			t.Errorf("Can't parse PAT table %v", err)
		}

		gotStreams := pat.NumPrograms()
		if test.wantStreams != gotStreams {
			t.Errorf("Invalid number of streams! got %v, want %v", gotStreams, test.wantStreams)
		}
	}
}

// ProgramMapPid is DEPRECATED - new code should use ProgramMap
func TestProgramMapPid(t *testing.T) {
	for _, test := range testData {
		pat_bytes, _ := hex.DecodeString(test.patBytes)

		pat, err := NewPAT(pat_bytes)
		if err != nil {
			t.Errorf("Can't parse PAT table %v", err)
		}

		gotPMPID := pat.ProgramMapPid()
		if test.wantProgramMapPID != gotPMPID {
			t.Errorf("Wrong Program Map PID! got %v, want %v", gotPMPID, test.wantProgramMapPID)
		}
	}
}

// ProgramNumber is DEPRECATED - new code should use ProgramMap
func TestProgramNumber(t *testing.T) {
	for _, test := range testData {
		pat_bytes, _ := hex.DecodeString(test.patBytes)

		pat, err := NewPAT(pat_bytes)
		if err != nil {
			t.Errorf("Can't parse PAT table %v", err)
		}

		gotPN := pat.ProgramNumber()
		if test.wantProgramNumber != gotPN {
			t.Errorf("Wrong Program Number! got %v, want %v", gotPN, test.wantProgramNumber)
		}
	}
}

func TestProgramMap(t *testing.T) {
	for _, test := range testData {
		pat_bytes, _ := hex.DecodeString(test.patBytes)

		pat, err := NewPAT(pat_bytes)
		if err != nil {
			t.Errorf("Can't parse PAT table %v", err)
		}

		gotMap := pat.ProgramMap()
		if !reflect.DeepEqual(test.wantProgramMap, gotMap) {
			t.Errorf("Wrong Program Map! got %v, want %v", gotMap, test.wantProgramMap)
		}
	}
}

func TestReadPATForSmoke(t *testing.T) {
	// requires full packets so cannot use test data above
	bs, _ := hex.DecodeString("474000100000b00d0001c100000001e256f803e71bfffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ff474256100002b0300001c10000e131f0060504435545491be121f0042a027e1" +
		"f86e225f00f52012a9700e9080c001f41850fa041ee3f6580ffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffff")
	r := bytes.NewReader(bs)
	pat, err := ReadPAT(r)
	if err != nil {
		t.Errorf("Unexpected error reading PAT: %v", err)
	}
	// sanity check (tests integration a bit)
	gotMap := pat.ProgramMap()
	wantMap := map[uint16]uint16{1: 598}
	if !reflect.DeepEqual(wantMap, gotMap) {
		t.Errorf("PAT read is invalid, did not have expected program map")
	}
}

func TestReadPATIncomplete(t *testing.T) {
	bs, _ := hex.DecodeString("47400") // incomplete PAT packet
	r := bytes.NewReader(bs)

	_, err := ReadPAT(r)
	if err == nil {
		t.Errorf("Expected to get error reading PAT, but did not")
	}
}
