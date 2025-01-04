package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	interval  = 100 * time.Millisecond
	localName = "Sump Sensor"
)

var (
	adapter = bluetooth.DefaultAdapter

	btHomeAdvInterval = bluetooth.NewDuration(interval)
	btHomeServiceUUID = bluetooth.New16BitUUID(0xFCD2)
	btHomeServiceData = []byte{0x40, 0x02, 0x00, 0x00, 0x20, 0x01}
)

func main() {

	// time.Sleep(time.Second)
	// println("starting")

	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	readPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	configureBLE()

	for lastValue := false; ; {
		value := readPin.Get()
		if value != lastValue {
			println("pin value:", value)
		}
		ledPin.Set(value)
		lastValue = value
		if value {
			btHomeServiceData[5] = 1
		} else {
			btHomeServiceData[5] = 0
		}
		advertise(interval)
		time.Sleep(interval)
		bluetooth.DefaultAdapter.DefaultAdvertisement().Stop()
	}

}

func configureBLE() {
	println("starting")
	must("enable BLE stack", adapter.Enable())
}

func advertise(interval time.Duration) {
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Sump Sensor",
		Interval:  bluetooth.NewDuration(interval),
		ServiceData: []bluetooth.ServiceDataElement{
			{UUID: btHomeServiceUUID, Data: btHomeServiceData},
		},
	}))
	must("start adv", adv.Start())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
