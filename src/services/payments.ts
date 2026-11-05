import { storage } from './storage'
import type { Payment, RegisterPaymentDTO } from '@/types'

export class PaymentService {
  async getAll(): Promise<Payment[]> {
    return storage.getPayments()
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
    const dateFormatted = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

    let updatedPayment: Payment

    if (pendingIdx !== -1) {
      updatedPayment = {
        ...list[pendingIdx],
        status: 'Pagado',
        amount: dto.amount ?? list[pendingIdx].amount,
        paymentDate: dto.paymentDate || dateFormatted,
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
        amount: dto.amount || null,
        status: 'Pagado',
        paymentDate: dto.paymentDate || dateFormatted,
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

  async sendPaymentLinkSimulation(studentId: string, paymentId?: string): Promise<{ success: boolean; link: string; message: string }> {
    const student = storage.getStudents().find(s => s.id === studentId)
    if (!student) {
      throw new Error('Alumno no encontrado')
    }

    const token = Math.random().toString(36).substring(2, 10).toUpperCase()
    const link = `https://demo.inverfit.invalid/checkout/${token}`

    // Update payment reminder note if payment exists
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
    const headers = ['ID', 'Alumno', 'Plan', 'Frecuencia', 'Vencimiento', 'Monto', 'Estado', 'Medio de Pago', 'Fecha Pago', 'Comprobante Enviado']
    const rows = payments.map(p => [
      p.id,
      `"${p.studentName}"`,
      `"${p.planName}"`,
      `"${p.frequency}"`,
      `"${p.dueDate}"`,
      p.amount !== null ? p.amount : 'Por definir',
      `"${p.status}"`,
      `"${p.paymentMethod}"`,
      `"${p.paymentDate || 'Pendiente'}"`,
      p.receiptSent ? 'Sí' : 'No'
    ])

    const csvContent = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    return csvContent
  }
}

export const paymentService = new PaymentService()
