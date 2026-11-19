import { httpClient } from '../api/httpClient'
import type { Payment, RegisterPaymentDTO, PaymentFilter, FinancialSummary } from '@/types'
import type { IPaymentService } from './index'

export class PaymentApiService implements IPaymentService {
  async getAll(filter?: PaymentFilter): Promise<Payment[]> {
    const params = new URLSearchParams()
    if (filter) {
      if (filter.studentId) params.append('studentId', filter.studentId)
      if (filter.status) params.append('status', filter.status)
      if (filter.search) params.append('search', filter.search)
      if (filter.from) params.append('from', filter.from)
      if (filter.to) params.append('to', filter.to)
    }

    const query = params.toString()
    const endpoint = query ? `/payments?${query}` : '/payments'
    return httpClient.get<Payment[]>(endpoint)
  }

  async getByStudentId(studentId: string): Promise<Payment[]> {
    return httpClient.get<Payment[]>(`/payments/student/${encodeURIComponent(studentId)}`)
  }

  async getById(id: string): Promise<Payment | null> {
    try {
      return await httpClient.get<Payment>(`/payments/${encodeURIComponent(id)}`)
    } catch {
      return null
    }
  }

  async registerPayment(dto: RegisterPaymentDTO): Promise<Payment | null> {
    return httpClient.post<Payment>('/payments', {
      studentId: dto.studentId,
      amount: dto.amount,
      dueDate: dto.dueDate,
      paymentDate: dto.paymentDate,
      paymentMethod: dto.paymentMethod,
      transactionReference: dto.transactionReference,
      voucherUrl: dto.voucherUrl,
      voucherFileName: dto.voucherFileName,
      notes: dto.notes
    })
  }

  async updatePayment(id: string, updates: Partial<Payment>): Promise<Payment | null> {
    return httpClient.put<Payment>(`/payments/${encodeURIComponent(id)}`, updates)
  }

  async archivePayment(id: string): Promise<boolean> {
    try {
      await httpClient.delete(`/payments/${encodeURIComponent(id)}`)
      return true
    } catch {
      return false
    }
  }

  async getFinancialSummary(): Promise<FinancialSummary> {
    return httpClient.get<FinancialSummary>('/payments/summary')
  }

  async sendPaymentLinkSimulation(_studentId: string, paymentId?: string): Promise<{ success: boolean; link: string; message: string }> {
    const token = Math.random().toString(36).substring(2, 10).toUpperCase()
    const link = `https://demo.inverfit.invalid/checkout/${token}`

    if (paymentId) {
      try {
        await httpClient.post<Payment>(`/payments/${encodeURIComponent(paymentId)}/simulate-reminder`)
      } catch {
        // Fallback simulation note
      }
    }

    return {
      success: true,
      link,
      message: 'Enlace de pago generado exitosamente (modo API).'
    }
  }

  async sendReminderSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    const res = await httpClient.post<Payment>(`/payments/${encodeURIComponent(paymentId)}/simulate-reminder`)
    return {
      success: true,
      message: `Recordatorio registrado exitosamente para ${res.studentName}.`
    }
  }

  async sendReceiptSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    const res = await httpClient.post<Payment>(`/payments/${encodeURIComponent(paymentId)}/simulate-receipt`)
    return {
      success: true,
      message: `Comprobante registrado exitosamente para ${res.studentName}.`
    }
  }

  async exportToCSV(filter?: PaymentFilter): Promise<string> {
    const baseUrl = (import.meta.env.VITE_API_BASE_URL as string) || 'http://localhost:8080/api/v1'
    const cleanBaseUrl = baseUrl.replace(/\/+$/, '')
    const params = new URLSearchParams()
    if (filter) {
      if (filter.studentId) params.append('studentId', filter.studentId)
      if (filter.status) params.append('status', filter.status)
      if (filter.search) params.append('search', filter.search)
      if (filter.from) params.append('from', filter.from)
      if (filter.to) params.append('to', filter.to)
    }
    const query = params.toString()
    const url = query ? `${cleanBaseUrl}/payments/export?${query}` : `${cleanBaseUrl}/payments/export`

    const res = await fetch(url)
    if (!res.ok) {
      throw new Error(`Error al exportar CSV: ${res.statusText}`)
    }
    return res.text()
  }
}

export const paymentApiService = new PaymentApiService()
