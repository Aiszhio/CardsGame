package deck

import (
	"fmt"
	"os"
	"strings"
)

func fillList[T Currents | Queue](deck T) []string {
	var saveDeck []string

	switch d := any(deck).(type) {
	case Queue:
		saveDeck = append(saveDeck, fmt.Sprint("Here are all the queue cards, \n"))

		for idx, card := range d.Cards {
			saveDeck = append(saveDeck, fmt.Sprint(idx+1, " is ", card.Rank, card.Suit))
		}
	case Currents:
		saveDeck = append(saveDeck, fmt.Sprint("Here are all current cards \n"))

		for idx, card := range d.Cards {
			saveDeck = append(saveDeck, fmt.Sprint(idx+1, " is ", card.Rank, " ", card.Suit, "\n"))
		}
	default:
		fmt.Println("Неизвестный тип")
	}

	return saveDeck
}

func SaveToFile[T Currents | Queue](filename string, deck T) error {
	saveDeck := fillList(deck)

	stringDeck := strings.Join(saveDeck, "")

	err := os.WriteFile(filename, []byte(stringDeck), 0666)
	if err != nil {
		return err
	}

	return nil
}
