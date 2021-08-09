package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "go"
	password = "gopassword"
	dbname   = "gamestatustracker"
)

type GameStatus struct {
	Title    string    `json:"title"`
	DateTime time.Time `json:"datetime"`
	Status   string    `json:"status"`
}

func main() {
	router := gin.Default()
	router.POST("/save-game-status", saveStatus)
	router.GET("/current-game-status/:game-name", getGameStatus)
	router.GET("/game-timeline/:game-name", getGameStatusTimeline)
	router.Run()
}

func saveStatus(context *gin.Context) {
	var record GameStatus

	if err := context.ShouldBindJSON(&record); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.DateTime = time.Now()

	psqlConnectionInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlConnectionInfo)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	insertStatement := `
	INSERT INTO gamestatus (Title, DateTime, Status)
	VALUES ($1, $2, $3)`

	_, err = db.Exec(insertStatement, record.Title, time.Now(), record.Status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{
		"message": "Yep it works!",
	})
}

func getGameStatus(context *gin.Context) {
	title := context.Param("game-name")
	psqlConnectionInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlConnectionInfo)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var status string
	selectStatement := `SELECT status FROM gamestatus WHERE title=$1 ORDER BY datetime DESC LIMIT 1`
	row := db.QueryRow(selectStatement, title)
	switch err := row.Scan(&status); err {
	case sql.ErrNoRows:
		context.String(http.StatusOK, "No status found")
		return
	case nil:
		context.String(http.StatusOK, status)
		return
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func getGameStatusTimeline(context *gin.Context) {
	title := context.Param("game-name")
	psqlConnectionInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlConnectionInfo)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var timeline []GameStatus
	selectStatement := `SELECT datetime, status FROM gamestatus WHERE title=$1 ORDER BY datetime`
	rows, err := db.Query(selectStatement, title)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		record := new(GameStatus)
		err = rows.Scan(&record.DateTime, &record.Status)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		timeline = append(timeline, *record)
	}

	err = rows.Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, timeline)
}
