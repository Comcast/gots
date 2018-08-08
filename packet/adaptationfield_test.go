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
	"testing"
)

func generatePacketAF(t *testing.T, AFString string) (*Packet, *AdaptationField) {
	p := createPacketEmptyAdaptationField(t, "47000030"+AFString)
	a, err := p.AdaptationField()
	if err != nil {
		t.Errorf("failed to get adaptation field. error: %s", err.Error())
	}
	if a == nil {
		t.Errorf("adaptation field does not exist")
	}
	return &p, a
}

func TestDiscontinuity(t *testing.T) {
	_, a := generatePacketAF(t, "0180")
	if discontinuity, err := a.Discontinuity(); !discontinuity || err != nil {
		t.Errorf("failed to read discontinuity correctly. expected false got true.")
	}
	_, a = generatePacketAF(t, "0190")
	if discontinuity, err := a.Discontinuity(); !discontinuity || err != nil {
		t.Errorf("failed to read discontinuity correctly. expected true got false.")
	}
	_, a = generatePacketAF(t, "0170")
	if discontinuity, err := a.Discontinuity(); discontinuity || err != nil {
		t.Errorf("failed to read discontinuity correctly. expected false got true.")
	}
}

func TestSetDiscontinuity(t *testing.T) {
	target, _ := generatePacketAF(t, "0180")
	generated, a := generatePacketAF(t, "0100")
	a.SetDiscontinuity(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Discontinuity to true has failed.", *generated, *target)
	}
	target, _ = generatePacketAF(t, "0100")
	generated, a = generatePacketAF(t, "0180")
	a.SetDiscontinuity(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Discontinuity to false has failed.", *generated, *target)
	}
}

func TestRandomAccess(t *testing.T) {
	_, a := generatePacketAF(t, "0140")
	if randomAccess, err := a.RandomAccess(); !randomAccess || err != nil {
		t.Errorf("failed to read RandomAccess correctly. expected true got false.")
	}
	_, a = generatePacketAF(t, "0130")
	if randomAccess, err := a.RandomAccess(); randomAccess || err != nil {
		t.Errorf("failed to read RandomAccess correctly. expected false got true.")
	}
}

func TestSetRandomAccess(t *testing.T) {
	target, _ := generatePacketAF(t, "0140")
	generated, a := generatePacketAF(t, "0100")
	a.SetRandomAccess(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the RandomAccess to true has failed.", *generated, *target)
	}
	target, _ = generatePacketAF(t, "0100")
	generated, a = generatePacketAF(t, "0140")
	a.SetRandomAccess(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the RandomAccess to false has failed.", *generated, *target)
	}
}

func TestElementaryStreamPriority(t *testing.T) {
	_, a := generatePacketAF(t, "0120")
	if esp, err := a.ElementaryStreamPriority(); !esp || err != nil {
		t.Errorf("failed to read ElementaryStreamPriority correctly. expected true got false.")
	}
	_, a = generatePacketAF(t, "0110")
	if esp, err := a.ElementaryStreamPriority(); esp || err != nil {
		t.Errorf("failed to read ElementaryStreamPriority correctly. expected false got true.")
	}
}

func TestSetElementaryStreamPriority(t *testing.T) {
	target, _ := generatePacketAF(t, "0120")
	generated, a := generatePacketAF(t, "0100")
	a.SetElementaryStreamPriority(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the ElementaryStreamPriority to true has failed.", *generated, *target)
	}
	target, _ = generatePacketAF(t, "0100")
	generated, a = generatePacketAF(t, "0120")
	a.SetElementaryStreamPriority(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the ElementaryStreamPriority to false has failed.", *generated, *target)
	}
}

func TestHasSplicingPoint(t *testing.T) {
	_, a := generatePacketAF(t, "0104")
	if hsp, err := a.HasSplicingPoint(); !hsp || err != nil {
		t.Errorf("failed to read HasSplicingPoint correctly. expected true got false.")
	}
	_, a = generatePacketAF(t, "0111")
	if hsp, err := a.HasSplicingPoint(); hsp || err != nil {
		t.Errorf("failed to read HasSplicingPoint correctly. expected false got true.")
	}
}

func TestSetHasSplicingPoint(t *testing.T) {
	target, _ := generatePacketAF(t, "0F04")
	generated, a := generatePacketAF(t, "0100")
	if a.SetHasSplicingPoint(true) == nil {
		t.Error("adaptation field cannot fit a splice countdown field but no error was returned")
	}
	generated, a = generatePacketAF(t, "0F00")
	a.SetHasSplicingPoint(true)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the HasSplicingPoint to true has failed.", *generated, *target)
	}
	target, _ = generatePacketAF(t, "0100")
	generated, a = generatePacketAF(t, "0104")
	a.SetHasSplicingPoint(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the HasSplicingPoint to false has failed.", *generated, *target)
	}
}

