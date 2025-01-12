package deck

//
//import (
//	"log"
//	"os"
//	"strings"
//)
//
//func isCardValid(card string) bool {
//	var flagSuit bool
//	var flagRank bool
//	for _, suit := range suits {
//		if strings.Fields(card)[2] == suit {
//			flagSuit = true
//		}
//	}
//	for _, rank := range ranks {
//		if strings.Fields(card)[0] == rank {
//			flagRank = true
//		}
//	}
//	return flagSuit && flagRank
//}
//
//func NewDeckFromFile(filename string) Deck {
//	bs, err := os.ReadFile(filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var deck Deck
//	strSlice := strings.Fields(string(bs))
//	var temp string
//	for _, word := range strSlice {
//		temp += word + " "
//		if len(strings.Fields(temp)) == 3 && isCardValid(temp) {
//			deck = append(deck, temp)
//			temp = ""
//		}
//	}
//	return deck
//}
