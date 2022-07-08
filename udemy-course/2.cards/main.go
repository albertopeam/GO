package main

func main() {
	cards := newDeck()
	cards.shuffle()
	cards.print()

	cards.saveToFile("my-deck")
	diskCards := newDeckFromFile("my-deck")
	diskCards.print()

	hand, remainingCards := deal(cards, 3)
	hand.print()
	remainingCards.print()
}
