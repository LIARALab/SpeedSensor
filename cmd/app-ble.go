package main

//
//import (
//	"fmt"
//	"github.com/go-ble/ble"
//	L "github.com/kevinchapron/BasicLogger/Logging"
//	"github.com/kevinchapron/FSHK/speedsensor/BLE"
//	"sort"
//	"time"
//)
//
////var targets = []string{"ce:b4:75:29:91:b2","c5:c8:af:cd:44:6b", "f1:ba:d3:a7:98:3b", "e1:15:11:ee:8f:02"}
//var connect_target = "ce:b4:75:29:91:b2"
//
//var BLEPeriphConnected = func(periph ble.Client) {
//	L.Info(fmt.Sprintf("Successfully connected to %s.", periph.Addr().String()))
//	defer func() {
//		err := periph.CancelConnection()
//		if err != nil {
//			L.Error(err)
//		} else {
//			L.Info(fmt.Sprintf("Successfully disconnected from %s.", periph.Addr().String()))
//		}
//	}()
//
//	profile, err := periph.DiscoverProfile(true)
//	if err != nil {
//		L.Error(err)
//	}
//
//	for _, s := range profile.Services {
//		L.Info(fmt.Sprintf("    Service %s %s, Handle (0x%02X)", s.UUID, ble.Name(s.UUID), s.Handle))
//
//		for _, c := range s.Characteristics {
//			L.Info(fmt.Sprintf("        Characteristic: %s %s, Property: 0x%02X, Handle(0x%02X), VHandle(0x%02X)",
//				c.UUID, ble.Name(c.UUID), c.Property, c.Handle, c.ValueHandle))
//
//			switch c.UUID.String() {
//			case "c2211dc241e211e9b210d663bd873d93":
//				// Battery Health
//				b, err := periph.ReadCharacteristic(c)
//				if err != nil {
//					L.Error(err)
//					break
//				}
//				L.Info(fmt.Sprintf("        -> Battery Health: %x | %d %%", b, b))
//				break
//			case "c2211f4841e211e9b210d663bd873d93":
//				// Battery value
//				b, err := periph.ReadCharacteristic(c)
//				if err != nil {
//					L.Error(err)
//					break
//				}
//				L.Info(fmt.Sprintf("        -> Battery Value:  %x | %d %%", b, b))
//				break
//			case "fd37495849d111e98646d663bd873d93":
//				// Activities
//				L.Info("        Subscribing to notifications ...")
//				h := func(req []byte) { L.Info(fmt.Sprintf("        -> Notified: %q [ % X ]", string(req), req)) }
//				if err := periph.Subscribe(c, false, h); err != nil {
//					L.Error(err)
//				}
//				time.Sleep(time.Second * 5)
//				if err := periph.Unsubscribe(c, false); err != nil {
//					L.Error(err)
//				}
//				L.Info("        Unsubscribed.")
//				break
//			default:
//				L.Warning("No UUID recognized for characteristics.")
//			}
//
//			for _, d := range c.Descriptors {
//				L.Info(fmt.Sprintf("            Descriptor: %s %s, Handle(0x%02x)",
//					d.UUID, ble.Name(d.UUID), d.Handle))
//			}
//		}
//	}
//}
//
//var BLEDeviceDiscovered = func(items []BLE.BLEItem) {
//	L.Info(fmt.Sprintf("Discovered %d items.", len(items)))
//	//
//	//if !speedsensor.BLE_CONNECT_TO_PERIPHERAL {
//	//	return
//	//}
//	//for _, item := range items {
//	//	if item.Addr != connect_target {
//	//		continue
//	//	}
//	//	periph, err := BLE.GetBLEDevice().ConnectTo(item)
//	//	if err != nil {
//	//		L.Error(err)
//	//		break
//	//	}
//	//	BLEPeriphConnected(periph)
//	//}
//}
//
//var BLEListenerRSSI = func(items []BLE.BLEItem) {
//	sort.Slice(items, func(i, j int) bool {
//		return BLE.GetBLEDevice().Scanner.GetNameFromAddr(items[i].Addr) < BLE.GetBLEDevice().Scanner.GetNameFromAddr(items[j].Addr)
//	})
//
//	for _, item := range items {
//		sort.Sort(sort.Reverse(sort.Float64Slice(item.AllRSSI)))
//		var mediane_rssis []float64
//		if item.NbRSSI >= 4{
//			mediane_rssis = item.AllRSSI[int(item.NbRSSI/4):int(item.NbRSSI/4)*3]
//		}else{
//			mediane_rssis = item.AllRSSI
//		}
//
//		L.Info(fmt.Sprintf("[%s] %s: (%.2f) %.2f cm || %.2f cm [%d]",BLE.GetBLEDevice().Scanner.GetNameFromAddr(item.Addr), item.Addr,item.RSSI,item.DistanceOfRSSIs(item.AllRSSI),item.DistanceOfRSSIs(mediane_rssis), len(item.AllRSSI)))
//		//L.Info(fmt.Sprintf(" -A-> %.2f cm (%d) %v",item.DistanceOfRSSIs(item.AllRSSI),int(item.NbRSSI),item.AllRSSI))
//		//L.Info(fmt.Sprintf(" -Q-> %.2f cm (%d) %v",item.DistanceOfRSSIs(first_quartile),len(first_quartile),first_quartile))
//	}
//	L.Info("--------------------")
//}
//
//func main() {
//	L.SetLevel(L.DEBUG)
//
//	L.Info("Getting BLE device ...")
//	dev := BLE.GetBLEDevice()
//
//	dev.Scanner.AllowDuplicates = true
//	dev.Scanner.TimeScanned = 10 * time.Second
//	dev.Scanner.Targets = targets
//	L.Info("Scanning ...")
//	dev.Scanner.ScanAllForConfig(&BLEDeviceDiscovered)
//
//	L.Info("Looking for RSSI's ...")
//	dev.Scanner.RealTimeWindowSize = 5 * time.Second
//	dev.Scanner.RealTimeOverlapping = 2500 * time.Millisecond // 50 % of Window Size
//	dev.Scanner.ScanRSSIs(&BLEListenerRSSI)
//}
