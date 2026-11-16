import { httpClient, HttpClientError } from '../api/httpClient'
import type { Student, CreateStudentDTO } from '@/types'

export interface StudentFilterParams {
  name?: string
  status?: string
  frequency?: string
}

export class StudentApiService {
  /**
   * Obtiene la lista de alumnos activos desde la API Go / MongoDB.
   * Admite filtros opcionales de nombre, estado y frecuencia.
   */
  async getAll(filter?: StudentFilterParams): Promise<Student[]> {
    const params: Record<string, string> = {}
    if (filter?.name) params.name = filter.name
    if (filter?.status) params.status = filter.status
    if (filter?.frequency) params.frequency = filter.frequency

    return await httpClient.get<Student[]>('students', params)
  }

  /**
   * Obtiene la ficha de un alumno por su ID (ObjectID hexadecimal).
   * Retorna null si el alumno no existe o está archivado.
   */
  async getById(id: string): Promise<Student | null> {
    try {
      return await httpClient.get<Student>(`students/${encodeURIComponent(id)}`)
    } catch (err: unknown) {
      if (err instanceof HttpClientError && err.statusCode === 404) {
        return null
      }
      throw err
    }
  }

  /**
   * Registra un nuevo alumno en MongoDB.
   * En caso de éxito retorna el Student creado con su ID generado por MongoDB.
   */
  async create(dto: CreateStudentDTO): Promise<Student> {
    return await httpClient.post<Student>('students', dto)
  }

  /**
   * Actualiza los datos de un alumno existente mediante PUT.
   */
  async update(id: string, updates: Partial<Student>): Promise<Student | null> {
    try {
      return await httpClient.put<Student>(`students/${encodeURIComponent(id)}`, updates)
    } catch (err: unknown) {
      if (err instanceof HttpClientError && err.statusCode === 404) {
        return null
      }
      throw err
    }
  }

  /**
   * Elimina lógicamente (archiva) un alumno en MongoDB.
   * Preserva el registro histórico en base de datos.
   */
  async delete(id: string): Promise<boolean> {
    try {
      await httpClient.delete(`students/${encodeURIComponent(id)}`)
      return true
    } catch (err: unknown) {
      if (err instanceof HttpClientError && err.statusCode === 404) {
        return false
      }
      throw err
    }
  }
}

export const studentApiService = new StudentApiService()
