package deck

var QueueDeck Deck
var MainDeck Deck

func (d *Deck) DealDeck() Deck {
	for idx, card := range *d {
		if idx <= 8 {
			MainDeck = append(MainDeck, card)
		} else {
			QueueDeck = append(QueueDeck, card)
		}
	}
	*d = MainDeck
	return *d
}
