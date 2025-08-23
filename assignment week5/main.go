package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Food struct
type Food struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Price int    `json:"price"`
}

var foods = []Food{
	{ID: "01", Name: "Meat steak", Type: "Main course", Price: 79},
	{ID: "02", Name: "Pepsi", Type: "Drinks", Price: 20},
	{ID: "03", Name: "Tomato soup", Type: "Potage soup", Price: 39},
	{ID: "04", Name: "Vanila cake", Type: "Dessert", Price: 99},
	{ID: "05", Name: "Pie", Type: "Dessert", Price: 79},
	{ID: "06", Name: "Vegetable salad", Type: "Salad", Price: 59},
	{ID: "07", Name: "Water", Type: "Drinks", Price: 10},
	{ID: "08", Name: "Fried rice", Type: "Main course", Price: 49},
	{ID: "09", Name: "Pork steak", Type: "Main course", Price: 79},
}

func getFoods(c *gin.Context) {
	priceQuery := c.Query("price")
	typeQuery := c.Query("type")
	idQuery := c.Query("id")
	nameQuery := c.Query("name")

	filter := []Food{}
	for _, food := range foods {
		match := true
		if idQuery != "" && food.ID != idQuery {
			match = false
		}
		if nameQuery != "" && food.Name != nameQuery {
			match = false
		}
		if priceQuery != "" && fmt.Sprint(food.Price) != priceQuery {
			match = false
		}
		if typeQuery != "" && food.Type != typeQuery {
			match = false
		}
		if match {
			filter = append(filter, food)
		}
	}

	c.JSON(http.StatusOK, filter)
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Restaurant!")
		c.String(http.StatusOK, " by 660710692 Kuntapong Maneekhum")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Ready to serve!"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/foods", getFoods)
	}
	r.Run(":3000")

}

//   http://localhost:3000/api/v1/foods?price=20
//   http://localhost:3000/api/v1/foods

//   lsof -i :3000
