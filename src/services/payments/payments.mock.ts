import { storage } from '../storage'
import type { Payment, RegisterPaymentDTO, PaymentFilter, FinancialSummary } from '@/types'
import { sanitizeCsvCell, getTodayDateString } from './utils'
import type { IPaymentService } from './index'

export class PaymentMockService implements IPaymentService {
  async getAll(filter?: PaymentFilter): Promise<Payment[]> {
    let list = storage.getPayments()
    if (filter) {
      if (filter.studentId) {
        list = list.filter(p => p.studentId === filter.studentId)
      }
      if (filter.status) {
        list = list.filter(p => p.status === filter.status)
      }
      if (filter.search) {
        const q = filter.search.toLowerCase().trim()
        list = list.filter(p => p.studentName.toLowerCase().includes(q))
      }
      if (filter.from) {
        list = list.filter(p => !p.paymentDate || p.paymentDate >= filter.from!)
      }
      if (filter.to) {
        list = list.filter(p => !p.paymentDate || p.paymentDate <= filter.to!)
      }
    }
    return list
  }

  async getByStudentId(studentId: string): Promise<Payment[]> {
    const list = storage.getPayments()
    return list.filter(p => p.studentId === studentId)
  }

  async getById(id: string): Promise<Payment | null> {
    const list = storage.getPayments()
    return list.find(p => p.id === id) || null
  }

  async registerPayment(dto: RegisterPaymentDTO): Promise<Payment | null> {
    const list = storage.getPayments()
    const student = storage.getStudents().find(s => s.id === dto.studentId)
    if (!student) return null

    // Look for existing pending payment for this student, or create a new payment
    const pendingIdx = list.findIndex(p => p.studentId === dto.studentId && p.status === 'Pendiente')
    const now = new Date()
    const dateFormatted = dto.paymentDate || getTodayDateString()

    let updatedPayment: Payment

    if (pendingIdx !== -1) {
      updatedPayment = {
        ...list[pendingIdx],
        status: 'Pagado',
        amount: dto.amount !== undefined ? dto.amount : list[pendingIdx].amount,
        paymentDate: dateFormatted,
        paymentMethod: dto.paymentMethod,
        transactionReference: `BOL-${now.getFullYear()}${String(now.getMonth() + 1).padStart(2, '0')}-${Date.now().toString().slice(-4)}`,
        notes: dto.notes || 'Pago confirmado y registrado por el entrenador.'
      }
      list[pendingIdx] = updatedPayment
    } else {
      updatedPayment = {
        id: `pay-${Date.now().toString().slice(-6)}`,
        studentId: student.id,
        studentName: student.fullName,
        planName: student.planName,
        frequency: student.frequency,
        dueDate: dto.dueDate || 'Por definir',
        amount: dto.amount ?? null,
        status: 'Pagado',
        paymentDate: dateFormatted,
        paymentMethod: dto.paymentMethod,
        transactionReference: `BOL-${now.getFullYear()}${String(now.getMonth() + 1).padStart(2, '0')}-${Date.now().toString().slice(-4)}`,
        reminderSentCount: 0,
        lastReminderAt: null,
        receiptSent: false,
        receiptSentAt: null,
        notes: dto.notes || 'Nuevo pago registrado manualmente.'
      }
      list.unshift(updatedPayment)
    }

    storage.setPayments([...list])

    // Update student membership status if it was "Por revisar"
    if (student.membershipStatus === 'Por revisar') {
      const students = storage.getStudents().map(s =>
        s.id === student.id ? { ...s, membershipStatus: 'Activo' as const } : s
      )
      storage.setStudents(students)
    }

    return updatedPayment
  }

  async updatePayment(id: string, updates: Partial<Payment>): Promise<Payment | null> {
    const list = storage.getPayments()
    const idx = list.findIndex(p => p.id === id)
    if (idx === -1) return null

    const updated: Payment = {
      ...list[idx],
      ...updates
    }
    list[idx] = updated
    storage.setPayments([...list])
    return updated
  }

  async archivePayment(id: string): Promise<boolean> {
    const list = storage.getPayments()
    const filtered = list.filter(p => p.id !== id)
    if (filtered.length === list.length) return false
    storage.setPayments(filtered)
    return true
  }

