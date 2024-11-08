package model

import (
	"gorm.io/gorm"
)

type APIInfo struct {
	gorm.Model
	Name        string
	Path        string
	Downstream  string
	Description string
}

func (md *APIInfo) GetID() uint { return md.ID }
