package cluster

import (
	"net/http"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/clusters", getClusters)
	r.POST("/clusters", createCluster)
	r.GET("/clusters/:id", getCluster)
	r.PUT("/clusters/:id", updateCluster)
	r.DELETE("/clusters/:id", deleteCluster)
	r.POST("/clusters/:id/sboms", addSBOM)
	r.POST("/clusters/:id/assets", addAsset)
}

func getClusters(c *gin.Context) {
	var clusters []models.Cluster
	db.DB.Preload("SBOMs").Preload("Assets").Find(&clusters)
	c.JSON(http.StatusOK, clusters)
}

func createCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cluster)
}

func getCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.Preload("SBOMs").Preload("Assets").First(&cluster, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}
	c.JSON(http.StatusOK, cluster)
}

func updateCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.First(&cluster, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Save(&cluster)
	c.JSON(http.StatusOK, cluster)
}

func deleteCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.First(&cluster, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	db.DB.Delete(&cluster)
	c.JSON(http.StatusOK, gin.H{"message": "Cluster deleted successfully"})
}

func addSBOM(c *gin.Context) {
	var sbom models.SBOM
	if err := c.ShouldBindJSON(&sbom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cluster ID"})
		return
	}
	sbom.ClusterID = uint(clusterID)
	if err := db.DB.Create(&sbom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sbom)
}

func addAsset(c *gin.Context) {
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cluster ID"})
		return
	}
	asset.ClusterID = uint(clusterID)
	if err := db.DB.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, asset)
}
