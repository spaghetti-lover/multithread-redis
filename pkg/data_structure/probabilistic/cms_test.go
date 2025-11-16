package probabilistic

import (
	"math"
	"testing"
)

func TestCalcCMSDim(t *testing.T) {
	w, d := CalcCMSDim(0.01, 0.001)

	if w == 0 || d == 0 {
		t.Errorf("Expected positive width and depth, got w=%d, d=%d", w, d)
	}
}

func TestIncrByAndCount(t *testing.T) {
	cms := NewCMS(100, 5).(*CMS)

	// Insert "apple" 10 times
	for i := 0; i < 10; i++ {
		cms.IncrBy("apple", 1)
	}

	countApple := cms.Count("apple")
	if countApple < 10 {
		t.Errorf("Expected count >= 10 for 'apple', got %d", countApple)
	}

	// Insert "banana" 5 times
	for i := 0; i < 5; i++ {
		cms.IncrBy("banana", 1)
	}

	countBanana := cms.Count("banana")
	if countBanana < 5 {
		t.Errorf("Expected count >= 5 for 'banana', got %d", countBanana)
	}

	// apple should have count >= banana
	if countApple < countBanana {
		t.Errorf("Expected apple >= banana, got apple=%d, banana=%d", countApple, countBanana)
	}
}

func TestOverflowProtection(t *testing.T) {
	cms := NewCMS(50, 3).(*CMS)

	// simulate near overflow
	huge := uint64(math.MaxUint64 - 5)
	cms.IncrBy("big", huge)

	// next increment should saturate at MaxUint64
	cms.IncrBy("big", 100)

	if cms.Count("big") != math.MaxUint64 {
		t.Errorf("Expected saturated MaxUint64, got %d", cms.Count("big"))
	}
}

func TestDifferentItemsIndependence(t *testing.T) {
	cms := NewCMS(200, 5).(*CMS)

	// Insert different items
	for i := 0; i < 20; i++ {
		cms.IncrBy("x", 1)
	}
	for i := 0; i < 5; i++ {
		cms.IncrBy("y", 1)
	}

	countX := cms.Count("x")
	countY := cms.Count("y")

	if countX < countY {
		t.Errorf("Expected x >= y, got x=%d, y=%d", countX, countY)
	}
}
