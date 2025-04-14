package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// no db so using in memory
var books = []Book{
	{ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
	{ID: "2", Title: "Introducing Go", Author: "Caleb Doxsey"},
}

func main() {
	// this creates a gin powered router
	router := gin.Default()
	// add authentication for some more fun
	// creating a default account here
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password123", // adhiraj please learn cyber security too
	}))

	// lets create a public access and private access route

	// public
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to public route"})
	})

	// admin only router
	authorized.GET("/private", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + user + ", you have full authority"})
	})
	// simple get to list books
	router.GET("/books", getBooks)
	// simple post to add book
	router.POST("/books", addBook)

	// start the router
	router.Run("localhost:8080")
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var newBook Book // empty struct

	// here i try to bind incoming json to the struct
	if err := c.BindJSON(&newBook); err != nil {
		// fancy if statement
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	// update in memory
	books = append(books, newBook)

	// lets respond with the confirmed book, good practice
	c.IndentedJSON(http.StatusCreated, newBook)
}
