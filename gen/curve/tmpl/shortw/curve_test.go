// CodeGenerationWarning

package shortw

import (
	"math/big"
	"testing"
)

var (
	cur = CURVENAME()
	ref = curvename.CurveParams
)

func TestCurveAdd(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		ex, ey := ref.Add(x1, y1, x2, y2)
		gx, gy := cur.Add(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveDouble(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)

		ex, ey := cur.Double(x1, y1)
		gx, gy := ref.Double(x1, y1)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarMult(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)
		x1, y1 := RandPoint(t)

		gx, gy := cur.ScalarMult(x1, y1, k.Bytes())
		ex, ey := ref.ScalarMult(x1, y1, k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarBaseMult(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		gx, gy := cur.ScalarBaseMult(k.Bytes())
		ex, ey := ref.ScalarBaseMult(k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveInverse(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		got := cur.Inverse(k)

		expect := new(big.Int).Set(k)
		expect.ModInverse(expect, ref.N)

		EqualInt(t, "inv", expect, got)
	}
}
