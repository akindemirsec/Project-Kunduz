package cve

import (
	"net/http"
	"os/exec"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/cves", getCVEs)
	r.POST("/cves", createCVE)
	r.POST("/cves/update", updateCVE)
}

func getCVEs(c *gin.Context) {
	var cves []models.CVE
	db.DB.Find(&cves)
	c.JSON(http.StatusOK, cves)
}

func createCVE(c *gin.Context) {
	var cve models.CVE
	if err := c.ShouldBindJSON(&cve); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&cve).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cve)
}

func updateCVE(c *gin.Context) {
	cmd := exec.Command("python3", "update_cve_database.py")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CVE database updated successfully"})
}
