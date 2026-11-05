import { storage } from './storage'
import type {
  ScheduleBlock,
  DayOfWeek,
  DayAvailability,
  ScheduleException,
  SessionOccurrence,
  AssignSessionDTO,
  AssignStudentDTO,
  CreateExceptionDTO
} from '@/types'

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

// Map day index (0=Dom, 1=Lun, ..., 6=Sáb) to DayOfWeek
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

export class ScheduleService {
  // ==========================================
  // DISPONIBILIDAD (Coach Availability)
  // ==========================================

  async getAvailability(): Promise<DayAvailability[]> {
    return storage.getAvailability()
  }

  async getAvailabilityForDay(day: DayOfWeek): Promise<DayAvailability> {
    const all = storage.getAvailability()
    const found = all.find(a => a.day === day)
    if (found) return found
    return { day, enabled: false, intervals: [] }
  }

  /**
   * Save coach availability. Checks if reducing intervals conflicts with existing sessions.
   */
  async saveAvailability(newAvailability: DayAvailability[]): Promise<{
    success: boolean
    message: string
    conflicts?: ScheduleBlock[]
  }> {
    const sessions = storage.getSchedule()
    const conflicts: ScheduleBlock[] = []

    for (const session of sessions) {
      const dayAvail = newAvailability.find(a => a.day === session.day)
      const sessionStart = timeToMinutes(session.startTime)
      const sessionEnd = timeToMinutes(session.endTime)

      if (!dayAvail || !dayAvail.enabled || dayAvail.intervals.length === 0) {
        conflicts.push(session)
        continue
      }

      // Must be within at least one interval
      const fits = dayAvail.intervals.some(interval => {
        const intStart = timeToMinutes(interval.startTime)
        const intEnd = timeToMinutes(interval.endTime)
        return sessionStart >= intStart && sessionEnd <= intEnd
      })

      if (!fits) {
        conflicts.push(session)
      }
    }

    if (conflicts.length > 0) {
      return {
        success: false,
        message: `La nueva disponibilidad deja ${conflicts.length} sesión(es) existente(s) fuera de horario. Resuelve los conflictos antes de guardar.`,
        conflicts
      }
    }

    storage.setAvailability(newAvailability)
    return {
      success: true,
      message: 'Disponibilidad del entrenador actualizada exitosamente.'
    }
  }

  // ==========================================
  // VALIDACIÓN DE SESIONES
  // ==========================================

  /**
   * Validate a session:
   * 1. End time must be strictly after start time.
   * 2. Must be within the coach's configured availability for that day.
   * 3. No overlap with other sessions.
   * 4. NO mandatory breaks (consecutive sessions allowed!).
   */
  validateSession(
    day: DayOfWeek,
    startTime: string,
    endTime: string,
    excludeBlockId?: string
  ): { valid: boolean; reason?: string } {
    const startM = timeToMinutes(startTime)
    const endM = timeToMinutes(endTime)

    if (endM <= startM) {
      return { valid: false, reason: 'La hora de término debe ser posterior a la hora de inicio.' }
    }

    if (endM - startM < 15) {
      return { valid: false, reason: 'La duración mínima de una sesión es de 15 minutos.' }
    }

    // Check availability
    const dayAvail = storage.getAvailability().find(a => a.day === day)
    if (!dayAvail || !dayAvail.enabled || dayAvail.intervals.length === 0) {
      return {
        valid: false,
        reason: `El entrenador no tiene disponibilidad configurada para el día ${day}. Configura su disponibilidad primero.`
      }
    }

    const fitsInAvailability = dayAvail.intervals.some(interval => {
      const intStart = timeToMinutes(interval.startTime)
      const intEnd = timeToMinutes(interval.endTime)
      return startM >= intStart && endM <= intEnd
    })

    if (!fitsInAvailability) {
      const intervalsDesc = dayAvail.intervals.map(i => `${i.startTime} a ${i.endTime}`).join(' y ')
      return {
        valid: false,
        reason: `El horario (${startTime} — ${endTime}) queda fuera de la disponibilidad configurada para el ${day} (${intervalsDesc}).`
      }
    }

    // Check overlap with other blocks on this day
    const blocks = storage.getSchedule().filter(b => b.day === day && b.id !== excludeBlockId)

    for (const b of blocks) {
      const bStartM = timeToMinutes(b.startTime)
      const bEndM = timeToMinutes(b.endTime)

      // Exact match is permitted if grouping students into the exact same block
      if (b.startTime === startTime && b.endTime === endTime) {
        continue
      }

      // REAL OVERLAP CHECK: [startM, endM] overlaps [bStartM, bEndM] if startM < bEndM && endM > bStartM
      if (startM < bEndM && endM > bStartM) {
        return {
          valid: false,
          reason: `Conflicto de horario: se superpone con la sesión existente (${b.startTime} — ${b.endTime}).`
        }
      }
      // Note: NO break validation! Consecutive blocks (startM == bEndM or endM == bStartM) are completely valid.
    }

    return { valid: true }
  }