func TestAdaptationField(t *testing.T) {
	p := createPacketEmptyPayload(t, "470000300102")
	a, err := p.AdaptationField()
	if err != nil {
		t.Errorf("error getting adaptation field: %s", err.Error())
	}
	if a == nil {
		t.Errorf("no adaptation field was returned")
	}

	p = createPacketEmptyPayload(t, "470000100002")
	a, err = p.AdaptationField()
	if err == nil {
		t.Error("no error was returned in trying to access a nonexistent adaptation field.")
	}
	if a != nil {
		t.Errorf("adaptation field does not exist but something was returned.")
	}
}

func TestAdaptationFieldFull(t *testing.T) {
	generated, a := generatePacketAF(t, "B700")
	target, _ := generatePacketAF(t, "B710000000007E01")
	err := a.SetHasPCR(true)
	if err != nil {
		t.Errorf("failed to set pcr flag. Error: %s", err.Error())
	}
	a.SetPCR(1)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PCR to 1 has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B718000000007E01000000007E02")
	a.SetHasOPCR(true)
	a.SetOPCR(2)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the OPCR to 2 has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B71A000000007E01000000007E020188")
	a.SetHasTransportPrivateData(true)
	a.SetTransportPrivateData([]byte{0x88})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Private Data to 0x88 has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B71B000000007E01000000007E0201880100")
	a.SetHasAdaptationFieldExtension(true)
	a.SetAdaptationFieldExtension([]byte{0x00})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Extension to 0x00 has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B71B000000007E01000000007E020266660100")
	a.SetTransportPrivateData([]byte{0x66, 0x66})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Private Data to 0x6666 has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B713000000007E010266660100")
	a.SetHasOPCR(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nremoving has failed.", *generated, *target)
	}

	target, _ = generatePacketAF(t, "B717000000007E01510266660100")
	a.SetHasSplicingPoint(true)
	a.SetSpliceCountdown(0x51)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Private Data to 0x6666 has failed.", *generated, *target)
	}

	a.SetHasOPCR(true)
	a.SetOPCR(2)

	hasPCR, err := a.HasPCR()
	if err != nil {
		t.Errorf("failed to get HasPCR. error: %s", err.Error())
	}
	if hasPCR != true {
		t.Errorf("generated HasPCR (%t) does not match expected HasPCR (%t)", hasPCR, true)
	}
	pcr, err := a.PCR()
	if err != nil {
		t.Errorf("failed to get PCR. error: %s", err.Error())
	}
	if pcr != 1 {
		t.Errorf("generated PCR (%d) does not match expected PCR (%d)", pcr, 1)
	}

	hasOPCR, err := a.HasOPCR()
	if err != nil {
		t.Errorf("failed to get HasPCR. error: %s", err.Error())
	}
	if hasOPCR != true {
		t.Errorf("generated HasPCR (%t) does not match expected HasPCR (%t)", hasOPCR, true)
	}
	opcr, err := a.OPCR()
	if err != nil {
		t.Errorf("failed to get OPCR. error: %s", err.Error())
	}
	if opcr != 2 {
		t.Errorf("generated OPCR (%d) does not match expected OPCR (%d)", opcr, 2)
	}

	hasSplicingPoint, err := a.HasSplicingPoint()
	if err != nil {
		t.Errorf("failed to get HasSplicingPoint. error: %s", err.Error())
	}
	if hasSplicingPoint != true {
		t.Errorf("generated hasSplicingPoint (%t) does not match expected hasSplicingPoint (%t)", hasSplicingPoint, true)
	}
	spliceCountdown, err := a.SpliceCountdown()
	if err != nil {
		t.Errorf("failed to get spliceCountdown. error: %s", err.Error())
	}
	if spliceCountdown != 0x51 {
		t.Errorf("generated spliceCountdown (0x%X) does not match expected spliceCountdown (0x%X)", spliceCountdown, 0x51)
	}

	hasTPD, err := a.HasTransportPrivateData()
	if err != nil {
		t.Errorf("failed to get hasTPD. error: %s", err.Error())
	}
	if hasTPD != true {
		t.Errorf("generated HasTransportPrivateData (%t) does not match expected HasTransportPrivateData (%t)", hasTPD, true)
	}
	tpd, err := a.TransportPrivateData()
	if err != nil {
		t.Errorf("failed to get TransportPrivateData. error: %s", err.Error())
	}
	if bytes.Equal(tpd, []byte{0x66, 0x66}) {
		t.Errorf("generated TransportPrivateData (0x%X) does not match expected TransportPrivateData (0x%X)", tpd, []byte{0x66, 0x66})
	}

	hasAFE, err := a.HasAdaptationFieldExtension()
	if err != nil {
		t.Errorf("failed to get HasAdaptationFieldExtension. error: %s", err.Error())
	}
	if hasAFE != true {
		t.Errorf("generated HasAdaptationFieldExtension (%t) does not match expected HasAdaptationFieldExtension (%t)", hasAFE, true)
	}
	afe, err := a.AdaptationFieldExtension()
	if err != nil {
		t.Errorf("failed to get AdaptationFieldExtension. error: %s", err.Error())
	}
	if bytes.Equal(afe, []byte{0x00}) {
		t.Errorf("generated AdaptationFieldExtension (0x%X) does not match expected AdaptationFieldExtension (0x%X)", tpd, []byte{0x00})
	}
}
