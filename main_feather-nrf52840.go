//go:build feather_nrf52840

package main

import "machine"

const (
	readPin machine.Pin = machine.A1
	ledPin  machine.Pin = machine.LED
)
