// Code generated by command: go run asm.go -out fp25519.s -stubs stub.go. DO NOT EDIT.

package ec3

func Add25519(x *[32]byte, y *[32]byte)

func Mul(z *[64]byte, x *[32]byte, y *[32]byte)