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
	"bytes"
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
		printlnf("Cannot access test asset %s.", fileName)
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
		println(err)
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

	pkt := make(packet.Packet, packet.PacketSize)
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
		if _, err := io.ReadFull(reader, pkt); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			println(err)
			return
		}
		numPackets++
		if *dumpSCTE35 {
			currPID, err := packet.Pid(pkt)
			if err != nil {
				printlnf("Cannot get packet PID for %d", currPID)
				continue
			}
			if scte35PIDs[currPID] {
				pay, err := packet.Payload(pkt)
				if err != nil {
					printlnf("Cannot get payload for packet number %d on PID %d Error=%s", numPackets, currPID, err)
					continue
				}
				msg, err := scte35.NewSCTE35(pay)
				if err != nil {
					printlnf("Cannot parse SCTE35 Error=%v", err)
					continue
				}
				printSCTE35(currPID, msg)

			}

		}
		if *showEbp {
			ebpBytes, err := adaptationfield.EncoderBoundaryPoint(pkt)
			if err != nil {
				// Not an EBP
				continue
			}
			buf := bytes.NewBuffer(ebpBytes)
			boundaryPoint, err := ebp.ReadEncoderBoundaryPoint(buf)
			if err != nil {
				fmt.Printf("EBP construction error %v", err)
				continue
			}
			ebps[numPackets] = boundaryPoint
			printlnf("Packet %d contains EBP %+v", numPackets, boundaryPoint)
		}
		if *showPacketNumberOfPID != 0 {
			pid := uint16(*showPacketNumberOfPID)
			pktPid, err := packet.Pid(pkt)
			if err != nil {
				continue
			}
			if pktPid == pid {
				printlnf("First Packet of PID %d contents: %x", pid, pkt)
				break
			}
		}
	}
	println()

}

func printSCTE35(pid uint16, msg scte35.SCTE35) {
	printlnf("SCTE35 Message on PID %d", pid)

	printSpliceCommand(msg.CommandInfo())

	insert, ok := msg.CommandInfo().(scte35.SpliceInsertCommand)
	if ok {

		printSpliceInsertCommand(insert)
	}
	for _, segdesc := range msg.Descriptors() {
		printSegDesc(segdesc)
	}

}

func printSpliceCommand(spliceCommand scte35.SpliceCommand) {
	printlnf("\tCommand Type %v", scte35.SpliceCommandTypeNames[spliceCommand.CommandType()])
	if spliceCommand.HasPTS() {

		printlnf("\tPTS %v", spliceCommand.PTS())

	}
}

func printSegDesc(segdesc scte35.SegmentationDescriptor) {
	if segdesc.IsIn() {

		printlnf("\t<--- IN Segmentation Descriptor")
	}
	if segdesc.IsOut() {

		printlnf("\t---> OUT Segmentation Descriptor")
	}

	printlnf("\t\tEvent ID %d", segdesc.EventID())
	printlnf("\t\tType %+v", scte35.SegDescTypeNames[segdesc.TypeID()])
	if segdesc.HasDuration() {

		printlnf("\t\t Duration %v", segdesc.Duration())
	}

}

func printSpliceInsertCommand(insert scte35.SpliceInsertCommand) {
	println("\tSplice Insert Command")
	printlnf("\t\tEvent ID %v", insert.EventID())
	if insert.HasDuration() {
		printlnf("\t\tDuration %v", insert.Duration())

	}
}

func printPmt(pn uint16, pmt psi.PMT) {
	printlnf("Program #%v PMT", pn)
	printlnf("\tPIDs %v", pmt.Pids())
	println("\tElementary Streams")
	for _, es := range pmt.ElementaryStreams() {
		printlnf("\t\tPid %v: StreamType %v: %v", es.ElementaryPid(), es.StreamType(), es.StreamTypeDescription())
		for _, d := range es.Descriptors() {
			printlnf("\t\t\t%+v", d)
		}
	}
}

func printPat(pat psi.PAT) {
	println("Pat")
	printlnf("\tPMT PIDs %v", pat.ProgramMap())
	printlnf("\tNumber of Programs %v", pat.NumPrograms())
}

func printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
