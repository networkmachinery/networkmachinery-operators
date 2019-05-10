package error

func NewZeroEventsError(message string) *ZeroEventsError {
	return &ZeroEventsError{
		message: message,
	}
}

type ZeroEventsError struct {
	message string
}

func (z *ZeroEventsError) Error() string {
	if len(z.message) == 0 {
		return "no network events occured in the past survillence period"
	}
	return z.message
}
