package main

import (
	"fmt"
	"math"
	"time"
	"io"
	"strings"
	"golang.org/x/tour/reader"
	"os"
	"image"
	"golang.org/x/tour/pic"	
	"image/color"
)

func main() {
	methodsOnStruct()
	methodsOnTypes()

	pointerReceivers()
	pointersAndFunctions()
	pointersVsValueReceiver()

	interfaces()
	interfacesWithUnderlyingNilValues()
	interfacesNil()
	interfacesEmpty()

	typeAssertions()
	typeSwitches()

	stringers()
	stringersExercise()

	errors()
	errorsExercise()

	readers()
	readersExercise()
	readersRoot13Exercise()

	images()
	imagesExercise()
}

/*
Go does not have classes. However, you can define methods on types.
A method is a function with a special receiver argument.
The receiver appears in its own argument list between the func keyword and the method name.
In this example, the Abs method has a receiver of type Vertex named v.
*/
type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func methodsOnStruct() {
	fmt.Println("methodsOnStruct---------------")
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
}

/*
You can declare a method on non-struct types, too.
In this example we see a numeric type MyFloat with an Abs method.
You can only declare a method with a receiver whose type is defined in the same package as the method. You cannot declare a method with a receiver whose type is defined in another package (which includes the built-in types such as int).
*/
type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func methodsOnTypes() {
	fmt.Println("methodsOnTypes---------------")
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
}

