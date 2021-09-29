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
	"io"

	"github.com/Comcast/gots"
)

// PacketWriter is subject to all rules governing implementations of io.Writer
//
// Additionally PacketWriter implementations must not modify or retain packet
// pointers, even temporarily.
type PacketWriter interface {
	WritePacket(p *Packet) (n int, err error)
}

// PacketWriteCloser is subject to all rules governing implementations of
// io.Writer
//
// Additionally PacketWriter implementations must not modify or retain packet
// array pointers, even temporarily.
type PacketWriteCloser interface {
	PacketWriter
	io.Closer
}

// Writer is the interface that groups PacketWriter and io.Writer methods
type Writer interface {
	PacketWriter
	io.Writer
}

// WriteCloser is the interface that groups PacketWriteCloser and io.Writer methods
type WriteCloser interface {
	PacketWriteCloser
	io.Writer
}

type packetWriter struct {
	PacketWriteCloser
	pkt Packet
}

func (pw *packetWriter) Write(p []byte) (n int, err error) {
	if len(p)%PacketSize != 0 {
		return 0, gots.ErrInvalidPacketLength
	}

	for i := 0; i < len(p); i += PacketSize {
		copy(pw.pkt[:], p[i:])
		if m, err := pw.WritePacket(&pw.pkt); err == nil {
			n += m
		} else {
			return n + m, err
		}
	}

	if len(p) > n {
		err = io.ErrShortWrite
	}
	return
}

func (pw *packetWriter) ReadFrom(r io.Reader) (n int64, err error) {
	buf := pw.pkt[:]
	for {
		nr, er := r.Read(buf)
		if nr == PacketSize {
			nw, ew := pw.WritePacket(&pw.pkt)
			if nw > 0 {
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		} else if nr > 0 && nr != PacketSize {
			err = gots.ErrInvalidPacketLength
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return n, err
}

// IOWriter returns a Writer with a default implementation of Write method
// wrapping the provided PacketWriter w.
func IOWriter(w PacketWriter) Writer {
	return &packetWriter{PacketWriteCloser: NopCloser(w)}
}

// IOWriteCloser returns a WriteCloser with a default implementation of Write
// method wrapping the provided PacketWriteCloser w.
func IOWriteCloser(w PacketWriteCloser) WriteCloser {
	return &packetWriter{PacketWriteCloser: w}
}

type nopCloser struct {
	PacketWriter
}

func (nopCloser) Close() error { return nil }

// NopCloser returns a PacketWriteCloser with a no-op Close method wrapping
// the provided PacketWriter r.
func NopCloser(r PacketWriter) PacketWriteCloser {
	return nopCloser{r}
}

// The PacketWriterFunc type is an adapter to allow the use of
// ordinary functions as PacketWriters. If f is a function
// with the appropriate signature, PacketWriterFunc(f) is a
// PacketWriter that calls f.
type PacketWriterFunc func(*Packet) (int, error)

// WritePacket calls f(p).
func (f PacketWriterFunc) WritePacket(p *Packet) (n int, err error) {
	return f(p)
}
