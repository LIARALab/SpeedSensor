package main

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK/speedsensor"
	"time"
)

func SetTimeout(t time.Duration, f interface{}) {
	time.Sleep(t)
	((f).(func()))()
}


func main() {
	Log.SetLevel(Log.DEBUG)
	Log.Info(fmt.Sprintf("Launching app."))

	Log.Info(fmt.Sprintf("Creating bot ..."))
	bot := speedsensor.CreateBot()
	Log.Info(fmt.Sprintf("Setting bot parameters ..."))
	bot.SetFrequency(speedsensor.FREQUENCY)
	bot.SetCallback(&speedsensor.ManageData)
	defer bot.Stop()

	go SetTimeout(time.Second*speedsensor.IR_SENSOR_CALIBRATION_TIME, func() {
		speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_RUNNING)
	})
	speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_CALIBRATING)

	Log.Info(bot.Driver.Name(), "Starting bot ...")
	err := bot.Start()
	if err != nil {
		Log.Error(err)
	}
}
