package scte35

import (
	"time"

	"github.com/Comcast/gots"
)

// Event is use to create a xml representation of the scte35 event
type Event struct {
	EventID       uint32            `xml:"eventID,attr"`
	EventTime     time.Time         `xml:"eventTime,attr"`
	PTS           gots.PTS          `xml:"pts,attr"`
	Command       SpliceCommandType `xml:"command,attr"`
	TypeID        SegDescType       `xml:"typeID,attr"`
	UPID          []byte            `xml:"upid,attr"`
	BreakDuration gots.PTS          `xml:"duration,attr"`
}
