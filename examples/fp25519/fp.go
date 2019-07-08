// Code generated by ec3. DO NOT EDIT.

package fp25519

// Size of a field element in bytes.
const Size = 32

// Elt is a field element.
type Elt [32]uint8

// Square computes z = x^2 (mod p).
func Square(z *Elt, x *Elt) {
	Mul(z, x)
}
