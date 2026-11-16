import { storage } from '../storage'
import type { Student, CreateStudentDTO, Payment, Evaluation, ContractConsent } from '@/types'

export class StudentMockService {
  async getAll(): Promise<Student[]> {
    return storage.getStudents()
  }

  async getById(id: string): Promise<Student | null> {
    const list = storage.getStudents()
    return list.find(s => s.id === id) || null
  }

  async create(dto: CreateStudentDTO): Promise<Student> {
    const list = storage.getStudents()
    const newId = `alu-${String(list.length + 1).padStart(2, '0')}-${Date.now().toString().slice(-4)}`

    const newStudent: Student = {
      id: newId,
      fullName: dto.fullName.trim(),
      email: dto.email?.trim() || `${dto.fullName.toLowerCase().replace(/\s+/g, '')}@ejemplo.com`,
      phone: dto.phone?.trim() || '+56 9 Por definir',
      planName: dto.planName || 'Plan estándar',
      frequency: dto.frequency,
      membershipStatus: dto.membershipStatus || 'Activo',
      fixedSchedule: dto.fixedSchedule || 'Días y horas por confirmar',
      membershipStartDate: dto.membershipStartDate || new Date().toISOString().split('T')[0],
      membershipDueDate: dto.membershipDueDate || 'Por definir',
      notes: dto.notes || '',
      healthConditions: dto.healthConditions || 'Sin observaciones de salud registradas',
      emergencyContact: dto.emergencyContact || 'Por registrar',
      createdAt: new Date().toISOString().split('T')[0]
    }

    // Save student
    storage.setStudents([newStudent, ...list])

    // Create synchronized initial payment record for mock demonstration
    const payments = storage.getPayments()
    const newPayment: Payment = {
      id: `pay-${Date.now().toString().slice(-6)}`,
      studentId: newStudent.id,
      studentName: newStudent.fullName,
      planName: newStudent.planName,
      frequency: newStudent.frequency,
      dueDate: 'Por definir',
      amount: null,
      status: 'Pendiente',
      paymentDate: null,
      paymentMethod: 'Por definir',
      reminderSentCount: 0,
      lastReminderAt: null,
      receiptSent: false,
      receiptSentAt: null,
      notes: 'Registro inicial generado al registrar alumno.'
    }
    storage.setPayments([newPayment, ...payments])

    // Create synchronized initial evaluation record for mock demonstration
    const evaluations = storage.getEvaluations()
    const newEval: Evaluation = {
      id: `eval-${Date.now().toString().slice(-6)}`,
      studentId: newStudent.id,
      studentName: newStudent.fullName,
      date: 'Por registrar',
      type: 'Evaluación inicial',
      status: 'Pendiente',
      metrics: [
        { name: 'Peso corporal', value: 'Por definir', unit: 'kg' },
        { name: 'Estatura', value: 'Por definir', unit: 'cm' },
        { name: 'Movilidad funcional', value: 'Por evaluar' },
        { name: 'Fuerza base', value: 'Por evaluar' }
      ],
      observations: 'Evaluación inicial programada al ingreso.',
      conductedBy: 'Pablo',
      createdAt: new Date().toISOString().split('T')[0]
    }
    storage.setEvaluations([newEval, ...evaluations])

    // Create synchronized contract consent record for mock demonstration
    const contracts = storage.getContracts()
    const newContract: ContractConsent = {
      id: `cnt-${Date.now().toString().slice(-6)}`,
      studentId: newStudent.id,
      studentName: newStudent.fullName,
      documentType: 'Términos y condiciones',
      status: 'Pendiente',
      date: 'Por registrar',
      termsVersion: 'v1.0 (2026)',
      signedCopyAvailable: false,
      notes: 'Pendiente de aceptación de términos del servicio.'
    }
    storage.setContracts([newContract, ...contracts])

    return newStudent
  }

  async update(id: string, updates: Partial<Student>): Promise<Student | null> {
    const list = storage.getStudents()
    const idx = list.findIndex(s => s.id === id)
    if (idx === -1) return null

    const oldStudent = list[idx]
    const updatedStudent: Student = { ...oldStudent, ...updates }
    list[idx] = updatedStudent
    storage.setStudents([...list])

    // If name changed, synchronize in payments, evaluations, and contracts
    if (updates.fullName && updates.fullName !== oldStudent.fullName) {
      const newName = updates.fullName

      const payments = storage.getPayments().map(p =>
        p.studentId === id ? { ...p, studentName: newName, planName: updatedStudent.planName, frequency: updatedStudent.frequency } : p
      )
      storage.setPayments(payments)

      const evaluations = storage.getEvaluations().map(e =>
        e.studentId === id ? { ...e, studentName: newName } : e
      )
      storage.setEvaluations(evaluations)

      const contracts = storage.getContracts().map(c =>
        c.studentId === id ? { ...c, studentName: newName } : c
      )
      storage.setContracts(contracts)
    }

    return updatedStudent
  }

  async delete(id: string): Promise<boolean> {
    const list = storage.getStudents()
    const filtered = list.filter(s => s.id !== id)
    if (filtered.length === list.length) return false
    storage.setStudents(filtered)

    // Remove student from any schedule blocks
    const schedule = storage.getSchedule().map(b => ({
      ...b,
      studentIds: b.studentIds.filter(sid => sid !== id)
    }))
    storage.setSchedule(schedule)

    // Remove related payments, evaluations, and contracts in mock mode
    storage.setPayments(storage.getPayments().filter(p => p.studentId !== id))
    storage.setEvaluations(storage.getEvaluations().filter(e => e.studentId !== id))
    storage.setContracts(storage.getContracts().filter(c => c.studentId !== id))

    return true
  }
}

export const studentMockService = new StudentMockService()
