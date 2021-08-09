package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GameStatus struct {
	Title    string    `json:"title"`
	DateTime time.Time `json:"datetime"`
}

func main() {
	router := gin.Default()
	router.POST("/game-status", saveStatus)
	router.Run()
}

func saveStatus(context *gin.Context) {
	var record GameStatus

	if err := context.ShouldBindJSON(&record); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.DateTime = time.Now()

	fmt.Println(record.Title)
	fmt.Println(record.DateTime)

	context.JSON(200, gin.H{
		"message": "Yep it works!",
	})
}
