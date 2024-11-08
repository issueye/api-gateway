package model

import (
	"gorm.io/gorm"
)

type TrafficStats struct {
	gorm.Model
	API        string // 被访问的API名称
	InTraffic  int64  // 入站流量大小（字节数）
	OutTraffic int64  // 出站流量大小（字节数）
}

func (md *TrafficStats) GetID() uint { return md.ID }
