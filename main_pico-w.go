//go:build pico && cyw43439

package main

import "machine"

const (
	readPin machine.Pin = machine.GP3
	ledPin  machine.Pin = machine.LED
)

func configureTemperatureSensor() error {
	return errTemperatureNotSupported
}

// readTemperature returns the temperature in celsius milli degrees (°C/1000)
func readTemperature() (t int32, err error) {
	return 0, errTemperatureNotSupported
}
