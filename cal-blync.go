package main

import (
	"fmt"
	"github.com/eddiewebb/goblync"
)

func main() {
	light := blync.NewBlyncLight()
	light.FlashOrder()
	fmt.Println("Flashing")
}