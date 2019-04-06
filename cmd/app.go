package main

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK/speedsensor"
	"github.com/kevinchapron/FSHK/speedsensor/BLE"
	"time"
)

var targets = []string{"ce:b4:75:29:91:b2", "c5:c8:af:cd:44:6b", "f1:ba:d3:a7:98:3b", "e1:15:11:ee:8f:02"}

var BLEDeviceDiscovered = func(items []BLE.BLEItem) {
	Log.Info(fmt.Sprintf("Discovered %d items.", len(items)))
}

var EventCallback = func(event speedsensor.FullEvent) {
	// BLE  - Know the nearest wristband
	var data_ble = make(map[string]*BLE.BLEItem)
	for _, item := range *event.BLEData {
		if data_ble[item.Addr] == nil {
			data_ble[item.Addr] = &BLE.BLEItem{
				Addr:        item.Addr,
				Name:        item.Name,
				Connectable: item.Connectable,
			}
		}
		data_ble[item.Addr].AllRSSI = append(data_ble[item.Addr].AllRSSI, item.RSSI)
	}
	min_distance := 9999.9
	min_distance_addr := ""
	for addr, item := range data_ble {
		current_distance := item.DistanceOfRSSIs(item.AllRSSI)
		if min_distance > current_distance {
			min_distance = current_distance
			min_distance_addr = addr
		}
	}
	event.ChosenAddr = min_distance_addr

	// Sensor - Know the distance traveled
	distance_traveled := 0
	for i, sensor := range event.Sensors {
		if i == 0 {
			continue
		}
		distance_traveled += speedsensor.IntAbs(int(sensor)-int(event.Sensors[i-1])) * speedsensor.IR_SENSOR_DISTANCE_BETWEEN
	}

	elapsed_time := event.EndTime.Sub(event.StartTime).Seconds()
	speed_ms := (float64(distance_traveled) / 100.0) / elapsed_time
	speed_kmh := ((float64(distance_traveled) / 100000) / elapsed_time) * 3600

	event.DistanceTraveled = distance_traveled
	event.SpeedKMH = speed_kmh
	event.SpeedMS = speed_ms

	Log.Debug("--------------------------")
	Log.Debug("Sensors:", event.Sensors)
	Log.Debug("Distance:", distance_traveled)
	Log.Debug("-- Starting Time:", event.StartTime)
	Log.Debug("-- Ending Time:", event.EndTime)
	Log.Debug(fmt.Sprintf("--> Speed: %.3f (km/h) // %.3f (m/s)", speed_kmh, speed_ms))
	Log.Debug("Wristband:", event.ChosenAddr)

	err := speedsensor.SaveEventToFile(&event)
	if err != nil {
		Log.Error("Error while writing Event in File:", err)
	} else {
		Log.Info("Event stored in file !")
	}
}

func main() {
	Log.SetLevel(Log.DEBUG)
	Log.Info(fmt.Sprintf("Launching app with the following configuration: "))
	Log.Info(speedsensor.GetConfiguration().String())
	err := speedsensor.CreateEventFile()
	if err != nil {
		Log.Error(err)
		return
	}

	Log.Info(fmt.Sprintf("Creating bot ..."))
	bot := speedsensor.CreateBot()
	Log.Info(fmt.Sprintf("Setting bot parameters ..."))
	Log.Info(fmt.Sprintf("Sensors are set at [ %d : %d ]", speedsensor.IR_SENSOR_MIN_DISTANCE, speedsensor.IR_SENSOR_MAX_DISTANCE))

	// Initialization of the BLE interface
	Log.Info("Starting BLE interface ...")
	dev := BLE.GetBLEDevice()
	dev.Scanner.AllowDuplicates = true
	dev.Scanner.TimeScanned = 10 * time.Second
	dev.Scanner.Targets = targets
	dev.Scanner.ScanAllForConfig(&BLEDeviceDiscovered)
	<-time.After(dev.Scanner.TimeScanned)

	Log.Info("Looking for RSSI's ...")
	go dev.Scanner.ScanRSSIContinuously(time.Second * 60)

	// Initialization of the main bot
	bot.SetFrequency(speedsensor.FREQUENCY)
	bot.SetCallback(&speedsensor.ManageData)
	speedsensor.GetAnalyzer().AnalyzerOfEvents.SetCallbackForEvents(&EventCallback)

	go speedsensor.SetTimeout(time.Second*speedsensor.IR_SENSOR_CALIBRATION_TIME, func() {
		speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_RUNNING)
	})
	speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_CALIBRATING)

	Log.Info(bot.Driver.Name(), "Starting bot ...")
	Log.Info(bot.Driver.Name(), fmt.Sprintf("Waiting %d seconds for calibration ...", speedsensor.IR_SENSOR_CALIBRATION_TIME))
	go func() {
		defer bot.Stop()
		err := bot.Start()
		if err != nil {
			Log.Error(err)
		}
	}()

	select {}

}
