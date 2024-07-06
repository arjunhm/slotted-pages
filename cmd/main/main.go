package main

import (
	"fmt"
	"page"
)

func main() {
	fmt.Println("Slotted Pages Implementation")

	//f := page.NewFile("temp.gob")
	//fmt.Print(f)
	s := page.NewSlot(1, 0, 0, 0)
	fmt.Print(s)
}
