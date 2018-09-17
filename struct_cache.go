package validator

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	structCache struct {
		mux sync.Mutex
		v   atomic.Value
	}
)

func newStructCache() *structCache {
	c := structCache{}
	c.v.Store(make(map[reflect.Type][]fieldCache))
	return &c
}

func (c *structCache) Load(k reflect.Type) ([]fieldCache, bool) {
	v, ok := c.v.Load().(map[reflect.Type][]fieldCache)[k]
	return v, ok
}

func (c *structCache) Store(k reflect.Type, fields []fieldCache) {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.Load(k)
	if ok {
		return
	}

	tmp := c.v.Load().(map[reflect.Type][]fieldCache)
	m := make(map[reflect.Type][]fieldCache, len(tmp)+1)
	for k, v := range tmp {
		m[k] = v
	}
	m[k] = fields
	c.v.Store(m)
}
