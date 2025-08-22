package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book struct
type Book struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Price int    `json:"price"`
	Page  int    `json:"page"`
}

var books = []Book{
	{ID: "01", Name: "Tip & Trick for freshman", Type: "Self-help", Price: 159, Page: 127},
	{ID: "02", Name: "Detective Shane", Type: "Crime fiction", Price: 399, Page: 256},
	{ID: "03", Name: "Dr.Who", Type: "Novel", Price: 429, Page: 420},
	{ID: "04", Name: "Stock Market in big 25", Type: "Marketing", Price: 99, Page: 121},
	{ID: "05", Name: "Jack Of All Trait:Season 1", Type: "Action fiction", Price: 199, Page: 233},
	{ID: "06", Name: "BiologyTech by Mike 5th", Type: "Science", Price: 560, Page: 369},
}

func getBooks(c *gin.Context) {
	priceQuery := c.Query("price")
	pageQuery := c.Query("page")
	typeQuery := c.Query("type")
	idQuery := c.Query("id")
	nameQuery := c.Query("name")

	filter := []Book{}
	for _, book := range books {
		match := true
		if idQuery != "" && book.ID != idQuery {
			match = false
		}
		if nameQuery != "" && book.Name != nameQuery {
			match = false
		}
		if priceQuery != "" && fmt.Sprint(book.Price) != priceQuery {
			match = false
		}
		if pageQuery != "" && fmt.Sprint(book.Page) != pageQuery {
			match = false
		}
		if typeQuery != "" && book.Type != typeQuery {
			match = false
		}
		if match {
			filter = append(filter, book)
		}
	}

	c.JSON(http.StatusOK, filter)
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Book Store!")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy!"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/books", getBooks)
	}
	r.Run(":8080")

}

//   http://localhost:8080/api/v1/books?price=159&type=Self-help
//   http://localhost:8080/api/v1/books

//   lsof -i :8080
