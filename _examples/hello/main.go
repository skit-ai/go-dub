package main

import (
	"fmt"

	"github.com/skit-ai/go-dub/audioop"
)

func main() {
	e := audioop.NewError("Hello, world: %d", 100)
	fmt.Println(e.Error())
}