  // ==========================================
  // SESIONES RECURRENTES (Schedule Blocks)
  // ==========================================

  async getAll(): Promise<ScheduleBlock[]> {
    return storage.getSchedule()
  }

  async getBlocksByDay(day: DayOfWeek): Promise<ScheduleBlock[]> {
    const list = storage.getSchedule()
    return list.filter(b => b.day === day)
  }

  async getBlocksByStudentId(studentId: string): Promise<ScheduleBlock[]> {
    const list = storage.getSchedule()
    return list.filter(b => b.studentIds.includes(studentId))
  }

  /**
   * Assign students to a session (new or existing block)
   */
  async assignSession(dto: AssignSessionDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }> {
    const endTime = dto.endTime || addMinutesToTime(dto.startTime, dto.durationMinutes || 75)
    const durationMinutes = timeToMinutes(endTime) - timeToMinutes(dto.startTime)

    const validation = this.validateSession(dto.day, dto.startTime, endTime, dto.blockId)
    if (!validation.valid) {
      return { success: false, message: validation.reason || 'Conflicto de horario' }
    }

    const list = storage.getSchedule()
    let targetBlock: ScheduleBlock | undefined

    if (dto.blockId) {
      targetBlock = list.find(b => b.id === dto.blockId)
    } else {
      // Find exact existing block at this day and startTime
      targetBlock = list.find(b => b.day === dto.day && b.startTime === dto.startTime && b.endTime === endTime)
    }

    if (targetBlock) {
      if (dto.blockId) {
        // Explicitly updating an existing block: set exact selected student IDs
        targetBlock.studentIds = [...dto.studentIds]
      } else {
        // Add students not already present
        dto.studentIds.forEach(sid => {
          if (!targetBlock!.studentIds.includes(sid)) {
            targetBlock!.studentIds.push(sid)
          }
        })
      }
      targetBlock.day = dto.day
      targetBlock.startTime = dto.startTime
      targetBlock.endTime = endTime
      targetBlock.durationMinutes = durationMinutes
      targetBlock.type = targetBlock.studentIds.length > 1 ? 'Sesión grupal' : 'Individual'
      if (dto.notes) targetBlock.notes = dto.notes
    } else {
      targetBlock = {
        id: `blk-${dto.day.toLowerCase()}-${dto.startTime.replace(':', '')}-${Date.now().toString().slice(-4)}`,
        day: dto.day,
        startTime: dto.startTime,
        endTime: endTime,
        durationMinutes: durationMinutes,
        type: dto.studentIds.length > 1 ? 'Sesión grupal' : 'Individual',
        studentIds: [...dto.studentIds],
        isRecurring: dto.isRecurring !== false,
        notes: dto.notes || `Sesión de ${dto.day} ${dto.startTime} a ${endTime}.`
      }
      list.push(targetBlock)
    }

    storage.setSchedule([...list])

    // Update fixed schedules for all affected students
    this.syncStudentFixedSchedules()

    return {
      success: true,
      message: `Sesión programada exitosamente para ${dto.day} de ${dto.startTime} a ${endTime}.`,
      block: targetBlock
    }
  }

