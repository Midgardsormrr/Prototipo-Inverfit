import { storage } from './storage'
import type { ContractConsent, ContractStatus, TermsVersion } from '@/types'

export class ContractService {
  async getAll(): Promise<ContractConsent[]> {
    return storage.getContracts()
  }

  async getByStudentId(studentId: string): Promise<ContractConsent[]> {
    const list = storage.getContracts()
    return list.filter(c => c.studentId === studentId)
  }

  async getById(id: string): Promise<ContractConsent | null> {
    const list = storage.getContracts()
    return list.find(c => c.id === id) || null
  }

  async updateStatus(id: string, status: ContractStatus): Promise<ContractConsent | null> {
    const list = storage.getContracts()
    const idx = list.findIndex(c => c.id === id)
    if (idx === -1) return null

    const now = new Date()
    const dateStr = status === 'Aceptado'
      ? `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
      : 'Por registrar'

    list[idx] = {
      ...list[idx],
      status,
      date: dateStr,
      signedCopyAvailable: status === 'Aceptado'
    }

    storage.setContracts([...list])
    return list[idx]
  }

  async getTerms(): Promise<TermsVersion[]> {
    return storage.getTerms()
  }

  async getCurrentTerms(): Promise<TermsVersion | null> {
    const list = storage.getTerms()
    return list.find(t => t.isCurrent) || list[0] || null
  }

  async updateCurrentTerms(version: string, title: string, content: string): Promise<TermsVersion> {
    const list = storage.getTerms()
    const now = new Date()
    const monthYear = `${now.toLocaleString('es-CL', { month: 'long' })} ${now.getFullYear()}`

    // Set existing current to false
    const updatedList = list.map(t => ({ ...t, isCurrent: false }))

    const newVersion: TermsVersion = {
      version,
      title,
      effectiveDate: monthYear,
      content,
      isCurrent: true
    }

    storage.setTerms([newVersion, ...updatedList])
    return newVersion
  }
}

export const contractService = new ContractService()
