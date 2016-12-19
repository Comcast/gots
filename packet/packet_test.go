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

package packet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
func TestPayloadUnitStartIndicatorTrue(t *testing.T) {
	packet, _ := hex.DecodeString(
		"474000130000b00d0001c700000001e0642273423bffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	expected := true
	if pusi, err := PayloadUnitStartIndicator(packet); pusi != expected || err != nil {
		t.Errorf("PayloadUnitStartIndicator() = %t, want %t err = %v", pusi, expected, err)
	}
}
func TestPayloadUnitStartIndicatorFalse(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4700673b7000ffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffd55e98825c522faa6d37de" +
			"cb2e84e04544cd54b9ffb497e788f2b308d1a71d4f2bce4c18c563b92ecd954b" +
			"558e7ca55796ca56ed4020812c21f1013ff5497f875897f58ffeeb1c")
	expected := false
	if pusi, err := PayloadUnitStartIndicator(packet); pusi != expected || err != nil {
		t.Errorf("PayloadUnitStartIndicator() = %t, want %t err = %v", pusi, expected, err)
	}
}

func TestPid(t *testing.T) {
	packet, _ := hex.DecodeString(
		"47406618000001c000f280800523fae5b8a3fff94c801d4010210994fd959f4b" +
			"6108806a912e4b972d025c92429595817016dca64a18e7fc5c271bb40a0f9150" +
			"3c0057776bdd66c0e9ab2ba7614de80ee468cc6e860846241710cfda6dabc569" +
			"79279065d30a93c8d5584d12b87b35938e18f2868f149f3dec38cae665db77bd" +
			"0ba9b7a659363d7347d22f835b4e53f6472f01be53d7df28ea7f1764972f5549" +
			"34096bd6bf42eabe1dff1c59e0cc55a716b6a40618b3305b45779c31")
	expected := uint16(102)
	if pid, err := Pid(packet); pid != expected || err != nil {
		t.Errorf("Pid() = %d, want %d err=%v", pid, expected, err)
	}
}

func TestPidGreaterThen255(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4701221B000001c000f280800523fae5b8a3fff94c801d4010210994fd959f4b" +
			"6108806a912e4b972d025c92429595817016dca64a18e7fc5c271bb40a0f9150" +
			"3c0057776bdd66c0e9ab2ba7614de80ee468cc6e860846241710cfda6dabc569" +
			"79279065d30a93c8d5584d12b87b35938e18f2868f149f3dec38cae665db77bd" +
			"0ba9b7a659363d7347d22f835b4e53f6472f01be53d7df28ea7f1764972f5549" +
			"34096bd6bf42eabe1dff1c59e0cc55a716b6a40618b3305b45779c31")
	expected := uint16(290)
	if pid, err := Pid(packet); pid != expected || err != nil {
		t.Errorf("Pid() = %d, want %d err=%v", pid, expected, err)
	}
}

func TestContainsPayloadTrue(t *testing.T) {
	packet, _ := hex.DecodeString(
		"47406618000001c000f280800523fae5b8a3fff94c801d4010210994fd959f4b" +
			"6108806a912e4b972d025c92429595817016dca64a18e7fc5c271bb40a0f9150" +
			"3c0057776bdd66c0e9ab2ba7614de80ee468cc6e860846241710cfda6dabc569" +
			"79279065d30a93c8d5584d12b87b35938e18f2868f149f3dec38cae665db77bd" +
			"0ba9b7a659363d7347d22f835b4e53f6472f01be53d7df28ea7f1764972f5549" +
			"34096bd6bf42eabe1dff1c59e0cc55a716b6a40618b3305b45779c31")
	expected := true
	if containsPayload, err := ContainsPayload(packet); containsPayload != expected || err != nil {
		t.Errorf("ContainsPayload() = %t, want %t err=%v", containsPayload, expected, err)
	}
}

