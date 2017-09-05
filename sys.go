// Copyright 2012 RVJ Callanan. All rights reserved.

package main

import (
	"fmt"
	"os"
)

func osSigHandler() {
	for sig := range osSigs {
		switch sig {
		case os.Interrupt:
			if active {
				if stepping {
					os.Exit(1)
				} else {
					fmt.Println("\nInterrupted\n")
					debugging = true
					stepping = true
				}
			}
		case os.Kill:
			fmt.Println("Forced Termination")
			os.Exit(1)
		}
	}
}
