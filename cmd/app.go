package main

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK/speedsensor"
	"time"
)

func main() {
	Log.SetLevel(Log.DEBUG)
	Log.Info(fmt.Sprintf("Launching app."))

	Log.Info(fmt.Sprintf("Creating bot ..."))
	bot := speedsensor.CreateBot()
	Log.Info(fmt.Sprintf("Setting bot parameters ..."))
	Log.Info(fmt.Sprintf("Sensors are set at [ %d : %d ]", speedsensor.IR_SENSOR_MIN_DISTANCE, speedsensor.IR_SENSOR_MAX_DISTANCE))

	bot.SetFrequency(speedsensor.FREQUENCY)
	bot.SetCallback(&speedsensor.ManageData)
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
