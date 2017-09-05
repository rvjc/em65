// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"fmt"
	"strings"
)

// fmtBool() customs formats a boolean value.	
func fmtBool(b bool, trueFmt string, falseFmt string) (f string) {
	if b {
		f = trueFmt
	} else {
		f = falseFmt
	}
	return
}

func fmtByte(val uint8) string  { return fmt.Sprintf("%02X", val) }
func fmtWord(val uint16) string { return fmt.Sprintf("%04X", val) }

func fmtCk() string { return fmt.Sprintf("%011d", ck) }
func fmtOp() string { return fmt.Sprintf("%02X", op) }
func fmtPc() string { return fmt.Sprintf("%04X", pc) }
func fmtAc() string { return fmt.Sprintf("%02X", ac) }
func fmtIx() string { return fmt.Sprintf("%02X", ix) }
func fmtIy() string { return fmt.Sprintf("%02X", iy) }
func fmtSp() string { return fmt.Sprintf("%02X", sp) }
func fmtSr() string { return fmt.Sprintf("%02X", sr) }

func fmtFlags() (s string) {
	s = ""
	s += fmtBool(tstN(), "N", "-")
	s += fmtBool(tstV(), "V", "-")
	s += fmtBool(tstU(), "*", "!")
	s += fmtBool(tstB(), "*", "!")
	s += fmtBool(tstD(), "D", "-")
	s += fmtBool(tstI(), "I", "-")
	s += fmtBool(tstZ(), "Z", "-")
	s += fmtBool(tstC(), "C", "-")
	return
}

func fmtSrc(addr uint16) (s string) {
	sd := src[addr]
	fb := ""
	for i := 0; i < sd.byteCount; i++ {
		b := readByte(pc + uint16(i))
		fb += fmtByte(b)
	}
	fb = fmt.Sprintf("%-6s", fb)[:6]
	fl := fmt.Sprintf("%-12s", sd.label)[:12]
	fm := fmt.Sprintf("%-3s", sd.mnem)
	fo := fmt.Sprintf("%-14s", sd.operand)[:14]
	return fb + " " + fl + " " + fm + " " + fo
}

func fmtState() string {
	fields := []string{
		fmtCk(),
		fmtAc(),
		fmtIx(),
		fmtIy(),
		fmtSp(),
		fmtFlags(),
		fmtPc(),
		fmtSrc(pc),
	}
	return strings.Join(fields, " ")
}
