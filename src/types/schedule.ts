export type DayOfWeek = 'Lun' | 'Mar' | 'Mié' | 'Jue' | 'Vie' | 'Sáb' | 'Dom'

// A single availability interval within a day (e.g. 08:30 to 12:00)
export interface AvailabilityInterval {
  id: string
  startTime: string // "HH:MM"
  endTime: string   // "HH:MM"
}

// Coach availability for a specific day of the week
export interface DayAvailability {
  day: DayOfWeek
  enabled: boolean
  intervals: AvailabilityInterval[]
}

// Scheduled session / block
export interface ScheduleBlock {
  id: string
  day: DayOfWeek
  startTime: string // "HH:MM"
  endTime: string   // "HH:MM"
  durationMinutes: number
  type: 'Sesión grupal' | 'Individual' | 'Evaluación'
  studentIds: string[]
  maxCapacity?: number // default 3 en INVERFIT
  isRecurring?: boolean // default true
  notes?: string
}

// Calendar exception for a specific date of a recurring session
export interface ScheduleException {
  id: string
  originalBlockId: string
  date: string // ISO date "YYYY-MM-DD" of the affected session instance
  type: 'rescheduled' | 'cancelled'
  newDate?: string      // "YYYY-MM-DD" if moved
  newDay?: DayOfWeek    // Day of the week for new date
  newStartTime?: string // "HH:MM"
  newEndTime?: string   // "HH:MM"
  notes?: string
  createdAt: string
}

// Session instance rendered in the calendar for a specific day/date
export interface SessionOccurrence {
  instanceId: string // unique for this week instance
  blockId: string
  day: DayOfWeek
  date: string // "YYYY-MM-DD"
  startTime: string
  endTime: string
  durationMinutes: number
  type: 'Sesión grupal' | 'Individual' | 'Evaluación'
  studentIds: string[]
  maxCapacity?: number // default 3
  isRecurring: boolean
  isException: boolean
  exceptionType?: 'rescheduled' | 'cancelled'
  exceptionId?: string
  notes?: string
}

// DTO for assigning / creating a session
export interface AssignSessionDTO {
  blockId?: string
  studentIds: string[]
  day: DayOfWeek
  startTime: string
  endTime: string
  durationMinutes?: number
  type?: 'Sesión grupal' | 'Individual' | 'Evaluación'
  maxCapacity?: number
  isRecurring?: boolean
  startDate?: string
  endDate?: string
  notes?: string
}

// Backward-compatible DTO
export interface AssignStudentDTO {
  blockId?: string
  studentId: string
  day: DayOfWeek
  startTime: string
  durationMinutes?: number
  endTime?: string
  isRecurring?: boolean
  startDate?: string
  endDate?: string
  notes?: string
}

export interface CreateExceptionDTO {
  blockId: string
  date: string
  type: 'rescheduled' | 'cancelled'
  newDate?: string
  newDay?: DayOfWeek
  newStartTime?: string
  newEndTime?: string
  notes?: string
}
