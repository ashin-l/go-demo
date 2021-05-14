package main

import (
	"fmt"
)

func main() {
	const (
		a = iota + 1
		b
		c = iota
		d
	)

	const (
		z = "bbb"
		e = "aaa"
		f = iota
		g
	)
	fmt.Println(a, b, c, d)
	fmt.Println(z, e, f, g)
}
