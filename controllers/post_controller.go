package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET / - Render Landing Page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcome to My SaaS",
	})
}

// GET /blog - Render Blog List
func IndexPosts(c *gin.Context) {
	var posts []models.BlogPost
	// Equivalent to Laravel: BlogPost::all()
	config.DB.Find(&posts)

	c.HTML(http.StatusOK, "blog.html", gin.H{
		"posts": posts,
	})
}

// POST /blog - Create a Post
func CreatePost(c *gin.Context) {
	// Get form data
	title := c.PostForm("title")
	content := c.PostForm("content")

	post := models.BlogPost{Title: title, Content: models.PostContent{{Type: models.ContentTypeParagraph, Text: content}}, Author: "Kevin"}

	// Equivalent to: $post->save()
	result := config.DB.Create(&post)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating post"})
		return
	}

	// Redirect back to blog page
	c.Redirect(http.StatusFound, "/blog")
}
