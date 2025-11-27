package main

import (
	"blog/config"
	"blog/controllers"
	"blog/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect to Database
	config.ConnectDatabase()

	// 2. Auto Migrate (Creates the table 'blog_posts' automatically)
	config.DB.AutoMigrate(&models.BlogPost{})

	// 3. Initialize Router
	r := gin.Default()

	// 4. Load HTML Templates
	r.LoadHTMLGlob("templates/*")

	// 5. Define Routes
	r.GET("/", controllers.Home)
	r.GET("/blog", controllers.IndexPosts)
	r.GET("/blog/:id", controllers.ShowPost)

	// Blog CRUD Endpoints
	r.GET("/blog/create", controllers.ShowCreate)
	r.POST("/blog/create", controllers.CreatePost)
	r.GET("/blog/:id/edit", controllers.ShowEdit)
	r.POST("/blog/:id/edit", controllers.UpdatePost)
	r.POST("/blog/:id/delete", controllers.DeletePost)

	// Upload Endpoint
	r.POST("/upload", controllers.UploadImage)

	// 6. Run Server (defaults to localhost:8080)
	r.Run()
}
