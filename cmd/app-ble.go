package main

import (
	"fmt"
	"github.com/go-ble/ble"
	L "github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK/speedsensor"
	"github.com/kevinchapron/FSHK/speedsensor/BLE"
	"time"
)

var targets = []string{"ce:b4:75:29:91:b2"}
var connect_target = "ce:b4:75:29:91:b2"

var BLEPeriphConnected = func(periph ble.Client) {
	L.Info(fmt.Sprintf("Successfully connected to %s.", periph.Addr().String()))
	defer func() {
		err := periph.CancelConnection()
		if err != nil {
			L.Error(err)
		} else {
			L.Info(fmt.Sprintf("Successfully disconnected from %s.", periph.Addr().String()))
		}
	}()

	profile, err := periph.DiscoverProfile(true)
	if err != nil {
		L.Error(err)
	}

	for _, s := range profile.Services {
		L.Info(fmt.Sprintf("    Service %s %s, Handle (0x%02X)", s.UUID, ble.Name(s.UUID), s.Handle))

		for _, c := range s.Characteristics {
			L.Info(fmt.Sprintf("        Characteristic: %s %s, Property: 0x%02X, Handle(0x%02X), VHandle(0x%02X)",
				c.UUID, ble.Name(c.UUID), c.Property, c.Handle, c.ValueHandle))

			switch c.UUID.String() {
			case "c2211dc241e211e9b210d663bd873d93":
				// Battery Health
				b, err := periph.ReadCharacteristic(c)
				if err != nil {
					L.Error(err)
					break
				}
				L.Info(fmt.Sprintf("        -> Battery Health: %x | %d %%", b, b))
				break
			case "c2211f4841e211e9b210d663bd873d93":
				// Battery value
				b, err := periph.ReadCharacteristic(c)
				if err != nil {
					L.Error(err)
					break
				}
				L.Info(fmt.Sprintf("        -> Battery Value:  %x | %d %%", b, b))
				break
			case "fd37495849d111e98646d663bd873d93":
				// Activities
				L.Info("        Subscribing to notifications ...")
				h := func(req []byte) { L.Info(fmt.Sprintf("        -> Notified: %q [ % X ]", string(req), req)) }
				if err := periph.Subscribe(c, false, h); err != nil {
					L.Error(err)
				}
				time.Sleep(time.Second * 5)
				if err := periph.Unsubscribe(c, false); err != nil {
					L.Error(err)
				}
				L.Info("        Unsubscribed.")
				break
			default:
				L.Warning("No UUID recognized for characteristics.")
			}

			for _, d := range c.Descriptors {
				L.Info(fmt.Sprintf("            Descriptor: %s %s, Handle(0x%02x)",
					d.UUID, ble.Name(d.UUID), d.Handle))
			}
		}
	}
}

var BLEDeviceDiscovered = func(items []BLE.BLEItem) {
	L.Info(fmt.Sprintf("Discovered %d items.", len(items)))
	if !speedsensor.BLE_CONNECT_TO_PERIPHERAL {
		return
	}
	for _, item := range items {
		if item.Addr != connect_target {
			continue
		}
		periph, err := BLE.GetBLEDevice().ConnectTo(item)
		if err != nil {
			L.Error(err)
			break
		}
		BLEPeriphConnected(periph)
	}
}

var BLEListenerRSSI = func(a ble.Advertisement) {
	L.Info(a.RSSI())
}

func main() {
	L.SetLevel(L.DEBUG)

	L.Info("Getting BLE device ...")
	dev := BLE.GetBLEDevice()

	dev.Scanner.AllowDuplicates = true
	dev.Scanner.TimeScanned = 5 * time.Second
	dev.Scanner.Targets = targets
	L.Info("Scanning ...")
	dev.Scanner.ScanAllForConfig(&BLEDeviceDiscovered)

	//L.Info("Looking for RSSI's ...")
	//dev.Scanner.RealTimeWindowSize = 10 * time.Second
	//dev.Scanner.RealTimeOverlapping = 5 * time.Second // 50 % of Window Size
	//dev.Scanner.ScanRSSIs(&BLEListenerRSSI)
}
