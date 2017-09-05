// Copyright 2012 RVJ Callanan. All rights reserved.

package main

// All of the functions in this file implement CPU instructions.
// Pointers to these functions are stored in the opFuncs array for
// fast access. The first triplet in the function name matches the
// instruction mnemonic and the second triplet identifies the
// addressing mode for the specific opcode as follows:
// 
// Imp: Implicit					e.g. NOP 
// Acc: Accumulator					e.g. LSR A 
// Imm: Immediate					e.g. LDA #10
// Zpg: Zero Page					e.g. LDA $00
// Zpx: Zero Page,X					e.g. STY $10,X
// Zpy: Zero Page,Y					e.g. LDX $10,Y
// Rel: Relative					e.g. BNE *+4
// Abs: Absolute					e.g. JMP $1234
// Abx: Absolute,X					e.g. STA $3000,X
// Aby: Absolute,Y					e.g. AND $4000,Y
// Ind: Indirect					e.g. JMP ($FFFC) 
// Idx: Indexed Indirect using X	e.g. LDA ($40,X) 
// Idy: Indirect Indexed using Y	e.g. LDA ($40),Y

func adcImm() {
	pc += 1
	data := readByte(pc)
	adcCore(data)
	pc += 1
	ck += 2
}

func adcZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	adcCore(data)
	pc += 1
	ck += 3
}

func adcZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	adcCore(data)
	pc += 1
	ck += 4
}

func adcAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	adcCore(data)
	pc += 2
	ck += 4
}

func adcAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	adcCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func adcAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	adcCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func adcIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	adcCore(data)
	pc += 1
	ck += 6
}

func adcIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	adcCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func andImm() {
	pc += 1
	data := readByte(pc)
	andCore(data)
	pc += 1
	ck += 2
}

func andZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	andCore(data)
	pc += 1
	ck += 3
}

func andZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	andCore(data)
	pc += 1
	ck += 4
}

func andAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	andCore(data)
	pc += 2
	ck += 4
}

func andAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	andCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func andAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	andCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func andIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	andCore(data)
	pc += 1
	ck += 6
}

func andIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	andCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func aslAcc() {
	aslCore(&ac)
	pc += 1
	ck += 2
}

func aslZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	aslCore(dst)
	pc += 1
	ck += 5
}

func aslZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	aslCore(dst)
	pc += 1
	ck += 6
}

func aslAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	aslCore(dst)
	pc += 2
	ck += 6
}

func aslAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	aslCore(dst)
	pc += 2
	ck += 7
}

func bccRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if !tstC() {
		braCore(offset)
	}
}

func bcsRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if tstC() {
		braCore(offset)
	}
}

func beqRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if tstZ() {
		braCore(offset)
	}
}

func bmiRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if tstN() {
		braCore(offset)
	}
}

func bneRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if !tstZ() {
		braCore(offset)
	}
}

func bplRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if !tstN() {
		braCore(offset)
	}
}

func bvcRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if !tstV() {
		braCore(offset)
	}
}

func bvsRel() {
	pc += 1
	offset := readByte(pc)
	pc += 1
	ck += 2
	if tstV() {
		braCore(offset)
	}
}

func bitZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	bitCore(data)
	pc += 1
	ck += 3
}

func bitAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	bitCore(data)
	pc += 2
	ck += 4
}

func brkImp() {
	pushWord(pc + 2)
	pushByte(sr)
	setI()
	pc = readWord(irqVec)
	ck += 7
}

func clcImp() {
	clrC()
	pc += 1
	ck += 2
}

func cldImp() {
	clrD()
	pc += 1
	ck += 2
}

func cliImp() {
	clrI()
	pc += 1
	ck += 2
}

func clvImp() {
	clrV()
	pc += 1
	ck += 2
}

func cmpImm() {
	pc += 1
	data := readByte(pc)
	cmpCore(data)
	pc += 1
	ck += 2
}

func cmpZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	cmpCore(data)
	pc += 1
	ck += 3
}

func cmpZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	cmpCore(data)
	pc += 1
	ck += 4
}

func cmpAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	cmpCore(data)
	pc += 2
	ck += 4
}

func cmpAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	cmpCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func cmpAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	cmpCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func cmpIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	cmpCore(data)
	pc += 1
	ck += 6
}

func cmpIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	cmpCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func cpxImm() {
	pc += 1
	data := readByte(pc)
	cpxCore(data)
	pc += 1
	ck += 2
}

func cpxZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	cpxCore(data)
	pc += 1
	ck += 3
}

func cpxAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	cpxCore(data)
	pc += 2
	ck += 4
}

func cpyImm() {
	pc += 1
	data := readByte(pc)
	cpyCore(data)
	pc += 1
	ck += 2
}

func cpyZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	cpyCore(data)
	pc += 1
	ck += 3
}

func cpyAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	cpyCore(data)
	pc += 2
	ck += 4
}

func decZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	decCore(dst)
	pc += 1
	ck += 5
}

func decZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	decCore(dst)
	pc += 1
	ck += 6
}

func decAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	decCore(dst)
	pc += 2
	ck += 6
}

func decAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	decCore(dst)
	pc += 2
	ck += 7
}

func dexImp() {
	ix -= 1
	chgZ(ix == 0)
	chgN(ix > 127)
	pc += 1
	ck += 2
}

func deyImp() {
	iy -= 1
	chgZ(iy == 0)
	chgN(iy > 127)
	pc += 1
	ck += 2
}

func eorImm() {
	pc += 1
	data := readByte(pc)
	eorCore(data)
	pc += 1
	ck += 2
}

func eorZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	eorCore(data)
	pc += 1
	ck += 3
}

func eorZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	eorCore(data)
	pc += 1
	ck += 4
}

func eorAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	eorCore(data)
	pc += 2
	ck += 4
}

func eorAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	eorCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func eorAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	eorCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func eorIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	eorCore(data)
	pc += 1
	ck += 6
}

func eorIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	eorCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func incZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	incCore(dst)
	pc += 1
	ck += 5
}

func incZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	incCore(dst)
	pc += 1
	ck += 6
}

func incAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	incCore(dst)
	pc += 2
	ck += 6
}

func incAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	incCore(dst)
	pc += 2
	ck += 7
}

func inxImp() {
	ix += 1
	chgZ(ix == 0)
	chgN(ix > 127)
	pc += 1
	ck += 2
}

func inyImp() {
	iy += 1
	chgZ(iy == 0)
	chgN(iy > 127)
	pc += 1
	ck += 2
}

func jmpAbs() {
	pc += 1
	pc = readWord(pc)
	ck += 3
}

func jmpInd() {
	pc += 1
	addr := readWord(pc)
	pc = readWord(addr)
	ck += 5
}

func jsrAbs() {
	pc += 1
	pushWord(pc + 1)
	pc = readWord(pc)
	ck += 6
}

func ldaImm() {
	pc += 1
	data := readByte(pc)
	ldaCore(data)
	pc += 1
	ck += 2
}

func ldaZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	ldaCore(data)
	pc += 1
	ck += 3
}

func ldaZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	ldaCore(data)
	pc += 1
	ck += 4
}

func ldaAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	ldaCore(data)
	pc += 2
	ck += 4
}

func ldaAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	ldaCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func ldaAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	ldaCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func ldaIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	ldaCore(data)
	pc += 1
	ck += 6
}

func ldaIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	ldaCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func ldxImm() {
	pc += 1
	data := readByte(pc)
	ldxCore(data)
	pc += 1
	ck += 2
}

func ldxZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	ldxCore(data)
	pc += 1
	ck += 3
}

func ldxZpy() {
	pc += 1
	addr := uint16(iy + readByte(pc))
	data := readByte(addr)
	ldxCore(data)
	pc += 1
	ck += 4
}

func ldxAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	ldxCore(data)
	pc += 2
	ck += 4
}

func ldxAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	ldxCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func ldyImm() {
	pc += 1
	data := readByte(pc)
	ldyCore(data)
	pc += 1
	ck += 2
}

func ldyZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	ldyCore(data)
	pc += 1
	ck += 3
}

func ldyZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	ldyCore(data)
	pc += 1
	ck += 4
}

func ldyAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	ldyCore(data)
	pc += 2
	ck += 4
}

func ldyAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	ldyCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func lsrAcc() {
	lsrCore(&ac)
	pc += 1
	ck += 2
}

func lsrZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	lsrCore(dst)
	pc += 1
	ck += 5
}

func lsrZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	lsrCore(dst)
	pc += 1
	ck += 6
}

func lsrAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	lsrCore(dst)
	pc += 2
	ck += 6
}

func lsrAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	lsrCore(dst)
	pc += 2
	ck += 7
}

func nopImp() {
	pc += 1
	ck += 2
}

func oraImm() {
	pc += 1
	data := readByte(pc)
	oraCore(data)
	pc += 1
	ck += 2
}

func oraZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	oraCore(data)
	pc += 1
	ck += 3
}

func oraZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	oraCore(data)
	pc += 1
	ck += 4
}

func oraAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	oraCore(data)
	pc += 2
	ck += 4
}

func oraAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	oraCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func oraAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	oraCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func oraIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	oraCore(data)
	pc += 1
	ck += 6
}

func oraIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	oraCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func phaImp() {
	pushByte(ac)
	pc += 1
	ck += 3
}

func phpImp() {
	pushByte(sr)
	pc += 1
	ck += 3
}

func plaImp() {
	ac = popByte()
	chgZ(ac == 0)
	chgN(ac > 127)
	pc += 1
	ck += 4
}

func plpImp() {
	sr = popByte() | maskU | maskB
	pc += 1
	ck += 4
}

