package model

import (
	"gorm.io/gorm"
)

type Downstream struct {
	gorm.Model
	Name string `gorm:"unique"`
	URL  string
}

func (md *Downstream) GetID() uint { return md.ID }
