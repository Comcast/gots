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
		if err != nil {
			fmt.Println("Cannot close File", file.Name(), err)
		}
	}(tsFile)
	// Verify if sync-byte is present and seek to the first sync-byte
	syncIndex, err := sync(tsFile)
	if err == nil {
		_, err = tsFile.Seek(syncIndex, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(err)
		return
	}
	pat, err := psi.ReadPAT(tsFile)
	if err != nil {
		println(err)
		return
	}
	printPat(pat)

	if *showPmt {
		pm := pat.ProgramMap()
		for pn, pid := range pm {
			pmt, err := psi.ReadPMT(tsFile, pid)
			if err != nil {
				panic(err)
			}
			printPmt(pn, pmt)
		}
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

func printPmt(pn uint16, pmt psi.PMT) {
	printlnf("Program #%v PMT", pn)
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
	printlnf("\tPMT PIDs %v", pat.ProgramMap())
	printlnf("\tNumber of Programs %v", pat.NumPrograms())
}

func printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func sync(buf io.Reader) (int64, error) {
	// function find the first sync byte of the array
	data := make([]byte, 1)
	for i := int64(0); ; i++ {
		read, err := buf.Read(data)
		if err != nil && err != io.EOF {
			println(err)
		}
		if read == 0 {
			break
		}
		if int(data[0]) == packet.SyncByte {
			// check next 188th byte
			nextData := make([]byte, packet.PacketSize)
			nextRead, err := buf.Read(nextData)
			if err != nil && err != io.EOF {
				println(err)
			}
			if nextRead == 0 {
				break
			}
			if nextData[187] == packet.SyncByte {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("Sync-byte not found.")
}