  // Backward-compatible assignStudent wrapper
  async assignStudent(dto: AssignStudentDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }> {
    const endTime = dto.endTime || addMinutesToTime(dto.startTime, dto.durationMinutes || 75)
    return this.assignSession({
      blockId: dto.blockId,
      studentIds: [dto.studentId],
      day: dto.day,
      startTime: dto.startTime,
      endTime: endTime,
      isRecurring: dto.isRecurring !== false,
      notes: dto.notes
    })
  }

  async updateSession(
    blockId: string,
    updates: Partial<ScheduleBlock>
  ): Promise<{ success: boolean; message: string; block?: ScheduleBlock }> {
    const list = storage.getSchedule()
    const blockIndex = list.findIndex(b => b.id === blockId)
    if (blockIndex === -1) {
      return { success: false, message: 'Sesión no encontrada.' }
    }

    const current = list[blockIndex]
    const newDay = updates.day || current.day
    const newStart = updates.startTime || current.startTime
    const newEnd = updates.endTime || current.endTime

    const validation = this.validateSession(newDay, newStart, newEnd, blockId)
    if (!validation.valid) {
      return { success: false, message: validation.reason || 'Conflicto de horario' }
    }

    const durationMinutes = timeToMinutes(newEnd) - timeToMinutes(newStart)

    const updated: ScheduleBlock = {
      ...current,
      ...updates,
      day: newDay,
      startTime: newStart,
      endTime: newEnd,
      durationMinutes,
      type: (updates.studentIds || current.studentIds).length > 1 ? 'Sesión grupal' : 'Individual'
    }

    list[blockIndex] = updated
    storage.setSchedule([...list])
    this.syncStudentFixedSchedules()

    return {
      success: true,
      message: 'Sesión actualizada exitosamente.',
      block: updated
    }
  }

  async removeStudentFromBlock(blockId: string, studentId: string): Promise<boolean> {
    const list = storage.getSchedule()
    const block = list.find(b => b.id === blockId)
    if (!block) return false

    block.studentIds = block.studentIds.filter(id => id !== studentId)
    if (block.studentIds.length <= 1) {
      block.type = 'Individual'
    }

    storage.setSchedule([...list])
    this.syncStudentFixedSchedules()
    return true
  }

  async deleteBlock(blockId: string): Promise<boolean> {
    const list = storage.getSchedule()
    const filtered = list.filter(b => b.id !== blockId)
    storage.setSchedule(filtered)

    // Remove any exceptions tied to this block
    const exceptions = storage.getExceptions().filter(e => e.originalBlockId !== blockId)
    storage.setExceptions(exceptions)

    this.syncStudentFixedSchedules()
    return true
  }

  // ==========================================
  // EXCEPCIONES DE CALENDARIO (Calendar Exceptions)
  // ==========================================

  async getExceptions(): Promise<ScheduleException[]> {
    return storage.getExceptions()
  }

  async createException(dto: CreateExceptionDTO): Promise<ScheduleException> {
    const exceptions = storage.getExceptions()

    // Remove any prior exception for the exact same block and date
    const filtered = exceptions.filter(e => !(e.originalBlockId === dto.blockId && e.date === dto.date))

    const newException: ScheduleException = {
      id: `exc-${Date.now()}-${Math.random().toString(36).substring(2, 6)}`,
      originalBlockId: dto.blockId,
      date: dto.date,
      type: dto.type,
      newDate: dto.newDate,
      newDay: dto.newDay,
      newStartTime: dto.newStartTime,
      newEndTime: dto.newEndTime,
      notes: dto.notes,
      createdAt: new Date().toISOString()
    }

    filtered.push(newException)
    storage.setExceptions(filtered)
    return newException
  }

  async deleteException(exceptionId: string): Promise<boolean> {
    const exceptions = storage.getExceptions().filter(e => e.id !== exceptionId)
    storage.setExceptions(exceptions)
    return true
  }

