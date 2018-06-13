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
package pes

import (
	"encoding/hex"
	"testing"

	"github.com/Comcast/gots/packet"
)

func TestPESHeader(t *testing.T) {

	pkt, _ := hex.DecodeString(
		"4740661a000001c006ff80800521dee9ca57fff94c801d2000210995341d9d43" +
			"61089848180b0884626048901425ddc09249220129d2fce728111c987e67ecb7" +
			"4284af5099181d8cd095b841b0c7539ad6c06260536e137615560052369fc984" +
			"0b3af532418b3924a6b28d6208a6a9e3d22d533ec89646246c734a696407a95e" +
			"3bf230404a4ad0000201038cf0c6a2e32abda45b7effe9f79280a137ed120fd3" +
			"bd8252e07cddadbe6d2b60084208500e06f6ceb6acf2c43011c5938b")
	pesBytes, err := packet.PESHeader(pkt)
	if err != nil {
		t.Errorf("Expected that this packet is a PES")
	}
	pes, err := NewPESHeader(pesBytes)
	if err != nil {
		t.Error(err)
	}

	expected := uint64(934962475)
	if pes.PTS() != expected {
		t.Errorf("Invalid pts. Expected: %d, Actual: %d", expected, pes.PTS())
	}
	if pes.StreamId() != 0xc0 {
		t.Errorf("Invalid stream id. Expected: %d, Actual: %d", 0xc0, pes.StreamId())
	}
	if pes.PacketStartCodePrefix() != uint32(0x000001) {
		t.Errorf("Invalid start code prefix. Expected: %d, Actual: %d", uint32(0x000001), pes.PacketStartCodePrefix())
	}
}

func TestPESHeader2(t *testing.T) {

	pkt, _ := hex.DecodeString(
		"4740651C000001E0000084C00A39EFF33A7519EFF30B89000000010950000000" +
			"01060104001A20100411B500314741393403C2FFFD8080FC942FFF8000000001" +
			"21A81C29145C6FEB86EB239E2EE231302CF5163D32D183B7822FE37E7FB84549" +
			"DC1D08780834029F139BDD36E9BBC25B18AE4DE5F508036AEDB9E8A321B93288" +
			"4EEF1482E6C77B31E92ADF3BC0D275E5D40864FD3A9806ABC74B98B0E3255EC1" +
			"B1C157068EF46E15ED82E7D7C1C0538C4B5B7AF39AEC09386939FE1C")

	pesBytes, err := packet.PESHeader(pkt)
	if err != nil {
		t.Errorf("Expected that this packet is a PES")
	}
	pes, err := NewPESHeader(pesBytes)
	if err != nil {
		t.Error(err)
	}

	expected := uint64(5301378362)
	if pes.PTS() != expected {
		t.Errorf("Invalid pts. Expected: %d, Actual: %d", expected, pes.PTS())
	}
	if pes.StreamId() != 0xE0 {
		t.Errorf("Invalid stream id. Expected: %d, Actual: %d", 0xe0, pes.StreamId())
	}
	if pes.DataAligned() != true {
		t.Error("PES header read incorrect data alignment flag")
	}
}

func TestNewPESHeaderMissingBytes(t *testing.T) {

	// Actual data from Cisco Transcoder (AMC channel).  Below packet was causing
	// index out of bounds exception.  It has the PES prefix code but we were not
	// checking to see if it's a PUSI to begin with
	pkt, _ := hex.DecodeString(
		"47006531b300ffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffff0000010b")
	_, err := packet.PESHeader(pkt)
	if err == nil {
		t.Errorf("Expected that this packet is not a PES since PUSI is 0")
	}
}

func TestPESHeaderTS(t *testing.T) {

	pkt, _ := hex.DecodeString(
		"4752a31c000001e0000080c00a210005bf21210005a7ab000001000697fffb80" +
			"000001b5844ffb9400000001b24741393403d4fffc8080fd8fdffa0000fa0000" +
			"fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa0000fa00" +
			"00fa0000fa0000fa0000fa0000fa0000ff000001014a24afffa4e8b836d7eeee" +
			"4dafded260dab9688b2a0d89bed7fd3ad106c1b6bfe5a24a20d89b572ca92544" +
			"389b572ca7b441b176bffebd06c5daffe8bd06c9b5fbb8364da6ffad")

	pesBytes, err := packet.PESHeader(pkt)
	if err != nil {
		t.Errorf("Expected that this packet is a PES")
	}
	pes, err := NewPESHeader(pesBytes)
	if err != nil {
		t.Error(err)
	}

	expectedPTS := uint64(90000)
	if pes.PTS() != expectedPTS {
		t.Errorf("Invalid pts. Expected: %d, Actual: %d", expectedPTS, pes.PTS())
	}

	expectedDTS := uint64(86997)
	if !pes.HasDTS() {
		t.Errorf("Invalid dts indicator.")
	}
	if pes.DTS() != expectedDTS {
		t.Errorf("Invalid dts. Expected: %d, Actual: %d", expectedDTS, pes.DTS())
	}
}
