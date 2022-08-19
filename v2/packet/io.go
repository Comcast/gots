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
	"encoding/binary"
	"io"

	"github.com/Comcast/gots/v2"
)

// Peeker wraps the Peek method.
type Peeker interface {
	// Peek returns the next n bytes without advancing the reader.
	Peek(n int) ([]byte, error)
}

// PeekScanner is an extended io.ByteScanner with peek capacity.
type PeekScanner interface {
	io.ByteScanner
	Peeker
}

// Sync finds the offset of the next packet header and advances the reader
// to the packet start. It returns the offset of the sync relative to the
// original reader position.
//
// Sync uses IsSynced to determine whether a position is at packet header.
func Sync(r PeekScanner) (off int64, err error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return off, gots.ErrSyncByteNotFound
			}
			return off, err
		}
		if b != SyncByte {
			off++
			continue
		}

		err = r.UnreadByte()
		if err != nil {
			return off, err
		}
		ok, err := IsSynced(r)
		if ok {
			return off, nil
		}
		if err != nil {
			if err == io.EOF {
				return off, gots.ErrSyncByteNotFound
			}
			return off, err
		}

		// Advance again. This is a consequence of not
		// duplicating IsSynced for 3 and 4 byte reads.
		_, err = r.ReadByte()
		// These errors should never happen since we
		// have already read this byte above.
		if err != nil {
			if err == io.EOF {
				return off, gots.ErrSyncByteNotFound
			}
			return off, err
		}
	}
}

// IsSynced returns whether r is synced to a packet boundary.
//
// IsSynced checks whether the first byte is the MPEG-TS sync byte,
// the PID is not within the reserved range of 0x4-0xf and that
// the AFC is not the reserved value 0x0.
func IsSynced(r Peeker) (ok bool, err error) {
	b, err := r.Peek(4)
	if err != nil {
		return false, err
	}
	// Check that the first byte is the sync byte.
	if b[0] != SyncByte {
		return false, nil
	}

	const (
		pidMask = 0x1fff << 8
		afcMask = 0x3 << 4
	)
	header := binary.BigEndian.Uint32(b)

	// Check that the AFC is not zero (reserved).
	afc := header & afcMask
	if afc == 0x0 {
		return false, nil
	}

	// Check that the PID is not 0x4-0xf (reserved).
	pid := (header & pidMask) >> 8
	return pid < 0x4 || 0xf < pid, nil
}
