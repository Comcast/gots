package psi

import (
	"bytes"
	"testing"
)

func TestDefaultPSIData(t *testing.T) {
	th := NewTableHeader()
	generated := th.Bytes()
	target := []byte{0x00, 0x30, 0x00}
	if !bytes.Equal(target, generated) {
		t.Errorf("NewPSI does not produce expected Data. \nExpected: %X \n     Got: %X ", target, generated)
	}
}

func TestPSIFromBytes(t *testing.T) {
	target := []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0x18, 0xB1, 0xFF}
	th := TableHeaderFromBytes(target[5:])

	if th.TableID != 0x18 {
		t.Errorf("PSIFromBytes does not produce expected TableID. \nExpected: %X \n     Got: %X ", 0x18, th.TableID)
	}
	if th.SectionLength != 0x1FF {
		t.Errorf("PSIFromBytes does not produce expected TableID. \nExpected: %X \n     Got: %X ", 0x1FF, th.SectionLength)
	}
	if th.PrivateIndicator {
		t.Errorf("PSIFromBytes does not produce expected PrivateIndicator. \nExpected: %t \n     Got: %t ", true, th.PrivateIndicator)
	}
	if !th.SectionSyntaxIndicator {
		t.Errorf("PSIFromBytes does not produce expected PrivateIndicator. \nExpected: %t \n     Got: %t ", false, th.SectionSyntaxIndicator)
	}
	generated := append(NewPointerField(4), th.Bytes()...)
	if !bytes.Equal(target, generated) {
		t.Errorf("Data does not produce same bytes. \nExpected: %X \n     Got: %X ", target, generated)
	}
}

func TestPSICreate(t *testing.T) {
	target := []byte{0x05, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x18, 0x70, 0xFF}
	th := NewTableHeader()

	th.TableID = 0x18
	th.SectionLength = 0x0FF
	th.PrivateIndicator = true
	th.SectionSyntaxIndicator = false
	generated := append(NewPointerField(5), th.Bytes()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated PSI does not produce expected bytes. \nExpected: %X \n     Got: %X ", target, generated)
	}
}
