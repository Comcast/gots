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
	"fmt"
	"testing"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

func parseHexString(h string) *packet.Packet {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic("bad test: " + h)
	}
	pkt := new(packet.Packet)
	if copy(pkt[:], b) != packet.PacketSize {
		panic("bad test (wrong length): " + h)
	}
	return pkt
}

type testPmtElementaryStream struct {
	elementaryPid       uint16
	streamType          uint8
	presentationLagsEbp bool
}

func (es *testPmtElementaryStream) ElementaryPid() uint16 {
	return es.elementaryPid
}

func (es *testPmtElementaryStream) StreamType() uint8 {
	return es.streamType
}

func (es *testPmtElementaryStream) StreamTypeDescription() string {
	return "Test Stream Description"
}

func (es *testPmtElementaryStream) IsAudioContent() bool {
	return es.presentationLagsEbp
}

func (es *testPmtElementaryStream) IsVideoContent() bool {
	return !es.presentationLagsEbp
}

func (es *testPmtElementaryStream) IsSCTE35Content() bool {
	return false
}

func (es *testPmtElementaryStream) IsStreamWherePresentationLagsEbp() bool {
	return es.presentationLagsEbp
}

func (es *testPmtElementaryStream) Descriptors() []PmtDescriptor {
	return make([]PmtDescriptor, 0)
}

func (es *testPmtElementaryStream) MaxBitRate() uint64 {
	return 0
}

func TestParseTable(t *testing.T) {
	byteArray, _ := hex.DecodeString("0002b02d0001cb0000e065f0060504435545491b" +
		"e065f0050e030004b00fe066f0060a04656e670086e06ef0" +
		"007fc9ad32ffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffff")

	pmt := pmt{}
	err := pmt.parseTables(byteArray)
	if err != nil {
		t.Errorf("Can't parse PMT table %v", err)
	}

	want := []uint16{101, 102, 110}
	got := pmt.Pids()
	if len(want) != len(got) {
		t.Errorf("ES count does not match want:%v got:%v", len(want), len(got))
	}

	for i, pid := range got {
		if want[i] != pid {
			t.Errorf("PID does not match Want:%v Got:%v", want[i], pid)
		}
	}
}

func TestParseMultipleTables(t *testing.T) {
	byteArray, _ := hex.DecodeString("00c00015000297000b0000000000000100000000006ed1581a" +
		"02b02d0001cb0000e065f0060504435545491b" +
		"e065f0050e030004b00fe066f0060a04656e670086e06ef0" +
		"007fc9ad32ffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffff")

	pmt := pmt{}
	err := pmt.parseTables(byteArray)
	if err != nil {
		t.Errorf("Can't parse PMT table %v", err)
	}

	want := []uint16{101, 102, 110}
	got := pmt.Pids()
	if len(want) != len(got) {
		t.Errorf("ES count does not match want:%v got:%v", len(want), len(got))
	}

	for i, pid := range got {
		if want[i] != pid {
			t.Errorf("PID does not match Want:%v Got:%v", want[i], pid)
		}
	}
}

func TestBuildPMT(t *testing.T) {
	pkt, _ := hex.DecodeString("474064100002b02d0001cb0000e065f0060504435545491b" +
		"e065f0050e030004b00fe066f0060a04656e670086e06ef0" +
		"007fc9ad32ffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffff")
	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, err := acc.Add(pkt)

	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Errorf("Single packet PMT expected. This means your doneFunc is probably bad.")
	}
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	want := []uint16{101, 102, 110}
	got := pmt.Pids()
	if len(want) != len(got) {
		t.Errorf("PID lengths do not match Expected %d: Got %d", len(want), len(got))
	}
	for i, pid := range got {
		if want[i] != pid {
			t.Errorf("PIDs incorrect in PMT Want %d: Got %d", want[i], pid)
		}
	}
}

func TestPIDExists(t *testing.T) {
	pkt, _ := hex.DecodeString("474064100002b02d0001cb0000e065f0060504435545491b" +
		"e065f0050e030004b00fe066f0060a04656e670086e06ef0" +
		"007fc9ad32ffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffff")
	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, err := acc.Add(pkt)

	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Errorf("Single packet PMT expected. This means your doneFunc is probably bad.")
	}
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	if !pmt.PIDExists(101) {
		t.Error("Expected true for PID 101 to exist in PMT.")
	}

	if pmt.PIDExists(199) {
		t.Error("Expected false for PID 199 to exist in PMT.")
	}
}

