/*
The EM65 Emulator.

The EM65 emulator supports the 6502 CPU family.

Operations

Load BIN

Loads raw binary data into memory.

Since binary files have no address information, binary data is assumed
to be in sequence with the last byte at $FFFF. The start address can
thus be determined from the total number of bytes in the file. This
assumption is reasonable because the exception vectors are located are
located at the top-of-memory. To ensure that the binary file ends at $FFFF,
specify all three vectors in the source file.

Load LST

Loads source code from the assembler LST file output.

Any source code extracted in this way is shown when the emulator is
stepping.

For example, the as65 -l option includes each line of original_source 
along with its address and up to 3 code bytes in the following format:

<hex_address> : <hex_code_pairs> <original_source>

Example:

f00d : 2016f0  tst_ok   jsr done  ; success

Any line that complies with the above format is loaded once its address
and code bytes match the previously loaded binary file. All other lines,
including those with pseudo-instructions are ignored. However a warning
is issued if code errors are found. This is a useful verification step 
to flag potential alignment problems in the binary file.

In debug mode, performance is sacrificed slightly in order to detect
breakpoints. When stepping in debug mode, the CPU state is printed at
each step followed by a user command prompt. The post action code returned
by each command determines whether to hold the current CPU state, execute
the current instruction or quit the emulator altogether. An illegal
instruction will cause the emulator to quit whether or not debug mode. 
*/
package documentation
