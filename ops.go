package main

// constants for instruction packing
const (
	PackWordLen   = 32
	PackOpLen     = 6
	PackOpMask    = 0b111111
	PackOpShift   = PackWordLen - PackOpLen
	PackRegLen    = 5
	PackRegMask   = 0b11111
	PackImm12Len  = 12
	PackImm12Mask = 0b111111111111
)

func packOp(op opcode) word {
	if op >= OpNull {
		panic("illegal opcode")
	}
	return word(op) << PackOpShift
}

func unpackOp(from word) opcode {
	return opcode(from>>PackOpShift) & PackOpMask
}

func packReg(into word, at uint, reg regist) word {
	if reg >= RegNull {
		panic("illegal register")
	}
	return into | word(reg)<<((PackWordLen-at)-PackRegLen)
}

func unpackReg(from word, at uint) regist {
	return regist(from>>((PackWordLen-at)-PackRegLen)) & PackRegMask
}

func packImm12(into word, at uint, n imm12) word {
	nw := word(uint16(n)) & PackImm12Mask
	return into | (nw << ((PackWordLen - at) - PackImm12Len))
}

func unpackImm12(from word, at uint) imm12 {
	ui := (from >> ((PackWordLen - at) - PackImm12Len)) & PackImm12Mask
	if ui>>11 == 1 {
		ui = ui | (0b1111 << 12)
	}
	return imm12(ui)
}

func makeHlt() word {
	return packOp(OpHlt)
}

func makePut(val, bas, at regist) word {
	i := packOp(OpPut)
	i = packReg(i, 8, val)
	i = packReg(i, 16, bas)
	i = packReg(i, 24, at)
	return i
}

func readPut(op word) (val, bas, at regist) {
	return unpackReg(op, 8), unpackReg(op, 16), unpackReg(op, 24)
}

func makeGet(val, bas, at regist) word {
	i := packOp(OpPut)
	i = packReg(i, 8, val)
	i = packReg(i, 16, bas)
	i = packReg(i, 24, at)
	return i
}

func readGet(op word) (val, bas, at regist) {
	return unpackReg(op, 8), unpackReg(op, 16), unpackReg(op, 24)
}

func makeAdd(rd, rs1, rs2 regist) word {
	i := packOp(OpAdd)
	i = packReg(i, 6, rd)
	i = packReg(i, 11, rs1)
	i = packReg(i, 16, rs2)
	return i
}

func readAdd(op word) (rd, rs1, rs2 regist) {
	return unpackReg(op, 6), unpackReg(op, 11), unpackReg(op, 16)
}

func makeMlt(rd, rs1, rs2 regist) word {
	i := packOp(OpMlt)
	i = packReg(i, 6, rd)
	i = packReg(i, 11, rs1)
	i = packReg(i, 16, rs2)
	return i
}

func readMlt(op word) (rd, rs1, rs2 regist) {
	return unpackReg(op, 6), unpackReg(op, 11), unpackReg(op, 16)
}

func makeAddi(rd, rs1 regist, n imm12) word {
	i := packOp(OpAddi)
	i = packReg(i, 6, rd)
	i = packReg(i, 11, rs1)
	i = packImm12(i, 16, n)
	return i
}

func readAddi(op word) (rd, rs regist, i imm12) {
	return unpackReg(op, 6), unpackReg(op, 11), unpackImm12(op, 16)
}

func makeSlti(rd, rs1 regist, n imm12) word {
	i := packOp(OpSlti)
	i = packReg(i, 6, rd)
	i = packReg(i, 11, rs1)
	i = packImm12(i, 16, n)
	return i
}

func readSlti(op word) (rd, rs regist, i imm12) {
	return unpackReg(op, 6), unpackReg(op, 11), unpackImm12(op, 16)
}

func makeJmp(n imm12) word {
	i := packOp(OpJmp)
	i = packImm12(i, 8, n)
	return i
}

func readJmp(op word) (n imm12) {
	return unpackImm12(op, 8)
}

func makeBne(rs1, rs2 regist, n imm12) word {
	i := packOp(OpBne)
	i = packReg(i, 6, rs1)
	i = packReg(i, 11, rs2)
	i = packImm12(i, 16, n)
	return i
}

func readBne(op word) (rs1, rs2 regist, i imm12) {
	return unpackReg(op, 6), unpackReg(op, 11), unpackImm12(op, 16)
}

func makeFoo() word {
	return packOp(OpFoo)
}
