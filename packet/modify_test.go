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
	"encoding/hex"
	"github.com/Comcast/gots"
	"testing"
)

const (
	testPacket = "4700001" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000"
)

func createPacketEmptyPayload(t *testing.T, header string) (p Packet) {
	headerBytes, _ := hex.DecodeString(header)
	bodyBytes := make([]byte, 188-len(headerBytes))
	packetBytes := Packet(append(headerBytes, bodyBytes...))

	p, err := FromBytes(packetBytes)
	if err != nil {
		t.Error("packet error checking failed")
	}
	return
}

func createPacketEmptyAdaptationField(t *testing.T, header string) (p Packet) {
	headerBytes, _ := hex.DecodeString(header)
	AFBytes := make([]byte, 188)
	AFBytes[4] = 183
	AFBytes[5] = 0
	for i := 6; i < len(AFBytes); i++ {
		AFBytes[i] = 0xFF
	}
	AFBytes = AFBytes[len(headerBytes):188]
	packetBytes := Packet(append(headerBytes, AFBytes...))

	p, err := FromBytes(packetBytes)
	if err != nil {
		t.Error("packet error checking failed")
	}
	return
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
	generated := NewPacket()
	if err := generated.CheckErrors(); err != nil {
		t.Error("Default packet has errors.")
	}
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nCreating a new packet failed.", generated, target)
	}
}

func TestSetTransportErrorIndicator(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyPayload(t, "479FFF10")
	err := generated.SetTransportErrorIndicator(true)
	if err != nil {
		t.Errorf("Failed to set Transport Error indicator: %s", err.Error())
	}
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport error indicator to true has failed.", generated, target)
	}

	target = createPacketEmptyPayload(t, "471FFF10")
	err = generated.SetTransportErrorIndicator(false)
	if err != nil {
		t.Errorf("Failed to set Transport Error indicator: %s", err.Error())
	}
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the transport error indicator to false has failed.", generated, target)
	}
}

func TestTransportErrorIndicator(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "479FFF10")
	tei, err := pkt.TransportErrorIndicator()
	if err != nil {
		t.Error("Failed to read TransportErrorIndicator flag.")
		return
	}
	if tei != true {
		t.Error("TransportErrorIndicator does not match expected.")
	}
}

func TestSetPayloadUnitStartIndicator(t *testing.T) {
	generated := NewPacket()

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
	pusi, err := pkt.PayloadUnitStartIndicator()
	if err != nil {
		t.Error("Failed to read PayloadUnitStartIndicator flag.")
		return
	}
	if pusi != true {
		t.Error("PayloadUnitStartIndicator does not match expected.")
	}
}

func TestSetTP(t *testing.T) {
	generated := NewPacket()

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
	tp, err := pkt.TransportPriority()
	if err != nil {
		t.Error("Failed to read TransportPriority flag.")
		return
	}
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

	generated = NewPacket()

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
	pid, err := pkt.PID()
	if err != nil {
		t.Error("failed to read PID")
		return
	}
	if pid != 0x1444 {
		t.Errorf("reads different PID than expected. expected: %d", 0x1444)
	}
}

func TestSetTransportScramblingControl(t *testing.T) {
	generated := NewPacket()

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
	tsc, err := pkt.TransportScramblingControl()
	if err != nil {
		t.Error("failed to read transport scrambling control")
		return
	}
	if tsc != ScrambleEvenKeyFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: ScrambleEvenKeyFlag")
	}

	pkt = createPacketEmptyPayload(t, "471FFFD0")
	tsc, err = pkt.TransportScramblingControl()
	if err != nil {
		t.Error("failed to read transport scrambling control")
		return
	}
	if tsc != ScrambleOddKeyFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: ScrambleOddKeyFlag")
	}

	pkt = createPacketEmptyPayload(t, "471FFF10")
	tsc, err = pkt.TransportScramblingControl()
	if err != nil {
		t.Error("failed to read transport scrambling control")
		return
	}
	if tsc != NoScrambleFlag {
		t.Errorf("reads different transport scrambling control than expected. expected: NoScrambleFlag")
	}
}

func TestSetAdaptationFieldControl(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyPayload(t, "471FFF10")
	generated.SetAdaptationFieldControl(PayloadFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to PayloadFlag has failed.", generated, target)
	}

	target = createPacketEmptyAdaptationField(t, "471FFF30")
	generated.SetAdaptationFieldControl(PayloadAndAdaptationFieldFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to PayloadAndAdaptationFieldFlag has failed.", generated, target)
	}

	target = createPacketEmptyAdaptationField(t, "471FFF20")
	generated.SetAdaptationFieldControl(AdaptationFieldFlag)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Control to AdaptationFieldFlag has failed.", generated, target)
	}
}

func TestAdaptationFieldControl(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF9000")
	asc, err := pkt.AdaptationFieldControl()
	if err != nil {
		t.Error("failed to read adaptation field control")
		return
	}
	if asc != PayloadFlag {
		t.Errorf("reads different adaptation field control than expected. expected: PayloadFlag")
	}
	pkt = createPacketEmptyAdaptationField(t, "471FFFB001")
	asc, err = pkt.AdaptationFieldControl()
	if err != nil {
		t.Error("failed to read adaptation field control")
		return
	}
	if asc != PayloadAndAdaptationFieldFlag {
		t.Errorf("reads different adaptation field control than expected. expected: PayloadAndAdaptationFieldFlag")
	}
	pkt = createPacketEmptyAdaptationField(t, "471FFFA001")
	asc, err = pkt.AdaptationFieldControl()
	if err != nil {
		t.Error("failed to read adaptation field control")
		return
	}
	if asc != AdaptationFieldFlag {
		t.Errorf("reads different adaptation field control than expected. expected: AdaptationFieldFlag")
	}
}

