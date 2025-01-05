//go:build arduino_nano33

package main

import "machine"

const (
	readPin machine.Pin = machine.D3

	ledPin machine.Pin = machine.LED
)
