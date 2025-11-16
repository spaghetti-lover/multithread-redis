package sorted_set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipListIndex_Add(t *testing.T) {
	ss, err := NewSortedSet(IndexConfig{
		Type:   IndexTypeSkipList,
		Degree: 16,
	})
	if err != nil {
		t.Fatal("Failed to create SortedSet with SkipListIndex:", err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, ss)
	assert.NotNil(t, ss)

	// Test adding new members
	result := ss.Add(10.0, "member1")
	assert.EqualValues(t, 1, result)

	result = ss.Add(20.0, "member2")
	assert.EqualValues(t, 1, result)

	result = ss.Add(30.0, "member3")
	assert.EqualValues(t, 1, result)

	// Test updating existing member
	result = ss.Add(15.0, "member1")
	assert.EqualValues(t, 0, result) // Should return 0 for update

	// Test adding with same score
	result = ss.Add(20.0, "member4")
	assert.EqualValues(t, 1, result)

	// Test edge cases
	result = ss.Add(0.0, "zero")
	assert.EqualValues(t, 1, result)

	result = ss.Add(-10.0, "negative")
	assert.EqualValues(t, 1, result)

	// Test empty member name
	result = ss.Add(100.0, "")
	assert.EqualValues(t, 0, result) // Should handle empty string
}

func TestSkipListIndex_GetRank(t *testing.T) {
	skiplist := NewSkipListIndex(16)

	// Add test data in non-sorted order
	skiplist.Add(30.0, "member3")
	skiplist.Add(10.0, "member1")
	skiplist.Add(40.0, "member4")
	skiplist.Add(20.0, "member2")

	// Test ranks (should be sorted by score)
	rank := skiplist.GetRank("member1") // score: 10.0
	assert.EqualValues(t, 0, rank)

	rank = skiplist.GetRank("member2") // score: 20.0
	assert.EqualValues(t, 1, rank)

	rank = skiplist.GetRank("member3") // score: 30.0
	assert.EqualValues(t, 2, rank)

	rank = skiplist.GetRank("member4") // score: 40.0
	assert.EqualValues(t, 3, rank)

	// Test non-existing member
	rank = skiplist.GetRank("nonexistent")
	assert.EqualValues(t, -1, rank)

	// Test after adding member with same score (should sort by member name)
	skiplist.Add(20.0, "member2b") // Same score as member2

	rank = skiplist.GetRank("member2") // "member2" < "member2b" lexicographically
	assert.EqualValues(t, 1, rank)

	rank = skiplist.GetRank("member2b")
	assert.EqualValues(t, 2, rank)

	// Test after score update
	skiplist.Add(5.0, "member4") // Update member4 to lowest score
	rank = skiplist.GetRank("member4")
	assert.EqualValues(t, 0, rank) // Should be first now

	rank = skiplist.GetRank("member1") // member1 should move to rank 1
	assert.EqualValues(t, 1, rank)

	// Test empty member name
	rank = skiplist.GetRank("")
	assert.EqualValues(t, -1, rank)
}

func TestSkipListIndex_EmptySet(t *testing.T) {
	skiplist := NewSkipListIndex(16)

	// Test operations on empty skiplist
	rank := skiplist.GetRank("anything")
	assert.EqualValues(t, -1, rank)
}

func TestSkipListIndex_SameScoreLexicographicOrder(t *testing.T) {
	skiplist := NewSkipListIndex(16)

	// Add members with same score
	skiplist.Add(10.0, "zebra")
	skiplist.Add(10.0, "alpha")
	skiplist.Add(10.0, "beta")

	// Should be ordered lexicographically when scores are equal
	rank := skiplist.GetRank("alpha")
	assert.EqualValues(t, 0, rank)

	rank = skiplist.GetRank("beta")
	assert.EqualValues(t, 1, rank)

	rank = skiplist.GetRank("zebra")
	assert.EqualValues(t, 2, rank)
}

func TestSkipListIndex_Integration(t *testing.T) {
	skiplist := NewSkipListIndex(16)

	// Add multiple members
	members := []struct {
		score  float64
		member string
	}{
		{50.0, "k5"},
		{10.0, "k1"},
		{30.0, "k3"},
		{20.0, "k2"},
		{40.0, "k4"},
	}

	for _, m := range members {
		result := skiplist.Add(m.score, m.member)
		assert.EqualValues(t, 1, result)
	}

	// Test all GetRank operations
	expectedRanks := map[string]int{
		"k1": 0,
		"k2": 1,
		"k3": 2,
		"k4": 3,
		"k5": 4,
	}

	for member, expectedRank := range expectedRanks {
		rank := skiplist.GetRank(member)
		assert.EqualValues(t, expectedRank, rank)
	}
}
