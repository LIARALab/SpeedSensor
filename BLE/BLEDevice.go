package BLE

import (
	"context"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/kevinchapron/BasicLogger/Logging"
	"time"
)

var device *BLEDevice

type BLEDevice struct {
	dev     ble.Device
	Scanner BLEScanner
}

func (bleDevice BLEDevice) ConnectTo(item BLEItem) (ble.Client, error) {
	filter := func(a ble.Advertisement) bool {
		return a.Addr().String() == item.Addr
	}
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), time.Second*10))
	periph, err := ble.Connect(ctx, filter)
	if err != nil {
		return nil, err
	}
	return periph, nil
}

func GetBLEDevice() *BLEDevice {
	if device == nil {
		bledevice, err := dev.NewDevice("default")
		if err != nil {
			Logging.Error(err)
		}

		device = &BLEDevice{
			dev: bledevice,
			Scanner: BLEScanner{
				TimeScanned:     0,
				AllowDuplicates: true,
				dev:             &bledevice,
				confChannel:     make(chan BLEItem),
				mainChannel:     make(chan BLEItem),
			},
		}
		device.Scanner.parent = device
		ble.SetDefaultDevice(device.dev)
	}
	return device
}
