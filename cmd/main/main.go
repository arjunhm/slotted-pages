package main

import (
	"fmt"
	"page"
)

func main() {
	fmt.Println("hello")
	p := page.AllocPage(1)
	fmt.Println(p.Data[4:8])
}
