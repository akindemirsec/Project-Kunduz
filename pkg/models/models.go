package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Cluster struct {
	gorm.Model
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
	SBOMs  []SBOM  `json:"sboms"`
}

type Asset struct {
	gorm.Model
	Name      string `json:"name"`
	ClusterID uint   `json:"cluster_id"`
}

type SBOM struct {
	gorm.Model
	Name      string `json:"name"`
	ClusterID uint   `json:"cluster_id"`
}

type CVE struct {
	gorm.Model
	CVEID       string `json:"cve_id"`
	Description string `json:"description"`
}

type Alarm struct {
	gorm.Model
	Message   string `json:"message"`
	ClusterID uint   `json:"cluster_id"`
}
