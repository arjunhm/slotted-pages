package main

import (
	"fmt"
	"page"
)

func main() {
	fmt.Println("slotted pages")

	p := page.AllocPage(1)
	p.DisplayHeaderDetails()

	p.AddData("aged", "12")
	p.DisplayHeaderDetails()

	p.AddData("abc", "10")
	p.DisplayHeaderDetails()

	p.AddData("def", "99")
	p.DisplayHeaderDetails()

	p.AddData("xyz", "xz")
	p.DisplayHeaderDetails()

	fmt.Println(p.Data[12:60])
}
