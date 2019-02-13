package speedsensor

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"gonum.org/v1/gonum/stat"
	"strings"
	"time"
)

var _analyzer *Analyzer

type Analyzer struct {
	type_analysis 		AnalyzerType

	calibrating_array []*ADSxDATA

	calibrating_values_avg 	map[uint]float64
	calibrating_values_std 	map[uint]float64
}

type AnalyzerType int
func (t *AnalyzerType) String() string{
	switch( (int)(*t) ){
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
	Log.Info(fmt.Sprintf("Analyzer type set to \"%v\"",t))
	if a.type_analysis == ANALYZER_TYPE_CALIBRATING && t == ANALYZER_TYPE_RUNNING{
		a.finishCalibration()
	}
	a.type_analysis = t
}
func (a *Analyzer) AddData(b *ADSxDATA) {
	switch(a.type_analysis){
		case ANALYZER_TYPE_RUNNING:
			//Log.Debug(b)
			a.dataReceived(b)
		break
		case ANALYZER_TYPE_CALIBRATING:
			a.calibrating_array = append(a.calibrating_array, b)
		break
	}
}
func (a *Analyzer) finishCalibration(){
	var data = make(map[uint][]float64)
	var avg = make(map[uint]float64)
	var std = make(map[uint]float64)

	for _, calibration_value := range a.calibrating_array{
		for index, value := range calibration_value.Values{
			if value.Error != nil {
				continue
			}
			if data[index] == nil{
				data[index] = []float64{}
			}
			data[index] = append(data[index],value.GetValue())
		}
	}

	var i uint
	for i = 0;i<3;i++{
		if data[i] == nil {
			data[i] = []float64{0}
		}
	}

	for index, datas := range data{
		avg[index] = stat.Mean(datas,nil)
		std[index] = stat.StdDev(datas,nil)

		if avg[index] < IR_SENSOR_MIN_DISTANCE || avg[index] > IR_SENSOR_MAX_DISTANCE{
			avg[index] = 0
			std[index] = 0
		}
	}

	a.calibrating_values_avg = avg
	a.calibrating_values_std = std

	var s = []string{}
	for index, _ := range avg{
		s = append(s,fmt.Sprintf("%d => (%.2f ± %.2f)",index,avg[index],std[index]))
	}
	Log.Info(fmt.Sprintf("Calibration finished with data : %v",strings.Join(s," || ")))
}

func (a *Analyzer) dataReceived(b *ADSxDATA){

	for index, value := range b.Values{
		if value.Error != nil {
			continue
		}
		var min_value = (a.calibrating_values_avg[index] - a.calibrating_values_std[index]) * (1-IR_SENSOR_TOLERATED_POURCENTAGE)
		var max_value = (a.calibrating_values_avg[index] + a.calibrating_values_std[index]) * (1+IR_SENSOR_TOLERATED_POURCENTAGE)

		if !(min_value <= value.Value && max_value >= value.Value){
			Log.Info("Someone passing in front of",index,"at",time.Now().Format("2006-01-02 15:04:05")," =>",min_value,"<",value.Value,"<",max_value)
		}

		//Log.Debug(min_value, "<",value.Value,"<",max_value)
	}
}

func GetAnalyzer() *Analyzer {
	if _analyzer == nil {
		_analyzer = &Analyzer{}
	}
	return _analyzer
}

var ManageData = func(b *ADSxDATA) {
	GetAnalyzer().AddData(b)
}
