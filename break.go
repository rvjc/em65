// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import "fmt"

// chkBreak() checks the system state against pre-defined
// break parameters. It is called at the start of each
// operation loop when in debug mode. If a break match is
// found, the emulator is paused by reverting to step mode.
// Current break functionality is very basic and must be
// set programmatically
func chkBreak() {
	// As there is currently no keyboard interrupt,
	// a regular break is performed every 100 virtual seconds
	if ck > brkCK {
		fmt.Println("\nBreak on time check\n")
		brkCK += 100000000
		stepping = true
	}
	// It is also useful to break on a single line endless loop.
	// This is the standard way for the error program to terminate
	// both in the case of failure and success
	if pc == brkPC {
		fmt.Println("\nBreak on endless loop\n")
		stepping = true
	}
	brkPC = pc
	// The following line can be modified to break on a custom
	// PC value (see assembler listing output). If not used,
	// set to 0xFFFF
	if pc == 0xFFFF {
		fmt.Println("\nBreak on PC\n")
		stepping = true
	}
}
