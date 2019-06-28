package addchain

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/mmcloughlin/ec3/internal/bigint"

	"github.com/mmcloughlin/ec3/internal/bigints"
)

// References:
//
//	[boscoster]                Bos, Jurjen and Coster, Matthijs. Addition Chain Heuristics. In Advances in
//	                           Cryptology --- CRYPTO' 89 Proceedings, pages 400--407. 1990.
//	                           https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf
//	[github:kwantam/addchain]  Riad S. Wahby. kwantam/addchain. Github Repository. Apache License, Version 2.0.
//	                           2018. https://github.com/kwantam/addchain
//	[hehcc:exp]                Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                           Cryptography, chapter 9. 2006.
//	                           https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
//	[modboscoster]             Ayan Nandy. Modifications of Bos and Coster’s Heuristics in search of a
//	                           shorter addition chain for faster exponentiation. Masters thesis, Indian
//	                           Statistical Institute Kolkata. 2011.
//	                           http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
//	[mpnt]                     F. L. Ţiplea, S. Iftene, C. Hriţcu, I. Goriac, R. Gordân and E. Erbiceanu.
//	                           MpNT: A Multi-Precision Number Theory Package, Number Theoretical Algorithms
//	                           (I). Technical Report TR03-02, Faculty of Computer Science, "Alexandru Ioan
//	                           Cuza" University, Iasi. 2003. https://profs.info.uaic.ro/~tr/tr03-02.pdf
//	[speedsubgroup]            Stam, Martijn. Speeding up subgroup cryptosystems. PhD thesis, Technische
//	                           Universiteit Eindhoven. 2003. https://cr.yp.to/bib/2003/stam-thesis.pdf

// Heuristic suggests insertions given a current protosequence.
type Heuristic interface {
	fmt.Stringer
	Suggest(f []*big.Int, target *big.Int) []*big.Int
}

// HeuristicAlgorithm searches for an addition sequence using a heuristic at
// each step. This implements the framework given in [mpnt], page 63, with the
// heuristic playing the role of the "newnumbers" function.
type HeuristicAlgorithm struct {
	heuristic Heuristic
}

// NewHeuristicAlgorithm builds a heuristic algorithm.
func NewHeuristicAlgorithm(h Heuristic) *HeuristicAlgorithm {
	return &HeuristicAlgorithm{
		heuristic: h,
	}
}

func (h HeuristicAlgorithm) String() string {
	return fmt.Sprintf("heuristic(%v)", h.heuristic)
}

// FindSequence searches for an addition sequence for the given targets.
func (h HeuristicAlgorithm) FindSequence(targets []*big.Int) (Chain, error) {
	// Skip the special case when targets is just {1}.
	if len(targets) == 1 && bigint.EqualInt64(targets[0], 1) {
		return targets, nil
	}

	// Initialize protosequence.
	leader := bigints.Int64s(1, 2)
	proto := append(leader, targets...)
	bigints.Sort(proto)
	proto = bigints.Unique(proto)
	c := []*big.Int{}

	for len(proto) > 2 {
		// Pop the target element.
		top := len(proto) - 1
		target := proto[top]
		proto = proto[:top]
		c = bigints.InsertSortedUnique(c, target)

		// Apply heuristic.
		insert := h.heuristic.Suggest(proto, target)
		if insert == nil {
			return nil, errors.New("failed to find sequence")
		}

		// Update protosequence.
		proto = bigints.MergeUnique(proto, insert)
	}

	// Prepare the chain to return.
	c = bigints.MergeUnique(leader, c)

	return Chain(c), nil
}

// DeltaLargest implements the simple heuristic of adding the delta between the
// largest two entries in the protosequence.
type DeltaLargest struct{}

func (DeltaLargest) String() string { return "delta_largest" }

// Suggest proposes inserting target-max(f).
func (DeltaLargest) Suggest(f []*big.Int, target *big.Int) []*big.Int {
	n := len(f)
	delta := new(big.Int).Sub(target, f[n-1])
	if delta.Sign() <= 0 {
		panic("delta must be positive")
	}
	return []*big.Int{delta}
}

// Halving is the "Halving" heuristic from [boscoster].
type Halving struct{}

func (Halving) String() string { return "halving" }

// Suggest applies when the target is at least twice as big as the next largest.
// If so it will return a sequence of doublings to insert. Otherwise it will
// return nil.
func (Halving) Suggest(f []*big.Int, target *big.Int) []*big.Int {
	n := len(f)
	max, next := target, f[n-1]

	// Check the condition f / f_1 ⩾ 2ᵘ
	r := new(big.Int).Div(max, next)
	if r.BitLen() < 2 {
		return nil
	}
	u := r.BitLen() - 1

	// Compute k = floor( f / 2ᵘ ).
	k := new(big.Int).Rsh(max, uint(u))

	// Proposal to insert:
	// Delta d = f - k*2ᵘ
	// Sequence k, 2*k, ..., k*2ᵘ
	kshifts := []*big.Int{}
	for e := 0; e <= u; e++ {
		kshift := new(big.Int).Lsh(k, uint(e))
		kshifts = append(kshifts, kshift)
	}
	d := new(big.Int).Sub(max, kshifts[u])
	if bigint.IsZero(d) {
		return kshifts[:u]
	}

	return bigints.InsertSortedUnique(kshifts, d)
}

// UseFirstHeuristic is a compositite heuristic that will make the first non-nil suggestion from the sub-heuristics.
type UseFirstHeuristic []Heuristic

func (h UseFirstHeuristic) String() string {
	names := []string{}
	for _, sub := range h {
		names = append(names, sub.String())
	}
	return "use_first(" + strings.Join(names, ",") + ")"
}

// Suggest delegates to each sub-heuristic in turn and returns the first non-nil suggestion.
func (h UseFirstHeuristic) Suggest(f []*big.Int, target *big.Int) []*big.Int {
	for _, heuristic := range h {
		if insert := heuristic.Suggest(f, target); insert != nil {
			return insert
		}
	}
	return nil
}
