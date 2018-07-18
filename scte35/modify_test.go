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

package scte35

import (
	"bytes"
	"encoding/base64"
	//"encoding/hex"
	"testing"
)

func TestCreate(t *testing.T) {
	//base64Bytes, _ := hex.DecodeString("00FC302F0000CFA9798200FFFFFF05620020027FEFFF58EDE344FE007B98A003350000000A0008435545490038323151C630E9") // set len of descriptors
	base64Bytes, _ := base64.StdEncoding.DecodeString("APwwNQAAAAAAAAD/8AEAACQCIkNVRUnAAAAAf78BEzU5MzkwMjY1NjUxNzc3OTIxNjMBAQHrr2Ob")
	target := base64Bytes
	scte, _ := NewSCTE35(target)
	generated := scte.Bytes()
	if !bytes.Equal(target, generated) {
		t.Errorf("\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestCreate2(t *testing.T) {
	target := vss
	scte, _ := NewSCTE35(target)
	generated := scte.Bytes()
	if !bytes.Equal(target, generated) {
		t.Errorf("\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}
