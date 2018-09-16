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
	c.v.Store(make(map[string][]Tag))
	return &c
}

func (c *tagCache) Load(k string) ([]Tag, bool) {
	v, ok := c.v.Load().(map[string][]Tag)[k]
	return v, ok
}

func (c *tagCache) Store(k string, tags []Tag) {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.Load(k)
	if ok {
		return
	}

	tmp := c.v.Load().(map[string][]Tag)
	m := make(map[string][]Tag, len(tmp)+1)
	for k, v := range tmp {
		m[k] = v
	}
	m[k] = tags
	c.v.Store(m)
}
