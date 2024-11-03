package main

import (
	"cardsgame/deck"
)

func main() {
	var playDeck deck.Deck
	playDeck.Create()
	playDeck.Print()
	deck.QueueDeck.Print()
}