func rolAcc() {
	rolCore(&ac)
	pc += 1
	ck += 2
}

func rolZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	rolCore(dst)
	pc += 1
	ck += 5
}

func rolZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	rolCore(dst)
	pc += 1
	ck += 6
}

func rolAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	rolCore(dst)
	pc += 2
	ck += 6
}

func rolAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	rolCore(dst)
	pc += 2
	ck += 7
}

func rorAcc() {
	rorCore(&ac)
	pc += 1
	ck += 2
}

func rorZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	rorCore(dst)
	pc += 1
	ck += 5
}

func rorZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	rorCore(dst)
	pc += 1
	ck += 6
}

func rorAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	rorCore(dst)
	pc += 2
	ck += 6
}

func rorAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	rorCore(dst)
	pc += 2
	ck += 7
}

func rtiImp() {
	sr = popByte() | maskU | maskB
	pc = popWord()
	ck += 6
}

func rtsImp() {
	pc = popWord()
	pc += 1
	ck += 6
}

func sbcImm() {
	pc += 1
	data := readByte(pc)
	sbcCore(data)
	pc += 1
	ck += 2
}

func sbcZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	data := readByte(addr)
	sbcCore(data)
	pc += 1
	ck += 3
}

func sbcZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	data := readByte(addr)
	sbcCore(data)
	pc += 1
	ck += 4
}

func sbcAbs() {
	pc += 1
	addr := readWord(pc)
	data := readByte(addr)
	sbcCore(data)
	pc += 2
	ck += 4
}

func sbcAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	data := readByte(addr)
	sbcCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func sbcAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	data := readByte(addr)
	sbcCore(data)
	pc += 2
	if (base >> 8) == (addr >> 8) {
		ck += 4
	} else {
		ck += 5
	}
}

func sbcIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	data := readByte(addr)
	sbcCore(data)
	pc += 1
	ck += 6
}

func sbcIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	data := readByte(addr)
	sbcCore(data)
	pc += 1
	if (base >> 8) == (addr >> 8) {
		ck += 5
	} else {
		ck += 6
	}
}

func secImp() {
	setC()
	pc += 1
	ck += 2
}

func sedImp() {
	setD()
	pc += 1
	ck += 2
}

func seiImp() {
	setI()
	pc += 1
	ck += 2
}

func staZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	*dst = ac
	pc += 1
	ck += 3
}

func staZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	*dst = ac
	pc += 1
	ck += 4
}

func staAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	*dst = ac
	pc += 2
	ck += 4
}

func staAbx() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(ix)
	dst := refByte(addr)
	*dst = ac
	pc += 2
	ck += 5
}

func staAby() {
	pc += 1
	base := readWord(pc)
	addr := base + uint16(iy)
	dst := refByte(addr)
	*dst = ac
	pc += 2
	ck += 5
}

func staIdx() {
	pc += 1
	vec := uint16(ix + readByte(pc))
	addr := readWord(vec)
	dst := refByte(addr)
	*dst = ac
	pc += 1
	ck += 6
}

func staIdy() {
	pc += 1
	vec := uint16(readByte(pc))
	base := readWord(vec)
	addr := base + uint16(iy)
	dst := refByte(addr)
	*dst = ac
	pc += 1
	ck += 6
}

func stxZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	*dst = ix
	pc += 1
	ck += 3
}

func stxZpy() {
	pc += 1
	addr := uint16(iy + readByte(pc))
	dst := refByte(addr)
	*dst = ix
	pc += 1
	ck += 4
}

func stxAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	*dst = ix
	pc += 2
	ck += 4
}

func styZpg() {
	pc += 1
	addr := uint16(readByte(pc))
	dst := refByte(addr)
	*dst = iy
	pc += 1
	ck += 3
}

func styZpx() {
	pc += 1
	addr := uint16(ix + readByte(pc))
	dst := refByte(addr)
	*dst = iy
	pc += 1
	ck += 4
}

func styAbs() {
	pc += 1
	addr := readWord(pc)
	dst := refByte(addr)
	*dst = iy
	pc += 2
	ck += 4
}

func taxImp() {
	ix = ac
	chgZ(ix == 0)
	chgN(ix > 127)
	pc += 1
	ck += 2
}

func tayImp() {
	iy = ac
	chgZ(iy == 0)
	chgN(iy > 127)
	pc += 1
	ck += 2
}

func tsxImp() {
	ix = sp
	chgZ(ix == 0)
	chgN(ix > 127)
	pc += 1
	ck += 2
}

func txaImp() {
	ac = ix
	chgZ(ac == 0)
	chgN(ac > 127)
	pc += 1
	ck += 2
}

func txsImp() {
	sp = ix
	pc += 1
	ck += 2
}

func tyaImp() {
	ac = iy
	chgZ(ac == 0)
	chgN(ac > 127)
	pc += 1
	ck += 2
}
