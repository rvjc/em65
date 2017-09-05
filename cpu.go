// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"time"
)

// Status Register flag operations

func tstN() bool { return sr&maskN > 0 }
func tstV() bool { return sr&maskV > 0 }
func tstU() bool { return sr&maskU > 0 }
func tstB() bool { return sr&maskB > 0 }
func tstD() bool { return sr&maskD > 0 }
func tstI() bool { return sr&maskI > 0 }
func tstZ() bool { return sr&maskZ > 0 }
func tstC() bool { return sr&maskC > 0 }

func clrN() { sr &= ^maskN }
func clrV() { sr &= ^maskV }
func clrU() { sr &= ^maskU }
func clrB() { sr &= ^maskB }
func clrD() { sr &= ^maskD }
func clrI() { sr &= ^maskI }
func clrZ() { sr &= ^maskZ }
func clrC() { sr &= ^maskC }

func setN() { sr |= maskN }
func setV() { sr |= maskV }
func setU() { sr |= maskU }
func setB() { sr |= maskB }
func setD() { sr |= maskD }
func setI() { sr |= maskI }
func setZ() { sr |= maskZ }
func setC() { sr |= maskC }

func chgN(val bool) {
	if val {
		sr |= maskN
	} else {
		sr &= ^maskN
	}
}

func chgV(val bool) {
	if val {
		sr |= maskV
	} else {
		sr &= ^maskV
	}
}

func chgU(val bool) {
	if val {
		sr |= maskU
	} else {
		sr &= ^maskU
	}
}

func chgB(val bool) {
	if val {
		sr |= maskB
	} else {
		sr &= ^maskB
	}
}

func chgD(val bool) {
	if val {
		sr |= maskD
	} else {
		sr &= ^maskD
	}
}

func chgI(val bool) {
	if val {
		sr |= maskI
	} else {
		sr &= ^maskI
	}
}

func chgZ(val bool) {
	if val {
		sr |= maskZ
	} else {
		sr &= ^maskZ
	}
}

func chgC(val bool) {
	if val {
		sr |= maskC
	} else {
		sr &= ^maskC
	}
}

// reset() clears CPU state and starts running from the reset vector.
func reset() {
	op = 0x00
	ac = 0x00
	ix = 0x00
	iy = 0x00
	sr = 0x00 | maskU | maskB
	sp = spMax
	pc = readWord(rstVec)
	ck = 0
	time.Sleep(minSleep)
	syncRefReal = time.Now()
	syncRefCk = 0
	syncNextCk = ticksPerSync
	syncCount = 0
}

// readByte() reads a byte from memory.
// Unused memory areas read as all-ones to simulate pull-up
// resistors on a typical data bus.
func readByte(addr uint16) (data uint8) {
	switch {
	case addr >= romMin && addr <= romMax:
		data = rom[addr-romMin]
	case addr >= ramMin && addr <= ramMax:
		data = ram[addr-ramMin]
	default:
		data = 0xFF
	}
	return
}

// writeByte() writes a byte to memory.
// Unused or ROM addresses are normally ignored, as would be 
// the case with typical hardware. However, ROM data will be 
// over-written when the global flash flag is set
func writeByte(addr uint16, data uint8) {
	switch {
	case addr >= ramMin && addr <= ramMax:
		ram[addr-ramMin] = data
	case addr >= romMin && addr <= romMax:
		if flashing {
			rom[addr-romMin] = data
		}
	}
}

// readWord() reads a word from memory as two bytes.
// The lower byte is read from the specified address.
// The upper byte is read from the following address. 
func readWord(addr uint16) uint16 {
	lo := uint16(readByte(addr))
	hi := uint16(readByte(addr + 1))
	return lo | (hi << 8)
}

// writeWord() writes a word to memory as two bytes.
// The lower byte is written to the specified address.
// The upper byte is written to the following address. 
func writeWord(addr uint16, data uint16) {
	lo := uint8(data & 0xFF)
	hi := uint8(data >> 8)
	writeByte(addr, lo)
	writeByte(addr+1, hi)
}

