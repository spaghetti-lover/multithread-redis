package probabilistic

import (
	"math/rand"
	"testing"
	"time"
)

func TestBloomBasic(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01).(*Bloom)

	// Thêm 1 phần tử
	bf.Add("golang")

	// Kiểm tra tồn tại
	if !bf.Exist("golang") {
		t.Errorf("expected 'golang' to exist in Bloom filter")
	}

	// Kiểm tra phần tử chưa thêm
	if bf.Exist("python") {
		t.Errorf("did not expect 'python' to exist in Bloom filter")
	}
}

func TestBloomMultipleEntries(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01).(*Bloom)
	words := []string{"apple", "banana", "cherry", "durian", "elderberry"}

	// Thêm nhiều từ
	for _, w := range words {
		bf.Add(w)
	}

	// Tất cả phải tồn tại
	for _, w := range words {
		if !bf.Exist(w) {
			t.Errorf("expected %s to exist in Bloom filter", w)
		}
	}

	// Một số từ khác chưa thêm → có thể false positive nhưng xác suất thấp
	if bf.Exist("nonexistent") {
		t.Logf("warning: 'nonexistent' returned true (false positive)")
	}
}

func TestBloomFalsePositiveRate(t *testing.T) {
	n := 10000
	errorRate := 0.01
	bf := NewBloomFilter(uint64(n), errorRate).(*Bloom)

	// Thêm n phần tử
	for i := 0; i < n; i++ {
		bf.Add(string(rune(i)))
	}

	// Kiểm tra false positives trên tập khác
	rand.Seed(time.Now().UnixNano())
	falsePositives := 0
	trials := 10000
	for i := 0; i < trials; i++ {
		word := string(rune(n + i))
		if bf.Exist(word) {
			falsePositives++
		}
	}

	fpRate := float64(falsePositives) / float64(trials)
	t.Logf("false positive rate observed = %.4f (expected ≤ %.2f)", fpRate, errorRate)

	if fpRate > errorRate*1.5 { // cho phép dao động
		t.Errorf("false positive rate too high: got %.4f, want ≤ %.2f", fpRate, errorRate)
	}
}
