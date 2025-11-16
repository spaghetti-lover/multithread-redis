package hash_table

import (
	"log"
	"time"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
)

type Obj struct {
	Value          interface{}
	LastAccessTime uint32
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func CreateDict() *Dict {
	res := Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
	return &res
}

func (d *Dict) GetExpireDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) GetDictStore() map[string]*Obj {
	return d.dictStore
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

// Remove TTL of a key (like PERSIST)
func (d *Dict) Persist(key string) {
	delete(d.expiredDictStore, key)
}

func (d *Dict) NewObj(key string, value interface{}, ttlMs int64) *Obj {
	obj := &Obj{
		Value:          value,
		LastAccessTime: now(),
	}
	if ttlMs > 0 {
		d.SetExpiry(key, ttlMs)
	} else {
		// Ensure old TTL is removed when overwriting without TTL. For example, SET key value EX 10 followed by SET key value
		d.Persist(key)
	}
	return obj
}

func (d *Dict) GetExpiry(key string) (uint64, bool) {
	exp, exist := d.expiredDictStore[key]
	return exp, exist
}

func (d *Dict) SetExpiry(key string, ttlMs int64) {
	d.expiredDictStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) HasExpired(key string) bool {
	exp, exist := d.expiredDictStore[key]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (d *Dict) Get(k string) *Obj {
	v := d.dictStore[k]
	if v != nil {
		v.LastAccessTime = now()
		if d.HasExpired(k) {
			d.Del(k)
			return nil
		}
	}
	return v
}

func (d *Dict) Set(k string, obj *Obj) {
	if len(d.dictStore) >= config.MaxKeyNumber {
		d.evict()
	}
	d.dictStore[k] = obj
}

func (d *Dict) evictRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("Trigger random eviction, evict count: ", evictCount)
	for k := range d.dictStore {
		d.Del(k)
		evictCount--
		if evictCount == 0 {
			break
		}
	}
}

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	default:
		d.evictRandom()
		log.Printf("Warning: Unknown eviction policy %s, defaulting to allkeys-random", config.EvictionPolicy)
	}
}

func (d *Dict) Del(k string) bool {
	if _, exist := d.dictStore[k]; exist {
		delete(d.dictStore, k)
		delete(d.expiredDictStore, k)
		return true
	}
	return false
}

// ExpiringKeysCount returns number of keys that currently have a TTL and are not expired.
// It does not clean up stale TTLs and expired keys.
func (d *Dict) ExpiringKeysCount() int {
	now := uint64(time.Now().UnixMilli())
	count := 0
	for _, exp := range d.expiredDictStore {
		if exp > now {
			count++
		}
	}
	return count
}

// TLL_Avg returns average TTL of keys that currently have a TTL and are not expired.
// It does not clean up stale TTLs and expired keys.
func (d *Dict) TLL_Avg() uint64 {
	now := uint64(time.Now().UnixMilli())
	res := uint64(0)
	cnt := uint64(0)
	for _, exp := range d.expiredDictStore {
		if exp > now {
			res = res + (exp - now)
			cnt++
		}
	}
	if cnt == 0 {
		return 0
	}
	return res / cnt
}
