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
	"bufio"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSyncForSmoke(t *testing.T) {
	bs, _ := hex.DecodeString("474000100000b00d0001c100000001e256f803e71bfffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ff4742")
	r := bufio.NewReader(bytes.NewReader(bs))

	offset, err := Sync(r)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if offset != 0 {
		t.Errorf("Incorrect offset returned for next sync marker")
	}
}

func TestSyncNonZeroOffset(t *testing.T) {
	bs, _ := hex.DecodeString("ffffff474000100000b00d0001c100000001e256f803e71bfffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ff4742")
	r := bufio.NewReader(bytes.NewReader(bs))

	offset, err := Sync(r)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if offset != 3 {
		t.Errorf("Incorrect offset returned for next sync marker")
	}
}

func TestSyncNotFound(t *testing.T) {
	// no sync byte here
	bs, _ := hex.DecodeString("ff4000100000b00d0001c100000001e256f803e71bfffffff")
	r := bufio.NewReader(bytes.NewReader(bs))

	_, err := Sync(r)
	if err == nil {
		t.Errorf("Expected there to be an error, but there was not")
	}
}

func TestSyncReaderPosAtPacketStart(t *testing.T) {
	bs, _ := hex.DecodeString("ffffff474000100000b00d0001c100000001e256f803e71bfffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
		"ff4742")
	r := bufio.NewReader(bytes.NewReader(bs))

	_, err := Sync(r)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	bytes := make([]byte, 3)
	_, err = r.Read(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	wantHex := "474000"
	gotHex := hex.EncodeToString(bytes)
	if gotHex != wantHex {
		t.Errorf("Reader not left in correct spot. Wanted next read %v, got %v", wantHex, gotHex)
	}
}
