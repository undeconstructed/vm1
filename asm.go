package main

import (
	"errors"
	"strconv"
	"strings"
)

func parseReg(s string) (regist, error) {
	if strings.HasPrefix(s, "x") {
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			panic(err.Error())
		}
		n1 := regist(n)
		if n1 < 0 || n1 >= RegPC {
			panic("invalid register #")
		}
		return n1, nil
	}
	panic("invalid register")
}

func parseImm12(s string) (imm12, error) {
	i, err := strconv.ParseInt(s, 10, 12)
	if err != nil {
		panic(err.Error())
	}
	return imm12(i), nil
}

func parseImm20(s string) (imm20, error) {
	i, err := strconv.ParseInt(s, 10, 20)
	if err != nil {
		panic(err.Error())
	}
	return imm20(i), nil
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
			i = makeAddi(regist(0), regist(0), imm12(0))
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
			rd, _ := parseReg(ws[1])
			n, _ := parseImm12(ws[2])
			i = makeAddi(rd, RegG0, n)
		case "add":
			rd, _ := parseReg(ws[1])
			rs1, _ := parseReg(ws[2])
			rs2, _ := parseReg(ws[3])
			i = makeAdd(rd, rs1, rs2)
		case "mlt":
			rd, _ := parseReg(ws[1])
			rs1, _ := parseReg(ws[2])
			rs2, _ := parseReg(ws[3])
			i = makeMlt(rd, rs1, rs2)
		case "mov":
			rd, _ := parseReg(ws[1])
			rs, _ := parseReg(ws[2])
			i = makeAddi(rd, rs, imm12(0))
		case "addi":
			rd, _ := parseReg(ws[1])
			rs, _ := parseReg(ws[2])
			n, _ := parseImm12(ws[3])
			i = makeAddi(rd, rs, n)
		case "slti":
			rd, _ := parseReg(ws[1])
			rs, _ := parseReg(ws[2])
			n, _ := parseImm12(ws[3])
			i = makeSlti(rd, rs, n)
		case "jal":
			rd, _ := parseReg(ws[1])
			n, _ := parseImm20(ws[2])
			i = makeJal(rd, n)
		case "bne":
			rs1, _ := parseReg(ws[1])
			rs2, _ := parseReg(ws[2])
			n, _ := parseImm12(ws[3])
			i = makeBne(rs1, rs2, n)
		case "lui":
			rd, _ := parseReg(ws[1])
			n, _ := parseImm20(ws[2])
			i = makeLui(rd, n)
		case "foo":
			i = makeFoo()
		case "#":
			// comment
			continue
		case "":
			// blank line
			continue
		default:
			return nil, errors.New("unknown op: " + ws[0])
		}
		out = append(out, i)
	}
	return out, nil
}
