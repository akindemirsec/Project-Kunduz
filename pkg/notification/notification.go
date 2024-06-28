package notification

import (
	"net/http"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/alarms", getAlarms)
}

func getAlarms(c *gin.Context) {
	var alarms []models.Alarm
	if err := db.DB.Find(&alarms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alarms)
}
