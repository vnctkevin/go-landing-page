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
	r.POST("/blog", controllers.CreatePost)

	// 6. Run Server (defaults to localhost:8080)
	r.Run()
}
