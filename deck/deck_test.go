package deck

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	var testDeck Deck
	testDeck.Create()

	if len(testDeck) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(testDeck))
	}
}

func TestShuffle(t *testing.T) {
	var testDeck Deck

	testDeck.Create()
	testDeck.Shuffle()

	if len(testDeck) != 36 {
		t.Error("Expected 36 cards in deck, but got ", len(testDeck))
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	var testDeck Deck
	os.Remove("_decktesting")

	testDeck.Create()
	testDeck.SaveToFile("_decktesting")

	loadedDeck := NewDeckFromFile("_decktesting")

	if len(loadedDeck) != 8 {
		t.Error("Expected 8 cards in deck, but got ", len(loadedDeck))
	}

	os.Remove("_decktesting")
}
