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

// package main contains CLI utilities for testing
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Comcast/gots/ebp"
	"github.com/Comcast/gots/packet"
	"github.com/Comcast/gots/packet/adaptationfield"
	"github.com/Comcast/gots/psi"
	"github.com/Comcast/gots/scte35"
)

// main parses a ts file that is provided with the -f flag
func main() {
	fileName := flag.String("f", "", "Required: Path to TS file to read")
	showPmt := flag.Bool("pmt", true, "Output PMT info")
	showEbp := flag.Bool("ebp", false, "Output EBP info. This is a lot of info")
	dumpSCTE35 := flag.Bool("scte35", false, "Output SCTE35 signals and info.")
	showPacketNumberOfPID := flag.Int("pid", 0, "Dump the contents of the first packet encountered on PID to stdout")
	flag.Parse()
	if *fileName == "" {
		flag.Usage()
		return
	}
	tsFile, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Cannot access test asset %s.\n", *fileName)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Cannot close File", file.Name(), err)
		}
	}(tsFile)
	// Verify if sync-byte is present and seek to the first sync-byte
	reader := bufio.NewReader(tsFile)
	_, err = packet.Sync(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	pat, err := psi.ReadPAT(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	printPat(pat)

	var pmts []psi.PMT
	pm := pat.ProgramMap()
	for pn, pid := range pm {
		pmt, err := psi.ReadPMT(reader, pid)
		if err != nil {
			panic(err)
		}
		pmts = append(pmts, pmt)
		if *showPmt {
			printPmt(pn, pmt)
		}
	}

	var pkt packet.Packet
	var numPackets uint64
	ebps := make(map[uint64]ebp.EncoderBoundaryPoint)
	scte35PIDs := make(map[uint16]bool)
	if *dumpSCTE35 {
		for _, pmt := range pmts {
			for _, es := range pmt.ElementaryStreams() {
				if es.StreamType() == psi.PmtStreamTypeScte35 {
					scte35PIDs[es.ElementaryPid()] = true
					break
				}

			}
		}
	}

	for {
		if _, err := io.ReadFull(reader, pkt[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Println(err)
			return
		}
		numPackets++
		if *dumpSCTE35 {
			currPID := packet.Pid(&pkt)
			if scte35PIDs[currPID] {
				pay, err := packet.Payload(&pkt)
				if err != nil {
					fmt.Printf("Cannot get payload for packet number %d on PID %d Error=%s\n", numPackets, currPID, err)
					continue
				}
				msg, err := scte35.NewSCTE35(pay)
				if err != nil {
					fmt.Printf("Cannot parse SCTE35 Error=%v\n", err)
					continue
				}
				printSCTE35(currPID, msg)

			}

		}
		if *showEbp {
			ebpBytes, err := adaptationfield.EncoderBoundaryPoint(&pkt)
			if err != nil {
				// Not an EBP
				continue
			}
			boundaryPoint, err := ebp.ReadEncoderBoundaryPoint(ebpBytes)
			if err != nil {
				fmt.Printf("EBP construction error %v", err)
				continue
			}
			ebps[numPackets] = boundaryPoint
			fmt.Printf("Packet %d contains EBP %+v\n\n", numPackets, boundaryPoint)
		}
		if *showPacketNumberOfPID != 0 {
			pid := uint16(*showPacketNumberOfPID)
			pktPid := packet.Pid(&pkt)
			if pktPid == pid {
				fmt.Printf("First Packet of PID %d contents: %x\n", pid, pkt)
				break
			}
		}
	}
	fmt.Println()
}

func printSCTE35(pid uint16, msg scte35.SCTE35) {
	fmt.Printf("SCTE35 Message on PID %d\n", pid)
	printSpliceCommand(msg.CommandInfo())

	if insert, ok := msg.CommandInfo().(scte35.SpliceInsertCommand); ok {
		printSpliceInsertCommand(insert)
	}
	for _, segdesc := range msg.Descriptors() {
		printSegDesc(segdesc)
	}

}

func printSpliceCommand(spliceCommand scte35.SpliceCommand) {
	fmt.Printf("\tCommand Type %v\n", scte35.SpliceCommandTypeNames[spliceCommand.CommandType()])

	if spliceCommand.HasPTS() {
		fmt.Printf("\tPTS %v\n", spliceCommand.PTS())
	}
}

func printSegDesc(segdesc scte35.SegmentationDescriptor) {
	if segdesc.IsIn() {
		fmt.Printf("\t<--- IN Segmentation Descriptor\n")
	}
	if segdesc.IsOut() {
		fmt.Printf("\t---> OUT Segmentation Descriptor\n")
	}

	fmt.Printf("\t\tEvent ID %d\n", segdesc.EventID())
	fmt.Printf("\t\tType %+v\n", scte35.SegDescTypeNames[segdesc.TypeID()])
	if segdesc.HasDuration() {
		fmt.Printf("\t\t Duration %v\n", segdesc.Duration())
	}

}

func printSpliceInsertCommand(insert scte35.SpliceInsertCommand) {
	fmt.Println("\tSplice Insert Command")
	fmt.Printf("\t\tEvent ID %v\n", insert.EventID())

	if insert.HasDuration() {
		fmt.Printf("\t\tDuration %v\n", insert.Duration())
	}
}

func printPmt(pn uint16, pmt psi.PMT) {
	fmt.Printf("Program #%v PMT\n", pn)
	fmt.Printf("\tPIDs %v\n", pmt.Pids())
	fmt.Println("\tElementary Streams")

	for _, es := range pmt.ElementaryStreams() {
		fmt.Printf("\t\tPid %v: StreamType %v: %v\n", es.ElementaryPid(), es.StreamType(), es.StreamTypeDescription())

		for _, d := range es.Descriptors() {
			fmt.Printf("\t\t\t%+v\n", d)
		}
	}
}

func printPat(pat psi.PAT) {
	fmt.Println("Pat")
	fmt.Printf("\tPMT PIDs %v\n", pat.ProgramMap())
	fmt.Printf("\tNumber of Programs %v\n", pat.NumPrograms())
}
