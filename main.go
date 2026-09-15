package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server on :8080")
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// Home route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home", gin.H{"title": "Home - Ek Nayi Pahal NGO"})
	})

	// About route
	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about", gin.H{"title": "About Us - Ek Nayi Pahal NGO"})
	})

	// Contact route
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact", gin.H{"title": "Contact Us - Ek Nayi Pahal NGO"})
	})

	r.Run(":8080")
}
