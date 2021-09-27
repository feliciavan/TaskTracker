package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tsk []Task

// album represents data about a record album.
type Task struct {
	ID       int64  `json:"id"`
	TEXT     string `json:"text"`
	DAY      string `json:"day"`
	REMINDER bool   `json:"reminder"`
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "taskdb",
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// albums slice to seed record album data.

	tsk, err = getAll()
	if err != nil {
		log.Fatal(err)
	}

	//router := gin.Default()
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/tasks", getTasks)
	router.DELETE("/tasks/:id", getTaskByID)
	router.PUT("/tasks/:id", changeTask)
	router.POST("/tasks", postTasks)
	router.Run("localhost:8080")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//retrieve all records
func getAll() ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM task")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var onetsk Task
		if err := rows.Scan(&onetsk.ID, &onetsk.TEXT, &onetsk.DAY, &onetsk.REMINDER); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tasks = append(tasks, onetsk)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return tasks, nil
}

// getAlbums responds with the list of all albums as JSON.
func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tsk)
}

// postAlbums adds an album from JSON received in the request body.
func postTasks(c *gin.Context) {
	var newTask Task
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newTask); err != nil {
		return
	}
	// Add the new album to the slice.
	addID, error := addTask(newTask)
	fmt.Print(addID)
	fmt.Print(error)
	tsk = append(tsk, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func addTask(tsk Task) (int64, error) {
	result, err := db.Exec("INSERT INTO task (text, day, reminder) VALUES (?, ?, ?)", tsk.TEXT, tsk.DAY, tsk.REMINDER)
	if err != nil {
		return 0, fmt.Errorf("addTask: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addTask: %v", err)
	}
	return id, nil
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getTaskByID(c *gin.Context) {
	idstr := c.Param("id")
	idget, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	deleteTask(c, idget)
}

func deleteTask(c *gin.Context, id int64) {
	_, err := db.Exec("DELETE FROM task WHERE id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, id)
}

func changeTask(c *gin.Context) {
	idstr := c.Param("id")
	idget, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	changeReminder(c, idget)
}

func changeReminder(c *gin.Context, id int64) {
	row := db.QueryRow("SELECT reminder FROM task WHERE id = ?", id)
	var reminder bool

	if err := row.Scan(&reminder); err != nil {
		fmt.Println(err)
	}
	reminder = !reminder

	_, err := db.Exec("UPDATE task SET REMINDER = ? WHERE id = ?", reminder, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, id)
}
