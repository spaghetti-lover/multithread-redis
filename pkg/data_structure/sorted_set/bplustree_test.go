package sorted_set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSet_Add(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

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

func TestSortedSet_GetScore(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

	// Add test data
	ss.Add(10.0, "member1")
	ss.Add(20.0, "member2")
	ss.Add(30.0, "member3")

	// Test existing members
	score, exists := ss.GetScore("member1")
	assert.True(t, exists)
	assert.EqualValues(t, 10.0, score)

	score, exists = ss.GetScore("member2")
	assert.True(t, exists)
	assert.EqualValues(t, 20.0, score)

	score, exists = ss.GetScore("member3")
	assert.True(t, exists)
	assert.EqualValues(t, 30.0, score)

	// Test non-existing member
	score, exists = ss.GetScore("nonexistent")
	assert.False(t, exists)
	assert.EqualValues(t, 0.0, score)

	// Test after update
	ss.Add(25.0, "member1")
	score, exists = ss.GetScore("member1")
	assert.True(t, exists)
	assert.EqualValues(t, 25.0, score)

	// Test empty member name
	score, exists = ss.GetScore("")
	assert.False(t, exists)
	assert.EqualValues(t, 0.0, score)
}

func TestSortedSet_GetRank(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

	// Add test data in non-sorted order
	ss.Add(30.0, "member3")
	ss.Add(10.0, "member1")
	ss.Add(40.0, "member4")
	ss.Add(20.0, "member2")

	// Test ranks (should be sorted by score)
	rank := ss.GetRank("member1") // score: 10.0
	assert.EqualValues(t, 0, rank)

	rank = ss.GetRank("member2") // score: 20.0
	assert.EqualValues(t, 1, rank)

	rank = ss.GetRank("member3") // score: 30.0
	assert.EqualValues(t, 2, rank)

	rank = ss.GetRank("member4") // score: 40.0
	assert.EqualValues(t, 3, rank)

	// Test non-existing member
	rank = ss.GetRank("nonexistent")
	assert.EqualValues(t, -1, rank)

	// Test after adding member with same score (should sort by member name)
	ss.Add(20.0, "member2b") // Same score as member2

	rank = ss.GetRank("member2") // "member2" < "member2b" lexicographically
	assert.EqualValues(t, 1, rank)

	rank = ss.GetRank("member2b")
	assert.EqualValues(t, 2, rank)

	// Test after score update
	ss.Add(5.0, "member4") // Update member4 to lowest score
	rank = ss.GetRank("member4")
	assert.EqualValues(t, 0, rank) // Should be first now

	rank = ss.GetRank("member1") // member1 should move to rank 1
	assert.EqualValues(t, 1, rank)

	// Test empty member name
	rank = ss.GetRank("")
	assert.EqualValues(t, -1, rank)
}

func TestSortedSet_EmptySet(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

	// Test operations on empty set
	score, exists := ss.GetScore("anything")
	assert.False(t, exists)
	assert.EqualValues(t, 0.0, score)

	rank := ss.GetRank("anything")
	assert.EqualValues(t, -1, rank)
}

func TestSortedSet_SameScoreLexicographicOrder(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

	// Add members with same score
	ss.Add(10.0, "zebra")
	ss.Add(10.0, "alpha")
	ss.Add(10.0, "beta")

	// Should be ordered lexicographically when scores are equal
	rank := ss.GetRank("alpha")
	assert.EqualValues(t, 0, rank)

	rank = ss.GetRank("beta")
	assert.EqualValues(t, 1, rank)

	rank = ss.GetRank("zebra")
	assert.EqualValues(t, 2, rank)

	// Verify scores
	score, exists := ss.GetScore("alpha")
	assert.True(t, exists)
	assert.EqualValues(t, 10.0, score)

	score, exists = ss.GetScore("beta")
	assert.True(t, exists)
	assert.EqualValues(t, 10.0, score)

	score, exists = ss.GetScore("zebra")
	assert.True(t, exists)
	assert.EqualValues(t, 10.0, score)
}

func TestSortedSet_Integration(t *testing.T) {
	config := IndexConfig{
		Type:   IndexTypeBTree,
		Degree: 4,
	}
	ss, err := NewSortedSet(config)
	assert.NoError(t, err)

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
		result := ss.Add(m.score, m.member)
		assert.EqualValues(t, 1, result)
	}

	// Test all GetScore operations
	expectedScores := map[string]float64{
		"k1": 10.0,
		"k2": 20.0,
		"k3": 30.0,
		"k4": 40.0,
		"k5": 50.0,
	}

	for member, expectedScore := range expectedScores {
		score, exists := ss.GetScore(member)
		assert.True(t, exists)
		assert.EqualValues(t, expectedScore, score)
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
		rank := ss.GetRank(member)
		assert.EqualValues(t, expectedRank, rank)
	}
}
