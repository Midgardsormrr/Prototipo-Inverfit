import {
  initialStudents,
  initialPayments,
  initialScheduleBlocks,
  initialCoachAvailability,
  initialScheduleExceptions,
  initialEvaluations,
  initialContracts,
  initialTermsVersions
} from './mockData'
import type {
  Student,
  Payment,
  ScheduleBlock,
  DayAvailability,
  ScheduleException,
  Evaluation,
  ContractConsent,
  TermsVersion
} from '@/types'

const STORAGE_KEYS = {
  STUDENTS: 'inverfit_students_v1',
  PAYMENTS: 'inverfit_payments_v1',
  SCHEDULE: 'inverfit_schedule_v1',
  AVAILABILITY: 'inverfit_availability_v1',
  EXCEPTIONS: 'inverfit_schedule_exceptions_v1',
  EVALUATIONS: 'inverfit_evaluations_v1',
  CONTRACTS: 'inverfit_contracts_v1',
  TERMS: 'inverfit_terms_v1'
}

class LocalStorageService {
  private memoryStore: Record<string, string> = {}
  private listeners: Set<() => void> = new Set()

  private hasStorage(): boolean {
    return typeof localStorage !== 'undefined'
  }

  private getItem(key: string): string | null {
    if (this.hasStorage()) {
      return localStorage.getItem(key)
    }
    return this.memoryStore[key] || null
  }

  private setItem(key: string, value: string): void {
    if (this.hasStorage()) {
      localStorage.setItem(key, value)
    } else {
      this.memoryStore[key] = value
    }
  }

  private removeItem(key: string): void {
    if (this.hasStorage()) {
      localStorage.removeItem(key)
    } else {
      delete this.memoryStore[key]
    }
  }

  public subscribe(cb: () => void): () => void {
    this.listeners.add(cb)
    return () => this.listeners.delete(cb)
  }

  private notify() {
    this.listeners.forEach(cb => cb())
  }

  // Migrate any previous local data to preserve user state
  private getMigratedOrInit<T>(key: string, initialDefault: T) {
    if (!this.getItem(key)) {
      if (this.hasStorage()) {
        const suffix = key.replace(/^inverfit/, '')
        for (let i = 0; i < localStorage.length; i++) {
          const existingKey = localStorage.key(i)
          if (existingKey && existingKey.endsWith(suffix) && existingKey !== key) {
            const legacyData = localStorage.getItem(existingKey)
            if (legacyData) {
              this.setItem(key, legacyData)
              this.removeItem(existingKey)
              return
            }
          }
        }
      }
      this.set(key, initialDefault)
    }
  }

  public cleanupLegacyDemoData(): void {
    if (!this.hasStorage()) return
    const rawContracts = this.getItem(STORAGE_KEYS.CONTRACTS)
    const rawStudents = this.getItem(STORAGE_KEYS.STUDENTS)
    const rawEvals = this.getItem(STORAGE_KEYS.EVALUATIONS)
    if (
      (rawContracts && (rawContracts.includes('Camila Morales') || rawContracts.includes('Valeria Soto') || rawContracts.includes('Carlos Valenzuela') || rawContracts.includes('Nelson Test Onboarding') || rawContracts.includes('cnt-0'))) ||
      (rawStudents && (rawStudents.includes('alu-0') || rawStudents.includes('Camila Morales') || rawStudents.includes('Carlos Valenzuela') || rawStudents.includes('Nelson Test Onboarding'))) ||
      (rawEvals && (rawEvals.includes('Camila Morales') || rawEvals.includes('Carlos Valenzuela') || rawEvals.includes('Nelson Test Onboarding') || rawEvals.includes('eval-0') || rawEvals.includes('eval-nelson')))
    ) {
      this.set(STORAGE_KEYS.STUDENTS, [])
      this.set(STORAGE_KEYS.PAYMENTS, [])
      this.set(STORAGE_KEYS.SCHEDULE, [])
      this.set(STORAGE_KEYS.EXCEPTIONS, [])
      this.set(STORAGE_KEYS.EVALUATIONS, [])
      this.set(STORAGE_KEYS.CONTRACTS, [])
    }
  }

  public initData() {
    this.cleanupLegacyDemoData()
    this.getMigratedOrInit(STORAGE_KEYS.STUDENTS, initialStudents)
    this.getMigratedOrInit(STORAGE_KEYS.PAYMENTS, initialPayments)
    this.getMigratedOrInit(STORAGE_KEYS.SCHEDULE, initialScheduleBlocks)
    this.getMigratedOrInit(STORAGE_KEYS.AVAILABILITY, initialCoachAvailability)
    this.getMigratedOrInit(STORAGE_KEYS.EXCEPTIONS, initialScheduleExceptions)
    this.getMigratedOrInit(STORAGE_KEYS.EVALUATIONS, initialEvaluations)
    this.getMigratedOrInit(STORAGE_KEYS.CONTRACTS, initialContracts)
    this.getMigratedOrInit(STORAGE_KEYS.TERMS, initialTermsVersions)
  }