func TestBuildPMT_ExpectsAnotherPacket(t *testing.T) {
	pkt, _ := hex.DecodeString(
		"4740271A0002B0BA0001F70000E065F00C0F04534150530504435545491BE065" +
			"F03028046400283F2A0FFF7F00000001000001C2000003E99F0E03C039219700" +
			"E90710830A41850241860701656E677EFFFF0FE066F0160A04656E67000E03C0" +
			"00F09700E90710830A408502400FE067F0160A04737061000E03C000F09700E9" +
			"0710830A4085024087E068F0160A04656E67000E03C001E09700E90710830A40" +
			"85024087E069F0160A04737061000E03C000F09700E90710830A4085")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, _ := acc.Add(pkt)
	if done {
		t.Errorf("Expected Error because not enough packets are present to create PMT")
	}
}
func TestBuildPMT_LargePointerFieldGood(t *testing.T) {
	pkt, _ := hex.DecodeString("474064108700000000000000000000000000000000000000" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"ffffffffffffffffffffffffffffffffffffffff02b02d00" +
		"01cb0000e065f0060504435545491be065f0050e030004b0" +
		"0fe066f0060a04656e670086e06ef0007fc9ad32")
	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, err := acc.Add(pkt)

	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Errorf("Single packet PMT expected. This means your doneFunc is probably bad.")
	}
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	want := []uint16{101, 102, 110}
	got := pmt.Pids()
	if len(want) != len(got) {
		t.Errorf("PID lengths do not match Expected %d: Got %d", len(want), len(got))
	}
	for i, pid := range got {
		if want[i] != pid {
			t.Errorf("PIDs incorrect in PMT Want %d: Got %d", want[i], pid)
		}
	}
}
func TestBuildPMT_LargePointerFieldExpectsAnotherPacket(t *testing.T) {
	pkt, _ := hex.DecodeString("474064108800000000000000000000000000000000000000" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"0102030405060708090a0b0c0d0e0f101112131415161718" +
		"ffffffffffffffffffffffffffffffffffffffffff02b02d" +
		"0001cb0000e065f0060504435545491be065f0050e030004" +
		"b00fe066f0060a04656e670086e06ef0007fc9ad")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, _ := acc.Add(pkt)
	if done {
		t.Errorf("Expected Error because not enough packets are present to create PMT")
	}
}
func TestBuildMultiPacketPMT(t *testing.T) {
	firstPacketBytes, _ := hex.DecodeString("474064100002b0ba0001c10000e065f00b0504435545490e03c03dd01be065f016970028046400283fe907108302808502800e03c0392087e066f0219700050445414333cc03c0c2100a04656e6700e907108302808502800e03c000f087e067f0219700050445414333cc03c0c4100a0473706100e907108302808502800e03c001e00fe068f01697000a04656e6700e907108302808502800e03c000f00fe069f01697000a0473706100e907108302808502800e03c000f086e0dc")

	secondPacketBytes, _ := hex.DecodeString("47006411f0002b59bc22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(firstPacketBytes)
	done, err := acc.Add(secondPacketBytes)
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Fatal("PMT should have been done after 2 packets and it is not")
	}
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}
	wantedPids := []uint16{101, 102, 103, 104, 105, 220}
	if len(wantedPids) != len(pmt.Pids()) {
		t.Errorf("PID length do not match expected %d Got %d", len(wantedPids), len(pmt.Pids()))
		t.FailNow()
	}
	for i, pid := range wantedPids {
		if pid != pmt.Pids()[i] {
			t.Errorf("Pids do not match expected %d, Got %d", pid, pmt.Pids()[i])
		}
	}
}

