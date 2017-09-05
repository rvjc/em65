// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"os"
	"time"
)

// Exception Vectors	
const (
	nmiVec uint16 = 0xFFFA // Non-Maskable Interrupt
	rstVec uint16 = 0xFFFC // Reset
	irqVec uint16 = 0xFFFE // Interrupt Request (or Break)
)

// Status Register Flag Masks
const (
	maskC uint8 = 0x01 // Carry Flag
	maskZ uint8 = 0x02 // Zero Flag
	maskI uint8 = 0x04 // Interrupt Disable
	maskD uint8 = 0x08 // Decimal Mode
	maskB uint8 = 0x10 // Break Command
	maskU uint8 = 0x20 // Unused Flag (always set to 1)
	maskV uint8 = 0x40 // Overflow Flag
	maskN uint8 = 0x80 // Negative Flag
)

// Memory Parameters
const (
	memMin  uint16 = 0x0000
	memMax  uint16 = 0xFFFF
	memSize uint32 = 0x10000
	romSize uint16 = 0x8000 // 32KB
	romMax  uint16 = memMax
	romMin  uint16 = romMax - romSize + 1
	ramSize uint16 = 0x8000 // 32KB
	ramMin  uint16 = memMin
	ramMax  uint16 = ramMin + ramSize - 1
	spMin   uint8  = 0x00
	spMax   uint8  = 0xFF
	saMin   uint16 = 0x0100
	saMax   uint16 = saMin + uint16(spMax)
)

// Command post action return codes
const (
	postActionHold     = iota // Hold current CPU state
	postActionContinue        // Execute current instruction
	postActionRefetch         // Fetch different instruction
	postActionQuit            // Quit emulator
)

// Sync variables
const (
	cpuFreq      uint64        = 1E6           // 1MHz
	syncFreq     uint64        = 60            // 60Hz
	gcInterval   uint64        = 10E9          // 10S
	resyncThresh time.Duration = 1E9           // 1S
	minSleep     time.Duration = 1E6           // 1mS
	cpuTick      uint64        = 1E9 / cpuFreq // in nS units
	ticksPerSync uint64        = cpuFreq / syncFreq
	syncsPerGc   uint64        = (syncFreq * gcInterval) / 1E9
)

// srcData contains parsed original source for a given line/address.
type srcData struct {
	byteCount int
	label     string
	mnem      string
	operand   string
	comment   string
}

type postAction int

// mnems is a fast lookup table for 6502 mnemonics.
var mnems = map[string]bool{
	"adc": true, "and": true, "asl": true, "bcc": true, "bcs": true,
	"beq": true, "bit": true, "bmi": true, "bne": true, "bpl": true,
	"brk": true, "bvc": true, "bvs": true, "clc": true, "cld": true,
	"cli": true, "clv": true, "cmp": true, "cpx": true, "cpy": true,
	"dec": true, "dex": true, "dey": true, "eor": true, "inc": true,
	"inx": true, "iny": true, "jmp": true, "jsr": true, "lda": true,
	"ldx": true, "ldy": true, "lsr": true, "nop": true, "ora": true,
	"pha": true, "php": true, "pla": true, "plp": true, "rol": true,
	"ror": true, "rti": true, "rts": true, "sbc": true, "sec": true,
	"sed": true, "sei": true, "sta": true, "stx": true, "sty": true,
	"tax": true, "tay": true, "tsx": true, "txa": true, "txs": true,
	"tya": true,
}

// opFuncs is a fast lookup table for implemented opcode functions.
var opFuncs []func() = make([]func(), 256)

// Memory-mapped devices
var rom = make([]uint8, romSize) // ROM Memory
var ram = make([]uint8, ramSize) // RAM Memory

// Program data
var binStart uint16 // First address of binary data in ROM
var src = make(map[uint16]srcData)

// System variables
var active bool           // Emulator is running or stepping through code
var osSigs chan os.Signal // Operating System Signals

// Monitor variables
var debugging bool // Currently in debug mode
var stepping bool  // Currently stepping (in debug mode)
var flashing bool  // Currently allowing writes to ROM

// CPU State
var ck uint64 // CPU Cycle Clock
var op uint8  // Current opcode
var pc uint16 // Program Counter
var ac uint8  // Accumulator
var ix uint8  // Index Register X
var iy uint8  // Index Register Y
var sp uint8  // Stack Pointer (Offset)
var sr uint8  // Status Register

//Sync Variables
var syncRefReal time.Time
var syncRefCk uint64
var syncNextCk uint64
var syncCount uint64

// Breakpoint Variables
var brkCK uint64 // CPU Cycle Clock
var brkPC uint16 // previous program counter 
