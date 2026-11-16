import type { Payment, RegisterPaymentDTO, PaymentFilter, FinancialSummary } from '@/types'
import { paymentMockService, PaymentMockService } from './payments.mock'
import { paymentApiService, PaymentApiService } from './payments.api'
import { formatCurrencyCLP, sanitizeCsvCell, getTodayDateString } from './utils'

export interface IPaymentService {
  getAll(filter?: PaymentFilter): Promise<Payment[]>
  getById(id: string): Promise<Payment | null>
  getByStudentId(studentId: string): Promise<Payment[]>
  registerPayment(dto: RegisterPaymentDTO): Promise<Payment | null>
  updatePayment(id: string, updates: Partial<Payment>): Promise<Payment | null>
  archivePayment(id: string): Promise<boolean>
  getFinancialSummary(): Promise<FinancialSummary>
  sendPaymentLinkSimulation(studentId: string, paymentId?: string): Promise<{ success: boolean; link: string; message: string }>
  sendReminderSimulation(paymentId: string): Promise<{ success: boolean; message: string }>
  sendReceiptSimulation(paymentId: string): Promise<{ success: boolean; message: string }>
  exportToCSV(filter?: PaymentFilter): string | Promise<string>
}

/**
 * Determina si la aplicación utiliza la API real Go + MongoDB o el almacenamiento local (mock) para pagos.
 * Valor predeterminado: 'mock' (seguro para Cloudflare Pages y demo offline).
 */
export function getPaymentsDataSource(): 'mock' | 'api' {
  const source = import.meta.env.VITE_PAYMENTS_DATA_SOURCE
  return source === 'api' ? 'api' : 'mock'
}

export function isPaymentsApiMode(): boolean {
  return getPaymentsDataSource() === 'api'
}

class PaymentServiceDelegator implements IPaymentService {
  private get activeService(): IPaymentService {
    return isPaymentsApiMode() ? paymentApiService : paymentMockService
  }

  getAll(filter?: PaymentFilter): Promise<Payment[]> {
    return this.activeService.getAll(filter)
  }

  getById(id: string): Promise<Payment | null> {
    return this.activeService.getById(id)
  }

  getByStudentId(studentId: string): Promise<Payment[]> {
    return this.activeService.getByStudentId(studentId)
  }

  registerPayment(dto: RegisterPaymentDTO): Promise<Payment | null> {
    return this.activeService.registerPayment(dto)
  }

  updatePayment(id: string, updates: Partial<Payment>): Promise<Payment | null> {
    return this.activeService.updatePayment(id, updates)
  }

  archivePayment(id: string): Promise<boolean> {
    return this.activeService.archivePayment(id)
  }

  getFinancialSummary(): Promise<FinancialSummary> {
    return this.activeService.getFinancialSummary()
  }

  sendPaymentLinkSimulation(studentId: string, paymentId?: string): Promise<{ success: boolean; link: string; message: string }> {
    return this.activeService.sendPaymentLinkSimulation(studentId, paymentId)
  }

  sendReminderSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    return this.activeService.sendReminderSimulation(paymentId)
  }

  sendReceiptSimulation(paymentId: string): Promise<{ success: boolean; message: string }> {
    return this.activeService.sendReceiptSimulation(paymentId)
  }

  exportToCSV(filter?: PaymentFilter): string | Promise<string> {
    return this.activeService.exportToCSV(filter)
  }
}

export const paymentService = new PaymentServiceDelegator()
export { PaymentMockService, PaymentApiService, formatCurrencyCLP, sanitizeCsvCell, getTodayDateString }