func TestBuildMultiPacketPMT2(t *testing.T) {
	firstPacket, _ := hex.DecodeString("4741E03001000002B1790001C10000E1E1F00B0504435545490E03C038F31BE1E1F016970028046400293FE907108302808502800E03C024DF0FE1EEF01697000A04656E6700E907108302808502800E03C001700FE1EFF01697000A0473706100E907108302808502800E03C001700FE1F0F01697000A04706F7200E907108302808502800E03C0017087E1E2F0219700050445414333CC03C0C4100A04656E6700E907108302808502800E03C002C287E1E3F02197000504454143")

	secondPacket, _ := hex.DecodeString("4701E031010033CC03C0C2100A0473706100E907108302808502800E03C0017E87E1E4F0219700050445414333CC03C0D2100A04656E6700E907108302808502800E03C0017E81E1E8F0259700050441432D33810706380FFF1F013F0A04656E6700E907108302808502800E03C0054A81E1E9F0259700050441432D338107062005FF1F013F0A0473706100E907108302808502800E03C001EA81E1EAF0259700050441432D338107062045FF00013F0A04656E6703E90710830280")

	thirdPacket, _ := hex.DecodeString("4701E03201008502800E03C001EA86E1F4F00096A58F55FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, err := acc.Add(firstPacket)
	if err != nil {
		t.Error(err)
	}
	if done {
		t.Error("Added first packet of multi-packet and already indicating done. That's not right.")
		t.FailNow()
	}

	done, err = acc.Add(secondPacket)
	if err != nil {
		t.Error(err)
	}
	if done {
		t.Error("Added second packet of multi-packet and already indicating done. That's not right.")
		t.FailNow()
	}

	done, err = acc.Add(thirdPacket)
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("Added third and final packet of multi-packet but indicating not done. That's not right.")
		t.FailNow()
	}

	bytes, parseErr := acc.Parse()
	if parseErr != nil {
		fmt.Printf("%v\n", parseErr)
		return
	}

	pmt, err := NewPMT(bytes)
	if err != nil {
		t.Error(err)
	}
	wantedPids := []uint16{481, 482, 483, 484, 488, 489, 490, 494, 495, 496, 500}
	if len(wantedPids) != len(pmt.Pids()) {
		t.Errorf("PID length do not match expected %d Got %d", len(wantedPids), len(pmt.Pids()))
		t.FailNow()
	}
	for _, wpid := range wantedPids {
		found := false
		for _, pid := range pmt.Pids() {
			if wpid == pid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PIDs not found %d", wpid)
		}
	}
}

func TestElementaryStreams(t *testing.T) {
	pids := []uint16{101, 102, 103}
	want := []PmtElementaryStream{
		&testPmtElementaryStream{101, 27, true},
		&testPmtElementaryStream{102, 15, true},
		&testPmtElementaryStream{103, 15, true},
	}
	pmt := &pmt{pids: pids, elementaryStreams: want}

	got := pmt.ElementaryStreams()

	for i, es := range want {
		if es.ElementaryPid() != got[i].ElementaryPid() {
			t.Errorf("PIDs do not match Want %d: Got %d", es.ElementaryPid(), got[i].ElementaryPid())
		}
	}
	pid := uint16(102)
	if !pmt.IsPidForStreamWherePresentationLagsEbp(pid) {
		t.Errorf("PID %d: presentation should lag EBP", pid)
	}
}
func TestIsPidForStreamWherePresentationLagsEbp(t *testing.T) {
	pids := []uint16{101, 102, 103}
	streams := []PmtElementaryStream{&testPmtElementaryStream{102, 15, true}}
	pmt := &pmt{pids: pids, elementaryStreams: streams}
	if !pmt.IsPidForStreamWherePresentationLagsEbp(102) {
		t.Errorf("Expected Presentation to lag EBP")
	}
}

func TestIsNotPidForStreamWherePresentationLagsEbp(t *testing.T) {
	pids := []uint16{101, 102, 103}
	streams := []PmtElementaryStream{&testPmtElementaryStream{102, 15, false}}
	pmt := &pmt{pids: pids, elementaryStreams: streams}

	if pmt.IsPidForStreamWherePresentationLagsEbp(102) {
		t.Errorf("Expected Presentation to NOT lag EBP")
	}
}

