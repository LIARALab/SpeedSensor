package BLE

import (
	"math"
	"time"
)

type BLEItem struct {
	Addr        string `json:"-"`
	Name        string `json:"-"`
	Connectable bool   `json:"-"`

	RSSI float64

	AllRSSI []float64 `json:"-"`
	NbRSSI  float64   `json:"-"`

	Time time.Time
}

func (item *BLEItem) Distance() float64 {
	return item.DistanceOfRSSI(item.RSSI)
}
func (item *BLEItem) DistanceOfRSSIs(float64s []float64) float64 {
	var avg float64
	for _, t := range float64s {
		avg += t
	}
	return item.DistanceOfRSSI(avg / float64(len(float64s)))
}

func (item *BLEItem) DistanceOfRSSI(rssi float64) float64 {
	return 0.2038*math.Pow(rssi, 2) + 12.487*rssi + 205.18
}
