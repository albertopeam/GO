package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	sut := newDeck()

	if len(sut) != 16 {
		t.Errorf("newDeck should contain 16 cards but it contains %v", len(sut))
	}
	if sut[0] != "Ace of ♣" {
		t.Errorf("newDeck first card is not Ace of ♣, is %v", sut[0])
	}
	if sut[len(sut)-1] != "Four of ♦" {
		t.Errorf("newDeck last card is not Four of ♦, is %v", sut[len(sut)-1])
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	var filename = "_deck.test"
	os.Remove(filename)

	sut := newDeck()
	sut.saveToFile(filename)
	d := newDeckFromFile(filename)
	if len(sut) != len(d) {
		t.Errorf("saveToFile and newDeckFromFile doesn't have the same length %v != %v", len(sut), len(d))
	}

	os.Remove(filename)
}
