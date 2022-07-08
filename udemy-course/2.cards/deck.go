package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck []string

func newDeck() deck {
	cards := deck{}
	suits := [4]string{"♣", "♠", "♥", "♦"}
	values := []string{"Ace", "Two", "Three", "Four"}
	for _, suit := range suits {
		for _, value := range values {
			card := value + " of " + suit
			cards = append(cards, card)
		}
	}
	return cards
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

func (d deck) toString() string {
	return strings.Join([]string(d), ",") //https://pkg.go.dev/strings@go1.18.3
}

func (d deck) saveToFile(filename string) error {
	var data []byte = []byte(d.toString())
	var mode fs.FileMode = 0666
	return ioutil.WriteFile(filename, data, mode) //https://pkg.go.dev/io/ioutil@go1.18.3
}

func newDeckFromFile(filename string) deck {
	bytes, err := ioutil.ReadFile(filename) // https://pkg.go.dev/io/ioutil@go1.18.3#ReadFile
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	var s []string = strings.Split(string(bytes), ",") // https://pkg.go.dev/strings#Split
	return deck(s)
}

func (d deck) shuffle() {
	seed := time.Now().UnixNano() // https://pkg.go.dev/time#Time.UnixNano
	souce := rand.NewSource(seed) // https://pkg.go.dev/math/rand#NewSource
	r := rand.New(souce)          // https://pkg.go.dev/math/rand#New
	length := len(d)
	for i := range d {
		//newPosition := rand.Intn(length) //https: //pkg.go.dev/math/rand@go1.18.3#Intn
		newPosition := r.Intn(length) // https://pkg.go.dev/math/rand#Rand.Intn
		d[i], d[newPosition] = d[newPosition], d[i]
	}
}
