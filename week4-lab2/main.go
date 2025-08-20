package main

import (
	"fmt"
)

func main() {
	// var name string = "Kuntapong"
	var age int = 21
	email := "maneekhum_k@su.ca.th"
	gpa := 3.50

	firstname, lastname := "Kuntapong", "Maneekhum"

	fmt.Printf("Name %s %s, age %d, email %s, gpa %.2f", firstname, lastname, age, email, gpa)
}
