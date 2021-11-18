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
	"testing"
)

const (
	testPacket = "47000010" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"00000000000000000000000000000000"

	testPacket2 = "471FFF30" +
		"530246999999999999999999999999999999999999999999" +
		"999999999999999999999999999999999999999999999999" +
		"999999999999999999999999999999999999999999999999" +
		"99FFFFFFFFFFFFFFFFFFFFFF777777777777777777777777" +
		"777777777777777777777777777777777777777777777777" +
		"777777777777777777777777777777777777777777777777" +
		"777777777777777777777777777777777777777777777777" +
		"77777777777777777777777777777777"
)

func createPacketEmptyPayload(t *testing.T, header string) *Packet {
	headerBytes, _ := hex.DecodeString(header)
	bodyBytes := make([]byte, 188-len(headerBytes))
	packetBytes := append(headerBytes, bodyBytes...)

	p, err := FromBytes(packetBytes)
	if err != nil {
		t.Error("packet error checking failed")
	}
	return p
}

func createPacketEmptyAdaptationField(t *testing.T, header string) *Packet {
	headerBytes, _ := hex.DecodeString(header)
	AFBytes := make([]byte, 188)
	AFBytes[4] = 183
	AFBytes[5] = 0
	for i := 6; i < len(AFBytes); i++ {
		AFBytes[i] = 0xFF
	}
	AFBytes = AFBytes[len(headerBytes):188]
	packetBytes := append(headerBytes, AFBytes...)

	p, err := FromBytes(packetBytes)
	if err != nil {
		t.Error("packet error checking failed")
	}
	return p
}

func TestFromBytes(t *testing.T) {
	bytes, _ := hex.DecodeString(testPacket)
	p, err := FromBytes(bytes)
	if err != nil {
		t.Error("eacket error checking failed")
	}
	if (len(bytes) != 188) && (len(p) != 188) {
		t.Error("packets are not 188 bytes")
		return
	}
	for i := range bytes {
		if bytes[i] != p[i] {
			t.Error("packet generated with FromBytes did not copy bytes correctly.")
			return
		}
	}
}

func TestNewPacket(t *testing.T) {
	target := createPacketEmptyPayload(t, "471FFF10")
	generated := New()
	if err := generated.CheckErrors(); err != nil {
		t.Error("Default packet has errors.")
	}
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nCreating a new packet failed.", generated, target)
	}
}

func TestSetTransportErrorIndicator(t *testing.T) {
	generated := New()

	target := createPacketEmptyPayload(t, "479FFF10")
	generated.SetTransportErrorIndicator(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport error indicator to true has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF10")
	generated.SetTransportErrorIndicator(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport error indicator to false has failed.", generated, target)
	}
}

func TestTransportErrorIndicator(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "479FFF10")
	tei := pkt.TransportErrorIndicator()
	if tei != true {
		t.Error("TransportErrorIndicator does not match expected.")
	}
}

func TestSetPayloadUnitStartIndicator(t *testing.T) {
	generated := New()

	target := createPacketEmptyPayload(t, "475FFF10")
	generated.SetPayloadUnitStartIndicator(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PUSI to true has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF10")
	generated.SetPayloadUnitStartIndicator(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PUSI to false has failed.", generated, target)
	}
}

func TestPayloadUnitStartIndicator(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "475FFF10")
	pusi := pkt.PayloadUnitStartIndicator()
	if pusi != true {
		t.Error("PayloadUnitStartIndicator does not match expected.")
	}
}

func TestSetTP(t *testing.T) {
	generated := New()

	target := createPacketEmptyPayload(t, "473FFF10")
	generated.SetTransportPriority(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport priority to true has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF10")
	generated.SetTransportPriority(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport priority to false has failed.", generated, target)
	}
}

func TestTP(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "473FFF10")
	tp := pkt.TransportPriority()
	if tp != true {
		t.Error("TransportPriority does not match expected.")
	}
}

func TestSetPID(t *testing.T) {
	generated := createPacketEmptyPayload(t, "47e09010")

	target := createPacketEmptyPayload(t, "47f76A10")
	generated.SetPID(0x176A)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PID to 0x176A has failed.", generated, target)
	}

	generated = New()

	target = createPacketEmptyPayload(t, "47000010")
	generated.SetPID(0x0000)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PID to 0x0000 has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471fec10")
	generated.SetPID(0x1fec)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PID to 0x1fec has failed.", generated, target)
	}
}

func TestPID(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "47344410")
	pid := pkt.PID()
	if pid != 0x1444 {
		t.Errorf("reads different PID than expected. expected: %d", 0x1444)
	}
}

