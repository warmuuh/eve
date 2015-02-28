minimal experimental event bus for dispatching non-blocking events to listeners.

sample usage:

```go
package main

import (
	"github.com/warmuuh/eve"
	"log"
	"time"
)

func main() {

	b := eve.Bus()

	go func() {
		for {
			evt := <-b.From("start")
			log.Println("1Event received: ", evt)
		}
	}()

	b.To("start") <- "test"

	time.Sleep(2 * time.Second)
}
```