import { storage } from './storage'
import type { Evaluation, CreateEvaluationDTO } from '@/types'

export class EvaluationService {
  async getAll(): Promise<Evaluation[]> {
    return storage.getEvaluations()
  }

  async getByStudentId(studentId: string): Promise<Evaluation[]> {
    const list = storage.getEvaluations()
    return list.filter(e => e.studentId === studentId)
  }

  async getById(id: string): Promise<Evaluation | null> {
    const list = storage.getEvaluations()
    return list.find(e => e.id === id) || null
  }

  async create(dto: CreateEvaluationDTO): Promise<Evaluation> {
    const list = storage.getEvaluations()
    const student = storage.getStudents().find(s => s.id === dto.studentId)
    const studentName = student ? student.fullName : 'Alumno no especificado'

    const newEval: Evaluation = {
      id: `eval-${Date.now().toString().slice(-6)}`,
      studentId: dto.studentId,
      studentName,
      date: dto.date || new Date().toISOString().split('T')[0],
      type: dto.type,
      status: dto.status || 'Registrada',
      metrics: dto.metrics || [],
      observations: dto.observations || '',
      conductedBy: 'Pablo',
      createdAt: new Date().toISOString().split('T')[0]
    }

    storage.setEvaluations([newEval, ...list])
    return newEval
  }

  async update(id: string, updates: Partial<Evaluation>): Promise<Evaluation | null> {
    const list = storage.getEvaluations()
    const idx = list.findIndex(e => e.id === id)
    if (idx === -1) return null

    const updated = { ...list[idx], ...updates }
    list[idx] = updated
    storage.setEvaluations([...list])
    return updated
  }

  async delete(id: string): Promise<boolean> {
    const list = storage.getEvaluations()
    const filtered = list.filter(e => e.id !== id)
    if (filtered.length === list.length) return false
    storage.setEvaluations(filtered)
    return true
  }
}

export const evaluationService = new EvaluationService()
