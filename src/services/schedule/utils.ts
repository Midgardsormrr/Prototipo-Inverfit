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
