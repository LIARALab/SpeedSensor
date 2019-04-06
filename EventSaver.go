package speedsensor

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/FSHK/speedsensor/BLE"
	"os"
	"strings"
)

const FOLDER_NAME = "Experimentations/"

var Config *Configuration

type Configuration struct {
	UserName string
}

func (configuration *Configuration) String() string {
	b, _ := json.Marshal(configuration)
	return string(b)
}

func GetConfiguration() *Configuration {
	if Config == nil {
		file, err := os.Open("data.json")
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&Config)
		if err != nil {
			panic(err)
		}
	}
	return Config
}

func CreateEventFile() error {
	cfg := GetConfiguration()
	headers := []string{"StartTime", "EndTime", "Sensors", "DistanceTraveled", "SpeedKMH", "SpeedMS", "ChosenAddr", "BleData_DEFAULT", "BleData_TEST_1", "BleData_TEST_2", "BleData_TEST_3"}
	file, err := os.Create(FOLDER_NAME + "DATA_" + cfg.UserName + ".csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(headers)
	return nil
}

func SaveEventToFile(event *FullEvent) error {
	cfg := GetConfiguration()
	sensing_through := []string{}
	for _, i := range event.Sensors {
		sensing_through = append(sensing_through, fmt.Sprint(i))
	}

	var map_ble_data = make(map[string][]BLE.BLEItem)
	var map_ble_str = make(map[string]string)
	for _, d := range *event.BLEData {
		if map_ble_data[d.Addr] == nil {
			map_ble_data[d.Addr] = []BLE.BLEItem{}
		}
		map_ble_data[d.Addr] = append(map_ble_data[d.Addr], d)
	}

	for addr, tmp_data := range map_ble_data {
		bytesTmp, err := json.Marshal(tmp_data)
		if err != nil {
			map_ble_str[addr] = "[]"
		} else {
			map_ble_str[addr] = strings.Replace(string(bytesTmp), "\"\"", "\"", -1)
		}
	}

	data := []string{
		event.StartTime.Format("2006-01-02 15:04:05.999999"),
		event.EndTime.Format("2006-01-02 15:04:05.999999"),
		strings.Join(sensing_through, ""),

		fmt.Sprint(event.DistanceTraveled),
		fmt.Sprintf("%.3f", event.SpeedKMH),
		fmt.Sprintf("%.3f", event.SpeedMS),

		event.ChosenAddr,
		map_ble_str[BLE_DEVICE_DEFAULT],
		map_ble_str[BLE_DEVICE_TEST_1],
		map_ble_str[BLE_DEVICE_TEST_2],
		map_ble_str[BLE_DEVICE_TEST_3],
	}

	f, err := os.OpenFile(FOLDER_NAME+"DATA_"+cfg.UserName+".csv", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Write(data)

	return nil
}
