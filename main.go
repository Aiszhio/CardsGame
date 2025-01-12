package main

import (
	"cardsgame/deck"
)

func main() {
	var playDeck deck.Deck
	playDeck.Create()
	playDeck.ShuffleAndTake()
	playDeck.Currents.Print()
	playDeck.Queue.Print()
}
