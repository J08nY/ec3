package mp

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

// Int represents a multi-precision integer.
type Int []operand.Op

// NewInt builds an empty integer with k limbs.
func NewInt(k int) Int {
	return make(Int, k)
}

// NewIntLimb64 builds multi-precision integer with k 64-bit limbs.
func NewIntLimb64(ctx *build.Context, k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = ctx.GP64()
	}
	return x
}

// NewIntFromMem builds a multi-precision integer referencing the k 64-bit limbs
// at memory address m.
func NewIntFromMem(m operand.Mem, k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = m.Offset(8 * i)
	}
	return x
}

// Mul does a full multiply z = x*y.
func Mul(ctx *build.Context, z, x, y Int) {
	// TODO(mbm): multi-precision multiply is ugly

	acc := make([]operand.Op, len(z))
	zero := ctx.GP64()

	for j := 0; j < len(y); j++ {
		ctx.Commentf("y[%d]", j)
		ctx.MOVQ(y[j], reg.RDX)
		ctx.XORQ(zero, zero) // clears flags
		carryinto := [2]int{-1, -1}
		for i := 0; i < len(x); i++ {
			k := i + j
			ctx.Commentf("x[%d] * y[%d] -> z[%d]", i, j, k)

			// Determine where the results should go.
			var product [2]operand.Op
			var add [2]bool
			for b := 0; b < 2; b++ {
				if acc[k+b] == nil {
					acc[k+b] = ctx.GP64()
					product[b] = acc[k+b]
				} else {
					product[b] = ctx.GP64()
					add[b] = true
				}
			}

			// Do the multiply.
			ctx.MULXQ(x[i], product[0], product[1])

			// Do the adds.
			if add[0] {
				ctx.ADCXQ(product[0], acc[k])
				carryinto[0] = k + 1
			}
			if add[1] {
				ctx.ADOXQ(product[1], acc[k+1])
				carryinto[1] = k + 2
			}
		}

		if carryinto[0] > 0 {
			ctx.ADCXQ(zero, acc[carryinto[0]])
		}
		if carryinto[1] > 0 {
			ctx.ADOXQ(zero, acc[carryinto[1]])
		}

		//
		ctx.MOVQ(acc[j], z[j])
	}

	for j := len(y); j < len(z); j++ {
		ctx.MOVQ(acc[j], z[j])
	}
}