// gpio/mem is not available for non-Linux, so always return test mode ON..

// +build linux

package helpers

import (
	"fmt"

	"github.com/warthog618/gpio"
)

var (
	pinIn,
	pinOut *gpio.Pin
)

func SetUp(gpioIn, gpioOut int) error {
	if (gpioIn + gpioOut) == 0 {
		// all fake, no need to open gpio
		pinIn = nil
		pinOut = nil
		return nil
	}

	err := gpio.Open()
	if err != nil {
		return err
	}

	if gpioIn > 0 {
		pinIn = gpio.NewPin(gpioIn)
		if pinIn == nil {
			gpio.Close()
			return fmt.Errorf("invalid GPIO %d", gpioIn)
		}
	}

	if gpioOut > 0 {
		pinOut = gpio.NewPin(gpioOut)
		if pinOut == nil {
			gpio.Close()
			return fmt.Errorf("invalid GPIO %d", gpioOut)
		}
		pinOut.Output()
	}

	return nil
}

func ReadFeedback() bool {
	if pinIn == nil {
		return ReadFakeFeedback()
	}
	return bool(pinIn.Read())
}

func SetTrigger(value bool) {
	if pinOut == nil {
		SetFakeTrigger(value)
	} else {
		pinOut.Write(gpio.Level(value))
	}

}

func Cleanup() {
	if pinIn == nil && pinOut == nil {
		return
	}

	if pinOut != nil {
		pinOut.Input()
	}
	gpio.Close()
}
