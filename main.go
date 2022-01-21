package main

import "fmt"

/*
x1 * x2 -> x8 = x1 loops of adding x2 to x8
*/
const src1 = `
# test
jal x0 4 # jump next word
i32 1000 # data at word 1
lw x1 x0 4 # load 0+4

# setup for loop
set x1 10
set x2 5
set x8 0
slti x3 x1 1
bne x3 x0 12 # 3ops*4bytes
addi x1 x1 -1
add x8 x8 x2
jal x0 -20 # 5ops*4bytes
set x7 400
sw x7 0 x8

# setup for multiply
set x1 10
set x2 5
mlt x8 x1 x2
set x7 400
sw x7 4 x8

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
