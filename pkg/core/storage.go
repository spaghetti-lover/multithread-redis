package core

import (
	"github.com/spaghetti-lover/multithread-redis/internal/data_structure/hash_table"
	"github.com/spaghetti-lover/multithread-redis/internal/data_structure/probabilistic"
	"github.com/spaghetti-lover/multithread-redis/internal/data_structure/simple_set"
	"github.com/spaghetti-lover/multithread-redis/internal/data_structure/sorted_set"
)

var dictStore *hash_table.Dict
var zsetStore map[string]*sorted_set.SortedSet
var setStore map[string]*simple_set.SimpleSet
var cmsStore map[string]probabilistic.FrequencyEstimator

func init() {
	dictStore = hash_table.CreateDict()
	zsetStore = make(map[string]*sorted_set.SortedSet)
	setStore = make(map[string]*simple_set.SimpleSet)
	cmsStore = make(map[string]probabilistic.FrequencyEstimator)
}
