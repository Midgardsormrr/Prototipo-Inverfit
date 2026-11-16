package students

import "errors"

var (
	// ErrStudentNotFound indica que el alumno consultado no existe o se encuentra archivado.
	ErrStudentNotFound = errors.New("alumno no encontrado")

	// ErrInvalidID indica que el formato de ID proporcionado no es un ObjectID hexadecimal válido de MongoDB.
	ErrInvalidID = errors.New("identificador de alumno inválido")

	// ErrEmptyName indica que el nombre del alumno es requerido y no fue provisto.
	ErrEmptyName = errors.New("el nombre del alumno es obligatorio")

	// ErrInvalidFrequency indica que la frecuencia semanal informada no es una de las permitidas.
	ErrInvalidFrequency = errors.New("frecuencia semanal inválida (valores permitidos: '1 vez/sem.', '2 veces/sem.', '3 veces/sem.', '4 veces/sem.')")

	// ErrInvalidStatus indica que el estado de membresía informado no es uno de los permitidos.
	ErrInvalidStatus = errors.New("estado de membresía inválido (valores permitidos: 'Activo', 'Por revisar', 'Inactivo')")

	// ErrDatabaseUnavailable indica que la base de datos no está disponible para persistir la operación.
	ErrDatabaseUnavailable = errors.New("base de datos no disponible")
)
