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

	// Serve index.html at the root URL
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Auth routes
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", auth.Login)
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/register", auth.Register)

	// Protected routes
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		db.ConnectDatabase()
		createDefaultAdmin() // Default admin kullanıcıyı oluştur

		auth.RegisterRoutes(protected)
		cluster.RegisterRoutes(protected)
		cve.RegisterRoutes(protected)
		notification.RegisterRoutes(protected)
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
