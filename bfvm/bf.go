package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func main() {
	inputfile := "main.bf"
	input, err := os.ReadFile(inputfile)
	if err != nil {
		panic(err)
	}

	compiled := Compile(string(input))
	bufIn := bufio.NewReader(os.Stdin)
	bufOut := bufio.NewWriter(os.Stdout)
	defer bufOut.Flush()
	err = Run(compiled, bufIn, bufOut)
	if err != nil {
		panic(err)
	}
}

type Instruction struct {
	Opcode  byte
	Operand uint32
}

func Compile(code string) []Instruction {
	var instructions []Instruction
	var program_counter uint32
	var jump_Stack []uint32

	for {
		if len(code) == 0 {
			break
		}

		const TOKEN_INC_PTR = ">"
		const TOKEN_DEC_PTR = "<"
		const TOKEN_INC_VAL = "+"
		const TOKEN_DEC_VAL = "-"
		const TOKEN_OUTPUT = "."
		const TOKEN_INPUT = ","
		const TOKEN_JUMP_FWD = "["
		const TOKEN_JUMP_BACK = "]"

		switch {
		case strings.HasPrefix(code, TOKEN_INC_PTR):
			instructions = append(instructions, Instruction{'>', 0})
			code = code[len(TOKEN_INC_PTR):]
		case strings.HasPrefix(code, TOKEN_DEC_PTR):
			instructions = append(instructions, Instruction{'<', 0})
			code = code[len(TOKEN_DEC_PTR):]
		case strings.HasPrefix(code, TOKEN_INC_VAL):
			instructions = append(instructions, Instruction{'+', 0})
			code = code[len(TOKEN_INC_VAL):]
		case strings.HasPrefix(code, TOKEN_DEC_VAL):
			instructions = append(instructions, Instruction{'-', 0})
			code = code[len(TOKEN_DEC_VAL):]
		case strings.HasPrefix(code, TOKEN_OUTPUT):
			instructions = append(instructions, Instruction{'.', 0})
			code = code[len(TOKEN_OUTPUT):]
		case strings.HasPrefix(code, TOKEN_INPUT):
			instructions = append(instructions, Instruction{',', 0})
			code = code[len(TOKEN_INPUT):]
		case strings.HasPrefix(code, TOKEN_JUMP_FWD):
			code = code[len(TOKEN_JUMP_FWD):]
			instructions = append(instructions, Instruction{'[', 0})
			//println("[ at ", program_counter)
			jump_Stack = append(jump_Stack, program_counter)
		case strings.HasPrefix(code, TOKEN_JUMP_BACK):
			code = code[len(TOKEN_JUMP_BACK):]
			if len(jump_Stack) == 0 {
				panic("Unmatched ]")
			}
			jmp_pos := jump_Stack[len(jump_Stack)-1]
			jump_Stack = jump_Stack[:len(jump_Stack)-1]
			instructions = append(instructions, Instruction{']', jmp_pos})
			instructions[jmp_pos].Operand = program_counter
		default:
			// skip
			code = code[1:]
			program_counter--
		}

		program_counter++
	}
	return instructions
}

func Run(code []Instruction, r io.Reader, w io.Writer) error {
	var readBuffer []byte = make([]byte, 1)

	var memory []byte = make([]byte, 65536)
	var memory_pointer uint32 = 0
	var program_counter uint32 = 0

	for {
		if program_counter >= uint32(len(code)) {
			break
		}
		// /println("opcode: ", string([]byte{code[program_counter].Opcode}))

		switch code[program_counter].Opcode {
		case '>':
			memory_pointer++
			memory_pointer %= uint32(len(memory))
		case '<':
			memory_pointer--
			memory_pointer %= uint32(len(memory))
		case '+':
			memory[memory_pointer]++
		case '-':
			memory[memory_pointer]--
		case '.':
			w.Write([]byte{memory[memory_pointer]})
		case ',':
			_, err := r.Read(readBuffer)
			if err != nil {
				return err
			}
			memory[memory_pointer] = readBuffer[0]
		case '[':
			if memory[memory_pointer] == 0 {
				program_counter = code[program_counter].Operand
				continue
			}
		case ']':
			if memory[memory_pointer] != 0 {
				program_counter = code[program_counter].Operand
				continue
			}
		}

		program_counter++
	}

	return nil
}
