package probabilistic

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const (
	Ln2       float64 = 0.693147180559945
	Ln2Square float64 = 0.480453013918201
	ABigSeed  uint32  = 0x9747b28c
)

type Bloom struct {
	Hashes      int // number of hash functions
	Entries     uint64
	Error       float64
	bitPerEntry float64
	bf          []uint64
	bits        uint64 // size of bf in bits
	words       uint64 // number of 64-bit words in bf
}

type HashValue struct {
	a uint64
	b uint64
}

func calcBpe(err float64) float64 {
	num := math.Log(err)
	return math.Abs(-(num / Ln2Square))
}

func NewBloomFilter(entries uint64, errorRate float64) MembershipTester {
	bloom := &Bloom{
		Entries: entries,
		Error:   errorRate,
	}
	bloom.bitPerEntry = calcBpe(errorRate)
	// Calculate necessary number of bits
	bits := uint64(float64(entries) * bloom.bitPerEntry)
	// round to lowest common ancestor of 64
	if bits%64 != 0 {
		bloom.bits = ((bits / 64) + 1) * 64
	} else {
		bloom.bits = bits
	}

	bloom.words = bloom.bits / 64
	bloom.Hashes = int(math.Ceil(Ln2 * bloom.bitPerEntry))
	bloom.bf = make([]uint64, bloom.words)

	return bloom
}

func (b *Bloom) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x, y := hasher.Sum128()
	return HashValue{a: x, b: y}
}

func (b *Bloom) Add(entry string) {
	initHash := b.CalcHash(entry)

	for i := 0; i < b.Hashes; i++ {
		hash := (initHash.a + initHash.b*uint64(i)) % b.bits
		word := hash >> 6  // chia 64
		bit := hash & 0x3F // mod 64 (mask 6 bit)
		b.bf[word] |= 1 << bit
	}
}

func (b *Bloom) Exist(entry string) bool {
	initHash := b.CalcHash(entry)
	for i := 0; i < b.Hashes; i++ {
		hash := (initHash.a + initHash.b*uint64(i)) % b.bits
		word := hash >> 6
		bit := hash & 0x3F
		if (b.bf[word] & (1 << bit)) == 0 {
			return false
		}
	}
	return true
}

func (b *Bloom) AddHash(initHash HashValue) {
	for i := 0; i < b.Hashes; i++ {
		hash := (initHash.a + initHash.b*uint64(i)) % b.bits
		word := hash >> 6
		bit := hash & 0x3F
		b.bf[word] |= 1 << bit
	}
}

func (b *Bloom) ExistHash(initHash HashValue) bool {
	for i := 0; i < b.Hashes; i++ {
		hash := (initHash.a + initHash.b*uint64(i)) % b.bits
		word := hash >> 6
		bit := hash & 0x3F
		if (b.bf[word] & (1 << bit)) == 0 {
			return false
		}
	}
	return true
}
