package sorted_set

// BTreeNode represents a node in the B+ Tree
type BTreeNode struct {
	Items    []*Item
	Children []*BTreeNode
	IsLeaf   bool
	Parent   *BTreeNode
	Next     *BTreeNode
}

// BTreeIndex implements OrderedIndex using B+ Tree
type BTreeIndex struct {
	Root   *BTreeNode
	Degree int
}

// NewBTreeIndex creates a new B+ Tree index
func NewBTreeIndex(degree int) OrderedIndex {
	return &BTreeIndex{
		Root:   &BTreeNode{IsLeaf: true},
		Degree: degree,
	}
}

func (t *BTreeIndex) Add(score float64, member string) int {
	item := &Item{Score: score, Member: member}

	if len(member) == 0 {
		return 0
	}

	// Find the correct leaf to insert into
	node := t.Root
	for !node.IsLeaf {
		i := 0
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// Check if the member already exists in the leaf node
	for i, existingItem := range node.Items {
		if existingItem.Member == member {
			node.Items[i].Score = score
			return 0 // Updated existing item
		}
	}

	// Member does not exist, insert it into the sorted position
	i := 0
	for i < len(node.Items) && item.CompareTo(node.Items[i]) > 0 {
		i++
	}
	node.Items = append(node.Items[:i], append([]*Item{item}, node.Items[i:]...)...)

	// Split the node if it's over capacity
	if len(node.Items) > t.Degree-1 {
		t.splitNode(node)
	}
	return 1 // Added new item
}

func (t *BTreeIndex) RemoveByScore(score float64, member string) int {
	// Optimized removal with known score

	// Find the correct leaf node using score (O(log n))
	node := t.Root
	for !node.IsLeaf {
		i := 0
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// Find and remove the item in the leaf node
	for i, item := range node.Items {
		if item.Member == member && item.Score == score {
			// Remove the item
			node.Items = append(node.Items[:i], node.Items[i+1:]...)

			// Handle underflow if necessary
			if len(node.Items) < (t.Degree-1)/2 && node.Parent != nil {
				t.handleUnderflow(node)
			}

			return 1 // Successfully removed
		}
	}

	return 0 // Item not found
}

func (t *BTreeIndex) handleUnderflow(node *BTreeNode) {
	// Root node doesn't need to handle underflow
	if node.Parent == nil {
		return
	}

	minItems := (t.Degree - 1) / 2
	if len(node.Items) >= minItems {
		return // No underflow
	}

	parent := node.Parent
	nodeIndex := t.findChildIndex(parent, node)

	// Try to borrow from left sibling
	if nodeIndex > 0 {
		leftSibling := parent.Children[nodeIndex-1]
		if len(leftSibling.Items) > minItems {
			t.borrowFromLeftSibling(node, leftSibling, nodeIndex-1)
			return
		}
	}

	// Try to borrow from right sibling
	if nodeIndex < len(parent.Children)-1 {
		rightSibling := parent.Children[nodeIndex+1]
		if len(rightSibling.Items) > minItems {
			t.borrowFromRightSibling(node, rightSibling, nodeIndex)
			return
		}
	}

	// Cannot borrow, must merge
	if nodeIndex > 0 {
		// Merge with left sibling
		leftSibling := parent.Children[nodeIndex-1]
		t.mergeWithLeftSibling(node, leftSibling, nodeIndex-1)
	} else if nodeIndex < len(parent.Children)-1 {
		// Merge with right sibling
		rightSibling := parent.Children[nodeIndex+1]
		t.mergeWithRightSibling(node, rightSibling, nodeIndex)
	}
}

func (t *BTreeIndex) findChildIndex(parent *BTreeNode, child *BTreeNode) int {
	for i, c := range parent.Children {
		if c == child {
			return i
		}
	}
	return -1
}

func (t *BTreeIndex) borrowFromLeftSibling(node *BTreeNode, leftSibling *BTreeNode, separatorIndex int) {
	parent := node.Parent

	if node.IsLeaf {
		// Borrow the last item from left sibling
		borrowedItem := leftSibling.Items[len(leftSibling.Items)-1]
		leftSibling.Items = leftSibling.Items[:len(leftSibling.Items)-1]

		// Insert borrowed item at the beginning of current node
		node.Items = append([]*Item{borrowedItem}, node.Items...)

		// Update separator in parent (first item of current node)
		parent.Items[separatorIndex] = node.Items[0]
	} else {
		// For internal nodes, move separator down and borrow from sibling
		separatorItem := parent.Items[separatorIndex]
		borrowedItem := leftSibling.Items[len(leftSibling.Items)-1]
		borrowedChild := leftSibling.Children[len(leftSibling.Children)-1]

		// Remove from left sibling
		leftSibling.Items = leftSibling.Items[:len(leftSibling.Items)-1]
		leftSibling.Children = leftSibling.Children[:len(leftSibling.Children)-1]

		// Add to current node
		node.Items = append([]*Item{separatorItem}, node.Items...)
		node.Children = append([]*BTreeNode{borrowedChild}, node.Children...)
		borrowedChild.Parent = node

		// Update separator in parent
		parent.Items[separatorIndex] = borrowedItem
	}
}

func (t *BTreeIndex) borrowFromRightSibling(node *BTreeNode, rightSibling *BTreeNode, separatorIndex int) {
	parent := node.Parent

	if node.IsLeaf {
		// Borrow the first item from right sibling
		borrowedItem := rightSibling.Items[0]
		rightSibling.Items = rightSibling.Items[1:]

		// Insert borrowed item at the end of current node
		node.Items = append(node.Items, borrowedItem)

		// Update separator in parent (first item of right sibling)
		if len(rightSibling.Items) > 0 {
			parent.Items[separatorIndex] = rightSibling.Items[0]
		}
	} else {
		// For internal nodes, move separator down and borrow from sibling
		separatorItem := parent.Items[separatorIndex]
		borrowedItem := rightSibling.Items[0]
		borrowedChild := rightSibling.Children[0]

		// Remove from right sibling
		rightSibling.Items = rightSibling.Items[1:]
		rightSibling.Children = rightSibling.Children[1:]

		// Add to current node
		node.Items = append(node.Items, separatorItem)
		node.Children = append(node.Children, borrowedChild)
		borrowedChild.Parent = node

		// Update separator in parent
		parent.Items[separatorIndex] = borrowedItem
	}
}

func (t *BTreeIndex) mergeWithLeftSibling(node *BTreeNode, leftSibling *BTreeNode, separatorIndex int) {
	parent := node.Parent

	if node.IsLeaf {
		// Merge all items from current node to left sibling
		leftSibling.Items = append(leftSibling.Items, node.Items...)
		leftSibling.Next = node.Next

		// Remove separator from parent and current node from children
		parent.Items = append(parent.Items[:separatorIndex], parent.Items[separatorIndex+1:]...)
		parent.Children = append(parent.Children[:separatorIndex+1], parent.Children[separatorIndex+2:]...)
	} else {
		// For internal nodes, bring separator down and merge
		separatorItem := parent.Items[separatorIndex]

		// Merge items and children
		leftSibling.Items = append(leftSibling.Items, separatorItem)
		leftSibling.Items = append(leftSibling.Items, node.Items...)
		leftSibling.Children = append(leftSibling.Children, node.Children...)

		// Update parent pointers
		for _, child := range node.Children {
			child.Parent = leftSibling
		}

		// Remove separator from parent and current node from children
		parent.Items = append(parent.Items[:separatorIndex], parent.Items[separatorIndex+1:]...)
		parent.Children = append(parent.Children[:separatorIndex+1], parent.Children[separatorIndex+2:]...)
	}

	// Check if parent needs handling
	if parent != t.Root {
		t.handleUnderflow(parent)
	} else if len(parent.Items) == 0 && len(parent.Children) == 1 {
		// Root has no items and only one child, make the child the new root
		t.Root = parent.Children[0]
		t.Root.Parent = nil
	}
}

func (t *BTreeIndex) mergeWithRightSibling(node *BTreeNode, rightSibling *BTreeNode, separatorIndex int) {
	parent := node.Parent

	if node.IsLeaf {
		// Merge all items from right sibling to current node
		node.Items = append(node.Items, rightSibling.Items...)
		node.Next = rightSibling.Next

		// Remove separator from parent and right sibling from children
		parent.Items = append(parent.Items[:separatorIndex], parent.Items[separatorIndex+1:]...)
		parent.Children = append(parent.Children[:separatorIndex+1], parent.Children[separatorIndex+2:]...)
	} else {
		// For internal nodes, bring separator down and merge
		separatorItem := parent.Items[separatorIndex]

		// Merge items and children
		node.Items = append(node.Items, separatorItem)
		node.Items = append(node.Items, rightSibling.Items...)
		node.Children = append(node.Children, rightSibling.Children...)

		// Update parent pointers
		for _, child := range rightSibling.Children {
			child.Parent = node
		}

		// Remove separator from parent and right sibling from children
		parent.Items = append(parent.Items[:separatorIndex], parent.Items[separatorIndex+1:]...)
		parent.Children = append(parent.Children[:separatorIndex+1], parent.Children[separatorIndex+2:]...)
	}

	// Check if parent needs handling
	if parent != t.Root {
		t.handleUnderflow(parent)
	} else if len(parent.Items) == 0 && len(parent.Children) == 1 {
		// Root has no items and only one child, make the child the new root
		t.Root = parent.Children[0]
		t.Root.Parent = nil
	}
}

func (t *BTreeIndex) GetRank(member string) int {
	rank := 0

	// Find the first leaf node
	node := t.Root
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse all leaf nodes from the beginning
	for node != nil {
		for _, item := range node.Items {
			if item.Member == member {
				return rank
			}
			rank++
		}
		node = node.Next
	}

	return -1
}

func (t *BTreeIndex) GetByRank(rank int) *Item {
	if rank < 0 {
		return nil
	}

	currentRank := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes
	for node != nil {
		for _, item := range node.Items {
			if currentRank == rank {
				return item
			}
			currentRank++
		}
		node = node.Next
	}

	return nil
}

func (t *BTreeIndex) GetRange(min, max float64) []*Item {
	var result []*Item
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes and collect items in range
	for node != nil {
		for _, item := range node.Items {
			if item.Score >= min && item.Score <= max {
				result = append(result, item)
			}
			if item.Score > max {
				return result
			}
		}
		node = node.Next
	}

	return result
}

func (t *BTreeIndex) GetRangeByRank(start, end int) []*Item {
	var result []*Item
	if start < 0 || end < start {
		return result
	}

	currentRank := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes
	for node != nil {
		for _, item := range node.Items {
			if currentRank >= start && currentRank <= end {
				result = append(result, item)
			}
			if currentRank > end {
				return result
			}
			currentRank++
		}
		node = node.Next
	}

	return result
}

func (t *BTreeIndex) Count() int {
	count := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Count items in all leaf nodes
	for node != nil {
		count += len(node.Items)
		node = node.Next
	}

	return count
}

func (t *BTreeIndex) Clear() {
	t.Root = &BTreeNode{IsLeaf: true}
}

// Helper methods for B+ Tree operations
func (t *BTreeIndex) splitNode(node *BTreeNode) {
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BTreeIndex) splitLeaf(node *BTreeNode) {
	medianIndex := len(node.Items) / 2

	newLeaf := &BTreeNode{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
	}

	newLeaf.Items = append(newLeaf.Items, node.Items[medianIndex:]...)
	node.Items = node.Items[:medianIndex]
	node.Next = newLeaf

	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*BTreeNode{newLeaf}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BTreeIndex) splitInternal(node *BTreeNode) {
	medianIndex := len(node.Items) / 2

	newInternal := &BTreeNode{
		IsLeaf: false,
		Parent: node.Parent,
	}

	promotedItem := node.Items[medianIndex]

	newInternal.Items = append(newInternal.Items, node.Items[medianIndex+1:]...)
	newInternal.Children = append(newInternal.Children, node.Children[medianIndex+1:]...)

	node.Items = node.Items[:medianIndex]
	node.Children = node.Children[:medianIndex+1]

	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	parent := node.Parent
	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*BTreeNode{newInternal}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BTreeIndex) splitRoot() {
	oldRoot := t.Root
	newRoot := &BTreeNode{}

	t.Root = newRoot
	oldRoot.Parent = newRoot
	newRoot.Children = append(newRoot.Children, oldRoot)

	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}
