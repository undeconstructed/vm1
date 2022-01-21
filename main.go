package main

import "fmt"

/*
x1 * x2 -> x8 = x1 loops of adding x2 to x8
*/
const src1 = `
# setup for loop
set x1 10
set x2 5
set x8 0
slti x3 x1 1
bne x3 x0 3
addi x1 x1 -1
add x8 x8 x2
jmp -5
set x7 100
put x8 x0 x7

# setup for multiply
set x1 10
set x2 5
mlt x8 x1 x2
set x7 101
put x8 x0 x7

# nonsense
foo
hlt
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
