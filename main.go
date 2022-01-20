package main

import "fmt"

const src1 = `
set g2 5
set g1 -1
set g0 5
br0 3
add g2 g3 g3
add g0 g1 g0
jmp -4
set g0 100
set g1 0
put g3 g0 g1
set g2 5
mlt g2 g2 g3
set g1 +1
put g3 g0 g1
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
