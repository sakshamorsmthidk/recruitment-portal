package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Recruitment Portal",
		})
	})

	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
