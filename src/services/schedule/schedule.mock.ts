import { storage } from '../storage'
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
import {
  timeToMinutes,
  addMinutesToTime,
  DAY_NAMES,
  type WeekDayInfo
} from './utils'

export class ScheduleMockService {
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

      if (b.startTime === startTime && b.endTime === endTime) {
        continue
      }

      // Overlap formula: startM < bEndM && endM > bStartM
      if (startM < bEndM && endM > bStartM) {
        return {
          valid: false,
          reason: `Conflicto de horario: se superpone con la sesión existente (${b.startTime} — ${b.endTime}).`
        }
      }
    }

    return { valid: true }
  }

  // ==========================================
  // SESIONES RECURRENTES
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
      targetBlock = list.find(b => b.day === dto.day && b.startTime === dto.startTime && b.endTime === endTime)
    }

    if (targetBlock) {
      if (dto.blockId) {
        targetBlock.studentIds = [...dto.studentIds]
      } else {
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
    this.syncStudentFixedSchedules()

    return {
      success: true,
      message: `Sesión programada exitosamente para ${dto.day} de ${dto.startTime} a ${endTime}.`,
      block: targetBlock
    }
  }

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

    const exceptions = storage.getExceptions().filter(e => e.originalBlockId !== blockId)
    storage.setExceptions(exceptions)

    this.syncStudentFixedSchedules()
    return true
  }

  // ==========================================
  // EXCEPCIONES DE CALENDARIO
  // ==========================================

  async getExceptions(): Promise<ScheduleException[]> {
    return storage.getExceptions()
  }

  async createException(dto: CreateExceptionDTO): Promise<ScheduleException> {
    const exceptions = storage.getExceptions()
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

  async deleteException(exceptionId: string, _blockId?: string): Promise<boolean> {
    const exceptions = storage.getExceptions().filter(e => e.id !== exceptionId)
    storage.setExceptions(exceptions)
    return true
  }

  // ==========================================
  // OCURRENCIAS SEMANALES
  // ==========================================

  getWeekOccurrences(
    weekDays: WeekDayInfo[],
    blocksOverride?: ScheduleBlock[],
    exceptionsOverride?: ScheduleException[]
  ): SessionOccurrence[] {
    const blocks = blocksOverride || storage.getSchedule()
    const exceptions = exceptionsOverride || storage.getExceptions()
    const occurrences: SessionOccurrence[] = []

    for (const dayInfo of weekDays) {
      const dayBlocks = blocks.filter(b => b.day === dayInfo.day)

      for (const block of dayBlocks) {
        const blockException = exceptions.find(e => e.originalBlockId === block.id && e.date === dayInfo.dateStr)

        if (blockException) {
          if (blockException.type === 'cancelled') {
            continue
          }
          if (blockException.type === 'rescheduled') {
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

    return occurrences.sort((a, b) => timeToMinutes(a.startTime) - timeToMinutes(b.startTime))
  }

  async getOccurrences(from: string, to: string, studentId?: string): Promise<SessionOccurrence[]> {
    const blocks = storage.getSchedule()
    const filteredBlocks = studentId ? blocks.filter(b => b.studentIds.includes(studentId)) : blocks
    // Project using mock data
    const weekDays: WeekDayInfo[] = []
    const start = new Date(from + 'T12:00:00Z')
    const end = new Date(to + 'T12:00:00Z')
    const cur = new Date(start)
    const dayMap: DayOfWeek[] = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']

    while (cur <= end) {
      const y = cur.getUTCFullYear()
      const m = String(cur.getUTCMonth() + 1).padStart(2, '0')
      const d = String(cur.getUTCDate()).padStart(2, '0')
      const iso = `${y}-${m}-${d}`
      weekDays.push({
        day: dayMap[cur.getUTCDay()],
        dateStr: iso,
        dayNumber: cur.getUTCDate(),
        monthLabel: '',
        isToday: false,
        dateObj: new Date(cur)
      })
      cur.setUTCDate(cur.getUTCDate() + 1)
    }

    return this.getWeekOccurrences(weekDays, filteredBlocks)
  }

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

export const scheduleMockService = new ScheduleMockService()
