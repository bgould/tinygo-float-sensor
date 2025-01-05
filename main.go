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
	// advertise(interval)

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
		stopAdvertisement()
	}

}

func configureBLE() {
	must("enable BLE stack", adapter.Enable())
}

func advertise(interval time.Duration) {
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: localName,
		Interval:  btHomeAdvInterval,
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
