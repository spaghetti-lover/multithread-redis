package sorted_set

import (
	"math/rand"
	"strings"
	"time"
)

const SkiplistMaxLevel = 32

type SkiplistLevel struct {
	forward *SkiplistNode
	// span is number of nodes between current node and node->forward at current level
	span uint32
}

type SkiplistNode struct {
	ele      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type SkipListIndex struct {
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint32
	level  int
	rnd    *rand.Rand
}

func (sl *SkipListIndex) randomLevel() int {
	level := 1
	for sl.rnd.Intn(2) == 1 {
		level++
	}
	if level > SkiplistMaxLevel {
		return SkiplistMaxLevel
	}
	return level
}

func (sl *SkipListIndex) createNode(level int, score float64, ele string) *SkiplistNode {
	res := &SkiplistNode{
		ele:      ele,
		score:    score,
		backward: nil,
	}
	res.levels = make([]SkiplistLevel, level)
	return res
}

// NewSkipListIndex creates a new SkipList index that implements OrderedIndex
func NewSkipListIndex(maxLevel int) OrderedIndex {
	sl := &SkipListIndex{
		length: 0,
		level:  1,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	sl.head = sl.createNode(SkiplistMaxLevel, 0, "")
	sl.head.backward = nil
	sl.tail = nil

	return sl
}

// Add implements OrderedIndex.Add
func (sl *SkipListIndex) Add(score float64, member string) int {
	// Simplified - just insert, let SortedSet handle duplicates
	sl.insert(score, member)
	return 1
}

// GetRank implements OrderedIndex.GetRank
func (sl *SkipListIndex) GetRank(member string) int {
	node := sl.findMember(member)
	if node == nil {
		return -1
	}

	rank := sl.getRankByScoreAndMember(node.score, member)
	if rank > 0 {
		return int(rank - 1) // Convert to 0-based index
	}
	return -1
}

// Private helper methods

func (sl *SkipListIndex) findMember(member string) *SkiplistNode {
	x := sl.head.levels[0].forward
	for x != nil {
		if x.ele == member {
			return x
		}
		x = x.levels[0].forward
	}
	return nil
}

func (sl *SkipListIndex) insert(score float64, ele string) *SkiplistNode {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	rank := [SkiplistMaxLevel]uint32{}
	x := sl.head

	// Traverse the skiplist from the highest level down to find the insertion point
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			rank[i] += x.levels[i].span
			x = x.levels[i].forward
		}
		update[i] = x
	}

	// Determine the level of the new node
	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].levels[i].span = sl.length
		}
		sl.level = level
	}

	// Create new node and insert
	x = sl.createNode(level, score, ele)
	for i := 0; i < level; i++ {
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x
		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// Increase span for untouched levels
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	// Update backward pointers
	if update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x
	} else {
		sl.tail = x
	}

	sl.length++
	return x
}

func (sl *SkipListIndex) getRankByScoreAndMember(score float64, ele string) uint32 {
	x := sl.head
	var rank uint32 = 0

	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score &&
				strings.Compare(x.levels[i].forward.ele, ele) <= 0)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x.score == score && strings.Compare(x.ele, ele) == 0 {
			return rank
		}
	}
	return 0
}

func (sl *SkipListIndex) deleteByNode(node *SkiplistNode) {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	x := sl.head

	// Find the node to delete
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < node.score ||
			(x.levels[i].forward.score == node.score &&
				strings.Compare(x.levels[i].forward.ele, node.ele) == -1)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}

	sl.deleteNode(node, update)
}

func (sl *SkipListIndex) deleteNode(x *SkiplistNode, update [SkiplistMaxLevel]*SkiplistNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span += x.levels[i].span - 1
			update[i].levels[i].forward = x.levels[i].forward
		} else {
			update[i].levels[i].span--
		}
	}

	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

// Remove removes member with known score instead of remove member only (more efficient - O(log N))
func (sl *SkipListIndex) RemoveByScore(score float64, member string) int {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	x := sl.head

	// Find the node to delete
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score &&
				strings.Compare(x.levels[i].forward.ele, member) == -1)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}

	x = x.levels[0].forward
	if x != nil && x.score == score && x.ele == member {
		sl.deleteNode(x, update)
		return 1
	}
	return 0
}