/*
You can declare methods with pointer receivers.
This means the receiver type has the literal syntax *T for some type T. (Also, T cannot itself be a pointer such as *int.)
For example, the Scale method here is defined on *Vertex.
Methods with pointer receivers can modify the value to which the receiver points (as Scale does here). Since methods often need to modify their receiver, pointer receivers are more common than value receivers.
Try removing the * from the declaration of the Scale function on line 16 and observe how the program's behavior changes.
With a value receiver, the Scale method operates on a copy of the original Vertex value. (This is the same behavior as for any other function argument.) The Scale method must have a pointer receiver to change the Vertex value declared in the main function.
*/
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v Vertex) ScaleNotMutating(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func ScaleVertexPointer(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func ScaleVertex(v Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func pointerReceivers() {
	fmt.Println("pointerReceivers---------------")
	v := Vertex{3, 4}
	//For the statement v.Scale(5), even though v is a value and not a pointer, the method with the pointer receiver is called automatically. 
	// That is, as a convenience, Go interprets the statement v.Scale(5) as (&v).Scale(5) since the Scale method has a pointer receiver.
	v.Scale(5)				// receiver is a pointer so it mutates original. shortcut can be used instead of &v
	(&v).Scale(5) 			// same than previous line. explicit
	fmt.Println(v.Abs())

	vnm := Vertex{3, 4}
	vnm.ScaleNotMutating(10) // without pointer the receiver is passed as value so copy is mutated but original not
	fmt.Println(vnm.Abs())

	vsp := Vertex{3, 4}
	ScaleVertexPointer(&vsp, 10)			// needs a pointer as arg
	// ScaleVertexPointer(vsp, 10)			// doesn't compile, expected pointer arg
	// ScaleVertex(&vsp, 10) 				// doesn't compile, expecter value arg
	fmt.Println(vsp.Abs())
}

/*
Here we see the Abs and Scale methods rewritten as functions.
*/
func Scale(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func pointersAndFunctions() {
	fmt.Println("pointersAndFunctions---------------")
	v := Vertex{3, 4}
	Scale(&v, 10)
	fmt.Println(Abs(v))

	v2 := Vertex{3, 4}
	fmt.Println(v2.Abs())
	fmt.Println(Abs(v2))

	p := &Vertex{4, 3}
	fmt.Println(p.Abs())
	fmt.Println((*p).Abs()) // same as previous line, explicit
	fmt.Println(Abs(*p))
	// fmt.Println(Abs(p)) // doesn't compile as type pointer doesn't match
}

/*
There are two reasons to use a pointer receiver.
The first is so that the method can modify the value that its receiver points to.
The second is to avoid copying the value on each method call. This can be more efficient if the receiver is a large struct, for example.
In this example, both Scale and Abs are with receiver type *Vertex, even though the Abs method needn't modify its receiver.
In general, all methods on a given type should have either value or pointer receivers, but not a mixture of both. (We'll see why over the next few pages.)
*/
func (v *Vertex) ScaleX(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v *Vertex) AbsX() float64 { // I wouldn't use a pointer if not mutating
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func pointersVsValueReceiver() {
	fmt.Println("pointersVsValueReceiver---------------")
	v := &Vertex{3, 4}
	fmt.Printf("Before scaling: %+v, Abs: %v\n", v, v.AbsX())
	v.ScaleX(5)
	fmt.Printf("After scaling: %+v, Abs: %v\n", v, v.AbsX())
}

/*
An interface type is defined as a set of method signatures.

A type implements an interface by implementing its methods. There is no explicit declaration of intent, no "implements" keyword.
Implicit interfaces decouple the definition of an interface from its implementation, which could then appear in any package without prearrangement.

Under the hood, interface values can be thought of as a tuple of a value and a concrete type: (value, type)
An interface value holds a value of a specific underlying concrete type.
Calling a method on an interface value executes the method of the same name on its underlying type.
*/
type Abser interface {
	Abs() float64
}

type Printable interface {
	Stdout()
}

type T struct {
	S string
}

// This method means type T implements the interface I,
// but we don't need to explicitly declare that it does so.
func (t T) Stdout() {
	fmt.Println(t.S)
}

func describe(i Abser) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func interfaces() {
	fmt.Println("interfaces---------------")
	var a Abser
	f := MyFloat(-math.Sqrt2)
	describe(f)
	v := Vertex{3, 4}
	describe(v)

	a = f  // a MyFloat implements Abser
	fmt.Println(a.Abs())
	a = &v // a *Vertex implements Abser
	fmt.Println(a.Abs())

	// In the following line, v is a Vertex (not *Vertex)
	// and does NOT implement Abser.
	// a = v

	var i Printable = T{"hello"}
	i.Stdout()
}

/*
If the concrete value inside the interface itself is nil, the method will be called with a nil receiver.
In some languages this would trigger a null pointer exception, but in Go it is common to write methods that gracefully handle being called with a nil receiver (as with the method M in this example.)
Note that an interface value that holds a nil concrete value is itself non-nil.
*/
type Printer interface{
	print()
}

type Data struct {
	Name string
}

func (data *Data) print() {
	if data == nil {
		fmt.Println("data is <nil>")
	} else {
		fmt.Printf("(%v, %T)\n", data, data)
	}
}

func interfacesWithUnderlyingNilValues() {
	fmt.Println("interfacesWithUnderlyingNilValues---------------")
	
	var p Printer

	// nil underlying value scenario
	var data *Data
	p = data // to handle nilable stuff we need to have a method with a pointer -> func (data *Data) print(), otherwise it won´t compile
	p.print()

	var data1 *Data
	data1 = &Data{Name: "Yoshie"}
	p = data1
	p.print()

	var data2 Data
	data2 = Data{Name: "Caren"}
	p = &data2
	p.print()
}

/*
A nil interface value holds neither value nor concrete type.
Calling a method on a nil interface is a run-time error because there is no type inside the interface tuple to indicate 
which concrete method to call.
*/
type I interface {
	M()
}

func interfacesNil() {
	fmt.Println("interfacesNil---------------")

	var i I	
	fmt.Printf("(%v, %T)\n", i, i)
	//i.M() // panic: runtime error: invalid memory address or nil pointer dereference
}

/*
The interface type that specifies zero methods is known as the empty interface:
interface{}
An empty interface may hold values of any type. (Every type implements at least zero methods.)
Empty interfaces are used by code that handles values of unknown type. 
For example, fmt.Print takes any number of arguments of type interface{}.
*/
func describeInterface(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i) // fmt takes args of any type interface{}
}

func interfacesEmpty() {
	fmt.Println("interfacesEmpty---------------")

	var x interface {}
	describeInterface(x)

	x = 1
	describeInterface(x)
}

/*
A type assertion provides access to an interface value's underlying concrete value.
t := i.(T)
This statement asserts that the interface value i holds the concrete type T and assigns the underlying T value to the variable t.
If i does not hold a T, the statement will trigger a panic.
To test whether an interface value holds a specific type, a type assertion can return two values: the underlying value and a boolean value that reports whether the assertion succeeded.
t, ok := i.(T)
If i holds a T, then t will be the underlying value and ok will be true.
If not, ok will be false and t will be the zero value of type T, and no panic occurs.
Note the similarity between this syntax and that of reading from a map.
*/
func typeAssertions() {
	fmt.Println("typeAssertions---------------")
	var i interface{} = "hello"

	s := i.(string)
	fmt.Println(s)

	f, ok := i.(float64) // type assertion returns a tuple to validate if the type exists. f is the data, ok is a bool
	fmt.Println(f, ok)

	// if type assertion is used without control tuple and fails it will throw a panic. so always use the tuple
	// f = i.(float64) // panic: interface conversion: interface {} is string, not float64
}

/*
A type switch is a construct that permits several type assertions in series.
A type switch is like a regular switch statement, but the cases in a type switch specify types (not values), and those values are compared against the type of the value held by the given interface value.
switch v := i.(type) {
case T:
    // here v has type T
case S:
    // here v has type S
default:
    // no match; here v has the same type as i
}
The declaration in a type switch has the same syntax as a type assertion i.(T), but the specific type T is replaced with the keyword type.
This switch statement tests whether the interface value i holds a value of type T or S. In each of the T and S cases, the variable v will be of type T or S respectively and hold the value held by i. In the default case (where there is no match), the variable v is of the same interface type and value as i.
*/
func do(i interface{}) {
	switch v:= i.(type) {	// type replaces concrete type cast
	case int: fmt.Printf("Twice %v is %v\n", v, v*2)
	case string: fmt.Printf("%q is %v bytes long\n", v, len(v))
	default: fmt.Printf("I don't know about type %T!\n", v)
	}
}

func typeSwitches() {
	fmt.Println("typeSwitches---------------")
	do(21)
	do("HI!")
	do(false)
}

/*
One of the most ubiquitous interfaces is Stringer defined by the fmt package.
type Stringer interface {
    String() string
}
A Stringer is a type that can describe itself as a string. The fmt package (and many others) look for this interface to print values.
*/
type Person struct {
	Name string
	Age int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)\n", p.Name, p.Age)
}

func stringers() {
	fmt.Println("stringers---------------")
	a := Person{"Arthur Dent", 42}
	var t Person = Person{"Zaphod Beeblebrox", 9001}
	var z *Person = &t
	fmt.Println(a, z)
}

/*
Make the IPAddr type implement fmt.Stringer to print the address as a dotted quad.
For instance, IPAddr{1, 2, 3, 4} should print as "1.2.3.4".
*/
type IPAddr [4]byte

func (addr IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3])
}

func stringersExercise() {
	hosts := map[string]IPAddr {
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}

/*
Go programs express error state with error values.
The error type is a built-in interface similar to fmt.Stringer:
type error interface {
    Error() string
}
(As with fmt.Stringer, the fmt package looks for the error interface when printing values.)
Functions often return an error value, and calling code should handle errors by testing whether the error equals nil.
i, err := strconv.Atoi("42")
if err != nil {
    fmt.Printf("couldn't convert number: %v\n", err)
    return
}
fmt.Println("Converted integer:", i)
A nil error denotes success; a non-nil error denotes failure.
*/
type MyError struct {
	When time.Time
	What string
}

// impl error interface
func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

// function that triggers an error as output
func run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}


func errors() {
	fmt.Println("errors---------------")
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

/*
Copy your Sqrt function from the earlier exercise and modify it to return an error value.

Sqrt should return a non-nil error value when given a negative number, as it doesn't support complex numbers.

Create a new type

type ErrNegativeSqrt float64
and make it an error by giving it a

func (e ErrNegativeSqrt) Error() string
method such that ErrNegativeSqrt(-2).Error() returns "cannot Sqrt negative number: -2".

Note: A call to fmt.Sprint(e) inside the Error method will send the program into an infinite loop. You can avoid this by converting e first: fmt.Sprint(float64(e)). Why?

Change your Sqrt function to return an ErrNegativeSqrt value when given a negative number.
*/
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	val := float64(e) // NEEDS TO explicitly be float64 to pass to fmt.Sprintf otherwise it will genereate an infinity loop as it will try to invoke Error again
	return fmt.Sprintf("negative number %v is not supported", val)
}

func Sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, ErrNegativeSqrt(n)
	} else {		
		return 1.41, nil //TODO: hardcoded result, not the purpose of the exercise
	}
}

