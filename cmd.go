// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"fmt"
	"strings"
)

func execCmd(cmd string) (pa postAction) {
	cmd = strings.TrimSpace(cmd)
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "":
		pa = cmdStep()
	case "g":
		pa = cmdGo()
	case "l":
		pa = cmdList()
	case "m":
		pa = cmdMem()
	case "r":
		pa = cmdReset()
	case "s":
		pa = cmdStack()
	case "z":
		pa = cmdZero()
	case "q":
		pa = cmdQuit()
	default:
		pa = cmdErr()
	}
	return
}

func cmdStep() (pa postAction) {
	pa = postActionContinue
	return
}

func cmdGo() (pa postAction) {
	stepping = false
	fmt.Println("\nRunning (press Ctrl-C to interrupt)...")
	pa = postActionContinue
	return
}

func cmdQuit() (pa postAction) {
	fmt.Println("\nQuitting...")
	pa = postActionQuit
	return
}

func cmdList() (pa postAction) {
	fmt.Println("\nListing...\n")
	addr := binStart
	for {
		cmd := ""
		fmt.Print(fmtWord(addr) + " " + fmtSrc(addr) + " >")
		fmt.Scanln(&cmd)
		if cmd == "x" {
			break
		}
		byteCount := src[addr].byteCount
		// For non-source lines, increment address by one
		if byteCount == 0 {
			byteCount = 1
		}
		nextAddr := addr + uint16(byteCount)
		// overflow causes wraparound
		if nextAddr < addr {
			break
		}
		addr = nextAddr
	}
	fmt.Println("\nEnd of Listing\n")
	pa = postActionHold
	return
}

func cmdMem() (pa postAction) {
	fmt.Println("\nMemory Dump...")
	dumpMem(0x200, 0x2ff)
	fmt.Println("End of Memory Dump\n")
	pa = postActionHold
	return
}

func cmdReset() (pa postAction) {
	fmt.Println("\nResetting...")
	reset()
	fmt.Println("\nReady\n")
	pa = postActionRefetch
	return
}

func cmdStack() (pa postAction) {
	fmt.Println("\nStack Dump...")
	dumpMem(saMin, saMax)
	fmt.Println("End of Stack Dump\n")
	pa = postActionHold
	return
}

func cmdZero() (pa postAction) {
	fmt.Println("\nZero Page Dump...")
	dumpMem(0x00, 0xFF)
	fmt.Println("End of Zero Page Dump\n")
	pa = postActionHold
	return
}

func cmdErr() (pa postAction) {
	fmt.Println("\n*** UNKNOWN COMMAND ***\n")
	pa = postActionHold
	return
}

func dumpMem(start uint16, end uint16) {
	for addr := start; addr <= end; addr++ {
		data := readByte(addr)
		if addr&0x000F == 0 {
			fmt.Println()
			fmt.Print(fmtWord(addr), ": ")
		}
		fmt.Print(fmtByte(data) + " ")
	}
	fmt.Println("\n")
}
