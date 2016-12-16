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

// Package ebp is used to detect and extract EBP information. EBP information is stored in a packet adaptation field in the private data section.
package ebp

import "time"

// EBP tags
const (
	ComcastEbpTag             = uint8(0xA9)
	CableLabsEbpTag           = uint8(0xDF)
	CableLabsFormatIdentifier = 0x45425030
)

// EncoderBoundaryPoint represents shared operations available on an all EBPs.
type EncoderBoundaryPoint interface {
	SegmentFlag() bool
	FragmentFlag() bool
	TimeFlag() bool
	EBPTime() time.Time            // time set in the EBP
	EBPSuccessReadTime() time.Time // time the EBP was successfully read
	SapFlag() bool
	Sap() byte
	ExtensionFlag() bool
	EBPType() byte
}
