export type PaymentStatus = 'Pagado' | 'Pendiente' | 'Por definir'

export type PaymentMethod =
  | 'SumUp (Link de pago)'
  | 'SumUp (POS Tarjeta)'
  | 'Transferencia Bancaria'
  | 'Efectivo'
  | 'Webpay'
  | 'Crédito'
  | 'Débito'
  | 'Servipag'
  | 'Por definir'

export interface Payment {
  id: string
  studentId: string
  studentName: string
  planName: string
  frequency: string
  dueDate: string
  amount: number | null
  status: PaymentStatus
  paymentDate: string | null
  paymentMethod: PaymentMethod
  transactionReference?: string
  voucherUrl?: string
  voucherFileName?: string
  reminderSentCount: number
  lastReminderAt: string | null
  receiptSent: boolean
  receiptSentAt: string | null
  notes?: string
}

export interface RegisterPaymentDTO {
  studentId: string
  amount?: number | null
  dueDate?: string
  paymentDate?: string
  paymentMethod: PaymentMethod
  transactionReference?: string
  voucherUrl?: string
  voucherFileName?: string
  notes?: string
}

export interface FinancialSummary {
  totalPaidAmount: number
  totalPendingAmount: number
  paidCount: number
  pendingCount: number
  totalCount: number
}

export interface PaymentFilter {
  studentId?: string
  status?: string
  search?: string
  from?: string
  to?: string
}
