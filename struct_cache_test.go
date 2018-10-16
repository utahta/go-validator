package validator

import (
	"reflect"
	"sync"
	"testing"
)

func TestStructCache_StoreParallel(t *testing.T) {
	var wg sync.WaitGroup
	sc := newStructCache()
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Store(reflect.TypeOf("test"), []fieldCache{})
		}()
	}
	wg.Wait()
}
