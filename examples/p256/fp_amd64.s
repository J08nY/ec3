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
	ADDQ    DI, DX
	ADCQ    R8, BX
	ADCQ    R9, BP
	ADCQ    CX, SI
	ADCQ    $0x00000000, R10
	MOVQ    DX, CX
	MOVQ    BX, DI
	MOVQ    BP, R8
	MOVQ    SI, R9
	SUBQ    p<>+0(SB), CX
	SBBQ    p<>+8(SB), DI
	SBBQ    p<>+16(SB), R8
	SBBQ    p<>+24(SB), R9
	SBBQ    $0x00000000, R10
	CMOVQCC CX, DX
	CMOVQCC DI, BX
	CMOVQCC R8, BP
	CMOVQCC R9, SI
	MOVQ    DX, (AX)
	MOVQ    BX, 8(AX)
	MOVQ    BP, 16(AX)
	MOVQ    SI, 24(AX)
	RET

DATA p<>+0(SB)/8, $0xffffffffffffffff
DATA p<>+8(SB)/8, $0x00000000ffffffff
DATA p<>+16(SB)/8, $0x0000000000000000
DATA p<>+24(SB)/8, $0xffffffff00000001
GLOBL p<>(SB), RODATA|NOPTR, $32

// func Mul(z *Elt, x *Elt, y *Elt)
TEXT ·Mul(SB), NOSPLIT, $64-24
	MOVQ z+0(FP), AX
	MOVQ x+8(FP), AX
	MOVQ y+16(FP), CX

	// y[0]
	MOVQ (CX), DX
	XORQ BX, BX

	// x[0] * y[0] -> z[0]
	MULXQ (AX), BP, SI

	// x[1] * y[0] -> z[1]
	MULXQ 8(AX), DI, R8
	ADCXQ DI, SI

	// x[2] * y[0] -> z[2]
	MULXQ 16(AX), DI, R9
	ADCXQ DI, R8

	// x[3] * y[0] -> z[3]
	MULXQ 24(AX), DX, DI
	ADCXQ DX, R9
	ADCXQ BX, DI
	MOVQ  BP, (SP)

	// y[1]
	MOVQ 8(CX), DX
	XORQ BX, BX

	// x[0] * y[1] -> z[1]
	MULXQ (AX), BP, R10
	ADCXQ BP, SI
	ADOXQ R10, R8

	// x[1] * y[1] -> z[2]
	MULXQ 8(AX), BP, R10
	ADCXQ BP, R8
	ADOXQ R10, R9

	// x[2] * y[1] -> z[3]
	MULXQ 16(AX), BP, R10
	ADCXQ BP, R9
	ADOXQ R10, DI

	// x[3] * y[1] -> z[4]
	MULXQ 24(AX), DX, BP
	ADCXQ DX, DI
	ADCXQ BX, BP
	ADOXQ BX, BP
	MOVQ  SI, 8(SP)

	// y[2]
	MOVQ 16(CX), DX
	XORQ BX, BX

	// x[0] * y[2] -> z[2]
	MULXQ (AX), SI, R10
	ADCXQ SI, R8
	ADOXQ R10, R9

	// x[1] * y[2] -> z[3]
	MULXQ 8(AX), SI, R10
	ADCXQ SI, R9
	ADOXQ R10, DI

	// x[2] * y[2] -> z[4]
	MULXQ 16(AX), SI, R10
	ADCXQ SI, DI
	ADOXQ R10, BP

	// x[3] * y[2] -> z[5]
	MULXQ 24(AX), DX, SI
	ADCXQ DX, BP
	ADCXQ BX, SI
	ADOXQ BX, SI
	MOVQ  R8, 16(SP)

	// y[3]
	MOVQ 24(CX), DX
	XORQ BX, BX

	// x[0] * y[3] -> z[3]
	MULXQ (AX), CX, R8
	ADCXQ CX, R9
	ADOXQ R8, DI

	// x[1] * y[3] -> z[4]
	MULXQ 8(AX), CX, R8
	ADCXQ CX, DI
	ADOXQ R8, BP

	// x[2] * y[3] -> z[5]
	MULXQ 16(AX), CX, R8
	ADCXQ CX, BP
	ADOXQ R8, SI

	// x[3] * y[3] -> z[6]
	MULXQ 24(AX), AX, CX
	ADCXQ AX, SI
	ADCXQ BX, CX
	ADOXQ BX, CX
	MOVQ  R9, 24(SP)
	MOVQ  DI, 32(SP)
	MOVQ  BP, 40(SP)
	MOVQ  SI, 48(SP)
	MOVQ  CX, 56(SP)

	// Reduction.
	RET
