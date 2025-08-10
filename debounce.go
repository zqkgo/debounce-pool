package debounce

import (
	"sync"
	"time"

	"github.com/bep/debounce"
)

type meta struct {
	fn  func(func())
	exp time.Time
}

// debouncePool 实现防抖处理函数的缓存和自动清理。
type debouncePool struct {
	mu sync.Mutex
	// key -> meta
	debounces sync.Map
	ttlMs     int
}

func NewPool(ttlMs int) *debouncePool {
	if ttlMs == 0 {
		ttlMs = 5000
	}
	p := &debouncePool{
		ttlMs: ttlMs,
	}
	go p.cleanLoop()
	return p
}

func (dp *debouncePool) Get(key string, after time.Duration) func(func()) {
	v, ok := dp.debounces.Load(key)
	if ok {
		return v.(meta).fn
	}

	dp.mu.Lock()
	defer dp.mu.Unlock()
	// 二次检查。
	v, ok = dp.debounces.Load(key)
	if ok {
		return v.(meta).fn
	}
	d := debounce.New(after)
	dp.debounces.Store(key, meta{
		fn:  d,
		exp: time.Now().Add(time.Millisecond * time.Duration(dp.ttlMs)),
	})
	return d
}

func (dp *debouncePool) cleanLoop() {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for range t.C {
		dp.debounces.Range(func(key, value any) bool {
			m := value.(meta)
			if time.Now().After(m.exp) {
				dp.debounces.Delete(key)
			}
			return true
		})

	}
}
