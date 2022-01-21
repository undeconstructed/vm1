package main

import "testing"

func TestOpCodePacking(t *testing.T) {
	in := OpLui
	packed := packOp(in)
	t.Logf("packed: %b", packed)
	out := unpackOp(packed)
	if in != out {
		t.Errorf("%d != %d", in, out)
	}
}

func TestRegPacking(t *testing.T) {
	in := regist(5)
	packed := packReg(0, 10, in)
	t.Logf("packed: %b", packed)
	out := unpackReg(packed, 10)
	if in != out {
		t.Errorf("%d != %d", in, out)
	}
}

func TestImm12Packing(t *testing.T) {
	in := imm12(-1)
	packed := packImm12(0, 16, in)
	t.Logf("packed: %b", packed)
	out := unpackImm12(packed, 16)
	if in != out {
		t.Errorf("%d != %d", in, out)
	}
}

func TestImm20Packing(t *testing.T) {
	in := imm20(-1)
	packed := packImm20(0, 7, in)
	t.Logf("packed: %b", packed)
	out := unpackImm20(packed, 7)
	if in != out {
		t.Errorf("%d != %d", in, out)
	}
}
