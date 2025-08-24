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
type Table struct {
	ID     string `json:"id"`
	Seats  int    `json:"seats"`
	Status string `json:"status"`
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
	{ID: "10", Name: "Egg bacon grilled", Type: "Appetizer", Price: 35},
	{ID: "11", Name: "Breadroll", Type: "Dessert", Price: 15},
}
var tables = []Table{
	{ID: "T01", Seats: 2, Status: "available"},
	{ID: "T02", Seats: 4, Status: "available"},
	{ID: "T03", Seats: 4, Status: "available"},
	{ID: "T04", Seats: 6, Status: "available"},
	{ID: "T05", Seats: 8, Status: "reserved"},
}

func reserveTable(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table id is required"})
		return
	}

	for i, table := range tables {
		if table.ID == id {
			if table.Status == "reserved" {
				c.JSON(http.StatusConflict, gin.H{"message": "Table already reserved"})
				return
			}
			tables[i].Status = "reserved"
			c.JSON(http.StatusOK, gin.H{"message": "Table reserved successfully", "table": tables[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
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

	r.GET("/tables", reserveTable)

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Ready to serve!"})
	})

	api := r.Group("/order")
	{
		api.GET("/foods", getFoods)
	}

	r.Run(":3000")

}

//   http://localhost:3000/order/foods?price=20
//   http://localhost:3000/order/foods?type=Main%20course&price=79
//   http://localhost:3000/order/foods

// 	 http://localhost:3000/tables
//   http://localhost:3000/tables?id=T01
//   http://localhost:3000/tables?id=T05

//   lsof -i :3000
