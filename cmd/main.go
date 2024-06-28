package main

import (
	"net/http"
	"project-kunduz/pkg/auth"
	"project-kunduz/pkg/cluster"
	"project-kunduz/pkg/cve"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"
	"project-kunduz/pkg/notification"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	r := gin.Default()

	// Static files
	r.Static("/static", "./web/static")

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*")

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", auth.Login)

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/register", auth.Register)

	db.ConnectDatabase()
	createDefaultAdmin()

	// Protected routes
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/clusters", func(c *gin.Context) {
			c.HTML(http.StatusOK, "clusters.html", nil)
		})
		protected.GET("/cves", func(c *gin.Context) {
			c.HTML(http.StatusOK, "cves.html", nil)
		})
		protected.GET("/alarms", func(c *gin.Context) {
			c.HTML(http.StatusOK, "alarms.html", nil)
		})
		protected.GET("/scan", func(c *gin.Context) {
			c.HTML(http.StatusOK, "scan.html", nil)
		})
		protected.POST("/logout", auth.Logout)

		// Register additional protected routes for API endpoints
		cluster.RegisterRoutes(protected.Group("/api"))
		cve.RegisterRoutes(protected.Group("/api"))
		notification.RegisterRoutes(protected.Group("/api"))
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

func createDefaultAdmin() {
	var admin models.User
	db.DB.Where("username = ?", "admin").First(&admin)
	if admin.ID == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin = models.User{Username: "admin", Password: string(hashedPassword)}
		db.DB.Create(&admin)
	}
}
