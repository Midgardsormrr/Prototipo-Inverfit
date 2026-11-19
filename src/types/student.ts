export type FrequencyType = '1 vez/sem.' | '2 veces/sem.' | '3 veces/sem.' | '4 veces/sem.'

export type StudentStatus = 'Activo' | 'Por revisar' | 'Inactivo'

/**
 * Perfil completo de ingreso y salud de Pablo (INVERFIT).
 * Mapea las 20 preguntas del formulario oficial + consentimiento y firma digital táctil.
 */
export interface StudentOnboarding {
  // 1. Demográficos
  rut?: string
  birthDate?: string
  age?: number | string
  gender?: 'Femenino' | 'Masculino' | 'Otro'
  comuna?: string
  occupation?: string // 'Estudiante' | 'Trabajador' | 'Deportista' | 'Otros'
  hasChildren?: 'Si' | 'No'
  childrenCount?: string

  // 2. Perfil Deportivo y Objetivos
  primaryGoal?: 'Salud, funcionalidad y longevidad' | 'Estética' | 'Rendimiento deportivo' | 'Alto Rendimiento' | 'Otros'
  goalDetails?: string
  hasTrainedBefore?: 'Si' | 'No'
  trainingExperience?: '1-6 meses' | '6-12 meses' | 'Más de 1 año' | 'Sin experiencia previa'
  sportsPracticed?: string
  availableDaysWeekly?: '1-2 veces a la semana' | '2-4 veces a la semana' | '4-6 veces a la semana' | 'Todos los días'

  // 3. Ficha Clínica y Médica (PAR-Q+ INVERFIT)
  hadCovid?: 'Si' | 'No'
  covidSequelae?: string
  chronicDiseases?: string // Asma, hipertensión, diabetes, cardiovasculares, etc.
  surgeriesInjuries?: string // Operaciones, lesiones crónicas y fechas
  injuriesPreventExercise?: 'Si' | 'No'
  currentDifficulties?: string // Dificultades para comenzar a entrenar con normalidad
  frequentConsumption?: string // Bebidas alcohólicas / drogas / fumador

  // 4. Consentimiento y Firma Digital
  consentAccepted: boolean
  consentTimestamp?: string
  digitalSignature?: string // Data URL base64 de la firma táctil
  termsVersion?: string
}

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
  onboarding?: StudentOnboarding
  createdAt: string
  updatedAt?: string
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
  onboarding?: StudentOnboarding
}
