package main

import "testing"

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
	packed := packImm12(0, 10, in)
	t.Logf("packed: %b", packed)
	out := unpackImm12(packed, 10)
	if in != out {
		t.Errorf("%d != %d", in, out)
	}
}
