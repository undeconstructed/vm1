package main

// constants for instruction packing
const (
	PackWordLen   = 32
	PackOpLen     = 7
	PackOpMask    = 0b1111111
	PackOpShift   = PackWordLen - PackOpLen
	PackRegLen    = 5
	PackRegMask   = 0b11111
	PackF3Mask    = 0b111
	PackF7Mask    = 0b1111111
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

func packF3(into word, at uint, f3 funct3) word {
	return into | (word(f3&PackF3Mask) << at)
}

func unpackF3(from word, at uint) funct3 {
	return funct3((from >> at) & PackF3Mask)
}

func packF7(into word, at uint, f7 funct7) word {
	return into | (word(f7&PackF7Mask) << at)
}

func unpackF7(from word, at uint) funct7 {
	return funct7((from >> at) & PackF7Mask)
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

func encodeRType(opc opcode, f3 funct3, f7 funct7, rd, rs1, rs2 regist) word {
	i := packOp(opc)
	i = packReg(i, 7, rd)
	i = packF3(i, 12, f3)
	i = packReg(i, 15, rs1)
	i = packReg(i, 20, rs2)
	i = packF7(i, 25, f7)
	return i
}

func decodeRType(op word) (opc opcode, f3 funct3, f7 funct7, rd, rs1, rs2 regist) {
	return unpackOp(op), unpackF3(op, 12), unpackF7(op, 25), unpackReg(op, 7), unpackReg(op, 15), unpackReg(op, 20)
}

func encodeIType(opc opcode, f3 funct3, rd, rs1 regist, n imm12) word {
	i := packOp(opc)
	i = packReg(i, 7, rd)
	i = packF3(i, 12, f3)
	i = packReg(i, 15, rs1)
	i = packImm12(i, 20, n)
	return i
}

func decodeIType(op word) (opc opcode, f3 funct3, rd, rs1 regist, n imm12) {
	return unpackOp(op), unpackF3(op, 12), unpackReg(op, 7), unpackReg(op, 15), unpackImm12(op, 20)
}

func encodeSType(opc opcode, f3 funct3, rs1, rs2 regist, n imm12) word {
	// TODO - spec
	i := packOp(opc)
	i = packReg(i, 7, rs1)
	i = packReg(i, 12, rs2)
	i = packF3(i, 17, f3)
	i = packImm12(i, 20, n)
	return i
}

func decodeSType(op word) (opc opcode, f3 funct3, rs1, rs2 regist, n imm12) {
	return unpackOp(op), unpackF3(op, 17), unpackReg(op, 7), unpackReg(op, 12), unpackImm12(op, 20)
}

func encodeBType(opc opcode, f3 funct3, rs1, rs2 regist, n imm12) word {
	// TODO - spec
	return encodeSType(opc, f3, rs1, rs2, n)
}

func decodeBType(op word) (opc opcode, f3 funct3, rs1, rs2 regist, n imm12) {
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

func encodeJType(opc opcode, rd regist, n imm20) word {
	// TODO - spec
	return encodeUType(opc, rd, n)
}

func decodeJType(op word) (opc opcode, rd regist, n imm20) {
	return decodeUType(op)
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
	return encodeRType(OpOp, Funct3Add, Funct7Zeros, rd, rs1, rs2)
}

func readAdd(op word) (rd, rs1, rs2 regist) {
	_, _, _, rd, rs1, rs2 = decodeRType(op)
	return rd, rs1, rs2
}

func makeMlt(rd, rs1, rs2 regist) word {
	return encodeRType(OpMlt, 0, 0, rd, rs1, rs2)
}

func readMlt(op word) (rd, rs1, rs2 regist) {
	_, _, _, rd, rs1, rs2 = decodeRType(op)
	return rd, rs1, rs2
}

func makeAddi(rd, rs1 regist, n imm12) word {
	return encodeIType(OpImm, Funct3Addi, rd, rs1, n)
}

// func readAddi(op word) (rd, rs regist, n imm12) {
// 	_, _, rd, rs, n = decodeIType(op)
// 	return rd, rs, n
// }

func makeSlti(rd, rs1 regist, n imm12) word {
	return encodeIType(OpImm, Funct3Slti, rd, rs1, n)
}

// func readSlti(op word) (rd, rs regist, n imm12) {
// 	_, _, rd, rs, n = decodeIType(op)
// 	return rd, rs, n
// }

func makeJal(rd regist, n imm20) word {
	return encodeJType(OpJal, rd, n)
}

func readJal(op word) (rd regist, n imm20) {
	_, rd, n = decodeJType(op)
	return rd, n
}

func makeBne(rs1, rs2 regist, n imm12) word {
	return encodeBType(OpBranch, Funct3Bne, rs1, rs2, n)
}

func readBne(op word) (rs1, rs2 regist, n imm12) {
	_, _, rs1, rs2, n = decodeBType(op)
	return rs1, rs2, n
}

func makeFoo() word {
	return packOp(OpFoo)
}

func makeLui(rd regist, n imm20) word {
	return encodeUType(OpLui, rd, n)
}

func readLui(op word) (rd regist, n imm20) {
	_, rd, n = decodeUType(op)
	return rd, n
}

func makeLw(rd regist, rs1 regist, n imm12) word {
	return encodeIType(OpLoad, Funct3Lw, rd, rs1, n)
}

func readLw(op word) (rd regist, rs1 regist, n imm12) {
	_, _, rd, rs1, n = decodeIType(op)
	return rd, rs1, n
}

func makeSw(rs1 regist, n imm12, rs2 regist) word {
	return encodeSType(OpStore, Funct3Sw, rs1, rs2, n)
}

// func readSw(op word) (rs1 regist, n imm12, rs2 regist) {
// 	_, _, rs1, rs2, n = decodeSType(op)
// 	return rs1, n, rs2
// }
