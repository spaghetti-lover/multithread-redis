package sorted_set

import "fmt"

// NewOrderedIndex creates a new OrderedIndex based on the configuration
func NewOrderedIndex(config IndexConfig) (OrderedIndex, error) {
	switch config.Type {
	case IndexTypeBTree:
		if config.Degree <= 0 {
			config.Degree = 4 // Default degree
		}
		return NewBTreeIndex(config.Degree), nil
	case IndexTypeSkipList:
		return NewSkipListIndex(config.MaxLevel), nil
	default:
		return nil, fmt.Errorf("unsupported index type: %s", config.Type)
	}
}
