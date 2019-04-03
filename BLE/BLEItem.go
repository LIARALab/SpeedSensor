package BLE

import (
	"math"
	"time"
)

type BLEItem struct {
	Addr        string
	Name        string
	Connectable bool
	RSSI        float64

	Time time.Time
}

func (item *BLEItem) Distance() float64 {
	return 0.2038*math.Pow(float64(item.RSSI), 2) + 12.487*float64(item.RSSI) + 205.18
}