func TestContainsPayloadFalse(t *testing.T) {
	packet, _ := hex.DecodeString(
		"47006523b7103f5c99597ef7ffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	expected := false
	if containsPayload, err := ContainsPayload(packet); containsPayload != expected || err != nil {
		t.Errorf("ContainsPayload() = %t, want %t", containsPayload, expected)
	}
}

func TestContinuityCounter(t *testing.T) {
	packet, _ := hex.DecodeString(
		"47006518dc0eff960f094176e794721d00cfedc13c1b039abf71e0f16bfeef88" +
			"de1d1901a576793da53551cfc53363e00be1417c08383ce8bc51efda4c4a465c" +
			"9aee27f76997169968829cf3343253c16243f7c21602cb2161767fda0485d4de" +
			"87bf26954b148497393886b854288985a50059dd5cbfdf61b3f701793cd3bdf0" +
			"43e5d197998e05e2a50f590f6923d45490c81750d5f603643f974a2bde5f4812" +
			"749dd96f61281aca2cb496c01f01e3152fcba48c2ec78314ab21da3b")
	expected := uint8(8)
	if cc, err := ContinuityCounter(packet); cc != expected || err != nil {
		t.Errorf("ContinuityCounter() = %d, want %d", cc, expected)
	}
}

func TestZeroLenthAdaptationField(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4701e1320034fcabf65d866a87eca0195db5ce1dcb6e0f75ba45a351722714db" +
			"a013cea9665e9e1866b13431429454a37cb5663ea00353624c5d1f84c9463651" +
			"634497dd837080b99ddf4bb26242f18d22ecd74dde47cd84041e5df3f0c57c40" +
			"e1c6782394c5dbb02a1896ca6712dc232a53958f47596c570c70f90e44303188" +
			"89ace8b8b378a515b088341942220c44578c157ee4313d123db73ec2a2726a29" +
			"9ab1852b9314ae15fad86177607b75be718f0c07d22400845160d980")

	expected := true

	if hasadapt, err := ContainsAdaptationField(packet); hasadapt != expected || err != nil {
		t.Errorf("ContainsAdaptationField() = %v, want %v (%v)", hasadapt, expected, err)
	}
}

func TestPayloadWhenPacketHasNoAdaptationField(t *testing.T) {
	packet, _ := hex.DecodeString(
		"47006518dc0eff960f094176e794721d00cfedc13c1b039abf71e0f16bfeef88" +
			"de1d1901a576793da53551cfc53363e00be1417c08383ce8bc51efda4c4a465c" +
			"9aee27f76997169968829cf3343253c16243f7c21602cb2161767fda0485d4de" +
			"87bf26954b148497393886b854288985a50059dd5cbfdf61b3f701793cd3bdf0" +
			"43e5d197998e05e2a50f590f6923d45490c81750d5f603643f974a2bde5f4812" +
			"749dd96f61281aca2cb496c01f01e3152fcba48c2ec78314ab21da3b")

	expected, _ := hex.DecodeString(
		"dc0eff960f094176e794721d00cfedc13c1b039abf71e0f16bfeef88de1d1901" +
			"a576793da53551cfc53363e00be1417c08383ce8bc51efda4c4a465c9aee27f7" +
			"6997169968829cf3343253c16243f7c21602cb2161767fda0485d4de87bf2695" +
			"4b148497393886b854288985a50059dd5cbfdf61b3f701793cd3bdf043e5d197" +
			"998e05e2a50f590f6923d45490c81750d5f603643f974a2bde5f4812749dd96f" +
			"61281aca2cb496c01f01e3152fcba48c2ec78314ab21da3b")

	if payload, err := Payload(packet); !(bytes.Equal(payload, expected)) || err != nil {
		t.Errorf("Payload() = %x, want %x err=%v", payload, expected, err)
	}
}

func TestPayloadWhenPacketHasAdaptationField(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	expected, _ := hex.DecodeString(
		"000001e0000084c00a33faf9760713faf900b900000001091000000001274d40" +
			"1f9a6281004b602d1000003e90000ea60e8601d400057e4bbcb8280000000128ee38" +
			"8000000001060007818a378085f8c00104007820100601c40411b500314741393403" +
			"c2fffd2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f" +
			"78053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	if payload, err := Payload(packet); !(bytes.Equal(payload, expected)) || err != nil {
		t.Errorf("Payload() = %x, want %x err=%v", payload, expected, err)
	}
}