func TestStringFormat(t *testing.T) {
	bytes := []byte{
		0x47, 0x40, 0x64, 0x10, 0x00, 0x02, 0xb0, 0x2d, 0x00, 0x01, 0xcb, 0x00,
		0x00, 0xe0, 0x65, 0xf0, 0x06, 0x05, 0x04, 0x43, 0x55, 0x45, 0x49, 0x1b,
		0xe0, 0x65, 0xf0, 0x05, 0x0e, 0x03, 0x00, 0x04, 0xb0, 0x0f, 0xe0, 0x66,
		0xf0, 0x06, 0x0a, 0x04, 0x65, 0x6e, 0x67, 0x00, 0x86, 0xe0, 0x6e, 0xf0,
		0x00, 0x7f, 0xc9, 0xad, 0x32, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(bytes)
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	want := "PMT[ElementaryStream[pid=101,streamType=27,descriptor0='Maximum Bit-Rate (1200)'],ElementaryStream[pid=102,streamType=15,descriptor0='ISO 639 Language (code=eng, audioType=0)'],ElementaryStream[pid=110,streamType=134]]"
	got := fmt.Sprintf("%v", pmt)
	if want != got {
		t.Errorf("String format for PMT failed. Want: %s: Got %s", want, got)
	}
}

func TestFilterPMTPacketsToPids_SinglePacketPMT(t *testing.T) {
	bytes := packet.Packet{
		0x47, 0x40, 0x64, 0x10, 0x00, 0x02, 0xb0, 0x2d, 0x00, 0x01, 0xcb, 0x00,
		0x00, 0xe0, 0x65, 0xf0, 0x06, 0x05, 0x04, 0x43, 0x55, 0x45, 0x49, 0x1b,
		0xe0, 0x65, 0xf0, 0x05, 0x0e, 0x03, 0x00, 0x04, 0xb0, 0x0f, 0xe0, 0x66,
		0xf0, 0x06, 0x0a, 0x04, 0x65, 0x6e, 0x67, 0x00, 0x86, 0xe0, 0x6e, 0xf0,
		0x00, 0x7f, 0xc9, 0xad, 0x32, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(bytes[:])
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	unfilteredPmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	pids := unfilteredPmt.Pids()
	pids = pids[:len(pids)-1]

	filteredPmtPackets, err := FilterPMTPacketsToPids([]*packet.Packet{&bytes}, pids)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err.Error())
	}

	acc = packet.NewAccumulator(PmtAccumulatorDoneFunc)
	for _, p := range filteredPmtPackets {
		acc.Add(p[:])
	}
	payload, err = acc.Parse()
	filteredPmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}
	for i, pid := range filteredPmt.Pids() {
		if pids[i] != pid {
			t.Errorf("PIDs do not match Expected:%d Got %d", pids[i], pid)
		}
	}
}

func TestFilterPMTPacketsToPids_MultiPacketPMT(t *testing.T) {
	firstPacketBytes := parseHexString("474064100002b0ba0001c10000e065f00b0504435545490e03c03dd01be065f016970028046400283fe907108302808502800e03c0392087e066f0219700050445414333cc03c0c2100a04656e6700e907108302808502800e03c000f087e067f0219700050445414333cc03c0c4100a0473706100e907108302808502800e03c001e00fe068f01697000a04656e6700e907108302808502800e03c000f00fe069f01697000a0473706100e907108302808502800e03c000f086e0dc")

	secondPacketBytes := parseHexString("47006411f0002b59bc22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(firstPacketBytes[:])
	acc.Add(secondPacketBytes[:])
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}

	wantedPids := []uint16{101, 102, 103, 104, 105, 220}

	filteredPids := wantedPids[:len(wantedPids)-1]
	filteredPMTPackets, err := FilterPMTPacketsToPids([]*packet.Packet{firstPacketBytes, secondPacketBytes}, filteredPids)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err.Error())
	}

	acc = packet.NewAccumulator(PmtAccumulatorDoneFunc)
	for _, p := range filteredPMTPackets {
		acc.Add(p[:])
	}

	wantedPids = []uint16{101, 102, 103, 104, 105}
	payload, err = acc.Parse()
	if err != nil {
		t.Error(err)
	}
	filteredPMT, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}
	if len(wantedPids) != len(filteredPMT.Pids()) {
		t.Errorf("PID Length do not match wanted:%d got %d", len(wantedPids), len(filteredPMT.Pids()))
	}
	for i, pid := range filteredPMT.Pids() {
		if wantedPids[i] != pid {
			t.Errorf("PIDs do not match Expected:%d Got %d", wantedPids[i], pid)
		}
	}
}

