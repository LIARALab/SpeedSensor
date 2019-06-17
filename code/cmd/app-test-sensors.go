package main

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	speedsensor "github.com/kevinchapron/FSHK/speedsensor/code"
	"time"
)

var PrintData = func(b *speedsensor.ADSxDATA) {
	Log.Debug(fmt.Sprintf("||\t%s\t||\t%s\t||\t%s\t||", b.Values[0].String(), b.Values[1].String(), b.Values[2].String()))
}

func main() {
	Log.SetLevel(Log.DEBUG)
	Log.Info(fmt.Sprintf("Launching app."))

	Log.Info(fmt.Sprintf("Creating bot ..."))
	bot := speedsensor.CreateBot()
	Log.Info(fmt.Sprintf("Setting bot parameters ..."))
	Log.Info(fmt.Sprintf("Sensors are set at [ %d : %d ]", speedsensor.IR_SENSOR_MIN_DISTANCE, speedsensor.IR_SENSOR_MAX_DISTANCE))

	bot.SetFrequency(speedsensor.FREQUENCY)
	bot.SetCallback(&PrintData)
	defer bot.Stop()

	go speedsensor.SetTimeout(time.Second*speedsensor.IR_SENSOR_CALIBRATION_TIME, func() {
		speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_RUNNING)
	})
	speedsensor.GetAnalyzer().SetType(speedsensor.ANALYZER_TYPE_CALIBRATING)

	Log.Info(bot.Driver.Name(), "Starting bot ...")
	Log.Info(bot.Driver.Name(), fmt.Sprintf("Waiting %d seconds for calibration ...", speedsensor.IR_SENSOR_CALIBRATION_TIME))

	err := bot.Start()
	if err != nil {
		Log.Error(err)
	}
}
