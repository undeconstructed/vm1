package main

import "fmt"

const src1 = `
nop
set 5
br0 2
add -1
jmp -3
foo
halt
`

func main() {
	prog, err := assemble(src1)

	if err != nil {
		fmt.Printf("assembly error: %v\n", err)
		return
	}

	vm := newMachine()

	vm.load(0, prog)

	vm.print()
	vm.run(100)
	vm.print()
}
