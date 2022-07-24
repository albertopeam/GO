package main

import "fmt"

type person struct {
	name    string
	surname string
	contact contactInfo
	// `contactInfo`` -> equivalent to `contactInfo contactInfo`
}

type contactInfo struct {
	email   string
	zipCode int
}

func (p person) print() {
	fmt.Printf("%+v\n", p)
}

func (p person) updateName(name string) { // not works as `p` is passed as value
	p.name = name
}

func (p *person) updateNme(name string) {
	// (*p).name = name // explicit
	p.name = name // implicit
}

func updateSlice(s []string) { // s is mutated, under the hood is a pointer
	s[0] = "Bye"
}

func main() {
	var emptyContactInfo contactInfo

	notMe := person{"Alberto", "Sanchez", emptyContactInfo} // default initializer, we are forced to pass some contactInfo
	fmt.Println(notMe)

	me := person{name: "Alberto", surname: "Penas Amor"} // initializer with property names
	fmt.Println(me)

	var alex person // zero value initialization
	fmt.Println(alex)
	fmt.Printf("%+v\n", alex) // %+v forces to print property names

	var peter person
	peter.name = "Peter"
	peter.surname = "Jackson"
	fmt.Println(peter)

	josep := person{
		name:    "Josep",
		surname: "Taradellas",
		contact: contactInfo{
			email:   "josep@info.com",
			zipCode: 28080,
		},
	}
	josep.print()

	// mutated as value doesn't work
	josep.updateName("Mike")
	josep.print()

	// mutated as a reference work
	pointerToJosep := &josep
	pointerToJosep.updateNme("new josep")
	pointerToJosep.print()

	// shortcut/implicit pointer conversion from value to pointer, updateNme expects a pointer
	josep.updateNme("Yusepe")
	josep.print()

	// TIP to use pointers
	// turn VALUE into ADDRESS with &VALUE
	// turn ADDRESS into VALUE with *ADDRESS

	slice := []string{"Hi", "there", "!"}
	updateSlice(slice) // slice is a struct but is not passed as a value, passed as reference so it is mutated
	fmt.Println(slice)
}
