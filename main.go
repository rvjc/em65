// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"fmt"
)

// main() starts up emulator.
func main() {

	fmt.Println("\nEmulator Initialising\n")

	initAll()
	load("test")
	reset()
	active = true
	opLoop()
	active = false

	fmt.Println("\nEmulator Terminated\n")
}

// opLoop() reads and executes each CPU instruction in turn
// stopping every so often to synchronise.
func opLoop() {

getOp:
	for {
		if ck >= syncNextCk {
			sync()
		}
		op = readByte(pc)
		if debugging {
			chkBreak()
			if stepping {
			getCmd:
				for {
					cmd := ""
					fmt.Print(fmtState() + " >")
					fmt.Scanln(&cmd)
					pa := execCmd(cmd)
					switch pa {
					case postActionHold:
						continue getCmd
					case postActionContinue:
						break getCmd
					case postActionRefetch:
						continue getOp
					case postActionQuit:
						break getOp
					}
				}
			}
		}
		opFunc := opFuncs[op]
		if opFunc == nil {
			fmt.Println("ILLEGAL INSTRUCTION at", fmtPc(), ":", fmtOp())
			break getOp
		}
		opFunc()
	}
}
