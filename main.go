package main

import (
	"cardsgame/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowAllOrigins:  true,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	router.Group("/")
	{
		router.GET("/home", handlers.GetHome)
		router.POST("/api/deck/create", handlers.CreateDeck)
		router.GET("/api/ws/:id", handlers.WebSocketHandler)
		router.POST("/api/deck/leave", handlers.LeaveCards)
		router.POST("/api/deck/ai", handlers.AIResponse)
		router.POST("/api/deck/status", handlers.WatchStatus)
		router.POST("/api/deck/save", handlers.SaveGameState)
		router.GET("/api/deck/load", handlers.LoadDeckFromFileHandler)
		router.POST("/api/deck/resume", handlers.ResumeGame)
	}

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