func errorsExercise() {
	fmt.Println("errors exercise---------------")
	sqrt := func(n float64) {
		if result, err := Sqrt(n); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("square root of", n, "is", result)
		}
	} 
	sqrt(-2)
	sqrt(2)
}

/*
The io package specifies the io.Reader interface, which represents the read end of a stream of data.
The Go standard library contains many implementations of this interface, including files, network connections, compressors, ciphers, and others.
The io.Reader interface has a Read method:
func (T) Read(b []byte) (n int, err error)
Read populates the given byte slice with data and returns the number of bytes populated and an error value. It returns an io.EOF error when the stream ends.
The example code creates a strings.Reader and consumes its output 8 bytes at a time.
*/
func readers() {
	fmt.Println("readers---------------")
	reader := strings.NewReader("Hello, world!")
	buffer := make([]byte, 8)

	for {
		size, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		fmt.Printf("size = %v err = %v b = %v\n", size, err, buffer)
		fmt.Printf("buffer[:size] = %q\n", buffer[:size])		
	}
}

/*
Exercise: Readers
Implement a Reader type that emits an infinite stream of the ASCII character 'A'.
*/
type MyReader struct{}

func (m MyReader) Read(b []byte) (int, error) {
	length := len(b)
	bytes := []byte("A")
	for i := 0; i < length; i++ {
		b = append(b, bytes...)
	}
	return length, nil
}

 func readersExercise() {
	fmt.Println("readersExercise---------------")
	reader.Validate(MyReader{})
 }

 /*
 Exercise: rot13Reader
A common pattern is an io.Reader that wraps another io.Reader, modifying the stream in some way.
For example, the gzip.NewReader function takes an io.Reader (a stream of compressed data) and returns a *gzip.Reader that also implements io.Reader (a stream of the decompressed data).
Implement a rot13Reader that implements io.Reader and reads from an io.Reader, modifying the stream by applying the rot13 substitution cipher to all alphabetical characters.
The rot13Reader type is provided for you. Make it an io.Reader by implementing its Read method.
 */

 type rot13Reader struct {
	r io.Reader
}

