package main

import (
	"fmt"

	"vintz.fr/nonogram/level"
)

func main() {
	real_main()
}

func real_main() {
	lvl := level.NewDefault()
	fmt.Println(lvl)
}
