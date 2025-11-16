package hash_table

import (
	"log"
	"sort"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
)

type LruEvictionCandidate struct {
	key            string
	lastAccessTime uint32
}

type LruEvictionPool struct {
	pool []*LruEvictionCandidate
}

type ByLastAccessTime []*LruEvictionCandidate

func (a ByLastAccessTime) Len() int {
	return len(a)
}

func (a ByLastAccessTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByLastAccessTime) Less(i, j int) bool {
	return a[i].lastAccessTime < a[j].lastAccessTime
}

// Push adds a new item to the pool, maintains the lastAccessTime accenting order (old items are on the left).
// If pool size > EpoolMaxSize, removes the newest item.
func (p *LruEvictionPool) Push(key string, lastAccessTime uint32) {
	newItem := &LruEvictionCandidate{
		key:            key,
		lastAccessTime: lastAccessTime,
	}
	// Note: In Redis implementation, it does not explicitly check if a key is already in the eviction pool
	// before attempting to insert it. This could lead to a key being in the pool twice
	// if it's sampled and inserted a second time. However, since the eviction pool is very small (EpoolMaxSize is 16)
	// and the random sampling is just a small fraction of the total keys, the probability of this happening
	// is extremely low
	// Ref: https://github.com/redis/redis/blob/unstable/src/evict.c#L126
	exist := false
	for i := 0; i < len(p.pool); i++ {
		if p.pool[i].key == key {
			exist = true
			p.pool[i] = newItem
		}
	}
	if !exist {
		p.pool = append(p.pool, newItem)
	}
	sort.Sort(ByLastAccessTime(p.pool))
	if len(p.pool) > config.EpoolMaxSize {
		lastIndex := len(p.pool) - 1
		p.pool = p.pool[:lastIndex]
	}
}

// Pop returns the oldest item in the pool
func (p *LruEvictionPool) Pop() *LruEvictionCandidate {
	if len(p.pool) == 0 {
		return nil
	}
	oldestItem := p.pool[0]
	p.pool = p.pool[1:]
	return oldestItem
}

func newEpool(size int) *LruEvictionPool {
	return &LruEvictionPool{
		pool: make([]*LruEvictionCandidate, size),
	}
}

func (d *Dict) populateEpool() {
	remain := config.EpoolLRUSampleSize
	for k := range d.dictStore {
		ePool.Push(k, d.dictStore[k].LastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}
	log.Println("Epool:")
	for _, item := range ePool.pool {
		log.Println(item.key, item.lastAccessTime)
	}
}

func (d *Dict) evictLru() {
	d.populateEpool()
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("Trigger LRU eviction, evict count: ", evictCount)
	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Del(item.key)
		}
	}
}

var ePool *LruEvictionPool = newEpool(0)
