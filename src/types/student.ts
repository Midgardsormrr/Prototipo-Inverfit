export type FrequencyType = '1 vez/sem.' | '2 veces/sem.' | '3 veces/sem.' | '4 veces/sem.'

export type StudentStatus = 'Activo' | 'Por revisar' | 'Inactivo'

export interface Student {
  id: string
  fullName: string
  email: string
  phone: string
  planName: string
  frequency: FrequencyType
  membershipStatus: StudentStatus
  fixedSchedule: string
  membershipStartDate: string
  membershipDueDate: string
  notes?: string
  healthConditions?: string
  emergencyContact?: string
  createdAt: string
}

export interface CreateStudentDTO {
  fullName: string
  email?: string
  phone?: string
  planName?: string
  frequency: FrequencyType
  membershipStatus?: StudentStatus
  fixedSchedule?: string
  membershipStartDate?: string
  membershipDueDate?: string
  notes?: string
  healthConditions?: string
  emergencyContact?: string
}
