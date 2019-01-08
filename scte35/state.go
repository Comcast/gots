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
	"github.com/Comcast/gots"
	"strings"
)

const receivedRingLen = 10

type receivedElem struct {
	pts   gots.PTS
	descs []SegmentationDescriptor
}

type state struct {
	open         []SegmentationDescriptor
	received     []*receivedElem
	receivedHead int
	blackoutIdx  int
	inBlackout   bool
}

// NewState returns an initialized state object
func NewState() State {
	return &state{received: make([]*receivedElem, receivedRingLen)}
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
	// check if desc has a pts because we can't handle if it doesn't
	if !desc.SCTE35().HasPTS() {
		return nil, gots.ErrSCTE35UnsupportedSpliceCommand
	}
	// check if this is a duplicate - if not, add it to the received list and
	// drop the old received if we're over the length limit
	descAdded := false
	pts := desc.SCTE35().PTS()
	for _, e := range s.received {
		if e != nil {
			for _, d := range e.descs {
				if e.pts == pts {
					if desc.Equal(d) {
						// Duplicate desc found
						return nil, gots.ErrSCTE35DuplicateDescriptor
					}
					e.descs = append(e.descs, desc)
					descAdded = true
				}
				// check if we have seen a VSS signal with the same signalId and
				// same eventId before.
				if desc.EventID() == d.EventID() &&
					d.TypeID() == SegDescUnscheduledEventStart && desc.TypeID() == SegDescUnscheduledEventStart {
					descStreamSwitchSignalId, err := desc.StreamSwitchSignalId()
					if err != nil {
						return nil, err
					}

					dStreamSwitchSignalId, err := d.StreamSwitchSignalId()
					if err != nil {
						return nil, err
					}

					if strings.Compare(descStreamSwitchSignalId, dStreamSwitchSignalId) == 0 &&
						(d.EventID() == desc.EventID()) {
						// desc and d contain same signalId and same eventID
						// we should not be processing this desc.
						return nil, gots.ErrSCTE35DuplicateDescriptor
					}
					descAdded = true
				}
			}
		}
	}
	if !descAdded {
		s.received[s.receivedHead] = &receivedElem{pts: pts, descs: []SegmentationDescriptor{desc}}
		s.receivedHead = (s.receivedHead + 1) % receivedRingLen
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

	// validation logic
	switch desc.TypeID() {
	// breakaway handling
	case SegDescProgramBreakaway:
		s.inBlackout = true
		s.blackoutIdx = len(s.open)
		// append breakaway to match against resumption even though it's an in
		s.open = append(s.open, desc)
	case SegDescProgramResumption:
		if s.inBlackout {
			s.inBlackout = false
			s.open = s.open[0:s.blackoutIdx]
			// TODO: verify that there is a program start that has a matching event id
		} else {
			// ProgramResumption can only come after a breakaway
			err = gots.ErrSCTE35InvalidDescriptor
		}
		fallthrough
	// out signals
	case SegDescProgramStart,
		SegDescChapterStart,
		SegDescProviderAdvertisementStart,
		SegDescDistributorAdvertisementStart,
		SegDescProviderPOStart,
		SegDescDistributorPOStart,
		SegDescUnscheduledEventStart,
		SegDescNetworkStart,
		SegDescProgramOverlapStart,
		SegDescProgramStartInProgress:
		s.open = append(s.open, desc)

	// in signals
	// SegDescProgramEnd treated individually since it is expected to normally
	// close program resumption AND program start
	case SegDescProgramEnd:
		if len(closed) == 0 {
			err = gots.ErrSCTE35MissingOut
			break
		}
		for _, d := range closed {
			if d.TypeID() != SegDescProgramStart &&
				d.TypeID() != SegDescProgramResumption {
				err = gots.ErrSCTE35MissingOut
				break
			}
		}
	case SegDescChapterEnd,
		SegDescProviderAdvertisementEnd,
		SegDescProviderPOEnd,
		SegDescDistributorAdvertisementEnd,
		SegDescDistributorPOEnd,
		SegDescUnscheduledEventEnd,
		SegDescNetworkEnd:
		var openDesc SegmentationDescriptor
		// We already closed a descriptor
		// and have no other open descriptors
		// so break and return closed descriptors
		if len(closed) != 0 && len(s.open) == 0 {
			break
		}
		// descriptor matches out, but doesn't close it.  Check event id against open
		if len(closed) == 0 || closed[len(closed)-1].TypeID() != desc.TypeID()-1 {
			if len(s.open) == 0 {
				err = gots.ErrSCTE35MissingOut
				break
			} else {
				openDesc = s.open[len(s.open)-1]
			}
		} else {
			openDesc = closed[len(closed)-1]
		}
		if openDesc.EventID() != desc.EventID() {
			err = gots.ErrSCTE35MissingOut
		}
	default:
		// no validating
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
