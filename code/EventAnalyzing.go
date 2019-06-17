package code

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK/speedsensor/code/BLE"
	"github.com/montanaflynn/stats"
	"sync"
	"time"
)

var mutex_Analyzer sync.Mutex

type FullEvent struct {
	StartTime time.Time
	EndTime   time.Time
	Sensors   []uint

	EventData *[]*AnalyzerEvent
	BLEData   *[]BLE.BLEItem

	ChosenAddr       string
	DistanceTraveled int
	SpeedKMH         float64
	SpeedMS          float64
}

type AnalyzerEvent struct {
	data      *ADSxValue
	Sensor    uint
	Timestamp time.Time
}
type EventsAnalyzer struct {
	new_event     chan *AnalyzerEvent
	events        []*AnalyzerEvent
	cross_through map[uint]int

	callbackEvent *func(event FullEvent)
}

func (a *EventsAnalyzer) addEvent(event *AnalyzerEvent) {
	mutex_Analyzer.Lock()

	//Logging.Debug(fmt.Sprintf("EVENT #%d (%.2f), at %v. -- ADDED",event.sensor,event.data.Value,event.timestamp.Format("2006-01-02 15:04:05.123456")))
	if len(a.events) > 0 && a.events[len(a.events)-1].Sensor != event.Sensor {
		last := a.events[len(a.events)-1]
		if last.Sensor+1 != event.Sensor && last.Sensor-1 != event.Sensor {
			a.events = a.events[:len(a.events)-1]
		}
	}
	//Logging.Debug(a.events)
	a.events = append(a.events, event)
	a.cross_through[event.Sensor]++

	mutex_Analyzer.Unlock()
}

func (a *EventsAnalyzer) reset() {
	mutex_Analyzer.Lock()
	a.events = []*AnalyzerEvent{}

	for i := 0; i < 3; i++ {
		a.cross_through[uint(i)] = 0
	}
	mutex_Analyzer.Unlock()
	//Logging.Debug("----")
}

func (a *EventsAnalyzer) start() {
	var after_time *time.Timer
	for {
		after_time = time.NewTimer(DETECTION_MAX_SECONDS_BETWEEN)
		select {
		case res := <-a.new_event:
			a.addEvent(res)
			after_time.Reset(DETECTION_MAX_SECONDS_BETWEEN)

			break
		case <-after_time.C:
			mutex_Analyzer.Lock()
			exited := false
			for i := 0; i < 3; i++ {
				if a.cross_through[uint(i)] < DETECTION_MIN_VALUES {
					mutex_Analyzer.Unlock()
					exited = true
					break
				}
				if i == 2 {
					a.computeSpeed()
				}
			}
			if !exited {
				mutex_Analyzer.Unlock()
			}

			a.reset()
			break
		}
	}
}

func (a *EventsAnalyzer) computeSpeed() {
	var data = make([]*AnalyzerEvent, len(a.events))
	copy(data, a.events)

	// Verification: Suppression des doublons (ex: [1 1 1 1] => [1]
	Logging.Debug("STEP-1")
	var sensors []uint
	for _, v := range data {
		if len(sensors) == 0 {
			sensors = append(sensors, v.Sensor)
		} else {
			if sensors[len(sensors)-1] != v.Sensor {
				sensors = append(sensors, v.Sensor)
			}
		}
	}

	// Modification: Trier les fausses détections s'il y en a (en fonction de la médiane)
	Logging.Debug("STEP-2")
	values := []float64{}
	to_delete := []int{}
	for i := 0; i < len(data); i++ {
		values = append(values, data[i].data.GetValue())
	}
	median, _ := stats.Median(values)
	stdev, _ := stats.StandardDeviation(values)
	for i := 0; i < len(data); i++ {
		if data[i].data.GetValue() < median-stdev || data[i].data.GetValue() > median+stdev {
			to_delete = append(to_delete, i)
		}
	}
	for left, right := 0, len(to_delete)-1; left < right; left, right = left+1, right-1 {
		to_delete[left], to_delete[right] = to_delete[right], to_delete[left]
	}
	for _, i := range to_delete {
		data = append(data[:i], data[i+1:]...)
	}

	// Verification: Si il n'y a qu'un sensor => Delete.
	Logging.Debug("STEP-3")
	if len(sensors) == 1 {
		return
	}

	// Modification: S'assurer que le dernier sensor capté correspond à une extrémité
	Logging.Debug("STEP-4")
	for sensors[len(sensors)-1] != IR_SENSOR_ID_FIRST && sensors[len(sensors)-1] != IR_SENSOR_ID_LAST {
		sensors = sensors[:len(sensors)-1]
		to_delete := data[len(data)-1].Sensor
		for i := len(data) - 1; i >= 0; i-- {
			if data[i].Sensor == to_delete {
				continue
			}
			data = data[:i]
			break
		}
	}

	// Modification: S'assurer que le premier sensor capté correspond à une extrémité
	Logging.Debug("STEP-5")
	for sensors[0] != IR_SENSOR_ID_FIRST && sensors[0] != IR_SENSOR_ID_LAST {
		sensors = sensors[1:]
		to_delete := data[len(data)-1].Sensor
		for i := 0; i < len(data); i++ {
			if data[i].Sensor == to_delete {
				continue
			}
			data = data[i:]
			break
		}
	}

	// Modification: Si plus de 3 sensors, ignorer
	Logging.Debug("STEP-6")
	if len(sensors) > 3 {
		Logging.Warning("A detection have been made, but with too much false detections to analyze.")
		return
	}

	// Modification: S'assurer que le premier sensor est différent du dernier
	Logging.Debug("STEP-7")
	if sensors[0] == sensors[len(sensors)-1] {
		return
	}

	Logging.Debug("STEP-OK")
	starting_time := data[0].Timestamp
	ending_time := data[len(data)-1].Timestamp

	data_ble := BLE.GetBLEDevice().Scanner.GetDataBetweenTimes(starting_time, ending_time)

	(*a.callbackEvent)(FullEvent{
		BLEData:   &data_ble,
		EventData: &data,
		Sensors:   sensors,
	})
}
func (a *EventsAnalyzer) SetCallbackForEvents(f *func(event FullEvent)) {
	a.callbackEvent = f
}
