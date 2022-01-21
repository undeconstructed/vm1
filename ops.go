package main

// constants for instruction packing
const (
	PackWordLen   = 32
	PackOpLen     = 7
	PackOpMask    = 0b1111111
	PackOpShift   = PackWordLen - PackOpLen
	PackRegLen    = 5
	PackRegMask   = 0b11111
	PackImm12Len  = 12
	PackImm12Mask = 0b111111111111
	PackImm20Len  = 20
	PackImm20Mask = 0b11111111111111111111
)

func packOp(op opcode) word {
	if op >= OpNull {
		panic("illegal opcode")
	}
	return word(op)
}

func unpackOp(from word) opcode {
	return opcode(from) & PackOpMask
}

func packReg(into word, at uint, reg regist) word {
	if reg >= RegNull {
		panic("illegal register")
	}
	return into | (word(reg&PackRegMask) << at)
}

func unpackReg(from word, at uint) regist {
	return regist((from >> at) & PackRegMask)
}

func packImm12(into word, at uint, n imm12) word {
	return into | (word(n&PackImm12Mask) << at)
}

func unpackImm12(from word, at uint) imm12 {
	ui := (from >> at) & PackImm12Mask
	if ui>>(PackImm12Len-1) == 1 {
		ui = ui | (0b1111 << PackImm12Len)
	}
	return imm12(ui)
}

func packImm20(into word, at uint, n imm20) word {
	return into | (word(n&PackImm20Mask) << at)
}

func unpackImm20(from word, at uint) imm20 {
	ui := (from >> at) & PackImm20Mask
	if ui>>(PackImm20Len-1) == 1 {
		ui = ui | (0b111111111111 << PackImm20Len)
	}
	return imm20(ui)
}

func encodeRType(opc opcode, rd, rs1, rs2 regist) word {
	i := packOp(opc)
	i = packReg(i, 7, rd)
	i = packReg(i, 15, rs1)
	i = packReg(i, 20, rs2)
	return i
}

func decodeRType(op word) (opc opcode, rd, rs1, rs2 regist) {
	return unpackOp(op), unpackReg(op, 7), unpackReg(op, 15), unpackReg(op, 20)
}

func encodeIType(opc opcode, rd, rs1 regist, n imm12) word {
	i := packOp(opc)
	i = packReg(i, 7, rd)
	i = packReg(i, 15, rs1)
	i = packImm12(i, 20, n)
	return i
}

func decodeIType(op word) (opc opcode, rd, rs1 regist, n imm12) {
	return unpackOp(op), unpackReg(op, 7), unpackReg(op, 15), unpackImm12(op, 20)
}

func encodeSType(opc opcode, rs1, rs2 regist, n imm12) word {
	// TODO - spec
	i := packOp(opc)
	i = packReg(i, 7, rs1)
	i = packReg(i, 12, rs2)
	i = packImm12(i, 17, n)
	return i
}

func decodeSType(op word) (opc opcode, rs1, rs2 regist, n imm12) {
	return unpackOp(op), unpackReg(op, 7), unpackReg(op, 12), unpackImm12(op, 17)
}

func encodeBType(opc opcode, rs1, rs2 regist, n imm12) word {
	return encodeSType(opc, rs1, rs2, n)
}

func decodeBType(op word) (opc opcode, rs1, rs2 regist, n imm12) {
	return decodeSType(op)
}

func encodeUType(opc opcode, rd regist, n imm20) word {
	i := packOp(opc)
	i = packReg(i, 7, rd)
	i = packImm20(i, 12, n)
	return i
}

func decodeUType(op word) (opc opcode, rd regist, n imm20) {
	return unpackOp(op), unpackReg(op, 7), unpackImm20(op, 12)
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
	return encodeRType(OpAdd, rd, rs1, rs2)
}

func readAdd(op word) (rd, rs1, rs2 regist) {
	_, rd, rs1, rs2 = decodeRType(op)
	return rd, rs1, rs2
}

func makeMlt(rd, rs1, rs2 regist) word {
	return encodeRType(OpMlt, rd, rs1, rs2)
}

func readMlt(op word) (rd, rs1, rs2 regist) {
	_, rd, rs1, rs2 = decodeRType(op)
	return rd, rs1, rs2
}

func makeAddi(rd, rs1 regist, n imm12) word {
	return encodeIType(OpAddi, rd, rs1, n)
}

func readAddi(op word) (rd, rs regist, n imm12) {
	_, rd, rs, n = decodeIType(op)
	return rd, rs, n
}

func makeSlti(rd, rs1 regist, n imm12) word {
	return encodeIType(OpSlti, rd, rs1, n)
}

func readSlti(op word) (rd, rs regist, n imm12) {
	_, rd, rs, n = decodeIType(op)
	return rd, rs, n
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
	return encodeBType(OpBne, rs1, rs2, n)
}

func readBne(op word) (rs1, rs2 regist, n imm12) {
	_, rs1, rs2, n = decodeBType(op)
	return rs1, rs2, n
}

func makeFoo() word {
	return packOp(OpFoo)
}
