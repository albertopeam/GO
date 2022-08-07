package main

import (
	"fmt"
	"strings"
	"golang.org/x/tour/wc"
	"golang.org/x/tour/pic"
	"math"
)

type Vertex struct { // A struct is a collection of fields.
	X int
	Y int
}

type Location struct {
	Lat, Long float64
}

func main() {
	pointers()

	structs()
	pointersToStructs()
	structLiterals()

	arrays()

	slices()
	slicesLikeArrays()
	sliceLiterals()
	sliceDefaults()
	sliceLengthAndCapacity()
	slicesNil()
	slicesMake()
	sliceOfSlices()
	sliceAppend()

	rangeStatement()
	rangeStatementSkip()

	pic.Show(Pic) // EXERCISE

	mapStatement()
	mapLiterals()
	mapLiteralsSimplified()
	mutatingMaps()

	wc.Test(WordCount) // EXERCISE

	functionValues()
	functionClosures()

	fibonacciExercise() // EXERCISE
}

func pointers() { // Go has pointers. A pointer holds the memory address of a value.
	i, j := 42, 2701
	var x *int
	fmt.Println(x)		// nil

	x = &i				// point to i. The & operator generates a pointer to its operand.
	p := &i         	// point to i
	fmt.Println(*p, *x, p, x, &p, &x) // read i through the pointer p & x. address that p & x points(same addr). addresses of p and x(different addr)

	*p = 21         	// set i through the pointer. The * operator denotes the pointer's underlying value.
	fmt.Println(i, *p)	// see the new value of i

	p = &j         		// point to j
	*p = *p / 37   		// divide j through the pointer
	fmt.Println(j) 		// see the new value of j
}

func structs() {
	vertex := Vertex{1, 2}
	fmt.Println(vertex)

	vertex.X = 4		// Struct fields are accessed using a dot
	fmt.Println(vertex.X)
}

func pointersToStructs() {
	v := Vertex{1, 2}
	p := &v
	(*p).X = 5	// mutate struct value via pointer
	fmt.Println(v)
	p.X = 1e9	// the language permits us instead to write just p.X, without the explicit dereference.
	fmt.Println(v)
}

func structLiterals() {
	var (
		v1 = Vertex{1, 2}  // has type Vertex
		v2 = Vertex{X: 1}  // Y:0 is implicit
		v3 = Vertex{}      // X:0 and Y:0
		v4 = Vertex{Y: 5}  // X:0 and Y:5
		p  = &Vertex{1, 2} // has type *Vertex. The special prefix & returns a pointer to the struct value.
	)
	fmt.Println(v1, p, v2, v3, v4)
}

func arrays() { // Arrays are fized size
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)
}

func slices() {
	primes := [6]int{2, 3, 5, 7, 11, 13}

	var s []int = primes[1:4] // This selects a half-open range which includes the first element, but excludes the last one.
	fmt.Println(s)
}

func slicesLikeArrays() { // Slices are like references to arrays
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	a := names[0:2]	// A slice does not store any data, it just describes a section of an underlying array.
	b := names[1:3]	// Other slices that share the same underlying array will see those changes.

	fmt.Println(a, b)

	b[0] = "XXX" // Changing the elements of a slice modifies the corresponding elements of its underlying array.
	fmt.Println(a, b)
	fmt.Println(names)
}

func sliceLiterals() { // Slice literals
	q := []int{2, 3, 5, 7, 11, 13}	// A slice literal is like an array literal without the length.
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	// inline definition + initialization
	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s)
}

func sliceDefaults() { // When slicing, you may omit the high or low bounds to use their defaults instead.
	s := []int{2, 3, 5, 7, 11, 13}

	s = s[:] // 0:n-1, same as original slice
	fmt.Println(s)

	s = s[1:4]
	fmt.Println(s)

	s = s[:2] // The default is zero for the low bound and the length of the slice for the high bound.
	fmt.Println(s)

	s = s[1:] // The default is zero for the low bound and the length of the slice for the high bound.
	fmt.Println(s)
}

func sliceLengthAndCapacity() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)

	// panic: runtime error: slice bounds out of range [:6] with capacity 4
	// s = s[:6]	
	// printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func slicesNil() { // A nil slice has a length and capacity of 0 and has no underlying array.
	var s []int
	fmt.Println(s, len(s), cap(s))
	if s == nil {
		fmt.Println("nil!")
	}
}

func slicesMake() { // Slices can be created with the built-in make function; this is how you create dynamically-sized arrays.
	a := make([]int, 5)	// The make function allocates a zeroed array and returns a slice that refers to that array
	printSliceNamed("a", a)

	b := make([]int, 0, 5)	// To specify a capacity, pass a third argument to make:
	printSliceNamed("b", b)

	c := b[:2]
	printSliceNamed("c", c)

	d := c[2:5]
	printSliceNamed("d", d)
}

func printSliceNamed(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

func sliceOfSlices() { // Slices can contain any type, including other slices.
	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}

func sliceAppend() { // https://pkg.go.dev/builtin#append
	var s []int
	printSlice(s)

	// append works on nil slices.
	s = append(s, 0)
	printSlice(s)

	// The slice grows as needed.
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 2, 3, 4)
	printSlice(s)	

	// Append two slices
	p := []int{5,6, 7}
	s = append(s, p...) // If the backing array of s is too small to fit all the given values a bigger array will be allocated. The returned slice will point to the newly allocated array.
	printSlice(s)
}

func rangeStatement() { // The range form of the for loop iterates over a slice or map.
	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
	for i, v := range pow { // When ranging over a slice, two values are returned for each iteration. The first is the index, and the second is a copy of the element at that index.
		fmt.Printf("2**%d = %d\n", i, v)
	}

	items := []int{0, 1, 2, 3, 4}
	items = items[1:]
	for i, item := range items { // Item is a copy and can not mutate items underlaying data
		fmt.Printf("index %d; value %d\n", i, item)
	}
}

