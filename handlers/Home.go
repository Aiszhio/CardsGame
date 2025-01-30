package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"status": "success",
	})
}
