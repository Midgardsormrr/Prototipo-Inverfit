import { httpClient, HttpClientError } from '../api/httpClient'
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
  type WeekDayInfo
} from './utils'

export class ScheduleApiService {
  // ==========================================
  // DISPONIBILIDAD (Coach Availability)
  // ==========================================

  async getAvailability(): Promise<DayAvailability[]> {
    const list = await httpClient.get<DayAvailability[]>('/schedule/availability')
    const safeList = Array.isArray(list) ? list : []
    return safeList.map(d => ({
      day: d.day,
      enabled: Boolean(d.enabled),
      intervals: Array.isArray(d?.intervals) ? d.intervals : []
    }))
  }

  async getAvailabilityForDay(day: DayOfWeek): Promise<DayAvailability> {
    const res = await httpClient.get<DayAvailability>(`/schedule/availability/${encodeURIComponent(day)}`)
    return {
      day: res?.day || day,
      enabled: Boolean(res?.enabled),
      intervals: Array.isArray(res?.intervals) ? res.intervals : []
    }
  }

  async saveAvailability(newAvailability: DayAvailability[]): Promise<{
    success: boolean
    message: string
    conflicts?: ScheduleBlock[]
  }> {
    try {
      // Save each day sequentially to respect atomic validation
      for (const dayAvail of (newAvailability || [])) {
        await httpClient.put<DayAvailability>(
          `/schedule/availability/${encodeURIComponent(dayAvail.day)}`,
          {
            enabled: Boolean(dayAvail.enabled),
            intervals: (dayAvail.intervals || []).map(i => ({
              id: i.id,
              startTime: i.startTime,
              endTime: i.endTime
            }))
          }
        )
      }

      return {
        success: true,
        message: 'Disponibilidad del entrenador actualizada exitosamente en el servidor.'
      }
    } catch (err: unknown) {
      if (err instanceof HttpClientError && err.statusCode === 409) {
        return {
          success: false,
          message: err.message,
          conflicts: []
        }
      }
      throw err
    }
  }

  // Client-side quick check
  validateSession(
    _day: DayOfWeek,
    startTime: string,
    endTime: string,
    _excludeBlockId?: string
  ): { valid: boolean; reason?: string } {
    const startM = timeToMinutes(startTime)
    const endM = timeToMinutes(endTime)

    if (endM <= startM) {
      return { valid: false, reason: 'La hora de término debe ser posterior a la hora de inicio.' }
    }
    if (endM - startM < 15) {
      return { valid: false, reason: 'La duración mínima de una sesión es de 15 minutos.' }
    }
    return { valid: true }
  }

  // ==========================================
  // SESIONES RECURRENTES
  // ==========================================

  async getAll(): Promise<ScheduleBlock[]> {
    // In API mode, get all active session blocks for current week
    const occurrences = await httpClient.get<SessionOccurrence[]>('/schedule/sessions')
    const list = Array.isArray(occurrences) ? occurrences : []
    return list.map(o => ({
      id: o.blockId,
      day: o.day,
      startTime: o.startTime,
      endTime: o.endTime,
      durationMinutes: o.durationMinutes,
      type: o.type,
      studentIds: Array.isArray(o.studentIds) ? [...o.studentIds] : [],
      isRecurring: o.isRecurring,
      notes: o.notes
    }))
  }

  async getBlocksByDay(day: DayOfWeek): Promise<ScheduleBlock[]> {
    const all = await this.getAll()
    return all.filter(b => b.day === day)
  }

  async getBlocksByStudentId(studentId: string): Promise<ScheduleBlock[]> {
    const occurrences = await httpClient.get<SessionOccurrence[]>(`/schedule/sessions?studentId=${encodeURIComponent(studentId)}`)
    const list = Array.isArray(occurrences) ? occurrences : []
    return list.map(o => ({
      id: o.blockId,
      day: o.day,
      startTime: o.startTime,
      endTime: o.endTime,
      durationMinutes: o.durationMinutes,
      type: o.type,
      studentIds: Array.isArray(o.studentIds) ? [...o.studentIds] : [],
      isRecurring: o.isRecurring,
      notes: o.notes
    }))
  }

