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
import { scheduleMockService, ScheduleMockService } from './schedule.mock'
import { scheduleApiService, ScheduleApiService } from './schedule.api'
import {
  timeToMinutes,
  minutesToTime,
  addMinutesToTime,
  formatDateISO,
  DAY_NAMES,
  getWeekDates,
  type WeekDayInfo
} from './utils'

export interface IScheduleService {
  getAvailability(): Promise<DayAvailability[]>
  getAvailabilityForDay(day: DayOfWeek): Promise<DayAvailability>
  saveAvailability(newAvailability: DayAvailability[]): Promise<{
    success: boolean
    message: string
    conflicts?: ScheduleBlock[]
  }>
  validateSession(
    day: DayOfWeek,
    startTime: string,
    endTime: string,
    excludeBlockId?: string
  ): { valid: boolean; reason?: string }
  getAll(): Promise<ScheduleBlock[]>
  getBlocksByDay(day: DayOfWeek): Promise<ScheduleBlock[]>
  getBlocksByStudentId(studentId: string): Promise<ScheduleBlock[]>
  assignSession(dto: AssignSessionDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }>
  assignStudent(dto: AssignStudentDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }>
  updateSession(
    blockId: string,
    updates: Partial<ScheduleBlock>
  ): Promise<{ success: boolean; message: string; block?: ScheduleBlock }>
  removeStudentFromBlock(blockId: string, studentId: string): Promise<boolean>
  deleteBlock(blockId: string): Promise<boolean>
  getExceptions(from?: string, to?: string): Promise<ScheduleException[]>
  createException(dto: CreateExceptionDTO): Promise<ScheduleException>
  deleteException(exceptionId: string, blockId?: string): Promise<boolean>
  getOccurrences(from: string, to: string, studentId?: string): Promise<SessionOccurrence[]>
  getWeekOccurrences(
    weekDays: WeekDayInfo[],
    blocksOverride?: ScheduleBlock[],
    exceptionsOverride?: ScheduleException[]
  ): SessionOccurrence[]
}

/**
 * Determina si la aplicación utiliza la API real Go + MongoDB o el almacenamiento local (mock).
 * Valor predeterminado: 'mock' (seguro para Cloudflare Pages y demo offline).
 */
export function getScheduleDataSource(): 'mock' | 'api' {
  const source = import.meta.env.VITE_SCHEDULE_DATA_SOURCE
  return source === 'api' ? 'api' : 'mock'
}

export function isScheduleApiMode(): boolean {
  return getScheduleDataSource() === 'api'
}

class ScheduleServiceDelegator implements IScheduleService {
  private get activeService(): IScheduleService {
    return isScheduleApiMode() ? scheduleApiService : scheduleMockService
  }

  getAvailability(): Promise<DayAvailability[]> {
    return this.activeService.getAvailability()
  }

  getAvailabilityForDay(day: DayOfWeek): Promise<DayAvailability> {
    return this.activeService.getAvailabilityForDay(day)
  }

  saveAvailability(newAvailability: DayAvailability[]): Promise<{
    success: boolean
    message: string
    conflicts?: ScheduleBlock[]
  }> {
    return this.activeService.saveAvailability(newAvailability)
  }

  validateSession(
    day: DayOfWeek,
    startTime: string,
    endTime: string,
    excludeBlockId?: string
  ): { valid: boolean; reason?: string } {
    return this.activeService.validateSession(day, startTime, endTime, excludeBlockId)
  }

  getAll(): Promise<ScheduleBlock[]> {
    return this.activeService.getAll()
  }

  getBlocksByDay(day: DayOfWeek): Promise<ScheduleBlock[]> {
    return this.activeService.getBlocksByDay(day)
  }

  getBlocksByStudentId(studentId: string): Promise<ScheduleBlock[]> {
    return this.activeService.getBlocksByStudentId(studentId)
  }

  assignSession(dto: AssignSessionDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }> {
    return this.activeService.assignSession(dto)
  }

  assignStudent(dto: AssignStudentDTO): Promise<{
    success: boolean
    message: string
    block?: ScheduleBlock
  }> {
    return this.activeService.assignStudent(dto)
  }

  updateSession(
    blockId: string,
    updates: Partial<ScheduleBlock>
  ): Promise<{ success: boolean; message: string; block?: ScheduleBlock }> {
    return this.activeService.updateSession(blockId, updates)
  }

  removeStudentFromBlock(blockId: string, studentId: string): Promise<boolean> {
    return this.activeService.removeStudentFromBlock(blockId, studentId)
  }

  deleteBlock(blockId: string): Promise<boolean> {
    return this.activeService.deleteBlock(blockId)
  }

  getExceptions(from?: string, to?: string): Promise<ScheduleException[]> {
    return this.activeService.getExceptions(from, to)
  }

  createException(dto: CreateExceptionDTO): Promise<ScheduleException> {
    return this.activeService.createException(dto)
  }

  deleteException(exceptionId: string, blockId?: string): Promise<boolean> {
    return this.activeService.deleteException(exceptionId, blockId)
  }

  getOccurrences(from: string, to: string, studentId?: string): Promise<SessionOccurrence[]> {
    return this.activeService.getOccurrences(from, to, studentId)
  }

  getWeekOccurrences(
    weekDays: WeekDayInfo[],
    blocksOverride?: ScheduleBlock[],
    exceptionsOverride?: ScheduleException[]
  ): SessionOccurrence[] {
    return this.activeService.getWeekOccurrences(weekDays, blocksOverride, exceptionsOverride)
  }
}

export const scheduleService = new ScheduleServiceDelegator()
export {
  scheduleMockService,
  ScheduleMockService,
  scheduleApiService,
  ScheduleApiService,
  timeToMinutes,
  minutesToTime,
  addMinutesToTime,
  formatDateISO,
  DAY_NAMES,
  getWeekDates,
  type WeekDayInfo
}