func TestFilterPMTPacketsToPids_PIDNotFound(t *testing.T) {
	// PMT contains PIDs 101-105.
	pmtPacketBytes, _ := hex.DecodeString("4740641D0002B0940001DF0000E065F0050E03C015581BE065F0150E03C0109D2A027E1F9700E9080C001F418503E84187E066F01A0E03C00122050445414333CC07E0C2B0E8656E670A04656E670087E067F01A0E03C00122050445414333CC07E0C2B0E8656E670A04656E67000FE068F0100E03C001262B030102010A04656E67000FE069F0100E03C001262B030102010A04656E67002E9B5B71FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	// Copy into a packet.Packet for calling FilterPMTPacketsToPids.
	var pkt packet.Packet
	copy(pkt[:], pmtPacketBytes)

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done, err := acc.Add(pmtPacketBytes)

	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Errorf("Single packet PMT expected. This means your doneFunc is probably bad.")
	}

	// Test when some of the PIDs to be filtered exist in the PMT.
	pmtPkts, err := FilterPMTPacketsToPids([]*packet.Packet{&pkt}, []uint16{105, 106})
	if err.Error() != "PID(s) [106] not found in PMT." {
		t.Errorf("Expected missing PID error string, got %s", err.Error())
	}

	if pmtPkts == nil {
		t.Error("Expected non-nil packets since none of the filtered PIDs exist in the PMT.")
	}

	// Test when none of the PIDs to be filtered exist in the PMT.
	pmtPkts, err = FilterPMTPacketsToPids([]*packet.Packet{&pkt}, []uint16{106, 107})
	if err.Error() != "PID(s) [106 107] not found in PMT." {
		t.Errorf("Expected missing PID error string, got %s", err.Error())
	}

	if pmtPkts != nil {
		t.Error("Expected nil packets since none of the filtered PIDs exist in the PMT.")
	}
}

func TestPMTIsIFrameStreamPositive(t *testing.T) {
	firstPacketBytes, _ := hex.DecodeString("4741E03001000002B02D0001C10000E1E1F0050E03C003531BE1E1F016970028044D401F3FE907108301808501800E03C003175D027AA4FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(firstPacketBytes)
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}

	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	var isIFrameProfile bool
	for _, es := range pmt.ElementaryStreams() {
		for _, des := range es.Descriptors() {
			isIFrameProfile = des.IsIFrameProfile()
			if isIFrameProfile {
				break
			}
		}
		if isIFrameProfile {
			break
		}
	}
	if !isIFrameProfile {
		t.Errorf("Positive I-Frame Stream failed. Supposed to be an I-Frame stream.")
	}
}

func TestPMTIsIFrameStreamNegative(t *testing.T) {
	firstPacketBytes, _ := hex.DecodeString("4741E03001000002B0FB0001C10000E1E1F00B0504435545490E03C02FD31BE1E1F016970028046400293FE907108302808502800E03C024DF0FE1E2F01697000A04656E6700E907108302808502800E03C001700FE1E3F01697000A0473706100E907108302808502800E03C001700FE1E4F01697000A04706F7200E907108302808502800E03C0017087E1E5F0219700050445414333CC03C0C4100A04656E6700E907108302808502800E03C002C287E1E6F02197000504454143")

	secondPacketBytes, _ := hex.DecodeString("4701E031010033CC03C0C2100A0473706100E907108302808502800E03C0017E87E1E7F0219700050445414333CC03C0D2100A04656E6700E907108302808502800E03C0017E86E1F4F00013E8BFD4FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(firstPacketBytes)
	acc.Add(secondPacketBytes)
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}

	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}

	var isIFrameProfile bool
	for _, es := range pmt.ElementaryStreams() {
		for _, des := range es.Descriptors() {
			isIFrameProfile = des.IsIFrameProfile()
			if isIFrameProfile {
				break
			}
		}
		if isIFrameProfile {
			break
		}
	}
	if isIFrameProfile {
		t.Errorf("Negative I-Frame Stream failed. Not supposed to be an I-Frame stream.")
	}
}

