package probabilistic

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type CMS struct {
	width   uint64
	depth   uint64
	counter []uint64 // 1D slides: index = i*width + j for better cache locality than 2D array
}

func NewCMS(w uint64, d uint64) FrequencyEstimator {
	return &CMS{
		width:   w,
		depth:   d,
		counter: make([]uint64, w*d),
	}
}

// CalcCMSDim calculates width and depth from error rate and probability.
// Formula: w = ceil(e/ε), d = ceil(ln(1/δ)).
// e: epsilon, ε: errorRate
// δ: errorProbability
func CalcCMSDim(errRate float64, probRate float64) (uint64, uint64) {
	w := uint64(math.Ceil(math.E / errRate))
	d := uint64(math.Ceil(math.Log(1 / probRate)))
	return w, d
}

// getIndex computes index in 1D array.
func (c *CMS) getIndex(row uint64, col uint64) uint64 {
	return row*c.width + col
}

// calcHash calculates a 64-bit hash for the given item and seed.
func (c *CMS) calcHash(item string, seed uint32) uint64 {
	hasher := murmur3.New64WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum64()
}

// IncrBy increments an item by value and returns the estimated count.
func (c *CMS) IncrBy(item string, value uint64) uint64 {
	var minCount uint64 = math.MaxUint64
	for i := uint64(0); i < c.depth; i++ {
		hash := c.calcHash(item, uint32(i))
		j := hash % c.width
		idx := c.getIndex(i, j)

		// Prevent overflow
		if math.MaxUint64-c.counter[idx] < value {
			c.counter[idx] = math.MaxUint64
		} else {
			c.counter[idx] += value
		}
		if c.counter[idx] < minCount {
			minCount = c.counter[idx]
		}
	}

	return minCount
}

// Count returns the estimated count of an item.
func (c *CMS) Count(item string) uint64 {
	var minCount uint64 = math.MaxUint64
	for i := uint64(0); i < c.depth; i++ {
		hash := c.calcHash(item, uint32(i))
		j := hash % c.width
		idx := c.getIndex(i, j)

		if c.counter[idx] < minCount {
			minCount = c.counter[idx]
		}
	}
	return minCount
}
