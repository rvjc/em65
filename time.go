// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"runtime"
	"time"
)

func absDuration(t time.Duration) (a time.Duration) {
	if t >= 0 {
		a = t
	} else {
		a = -t
	}
	return
}

// sync() keeps cpu time synchronised with real time. It is
// called by the main program loop on each syncTicks interval.
// Synchronisation is achieved by forcing the program to sleep
// just long enough to let real time "catch up" with cpu time. 
// Excessive time drift (in either direction) will force a resync 
// from a new reference point in time. This will happen in the
// following scenarios:
// 1. When the user is stepping while in debug mode.
// 2. When the system is too slow or is over-loaded.
// 3. When the system clock is adjusted during program execution.
// The sync() function is an opportune time for garbage collection
// and scheduling pending goroutines. A minimal sleep is always
// enforced even when the cpu time is falling behind. This ensures
// that the emulator will not hog the system completely. And by
// being a good citizen, it is less likely to be pre-emptively
// interrupted by the OS. The obligatory sleep call also ensures
// that the sampled system time is up to date.
func sync() {
	syncCount += 1
	if syncCount%syncsPerGc == 0 {
		runtime.GC()
	}
	runtime.Gosched()
	time.Sleep(minSleep)
	now := time.Now()
	realTime := now.Sub(syncRefReal)
	cpuTime := time.Duration((ck - syncRefCk) * cpuTick)
	diffTime := cpuTime - realTime
	switch {
	case absDuration(diffTime) > resyncThresh:
		syncRefReal = now
		syncRefCk = ck
	case diffTime > minSleep:
		time.Sleep(diffTime)
	}
	syncNextCk += ticksPerSync
}