  // ==========================================
  // CÁLCULO DE OCURRENCIAS PARA LA SEMANA
  // ==========================================

  /**
   * Given the 7 days of a specific week, calculates the concrete session occurrences
   * applying recurring blocks and calendar exceptions. Accepts optional overrides
   * for reactive store bindings.
   */
  getWeekOccurrences(
    weekDays: WeekDayInfo[],
    blocksOverride?: ScheduleBlock[],
    exceptionsOverride?: ScheduleException[]
  ): SessionOccurrence[] {
    const blocks = blocksOverride || storage.getSchedule()
    const exceptions = exceptionsOverride || storage.getExceptions()
    const occurrences: SessionOccurrence[] = []

    for (const dayInfo of weekDays) {
      // 1. Recurring blocks scheduled for this day of week
      const dayBlocks = blocks.filter(b => b.day === dayInfo.day)

      for (const block of dayBlocks) {
        // Check if there is an exception for this specific date
        const blockException = exceptions.find(e => e.originalBlockId === block.id && e.date === dayInfo.dateStr)

        if (blockException) {
          if (blockException.type === 'cancelled') {
            // Cancelled for this date: do not render as active occurrence
            continue
          }
          if (blockException.type === 'rescheduled') {
            // If rescheduled away to another date, do not render here unless newDate is same day
            if (blockException.newDate && blockException.newDate !== dayInfo.dateStr) {
              continue
            }
          }
        }

        occurrences.push({
          instanceId: `inst-${block.id}-${dayInfo.dateStr}`,
          blockId: block.id,
          day: dayInfo.day,
          date: dayInfo.dateStr,
          startTime: block.startTime,
          endTime: block.endTime,
          durationMinutes: block.durationMinutes,
          type: block.type,
          studentIds: [...block.studentIds],
          isRecurring: block.isRecurring !== false,
          isException: false,
          notes: block.notes
        })
      }

      // 2. Add sessions that were rescheduled TO this day from another day
      const movedToThisDay = exceptions.filter(e => e.type === 'rescheduled' && e.newDate === dayInfo.dateStr)

      for (const exc of movedToThisDay) {
        const originalBlock = blocks.find(b => b.id === exc.originalBlockId)
        if (!originalBlock) continue

        const startTime = exc.newStartTime || originalBlock.startTime
        const endTime = exc.newEndTime || originalBlock.endTime
        const durationMinutes = timeToMinutes(endTime) - timeToMinutes(startTime)

        occurrences.push({
          instanceId: `inst-exc-${exc.id}`,
          blockId: originalBlock.id,
          day: dayInfo.day,
          date: dayInfo.dateStr,
          startTime,
          endTime,
          durationMinutes,
          type: originalBlock.type,
          studentIds: [...originalBlock.studentIds],
          isRecurring: true,
          isException: true,
          exceptionType: 'rescheduled',
          exceptionId: exc.id,
          notes: exc.notes || `Reagendado para esta semana desde ${originalBlock.day}.`
        })
      }
    }

    // Sort chronologically by start time
    return occurrences.sort((a, b) => timeToMinutes(a.startTime) - timeToMinutes(b.startTime))
  }

  // ==========================================
  // SYNC STUDENT FIXED SCHEDULES
  // ==========================================

  private syncStudentFixedSchedules(): void {
    const list = storage.getSchedule()
    const students = storage.getStudents()

    const updatedStudents = students.map(s => {
      const studentBlocks = list.filter(b => b.studentIds.includes(s.id))
      if (studentBlocks.length === 0) {
        return s
      }
      const summary = studentBlocks
        .sort((a, b) => DAY_NAMES.indexOf(a.day) - DAY_NAMES.indexOf(b.day))
        .map(b => `${b.day} ${b.startTime}`)
        .join(', ')
      return { ...s, fixedSchedule: summary || 'Días y horas por confirmar' }
    })

    storage.setStudents(updatedStudents)
  }
}

export const scheduleService = new ScheduleService()
