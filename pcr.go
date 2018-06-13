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

package gots

// ExtractPCR extracts a PCR time
// PCR is the Program Clock Reference.
// First 33 bits are PCR base.
// Next 6 bits are reserved.
// Final 9 bits are PCR extension.
func ExtractPCR(bytes []byte) uint64 {
	var a, b, c, d, e, f uint64
	a = uint64(bytes[0])
	b = uint64(bytes[1])
	c = uint64(bytes[2])
	d = uint64(bytes[3])
	e = uint64(bytes[4])
	f = uint64(bytes[5])
	pcrBase := (a << 25) | (b << 17) | (c << 9) | (d << 1) | (e >> 7)
	pcrExt := ((e & 0x1) << 8) | f
	return pcrBase*300 + pcrExt
}

// InsertPCR insterts a given pcr time into a byte slice.
func InsertPCR(b []byte, pcr uint64) {
	pcrBase := pcr / 300
	pcrExt := (pcr - pcrBase*300) & 0x1ff
	b[0] = byte(pcrBase >> 25)
	b[1] = byte(pcrBase >> 17)
	b[2] = byte(pcrBase >> 9)
	b[3] = byte(pcrBase >> 1)
	b[4] = byte((pcrBase << 7) | (pcrExt >> 8) | 0x7e)
	b[5] = byte(pcrExt & 0xff)
}
