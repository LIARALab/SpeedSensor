package code

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"github.com/montanaflynn/stats"
)

var _analyzer *Analyzer

type Analyzer struct {
	type_analysis AnalyzerType

	calibrating_array      []*ADSxDATA
	calibrating_values_avg map[uint]float64
	calibrating_values_std map[uint]float64

	AnalyzerOfEvents EventsAnalyzer
}
type AnalyzerType int

func (t *AnalyzerType) String() string {
	switch (int)(*t) {
	case 1:
		return "CALIBRATING"
	case 2:
		return "RUNNING"
	}
	return "---"
}

const ANALYZER_TYPE_CALIBRATING AnalyzerType = 1
const ANALYZER_TYPE_RUNNING AnalyzerType = 2

func (a *Analyzer) SetType(t AnalyzerType) {
	Log.Info(fmt.Sprintf("Analyzer type set to \"%s\"", t.String()))
	if a.type_analysis == ANALYZER_TYPE_CALIBRATING && t == ANALYZER_TYPE_RUNNING {
		a.finishCalibration()
	}
	a.type_analysis = t
}
func (a *Analyzer) AddData(b *ADSxDATA) {
	switch a.type_analysis {
	case ANALYZER_TYPE_RUNNING:
		//Log.Debug(b)
		a.dataReceived(b)
		break
	case ANALYZER_TYPE_CALIBRATING:
		a.calibrating_array = append(a.calibrating_array, b)
		break
	}
}
func (a *Analyzer) finishCalibration() {
	var data = make(map[uint][]float64)
	var avg = make(map[uint]float64)
	var std = make(map[uint]float64)

	for _, calibration_value := range a.calibrating_array {
		for index, value := range calibration_value.Values {
			if value.Error != nil {
				continue
			}
			if data[index] == nil {
				data[index] = []float64{}
			}
			data[index] = append(data[index], value.GetValue())
		}
	}

	var i uint
	for i = 0; i < 3; i++ {
		if data[i] == nil {
			data[i] = []float64{0}
		}
	}

	for index, datas := range data {
		avg[index], _ = stats.Mean(datas)
		std[index], _ = stats.StandardDeviation(datas)

		if avg[index] < IR_SENSOR_MIN_DISTANCE || avg[index] > IR_SENSOR_MAX_DISTANCE {
			avg[index] = 0
			std[index] = 0
		}
	}

	a.calibrating_values_avg = avg
	a.calibrating_values_std = std

	Log.Info("Calibration finished. See values below : ")
	for i = 0; i < 3; i++ {
		Log.Info(fmt.Sprintf("\t\t%d => (%.2f ± %.2f)", i, avg[i], std[i]))
	}
	go _analyzer.AnalyzerOfEvents.start()
}
func (a *Analyzer) dataReceived(b *ADSxDATA) {
	for index, value := range b.Values {
		if value.Error != nil {
			continue
		}
		var min_value = (a.calibrating_values_avg[index] - a.calibrating_values_std[index]) * (1 - IR_SENSOR_TOLERATED_POURCENTAGE)
		var max_value = (a.calibrating_values_avg[index] + a.calibrating_values_std[index]) * (1 + IR_SENSOR_TOLERATED_POURCENTAGE)

		if !(min_value <= value.Value && max_value >= value.Value) {
			a.eventReceived(AnalyzerEvent{
				data:      value,
				Sensor:    index,
				Timestamp: b.Timestamp,
			})
			break
		}
	}
}
func (a *Analyzer) eventReceived(event AnalyzerEvent) {
	a.AnalyzerOfEvents.new_event <- &event
	//Log.Debug("Data received: ",event.Sensor,event.data.GetValue())
}

func GetAnalyzer() *Analyzer {
	if _analyzer == nil {
		_analyzer = &Analyzer{}
		_analyzer.AnalyzerOfEvents.cross_through = make(map[uint]int)
		_analyzer.AnalyzerOfEvents.new_event = make(chan *AnalyzerEvent)
		_analyzer.AnalyzerOfEvents.reset()
	}
	return _analyzer
}

var ManageData = func(b *ADSxDATA) {
	GetAnalyzer().AddData(b)
}
