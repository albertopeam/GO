package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	forLoop()
	optionalForLoop()
	while()

	// infinitLoop()
	ifStatement(true)

	fmt.Println(pow(2, 4, 8))
	fmt.Println(pow2(2, 8))
	fmt.Println("sqrt(12)", Sqrt(100))

	switchStatement()
	switchEvaluationOrder()
	switchNoCondition()

	deferStatement()
	deferMultiple()
}

func forLoop() {
	sum := 0
	for i := 0; i < 10; i++ { // for loop
		sum += i
	}
	fmt.Println(sum)
}

func optionalForLoop() { // optional init/post statements
	sum := 1
	for ; sum < 1000; {
		sum += sum
	}
	fmt.Println(sum)
}

func while() { // while using a for loop. Go only has for, no while nor do while
	var sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)
}

func infinitLoop() { // infinite loop with for statement
	for {
		fmt.Println("Loop")
	}
}

func ifStatement(input bool) { // if statement
	if input {
		fmt.Println("TRUE")
	} else {
		fmt.Println("FALSE")
	}
}

func pow(x, n, lim float64) float64 { // if statement with local if var
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func pow2(x, lim float64) float64 { // if statement local var v can be accesed in else
	if v := math.Pow(x, 2); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here, though
	return lim
}

// As a way to play with functions and loops, let's implement a square root function: given a number x, we want to find the number z for which z² is most nearly x.
func Sqrt(x float64) float64 {
	var max = x
	var min = 0.0
	var sqrt float64
	const tolerance = 0.0001
	for {
		sqrt = (max - min) / 2.0 + min
		found := sqrt*sqrt
		if found < x {
			min  = sqrt
		} else {
			max = sqrt
		}
		if found >= x - tolerance && found <= x + tolerance {
			break
		}
	}
	return sqrt
}

func switchStatement() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os { // only executes the concrete case, not all the cases that follow.
	case "darwin":
		fmt.Println("OS X.") 
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}

func switchEvaluationOrder() {	
	today := time.Now().Weekday()
	fmt.Println("Today is ", today)
	fmt.Println("When's Saturday?")
	switch time.Saturday { // evaluates from top to bottom, stopping when a case succeeds
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}
}

func switchNoCondition() { // Switch without a condition is the same as switch true. This construct can be a clean way to write long if-then-else chains.
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

func deferStatement() { //A defer statement defers the execution of a function until the surrounding function returns. The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.
	tmp := 1
	defer fmt.Println("Deferred", tmp) // args evaluated before defer, but execution not triggered until body gets executed
	tmp += 1
	fmt.Println("Function deferStatement body")
}

func deferMultiple() { // Deferred function calls are pushed onto a stack. When a function returns, its deferred calls are executed in last-in-first-out order.
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i) // FIFO
	}

	fmt.Println("done")
}