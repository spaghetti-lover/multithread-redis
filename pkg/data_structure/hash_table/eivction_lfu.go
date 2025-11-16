package hash_table

import (
	"log"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
)

type LfuEvictionCandidate struct {
	key       string
	frequency uint32
}

type LfuEvictionPool struct {
	pool []*LfuEvictionCandidate
}

type ByFrequency []*LfuEvictionCandidate

func (d *Dict) evictLfu() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("Trigger LRU eviction, evict count: ", evictCount)

	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Del(item.key)
		}
	}
}
