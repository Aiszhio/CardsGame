package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	ID   string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	log.Println("Попытка подключения WebSocket")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	clientID := c.Param("id")
	log.Printf("Успешное подключение для ID: %s", clientID)

	client := &Client{
		Conn: ws,
		ID:   clientID,
		Send: make(chan []byte),
	}

	go ReadPump(client)
	go WritePump(client)
}

func ReadPump(c *Client) {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {

		}
	}(c.Conn)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Printf("Сообщение от пользователя %s: %s", c.ID, string(message))
	}
}

func WritePump(c *Client) {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(c.Conn)
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