func TestSetContinuityCounter(t *testing.T) {
	target := createPacketEmptyPayload(t, "471FFF1f")
	generated := NewPacket()

	generated.SetContinuityCounter(15)

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Continuity Counter to 15 has failed.", generated, target)
	}
}

func TestContinuityCounterModify(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF19")
	cc, err := pkt.ContinuityCounter()
	if err != nil {
		t.Error("failed to read continuity counter")
		return
	}
	if cc != 9 {
		t.Errorf("Reads different continuity counter than expected. ecpected: 9")
	}
}

func TestIncContinuityCounter(t *testing.T) {
	target := createPacketEmptyPayload(t, "471FFF10")
	generated := NewPacket()

	generated.SetContinuityCounter(0xFE) // cc = 14
	generated.IncContinuityCounter()     // cc = 15
	generated.IncContinuityCounter()     // cc = 0, overflow

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nContinuity Counter did not rollover as expected.", generated, target)
	}
}

func TestIsPAT(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "47000010")
	isPAT, err := pkt.IsPAT()
	if err != nil {
		t.Error("failed to read if packet is a PAT")
		return
	}
	if !isPAT {
		t.Errorf("packet should be a PAT")
	}
}

func TestIsNull(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF10")
	isNull, err := pkt.IsNull()
	if err != nil {
		t.Error("failed to read if packet is null")
		return
	}
	if !isNull {
		t.Errorf("packet should be Null")
	}
}

func TestHasAdaptationField(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF100100")
	hasAF, err := pkt.HasAdaptationField()
	if err != nil {
		t.Errorf("could not run has adaptation field. error: %s", err.Error())
		return
	}
	if hasAF {
		t.Errorf("Packet should not have Adaptation Field (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF200100")
	hasAF, err = pkt.HasAdaptationField()
	if err != nil {
		t.Errorf("could not run has adaptation field. error: %s", err.Error())
		return
	}
	if !hasAF {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF300100")
	hasAF, err = pkt.HasAdaptationField()
	if err != nil {
		t.Errorf("could not run has adaptation field. error: %s", err.Error())
		return
	}
	if !hasAF {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 11).")
	}
}

func TestHasPayload(t *testing.T) {
	pkt := createPacketEmptyPayload(t, "471FFF10")
	hasPayload, err := pkt.HasPayload()
	if err != nil {
		t.Errorf("could not run has payload. error: %s", err.Error())
		return
	}
	if !hasPayload {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF20")
	hasPayload, err = pkt.HasPayload()
	if err != nil {
		t.Errorf("could not run has payload. error: %s", err.Error())
		return
	}
	if hasPayload {
		t.Errorf("Packet should not have Payload (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyPayload(t, "471FFF30")
	hasPayload, err = pkt.HasPayload()
	if err != nil {
		t.Errorf("could not run has payload. error: %s", err.Error())
		return
	}
	if !hasPayload {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 11).")
	}
}

func TestHeaderComboBasic(t *testing.T) {
	target := createPacketEmptyPayload(t, "47EFA098")
	generated := NewPacket()

	if err := generated.SetContinuityCounter(7); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.IncContinuityCounter(); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetPID(4000); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetTransportErrorIndicator(true); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetPayloadUnitStartIndicator(true); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetTransportPriority(true); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetAdaptationFieldControl(PayloadFlag); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}
	if err := generated.SetTransportScramblingControl(ScrambleEvenKeyFlag); err != nil {
		t.Errorf("failed to set field in packet. error: %s", err.Error())
		return
	}

	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nFields did not set successfully", generated, target)
	}

	cc, err := generated.ContinuityCounter()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	pid, err := generated.PID()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	tei, err := generated.TransportErrorIndicator()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	pusi, err := generated.PayloadUnitStartIndicator()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	tp, err := generated.TransportPriority()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	afc, err := generated.AdaptationFieldControl()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}
	tsc, err := generated.TransportScramblingControl()
	if err != nil {
		t.Errorf("failed to read field in packet. error: %s", err)
		return
	}

	if cc != 8 {
		t.Errorf("Reads different ContinuityCounter than expected.")
	}
	if pid != 4000 {
		t.Errorf("Reads different PID than expected.")
	}
	if tei != true {
		t.Errorf("Reads different TransportErrorIndicator than expected.")
	}
	if pusi != true {
		t.Errorf("Reads different PayloadUnitStartIndicator than expected.")
	}
	if tp != true {
		t.Errorf("Reads different TP than expected.")
	}
	if afc != PayloadFlag {
		t.Errorf("Reads different AdaptationFieldControl than expected.")
	}
	if tsc != ScrambleEvenKeyFlag {
		t.Errorf("Reads different TransportScramblingControl than expected.")
	}
}

func TestNilSlicePacket(t *testing.T) {
	//target := createPacketEmptyPayload(t, "47EFA098")
	generated := Packet(nil)

	if err := generated.SetContinuityCounter(7); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if err := generated.SetPID(4000); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if err := generated.SetTransportErrorIndicator(true); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if err := generated.SetAdaptationFieldControl(PayloadFlag); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if _, err := generated.ContinuityCounter(); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if _, err := generated.PID(); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if _, err := generated.TransportErrorIndicator(); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
	if _, err := generated.AdaptationFieldControl(); err != gots.ErrInvalidPacketLength {
		t.Errorf("incorrect error returned. error: %s", err.Error())
		return
	}
}

func BenchmarkNewStyleAllFields(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// create everything
		p := NewPacket()
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
		pkt := NewPacket()
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
				uint16(13),
				WithHasPayloadFlag,
				WithPUSI),
			7)
	}
}

func BenchmarkNewStyleRead(b *testing.B) {
	pkt := NewPacket()
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
	pkt := NewPacket()
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
