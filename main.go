package main

import (
	"cardsgame/deck"
	"fmt"
)

func main() {
	playDeck := deck.NewDeckFromFile("test2")
	for _, card := range playDeck {
		fmt.Println(card)
	}
}
