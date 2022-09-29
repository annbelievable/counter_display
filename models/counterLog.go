package models

import "time"

type CounterLog struct {
	Value    uint8     `json:"value"`
	Datetime time.Time `json:"datetime"`
}