  public resetAll() {
    this.set(STORAGE_KEYS.STUDENTS, initialStudents)
    this.set(STORAGE_KEYS.PAYMENTS, initialPayments)
    this.set(STORAGE_KEYS.SCHEDULE, initialScheduleBlocks)
    this.set(STORAGE_KEYS.AVAILABILITY, initialCoachAvailability)
    this.set(STORAGE_KEYS.EXCEPTIONS, initialScheduleExceptions)
    this.set(STORAGE_KEYS.EVALUATIONS, initialEvaluations)
    this.set(STORAGE_KEYS.CONTRACTS, initialContracts)
    this.set(STORAGE_KEYS.TERMS, initialTermsVersions)
    this.notify()
  }

  public get<T>(key: string, defaultValue: T): T {
    try {
      const raw = this.getItem(key)
      if (!raw) return defaultValue
      return JSON.parse(raw) as T
    } catch (e) {
      console.error(`Error reading ${key} from storage:`, e)
      return defaultValue
    }
  }

  public set<T>(key: string, value: T): void {
    try {
      this.setItem(key, JSON.stringify(value))
      this.notify()
    } catch (e) {
      console.error(`Error writing ${key} to storage:`, e)
    }
  }

  // Getters for specific collections
  public getStudents(): Student[] {
    return this.get<Student[]>(STORAGE_KEYS.STUDENTS, initialStudents)
  }

  public setStudents(students: Student[]) {
    this.set(STORAGE_KEYS.STUDENTS, students)
  }

  public getPayments(): Payment[] {
    return this.get<Payment[]>(STORAGE_KEYS.PAYMENTS, initialPayments)
  }

  public setPayments(payments: Payment[]) {
    this.set(STORAGE_KEYS.PAYMENTS, payments)
  }

  public getSchedule(): ScheduleBlock[] {
    return this.get<ScheduleBlock[]>(STORAGE_KEYS.SCHEDULE, initialScheduleBlocks)
  }

  public setSchedule(blocks: ScheduleBlock[]) {
    this.set(STORAGE_KEYS.SCHEDULE, blocks)
  }

  public getAvailability(): DayAvailability[] {
    return this.get<DayAvailability[]>(STORAGE_KEYS.AVAILABILITY, initialCoachAvailability)
  }

  public setAvailability(availability: DayAvailability[]) {
    this.set(STORAGE_KEYS.AVAILABILITY, availability)
  }

  public getExceptions(): ScheduleException[] {
    return this.get<ScheduleException[]>(STORAGE_KEYS.EXCEPTIONS, initialScheduleExceptions)
  }

  public setExceptions(exceptions: ScheduleException[]) {
    this.set(STORAGE_KEYS.EXCEPTIONS, exceptions)
  }

  public getEvaluations(): Evaluation[] {
    const list = this.get<Evaluation[]>(STORAGE_KEYS.EVALUATIONS, initialEvaluations)
    return list.filter(e => 
      !e.studentName.includes('Camila Morales') && 
      !e.studentName.includes('Carlos Valenzuela') &&
      !e.studentName.includes('Nelson Test Onboarding') &&
      !e.id.startsWith('eval-0') &&
      !e.id.includes('nelson-official')
    )
  }

  public setEvaluations(evaluations: Evaluation[]) {
    this.set(STORAGE_KEYS.EVALUATIONS, evaluations)
  }

  public getContracts(): ContractConsent[] {
    const list = this.get<ContractConsent[]>(STORAGE_KEYS.CONTRACTS, initialContracts)
    return list.filter(c => 
      !c.studentName.toLowerCase().includes('alumno de ejemplo') &&
      !c.studentName.includes('Camila Morales') &&
      !c.studentName.includes('Felipe Arancibia') &&
      !c.studentName.includes('Valeria Soto') &&
      !c.studentName.includes('Carlos Valenzuela') &&
      !c.studentName.includes('Nelson Test Onboarding') &&
      !c.id.startsWith('cnt-0')
    )
  }

  public setContracts(contracts: ContractConsent[]) {
    this.set(STORAGE_KEYS.CONTRACTS, contracts)
  }

  public getTerms(): TermsVersion[] {
    return this.get<TermsVersion[]>(STORAGE_KEYS.TERMS, initialTermsVersions)
  }

  public setTerms(terms: TermsVersion[]) {
    this.set(STORAGE_KEYS.TERMS, terms)
  }
}

export const storage = new LocalStorageService()
storage.initData()
