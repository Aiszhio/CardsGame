package main

import (
	"cardsgame/handlers"
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	var testDeck handlers.Deck
	testDeck.Create()

	if len(testDeck.Currents.Cards)+len(testDeck.Queue.Cards) != 36 {
		t.Error("Expected 36 cards, but got ", len(testDeck.Currents.Cards)+len(testDeck.Queue.Cards))
	}

}

func TestShuffle(t *testing.T) {
	var testDeck handlers.Deck

	testDeck.Create()
	testDeck.ShuffleAndTake()

	if len(testDeck.Currents.Cards) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(testDeck.Currents.Cards))
	}

	if len(testDeck.Queue.Cards) != 28 {
		t.Error("Expected 28 cards in queue, but got ", len(testDeck.Queue.Cards))
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	var testDeck handlers.Deck

	testDeck.Create()
	testDeck.ShuffleAndTake()

	err := handlers.SaveToFile("_decktesting", testDeck.Currents)

	loadedDeck := handlers.NewDeckFromFile("_decktesting")

	if len(loadedDeck.Currents.Cards) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(loadedDeck.Currents.Cards))
	}

	err = os.Remove("_decktesting")
	if err != nil {
		t.Error(err)
	}
}
