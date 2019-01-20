package validator

import (
	"sync"
	"sync/atomic"
)

type (
	tagCache struct {
		mux sync.Mutex
		v   atomic.Value
	}
)

func newTagCache() *tagCache {
	c := tagCache{}
	c.v.Store(make(map[string]*tagChunk))
	return &c
}

func (c *tagCache) Load(k string) (*tagChunk, bool) {
	v, ok := c.v.Load().(map[string]*tagChunk)[k]
	return v, ok
}

func (c *tagCache) Store(k string, chunk *tagChunk) {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.Load(k)
	if ok {
		return
	}

	tmp := c.v.Load().(map[string]*tagChunk)
	m := make(map[string]*tagChunk, len(tmp)+1)
	for k, v := range tmp {
		m[k] = v
	}
	m[k] = chunk
	c.v.Store(m)
}
