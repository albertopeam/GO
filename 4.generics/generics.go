package main

import "fmt"

func main() {
	generics()
	genericsExercise()
}

/*
Type parameters
Go functions can be written to work on multiple types using type parameters. The type parameters of a function appear between brackets, before the function's arguments.
func Index[T comparable](s []T, x T) int
This declaration means that s is a slice of any type T that fulfills the built-in constraint comparable. x is also a value of the same type.
comparable is a useful constraint that makes it possible to use the == and != operators on values of the type. In this example, we use it to compare a value to all slice elements until a match is found. This Index function works for any type that supports comparison.
*/
// Index returns the index of x in s, or -1 if not found.
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

func generics() {
	fmt.Println("generics------------------")
	// Index works on a slice of ints
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	// Index also works on a slice of strings
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))
}

/*
Generic types
In addition to generic functions, Go also supports generic types. A type can be parameterized with a type parameter, which could be useful for implementing generic data structures.
This example demonstrates a simple type declaration for a singly-linked list holding any type of value.
As an exercise, add some functionality to this list implementation.
*/

// List represents a singly-linked list that holds
// values of any type.
type List[T comparable] struct {
	next *List[T]
	val T
}

func (l *List[T]) Next() *List[T] {
	return l.next
}

/// Stringer interface
func (l *List[T]) String() (res string) {
	list := l // copy to not mutate original pointer
	for list != nil {
		next := list.Next()
		if next != nil {
			res += fmt.Sprintf("%v -> ", list.val)
		} else {
			res += fmt.Sprintf("%v", list.val)
		}
		list = next
	}
	return
}

func genericsExercise() {
	fmt.Println("generics exercises------------------")
	list3 := List[int]{next: nil, val: 3}
	list2 := List[int]{next: &list3, val: 2}
	list1 := List[int]{next: &list2, val: 1}
	list0 := List[int]{next: &list1, val: 0}
	fmt.Println(&list0) // usage of Stringer interface
	fmt.Println("Original list0", list0) // assert that list0 pointer hasn't been mutated
}