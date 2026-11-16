import type { Student, CreateStudentDTO } from '@/types'
import { studentMockService, StudentMockService } from './students.mock'
import { studentApiService, StudentApiService, type StudentFilterParams } from './students.api'

export interface IStudentService {
  getAll(filter?: StudentFilterParams): Promise<Student[]>
  getById(id: string): Promise<Student | null>
  create(dto: CreateStudentDTO): Promise<Student>
  update(id: string, updates: Partial<Student>): Promise<Student | null>
  delete(id: string): Promise<boolean>
}

/**
 * Determina si la aplicación está configurada para operar con la API real Go + MongoDB
 * o en modo de demostración con localStorage.
 * Valor por defecto: 'mock' (seguro para Cloudflare Pages y prototipo offline).
 */
export function getStudentsDataSource(): 'mock' | 'api' {
  const source = import.meta.env.VITE_STUDENTS_DATA_SOURCE
  return source === 'api' ? 'api' : 'mock'
}

export function isStudentsApiMode(): boolean {
  return getStudentsDataSource() === 'api'
}

/**
 * Servicio unificado de alumnos.
 * Delega automáticamente a StudentApiService si VITE_STUDENTS_DATA_SOURCE=api,
 * o a StudentMockService si está en modo 'mock'.
 */
class StudentServiceDelegator implements IStudentService {
  private get activeService(): IStudentService {
    return isStudentsApiMode() ? studentApiService : studentMockService
  }

  getAll(filter?: StudentFilterParams): Promise<Student[]> {
    return this.activeService.getAll(filter)
  }

  getById(id: string): Promise<Student | null> {
    return this.activeService.getById(id)
  }

  create(dto: CreateStudentDTO): Promise<Student> {
    return this.activeService.create(dto)
  }

  update(id: string, updates: Partial<Student>): Promise<Student | null> {
    return this.activeService.update(id, updates)
  }

  delete(id: string): Promise<boolean> {
    return this.activeService.delete(id)
  }
}

export const studentService = new StudentServiceDelegator()
export { studentMockService, StudentMockService, studentApiService, StudentApiService }
