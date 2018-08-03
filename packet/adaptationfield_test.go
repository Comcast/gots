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
	"testing"
)

func generatePacketAF(t *testing.T, AFString string) (*Packet, AdaptationField) {
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
		t.Errorf("failed to read discontinuity correctly.")
	}
	_, a = generatePacketAF(t, "0190")
	if discontinuity, err := a.Discontinuity(); !discontinuity || err != nil {
		t.Errorf("failed to read discontinuity correctly.")
	}
	_, a = generatePacketAF(t, "0170")
	if discontinuity, err := a.Discontinuity(); discontinuity || err != nil {
		t.Errorf("failed to read discontinuity correctly.")
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
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the PCR to 1 has failed.", generated, target)
	}

	target, _ = generatePacketAF(t, "B718000000007E01000000007E02")
	a.SetHasOPCR(true)
	a.SetOPCR(2)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the OPCR to 2 has failed.", generated, target)
	}

	target, _ = generatePacketAF(t, "B71A000000007E01000000007E020188")
	a.SetHasTransportPrivateData(true)
	a.SetTransportPrivateData([]byte{0x88})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Private Data to 0x88 has failed.", generated, target)
	}

	target, _ = generatePacketAF(t, "B71B000000007E01000000007E0201880100")
	a.SetHasAdaptationFieldExtension(true)
	a.SetAdaptationFieldExtension([]byte{0x00})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Adaptation Field Extension to 0x00 has failed.", generated, target)
	}

	target, _ = generatePacketAF(t, "B71B000000007E01000000007E020266660100")
	a.SetTransportPrivateData([]byte{0x66, 0x66})
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nSetting the Transport Private Data to 0x6666 has failed.", generated, target)
	}

	target, _ = generatePacketAF(t, "B713000000007E010266660100")
	a.SetHasOPCR(false)
	if !Equal(generated, target) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X\nremoving has failed.", generated, target)
	}
}
