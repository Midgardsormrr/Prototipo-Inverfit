package payments

import "errors"

var (
	// ErrInvalidID indicates that the provided hexadecimal ID is not a valid MongoDB ObjectID.
	ErrInvalidID = errors.New("identificador de pago inválido")

	// ErrPaymentNotFound is returned when no payment matches the specified ID.
	ErrPaymentNotFound = errors.New("registro de pago no encontrado")

	// ErrStudentNotFound is returned when the associated student does not exist or is deleted.
	ErrStudentNotFound = errors.New("alumno no encontrado o no existe")

	// ErrInvalidAmount is returned when an amount is negative.
	ErrInvalidAmount = errors.New("el monto del pago debe ser mayor o igual a cero")

	// ErrInvalidStatus is returned when an unsupported payment status is specified.
	ErrInvalidStatus = errors.New("estado de pago inválido")

	// ErrInvalidPaymentMethod is returned when an unsupported payment method is specified.
	ErrInvalidPaymentMethod = errors.New("medio de pago inválido")

	// ErrInvalidDateFormat is returned when a date does not comply with the YYYY-MM-DD format.
	ErrInvalidDateFormat = errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
)
