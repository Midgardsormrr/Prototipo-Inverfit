package schedule

import "errors"

var (
	// ErrInvalidDayOfWeek indicates an unsupported day of week identifier.
	ErrInvalidDayOfWeek = errors.New("día de la semana no válido")

	// ErrInvalidTimeFormat indicates a time string not in HH:mm format.
	ErrInvalidTimeFormat = errors.New("formato de hora no válido (debe ser HH:mm)")

	// ErrEndTimeBeforeStartTime indicates end time is before or equal to start time.
	ErrEndTimeBeforeStartTime = errors.New("la hora de término debe ser posterior a la hora de inicio")

	// ErrSessionDurationTooShort indicates session duration is less than minimum allowed (15 min).
	ErrSessionDurationTooShort = errors.New("la duración mínima de una sesión es de 15 minutos")

	// ErrOutsideOperatingHours indicates time is outside gym operating hours (06:30 - 22:00).
	ErrOutsideOperatingHours = errors.New("el horario se encuentra fuera del rango operativo del gimnasio (06:30 a 22:00)")

	// ErrOverlappingIntervals indicates two availability intervals overlap on the same day.
	ErrOverlappingIntervals = errors.New("los intervalos de disponibilidad no pueden superponerse entre sí")

	// ErrSessionOutsideAvailability indicates session is not contained within coach's availability.
	ErrSessionOutsideAvailability = errors.New("la sesión se encuentra fuera de la disponibilidad configurada para el entrenador")

	// ErrSessionOverlap indicates session overlaps with another scheduled session.
	ErrSessionOverlap = errors.New("conflicto de horario: la sesión se superpone con otra sesión existente")

	// ErrAvailabilityConflict indicates availability update would leave existing sessions stranded.
	ErrAvailabilityConflict = errors.New("la nueva disponibilidad deja sesiones existentes fuera de horario")

	// ErrSessionNotFound indicates session was not found or is archived.
	ErrSessionNotFound = errors.New("la sesión solicitada no existe o se encuentra inactiva")

	// ErrExceptionNotFound indicates exception was not found.
	ErrExceptionNotFound = errors.New("la excepción de calendario solicitada no existe")

	// ErrDuplicateException indicates an exception already exists for this session and date.
	ErrDuplicateException = errors.New("ya existe una excepción registrada para esta fecha y sesión")

	// ErrInvalidDateFormat indicates date string not in YYYY-MM-DD format.
	ErrInvalidDateFormat = errors.New("formato de fecha no válido (debe ser YYYY-MM-DD)")

	// ErrInvalidDateRange indicates 'from' date is after 'to' date.
	ErrInvalidDateRange = errors.New("el rango de fechas no es válido: 'from' debe ser anterior o igual a 'to'")

	// ErrDuplicateStudentInSession indicates a student is listed multiple times in the same session.
	ErrDuplicateStudentInSession = errors.New("un alumno no puede estar duplicado en la misma sesión")

	// ErrStudentNotFoundOrArchived indicates one or more students do not exist or are archived.
	ErrStudentNotFoundOrArchived = errors.New("uno o más alumnos asignados no existen o se encuentran archivados")

	// ErrInvalidExceptionType indicates exception type is not 'cancelled' or 'rescheduled'.
	ErrInvalidExceptionType = errors.New("tipo de excepción no válido (debe ser 'cancelled' o 'rescheduled')")

	// ErrRescheduledMissingData indicates rescheduled exception is missing newDate, startTime or endTime.
	ErrRescheduledMissingData = errors.New("una sesión reagendada requiere indicar nueva fecha y horario de inicio y término")

	// ErrDatabaseUnavailable indicates database connection failure.
	ErrDatabaseUnavailable = errors.New("el servicio de base de datos no se encuentra disponible")
)
