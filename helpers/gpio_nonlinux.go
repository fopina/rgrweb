// gpio/mem is not available for non-Linux, so always return test mode ON..

// +build !linux

package helpers

func SetUp(gpioIn, gpioOut int) error {
	return nil
}

func ReadFeedback() bool {
	return ReadFakeFeedback()
}

func SetTrigger(value bool) {
	SetFakeTrigger(value)
}

func Cleanup() {

}