func TestIsPMT(t *testing.T) {
	patPkt := parseHexString("4740003001000000b00d0001c100000001e1e02d507804ffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	patPayload, _ := packet.Payload(patPkt)
	pat, _ := NewPAT(patPayload)

	if pat == nil {
		t.Error("Couldn't load the PAT")
	}

	pmt := parseHexString("4741e03001000002b0480001c10000e1e1f0050e03c004751be1e1f016970028" +
		"044d401f3fe907108302808502800e03c003350fe1e2f01697000a04656e6700" +
		"e907108302808502800e03c00104db121f57ffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	notPMT := parseHexString("4741e13117f200014307ff050fdf0d45425030c8dae4dd8000000000000001e0" +
		"000084d00d31000bab4111000b93cb80054700000001091000000001674d401f" +
		"ba202833f3e022000007d20001d4c1c040020f400041eb4d4601f18311200000" +
		"000168ebef20000000010600068232993c76c08000000001060447b500314741" +
		"393403d4fffc8080fd8080fa0000fa0000fa0000fa0000fa0000fa0000fa0000" +
		"fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa")

	if isPMTExpectTrue, _ := IsPMT(pmt, pat); isPMTExpectTrue == false {
		t.Error("PMT packet should be counted as a PMT")
	}

	if isPMTExpectFalse, _ := IsPMT(notPMT, pat); isPMTExpectFalse == true {
		t.Error("EBP packet should not be counted as a PMT")
	}
}

func TestIsPMTErrorConditions(t *testing.T) {
	// Test nil PAT

	pmt := parseHexString("4741e03001000002b0480001c10000e1e1f0050e03c004751be1e1f016970028" +
		"044d401f3fe907108302808502800e03c003350fe1e2f01697000a04656e6700" +
		"e907108302808502800e03c00104db121f57ffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	isPMTExpectFalse, errExpectInvalidArg := IsPMT(pmt, nil)
	if isPMTExpectFalse == true {
		t.Error("nil PAT should return false for any PMT")
	}

	if errExpectInvalidArg != gots.ErrNilPAT {
		t.Error("Nil Pat should return nil pat error")
	}

	patPkt := parseHexString("4740003001000000b00d0001c100000001e1e02d507804ffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	patPayload, _ := packet.Payload(patPkt)
	pat, _ := NewPAT(patPayload)

	if pat == nil {
		t.Error("Couldn't load the PAT")
	}
}

func TestReadPMTForSmoke(t *testing.T) {
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

	pid := uint16(598)
	pmt, err := ReadPMT(r, pid)
	if err != nil {
		t.Fatalf("Unexpected error reading PMT: %v", err)
	}
	// sanity check (tests integration a bit)
	if len(pmt.ElementaryStreams()) != 2 {
		t.Errorf("PMT read is invalid, did not have expected number of streams")
	}
}

func TestReadPMTIncomplete(t *testing.T) {
	bs, _ := hex.DecodeString("474000100000b00d0001c100000001e256f803e71bfffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ff4742") // incomplete PMT packet
	r := bytes.NewReader(bs)

	pid := uint16(598)
	_, err := ReadPMT(r, pid)
	if err == nil {
		t.Errorf("Expected to get error reading PMT, but did not")
	}
}

func TestReadPMTSCTE(t *testing.T) {
	bs, _ := hex.DecodeString("47403b1b00c0001500000100810000000000000100000000002f832c69ff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffff47403b1c0002b0b20001cb0000f13df01809044749e1" +
		"0b050441432d330504454143330504435545491bf13df0102a027e1f9700" +
		"e9080c001f418507d04181f13ef00f810706380fff1f003f0a04656e6700" +
		"81f13ff00f8107061003ff1f003f0a047370610086f140f00f8a01009700" +
		"e9080c001f418507d041c0f141f012050445545631a100e9080c001f4185" +
		"07d041c0f142f013050445545631a20100e9080c001f418507d041c0f164" +
		"f008bf06496e766964690bcfa64bffff") // two PMT packets, first with SCTE 0xc0 table only
	r := bytes.NewReader(bs)

	pid := uint16(59)
	pmt, err := ReadPMT(r, pid)
	if err != nil {
		t.Errorf("Unexpected error reading PMT: %v", err)
		return
	}
	// sanity check (tests integration a bit)
	if len(pmt.ElementaryStreams()) != 7 {
		t.Errorf("PMT read is invalid, did not have expected number of streams")
	}
}

func TestReadPMT_MultipleTables_MultiplePackets(t *testing.T) {
	bs, _ := hex.DecodeString("47403b1e00c0001500000100610000000000000100000000" +
		"0035e3e2d702b0b20001c50000eefdf01809044749e10b05" +
		"0441432d330504454143330504435545491beefdf0102a02" +
		"7e1f9700e9080c001f418507d04181eefef00f810706380f" +
		"ff1f003f0a04656e670081eefff00f8107061003ff1f003f" +
		"0a047370610086ef00f00f8a01009700e9080c001f418507" +
		"d041c0ef01f012050445545631a100e9080c001f418507d0" +
		"41c0ef02f013050445545631a20100e9080c001f47003b1f" +
		"418507d041c0ef03f008bf06496e76696469a5cff3afffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ffffffffffffffffffffffffffffffff") // two tables (0xc0 and 0x2) combined across two packets
	r := bytes.NewReader(bs)

	pid := uint16(59)
	pmt, err := ReadPMT(r, pid)
	if err != nil {
		t.Errorf("Unexpected error reading PMT: %v", err)
		return
	}
	if len(pmt.ElementaryStreams()) != 7 {
		t.Errorf("PMT read is invalid, did not have expected number of streams")
	}
}

func TestIsDolbyATMOS(t *testing.T) {
	d, _ := hex.DecodeString("00ff")
	desc := &pmtDescriptor{tag: uint8(204), data: d}
	if desc.IsDolbyATMOS() {
		t.Errorf("Expected not Dolby ATMOS")
	}

	d, _ = hex.DecodeString("ffffc000000000000000000000000000000000000000000000")
	desc = &pmtDescriptor{tag: uint8(204), data: d}
	if desc.IsDolbyATMOS() {
		t.Errorf("Expected not Dolby ATMOS")
	}

	d, _ = hex.DecodeString("0000ff0001")
	desc = &pmtDescriptor{tag: uint8(204), data: d}
	if !desc.IsDolbyATMOS() {
		t.Errorf("Expected Dolby ATMOS")
	}

	d, _ = hex.DecodeString("ffffc000000000000000000000000000000000000000000001")
	desc = &pmtDescriptor{tag: uint8(204), data: d}
	if !desc.IsDolbyATMOS() {
		t.Errorf("Expected Dolby ATMOS")
	}

	desc = &pmtDescriptor{tag: uint8(2), data: d}
	if desc.IsDolbyATMOS() {
		t.Errorf("Expected not Dolby ATMOS")
	}

	b1 := []byte{0x47, 0x41, 0xe0, 0x1f, 0x00, 0x02, 0xb0, 0xd7, 0x00,
		0x01, 0xc1, 0x00, 0x00, 0xe1, 0xe1, 0xf0, 0x0b, 0x05, 0x04, 0x47,
		0x41, 0x39, 0x34, 0x0e, 0x03, 0xc0, 0xbf, 0x0e, 0x24, 0xe1, 0xe1,
		0xf0, 0x25, 0x06, 0x01, 0x02, 0x38, 0x0f, 0x02, 0x20, 0x00, 0x00,
		0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x99, 0x9f, 0x1f, 0x1f,
		0x97, 0x00, 0x0e, 0x03, 0xc0, 0xb4, 0x70, 0xe9, 0x08, 0x0c, 0x00,
		0x0f, 0xa1, 0x83, 0x03, 0xe9, 0x41, 0x87, 0xe1, 0xe2, 0xf0, 0x26,
		0x97, 0x00, 0x0e, 0x03, 0xc0, 0x06, 0xe2, 0xe9, 0x08, 0x0c, 0x00,
		0x0f, 0xa1, 0x83, 0x03, 0xe9, 0x41, 0x05, 0x04, 0x41, 0x43, 0x2d,
		0x33, 0xcc, 0x07, 0x80, 0xc5, 0xa0, 0x65, 0x6e, 0x67, 0x01, 0x0a,
		0x04, 0x65, 0x6e, 0x67, 0x00, 0x87, 0xe1, 0xe3, 0xf0, 0x25, 0x97,
		0x00, 0x0e, 0x03, 0xc0, 0x01, 0x74, 0xe9, 0x08, 0x0c, 0x00, 0x0f,
		0xa1, 0x83, 0x03, 0xe9, 0x41, 0x05, 0x04, 0x41, 0x43, 0x2d, 0x33,
		0xcc, 0x06, 0x80, 0xc2, 0xa0, 0x73, 0x70, 0x61, 0x0a, 0x04, 0x73,
		0x70, 0x61, 0x00, 0x0f, 0xe1, 0xe8, 0xf0, 0x1b, 0x7c, 0x02, 0x59,
		0x00, 0x0a, 0x04, 0x65, 0x6e, 0x67, 0x00, 0x97, 0x00, 0x0e, 0x03,
		0xc0, 0x01, 0x25, 0xe9, 0x08, 0x0c, 0x00, 0x0f, 0xa1, 0x83, 0x03,
		0xe9, 0x41, 0x0f}
	b2 := []byte{0x47, 0x01, 0xe0, 0x10, 0xe1, 0xe9, 0xf0, 0x1b, 0x7c,
		0x02, 0x59, 0x00, 0x0a, 0x04, 0x73, 0x70, 0x61, 0x00, 0x97, 0x00,
		0x0e, 0x03, 0xc0, 0x01, 0x25, 0xe9, 0x08, 0x0c, 0x00, 0x0f, 0xa1,
		0x83, 0x03, 0xe9, 0x41, 0x38, 0x63, 0x37, 0x8f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff}

	acc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(b1)
	acc.Add(b2)
	payload, err := acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err := NewPMT(payload)
	if err != nil {
		t.Error(err)
	}
	var hasATMOS bool
	for _, es := range pmt.ElementaryStreams() {
		for _, des := range es.Descriptors() {
			hasATMOS = des.IsDolbyATMOS()
			if hasATMOS {
				break
			}
		}
		if hasATMOS {
			break
		}
	}
	if !hasATMOS {
		t.Errorf("Positive Dolby ATMOS Stream failed. Supposed to be a Dolby ATMOS stream.")
	}

	b := []byte{0x47, 0x41, 0xe0, 0x1f, 0x00, 0x02, 0xb0, 0x8d, 0x00, 0x01,
		0xc1, 0x00, 0x00, 0xe1, 0xe1, 0xf0, 0x0b, 0x05, 0x04, 0x47, 0x41,
		0x39, 0x34, 0x0e, 0x03, 0xc0, 0xb9, 0xe8, 0x24, 0xe1, 0xe1, 0xf0,
		0x25, 0x06, 0x01, 0x02, 0x38, 0x0f, 0x02, 0x20, 0x00, 0x00, 0x00,
		0xb0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x99, 0x9f, 0x1f, 0x1f, 0x97,
		0x00, 0x0e, 0x03, 0xc0, 0xb4, 0x70, 0xe9, 0x08, 0x0c, 0x00, 0x0f,
		0xa1, 0x83, 0x03, 0xe9, 0x41, 0x87, 0xe1, 0xe2, 0xf0, 0x26, 0x97,
		0x00, 0x0e, 0x03, 0xc0, 0x04, 0x54, 0xe9, 0x08, 0x0c, 0x00, 0x0f,
		0xa1, 0x83, 0x03, 0xe9, 0x41, 0x05, 0x04, 0x41, 0x43, 0x2d, 0x33,
		0xcc, 0x07, 0x80, 0xf4, 0xa0, 0x65, 0x6e, 0x67, 0x01, 0x0a, 0x04,
		0x65, 0x6e, 0x67, 0x00, 0x0f, 0xe1, 0xe3, 0xf0, 0x1b, 0x7c, 0x02,
		0x59, 0x00, 0x0a, 0x04, 0x65, 0x6e, 0x67, 0x00, 0x97, 0x00, 0x0e,
		0x03, 0xc0, 0x01, 0x25, 0xe9, 0x08, 0x0c, 0x00, 0x0f, 0xa1, 0x83,
		0x03, 0xe9, 0x41, 0xc3, 0x50, 0x69, 0xb4, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff}

	acc = packet.NewAccumulator(PmtAccumulatorDoneFunc)
	acc.Add(b)
	payload, err = acc.Parse()
	if err != nil {
		t.Error(err)
	}
	pmt, err = NewPMT(payload)
	if err != nil {
		t.Error(err)
	}
	hasATMOS = false
	for _, es := range pmt.ElementaryStreams() {
		for _, des := range es.Descriptors() {
			hasATMOS = des.IsDolbyATMOS()
			if hasATMOS {
				break
			}
		}
		if hasATMOS {
			break
		}
	}
	if !hasATMOS {
		t.Errorf("Positive Dolby ATMOS Stream failed. Supposed to be a Dolby ATMOS stream.")
	}
}
