package error

import "errors"

var (
	ErrOrderNotFound              = errors.New("order not found")
	ErrFieldScheduleAlreadyBooked = errors.New("field schedule already booked")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrFieldScheduleAlreadyBooked,
}
