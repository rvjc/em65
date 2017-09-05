// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"os"
	"os/signal"
)

// initAll() performs general program initialisation.
func initAll() {
	initSys()
	initDebug()
	initOps()
	initRom()
	initRam()
}

// initSys() performs system initialisation
func initSys() {
	active = false
	osSigs = make(chan os.Signal, 1)
	signal.Notify(osSigs, os.Interrupt, os.Kill)
	go osSigHandler()
}

// initDebug() initialises debug parameters.
func initDebug() {
	debugging = true
	stepping = true
	brkCK = 100000000
	brkPC = 0xFFFF
}

// initOps() initialises opFuncs table.
// Unused or unimplemented opcodes are mapped to nil.
func initOps() {

	opFuncs[0x69] = adcImm
	opFuncs[0x65] = adcZpg
	opFuncs[0x75] = adcZpx
	opFuncs[0x6D] = adcAbs
	opFuncs[0x7D] = adcAbx
	opFuncs[0x79] = adcAby
	opFuncs[0x61] = adcIdx
	opFuncs[0x71] = adcIdy

	opFuncs[0x29] = andImm
	opFuncs[0x25] = andZpg
	opFuncs[0x35] = andZpx
	opFuncs[0x2D] = andAbs
	opFuncs[0x3D] = andAbx
	opFuncs[0x39] = andAby
	opFuncs[0x21] = andIdx
	opFuncs[0x31] = andIdy

	opFuncs[0x0A] = aslAcc
	opFuncs[0x06] = aslZpg
	opFuncs[0x16] = aslZpx
	opFuncs[0x0E] = aslAbs
	opFuncs[0x1E] = aslAbx

	opFuncs[0x90] = bccRel
	opFuncs[0xB0] = bcsRel
	opFuncs[0xF0] = beqRel
	opFuncs[0x30] = bmiRel
	opFuncs[0xD0] = bneRel
	opFuncs[0x10] = bplRel
	opFuncs[0x50] = bvcRel
	opFuncs[0x70] = bvsRel

	opFuncs[0x24] = bitZpg
	opFuncs[0x2C] = bitAbs

	opFuncs[0x00] = brkImp

	opFuncs[0x18] = clcImp
	opFuncs[0xD8] = cldImp
	opFuncs[0x58] = cliImp
	opFuncs[0xB8] = clvImp

	opFuncs[0xC9] = cmpImm
	opFuncs[0xC5] = cmpZpg
	opFuncs[0xD5] = cmpZpx
	opFuncs[0xCD] = cmpAbs
	opFuncs[0xDD] = cmpAbx
	opFuncs[0xD9] = cmpAby
	opFuncs[0xC1] = cmpIdx
	opFuncs[0xD1] = cmpIdy

	opFuncs[0xE0] = cpxImm
	opFuncs[0xE4] = cpxZpg
	opFuncs[0xEC] = cpxAbs

	opFuncs[0xC0] = cpyImm
	opFuncs[0xC4] = cpyZpg
	opFuncs[0xCC] = cpyAbs

	opFuncs[0xC6] = decZpg
	opFuncs[0xD6] = decZpx
	opFuncs[0xCE] = decAbs
	opFuncs[0xDE] = decAbx

	opFuncs[0xCA] = dexImp
	opFuncs[0x88] = deyImp

	opFuncs[0x49] = eorImm
	opFuncs[0x45] = eorZpg
	opFuncs[0x55] = eorZpx
	opFuncs[0x4D] = eorAbs
	opFuncs[0x5D] = eorAbx
	opFuncs[0x59] = eorAby
	opFuncs[0x41] = eorIdx
	opFuncs[0x51] = eorIdy

	opFuncs[0xE6] = incZpg
	opFuncs[0xF6] = incZpx
	opFuncs[0xEE] = incAbs
	opFuncs[0xFE] = incAbx

	opFuncs[0xE8] = inxImp
	opFuncs[0xC8] = inyImp

	opFuncs[0x4C] = jmpAbs
	opFuncs[0x6C] = jmpInd

	opFuncs[0x20] = jsrAbs

	opFuncs[0xA9] = ldaImm
	opFuncs[0xA5] = ldaZpg
	opFuncs[0xB5] = ldaZpx
	opFuncs[0xAD] = ldaAbs
	opFuncs[0xBD] = ldaAbx
	opFuncs[0xB9] = ldaAby
	opFuncs[0xA1] = ldaIdx
	opFuncs[0xB1] = ldaIdy

	opFuncs[0xA2] = ldxImm
	opFuncs[0xA6] = ldxZpg
	opFuncs[0xB6] = ldxZpy
	opFuncs[0xAE] = ldxAbs
	opFuncs[0xBE] = ldxAby

	opFuncs[0xA0] = ldyImm
	opFuncs[0xA4] = ldyZpg
	opFuncs[0xB4] = ldyZpx
	opFuncs[0xAC] = ldyAbs
	opFuncs[0xBC] = ldyAbx

	opFuncs[0x4A] = lsrAcc
	opFuncs[0x46] = lsrZpg
	opFuncs[0x56] = lsrZpx
	opFuncs[0x4E] = lsrAbs
	opFuncs[0x5E] = lsrAbx

	opFuncs[0xEA] = nopImp

	opFuncs[0x09] = oraImm
	opFuncs[0x05] = oraZpg
	opFuncs[0x15] = oraZpx
	opFuncs[0x0D] = oraAbs
	opFuncs[0x1D] = oraAbx
	opFuncs[0x19] = oraAby
	opFuncs[0x01] = oraIdx
	opFuncs[0x11] = oraIdy

	opFuncs[0x48] = phaImp
	opFuncs[0x08] = phpImp
	opFuncs[0x68] = plaImp
	opFuncs[0x28] = plpImp

	opFuncs[0x2A] = rolAcc
	opFuncs[0x26] = rolZpg
	opFuncs[0x36] = rolZpx
	opFuncs[0x2E] = rolAbs
	opFuncs[0x3E] = rolAbx

	opFuncs[0x6A] = rorAcc
	opFuncs[0x66] = rorZpg
	opFuncs[0x76] = rorZpx
	opFuncs[0x6E] = rorAbs
	opFuncs[0x7E] = rorAbx

	opFuncs[0x40] = rtiImp
	opFuncs[0x60] = rtsImp

	opFuncs[0xE9] = sbcImm
	opFuncs[0xE5] = sbcZpg
	opFuncs[0xF5] = sbcZpx
	opFuncs[0xED] = sbcAbs
	opFuncs[0xFD] = sbcAbx
	opFuncs[0xF9] = sbcAby
	opFuncs[0xE1] = sbcIdx
	opFuncs[0xF1] = sbcIdy

	opFuncs[0x38] = secImp
	opFuncs[0xF8] = sedImp
	opFuncs[0x78] = seiImp

	opFuncs[0x85] = staZpg
	opFuncs[0x95] = staZpx
	opFuncs[0x8D] = staAbs
	opFuncs[0x9D] = staAbx
	opFuncs[0x99] = staAby
	opFuncs[0x81] = staIdx
	opFuncs[0x91] = staIdy

	opFuncs[0x86] = stxZpg
	opFuncs[0x96] = stxZpy
	opFuncs[0x8E] = stxAbs

	opFuncs[0x84] = styZpg
	opFuncs[0x94] = styZpx
	opFuncs[0x8C] = styAbs

	opFuncs[0xAA] = taxImp
	opFuncs[0xA8] = tayImp
	opFuncs[0xBA] = tsxImp
	opFuncs[0x8A] = txaImp
	opFuncs[0x9A] = txsImp
	opFuncs[0x98] = tyaImp
}

// initRom() initialises all ROM data to 0xFF.
// This simulates the unprogrammed state of ROMs
func initRom() {
	for i, _ := range rom {
		rom[i] = 0xFF
	}
	flashing = false
	binStart = 0x0000
}

// initRam() initialises all RAM data to 0xFF. 
func initRam() {
	for i, _ := range ram {
		ram[i] = 0xFF
	}
}