  async assignSession(dto: AssignSessionDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }> {
    const endTime = dto.endTime || addMinutesToTime(dto.startTime, dto.durationMinutes || 75)
    const durationMinutes = timeToMinutes(endTime) - timeToMinutes(dto.startTime)

    try {
      if (dto.blockId) {
        // Update existing session
        const updated = await httpClient.put<ScheduleBlock>(`/schedule/sessions/${encodeURIComponent(dto.blockId)}`, {
          day: dto.day,
          startTime: dto.startTime,
          endTime: endTime,
          studentIds: dto.studentIds,
          isRecurring: dto.isRecurring !== false,
          startDate: dto.startDate,
          endDate: dto.endDate,
          notes: dto.notes
        })
        return {
          success: true,
          message: `Sesión actualizada para ${dto.day} de ${dto.startTime} a ${endTime}.`,
          block: updated
        }
      } else {
        // Create new session
        const created = await httpClient.post<ScheduleBlock>('/schedule/sessions', {
          day: dto.day,
          startTime: dto.startTime,
          endTime: endTime,
          durationMinutes: durationMinutes,
          studentIds: dto.studentIds,
          isRecurring: dto.isRecurring !== false,
          startDate: dto.startDate,
          endDate: dto.endDate,
          notes: dto.notes
        })
        return {
          success: true,
          message: `Sesión programada exitosamente para ${dto.day} de ${dto.startTime} a ${endTime}.`,
          block: created
        }
      }
    } catch (err: unknown) {
      const msg = err instanceof HttpClientError ? err.message : (err instanceof Error ? err.message : 'Error al programar sesión')
      return {
        success: false,
        message: msg
      }
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
    try {
      const updated = await httpClient.put<ScheduleBlock>(`/schedule/sessions/${encodeURIComponent(blockId)}`, updates)
      return {
        success: true,
        message: 'Sesión actualizada exitosamente.',
        block: updated
      }
    } catch (err: unknown) {
      const msg = err instanceof HttpClientError ? err.message : (err instanceof Error ? err.message : 'Error al actualizar sesión')
      return {
        success: false,
        message: msg
      }
    }
  }

  async removeStudentFromBlock(blockId: string, studentId: string): Promise<boolean> {
    // Fetch session, remove student ID, update
    const session = await httpClient.get<ScheduleBlock>(`/schedule/sessions/${encodeURIComponent(blockId)}`)
    const updatedStudentIDs = session.studentIds.filter(id => id !== studentId)
    await httpClient.put<ScheduleBlock>(`/schedule/sessions/${encodeURIComponent(blockId)}`, {
      studentIds: updatedStudentIDs
    })
    return true
  }

  async deleteBlock(blockId: string): Promise<boolean> {
    await httpClient.delete<{ message: string; id: string }>(`/schedule/sessions/${encodeURIComponent(blockId)}`)
    return true
  }

  // ==========================================
  // EXCEPCIONES DE CALENDARIO
  // ==========================================

  async getExceptions(from?: string, to?: string): Promise<ScheduleException[]> {
    const params = new URLSearchParams()
    if (from) params.set('from', from)
    if (to) params.set('to', to)
    const qs = params.toString()
    const res = await httpClient.get<ScheduleException[]>(qs ? `/schedule/exceptions?${qs}` : '/schedule/exceptions')
    return Array.isArray(res) ? res : []
  }

  async createException(dto: CreateExceptionDTO): Promise<ScheduleException> {
    return httpClient.post<ScheduleException>(
      `/schedule/sessions/${encodeURIComponent(dto.blockId)}/exceptions`,
      {
        date: dto.date,
        type: dto.type,
        newDate: dto.newDate,
        newDay: dto.newDay,
        newStartTime: dto.newStartTime,
        newEndTime: dto.newEndTime,
        notes: dto.notes
      }
    )
  }

  async deleteException(exceptionId: string, blockId?: string): Promise<boolean> {
    // If blockId is provided, use standard route /sessions/:id/exceptions/:exceptionId
    const targetBlockId = blockId || 'any'
    await httpClient.delete<{ message: string }>(
      `/schedule/sessions/${encodeURIComponent(targetBlockId)}/exceptions/${encodeURIComponent(exceptionId)}`
    )
    return true
  }

  // ==========================================
  // OCURRENCIAS Y AGENDA SEMANAL
  // ==========================================

  async getOccurrences(from: string, to: string, studentId?: string): Promise<SessionOccurrence[]> {
    const params = new URLSearchParams({ from, to })
    if (studentId) params.set('studentId', studentId)
    const raw = await httpClient.get<SessionOccurrence[]>(`/schedule/sessions?${params.toString()}`)
    const list = Array.isArray(raw) ? raw : []
    return list.map(o => ({
      ...o,
      studentIds: Array.isArray(o.studentIds) ? [...o.studentIds] : []
    }))
  }

  getWeekOccurrences(
    weekDays: WeekDayInfo[],
    blocksOverride?: ScheduleBlock[],
    exceptionsOverride?: ScheduleException[]
  ): SessionOccurrence[] {
    // For local immediate reactivity when modifying drafts
    if (blocksOverride && blocksOverride.length > 0) {
      // Use client-side projection algorithm
      const blocks = blocksOverride
      const exceptions = exceptionsOverride || []
      const occurrences: SessionOccurrence[] = []

      for (const dayInfo of weekDays) {
        const dayBlocks = blocks.filter(b => b.day === dayInfo.day)

        for (const block of dayBlocks) {
          const blockException = exceptions.find(e => e.originalBlockId === block.id && e.date === dayInfo.dateStr)

          if (blockException) {
            if (blockException.type === 'cancelled') continue
            if (blockException.type === 'rescheduled' && blockException.newDate && blockException.newDate !== dayInfo.dateStr) {
              continue
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

    return []
  }
}

export const scheduleApiService = new ScheduleApiService()
