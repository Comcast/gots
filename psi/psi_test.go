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
	"testing"
)

func TestDefaultPSIData(t *testing.T) {
	th := NewTableHeader()
	generated := th.Data()
	target := []byte{0x00, 0x30, 0x00}
	if !bytes.Equal(target, generated) {
		t.Errorf("NewTableHeader does not produce expected Data. \nExpected: %X \n     Got: %X ", target, generated)
	}
}

func TestPSIFromBytes(t *testing.T) {
	target := []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0x18, 0xB1, 0xFF}
	th := TableHeaderFromBytes(target[5:])

	if th.TableID != 0x18 {
		t.Errorf("TableHeaderFromBytes does not produce expected TableID. \nExpected: %X \n     Got: %X ", 0x18, th.TableID)
	}
	if th.SectionLength != 0x1FF {
		t.Errorf("TableHeaderFromBytes does not produce expected SectionLength. \nExpected: %X \n     Got: %X ", 0x1FF, th.SectionLength)
	}
	if th.PrivateIndicator {
		t.Errorf("TableHeaderFromBytes does not produce expected PrivateIndicator. \nExpected: %t \n     Got: %t ", true, th.PrivateIndicator)
	}
	if !th.SectionSyntaxIndicator {
		t.Errorf("TableHeaderFromBytes does not produce expected SectionSyntaxIndicator. \nExpected: %t \n     Got: %t ", false, th.SectionSyntaxIndicator)
	}
	generated := append(NewPointerField(4), th.Data()...)
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
	generated := append(NewPointerField(5), th.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Pointer field and TableHeader do not produce expected bytes. \nExpected: %X \n     Got: %X ", target, generated)
	}
}
