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
	router.Static("/static", "static")

	router.Group("/")
	{
		router.GET("/home", handlers.GetHome)
		router.POST("/api/deck/create", handlers.CreateDeck)
		router.GET("/api/ws/:id", handlers.WebSocketHandler)
	}

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
