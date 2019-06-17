package BLE

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/kevinchapron/BasicLogger/Logging"
)

var device *BLEDevice

type BLEDevice struct {
	dev     ble.Device
	Scanner BLEScanner
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
