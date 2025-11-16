package sorted_set

type Item struct {
	Score  float64
	Member string
}

func (i *Item) CompareTo(other *Item) int {
	if i.Score < other.Score {
		return -1
	}
	if i.Score > other.Score {
		return 1
	}

	//If Score is equal. Use Member
	if i.Member < other.Member {
		return -1
	}
	if i.Member > other.Member {
		return 1
	}
	return 0
}
