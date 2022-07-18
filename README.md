# browsolate

A badly-named tool for opening multiple isolated instances of Chrome. ALL AT THE SAME TIME.

## How do?

Use it from your shell

```shell
browsolate https://google.com
```

Or from your go

```go
package main

import (
	"fmt"
	"github.com/patricksanders/browsolate"
)

func main() {
	opts := &browsolate.InstanceOpts{ChromePath: "/usr/bin/chrome"}
	err := browsolate.StartIsolatedChromeInstance("https://google.com", opts)
	if err != nil {
		fmt.Print(err)
	}
}
```

(but you probably shouldn't because I just threw it together)