// refByte() returns a reference to a byte in memory.
// This allows memory data to be accessed in place
// for efficient read-modify-write operations. If the
// target address is not writeable, a reference to a
// dummy byte is returned. If the non-writable addresses 
// is readable, its read value is copied to the dummy
// byte. If the target address is completely inaccessible,
// the dummy byte is set to 0xFF.
func refByte(addr uint16) (ref *uint8) {
	switch {
	case addr >= ramMin && addr <= ramMax:
		ref = &ram[addr-ramMin]
	case addr >= romMin && addr <= romMax:
		if flashing {
			ref = &rom[addr-romMin]
		} else {
			ref = new(uint8)
			*ref = rom[addr-romMin]
		}
	default:
		ref = new(uint8)
		*ref = 0xFF
	}
	return
}

// pushByte() saves byte to stack and decrements stack pointer
func pushByte(data uint8) {
	writeByte(saMin+uint16(sp), data)
	sp--
}

// popByte() increments stack pointer and reads byte from stack
func popByte() uint8 {
	sp++
	return readByte(saMin + uint16(sp))
}

// popWord() pops word from stack as two consecutive bytes
// The upper byte is pushed first because the SP is decrementing
func pushWord(data uint16) {
	lo := uint8(data & 0xFF)
	hi := uint8(data >> 8)
	pushByte(hi)
	pushByte(lo)
}

// popWord() pops word from stack as two consecutive bytes
// The lower byte is popped first because the SP is incrementing
func popWord() uint16 {
	lo := uint16(popByte())
	hi := uint16(popByte())
	return lo | (hi << 8)
}

// adcCore() performs the core operation common to all ADC instructions.
// The data argument is added to the accumulator taking account of the
// C flag. The D flag enables Binary-Coded Decimal (BCD) addition. The
// result is reflected in the N, V, Z and C flags. The 6502 has well-
// known arithmetic "quirks" which must be faithfully reproduced. This
// is especially true in BCD mode which is "bolted onto" binary mode. 

func adcCore(data uint8) {
	a := uint16(ac)
	d := uint16(data)
	c := uint16(sr & maskC)
	r := a + d + c
	chgZ(r&0xFF == 0)
	if tstD() {
		if (a&0x0F)+(d&0x0F)+c > 0x09 {
			r += 0x06
		}
		chgN(r&0x80 > 0)
		chgV(((a^^d)&(a^r))&0x80 > 0)
		if r > 0x99 {
			r += 0x60
		}
	} else {
		chgN(r&0x80 > 0)
		chgV(((a^^d)&(a^r))&0x80 > 0)
	}
	chgC(r > 0xFF)
	ac = uint8(r & 0xFF)
}

// andCore() performs the core operation common to all AND instructions.
// The N and Z flags are modified to reflect the result.
func andCore(data uint8) {
	ac &= data
	chgN(ac > 127)
	chgZ(ac == 0)
}

