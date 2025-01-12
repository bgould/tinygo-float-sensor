//go:build arduino_nano33

package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/lsm6ds3"
)

const (
	interval         = 100 * time.Millisecond
	tempReadInterval = 5 * time.Second
	localName        = "Sump Sensor"
)

var (
	tempSensor   = NewTemperatureSensor(machine.I2C0)
	floatSensor1 = NewFloatSensor(machine.D3)
	floatSensor2 = NewFloatSensor(machine.D5)
	bthome       = NewBTHome(interval, localName)
)

func init() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.I2C0.Configure(machine.I2CConfig{})

	tempSensor.Configure()
	floatSensor1.Configure()
	floatSensor2.Configure()
	bthome.Configure()
}

func main() {

	var lastValue1, lastValue2 bool
	for lastTemp, lastTempRead := int32(0), time.Unix(0, 0); ; {

		// float sensor 1
		value1 := floatSensor1.Get()
		if value1 != lastValue1 {
			println("float sensor 1:", value1)
		}
		machine.LED.Set(value1)
		bthome.SetBool(BTHomeDataBinarySensor1, value1)
		lastValue1 = value1

		// float sensor 2
		value2 := floatSensor2.Get()
		if value2 != lastValue2 {
			println("float sensor 2:", value2)
		}
		// ledPin.Set(value2)
		bthome.SetBool(BTHomeDataBinarySensor2, value2)
		lastValue2 = value2

		// read temperature sensor value
		if time.Since(lastTempRead) > tempReadInterval {
			lastTempRead = time.Now()
			if milliDegreesC, err := tempSensor.ReadTemperature(); err != nil {
				if lastTemp != 0 {
					println("error reading temperature", err.Error())
				}
				lastTemp = 0
				bthome.SetSignedInt16(BTHomeDataTemperature, 0)
			} else {
				if lastTemp != milliDegreesC {
					println("Degrees C", float32(milliDegreesC)/1000, "\n\n")
				}
				lastTemp = milliDegreesC
				bthome.SetSignedInt16(BTHomeDataTemperature, int16(milliDegreesC/100))
			}
		}

		// advertise values
		bthome.Advertise(interval)
		time.Sleep(interval)
		bthome.StopAdvertisement()

	}

}

type TemperatureSensor struct {
	accel *lsm6ds3.Device
}

func NewTemperatureSensor(i2c drivers.I2C) *TemperatureSensor {
	return &TemperatureSensor{accel: lsm6ds3.New(i2c)}
}

func (t *TemperatureSensor) Configure() error {
	err := t.accel.Configure(lsm6ds3.Configuration{})
	return err
}

// ReadTemperature returns the temperature in celsius milli degrees (°C/1000)
func (t *TemperatureSensor) ReadTemperature() (temp int32, err error) {
	x, err := t.accel.ReadTemperature()
	return x, err
}

type FloatSensor struct {
	gpio machine.Pin
	last bool
}

func NewFloatSensor(pin machine.Pin) *FloatSensor {
	return &FloatSensor{gpio: pin}
}

func (s *FloatSensor) Configure() {
	s.gpio.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

func (s *FloatSensor) Get() bool {
	return s.gpio.Get()
}

type BTHome struct {
	adapter  *bluetooth.Adapter
	interval bluetooth.Duration
	svcuuid  bluetooth.UUID
	svcdata  []byte
	localnm  string
}

type BTHomeData int

const (
	BTHomeDataBinarySensor1 BTHomeData = 2
	BTHomeDataBinarySensor2 BTHomeData = 4
	BTHomeDataTemperature   BTHomeData = 6
)

func NewBTHome(interval time.Duration, localName string) *BTHome {
	return &BTHome{
		adapter:  bluetooth.DefaultAdapter,
		interval: bluetooth.Duration(interval),
		svcuuid:  bluetooth.New16BitUUID(0xFCD2),
		svcdata:  []byte{0x40, 0x20, 0x01, 0x20, 0x01, 0x45, 0x00, 0x00},
		localnm:  localName,
	}
}

func (bt *BTHome) SetBool(index BTHomeData, val bool) {
	if val {
		bt.svcdata[index] = 1
	} else {
		bt.svcdata[index] = 1
	}
}

func (bt *BTHome) SetSignedInt16(index BTHomeData, val int16) {
	bt.svcdata[index+0] = byte(val >> 0)
	bt.svcdata[index+1] = byte(val >> 8)
}

func (bt *BTHome) Configure() error {
	return bt.adapter.Enable()
}

func (bt *BTHome) Advertise(interval time.Duration) error {
	adv := bt.adapter.DefaultAdvertisement()
	opts := bluetooth.AdvertisementOptions{
		AdvertisementType: bluetooth.AdvertisingTypeScanInd,
		Interval:          bt.interval,
		LocalName:         bt.localnm,
		ServiceData:       []bluetooth.ServiceDataElement{{UUID: bt.svcuuid, Data: bt.svcdata}},
	}
	if err := adv.Configure(opts); err != nil {
		return err
	}
	return adv.Start()
}

func (bt *BTHome) StopAdvertisement() error {
	return bluetooth.DefaultAdapter.DefaultAdvertisement().Stop()
}
