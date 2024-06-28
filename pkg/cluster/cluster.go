package cluster

import (
	"net/http"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/clusters", CreateCluster)
	r.GET("/clusters", GetClusters)
	r.GET("/clusters/:id", GetCluster)
	r.PUT("/clusters/:id", UpdateCluster)
	r.DELETE("/clusters/:id", DeleteCluster)
	r.POST("/clusters/:id/assets", AddAssetToCluster)
	r.POST("/clusters/:id/sboms", AddSBOMToCluster)
}

func CreateCluster(c *gin.Context) {
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

func GetClusters(c *gin.Context) {
	var clusters []models.Cluster
	if err := db.DB.Preload("Assets").Preload("SBOMs").Find(&clusters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clusters)
}

func GetCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.Preload("Assets").Preload("SBOMs").Where("id = ?", c.Param("id")).First(&cluster).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

func UpdateCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.Where("id = ?", c.Param("id")).First(&cluster).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

func DeleteCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := db.DB.Where("id = ?", c.Param("id")).First(&cluster).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	if err := db.DB.Delete(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cluster deleted successfully"})
}

func AddAssetToCluster(c *gin.Context) {
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

func AddSBOMToCluster(c *gin.Context) {
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
