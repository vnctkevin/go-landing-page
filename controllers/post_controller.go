package controllers

import (
	"blog/models"
	"blog/services"
	"encoding/json"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var postService = services.NewPostService()

// GET / - Render Landing Page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "ShipFlow",
	})
}

// GET /blog - Render Blog List
func IndexPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10

	posts, total, err := postService.GetAllPosts(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching posts"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.HTML(http.StatusOK, "blog.html", gin.H{
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
		"hasPrev":     page > 1,
		"hasNext":     page < totalPages,
		"prevPage":    page - 1,
		"nextPage":    page + 1,
	})
}

func ShowPost(c *gin.Context) {
	id := c.Param("id")
	post, err := postService.GetPost(id)
	if err != nil {
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
	metaTitle := c.PostForm("meta_title")
	metaDescription := c.PostForm("meta_description")
	keywords := c.PostForm("keywords")
	canonicalURL := c.PostForm("canonical_url")

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

	post := models.BlogPost{
		Title:           title,
		Content:         postContent,
		Author:          author,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		Keywords:        keywords,
		CanonicalURL:    canonicalURL,
	}

	if err := postService.CreatePost(&post); err != nil {
		c.JSON(500, gin.H{"error": "Error creating post"})
		return
	}

	c.Redirect(http.StatusFound, "/blog")
}

// GET /blog/:id/edit - Show Edit Post Form
func ShowEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := postService.GetPost(id)
	if err != nil {
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
	post, err := postService.GetPost(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	post.Title = c.PostForm("title")
	post.Slug = c.PostForm("slug")
	post.MetaTitle = c.PostForm("meta_title")
	post.MetaDescription = c.PostForm("meta_description")
	post.Keywords = c.PostForm("keywords")
	post.CanonicalURL = c.PostForm("canonical_url")

	contentJSON := c.PostForm("content_json")

	if contentJSON != "" {
		var postContent models.PostContent
		if err := json.Unmarshal([]byte(contentJSON), &postContent); err == nil {
			post.Content = postContent
		}
	}

	if err := postService.UpdatePost(post); err != nil {
		c.JSON(500, gin.H{"error": "Error updating post"})
		return
	}

	c.Redirect(http.StatusFound, "/blog")
}

// POST /blog/:id/delete - Delete a Post
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := postService.DeletePost(id); err != nil {
		c.JSON(500, gin.H{"error": "Error deleting post"})
		return
	}
	c.Redirect(http.StatusFound, "/blog")
}
