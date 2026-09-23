package main

import (
	"log"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

func main() {
	// Initialize templates at startup
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}

	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "home", gin.H{
			"title": "Home - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/about", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "about", gin.H{
			"title": "About Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/contact", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "contact", gin.H{
			"title": "Contact Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/gallery", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "gallery", gin.H{
			"title": "Gallery - Ek Nayi Pahal NGO",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
