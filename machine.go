package main

import "fmt"

const MemSize = 1000

type word uint32

type opcode uint8
type regist uint8
type funct3 uint8
type funct7 uint8
type imm12 int16
type imm20 int32

type op func(word)

const (
	OpHlt opcode = iota
	OpPut        // reg -> mem
	OpGet        // mem -> reg
	OpMlt        // reg, reg -> reg
	OpImm
	OpBranch
	OpLui
	// OpAuipc
	OpOp
	OpJal
	OpFoo
	OpNull
)

const (
	Funct3Addi  funct3 = 0
	Funct3Slti         = 1
	Funct3Sltiu        = 2

	Funct3Andi = 0
	Funct3Ori  = 1
	Funct3Xori = 2

	Funct3Slli = 0
	Funct3Srli = 1
	Funct3Srai = 2

	Funct3Add  = 0
	Funct3Slt  = 1
	Funct3Sltu = 2

	Funct3And = 0
	Funct3Or  = 1
	Funct3Xor = 2

	Funct3Sll = 0
	Funct3Srl = 1

	Funct3Sub = 0
	Funct3Sra = 1

	Funct3Beq  = 0
	Funct3Bne  = 1
	Funct3Blt  = 2
	Funct3Bltu = 3
	Funct3Bge  = 4
	Funct3Bgeu = 5
)

const (
	Funct7Zeros funct7 = 0
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
		vm._mlt,
		vm._imm,
		vm._branch,
		vm._lui,
		vm._op,
		vm._jal,
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

func (vm *machine) _op(i word) {
	_, f3, f7, rd, rs1, rs2 := decodeRType(i)
	fmt.Printf("opop f3=%d f7=%d\n", f3, f7)
	switch f7 {
	case Funct7Zeros:
		switch f3 {
		case Funct3Add:
			vm._add(rd, rs1, rs2)
		}
	}
}

func (vm *machine) _add(rd, rs1, rs2 regist) {
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
	fmt.Printf("opimm f3=%d\n", f3)
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

func (vm *machine) _jal(i word) {
	// TODO - program counter should be bytes not words
	rd, n := readJal(i)
	fmt.Printf("jal x%d %d\n", rd, n)
	pc := int32(vm.registers[RegPC])
	pc1 := pc + 1
	pc2 := pc + int32(n)
	vm.registers[RegPC] = word(pc2)
	vm.setRegister(rd, word(pc1))
}

func (vm *machine) _branch(i word) {
	_, f3, rs1, rs2, n := decodeBType(i)
	fmt.Printf("opbranch f3=%d\n", f3)
	switch f3 {
	case Funct3Bne:
		vm._bne(rs1, rs2, n)
	}
}

func (vm *machine) _bne(rs1, rs2 regist, n imm12) {
	fmt.Printf("bne x%d x%d %d\n", rs1, rs2, n)
	if vm.registers[rs1] != vm.registers[rs2] {
		pc := int32(vm.registers[RegPC])
		npc := pc + int32(n)
		vm.registers[RegPC] = word(npc)
	}
}

func (vm *machine) _lui(i word) {
	rd, n := readLui(i)
	fmt.Printf("lui x%d %d\n", rd, n)
	n1 := word(n) << 12
	vm.setRegister(rd, n1)
}

func (vm *machine) _foo(word) {
	fmt.Println("foo")
}
