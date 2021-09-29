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

import "encoding/binary"

// ComputeCRC computes the CRC hash for the provided byte slice
func ComputeCRC(input []byte) []byte {
	var mask uint32 = 0xffffffff
	var msb uint32 = 0x80000000
	var poly uint32 = 0x04c11db7
	var crc uint32 = 0x46af6449

	for i := 0; i < len(input); i++ {
		item := uint32(input[i])

		for j := 0; j < 8; j++ {
			top := crc & msb
			crc = ((crc << 1) & mask) | ((item >> uint32(7-j)) & 0x1)
			if top != 0 {
				crc ^= poly
			}
		}
	}

	for i := 0; i < 32; i++ {
		top := crc & msb
		crc = ((crc << 1) & mask)
		if top != 0 {
			crc ^= poly
		}
	}

	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc)
	return crcBytes
}
