//go:build arduino_nano33

package main

import (
	"machine"

	"tinygo.org/x/drivers/lsm6ds3"
)

const (
	readPin machine.Pin = machine.D3
	ledPin  machine.Pin = machine.LED
)

var (
	accel = lsm6ds3.New(machine.I2C0)
)

func configureTemperatureSensor() error {
	machine.I2C0.Configure(machine.I2CConfig{})
	err := accel.Configure(lsm6ds3.Configuration{})
	return err
}

// readTemperature returns the temperature in celsius milli degrees (°C/1000)
func readTemperature() (t int32, err error) {
	x, err := accel.ReadTemperature()
	// println("Degrees C", float32(x)/1000, "\n\n")
	return x, err
}