func rangeStatementSkip() {
	var items = []int{0, 1, 2, 3, 4}
	for _, e := range items { // Can be skipped both index or value
		fmt.Println(e)
	} 

	pow := make([]int, 10)	// Create slice
	for i := range pow {
		pow[i] = 1 << uint(i) // == 2**i	// Insert
	}
	fmt.Println(pow)
}

// Implement Pic. It should return a slice of length dy, each element of which is a slice of dx 8-bit unsigned integers. 
// When you run the program, it will display your picture, interpreting the integers as grayscale (well, bluescale) values. 
// The choice of image is up to you. Interesting functions include (x+y)/2, x*y, and x^y.
func Pic(dx, dy int) [][]uint8 {
	columns := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		rows := make([]uint8, dx)
		for x := 0; x < dx; x++ {
			rows[x] = uint8((x+y) / 2)
		}
		columns[y] = rows
	}
	return columns
}

func mapStatement() { // Map data structure
	var m map[string]Location	
	m = make(map[string]Location) // The make function returns a map of the given type, initialized and ready for use.
	m["Bell Labs"] = Location {
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])
	fmt.Println(m["Bell"]) // Not available key
	fmt.Println(m) // Print the entire map

	if value, exists := m["Bell Labs"]; exists { // Checks for existence and get its value
		fmt.Println("`Bell Labs` exists and its value is ", value)
	}
}

func mapLiterals() { // Map literals are like struct literals, but the keys are required.
	var m = map[string]Location{
		"Bell Labs": Location{40.68433, -74.39967},
		"Google": Location{37.42202, -122.08408},
	}
	fmt.Println(m)
}


func mapLiteralsSimplified() { // If the top-level type is just a type name, you can omit it from the elements of the literal.
	var m = map[string]Location{
		"Bell Labs": {40.68433, -74.39967},
		"Google": {37.42202, -122.08408},
	}
	fmt.Println(m)
}

func mutatingMaps() {
	m := make(map[string]int)

	m["Answer"] = 42 // Insert an element in map
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48 // Update an element in map
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer") // Delete an element in map
	fmt.Println("The value:", m["Answer"])

	// If elem or ok have not yet been declared you could use a short declaration form: `value, ok := m[...]``
	value, ok := m["Answer"] // If key is in m, ok is true. If not, ok is false. If key is not in the map, then elem is the zero value for the map's element type.
	fmt.Println("The value:", value, "Present?", ok)
}

//Implement WordCount. It should return a map of the counts of each “word” in the string s.
func WordCount(s string) map[string]int {
	wordsCount := make(map[string]int)

	words := strings.Split(s, " ") // doc https://pkg.go.dev/strings#Split
	
	for _ , word := range words {		
		if value, ok := wordsCount[word]; ok {
			wordsCount[word] = value + 1
		} else {
			wordsCount[word] = 1
		}
	}
	
	return wordsCount
}

func functionValues() { // Functions are values too. They can be passed around just like other values. Function values may be used as function arguments and return values.
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(5, 12))

	fmt.Println(compute(hypot))
	fmt.Println(compute(math.Pow))

	swap := func(a, b int) (sa, sb int) {
		sa, sb = b, a
		return
	}
	fmt.Println(swap(1, 5))

	addOne := func(base int) int {
		return base + 1
	}
	fmt.Println(basePlusOne(addOne))
}

func basePlusOne(function func(int) int) int {
	var value int = 5	
	return function(value)
}

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func functionClosures() { // A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables.
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

// Let's have some fun with functions. Implement a fibonacci function that returns a function (a closure) that returns successive fibonacci numbers (0, 1, 1, 2, 3, 5, ...).
func fibonacciExercise() {
	f := fibonacci()
	var result []int
	for i := 0; i < 10; i++ {
		result = append(result, f())
	}
	fmt.Println(result)

	f2 := fibonacci2()
	var result2 []int
	for i := 0; i < 10; i++ {
		result2 = append(result2, f2())
	}
	fmt.Println(result2)

	f3 := fibonacci3()
	var result3 []int
	for i := 0; i < 10; i++ {
		result3 = append(result3, f3())
	}
	fmt.Println(result3)

	fmt.Println(fibonacciRecursive(9))
}

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	previous := 0
	current := 0
	index := 0
	return func() int { // depending on the value we compute fibonnaci until that value
		sum := previous + current
		switch index {
			case 0:
				previous = 0
				current = 1
				index += 1
			case 1: 
				previous = 0
				current = 1
				index += 1
			case 2:
				previous = 1
				current = 1
				index += 1
			default:
				previous = current
				current = sum
		}		
		return sum
	}
}

func fibonacci2() func() int {
	previous := 0
	current := 0
	index := 0
	return func() int { // depending on the value we compute fibonnaci until that value
		sum := previous + current
		switch index {
			case 0:
				index += 1
				return 0
			case 1: 				
				index += 1
				return 1
			case 2:
				previous = 1
				current = 1
				index += 1
				return 1
			default:
				previous = current
				current = sum
		}		
		return sum
	}
}

func fibonacci3() func() int {
	result := [2]int{0, 1}
	index := 0
	incrIndex := func() {
		index += 1
	}
	return func() int { // depending on the value we compute fibonnaci until that value
		defer incrIndex()
		sum := 0		
		if index < 2 {
			return result[index]
		} else {
			sum = result[0] + result[1]
			result[index % 2] = sum			
		}
		return sum
	}
}

func fibonacciRecursive(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
	}
}