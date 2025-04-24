package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"cardsgame/handlers"
	"github.com/gin-gonic/gin"
)

// totalCards возвращает общее количество карт в Deck (рука + очередь)
func totalCards(d handlers.Deck) int {
	return len(d.Currents.Cards) + len(d.Queue.Cards)
}

// TestCreate_FullDeck проверяет, что Create() создаёт полный набор карт без дубликатов
func TestCreate_FullDeck(t *testing.T) {
	var d handlers.Deck
	d.Create()

	expected := len(handlers.Suits) * len(handlers.Ranks)
	if totalCards(d) != expected {
		t.Fatalf("Create(): ожидалось %d карт, получили %d", expected, totalCards(d))
	}

	seen := make(map[string]bool)
	for _, c := range d.Queue.Cards {
		key := c.Rank + "-" + c.Suit
		if seen[key] {
			t.Errorf("Create(): дублирование карты %s", key)
		}
		seen[key] = true
	}
}

// TestShuffleAndTake_Distribution проверяет, что после Shuffle() и раздачи по 8
// карт игроку и ИИ, в руке каждого по 8, а очередь уменьшилась на 16
func TestShuffleAndTake_Distribution(t *testing.T) {
	var d handlers.Deck
	d.Create()
	d.Shuffle()
	handlers.CardsDistribution(&d, &d)

	if len(d.Currents.Cards) != 8 {
		t.Errorf("CardsDistribution: ожидалось 8 карт в руке, получили %d", len(d.Currents.Cards))
	}

	left := len(handlers.Suits)*len(handlers.Ranks) - 16
	if len(d.Queue.Cards) != left {
		t.Errorf("CardsDistribution: ожидалось %d карт в очереди, получили %d", left, len(d.Queue.Cards))
	}
}

// TestSaveGameStateHandler проверяет эндпоинт POST /api/deck/save:
// 1) создаёт файл со стейтом игры;
// 2) возвращает ссылку downloadUrl;
// 3) NewDeckFromFile успешно парсит этот файл и восстанавливает карты.
func TestSaveGameStateHandler(t *testing.T) {
	state := handlers.GameState{
		HandCards:  []string{"Ace of Spades", "King of Hearts"},
		AICards:    []string{"Queen of Clubs"},
		TableCards: []string{"Ten of Diamonds", "Nine of Clubs"},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(state)
	c.Request = httptest.NewRequest("POST", "/api/deck/save", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.SaveGameState(c)

	if w.Code != http.StatusOK {
		t.Fatalf("SaveGameState: ожидался статус 200, получили %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("SaveGameState: не удалось распарсить JSON: %v", err)
	}
	url, ok := resp["downloadUrl"]
	if !ok || url == "" {
		t.Fatalf("SaveGameState: пустой downloadUrl в ответе")
	}

	filename := filepath.Base(url)
	path := filepath.Join("static", "saved", filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Чтение сгенерированного файла не удалось: %v", err)
	}

	if !bytes.Contains(data, []byte(state.HandCards[0])) {
		t.Errorf("Файл не содержит ожидаемую карту %s", state.HandCards[0])
	}

	hand, ai, table, err := handlers.NewDeckFromFile(path)
	if err != nil {
		t.Fatalf("NewDeckFromFile: ошибка загрузки файла: %v", err)
	}
	if len(hand) != len(state.HandCards) {
		t.Errorf("NewDeckFromFile: ожидалось %d карт в руке, получили %d", len(state.HandCards), len(hand))
	}
	if len(ai) != len(state.AICards) {
		t.Errorf("NewDeckFromFile: ожидалось %d карт ИИ, получили %d", len(state.AICards), len(ai))
	}
	if len(table) != len(state.TableCards) {
		t.Errorf("NewDeckFromFile: ожидалось %d карт на столе, получили %d", len(state.TableCards), len(table))
	}

	err = os.Remove(path)
	if err != nil {
		return
	}
}

// TestNewDeckFromFile_Invalid проверяет, что при попытке
// загрузить несуществующий файл возвращается ошибка
func TestNewDeckFromFile_Invalid(t *testing.T) {
	_, _, _, err := handlers.NewDeckFromFile("nonexistent_file.txt")
	if err == nil {
		t.Fatal("NewDeckFromFile: ожидалась ошибка для несуществующего файла")
	}
}
