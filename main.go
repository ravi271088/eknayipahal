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
		// Create a temporary template for this page that defines "content" as "home_content"
		// and then executes the "layout" template.
		t, err := tmpl.Clone()
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}
		t.New("content").Parse(`{{ template "home_content" . }}`)
		t.ExecuteTemplate(c.Writer, "layout", gin.H{
			"title": "Home - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/about", func(c *gin.Context) {
		t, err := tmpl.Clone()
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}
		t.New("content").Parse(`{{ template "about_content" . }}`)
		t.ExecuteTemplate(c.Writer, "layout", gin.H{
			"title": "About Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/contact", func(c *gin.Context) {
		t, err := tmpl.Clone()
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}
		t.New("content").Parse(`{{ template "contact_content" . }}`)
		t.ExecuteTemplate(c.Writer, "layout", gin.H{
			"title": "Contact Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/gallery", func(c *gin.Context) {
		t, err := tmpl.Clone()
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}
		t.New("content").Parse(`{{ template "gallery_content" . }}`)
		t.ExecuteTemplate(c.Writer, "layout", gin.H{
			"title": "Gallery - Ek Nayi Pahal NGO",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
