package main

import (
	"errors"
	"strconv"
	"strings"
)

func parseReg(s string) (regist, error) {
	if strings.HasPrefix(s, "g") {
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			panic(err.Error())
		}
		if n < 0 || n > 7 {
			panic("invalid register #")
		}
		return regist(n) + RegG0, nil
	}
	panic("invalid register")
}

func parseInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		panic(err.Error())
	}
	return int16(i), nil
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
		case "hlt":
			i = makeHlt()
		case "put":
			val, _ := parseReg(ws[1])
			bas, _ := parseReg(ws[2])
			at, _ := parseReg(ws[3])
			i = makePut(val, bas, at)
		case "get":
			val, _ := parseReg(ws[1])
			bas, _ := parseReg(ws[2])
			at, _ := parseReg(ws[3])
			i = makeGet(val, bas, at)
		case "set":
			r, _ := parseReg(ws[1])
			n, _ := parseInt16(ws[2])
			i = makeSet(r, n)
		case "add":
			a, _ := parseReg(ws[1])
			b, _ := parseReg(ws[2])
			r, _ := parseReg(ws[3])
			i = makeAdd(a, b, r)
		case "mlt":
			a, _ := parseReg(ws[1])
			b, _ := parseReg(ws[2])
			r, _ := parseReg(ws[3])
			i = makeMlt(a, b, r)
		case "mov":
			a, _ := parseReg(ws[1])
			b, _ := parseReg(ws[2])
			i = makeMov(a, b)
		case "jmp":
			n, _ := parseInt16(ws[1])
			i = makeJmp(int16(n))
		case "br0":
			n, _ := parseInt16(ws[1])
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