func TestSetTransportScramblingControl(t *testing.T) {
	generated := New()

	target := createPacketEmptyPayload(t, "471FFFD0")
	generated.SetTransportScramblingControl(ScrambleOddKeyFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Scrambling Control to ScrambleOddKeyFlag has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF90")
	generated.SetTransportScramblingControl(ScrambleEvenKeyFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Scrambling Control to ScrambleEvenKeyFlag has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF10")
	generated.SetTransportScramblingControl(NoScrambleFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Scrambling Control to NoScrambleFlag has failed.", generated, target)
	}
}

func TestTransportScramblingControl(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF90")
	tsc := pkt.TransportScramblingControl()
	if tsc != ScrambleEvenKeyFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: ScrambleEvenKeyFlag")
	}

	pkt = createPacketEmptyPayload(t, "471FFFD0")
	tsc = pkt.TransportScramblingControl()
	if tsc != ScrambleOddKeyFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: ScrambleOddKeyFlag")
	}

	pkt = createPacketEmptyPayload(t, "471FFF10")
	tsc = pkt.TransportScramblingControl()
	if tsc != NoScrambleFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: NoScrambleFlag")
	}
}

func TestSetAdaptationFieldControl(t *testing.T) {
	generated := New()

	target := createPacketEmptyPayload(t, "471FFF10")
	generated.SetAdaptationFieldControl(PayloadFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to PayloadFlag has failed.", generated, target)
	}

	target = createPacketEmptyAdaptationField(t, "471FFF30B6")
	generated.SetAdaptationFieldControl(PayloadAndAdaptationFieldFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to PayloadAndAdaptationFieldFlag has failed.", generated, target)
	}

	target = createPacketEmptyAdaptationField(t, "471FFF20")
	generated.SetAdaptationFieldControl(PayloadFlag)
	generated.SetAdaptationFieldControl(AdaptationFieldFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to AdaptationFieldFlag has failed.", generated, target)
	}
}

func TestAdaptationFieldControl(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF9000")
	asc := pkt.AdaptationFieldControl()
	if asc != PayloadFlag {
		t.Errorf("reads different adaptation field control than expected. expected: PayloadFlag")
	}
	pkt = createPacketEmptyAdaptationField(t, "471FFFB001")
	asc = pkt.AdaptationFieldControl()
	if asc != PayloadAndAdaptationFieldFlag {
		t.Errorf("reads different adaptation field control than expected. expected: PayloadAndAdaptationFieldFlag")
	}
	pkt = createPacketEmptyAdaptationField(t, "471FFFA001")
	asc = pkt.AdaptationFieldControl()
	if asc != AdaptationFieldFlag {
		t.Errorf("reads different adaptation field control than expected. expected: AdaptationFieldFlag")
	}
}

func TestSetContinuityCounter(t *testing.T) {
	target := createPacketEmptyPayload(t, "471FFF1f")
	generated := New()

	generated.SetContinuityCounter(15)

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Continuity Counter to 15 has failed.", generated, target)
	}
}

func TestContinuityCounterModify(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF19")
	cc := pkt.ContinuityCounter()
	if cc != 9 {
		t.Errorf("Reads different continuity counter than expected. ecpected: 9")
	}
}

func TestIncContinuityCounter(t *testing.T) {
	target := createPacketEmptyPayload(t, "471FFF10")
	generated := New()

	generated.SetContinuityCounter(0xFE) // cc = 14
	generated.IncContinuityCounter()     // cc = 15
	generated.IncContinuityCounter()     // cc = 0, overflow

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nContinuity Counter did not rollover as expected.", generated, target)
	}
}

func TestIsPAT(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "47000010")
	isPAT := pkt.IsPAT()
	if !isPAT {
		t.Errorf("packet should be a PAT")
	}
}

func TestIsNull(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF10")
	isNull := pkt.IsNull()
	if !isNull {
		t.Errorf("packet should be Null")
	}
}

