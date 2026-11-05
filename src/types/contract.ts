export type DocumentType = 'Términos y condiciones' | 'Consentimiento informado' | 'Contrato de servicios'
export type ContractStatus = 'Aceptado' | 'Pendiente'

export interface ContractConsent {
  id: string
  studentId: string
  studentName: string
  documentType: DocumentType
  status: ContractStatus
  date: string
  termsVersion: string
  signedCopyAvailable: boolean
  notes?: string
}

export interface TermsVersion {
  version: string
  title: string
  effectiveDate: string
  content: string
  isCurrent: boolean
}
