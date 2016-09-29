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
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Comcast/gots/ebp"
	"github.com/Comcast/gots/packet"
	"github.com/Comcast/gots/packet/adaptationfield"
	"github.com/Comcast/gots/psi"
)

// main parses a ts file that is provided with the -f flag
func main() {
	fileName := flag.String("f", "", "Required: Path to TS file to read")
	showPmt := flag.Bool("pmt", true, "Output PMT info")
	showEbp := flag.Bool("ebp", false, "Output EBP info. This is a lot of info")
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
		fmt.Println("Cannot close File", file.Name(), err)
	}(tsFile)
	pat, err := extractPat(tsFile)
	if err != nil {
		println(err)
		return
	}
	printPat(pat)
	if *showPmt {
		pmt, err := extractPmt(tsFile, pat.ProgramMapPid())
		if err != nil {
			println(err)
		}
		printPmt(pmt)
	}
	pkt := make([]byte, packet.PacketSize, packet.PacketSize)
	var offset int64
	var numPackets uint64
	ebps := make(map[uint64]ebp.EncoderBoundaryPoint)
	for {
		read, err := tsFile.ReadAt(pkt, offset)
		if err == io.EOF || read == 0 {
			break
		}
		offset += packet.PacketSize
		numPackets++
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

func printPmt(pmt psi.PMT) {
	println("PMT")
	printlnf("\tPIDs %v", pmt.Pids())
	println("\tElementary Streams")
	for _, es := range pmt.ElementaryStreams() {
		printlnf("\t\tPid %v : StreamType %v", es.ElementaryPid(), es.StreamType())
		for _, d := range es.Descriptors() {
			printlnf("\t\t\t%+v", d)
		}
	}
}

func printPat(pat psi.PAT) {
	println("Pat")
	printlnf("\tPMT PID %12v", pat.ProgramMapPid())
	printlnf("\tProgramNumber %6v", pat.ProgramNumber())
	printlnf("\tNumber of Programs %v", pat.ProgramNumber())
}

func extractPat(buf io.Reader) (psi.PAT, error) {
	pkt := make([]byte, packet.PacketSize)
	var pat psi.PAT
	for read, err := buf.Read(pkt); pat == nil; read, err = buf.Read(pkt) {
		if err != nil {
			return nil, err
		}
		if read <= 0 {
			return nil, fmt.Errorf("Reached EOF without PAT")
		}
		pid, err := packet.Pid(pkt)
		if err != nil {
			return nil, err
		}
		if pid == 0 {
			pay, err := packet.Payload(pkt)
			if err != nil {
				println(err)
				continue
			}
			cp := make([]byte, len(pay))
			copy(cp, pay)
			pat, err := psi.NewPAT(cp)
			if err != nil {
				println(err)
				continue
			}
			return pat, nil
		}
	}
	return nil, fmt.Errorf("No pat found")
}
func printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
func extractPmt(buf io.Reader, pid uint16) (psi.PMT, error) {
	pkt := make([]byte, packet.PacketSize)
	pmtAcc := packet.NewAccumulator(psi.PmtAccumulatorDoneFunc)
	var pmt psi.PMT
	for read, err := buf.Read(pkt); pmt == nil && read > 0; read, err = buf.Read(pkt) {
		if err != nil {
			return nil, err
		}
		currPid, err := packet.Pid(pkt)
		if err != nil {
			return nil, err
		}
		if currPid == pid {
			done, err := pmtAcc.Add(pkt)
			if err != nil {
				return nil, err
			}
			if done {
				b, err := pmtAcc.Parse()
				if err != nil {
					return nil, err
				}
				pmt, err = psi.NewPMT(b)
				if err != nil {
					return nil, err
				}

			}

		}
	}
	return pmt, nil
}
