package main

import (
	"cardsgame/deck"
)

func main() {
	var playDeck deck.Deck
	playDeck.Create()
	playDeck.ShuffleAndTake()
	playDeck.Currents.Print()
	err := deck.SaveToFile("aboba", playDeck.Currents)
	if err != nil {
		panic(err)
	}
	_ = deck.NewDeckFromFile("aboba")
}
