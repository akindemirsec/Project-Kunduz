package notification

import (
	"net/http"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/alarms", GetAlarms)
	r.POST("/scan", ScanClustersForVulnerabilities)
}

func GetAlarms(c *gin.Context) {
	var alarms []models.Alarm
	if err := db.DB.Find(&alarms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alarms)
}

func ScanClustersForVulnerabilities(c *gin.Context) {
	var clusters []models.Cluster
	if err := db.DB.Preload("Assets").Preload("SBOMs").Find(&clusters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var allAlarms []models.Alarm
	for _, cluster := range clusters {
		alarms, err := containsVulnerability(cluster.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allAlarms = append(allAlarms, alarms...)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scan completed", "alarms": allAlarms})
}

func containsVulnerability(clusterID uint) ([]models.Alarm, error) {
	var cluster models.Cluster
	if err := db.DB.Preload("Assets").Preload("SBOMs").Where("id = ?", clusterID).First(&cluster).Error; err != nil {
		return nil, err
	}

	var alarms []models.Alarm
	for _, asset := range cluster.Assets {
		var cves []models.CVE
		db.DB.Where("description LIKE ?", "%"+asset.Name+"%").Find(&cves)
		for _, cve := range cves {
			alarm := models.Alarm{
				Message:   "CVE " + cve.CVEID + " found for asset " + asset.Name + " in cluster " + cluster.Name,
				ClusterID: clusterID,
			}
			alarms = append(alarms, alarm)
			db.DB.Create(&alarm)
		}
	}

	for _, sbom := range cluster.SBOMs {
		var cves []models.CVE
		db.DB.Where("description LIKE ?", "%"+sbom.Name+"%").Find(&cves)
		for _, cve := range cves {
			alarm := models.Alarm{
				Message:   "CVE " + cve.CVEID + " found for SBOM " + sbom.Name + " in cluster " + cluster.Name,
				ClusterID: clusterID,
			}
			alarms = append(alarms, alarm)
			db.DB.Create(&alarm)
		}
	}

	return alarms, nil
}
