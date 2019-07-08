// Code generated by ec3. DO NOT EDIT.

#include "textflag.h"

// func Add(x *Elt, y *Elt)
TEXT ·Add(SB), NOSPLIT, $0-16
	MOVQ    x+0(FP), AX
	MOVQ    y+8(FP), CX
	MOVQ    (AX), DX
	MOVQ    8(AX), BX
	MOVQ    16(AX), BP
	MOVQ    24(AX), SI
	MOVQ    (CX), DI
	MOVQ    8(CX), R8
	MOVQ    16(CX), R9
	MOVQ    24(CX), CX
	XORQ    R10, R10
	MOVQ    $0x00000026, R11
	ADDQ    DI, DX
	ADCXQ   R8, BX
	ADCXQ   R9, BP
	ADCXQ   CX, SI
	MOVQ    R10, CX
	CMOVQCS R11, CX
	ADDQ    CX, DX
	ADCXQ   R10, BX
	ADCXQ   R10, BP
	ADCXQ   R10, SI
	MOVQ    R10, CX
	CMOVQCS R11, CX
	ADDQ    CX, DX
	MOVQ    DX, (AX)
	MOVQ    BX, 8(AX)
	MOVQ    BP, 16(AX)
	MOVQ    SI, 24(AX)
	RET

// func Mul(z *Elt, x *Elt, y *Elt)
TEXT ·Mul(SB), NOSPLIT, $64-24
	MOVQ z+0(FP), AX
	MOVQ x+8(FP), CX
	MOVQ y+16(FP), BX

	// y[0]
	MOVQ (BX), DX
	XORQ BP, BP

	// x[0] * y[0] -> z[0]
	MULXQ (CX), SI, DI

	// x[1] * y[0] -> z[1]
	MULXQ 8(CX), R8, R9
	ADCXQ R8, DI

	// x[2] * y[0] -> z[2]
	MULXQ 16(CX), R8, R12
	ADCXQ R8, R9

	// x[3] * y[0] -> z[3]
	MULXQ 24(CX), DX, R8
	ADCXQ DX, R12
	ADCXQ BP, R8
	MOVQ  SI, (SP)

	// y[1]
	MOVQ 8(BX), DX
	XORQ BP, BP

	// x[0] * y[1] -> z[1]
	MULXQ (CX), SI, R13
	ADCXQ SI, DI
	ADOXQ R13, R9

	// x[1] * y[1] -> z[2]
	MULXQ 8(CX), SI, R13
	ADCXQ SI, R9
	ADOXQ R13, R12

	// x[2] * y[1] -> z[3]
	MULXQ 16(CX), SI, R13
	ADCXQ SI, R12
	ADOXQ R13, R8

	// x[3] * y[1] -> z[4]
	MULXQ 24(CX), DX, SI
	ADCXQ DX, R8
	ADCXQ BP, SI
	ADOXQ BP, SI
	MOVQ  DI, 8(SP)

	// y[2]
	MOVQ 16(BX), DX
	XORQ BP, BP

	// x[0] * y[2] -> z[2]
	MULXQ (CX), DI, R13
	ADCXQ DI, R9
	ADOXQ R13, R12

	// x[1] * y[2] -> z[3]
	MULXQ 8(CX), DI, R13
	ADCXQ DI, R12
	ADOXQ R13, R8

	// x[2] * y[2] -> z[4]
	MULXQ 16(CX), DI, R13
	ADCXQ DI, R8
	ADOXQ R13, SI

	// x[3] * y[2] -> z[5]
	MULXQ 24(CX), DX, DI
	ADCXQ DX, SI
	ADCXQ BP, DI
	ADOXQ BP, DI
	MOVQ  R9, 16(SP)

	// y[3]
	MOVQ 24(BX), DX
	XORQ BP, BP

	// x[0] * y[3] -> z[3]
	MULXQ (CX), BX, R9
	ADCXQ BX, R12
	ADOXQ R9, R8

	// x[1] * y[3] -> z[4]
	MULXQ 8(CX), BX, R9
	ADCXQ BX, R8
	ADOXQ R9, SI

	// x[2] * y[3] -> z[5]
	MULXQ 16(CX), BX, R9
	ADCXQ BX, SI
	ADOXQ R9, DI

	// x[3] * y[3] -> z[6]
	MULXQ 24(CX), CX, DX
	ADCXQ CX, DI
	ADCXQ BP, DX
	ADOXQ BP, DX
	MOVQ  R12, 24(SP)
	MOVQ  R8, 32(SP)
	MOVQ  SI, 40(SP)
	MOVQ  DI, 48(SP)
	MOVQ  DX, 56(SP)

	// Reduction.
	XORQ    R10, R10
	MOVQ    $0x00000026, CX
	MOVQ    CX, DX
	XORQ    R11, R11
	MULXQ   32(SP), BP, BX
	ADCXQ   BP, R11
	MULXQ   40(SP), SI, BP
	ADCXQ   SI, BX
	MULXQ   48(SP), DI, SI
	ADCXQ   DI, BP
	MULXQ   56(SP), DI, DX
	ADCXQ   DI, SI
	ADOXQ   (SP), R11
	ADOXQ   8(SP), BX
	ADOXQ   16(SP), BP
	ADOXQ   24(SP), SI
	ADOXQ   R10, DX
	IMULQ   CX, DX
	ADCXQ   DX, R11
	ADCXQ   R10, BX
	ADCXQ   R10, BP
	ADCXQ   R10, SI
	MOVQ    R10, DX
	CMOVQCS CX, DX
	ADDQ    DX, R11
	MOVQ    R11, (AX)
	MOVQ    BX, 8(AX)
	MOVQ    BP, 16(AX)
	MOVQ    SI, 24(AX)
	RET