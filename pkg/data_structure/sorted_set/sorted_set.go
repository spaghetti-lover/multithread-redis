package sorted_set

type SortedSet struct {
	Index       OrderedIndex
	MemberScore map[string]float64
}

// NewSortedSet creates a new SortedSet with the specified index configuration
func NewSortedSet(config IndexConfig) (*SortedSet, error) {
	index, err := NewOrderedIndex(config)
	if err != nil {
		return nil, err
	}

	return &SortedSet{
		Index:       index,
		MemberScore: make(map[string]float64),
	}, nil
}

// NewSortedSetWithBTree creates a SortedSet with B+ Tree (convenience function)
func NewSortedSetWithBTree(degree int) (*SortedSet, error) {
	return NewSortedSet(IndexConfig{
		Type:   IndexTypeBTree,
		Degree: degree,
	})
}

func (ss *SortedSet) Add(score float64, member string) int {
	if member == "" {
		return 0
	}

	oldScore, exists := ss.MemberScore[member]
	if exists {
		if oldScore == score {
			return 0 // No change needed
		}
		// Remove old entry using score (efficient)
		ss.Index.RemoveByScore(oldScore, member)
		ss.Index.Add(score, member)
		ss.MemberScore[member] = score
		return 0
	}

	result := ss.Index.Add(score, member)
	if result == 1 {
		ss.MemberScore[member] = score
	}
	return result
}

func (ss *SortedSet) GetScore(member string) (float64, bool) {
	score, exists := ss.MemberScore[member]
	return score, exists
}

func (ss *SortedSet) GetRank(member string) int {
	return ss.Index.GetRank(member)
}

func (ss *SortedSet) Remove(member string) int {
	if member == "" {
		return 0
	}

	score, exists := ss.MemberScore[member]
	if !exists {
		return 0
	}

	// Use score for efficient removal
	result := ss.Index.RemoveByScore(score, member)
	if result == 1 {
		delete(ss.MemberScore, member)
	}
	return result
}
