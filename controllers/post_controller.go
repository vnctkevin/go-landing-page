package controllers

import (
	"blog/config"
	"blog/models"
	"encoding/json"
	"html/template"
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

func ShowPost(c *gin.Context) {
	id := c.Param("id")
	var post models.BlogPost
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.HTML(http.StatusOK, "blog-detail.html", gin.H{
		"post": post,
	})
}

// GET /blog/create - Show Create Post Form
func ShowCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create-blog.html", gin.H{
		"title": "Create New Post",
	})
}

// POST /blog/create - Create a Post
func CreatePost(c *gin.Context) {
	title := c.PostForm("title")
	contentJSON := c.PostForm("content_json")
	author := "Kevin" // Hardcoded for now

	var postContent models.PostContent

	if contentJSON != "" {
		// Try parsing structured JSON content
		if err := json.Unmarshal([]byte(contentJSON), &postContent); err != nil {
			// Fallback or error? Let's fallback to simple text if parsing fails
			postContent = models.PostContent{{Type: models.ContentTypeParagraph, Text: c.PostForm("content")}}
		}
	} else {
		// Simple text fallback
		postContent = models.PostContent{{Type: models.ContentTypeParagraph, Text: c.PostForm("content")}}
	}

	post := models.BlogPost{Title: title, Content: postContent, Author: author}

	if result := config.DB.Create(&post); result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating post"})
		return
	}

	c.Redirect(http.StatusFound, "/blog")
}

// GET /blog/:id/edit - Show Edit Post Form
func ShowEdit(c *gin.Context) {
	id := c.Param("id")
	var post models.BlogPost
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	// We need to marshal content to JSON string for the frontend editor to load
	contentBytes, _ := json.Marshal(post.Content)

	c.HTML(http.StatusOK, "edit-blog.html", gin.H{
		"title":       "Edit Post",
		"post":        post,
		"contentJson": template.JS(contentBytes),
	})
}

// POST /blog/:id/edit - Update a Post
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.BlogPost
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	post.Title = c.PostForm("title")
	post.Slug = c.PostForm("slug")
	contentJSON := c.PostForm("content_json")

	if contentJSON != "" {
		var postContent models.PostContent
		if err := json.Unmarshal([]byte(contentJSON), &postContent); err == nil {
			post.Content = postContent
		}
	} else {
		// Fallback for simple edits if needed, though we should enforce JSON for full edits
		// post.Content = models.PostContent{{Type: models.ContentTypeParagraph, Text: c.PostForm("content")}}
	}

	config.DB.Save(&post)
	c.Redirect(http.StatusFound, "/blog")
}

// POST /blog/:id/delete - Delete a Post
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.BlogPost{}, id)
	c.Redirect(http.StatusFound, "/blog")
}
