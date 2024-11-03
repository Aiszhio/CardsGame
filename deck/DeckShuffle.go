package deck

import (
	"math/rand"
	"time"
)

func (d *Deck) Shuffle() Deck {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for i := 0; i < len(*d); i++ {
		temp := rand.Intn(len(*d) - i)
		(*d)[i], (*d)[temp] = (*d)[temp], (*d)[i]
	}
	return *d
}
