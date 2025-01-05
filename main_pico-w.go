//go:build pico && cyw43439

package main

import "machine"

const (
	readPin machine.Pin = machine.GP3
	ledPin  machine.Pin = machine.LED
)
