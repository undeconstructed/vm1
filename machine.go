package main

import "fmt"

const MemSize = 1000

type word uint32

type opcode uint8
type regist uint8

type op func(word)

const (
	OpNop opcode = iota
	OpHlt
	OpPut // reg -> mem
	OpGet // mem -> reg
	OpSet // n -> reg
	OpAdd // reg, reg -> reg
	OpMlt // reg, reg -> reg
	OpMov // reg -> reg
	OpJmp
	OpBr0
	OpFoo
	OpNull
)

const (
	RegPC regist = iota
	RegStat
	RegG0
	RegG1
	RegG2
	RegG3
	RegG4
	RegG5
	RegG6
	RegG7
	RegNull
)

const (
	StatEmpty word = 0
	StatZero       = 1 << 0
)

type machine struct {
	ops       [OpNull]op
	registers [RegNull]word
	memory    []word
}

func newMachine() *machine {
	vm := &machine{}
	vm.ops = [OpNull]op{
		vm._nop,
		vm._hlt,
		vm._put,
		vm._get,
		vm._set,
		vm._add,
		vm._mlt,
		vm._mov,
		vm._jmp,
		vm._br0,
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
	vm.registers[r] = n
	vm.setFlags(n)
}

func (vm *machine) setFlags(n word) {
	stat := StatEmpty
	if n == 0 {
		stat |= StatZero
	}
	vm.registers[RegStat] = stat
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

func (vm *machine) _nop(word) {
	fmt.Println("nop")
}

func (vm *machine) _hlt(word) {
	fmt.Println("halt")
}

func (vm *machine) _put(i word) {
	val, bas, at := readPut(i)
	fmt.Printf("put %d %d %d\n", val, bas, at)
	addr := vm.registers[bas] + vm.registers[at]
	vm.setMemory(addr, vm.registers[val])
}

func (vm *machine) _get(i word) {
	val, bas, at := readPut(i)
	fmt.Printf("get %d %d %d\n", val, bas, at)
	addr := vm.registers[bas] + vm.registers[at]
	n := vm.memory[addr]
	vm.setRegister(val, n)
}

func (vm *machine) _set(i word) {
	r, n := readSet(i)
	fmt.Printf("set %d %d\n", r, n)
	vm.setRegister(r, word(n))
}

func (vm *machine) _add(i word) {
	a, b, res := readAdd(i)
	fmt.Printf("add %d %d %d\n", a, b, res)
	n := vm.registers[a] + vm.registers[b]
	vm.setRegister(res, n)
}

func (vm *machine) _mlt(i word) {
	a, b, res := readAdd(i)
	fmt.Printf("mlt %d %d %d\n", a, b, res)
	n := vm.registers[b] * vm.registers[b]
	vm.setRegister(res, n)
}

func (vm *machine) _mov(i word) {
	a, b := readMov(i)
	fmt.Printf("mlt %d %d\n", a, b)
	n := vm.registers[a]
	vm.setRegister(b, n)
}

func (vm *machine) _jmp(i word) {
	n := readJmp(i)
	fmt.Printf("jmp %d\n", n)
	pc := int32(vm.registers[RegPC])
	npc := pc + int32(n)
	vm.registers[RegPC] = word(npc)
}

func (vm *machine) _br0(i word) {
	n := readBr0(i)
	fmt.Printf("br0 %d\n", n)
	stat := vm.registers[RegStat]
	if stat&StatZero == StatZero {
		pc := int32(vm.registers[RegPC])
		npc := pc + int32(n)
		vm.registers[RegPC] = word(npc)
	}
}

func (vm *machine) _foo(word) {
	fmt.Println("foo")
}
