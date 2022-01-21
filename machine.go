package main

import "fmt"

const MemSize = 1000

type word uint32

type opcode uint8
type regist uint8
type funct3 uint8
type imm12 int16
type imm20 int32

type op func(word)

const (
	OpHlt opcode = iota
	OpPut        // reg -> mem
	OpGet        // mem -> reg
	OpAdd        // reg, reg -> reg
	OpMlt        // reg, reg -> reg
	OpImm
	OpJmp
	OpBne
	OpFoo
	OpNull
)

const (
	Funct3Addi funct3 = iota
	Funct3Slti
)

const (
	RegG0   regist = 0
	RegPC          = 32
	RegNull        = RegPC + 1
)

type machine struct {
	ops       [OpNull]op
	registers [RegNull]word
	memory    []word
}

func newMachine() *machine {
	vm := &machine{}
	vm.ops = [OpNull]op{
		vm._hlt,
		vm._put,
		vm._get,
		vm._add,
		vm._mlt,
		vm._imm,
		vm._jmp,
		vm._bne,
		vm._foo,
	}
	vm.registers = [RegNull]word{}
	vm.memory = make([]word, MemSize, MemSize)
	return vm
}

func (vm *machine) reset() {
}

func (vm *machine) load(at int, data []word) {
	for i, v := range data {
		vm.memory[at+i] = v
	}
}

func (vm *machine) set(n int, instruction word) {
	vm.memory[n] = instruction
}

func (vm *machine) setMemory(a word, n word) {
	vm.memory[a] = n
}

func (vm *machine) setRegister(r regist, n word) {
	if r == RegG0 {
		// always must be zero
		return
	}
	vm.registers[r] = n
	fmt.Printf("R: %v\n", vm.registers)
}

func (vm *machine) step() bool {
	pc := vm.registers[RegPC]
	op := vm.memory[pc]

	opCode := unpackOp(op)
	if opCode == OpHlt {
		return false
	}

	opFunc := vm.ops[opCode]

	opFunc(op)

	vm.registers[RegPC] += 1

	return true
}

func (vm *machine) run(n int) {
	for n > 0 {
		if !vm.step() {
			return
		}
		n--
	}
}

func (vm *machine) print() {
	fmt.Printf("R: %v\n", vm.registers)
	fmt.Printf("M: %v\n", vm.memory)
}

func (vm *machine) _hlt(word) {
	fmt.Println("halt")
}

func (vm *machine) _put(i word) {
	val, bas, at := readPut(i)
	fmt.Printf("put x%d x%d x%d\n", val, bas, at)
	addr := vm.registers[bas] + vm.registers[at]
	vm.setMemory(addr, vm.registers[val])
}

func (vm *machine) _get(i word) {
	val, bas, at := readPut(i)
	fmt.Printf("get x%d x%d x%d\n", val, bas, at)
	addr := vm.registers[bas] + vm.registers[at]
	n := vm.memory[addr]
	vm.setRegister(val, n)
}

func (vm *machine) _add(i word) {
	rd, rs1, rs2 := readAdd(i)
	fmt.Printf("add x%d x%d x%d\n", rd, rs1, rs2)
	n := vm.registers[rs1] + vm.registers[rs2]
	vm.setRegister(rd, n)
}

func (vm *machine) _mlt(i word) {
	rd, rs1, rs2 := readMlt(i)
	fmt.Printf("mlt x%d x%d x%d\n", rd, rs1, rs2)
	n := vm.registers[rs1] * vm.registers[rs2]
	vm.setRegister(rd, n)
}

func (vm *machine) _imm(i word) {
	_, f3, rd, rs, imm := decodeIType(i)
	switch f3 {
	case Funct3Addi:
		vm._addi(rd, rs, imm)
	case Funct3Slti:
		vm._slti(rd, rs, imm)
	}
}

func (vm *machine) _addi(rd, rs regist, v imm12) {
	fmt.Printf("addi x%d x%d %d\n", rd, rs, v)
	n := int32(vm.registers[rs]) + int32(v)
	vm.setRegister(rd, word(n))
}

func (vm *machine) _slti(rd, rs regist, n imm12) {
	fmt.Printf("slti x%d x%d %d\n", rd, rs, n)
	n0 := vm.registers[rs]
	flag := 0
	if int32(n0) < int32(n) {
		flag = 1
	}
	vm.setRegister(rd, word(flag))
}

func (vm *machine) _jmp(i word) {
	n := readJmp(i)
	fmt.Printf("jmp %d\n", n)
	pc := int32(vm.registers[RegPC])
	npc := pc + int32(n)
	vm.registers[RegPC] = word(npc)
}

func (vm *machine) _bne(i word) {
	rs1, rs2, n := readBne(i)
	fmt.Printf("bne x%d x%d %d\n", rs1, rs2, n)
	if vm.registers[rs1] != vm.registers[rs2] {
		pc := int32(vm.registers[RegPC])
		npc := pc + int32(n)
		vm.registers[RegPC] = word(npc)
	}
}

func (vm *machine) _foo(word) {
	fmt.Println("foo")
}
