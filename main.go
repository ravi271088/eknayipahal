package main

import (
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/home.html",
		))

		tmpl.ExecuteTemplate(c.Writer, "home", gin.H{
			"title": "Home - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/about", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/about.html",
		))

		tmpl.ExecuteTemplate(c.Writer, "about", gin.H{
			"title": "About Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/contact", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/contact.html",
		))
		tmpl.ExecuteTemplate(c.Writer, "contact", gin.H{
			"title": "Contact Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/gallery", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/gallery.html",
		))
		tmpl.ExecuteTemplate(c.Writer, "gallery", gin.H{
			"title": "Gallery - Ek Nayi Pahal NGO",
		})
	})

	r.Run(":8080")
}
