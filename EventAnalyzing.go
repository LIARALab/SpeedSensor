package speedsensor

import (
	"github.com/kevinchapron/FSHK/speedsensor/BLE"
	"sync"
	"time"
)

var mutex_Analyzer sync.Mutex

type FullEvent struct {
	StartTime time.Time
	EndTime   time.Time
	Sensors   []uint

	BLEData *[]BLE.BLEItem

	ChosenAddr       string
	DistanceTraveled int
	SpeedKMH         float64
	SpeedMS          float64
}

type AnalyzerEvent struct {
	data      *ADSxValue
	sensor    uint
	timestamp time.Time
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
	if len(a.events) > 0 && a.events[len(a.events)-1].sensor != event.sensor {
		last := a.events[len(a.events)-1]
		if last.sensor+1 != event.sensor && last.sensor-1 != event.sensor {
			a.events = a.events[:len(a.events)-1]
		}
	}
	//Logging.Debug(a.events)
	a.events = append(a.events, event)
	a.cross_through[event.sensor]++

	mutex_Analyzer.Unlock()
}

func (a *EventsAnalyzer) reset() {
	mutex_Analyzer.Lock()
	a.events = []*AnalyzerEvent{}

	for i := 0; i < 3; i++ {
		a.cross_through[uint(i)] = 0
	}
	mutex_Analyzer.Unlock()
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

	var sensors []uint
	for _, v := range data {
		if len(sensors) == 0 {
			sensors = append(sensors, v.sensor)
		} else {
			if sensors[len(sensors)-1] != v.sensor {
				sensors = append(sensors, v.sensor)
			}
		}
	}

	if len(sensors) == 1 {
		return
	}

	starting_time := data[0].timestamp
	ending_time := data[len(data)-1].timestamp

	data_ble := BLE.GetBLEDevice().Scanner.GetDataBetweenTimes(starting_time, ending_time)

	(*a.callbackEvent)(FullEvent{
		BLEData:   &data_ble,
		StartTime: starting_time,
		EndTime:   ending_time,
		Sensors:   sensors,
	})
}
func (a *EventsAnalyzer) SetCallbackForEvents(f *func(event FullEvent)) {
	a.callbackEvent = f
}
