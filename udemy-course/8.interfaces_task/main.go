package main

import "fmt"

type triangle struct {
	height float64
	base   float64
}
type square struct {
	sideLength float64
}
type shape interface {
	getArea() float64
}

func (t triangle) getArea() float64 {
	return t.base * t.height / 2.0
}

func (s square) getArea() float64 {
	return s.sideLength * s.sideLength
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}

func main() {
	t := triangle{height: 5, base: 2}
	printArea(t)

	s := square{sideLength: 5}
	printArea(s)
}
