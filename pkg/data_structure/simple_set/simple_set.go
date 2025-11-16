package simple_set

type SimpleSet struct {
	key  string
	dict map[string]struct{}
}

func NewSimpleSet(key string) *SimpleSet {
	return &SimpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

// SADD
func (s *SimpleSet) Add(members ...string) int {
	added := 0

	for _, m := range members {
		if _, exist := s.dict[m]; !exist {
			s.dict[m] = struct{}{}
			added += 1
		}
	}

	return added
}

// SREM
func (s *SimpleSet) Rem(members ...string) int {
	removed := 0
	for _, m := range members {
		if _, exist := s.dict[m]; exist {
			delete(s.dict, m)
			removed += 1
		}
	}
	return removed
}

// SISMEMBER
func (s *SimpleSet) IsMember(member string) int {
	_, exist := s.dict[member]

	if exist {
		return 1
	}
	return 0
}

// SMEMBERS
func (s *SimpleSet) Members() []string {
	m := make([]string, 0, len(s.dict))

	for k := range s.dict {
		m = append(m, k)
	}

	return m
}
