package model

import "time"

type HotRank struct {
	Title string
	Link string
	HotValue string
	Day string
	Source string
	CreateTime time.Time
}
