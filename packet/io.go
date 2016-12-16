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
	"io"

	"github.com/Comcast/gots"
)

// Sync finds the offset of the next packet sync byte and advances the reader
// to the packet start. It also checks the next 188th byte to ensure a sync is
// found. It returns the offset of the sync w.r.t. the original reader
// position.
func Sync(r *bufio.Reader) (int64, error) {
	data := make([]byte, 1)
	for i := int64(0); ; i++ {
		_, err := io.ReadFull(r, data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return 0, err
		}
		if int(data[0]) == SyncByte {
			// check next 188th byte
			rp := bufio.NewReaderSize(r, PacketSize) // extends only if needed
			nextData, err := rp.Peek(PacketSize)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			if err != nil {
				return 0, err
			}
			if nextData[187] == SyncByte {
				r.UnreadByte()
				return i, nil
			}
		}
	}
	return 0, gots.ErrSyncByteNotFound
}
