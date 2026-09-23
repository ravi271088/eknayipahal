package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

type GalleryImage struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var tmpl *template.Template

func loadImages() ([]GalleryImage, error) {
	data, err := os.ReadFile("gallery.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []GalleryImage{}, nil
		}
		return nil, err
	}
	var images []GalleryImage
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}
	return images, nil
}

func saveImages(images []GalleryImage) error {
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("gallery.json", data, 0644)
}

func main() {
	// Parse all templates at startup
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}

	r := gin.Default()

	r.Static("/static", "./static")

	// Common helper to render content inside layout
	renderPage := func(c *gin.Context, templateName string, title string, data gin.H) {
		// We use the predefined "layout" template and pass the a map
		// that tells the layout which "content" block to render.
		// However, since layout.html uses {{ template "content" . }},
		// we must ensure the current request's specific content is defined as "content".

		// The most reliable way with ParseGlob is to create a new template,
		// associate the layout, and then define the specific content.
		t, err := tmpl.Clone()
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}

		// We map the specific page content (e.g., "home_content") to the generic "content" name
		// so that layout.html's {{ template "content" . }} works.
		t.New("content").Parse(fmt.Sprintf(`{{ template "%s" . }}`, templateName))

		t.ExecuteTemplate(c.Writer, "layout", data)
	}

	r.GET("/", func(c *gin.Context) {
		renderPage(c, "home_content", "Home - Ek Nayi Pahal NGO", gin.H{
			"title": "Home - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/about", func(c *gin.Context) {
		renderPage(c, "about_content", "About Us - Ek Nayi Pahal NGO", gin.H{
			"title": "About Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/contact", func(c *gin.Context) {
		renderPage(c, "contact_content", "Contact Us - Ek Nayi Pahal NGO", gin.H{
			"title": "Contact Us - Ek Nayi Pahal NGO",
		})
	})

	r.GET("/gallery", func(c *gin.Context) {
		images, err := loadImages()
		if err != nil {
			c.String(500, "Error loading gallery: %v", err)
			return
		}

		pageStr := c.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}

		pageSize := 6
		totalImages := len(images)
		start := (page - 1) * pageSize
		end := start + pageSize

		if start > totalImages {
			start = totalImages
		}
		if end > totalImages {
			end = totalImages
		}

		paginatedImages := images[start:end]
		totalPages := int(math.Ceil(float64(totalImages) / float64(pageSize)))

		renderPage(c, "gallery_content", "Gallery - Ek Nayi Pahal NGO", gin.H{
			"title":       "Gallery - Ek Nayi Pahal NGO",
			"Images":      paginatedImages,
			"CurrentPage":  page,
			"TotalPages":   totalPages,
			"PrevPage":     page - 1,
			"NextPage":     page + 1,
		})
	})

	adminAuth := func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()
		if !hasAuth || user != "admin" || pass != "admin123" {
			c.Header("WWW-Authenticate", `Basic realm="Admin Panel"`)
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}

	adminGroup := r.Group("/admin", adminAuth)
	{
		adminGroup.GET("", func(c *gin.Context) {
			images, err := loadImages()
			if err != nil {
				c.String(500, "Error loading images: %v", err)
				return
			}

			renderPage(c, "content", "Admin Panel - Ek Nayi Pahal NGO", gin.H{
				"title":   "Admin Panel - Ek Nayi Pahal NGO",
				"Images":  images,
			})
		})

		adminGroup.POST("/add", func(c *gin.Context) {
			url := c.PostForm("url")
			title := c.PostForm("title")
			description := c.PostForm("description")

			images, _ := loadImages()
			id := 1
			if len(images) > 0 {
				id = images[len(images)-1].ID + 1
			}

			newImage := GalleryImage{
				ID:          id,
				URL:         url,
				Title:       title,
				Description: description,
			}

			images = append(images, newImage)
			if err := saveImages(images); err != nil {
				c.String(500, "Error saving image: %v", err)
				return
			}

			c.Redirect(302, "/admin")
		})

		adminGroup.POST("/admin/delete", func(c *gin.Context) {
			idStr := c.PostForm("id")
			id, _ := strconv.Atoi(idStr)

			images, _ := loadImages()
			var updatedImages []GalleryImage
			for _, img := range images {
				if img.ID != id {
					updatedImages = append(updatedImages, img)
				}
			}

			if err := saveImages(updatedImages); err != nil {
				c.String(500, "Error deleting image: %v", err)
				return
			}

			c.Redirect(302, "/admin")
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