func TestHasAdaptationField(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF100100")
	hasAF := pkt.HasAdaptationField()
	if hasAF {
		t.Errorf("Packet should not have Adaptation Field (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF200100")
	hasAF = pkt.HasAdaptationField()
	if !hasAF {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF300100")
	hasAF = pkt.HasAdaptationField()
	if !hasAF {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 11).")
	}
}

func TestHasPayload(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF10")
	hasPayload := pkt.HasPayload()
	if !hasPayload {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF20")
	hasPayload = pkt.HasPayload()
	if hasPayload {
		t.Errorf("Packet should not have Payload (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF30")
	hasPayload = pkt.HasPayload()
	if !hasPayload {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 11).")
	}
}

func TestHeaderBasic(t *testing.T) {
	target := createPacketEmptyPayload(t, "47EFA098")
	generated := New()

	generated.SetContinuityCounter(7)
	generated.IncContinuityCounter()
	generated.SetPID(4000)
	generated.SetTransportErrorIndicator(true)
	generated.SetPayloadUnitStartIndicator(true)
	generated.SetTransportPriority(true)
	generated.SetAdaptationFieldControl(PayloadFlag)
	generated.SetTransportScramblingControl(ScrambleEvenKeyFlag)

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nFields did not set successfully", generated, target)
	}

	if generated.ContinuityCounter() != 8 {
		t.Errorf("Reads different ContinuityCounter than expected.")
	}
	if generated.PID() != 4000 {
		t.Errorf("Reads different PID than expected.")
	}
	if generated.TransportErrorIndicator() != true {
		t.Errorf("Reads different TransportErrorIndicator than expected.")
	}
	if generated.PayloadUnitStartIndicator() != true {
		t.Errorf("Reads different PayloadUnitStartIndicator than expected.")
	}
	if generated.TransportPriority() != true {
		t.Errorf("Reads different TP than expected.")
	}
	if generated.AdaptationFieldControl() != PayloadFlag {
		t.Errorf("Reads different AdaptationFieldControl than expected.")
	}
	if generated.TransportScramblingControl() != ScrambleEvenKeyFlag {
		t.Errorf("Reads different TransportScramblingControl than expected.")
	}
}

func TestSetPayload(t *testing.T) {
	data, _ := hex.DecodeString(testPacket2)
	target, _ := FromBytes(data)
	payload := []byte{}
	tpd := []byte{}
	for i := 0; i < 188; i++ {
		payload = append(payload, 0x77)
	}
	for i := 0; i < 188; i++ {
		tpd = append(tpd, 0x99)
	}

	copyAF := NewAdaptationField()
	err := copyAF.SetHasTransportPrivateData(true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = copyAF.SetTransportPrivateData(tpd[:70])
	if err != nil {
		t.Error(err.Error())
		return
	}
	p := New()
	err = p.SetAdaptationFieldControl(AdaptationFieldFlag)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = p.SetAdaptationFieldControl(PayloadAndAdaptationFieldFlag)
	if err != nil {
		t.Error(err.Error())
		return
	}
	af, _ := p.AdaptationField()
	err = af.SetHasTransportPrivateData(true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = af.SetTransportPrivateData(tpd[:100])
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = p.SetPayload(payload[:70])
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = p.SetAdaptationField(copyAF)
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = p.SetPayload(payload[:100])
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !Equal(target, p) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the payload failed.", p, target)
	}
	payloadInPacket, _ := p.Payload()
	if !bytes.Equal(payload[:100], payloadInPacket) {
		t.Errorf("payload in packet is incorrect.")
	}
}

func BenchmarkNewStyleAllFields(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// create everything
		p := New()
		p.SetContinuityCounter(7)
		p.IncContinuityCounter()
		p.SetPID(4000)
		p.SetTransportErrorIndicator(false)
		p.SetPayloadUnitStartIndicator(true)
		p.SetTransportPriority(false)
		p.SetAdaptationFieldControl(AdaptationFieldFlag)
		p.SetTransportScramblingControl(ScrambleEvenKeyFlag)

		a, _ := p.AdaptationField()
		a.SetHasPCR(true)
		a.SetHasOPCR(true)
		a.SetHasSplicingPoint(true)
		a.SetHasTransportPrivateData(true)
		a.SetHasAdaptationFieldExtension(true)
		a.SetElementaryStreamPriority(true)
		a.SetRandomAccess(true)
		a.SetHasAdaptationFieldExtension(true)

		// read everything
		p.ContinuityCounter()
		p.HasAdaptationField()
		p.HasPayload()
		p.IsNull()
		p.IsPAT()
		p.IsPAT()
		p.PID()
		p.PayloadUnitStartIndicator()
		p.TransportErrorIndicator()
		p.TransportPriority()
		p.TransportScramblingControl()

		a, _ = p.AdaptationField()
		a.Length()
		a.Discontinuity()
		a.ElementaryStreamPriority()
		a.HasAdaptationFieldExtension()
		a.HasOPCR()
		a.HasPCR()
		a.HasSplicingPoint()
		a.HasTransportPrivateData()
		a.RandomAccess()
	}
}

func BenchmarkNewStyleCreate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		pkt := New()
		pkt.SetPID(13)
		pkt.SetAdaptationFieldControl(PayloadFlag)
		pkt.SetPayloadUnitStartIndicator(true)
		pkt.SetContinuityCounter(7)
	}
}

func BenchmarkOldStyleCreate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SetCC(
			Create(
				13,
				WithHasPayloadFlag,
				WithPUSI),
			7)
	}
}

func BenchmarkNewStyleRead(b *testing.B) {
	pkt := New()
	pkt.SetPID(13)
	pkt.SetAdaptationFieldControl(PayloadFlag)
	pkt.SetPayloadUnitStartIndicator(true)
	pkt.CheckErrors() // no errors possible

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pkt.ContinuityCounter()
		pkt.PayloadUnitStartIndicator()
		pkt.AdaptationFieldControl()
		pkt.PID()
	}
}

func BenchmarkOldStyleRead(b *testing.B) {
	pkt := New()
	pkt.SetPID(13)
	pkt.SetAdaptationFieldControl(PayloadFlag)
	pkt.SetPayloadUnitStartIndicator(true)
	pkt.CheckErrors() // no errors possible

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ContinuityCounter(pkt)
		PayloadUnitStartIndicator(pkt)
		ContainsPayload(pkt)
		Pid(pkt)
	}
}
