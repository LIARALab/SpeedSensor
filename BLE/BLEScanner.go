package BLE

import (
	"context"
	"github.com/go-ble/ble"
	"sync"
	"time"
)

type BLEScanner struct {
	TimeScanned     time.Duration
	AllowDuplicates bool
	Targets         []string

	RealTimeWindowSize  time.Duration
	RealTimeOverlapping time.Duration

	dev    *ble.Device
	parent *BLEDevice

	items            []BLEItem
	confChannel      chan BLEItem
	mainChannel      chan BLEItem
	callBackFunction *func(items []BLEItem)

	ContinuousData  []BLEItem
	ContinuousMutex sync.Mutex
}

func (bleScanner *BLEScanner) ScanAllForConfig(handler *func([]BLEItem)) {
	go bleScanner._runConfChannel()
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), bleScanner.TimeScanned))
	ble.Scan(ctx, bleScanner.AllowDuplicates, func(a ble.Advertisement) {
		bleScanner.confChannel <- BLEItem{
			Addr:        a.Addr().String(),
			Name:        a.LocalName(),
			Connectable: a.Connectable(),
		}
	}, nil)
	<-time.After(bleScanner.TimeScanned)
	(*handler)(bleScanner.items)
}
func (bleScanner *BLEScanner) _runConfChannel() {
	var fullBreak bool
	t := time.NewTimer(time.Second * 5)
	for {
		fullBreak = false
		select {
		case item := <-bleScanner.confChannel:
			t.Reset(time.Second * 5)
			found := false
			for index, sitem := range bleScanner.items {
				if sitem.Addr != item.Addr {
					continue
				}
				found = true

				if len(item.Addr) > len(sitem.Name) {
					bleScanner.items[index].Name = item.Name
				}

				break
			}
			if !found {
				bleScanner.items = append(bleScanner.items, item)
			}
			break
		case <-t.C:
			fullBreak = true
		}
		if fullBreak {
			break
		}
	}
}

func (bleScanner *BLEScanner) ScanRSSIs(f *func(items []BLEItem)) {
	bleScanner.ScanRSSIsFor(f, -1)
}
func (bleScanner *BLEScanner) ScanRSSIsFor(f *func(items []BLEItem), d time.Duration) {
	var ctx context.Context
	if d > 0 {
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), d))
	} else {
		ctx = context.Background()
	}
	bleScanner.callBackFunction = f

	go bleScanner._runMainChannel()

	ble.Scan(ctx, true, func(a ble.Advertisement) {
		found := false
		for _, s := range bleScanner.Targets {
			if s == a.Addr().String() {
				found = true
			}
		}
		if !found {
			return
		}

		bleScanner.mainChannel <- BLEItem{

			Addr:        a.Addr().String(),
			Name:        a.LocalName(),
			Connectable: a.Connectable(),
			RSSI:        float64(a.RSSI()),

			Time: time.Now(),
		}

	}, nil)
}

func (bleScanner *BLEScanner) ScanRSSIContinuously(storedTime time.Duration) {
	ctx := context.Background()
	ble.Scan(ctx, true, func(a ble.Advertisement) {
		found := false
		for _, s := range bleScanner.Targets {
			if s == a.Addr().String() {
				found = true
				break
			}
		}
		if !found {
			return
		}
		bleScanner.ContinuousMutex.Lock()

		for len(bleScanner.ContinuousData) > 0 && time.Since(bleScanner.ContinuousData[0].Time) > storedTime {
			if len(bleScanner.ContinuousData) == 1 {
				bleScanner.ContinuousData = nil
			} else {
				bleScanner.ContinuousData = bleScanner.ContinuousData[1:]
			}
		}

		bleScanner.ContinuousData = append(bleScanner.ContinuousData, BLEItem{
			Addr:        a.Addr().String(),
			Name:        bleScanner.GetNameFromAddr(a.Addr().String()),
			Connectable: a.Connectable(),
			RSSI:        float64(a.RSSI()),

			Time: time.Now(),
		})
		bleScanner.ContinuousMutex.Unlock()

	}, nil)
}

func (bleScanner *BLEScanner) GetDataBetweenTimes(from time.Time, to time.Time) []BLEItem {
	from_index := -1
	to_index := -1

	bleScanner.ContinuousMutex.Lock()
	//Logging.Debug("Looking for:", from,"and",to)
	for index, item := range bleScanner.ContinuousData {
		//Logging.Debug(index,":",item.Time.After(from),item.Time.After(to))
		if from_index == -1 && item.Time.After(from) {
			from_index = index
		}
		if to_index == -1 && item.Time.After(to) {
			to_index = index
		}
	}
	//Logging.Debug("Final values: ", from_index,"=>",to_index)

	if from_index == -1 || to_index == -1 {
		bleScanner.ContinuousMutex.Unlock()
		return []BLEItem{}
	}

	var items = make([]BLEItem, len(bleScanner.ContinuousData[from_index:to_index]))
	copy(items, bleScanner.ContinuousData[from_index:to_index])
	bleScanner.ContinuousMutex.Unlock()

	return items
}

func (bleScanner *BLEScanner) _runMainChannel() {
	var data []BLEItem
	var all_values = make(map[string]*BLEItem)

	t := time.NewTimer(bleScanner.RealTimeOverlapping)
	for {
		select {
		case item := <-bleScanner.mainChannel:
			data = append(data, item)
			break
		case <-t.C:
			for _, target_addr := range bleScanner.Targets {
				all_values[target_addr] = &BLEItem{
					Addr:    target_addr,
					NbRSSI:  0,
					AllRSSI: nil,
				}
			}

			time_to_cut := time.Now().Add(-bleScanner.RealTimeOverlapping)
			cut_index := -1
			for cutting_index, it := range data {
				all_values[it.Addr].AllRSSI = append(all_values[it.Addr].AllRSSI, float64(it.RSSI))
				all_values[it.Addr].NbRSSI += 1

				if cut_index == -1 && it.Time.Sub(time_to_cut) > 0 {
					cut_index = cutting_index
				}
			}

			var b_array []BLEItem
			for _, value := range all_values {
				b_array = append(b_array, *value)
			}
			(*bleScanner.callBackFunction)(b_array)

			if cut_index > 0 {
				data = data[cut_index:]
			} else {
				data = nil
			}

			t.Reset(bleScanner.RealTimeOverlapping)
			break
		}
	}
}

func (bleScanner *BLEScanner) GetNameFromAddr(addr string) string {
	for _, item := range bleScanner.items {
		if item.Addr == addr {
			return item.Name
		}
	}
	return "NONAME"
}
