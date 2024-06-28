package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Cluster struct {
	gorm.Model
	Name   string  `json:"name"`
	SBOMs  []SBOM  `json:"sboms"`
	Assets []Asset `json:"assets"`
}

type SBOM struct {
	gorm.Model
	Name      string `json:"name"`
	Data      string `json:"data"`
	ClusterID uint   `json:"cluster_id"`
}

type Asset struct {
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

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
