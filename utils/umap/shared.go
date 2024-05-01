package umap

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m map[string]any
}

// ShardMap used for segmentation cached data
type ShardMap []*Shard

func NewSharedMap(n int) ShardMap {
	shards := make([]*Shard, n)
	for i := 0; i < n; i++ {
		shards[i] = &Shard{m: make(map[string]any)}
	}
	return shards
}

func (m ShardMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17])
	return hash % len(m)
}

func (m ShardMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardMap) Get(key string) any {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardMap) Set(key string, val any) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = val
}

func (m ShardMap) Keys() []string {
	var (
		keys = make([]string, 0)
		mu   = sync.RWMutex{}
		wg   = sync.WaitGroup{}
	)

	wg.Add(len(m))

	for _, shard := range m {
		go func(sh *Shard) {
			sh.RLock()

			for key := range sh.m {
				mu.Lock()
				keys = append(keys, key)
				mu.Unlock()
			}

			sh.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return keys
}
