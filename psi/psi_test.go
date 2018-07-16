package psi

import (
	"bytes"
	"testing"
)

func TestDefaultPSIData(t *testing.T) {
	psi := NewPSI()
	generated := psi.Bytes()
	target := []byte{0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(target, generated) {
		t.Errorf("NewPSI does not produce expected Data. \nExpected: %X \n     Got: %X ", target, generated)
	}
}

func TestPSIFromBytes(t *testing.T) {
	target := []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0x18, 0x81, 0xFF}
	psi := PSIFromBytes(target)
	if psi.PointerField != 0x04 {
		t.Errorf("PSIFromBytes does not produce expected PointerField. \nExpected: %X \n     Got: %X ", 0x04, psi.PointerField)
	}
	if psi.TableID != 0x18 {
		t.Errorf("PSIFromBytes does not produce expected TableID. \nExpected: %X \n     Got: %X ", 0x18, psi.TableID)
	}
	if psi.SectionLength != 0x1FF {
		t.Errorf("PSIFromBytes does not produce expected TableID. \nExpected: %X \n     Got: %X ", 0x1FF, psi.SectionLength)
	}
	if psi.PrivateIndicator {
		t.Errorf("PSIFromBytes does not produce expected PrivateIndicator. \nExpected: %t \n     Got: %t ", true, psi.PrivateIndicator)
	}
	if !psi.SectionSyntaxIndicator {
		t.Errorf("PSIFromBytes does not produce expected PrivateIndicator. \nExpected: %t \n     Got: %t ", false, psi.SectionSyntaxIndicator)
	}
	generated := psi.Bytes()
	if !bytes.Equal(target, generated) {
		t.Errorf("Data does not produce same bytes. \nExpected: %X \n     Got: %X ", target, generated)
	}
}

func TestPSICreate(t *testing.T) {
	target := []byte{0x05, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x18, 0x40, 0xFF}
	psi := NewPSI()

	psi.PointerField = 5
	psi.TableID = 0x18
	psi.SectionLength = 0x0FF
	psi.PrivateIndicator = true
	psi.SectionSyntaxIndicator = false
	generated := psi.Bytes()

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated PSI does not produce expected bytes. \nExpected: %X \n     Got: %X ", target, generated)
	}
}
