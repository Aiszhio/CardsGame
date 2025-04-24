package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	ID   string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients = make(map[string]*Client)
	mu      sync.RWMutex
)

func register(c *Client) {
	mu.Lock()
	clients[c.ID] = c
	mu.Unlock()
}

func unregister(c *Client) {
	mu.Lock()
	delete(clients, c.ID)
	mu.Unlock()
}

func getClient(id string) *Client {
	mu.RLock()
	cl := clients[id]
	mu.RUnlock()
	return cl
}

func WebSocketHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	clientID := c.Param("id")
	client := &Client{
		Conn: ws,
		ID:   clientID,
		Send: make(chan []byte, 8),
	}

	register(client)

	log.Printf("WebSocket connected: %s", clientID)

	go ReadPump(client)
	go WritePump(client)
}

func ReadPump(c *Client) {
	defer func() {
		unregister(c)
		err := c.Conn.Close()
		if err != nil {
			return
		}
		log.Printf("WebSocket disconnected: %s", c.ID)
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func WritePump(c *Client) {
	defer func() {
		unregister(c)
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

func push(id string, payload any) {
	cl := getClient(id)
	if cl == nil {
		log.Printf("[WS] push to %q failed: client not found\n", id)
		return
	}
	b, _ := json.Marshal(payload)
	log.Printf("[WS] → %s: %s\n", id, string(b))
	select {
	case cl.Send <- b:
	default:
		log.Printf("[WS] buffer full, drop message for %s\n", id)
	}
}
