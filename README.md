# Example

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zqkgo/debounce-pool"
)

func main() {
	p := debounce.NewPool(1000)
	call1 := p.Get("key1", time.Second)
	call2 := p.Get("key2", 2*time.Second)

	var n1 int
	for i := 0; i < 100; i++ {
		call1(func() {
			n1++
		})
	}

	var n2 int
	for i := 0; i < 100; i++ {
		call2(func() {
			n2++
		})
	}

	// wait for call1 to actually execute
	time.Sleep(time.Second)
	fmt.Println(n1 == 1) // true
	fmt.Println(n2 == 1) // false

	// wait for call2 to actually execute
	time.Sleep(2 * time.Second)
	fmt.Println(n2 == 1) // true
}
```