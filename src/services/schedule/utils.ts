import type { DayOfWeek } from '@/types'

// Helper to convert "HH:MM" to minutes from midnight
export function timeToMinutes(timeStr: string): number {
  if (!timeStr) return 0
  const [h, m] = timeStr.split(':').map(Number)
  return (h || 0) * 60 + (m || 0)
}

// Helper to convert minutes to "HH:MM"
export function minutesToTime(minutes: number): string {
  const normalized = Math.max(0, Math.min(minutes, 24 * 60 - 1))
  const h = Math.floor(normalized / 60)
  const m = normalized % 60
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
}

// Helper to add minutes to a time string
export function addMinutesToTime(timeStr: string, minutesToAdd: number): string {
  const current = timeToMinutes(timeStr)
  return minutesToTime(current + minutesToAdd)
}

// Helper to format ISO date "YYYY-MM-DD"
export function formatDateISO(date: Date): string {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

// Canonical DayOfWeek list: Monday through Sunday
export const DAY_NAMES: DayOfWeek[] = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']

export interface WeekDayInfo {
  day: DayOfWeek
  dateStr: string // "YYYY-MM-DD"
  dayNumber: number
  monthLabel: string
  isToday: boolean
  dateObj: Date
}

// Helper to get all 7 days for a given week offset from today (0 = current week)
export function getWeekDates(weekOffset: number = 0, baseDate: Date = new Date()): WeekDayInfo[] {
  const today = new Date(baseDate)
  today.setHours(0, 0, 0, 0)
  const todayISO = formatDateISO(today)

  // Find Monday of the current week (Lun = 1, Dom = 0 -> shift 0 to 7)
  const dayOfWeek = today.getDay()
  const mondayOffset = dayOfWeek === 0 ? -6 : 1 - dayOfWeek

  const monday = new Date(today)
  monday.setDate(today.getDate() + mondayOffset + (weekOffset * 7))

  const monthNames = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']

  const result: WeekDayInfo[] = []
  for (let i = 0; i < 7; i++) {
    const d = new Date(monday)
    d.setDate(monday.getDate() + i)
    const iso = formatDateISO(d)
    result.push({
      day: DAY_NAMES[i],
      dateStr: iso,
      dayNumber: d.getDate(),
      monthLabel: monthNames[d.getMonth()],
      isToday: iso === todayISO,
      dateObj: d
    })
  }

  return result
}

// ==============================================================================
// BLOQUES FIJOS Y CONFIGURACIÓN REAL DE PABLO (INVERFIT)
// ==============================================================================

export interface FixedPresetBlock {
  id: string
  label: string
  startTime: string
  endTime: string
  shift: 'mañana' | 'tarde'
  isEvaluation?: boolean
}

/**
 * Plantilla fija de turnos semanales de Pablo (Lunes a Viernes).
 * Sesiones estructuradas de 75 minutos exactos según histórico oficial.
 */
export const PABLO_FIXED_PRESET_BLOCKS: FixedPresetBlock[] = [
  { id: 'b-0630', label: '06:30 – 07:45', startTime: '06:30', endTime: '07:45', shift: 'mañana' },
  { id: 'b-0800', label: '08:00 – 09:15', startTime: '08:00', endTime: '09:15', shift: 'mañana' },
  { id: 'b-0930', label: '09:30 – 10:45', startTime: '09:30', endTime: '10:45', shift: 'mañana' },
  { id: 'b-1100', label: '11:00 – 12:15', startTime: '11:00', endTime: '12:15', shift: 'mañana' },
  { id: 'b-1400', label: '14:00 – 15:15', startTime: '14:00', endTime: '15:15', shift: 'tarde' },
  { id: 'b-1600', label: '16:00 – 17:15', startTime: '16:00', endTime: '17:15', shift: 'tarde' },
  { id: 'b-1730', label: '17:30 – 18:45', startTime: '17:30', endTime: '18:45', shift: 'tarde' },
  { id: 'b-1900', label: '19:00 – 20:15', startTime: '19:00', endTime: '20:15', shift: 'tarde' },
  { id: 'b-2030', label: '20:30 – 21:45', startTime: '20:30', endTime: '21:45', shift: 'tarde' }
]

/**
 * Franjas exclusivas de Sábado reservadas para EVALUACIONES físicas y deportivas.
 */
export const PABLO_SATURDAY_EVALUATION_BLOCKS: FixedPresetBlock[] = [
  { id: 'eval-0630', label: '06:30 – 07:45 (Evaluación)', startTime: '06:30', endTime: '07:45', shift: 'mañana', isEvaluation: true },
  { id: 'eval-0800', label: '08:00 – 09:15 (Evaluación)', startTime: '08:00', endTime: '09:15', shift: 'mañana', isEvaluation: true },
  { id: 'eval-0930', label: '09:30 – 10:45 (Evaluación)', startTime: '09:30', endTime: '10:45', shift: 'mañana', isEvaluation: true },
  { id: 'eval-1100', label: '11:00 – 12:15 (Evaluación)', startTime: '11:00', endTime: '12:15', shift: 'mañana', isEvaluation: true }
]

/**
 * Límite estricto de alumnos por bloque en INVERFIT.
 * Permite control personalizado de técnica, cargas y seguridad.
 */
export const INVERFIT_MAX_STUDENTS_PER_BLOCK = 3

