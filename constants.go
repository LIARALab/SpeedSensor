package speedsensor

import "time"

const IR_SENSOR_MIN_DISTANCE = 20
const IR_SENSOR_MAX_DISTANCE = 100
const IR_SENSOR_DISTANCE_BETWEEN = 60 // Centimeter
const IR_SENSOR_CALIBRATION_TIME = 1  // Seconds

const FREQUENCY = 50
const IR_SENSOR_TOLERATED_POURCENTAGE = 0.05 // pourcentage

const DETECTION_MIN_VALUES = 2
const DETECTION_MAX_SECONDS_BETWEEN = time.Second * 3

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
