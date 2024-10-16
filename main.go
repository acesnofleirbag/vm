// specs: <https://www.jmeiners.com/lc3-vm/supplies/lc3-isa.pdf>
package main

import (
	"fmt"
	"os"
)

const START = 0x3000

type Layout [1 << 15]uint16

var MEM Layout
var REG [Reg__Count]Reg

type Reg int

const (
	Reg__R0 Reg = iota
	Reg__R1
	Reg__R2
	Reg__R3
	Reg__R4
	Reg__R5
	Reg__R6
	Reg__R7
	Reg__PC   // Program Counter
	Reg__Cond // Condition flags
	Reg__Count
)

type OpCode int

const (
	OpCode__BR   OpCode = iota // branch
	OpCode__ADD                // add
	OpCode__LD                 // load
	OpCode__ST                 // store
	OpCode__JSR                // jump register
	OpCode__AND                // bitwise and
	OpCode__LDR                // load register
	OpCode__STR                // store register
	OpCode__RTI                // unused
	OpCode__NOT                // bitwise not
	OpCode__LDI                // load indirect
	OpCode__STI                // store indirect
	OpCode__JMP                // jump
	OpCode__RES                // reserved (unused)
	OpCode__LEA                // load effective address
	OpCode__TRAP               // execute trap
)

type Flag int

const (
	Flag__POS Flag = 1 << iota // positive
	Flag__ZRO                  // zero
	Flag__NEG                  // negative
)

func SignExtend(num uint16, bitc int) uint16 {
	if (num >> (bitc - 1) & 1) != 0 {
		num |= 0xFFFF << bitc
	}

	return num
}

func SetFlag(reg Reg) {
	if REG[reg] == 0 {
		REG[Reg__Cond] = Reg(Flag__ZRO)
	} else if REG[reg]>>15 == 1 {
		// 1 in the left-most bit indicates a negative number
		REG[Reg__Cond] = Reg(Flag__NEG)
	} else {
		REG[Reg__Cond] = Reg(Flag__POS)
	}
}

// TODO
func MemRead(num uint16) uint16 {
	return uint16(2)
}

// Instruction Layout (ADD with imm flag off)
//
// 15     12 11     9 8      6    5     4      3 2       0
// +--------+--------+--------+--------+--------+--------+
// |        |        |        |        |        |        |
// |        |        |        |        |        |        |
// |   0001 |     DR |    SR1 |      0 |     00 |    SR2 |
// +--------+--------+--------+--------+--------+--------+

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("vm [image]...\n")

		os.Exit(2)
	}

	for i := 1; i < len(args); i++ {
		fmt.Printf("Failed to load image: %s\n", args[i])

		os.Exit(1)
	}

	REG[Reg__Cond] = Reg(Flag__ZRO)
	REG[Reg__PC] = Reg(START)

	running := 1

	for running > 0 {
		cmd := REG[Reg__PC]
		op := cmd >> 12

		switch op {
		case Reg(OpCode__ADD):
			r0 := (cmd >> 9) & 0x7      // destination register (DR)
			r1 := (cmd >> 6) & 0x7      // first operand (SR1)
			immFlag := (cmd >> 5) & 0x1 // whether we are in immediate mode

			if immFlag == 1 {
				imm5 := SignExtend(uint16(cmd&0x1F), 5)

				REG[r0] = Reg(uint16(REG[r1]) + imm5)
			} else {
				r2 := cmd & 0x7 // second operand (SR2)

				REG[r0] = REG[r1] + REG[r2]
			}

			SetFlag(r0)

			break
		case Reg(OpCode__AND):
			r0 := cmd >> 9 & 0x7        // destination register (DR)
			r1 := cmd >> 6 & 0x7        // first operand (SR1)
			immFlag := (cmd >> 5) & 0x1 // whether we are in immediate mode

			if immFlag == 1 {
				imm5 := SignExtend(uint16(cmd&0x1FF), 5)

				REG[r0] = Reg(uint16(REG[r1]) & imm5)
			} else {
				r2 := cmd & 0x7 // second operand (SR2)

				REG[r0] = REG[r1] + REG[r2]
			}

			break
		case Reg(OpCode__NOT):
			break
		case Reg(OpCode__BR):
			break
		case Reg(OpCode__JMP):
			break
		case Reg(OpCode__JSR):
			break
		case Reg(OpCode__LD):
			break
		case Reg(OpCode__LDI):
			r0 := (cmd >> 9) & 0x7

			offset := SignExtend(uint16(cmd&0x1FF), 9)

			REG[r0] = Reg(MemRead(MemRead(uint16(REG[Reg__PC]) + offset)))

			SetFlag(r0)

			break
		case Reg(OpCode__LDR):
			break
		case Reg(OpCode__LEA):
			break
		case Reg(OpCode__ST):
			break
		case Reg(OpCode__STI):
			break
		case Reg(OpCode__STR):
			break
		case Reg(OpCode__TRAP):
			break
		case Reg(OpCode__RES), Reg(OpCode__RTI):
			break
		default:
			break
		}

		REG[Reg__PC] += 1
	}
}