// lowercase means private to this package
func rot13Chyper(b byte) byte {  
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}	
	// b-a+13: input - a offset, we transform from the letters positional system to our 0 based index positional system. We sum 13 to apply rot13 algorythm.
	// z-a+1: represent the size of our representation system(0-25, length 26). offset 0 based, so +1
	// %: we can't overflow our size(z-a+1), module makes that the letters that are greather than 14 returns to the alphabet start
	// + a: we need to return the new character to the original non based 0 positional system
	return (b-a+13)%(z-a+1) + a
}

func (rot13 rot13Reader) Read(b []byte) (n int, err error) {
	n, err = rot13.r.Read(b)
	for i := 0; i < n; i++ {
		b[i] = rot13Chyper(b[i])
	} 
	return 
}

 func readersRoot13Exercise() {
	fmt.Println("readersRoot13Exercise---------------")
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r) // https://pkg.go.dev/io#Copy
	fmt.Println()
 }

 /*
Package image defines the Image interface:

package image
type Image interface {
    ColorModel() color.Model
    Bounds() Rectangle
    At(x, y int) color.Color
}

Note: the Rectangle return value of the Bounds method is actually an image.Rectangle, as the declaration is inside package image.
(See the documentation for all the details.)
The color.Color and color.Model types are also interfaces, but we'll ignore that by using the predefined implementations color.RGBA and color.RGBAModel. These interfaces and types are specified by the image/color package
 */
 func images() {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(m.Bounds())
	fmt.Println(m.At(0, 0).RGBA())
 }

 /*
 Remember the picture generator you wrote earlier? Let's write another one, but this time it will return an implementation of image.Image instead of a slice of data.
Define your own Image type, implement the necessary methods, and call pic.ShowImage.
Bounds should return a image.Rectangle, like image.Rect(0, 0, w, h).
ColorModel should return color.RGBAModel.
At should return a color; the value v in the last picture generator corresponds to color.RGBA{v, v, 255, 255} in this one.
 */

 type Image struct{
	 X, Y int
	 Width, Height int
 }

 // image interface: https://pkg.go.dev/image#Image

 // ColorModel returns the Image's color model.
 func (i Image) ColorModel() color.Model {
	return color.GrayModel
 }
 	
// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i Image) Bounds() image.Rectangle {
	min := image.Point{X: i.X, Y: i.Y}
	max := image.Point{X: i.X + i.Width, Y: i.Y + i.Height}
	return image.Rectangle{Min: min, Max: max}
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (i Image) At(x, y int) color.Color {
	return color.Gray{}
}

 func imagesExercise() {
	m := Image{X: 0, Y: 0, Width: 100, Height: 100}
	pic.ShowImage(m) // https://pkg.go.dev/golang.org/x/tour/pic#ShowImage
 }