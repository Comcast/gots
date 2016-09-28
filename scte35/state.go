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

import "github.com/Comcast/gots"

type state struct {
	open        []SegmentationDescriptor
	blackoutIdx int
	inBlackout  bool
}

// NewState returns an initialized state object
func NewState() State {
	return &state{}
}

func (s *state) Open() []SegmentationDescriptor {
	open := make([]SegmentationDescriptor, len(s.open))
	copy(open, s.open)
	if s.inBlackout {
		return append(open[0:s.blackoutIdx], open[s.blackoutIdx+1:]...)
	} else {
		return open
	}
}

func (s *state) ProcessDescriptor(desc SegmentationDescriptor) ([]SegmentationDescriptor, error) {
	var err error
	var closed []SegmentationDescriptor
	// check if this is a duplicate
	for _, d := range s.open {
		if desc.Equal(d) {
			return nil, gots.ErrSCTE35DuplicateDescriptor
		}
	}
	// first close signals until one returns false, then handle the breakaway
	for i := len(s.open) - 1; i >= 0; i-- {
		d := s.open[i]
		if desc.CanClose(d) {
			closed = append(closed, d)
		} else {
			break
		}
	}
	// remove all closed descriptors
	s.open = s.open[0 : len(s.open)-len(closed)]
	if desc.IsOut() {
		s.open = append(s.open, desc)
		// if there was anything in the close list, that means we missed the out
		// for that descriptor
		if len(closed) != 0 {
			err = gots.ErrSCTE35MissingOut
		}
	} else if len(closed) != 0 {
		// check for validity - if a descriptor is closed, the closing descriptor
		// must be of the appropriate closing type
		switch desc.TypeID() {
		case SegDescProgramEnd,
			SegDescChapterEnd,
			SegDescProviderAdvertisementEnd, SegDescProviderPOEnd,
			SegDescDistributorAdvertisementEnd, SegDescDistributorPOEnd,
			SegDescUnscheduledEventEnd, SegDescNetworkEnd:
			if closed[len(closed)-1].TypeID() != desc.TypeID()-1 {
				err = gots.ErrSCTE35MissingOut
			}
		}
	}
	// breakaway handling
	if desc.TypeID() == SegDescProgramBreakaway {
		s.inBlackout = true
		s.blackoutIdx = len(s.open)
		s.open = append(s.open, desc)
	} else if desc.TypeID() == SegDescProgramResumption {
		if s.inBlackout {
			s.inBlackout = false
		} else {
			// ProgramResumption can only come after a breakaway
			err = gots.ErrSCTE35InvalidDescriptor
		}
		s.open = s.open[0:s.blackoutIdx]
	}
	return closed, err
}

func (s *state) Close(desc SegmentationDescriptor) ([]SegmentationDescriptor, error) {
	// back off list until we reach the descriptor we are closing. If we don't
	// find it, return error
	var closed []SegmentationDescriptor
	for i := len(s.open) - 1; i >= 0; i-- {
		d := s.open[i]
		closed = append(closed, d)
		if desc.Equal(d) {
			// found our descriptor, remove it and everything after it
			s.open = s.open[0:i]
			return closed, nil
		}
	}
	return nil, gots.ErrSCTE35DescriptorNotFound
}
