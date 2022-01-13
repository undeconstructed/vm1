package main

import "fmt"

const MemSize = 1000

type word uint32

type op func(word)

const (
	OpNop = iota
	OpHalt
	OpSet
	OpAdd
	OpJmp
	OpBr0
	OpFoo
	OpNull
)

const (
	RegPC = iota
	RegAcc
	RegStat
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
		vm._halt,
		vm._set,
		vm._add,
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

func (vm *machine) setFlags() {
	acc := vm.registers[RegAcc]
	stat := StatEmpty
	if acc == 0 {
		stat |= StatZero
	}
	vm.registers[RegStat] = stat
}

func (vm *machine) step() bool {
	pc := vm.registers[RegPC]
	op := vm.memory[pc]

	opCode := op >> 20
	if opCode == OpHalt {
		return false
	}

	opFunc := vm.ops[opCode]

	opFunc(op)

	vm.setFlags()

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
	fmt.Printf("%v\n", vm)
}

func (vm *machine) _nop(word) {
	fmt.Println("nop")
}

func (vm *machine) _halt(word) {
	fmt.Println("halt")
}

func (vm *machine) _set(i word) {
	v := int32(int16(i & ((1 << 16) - 1)))
	fmt.Printf("set %d\n", v)
	vm.registers[RegAcc] = word(v)
}

func (vm *machine) _add(i word) {
	v := int32(int16(i & ((1 << 16) - 1)))
	fmt.Printf("add %d\n", v)
	a := int32(vm.registers[RegAcc])
	nv := a + v
	vm.registers[RegAcc] = word(nv)
}

func (vm *machine) _jmp(i word) {
	v := int32(int16(i & ((1 << 16) - 1)))
	fmt.Printf("jmp %d\n", v)
	pc := int32(vm.registers[RegPC])
	npc := pc + v
	vm.registers[RegPC] = word(npc)
}

func (vm *machine) _br0(i word) {
	v := int32(int16(i & ((1 << 16) - 1)))
	fmt.Printf("br0 %d\n", v)
	stat := vm.registers[RegStat]
	if stat&StatZero == StatZero {
		pc := int32(vm.registers[RegPC])
		npc := pc + v
		vm.registers[RegPC] = word(npc)
	}
}

func (vm *machine) _foo(word) {
	fmt.Println("foo")
}
