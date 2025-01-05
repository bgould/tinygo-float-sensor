package main

import (
	"errors"
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	interval         = 100 * time.Millisecond
	tempReadInterval = 5 * time.Second
	localName        = "Sump Sensor"

	binarySensorIdx = 2
	tempSensorIdx   = 4
)

var (
	adapter = bluetooth.DefaultAdapter

	btHomeAdvInterval = bluetooth.NewDuration(interval)
	btHomeServiceUUID = bluetooth.New16BitUUID(0xFCD2)
	btHomeServiceData = []byte{0x40, 0x20, 0x01, 0x45, 0x00, 0x00}

	errTemperatureNotSupported = errors.New("temperature not supported")
	errTemperatureNotAvailable = errors.New("temperature not available")
)

func main() {

	// time.Sleep(time.Second)
	// println("starting")

	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	readPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	configureBLE()

	temperatureSupported := false
	if err := configureTemperatureSensor(); err != nil {
		println("could not configure temperature sensor:", err.Error())
	} else {
		temperatureSupported = true
	}

	for lastValue, lastTemp, lastTempRead := false, int32(0), time.Unix(0, 0); ; {

		// read float sensor value
		value := readPin.Get()
		if value != lastValue {
			println("pin value:", value)
		}
		ledPin.Set(value)
		lastValue = value
		if value {
			btHomeServiceData[binarySensorIdx] = 1
		} else {
			btHomeServiceData[binarySensorIdx] = 0
		}

		// read temperature sensor value
		if temperatureSupported && time.Since(lastTempRead) > tempReadInterval {
			lastTempRead = time.Now()
			if milliDegreesC, err := readTemperature(); err != nil {
				if lastTemp != 0 {
					println("error reading temperature", err.Error())
				}
				lastTemp = 0
				btHomeServiceData[tempSensorIdx+0] = 0x0
				btHomeServiceData[tempSensorIdx+1] = 0x0
			} else {
				if lastTemp != milliDegreesC {
					println("Degrees C", float32(milliDegreesC)/1000, "\n\n")
				}
				lastTemp = milliDegreesC
				btHomeServiceData[tempSensorIdx+0] = byte((milliDegreesC / 100) >> 0)
				btHomeServiceData[tempSensorIdx+1] = byte((milliDegreesC / 100) >> 8)
			}
		}

		// advertise values
		advertise(interval)
		time.Sleep(interval)
		stopAdvertisement()
	}

}

func configureBLE() {
	must("enable BLE stack", adapter.Enable())
}

func advertise(interval time.Duration) {
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:         localName,
		AdvertisementType: bluetooth.AdvertisingTypeScanInd,
		Interval:          btHomeAdvInterval,
		ServiceData: []bluetooth.ServiceDataElement{
			{UUID: btHomeServiceUUID, Data: btHomeServiceData},
		},
	}))
	must("start adv", adv.Start())
}

func stopAdvertisement() {
	must("stop adv", bluetooth.DefaultAdapter.DefaultAdvertisement().Stop())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	} else {
		// println(time.Now().String(), action+" successful")
	}
}
