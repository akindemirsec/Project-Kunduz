package cve

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"project-kunduz/pkg/db"
	"project-kunduz/pkg/models"

	"github.com/gin-gonic/gin"
)

type CVEEntry struct {
	CVEID       string `xml:"id,attr"`
	Description string `xml:"desc"`
}

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/cves", GetCVEs)
	r.POST("/update_cves", UpdateCVEs)
}

func GetCVEs(c *gin.Context) {
	var cves []models.CVE
	if err := db.DB.Find(&cves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cves)
}

func UpdateCVEs(c *gin.Context) {
	err := downloadAndUpdateCVEData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CVE database updated successfully"})
}

func downloadAndUpdateCVEData() error {
	url := "https://cve.mitre.org/data/downloads/allitems.xml"

	tempFile, err := os.CreateTemp("", "cve_data_*.xml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download CVE data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download CVE data: status code %d", resp.StatusCode)
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write CVE data to file: %w", err)
	}

	body, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read CVE data: %w", err)
	}

	var cveList CVEList
	err = xml.Unmarshal(body, &cveList)
	if err != nil {
		return fmt.Errorf("failed to parse CVE data: %w", err)
	}

	for _, entry := range cveList.Entries {
		var cve models.CVE
		if err := db.DB.Where("cve_id = ?", entry.CVEID).First(&cve).Error; err != nil {
			cve = models.CVE{
				CVEID:       entry.CVEID,
				Description: entry.Description,
			}
			db.DB.Create(&cve)
		} else {
			cve.Description = entry.Description
			db.DB.Save(&cve)
		}
	}

	return nil
}

func containsVulnerability(clusterID uint) ([]models.Alarm, error) {
	var cluster models.Cluster
	if err := db.DB.Preload("Assets").Preload("SBOMs").Where("id = ?", clusterID).First(&cluster).Error; err != nil {
		return nil, fmt.Errorf("failed to find cluster: %w", err)
	}

	var alarms []models.Alarm
	for _, asset := range cluster.Assets {
		var cves []models.CVE
		db.DB.Where("description LIKE ?", "%"+asset.Name+"%").Find(&cves)
		for _, cve := range cves {
			alarm := models.Alarm{
				Message:   fmt.Sprintf("CVE %s found for asset %s in cluster %s", cve.CVEID, asset.Name, cluster.Name),
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
				Message:   fmt.Sprintf("CVE %s found for SBOM %s in cluster %s", cve.CVEID, sbom.Name, cluster.Name),
				ClusterID: clusterID,
			}
			alarms = append(alarms, alarm)
			db.DB.Create(&alarm)
		}
	}

	return alarms, nil
}
