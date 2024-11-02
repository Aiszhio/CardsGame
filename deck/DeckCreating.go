package deck

type Deck []string

var suits = []string{"Spades", "Hearts", "Diamonds", "Clubs"}
var ranks = []string{"Ace", "King", "Queen", "Jack", "Ten", "Nine", "Eight", "Seven", "Six"}

func Create() Deck {
	var playDeck Deck
	for _, suit := range suits {
		for _, rank := range ranks {
			playDeck = append(playDeck, rank+" of "+suit)
		}
	}
	return playDeck
}
