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
	// Initialize templates
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}

	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "layout", gin.H{
			"title": "Home - Ek Nayi Pahal NGO",
			"content": "home_content",
		})
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "layout", gin.H{
			"title": "About Us - Ek Nayi Pahal NGO",
			"content": "about_content",
		})
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "layout", gin.H{
			"title": "Contact Us - Ek Nayi Pahal NGO",
			"content": "contact_content",
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

		c.HTML(200, "layout", gin.H{
			"title":       "Gallery - Ek Nayi Pahal NGO",
			"content":     "gallery_content",
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

			c.HTML(200, "layout", gin.H{
				"title":   "Admin Panel - Ek Nayi Pahal NGO",
				"content": "content",
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

	// This is critical: Gin's c.HTML uses the HTMLRender interface.
	// We need to tell Gin how to use our standard text/template.
	r.SetHTMLTemplate(htmlRender{tmpl: tmpl})

	r.Run(":" + port)
}

// htmlRender implements gin.HTMLRender
type htmlRender struct {
	tmpl *template.Template
}

func (h htmlRender) Render(w http.ResponseWriter, name string, data interface{}) {
	// Since we are using a layout, we always execute the "layout" template
	// but we must ensure the 'content' block in layout.html is filled.
	// We do this by cloning the template and defining a new 'content' block
	// that points to the specific page content.

	// The data passed to layout must contain the 'content' name.
	// We assume the 'data' is a gin.H (map[string]interface{})
	dataMap, ok := data.(gin.H)
	if !ok {
		// If it's not a map, we can't dynamically set the content
		h.tmpl.ExecuteTemplate(w, name, data)
		return
	}

	contentName, ok := dataMap["content"].(string)
	if !ok {
		h.tmpl.ExecuteTemplate(w, name, data)
		return
	}

	t, err := h.tmpl.Clone()
	if err != nil {
		http.Error(w, "Template Clone Error", 500)
		return
	}

	// Map the layout's "content" block to the actual page content
	t.New("content").Parse(fmt.Sprintf(`{{ template "%s" . }}`, contentName))

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Need to add "net/http" to imports
