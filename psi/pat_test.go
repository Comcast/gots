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

	"github.com/Comcast/gots/v3"
	"github.com/Comcast/gots/v3/packet"
)

var testData = []struct {
	patBytes          string
	wantStreams       int
	wantProgramNumber int
	wantProgramMapPID int
	wantProgramMap    map[int]int
}{
	{"0000b00d0000c100000001e064dee0f320", 1, 1, 100, map[int]int{1: 100}},
	{"0000b0150000c100000001e0640002e0c80003e12ce8f16345", 3, 0, 8192, map[int]int{1: 100, 2: 200, 3: 300}},
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

func TestSPTSpmtPID(t *testing.T) {
	for _, test := range testData {
		pat_bytes, _ := hex.DecodeString(test.patBytes)

		pat, err := NewPAT(pat_bytes)
		if err != nil {
			t.Errorf("Can't parse PAT table %v", err)
		}

		gotPMTpid, err := pat.SPTSpmtPID()
		if test.wantStreams > 1 {
			if err == nil {
				t.Fatal("Expected error for MPTS")
			}
			continue
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}

		if test.wantProgramMapPID != gotPMTpid {
			t.Errorf("Wrong Program Map PID got %v, want %v", gotPMTpid, test.wantProgramMapPID)
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
	wantMap := map[int]int{1: 598}
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

func TestReadPATMultiplePackets(t *testing.T) {
	bs, _ := hex.DecodeString(
		"474000180000b0d9040fcb00000000e0102887f5182888f5222889f52c288af5" +
			"36288bf5402896f5722897f57328a0e06428a1e06e28a2e07828a3e08228a4e0" +
			"8c28a5e09628a6e0a028a8e0b428ace0c828ade0d228aee0dc28afe0e628b0e0" +
			"f028b1e0fa28b2e10428b3e10e28b4e11828b5e12228b6e12728bae12c28bbe1" +
			"3628bce14028bde14a28c0e19028c1e19a28c2e1a428c8e1f428c9e1fe28cae2" +
			"0828cbe21228cce21c28cde22628cee23028cfe23a28d0e24428d3e247000019" +
			"5828d4e26228d5e26c28d6e27628d7e28028d8e28a28d9e29428dae29e28dbe2" +
			"a854869793ffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffff")
	r := bytes.NewReader(bs)
	pat, err := ReadPAT(r)
	if err != nil {
		t.Errorf("Unexpected error reading PAT: %v", err)
	}
	wantMap := map[int]int{
		10375: 5400, 10376: 5410, 10377: 5420, 10378: 5430, 10379: 5440,
		10390: 5490, 10391: 5491, 10400: 100, 10401: 110, 10402: 120,
		10403: 130, 10404: 140, 10405: 150, 10406: 160, 10408: 180,
		10412: 200, 10413: 210, 10414: 220, 10415: 230, 10416: 240,
		10417: 250, 10418: 260, 10419: 270, 10420: 280, 10421: 290,
		10422: 295, 10426: 300, 10427: 310, 10428: 320, 10429: 330,
		10432: 400, 10433: 410, 10434: 420, 10440: 500, 10441: 510,
		10442: 520, 10443: 530, 10444: 540, 10445: 550, 10446: 560,
		10447: 570, 10448: 580, 10451: 600, 10452: 610, 10453: 620,
		10454: 630, 10455: 640, 10456: 650, 10457: 660, 10458: 670,
		10459: 680,
	}
	gotMap := pat.ProgramMap()
	if !reflect.DeepEqual(wantMap, gotMap) {
		t.Errorf("PAT read is invalid, did not have expected program map")
	}
}

func TestBuildPAT(t *testing.T) {
	pkt := parseHexString("474000100000b00d0001c100000001e256f803e71bffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffff")
	acc := packet.NewAccumulator(PatAccumulatorDoneFunc)

	_, err := acc.WritePacket(pkt)
	if err != gots.ErrAccumulatorDone {
		t.Errorf("Single packet PAT expected. This means your doneFunc is probably bad.")
	}

	pat, err := NewPAT(acc.Bytes())
	if err != nil {
		t.Error(err)
	}

	wantMap := map[int]int{1: 598}
	if !reflect.DeepEqual(wantMap, pat.ProgramMap()) {
		t.Errorf("PAT program map does not match want:%v got:%v", wantMap, pat.ProgramMap())
	}
}

func TestBuildMultiPacketPAT(t *testing.T) {
	firstPacket := parseHexString("474000180000b0d9040fcb00000000e0102887f5182888f5222889f52c288af5" +
		"36288bf5402896f5722897f57328a0e06428a1e06e28a2e07828a3e08228a4e0" +
		"8c28a5e09628a6e0a028a8e0b428ace0c828ade0d228aee0dc28afe0e628b0e0" +
		"f028b1e0fa28b2e10428b3e10e28b4e11828b5e12228b6e12728bae12c28bbe1" +
		"3628bce14028bde14a28c0e19028c1e19a28c2e1a428c8e1f428c9e1fe28cae2" +
		"0828cbe21228cce21c28cde22628cee23028cfe23a28d0e24428d3e2")
	secondPacket := parseHexString("470000195828d4e26228d5e26c28d6e27628d7e28028d8e28a28d9e29428dae2" +
		"9e28dbe2a854869793ffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	acc := packet.NewAccumulator(PatAccumulatorDoneFunc)

	_, err := acc.WritePacket(firstPacket)
	if err == gots.ErrAccumulatorDone {
		t.Error("Added first packet of multi-packet PAT and already indicating done. That's not right.")
	}

	_, err = acc.WritePacket(secondPacket)
	if err != gots.ErrAccumulatorDone {
		t.Error("Added second and final packet of multi-packet PAT but indicating not done. That's not right.")
	}

	pat, err := NewPAT(acc.Bytes())
	if err != nil {
		t.Error(err)
	}

	if got, want := pat.NumPrograms(), 52; got != want {
		t.Errorf("Wrong number of programs, want %d got %d", want, got)
	}
}

func TestBuildPAT_ExpectsAnotherPacket(t *testing.T) {
	// First packet of a multi-packet PAT on its own: the section_length spans
	// into the next packet, so the doneFunc must not yet report done.
	pkt := parseHexString("474000180000b0d9040fcb00000000e0102887f5182888f5222889f52c288af5" +
		"36288bf5402896f5722897f57328a0e06428a1e06e28a2e07828a3e08228a4e0" +
		"8c28a5e09628a6e0a028a8e0b428ace0c828ade0d228aee0dc28afe0e628b0e0" +
		"f028b1e0fa28b2e10428b3e10e28b4e11828b5e12228b6e12728bae12c28bbe1" +
		"3628bce14028bde14a28c0e19028c1e19a28c2e1a428c8e1f428c9e1fe28cae2" +
		"0828cbe21228cce21c28cde22628cee23028cfe23a28d0e24428d3e2")

	acc := packet.NewAccumulator(PatAccumulatorDoneFunc)
	_, err := acc.WritePacket(pkt)
	if err == gots.ErrAccumulatorDone {
		t.Errorf("Expected not done because not enough packets are present to create the PAT")
	}
}
