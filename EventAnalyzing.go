package speedsensor

import (
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"sync"
	"time"
)

var mutex_Analyzer sync.Mutex

type AnalyzerEvent struct {
	data      *ADSxValue
	sensor    uint
	timestamp time.Time
}
type EventsAnalyzer struct {
	new_event     chan *AnalyzerEvent
	events        []*AnalyzerEvent
	cross_through map[uint]int
}

func (a *EventsAnalyzer) addEvent(event *AnalyzerEvent) {
	mutex_Analyzer.Lock()

	//Log.Debug(fmt.Sprintf("EVENT #%d (%.2f), at %v. -- ADDED",event.sensor,event.data.Value,event.timestamp.Format("2006-01-02 15:04:05.123456")))
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

func (a *EventsAnalyzer) computeSpeed() {
	var data = make([]*AnalyzerEvent, len(a.events))
	copy(data, a.events)

	Log.Debug("Computing speed : ")
	var sensors []uint
	for _, v := range data {
		Log.Debug(fmt.Sprintf("\tSensor: %d\tDate: %v", v.sensor, v.timestamp.Format("2006-01-02 15:04:05")))

		if len(sensors) == 0 {
			sensors = append(sensors, v.sensor)
		} else {
			if sensors[len(sensors)-1] != v.sensor {
				sensors = append(sensors, v.sensor)
			}
		}
	}
	Log.Debug("Speed computed.")

	if len(sensors) == 1 {
		return
	}

	distance := 0
	for i, sensor := range sensors {
		if i == 0 {
			continue
		}
		distance += IntAbs(int(sensor)-int(sensors[i-1])) * IR_SENSOR_DISTANCE_BETWEEN
	}

	elapsed_time := data[len(data)-1].timestamp.Sub(data[0].timestamp).Seconds()
	speed_ms := (float64(distance) / 100.0) / elapsed_time
	speed_kmh := ((float64(distance) / 100000) / elapsed_time) * 3600

	Log.Debug("Sensors direction : ", sensors)
	Log.Debug("Starting_time :", data[0].timestamp.Format("2006-01-02 15:04:05"))
	Log.Debug("Ending_time :", data[len(data)-1].timestamp.Format("2006-01-02 15:04:05"))
	Log.Debug("-----------")
	Log.Debug("Distance : ", distance)
	Log.Debug(fmt.Sprintf("Time elapsed : %.3f", elapsed_time))
	Log.Debug(fmt.Sprintf("VITESSE : %.3f m/s (%.3f km/h).", speed_ms, speed_kmh))
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