func TestIncrementCC(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4700673b7000ffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffd55e98825c522faa6d37de" +
			"cb2e84e04544cd54b9ffb497e788f2b308d1a71d4f2bce4c18c563b92ecd954b" +
			"558e7ca55796ca56ed4020812c21f1013ff5497f875897f58ffeeb1c")
	packet[3] = byte(0x00)
	newPacket, err := IncrementCC(packet)
	if err != nil {
		t.Error(err)
	}
	expected := uint8(1)
	if expected != newPacket[3] {
		t.Errorf("CC= %x, want %x", newPacket[3], expected)
	}

}

func TestBadLength(t *testing.T) {
	packet, _ := hex.DecodeString("4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084")
	_, err := Header(packet)

	if err == nil {
		t.Errorf("BadLength, expected error from new packet")
	}
}

func TestIncrementCCFunc(t *testing.T) {
	for i := byte(0); i < 16; i++ {
		if i == 15 && increment4BitInt(i) != 0 {
			t.Errorf("IncrementingCC from 15 to rollover did not cause a 0")
		}
		if i == 15 {
			continue
		}
		res := increment4BitInt(i)
		if res != i+1 {
			t.Errorf("IncrementingCC from %d expected %d was %d", i, i+1, res)
		}
	}
}

func TestContainsAdaptationField(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4700663a7700ffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0b82aaa" +
			"bd5a13a6b27a23c4e556bd78ccdb05c72a3a5cac278d76d89aae3de5241728ea" +
			"ea79344f6ff95e63bd6060c9462c42b33a89b3fcff000480003e71c0")
	if exists, _ := ContainsAdaptationField(packet); !exists {
		t.Error("Did not correctly find presence of adaptation field")
	}
}

func TestEqualsNilPacket(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")
	if Equal(packet, nil) {
		t.Error("Nil packet should not be equal to a non-nil packet")
	}
}

func TestEqualsIdenticalPackets(t *testing.T) {
	packet, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")
	same := packet[:]
	if !Equal(packet, same) {
		t.Errorf("Identical packets are different p1%v p2%v", packet, same)
	}
}

func TestEqualsHeadersNotEqual(t *testing.T) {
	packet1, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	// Same as above, but with the MPEG-TS headers TEI bit flipped.
	packet2, _ := hex.DecodeString(
		"4780653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	if Equal(packet1, packet2) {
		t.Errorf("Packets should be different\n p1%v\n p2%v", packet1, packet2)
	}
}

func TestNullPacketIsNull(t *testing.T) {
	p, _ := hex.DecodeString(
		"471fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	if isNull, _ := IsNull(p); isNull == false {
		pid, e := Pid(p)
		t.Errorf("Packets with PID == %d should be null. PID was %v and error was %v", NullPacketPid, pid, e)
	}
}

func TestNonNullPacketIsNotNull(t *testing.T) {
	packet1, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	if isNull, _ := IsNull(packet1); isNull == true {
		t.Errorf("Packets with PID != %d should not be null.", NullPacketPid)
	}
}

func TestIsPat(t *testing.T) {
	pat, _ := hex.DecodeString(
		"4740001f0000b00d0031e100000001e064bfcd282fffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	if isPat, _ := IsPat(pat); isPat == false {
		t.Error("PAT packet should be counted as a PAT")
	}

	notPat, _ := hex.DecodeString(
		"4740653214723f5d09c67ec90ca90ad800d6ae02c11e66772d000001e0000084" +
			"c00a33faf9760713faf900b900000001091000000001274d401f9a6281004b60" +
			"2d1000003e90000ea60e8601d400057e4bbcb8280000000128ee388000000001" +
			"060007818a378085f8c00104007820100601c40411b500314741393403c2fffd" +
			"2980fc8080ff800000000125b80100017fb2c69de69e51f57c4a1b8623115f78" +
			"053598e7f47c066bf03c90c6233c0405369fd5f8e20957e40437f784")

	if isPat2, _ := IsPat(notPat); isPat2 == true {
		t.Error("Non PAT Packet shouldn't be counted as a PAT")
	}
}
