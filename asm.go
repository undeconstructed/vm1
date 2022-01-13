package main

import (
	"errors"
	"strconv"
	"strings"
)

func makeNop() word {
	return OpNop << 20
}

func makeHalt() word {
	return OpHalt << 20
}

func makeSet(n int16) word {
	return OpSet<<20 | word(uint16(n))
}

func makeAdd(n int16) word {
	return OpAdd<<20 | word(uint16(n))
}

func makeJmp(n int16) word {
	return OpJmp<<20 | word(uint16(n))
}

func makeBr0(n int16) word {
	return OpBr0<<20 | word(uint16(n))
}

func makeFoo() word {
	return OpFoo << 20
}

func assemble(src string) ([]word, error) {
	r := strings.Split(src, "\n")
	out := make([]word, 0, len(r))
	for _, l := range r {
		if l == "" {
			continue
		}

		ws := strings.Split(l, " ")

		var i word
		switch ws[0] {
		case "nop":
			i = makeNop()
		case "halt":
			i = makeHalt()
		case "set":
			n, err := strconv.Atoi(ws[1])
			if err != nil {
				return nil, err
			}
			i = makeSet(int16(n))
		case "add":
			n, err := strconv.Atoi(ws[1])
			if err != nil {
				return nil, err
			}
			i = makeAdd(int16(n))
		case "jmp":
			n, err := strconv.Atoi(ws[1])
			if err != nil {
				return nil, err
			}
			i = makeJmp(int16(n))
		case "br0":
			n, err := strconv.Atoi(ws[1])
			if err != nil {
				return nil, err
			}
			i = makeBr0(int16(n))
		case "foo":
			i = makeFoo()
		default:
			return nil, errors.New("unknown op: " + ws[0])
		}
		out = append(out, i)
	}
	return out, nil
}
