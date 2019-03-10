# resize

This package lets you to resize image using LERP algorism.

## Quick start

```go
package main

import (
	"image"
	"image/png"
	"os"

	"github.com/po3rin/resize"
)

func main() {
	img, _, _ := image.Decode(os.Stdin)
	dst := resize.Resize(img, 2, 2)
	png.Encode(os.Stdout, dst)
}
```
