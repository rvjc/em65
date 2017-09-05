// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// load() loads binary data and source code.
func load(name string) {
	loadBin(name + ".bin")
	loadLst(name + ".lst")
	fmt.Println()
}

// loadBin() loads raw binary data into memory.
// (See package documentation for details)
func loadBin(path string) {

	fmt.Println("Found binary file:", path)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]uint8, memSize)
	count, err := file.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	flashing = true
	binStart = memMax - uint16(count) + 1
	addr := binStart
	fmt.Println("Loading binary data from " + fmtWord(addr) + " to " + fmtWord(memMax) + "...")

	for i := 0; i < count; i++ {
		data := buf[i]
		writeByte(addr, data)
		addr++
	}

	flashing = false
	fmt.Println("Binary file loaded\n")
}

// loadLst() loads source code from the LST file. 
// (See package documentation for details)
func loadLst(path string) {

	var comment, label string
	var sd srcData

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Println("Found LST file:", path)
	fileBuf := bufio.NewReader(file)
	lineCount := 0
	srcCount := 0
	errCount := 0

getLine:
	for {
		line, err := fileBuf.ReadString('\n')
		if err != nil {
			break getLine
		}
		lineCount++

		// extract trailing comment first
		span := strings.TrimSpace(line)
		i := strings.LastIndex(span, ";")
		j := strings.LastIndex(span, "\"")
		if i > 0 && i > j {
			comment = strings.TrimSpace(span[i+1:])
			span = strings.TrimSpace(span[:i])
		} else {
			comment = ""
		}

		// extract fields to left of comment
		fields := strings.Fields(span)
		fieldCount := len(fields)
		if fieldCount < 3 {
			continue getLine
		}
		if fields[1] != ":" {
			continue getLine
		}

		// extract address
		hexAddr := fields[0]
		if len(hexAddr) != 4 {
			continue getLine
		}
		parAddr, err := strconv.ParseUint(hexAddr, 16, 16)
		if err != nil {
			continue getLine
		}
		addr := uint16(parAddr)

		// retrieve existing source for this address (if any)
		// e.g. there may already have been an address label
		// on the previous line with zero bytes of actual code.
		// If there is no existing source for this address then
		// we are starting off with blank source data.

		sd = src[addr]

		// concatentate comments from multiple lines
		if len(comment) > 0 {
			if len(sd.comment) > 0 {
				sd.comment += " - "
			}
			sd.comment += comment
		}

		// check for a label only line of source  with zero bytes
		// ignore > symbol masquerading as a label
		if len(fields) == 3 {

			label = fields[2]
			if label != ">" {
				sd.label = label
			}

		} else {

			// extract ByteCount and verify code bytes
			hexCode := fields[2]
			charCount := len(hexCode)
			if charCount > 6 {
				continue getLine
			}
			if charCount%2 != 0 {
				continue getLine
			}
			sd.byteCount = charCount / 2
			for i = 0; i < sd.byteCount; i++ {
				j = 2 * i
				pbyte, err := strconv.ParseUint(hexCode[j:j+2], 16, 8)
				if err != nil {
					continue getLine
				}
				if uint8(pbyte) != readByte(addr+uint16(i)) {
					errCount += 1
					continue getLine
				}
			}

			// extract label, mnemomic and operand
			if mnems[fields[3]] {
				sd.mnem = fields[3]
				i = 4
			} else {
				if fieldCount < 5 {
					continue getLine
				}
				if !mnems[fields[4]] {
					continue getLine
				}
				label = fields[3]
				if label != ">" {
					sd.label = label
				}
				sd.mnem = fields[4]
				i = 5
			}
			sd.operand = strings.Join(fields[i:], "")
		}
		src[addr] = sd
		srcCount++
	}

	fmt.Println(lineCount, "lines processed")
	fmt.Println(srcCount, "lines of valid source obtained\n")
	if errCount > 0 {
		fmt.Println("WARNING:", errCount, "code errors found.\n"+
			"Probable cause: Binary file mis-alignment.\n"+
			"Last address of binary data is assumed to be $FFFF.\n"+
			"For correct binary alignment, source must end at top-of-memory.\n"+
			"Solution: Specify last vector (IRQ) at $FFFE.\n")
	}
}
