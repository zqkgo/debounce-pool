package debounce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDebouncePool(t *testing.T) {
	ast := require.New(t)
	p := NewPool(1000)

	after := 150 * time.Millisecond
	call := p.Get("key0", after)

	var n int
	for i := 0; i < 100; i++ {
		call(func() {
			n++
		})
	}
	time.Sleep(after)
	// 只执行一次。
	ast.Equal(1, n)

	p.Get("key0", after)
	call1 := p.Get("key1", 500*time.Millisecond)
	p.Get("key1", 500*time.Millisecond)
	call2 := p.Get("key2", time.Second)
	p.Get("key2", time.Second)

	// 多次调用，但缓存三个实例。
	var num int
	p.debounces.Range(func(key, value any) bool {
		t.Logf("key: %s, value: %v", key, value.(meta))
		num++
		return true
	})
	ast.Equal(3, num)

	var c1, c2 bool
	call1(func() {
		t.Log("call1 execute", time.Now())
		c1 = true
	})
	call2(func() {
		t.Log("call2 execute", time.Now())
		c2 = true
	})
	// 不同时间点执行。
	time.Sleep(600 * time.Millisecond)
	ast.True(c1)
	ast.False(c2)
	time.Sleep(time.Second)
	ast.True(c2)

	// 自动过期。
	num = 0
	p.debounces.Range(func(key, value any) bool {
		num++
		return true
	})
	ast.Equal(0, num)
}