  async getFinancialSummary(): Promise<FinancialSummary> {
    const payments = storage.getPayments()
    let totalPaid = 0
    let totalPending = 0
    let paidCount = 0
    let pendingCount = 0

    for (const p of payments) {
      if (p.status === 'Pagado') {
        paidCount++
        if (p.amount) totalPaid += p.amount
      } else if (p.status === 'Pendiente') {
        pendingCount++
        if (p.amount) totalPending += p.amount
      }
    }

    return {
      totalPaidAmount: totalPaid,
      totalPendingAmount: totalPending,
      paidCount,
      pendingCount,
      totalCount: payments.length
    }
  }

  async sendPaymentLinkSimulation(studentId: string, paymentId?: string): Promise<{ success: boolean; link: string; message: string }> {
    const student = storage.getStudents().find(s => s.id === studentId)
    if (!student) {
      throw new Error('Alumno no encontrado')
    }

    const token = Math.random().toString(36).substring(2, 10).toUpperCase()
    const link = `https://demo.inverfit.invalid/checkout/${token}`

    const list = storage.getPayments()
    let pIdx = -1
    if (paymentId) {
      pIdx = list.findIndex(p => p.id === paymentId)
    } else {
      pIdx = list.findIndex(p => p.studentId === studentId && p.status === 'Pendiente')
    }

    const now = new Date()
    const timestamp = `${now.toLocaleDateString('es-CL')} ${now.toLocaleTimeString('es-CL', { hour: '2-digit', minute: '2-digit' })}`

    if (pIdx !== -1) {
      list[pIdx] = {
        ...list[pIdx],
        reminderSentCount: (list[pIdx].reminderSentCount || 0) + 1,
        lastReminderAt: timestamp
      }
      storage.setPayments([...list])
    }

    return {
      success: true,
      link,
      message: `Enlace de pago simulado generado exitosamente para ${student.fullName}.`
    }
  }

  async sendReminderSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    const list = storage.getPayments()
    const idx = list.findIndex(p => p.id === paymentId)
    if (idx === -1) throw new Error('Pago no encontrado')

    const now = new Date()
    const timestamp = `${now.toLocaleDateString('es-CL')} ${now.toLocaleTimeString('es-CL', { hour: '2-digit', minute: '2-digit' })}`

    list[idx] = {
      ...list[idx],
      reminderSentCount: (list[idx].reminderSentCount || 0) + 1,
      lastReminderAt: timestamp
    }
    storage.setPayments([...list])

    return {
      success: true,
      message: `Recordatorio simulado enviado exitosamente a ${list[idx].studentName}.`
    }
  }

  async sendReceiptSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    const list = storage.getPayments()
    const idx = list.findIndex(p => p.id === paymentId)
    if (idx === -1) throw new Error('Pago no encontrado')

    const now = new Date()
    const timestamp = `${now.toLocaleDateString('es-CL')} ${now.toLocaleTimeString('es-CL', { hour: '2-digit', minute: '2-digit' })}`

    list[idx] = {
      ...list[idx],
      receiptSent: true,
      receiptSentAt: timestamp
    }
    storage.setPayments([...list])

    return {
      success: true,
      message: `Comprobante de pago simulado enviado a ${list[idx].studentName}.`
    }
  }

  exportToCSV(): string {
    const payments = storage.getPayments()
    const headers = ['ID', 'Alumno', 'Plan', 'Frecuencia', 'Vencimiento', 'Monto', 'Estado', 'Medio de Pago', 'Fecha Pago', 'Comprobante Enviado', 'Notas']
    const rows = payments.map(p => [
      `"${sanitizeCsvCell(p.id)}"`,
      `"${sanitizeCsvCell(p.studentName)}"`,
      `"${sanitizeCsvCell(p.planName)}"`,
      `"${sanitizeCsvCell(p.frequency)}"`,
      `"${sanitizeCsvCell(p.dueDate)}"`,
      p.amount !== null ? p.amount : 'Por definir',
      `"${sanitizeCsvCell(p.status)}"`,
      `"${sanitizeCsvCell(p.paymentMethod)}"`,
      `"${sanitizeCsvCell(p.paymentDate || 'Pendiente')}"`,
      p.receiptSent ? 'Sí' : 'No',
      `"${sanitizeCsvCell(p.notes || '')}"`
    ])

    return [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
  }
}

export const paymentMockService = new PaymentMockService()
