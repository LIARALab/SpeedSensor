package code

import "time"

const IR_SENSOR_ID_FIRST = 0
const IR_SENSOR_ID_LAST = 2

const IR_SENSOR_MIN_DISTANCE = 20
const IR_SENSOR_MAX_DISTANCE = 150
const IR_SENSOR_DISTANCE_BETWEEN = 60 // Centimeter
const IR_SENSOR_CALIBRATION_TIME = 20 // Seconds

const FREQUENCY = 50
const IR_SENSOR_TOLERATED_POURCENTAGE = 0.05 // pourcentage

const DETECTION_MIN_VALUES = 2
const DETECTION_MAX_SECONDS_BETWEEN = time.Second * 3

const BLE_CONNECT_TO_PERIPHERAL = true
const BLE_DEVICE_DEFAULT = "ce:b4:75:29:91:b2"
const BLE_DEVICE_TEST_1 = "c5:c8:af:cd:44:6b"
const BLE_DEVICE_TEST_2 = "f1:ba:d3:a7:98:3b"
const BLE_DEVICE_TEST_3 = "e1:15:11:ee:8f:02"

func SetTimeout(t time.Duration, f interface{}) {
	time.Sleep(t)
	((f).(func()))()
}

func IntAbs(d int) int {
	if d < 0 {
		return -d
	}
	return d
}

type BLEFeatures struct {
	Addr   string
	Mean   float64
	Median float64
	Stdev  float64
}
