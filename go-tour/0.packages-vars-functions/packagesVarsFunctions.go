package main // Programs start running in package main

// This code groups the imports into a parenthesized, "factored" import statement.
import (
	"fmt"
	"math"
	"math/cmplx"
)

func add(x int, y int) int { // arg type after arg name
	return x + y
}

func add2(x, y int) int { // arg sharing type
	return x + y
}

func swap(x, y string) (string, string) { // return tuple
	return y, x
}

func split(sum int) (x, y int) { // named return values. they are treated as variables defined at the top of the function
	x = sum * 4 / 9
	y = sum - x
	return // A return statement without arguments returns the named return values. This is known as a "naked" return.
}

// module property. var <names> <type>
/*
 Variables declared without an explicit initial value are given their zero value.
The zero value is:
0 for numeric types,
false for the boolean type, and
"" (the empty string) for strings.
*/
var c, python, java bool // lowercase means not exported outside the module.

var j, k int = 1, 2 // initializers

// baasic types
var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

// constants, use the const keyword. Constants cannot be declared using the := syntax.
const Pi = 3.14

// Main
func main() {
	fmt.Printf("2ˆ7 = %g .\n", math.Exp2(7))

	// a name is exported if it begins with a capital letter. When importing a package, you can refer only to its exported names. Any "unexported" names are not accessible from outside the package.
	fmt.Println(math.Pi) 
	
	fmt.Println(add(3,52))
	fmt.Println(add2(3,52))

	a,b := swap("hello", "world") // asignattion
	fmt.Println(a, b)

	fmt.Println(split(17))

	// func property
	var i int
	fmt.Println(i, c, python, java)

	var goo, rust, kotlin = true, false, "kotlin"
	fmt.Println(j, k, goo, rust, kotlin)

	l := 3 // := short assignment statement can be used in place of a var declaration with implicit type. outside a function is not allowed
	fmt.Println(l)

	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	// assignment between items of different type requires an explicit conversion
	m := 42 // type inference
	n := float64(m)
	o := uint(n)
	fmt.Println(m, n, o)

	// const inside a function
	const World = "世界"
}