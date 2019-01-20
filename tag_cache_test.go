package validator

import (
	"sync"
	"testing"
)

func TestTagCache_StoreParallel(t *testing.T) {
	var wg sync.WaitGroup
	tc := newTagCache()
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tc.Store("test", &tagChunk{})
		}()
	}
	wg.Wait()
}