// aslCore() performs the core operation common to all ASL instructions.
// The N, Z and C flags are modified to reflect the result.
func aslCore(dst *uint8) {
	val := *dst
	chgC(val > 127)
	val = val << 1
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// bitCore() performs the core operation common to all BIT instructions.
// The mask pattern in the accumulator is ANDed with the data byte and
// the result is reflected in the Z flag, but the result is not kept.
// Bits 7 and 6 of the data byte are copied to the N and Z flags.
func bitCore(data uint8) {
	bits := ac & data
	chgZ(bits == 0)
	chgN(data&0x80 > 0)
	chgV(data&0x40 > 0)
}

// braCore() performs the core operation common to all BRA instructions.
// The program counter is shifted up or down by the byte offset.
// The CPU clock is incremented by 1 if the branch is on the same page
// but it is shifted by 2 if the branch crosses a page boundary.
func braCore(offset uint8) {
	oldpc := int(pc)
	shift := int(int8(offset))
	newpc := oldpc + shift
	if (newpc >> 8) == (oldpc >> 8) {
		ck += 1
	} else {
		ck += 2
	}
	pc = uint16(newpc)
}

// cmpCore() performs the core operation common to all CMP instructions.
// The data byte is subtracted from the value in the accumulator but
// the accumulator is not updated. The result is reflected in the C,N
// and Z flags. The subtration is similar to a binary SBC operation,
// however the carry flag is ignored and the V flag is not modified.
// The carry flag is modified in the same way as SBC. 
func cmpCore(data uint8) {
	a := uint16(ac)
	d := uint16(data)
	t := a - d
	chgC(t <= 0xFF)
	chgN(t&0x80 > 0)
	chgZ(t&0xFF == 0)
}

// cpxCore() performs the core operation common to all CPX instructions.
// The data byte is subtracted from the value in the X register but
// the X register is not updated. The result is reflected in the C,N
// and Z flags as per the CMP instruction.
func cpxCore(data uint8) {
	x := uint16(ix)
	d := uint16(data)
	t := x - d
	chgC(t <= 0xFF)
	chgN(t&0x80 > 0)
	chgZ(t&0xFF == 0)
}

// cpyCore() performs the core operation common to all CPY instructions.
// The data byte is subtracted from the value in the Y register but
// the X register is not updated. The result is reflected in the C,N
// and Z flags as per the CMP instruction.
func cpyCore(data uint8) {
	y := uint16(iy)
	d := uint16(data)
	t := y - d
	chgC(t <= 0xFF)
	chgN(t&0x80 > 0)
	chgZ(t&0xFF == 0)
}

// decCore() performs the core operation common to all DEC instructions.
// The destination is decremented by 1 using standard binary subtraction.
// The N and Z flags are modified to reflect the result.
func decCore(dst *uint8) {
	val := *dst
	val -= 1
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// eorCore() performs the core operation common to all EOR instructions.
// The N and Z flags are modified to reflect the result.
func eorCore(data uint8) {
	ac ^= data
	chgN(ac > 127)
	chgZ(ac == 0)
}

// incCore() performs the core operation common to all INC instructions.
// The destination is incremented by 1 using standard binary addition.
// The N and Z flags are modified to reflect the result.
func incCore(dst *uint8) {
	val := *dst
	val += 1
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// ldaCore() performs the core operation common to all LDA instructions.
// The N and Z flags are modified to reflect the loaded data byte.
func ldaCore(data uint8) {
	ac = data
	chgN(ac > 127)
	chgZ(ac == 0)
}

// ldxCore() performs the core operation common to all LDX instructions.
// The N and Z flags are modified to reflect the loaded data byte.
func ldxCore(data uint8) {
	ix = data
	chgN(ix > 127)
	chgZ(ix == 0)
}

// ldyCore() performs the core operation common to all LDY instructions.
// The N and Z flags are modified to reflect the loaded data byte.
func ldyCore(data uint8) {
	iy = data
	chgN(iy > 127)
	chgZ(iy == 0)
}

// lsrCore() performs the core operation common to all LSR instructions.
// The N, Z and C flags are modified to reflect the result.
func lsrCore(dst *uint8) {
	val := *dst
	chgC(val&0x01 > 0)
	val = val >> 1
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// oraCore() performs the core operation common to all ORA instructions.
// The N and Z flags are modified to reflect the result.
func oraCore(data uint8) {
	ac |= data
	chgN(ac > 127)
	chgZ(ac == 0)
}

// rolCore() performs the core operation common to all ROL instructions.
// Rotates all bits in the destination to the left by one position.
// The C Flag is "rotated" into Bit 0 and Bit 7 is "rotated" back to
// the C flag. The N and Z flags are modified to reflect the final
// value of the destination.
func rolCore(dst *uint8) {
	lsb := sr & maskC
	val := *dst
	chgC(val&0x80 > 0)
	val = (val << 1) | lsb
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// rorCore() performs the core operation common to all ROr instructions.
// Rotates all bits in the destination to the right by one position.
// The C Flag is "rotated" into Bit 7 and Bit 0 is "rotated" back to
// the C flag. The N and Z flags are modified to reflect the final
// value of the destination.
func rorCore(dst *uint8) {
	msb := sr & maskC << 7
	val := *dst
	chgC(val&0x01 > 0)
	val = (val >> 1) | msb
	chgN(val > 127)
	chgZ(val == 0)
	*dst = val
}

// sbcCore() performs the core operation common to all SBC instructions.
// The data argument is subtracted from the accumulator taking account
// of the C flag (which has the OPPOSITE meaning to ADC). The D flag
// enables BCD subtraction. The result of the operation is reflected
// in the N, V, Z and C flags. 6502 subtraction has even more "quirks"
// than addition which must be faithfully reproduced. 
func sbcCore(data uint8) {
	a := uint16(ac)
	d := uint16(data)
	c := uint16(^sr & maskC)
	r := a - d - c
	chgZ(r&0xFF == 0)
	chgN(r&0x80 > 0)
	chgV(((a^d)&(a^r))&0x80 > 0)
	if tstD() {
		if a&0x0F < (d&0x0F + c) {
			r -= 0x06
		}
		if r > 0x99 {
			r -= 0x60
		}
	}
	chgC(r <= 0xFF)
	ac = uint8(r & 0xFF)
}
