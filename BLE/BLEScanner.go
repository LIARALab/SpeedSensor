package BLE

import (
	"context"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/kevinchapron/BasicLogger/Logging"
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

	items       []BLEItem
	confChannel chan BLEItem
	mainChannel chan BLEItem
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

func (bleScanner *BLEScanner) ScanRSSIs(f *func(a ble.Advertisement)) {
	bleScanner.ScanRSSIsFor(f, -1)
}
func (bleScanner *BLEScanner) ScanRSSIsFor(f *func(a ble.Advertisement), d time.Duration) {
	var ctx context.Context
	if d > 0 {
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), d))
	} else {
		ctx = context.Background()
	}

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

func (bleScanner *BLEScanner) _runMainChannel() {
	var data []BLEItem
	var avg_values = make(map[string]float64)
	var nb_values = make(map[string]float64)

	t := time.NewTimer(bleScanner.RealTimeOverlapping)
	for {
		select {
		case item := <-bleScanner.mainChannel:
			data = append(data, item)
			break
		case <-t.C:
			for _, target_addr := range bleScanner.Targets {
				avg_values[target_addr] = 0
				nb_values[target_addr] = 0
			}

			time_to_cut := time.Now().Add(-bleScanner.RealTimeOverlapping)
			cut_index := -1
			for cutting_index, it := range data {
				nb_values[it.Addr]++
				avg_values[it.Addr] += float64(it.RSSI)

				if cut_index == -1 && it.Time.Sub(time_to_cut) > 0 {
					cut_index = cutting_index
				}
			}

			for addr, value := range avg_values {
				b := BLEItem{
					Addr: addr,
					RSSI: value / nb_values[addr],
				}
				Logging.Debug(fmt.Sprintf("Addr: %s\tDistance: %.2f cm.", addr, b.Distance()))
			}

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
