package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/game-status", saveStatus)
	router.Run()
}

func saveStatus(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "Yep it works!",
	})
}
