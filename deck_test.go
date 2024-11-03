package main

import (
	"cardsgame/deck"
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	var testDeck deck.Deck
	testDeck.Create()

	if len(testDeck) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(testDeck))
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	var testDeck deck.Deck
	os.Remove("_decktesting")

	testDeck.Create()
	testDeck.SaveToFile("_decktesting")

	loadedDeck := deck.NewDeckFromFile("_decktesting")

	if len(loadedDeck) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(loadedDeck))
	}

	os.Remove("_decktesting")
}
