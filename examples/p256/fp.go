// Code generated by ec3. DO NOT EDIT.

package p256

import "math/big"

// Size is the size of a field element in bytes.
const Size = 32

// Elt is a field element.
type Elt [32]uint8

// p is the field prime modulus as a big integer.
var p, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)

// prime is the prime field modulus as a field element.
var prime = Elt{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
}

// SetInt64 constructs a field element from an integer.
func (x *Elt) SetInt64(y int64) *Elt {
	x.SetInt(big.NewInt(y))
	return x
}

// SetInt constructs a field element from a big integer.
func (x *Elt) SetInt(y *big.Int) *Elt {
	// Reduce if outside range.
	if y.Sign() < 0 || y.Cmp(p) >= 0 {
		y = new(big.Int).Mod(y, p)
	}
	// Copy bytes into field element.
	b := y.Bytes()
	i := 0
	for ; i < len(b); i++ {
		x[i] = b[len(b)-1-i]
	}
	for ; i < Size; i++ {
		x[i] = 0
	}
	// Encode into the Montgomery domain.
	Encode(x, x)
	return x
}

// SetBytes constructs a field element from bytes in big-endian order.
func (x *Elt) SetBytes(b []byte) *Elt {
	x.SetInt(new(big.Int).SetBytes(b))
	return x
}

// Int converts to a big integer.
func (x *Elt) Int() *big.Int {
	var z Elt
	// Decode from the Montgomery domain.
	Decode(&z, x)
	// Endianness swap.
	for l, r := 0, Size-1; l < r; l, r = l+1, r-1 {
		z[l], z[r] = z[r], z[l]
	}
	// Build big.Int.
	return new(big.Int).SetBytes(z[:])
}

// SetInt64Raw constructs a field element from an integer.
// This raw variant sets the value directly, bypassing any encoding/decoding steps.
func (x *Elt) SetInt64Raw(y int64) *Elt {
	x.SetIntRaw(big.NewInt(y))
	return x
}

// SetIntRaw constructs a field element from a big integer.
// This raw variant sets the value directly, bypassing any encoding/decoding steps.
func (x *Elt) SetIntRaw(y *big.Int) *Elt {
	// Reduce if outside range.
	if y.Sign() < 0 || y.Cmp(p) >= 0 {
		y = new(big.Int).Mod(y, p)
	}
	// Copy bytes into field element.
	b := y.Bytes()
	i := 0
	for ; i < len(b); i++ {
		x[i] = b[len(b)-1-i]
	}
	for ; i < Size; i++ {
		x[i] = 0
	}
	return x
}

// SetBytesRaw constructs a field element from bytes in big-endian order.
// This raw variant sets the value directly, bypassing any encoding/decoding steps.
func (x *Elt) SetBytesRaw(b []byte) *Elt {
	x.SetIntRaw(new(big.Int).SetBytes(b))
	return x
}

// IntRaw converts to a big integer.
// This raw variant sets the value directly, bypassing any encoding/decoding steps.
func (x *Elt) IntRaw() *big.Int {
	z := *x
	// Endianness swap.
	for l, r := 0, Size-1; l < r; l, r = l+1, r-1 {
		z[l], z[r] = z[r], z[l]
	}
	// Build big.Int.
	return new(big.Int).SetBytes(z[:])
}

// one is the field element 1.
var one = Elt{0x1}

// Decode decodes from the Montgomery domain.
func Decode(z *Elt, x *Elt) {
	Mul(z, x, &one)
}

// r2 is the multiplier R^2 for encoding into the Montgomery domain.
var r2 = Elt{
	0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff,
	0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xfd, 0xff, 0xff, 0xff, 0x04,
}

// Encode encodes into the Montgomery domain.
func Encode(z *Elt, x *Elt) {
	Mul(z, x, &r2)
}

// Neg computes z = -x (mod p).
func Neg(z *Elt, x *Elt) {
	Sub(z, &prime, x)
}

// Inv computes z = 1/x (mod p).
func Inv(z *Elt, x *Elt) {
	// Inversion computation is derived from the addition chain:
	//
	// _10     = 2*1
	// _11     = 1 + _10
	// _110    = 2*_11
	// _111    = 1 + _110
	// _111000 = _111 << 3
	// _111111 = _111 + _111000
	// x12     = _111111 << 6 + _111111
	// x15     = x12 << 3 + _111
	// x16     = 2*x15 + 1
	// x32     = x16 << 16 + x16
	// i53     = x32 << 15
	// x47     = x15 + i53
	// i263    = ((i53 << 17 + 1) << 143 + x47) << 47
	// return    (x47 + i263) << 2 + 1
	//
	// Operations: 255 squares 12 multiplies

	// Allocate 2 temporaries.
	var t [2]Elt

	// Step 1: z = x^0x2.
	Sqr(z, x)

	// Step 2: z = x^0x3.
	Mul(z, x, z)

	// Step 3: z = x^0x6.
	Sqr(z, z)

	// Step 4: z = x^0x7.
	Mul(z, x, z)

	// Step 7: &t[0] = x^0x38.
	Sqr(&t[0], z)
	for s := 1; s < 3; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 8: &t[0] = x^0x3f.
	Mul(&t[0], z, &t[0])

	// Step 14: &t[1] = x^0xfc0.
	Sqr(&t[1], &t[0])
	for s := 1; s < 6; s++ {
		Sqr(&t[1], &t[1])
	}

	// Step 15: &t[0] = x^0xfff.
	Mul(&t[0], &t[0], &t[1])

	// Step 18: &t[0] = x^0x7ff8.
	for s := 0; s < 3; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 19: z = x^0x7fff.
	Mul(z, z, &t[0])

	// Step 20: &t[0] = x^0xfffe.
	Sqr(&t[0], z)

	// Step 21: &t[0] = x^0xffff.
	Mul(&t[0], x, &t[0])

	// Step 37: &t[1] = x^0xffff0000.
	Sqr(&t[1], &t[0])
	for s := 1; s < 16; s++ {
		Sqr(&t[1], &t[1])
	}

	// Step 38: &t[0] = x^0xffffffff.
	Mul(&t[0], &t[0], &t[1])

	// Step 53: &t[0] = x^0x7fffffff8000.
	for s := 0; s < 15; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 54: z = x^0x7fffffffffff.
	Mul(z, z, &t[0])

	// Step 71: &t[0] = x^0xffffffff00000000.
	for s := 0; s < 17; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 72: &t[0] = x^0xffffffff00000001.
	Mul(&t[0], x, &t[0])

	// Step 215: &t[0] = x^0x7fffffff80000000800000000000000000000000000000000000.
	for s := 0; s < 143; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 216: &t[0] = x^0x7fffffff800000008000000000000000000000007fffffffffff.
	Mul(&t[0], z, &t[0])

	// Step 263: &t[0] = x^0x3fffffffc00000004000000000000000000000003fffffffffff800000000000.
	for s := 0; s < 47; s++ {
		Sqr(&t[0], &t[0])
	}

	// Step 264: z = x^0x3fffffffc00000004000000000000000000000003fffffffffffffffffffffff.
	Mul(z, z, &t[0])

	// Step 266: z = x^0xffffffff00000001000000000000000000000000fffffffffffffffffffffffc.
	for s := 0; s < 2; s++ {
		Sqr(z, z)
	}

	// Step 267: z = x^0xffffffff00000001000000000000000000000000fffffffffffffffffffffffd.
	Mul(z, x, z)
}
