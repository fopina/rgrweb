package helpers

var (
	fakeStatus bool
)

func ReadFakeFeedback() bool {
	return fakeStatus
}

func SetFakeTrigger(value bool) {
	fakeStatus = value
}
