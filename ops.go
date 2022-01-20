package main

// constants for instruction packing
const (
	PackWordLen = 32
	PackOpLen   = 8
	PackOpShift = PackWordLen - PackOpLen
	PackRegLen  = 8
)

func packOp(op opcode) word {
	return word(op) << PackOpShift
}

func unpackOp(from word) opcode {
	return opcode(from >> PackOpShift)
}

func packReg(into word, at uint, reg regist) word {
	return into | word(reg)<<(PackWordLen-at-PackRegLen)
}

func unpackReg(from word, at uint) regist {
	return regist(from >> (PackWordLen - at - PackRegLen))
}

func packInt16(into word, at uint, n int16) word {
	nw := word(uint16(n))
	return into | nw<<(PackWordLen-at-16)
}

func unpackInt16(from word, at uint) int16 {
	return int16(from >> (PackWordLen - at - 16))
}

func makeNop() word {
	return packOp(OpNop)
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

func makeSet(r regist, n int16) word {
	i := packOp(OpSet)
	i = packReg(i, 8, r)
	i = packInt16(i, 16, n)
	return i
}

func readSet(op word) (r regist, n int16) {
	return unpackReg(op, 8), unpackInt16(op, 16)
}

func makeAdd(a, b, res regist) word {
	i := packOp(OpAdd)
	i = packReg(i, 8, a)
	i = packReg(i, 16, b)
	i = packReg(i, 24, res)
	return i
}

func readAdd(op word) (a, b, res regist) {
	return unpackReg(op, 8), unpackReg(op, 16), unpackReg(op, 24)
}

func makeMlt(a, b, res regist) word {
	i := packOp(OpMlt)
	i = packReg(i, 8, a)
	i = packReg(i, 16, b)
	i = packReg(i, 24, res)
	return i
}

func readMlt(op word) (a, b, res regist) {
	return unpackReg(op, 8), unpackReg(op, 16), unpackReg(op, 24)
}

func makeMov(a, b regist) word {
	i := packOp(OpMov)
	i = packReg(i, 8, a)
	i = packReg(i, 16, b)
	return i
}

func readMov(op word) (a, b regist) {
	return unpackReg(op, 8), unpackReg(op, 16)
}

func makeJmp(n int16) word {
	i := packOp(OpJmp)
	i = packInt16(i, 8, n)
	return i
}

func readJmp(op word) (n int16) {
	return unpackInt16(op, 8)
}

func makeBr0(n int16) word {
	i := packOp(OpBr0)
	i = packInt16(i, 8, n)
	return i
}

func readBr0(op word) (n int16) {
	return unpackInt16(op, 8)
}

func makeFoo() word {
	return packOp(OpFoo)
}
