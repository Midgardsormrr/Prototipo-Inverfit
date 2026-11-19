export type EvaluationType = 'Evaluación inicial' | 'Seguimiento'
export type EvaluationStatus = 'Registrada' | 'Pendiente'
export type EvaluationProtocol = 'clinico_deportivo' | 'adulto_mayor'

export interface EvaluationMetric {
  name: string
  value: string
  unit?: string
}

export interface VitalSignsData {
  spo2?: number
  spo2Alert?: string
  systolicBP?: number
  diastolicBP?: number
  bpClassification?: string
  restingHR?: number
  hrClassification?: string
}

export interface AnthropometryData {
  weightKg?: number
  heightCm?: number
  imc?: number
  imcClassification?: string
  fatPercentage?: number
  fatClassification?: string
  muscleMassKg?: number
}

export interface DynamometryData {
  rightHandKg?: number
  leftHandKg?: number
  asymmetryPercent?: number
  asymmetryClassification?: string
  strengthClassification?: string
}

export interface SeniorFitnessData {
  sitToStand5s?: number
  sitToStandAlert?: string // e.g. "Alto Riesgo de Caída (>13.7s)" vs "Normal"
  chairStandReps?: number
  armCurlReps?: number
  walk6MinMeters?: number
}

export interface VelocityEncoderData {
  exercise?: string
  loadKg?: number
  velocityMs?: number
  estimated1RM?: number
  trainingZone?: string // e.g. "Fuerza - Hipertrofia (0.50 - 0.75 m/s)"
}

export interface RomPatternsData {
  kneeToWallLeft?: string
  kneeToWallRight?: string
  aslrLeft?: string
  aslrRight?: string
  overheadSquat?: string
  lungeUnipodal?: string
  hipHinge?: string
  observations?: string
}

export interface IsometricStrengthData {
  imtpLeftN?: number
  imtpRightN?: number
  legExtensionLeftN?: number
  legExtensionRightN?: number
  legCurlLeftN?: number
  legCurlRightN?: number
  hipAbductionLeftN?: number
  hipAbductionRightN?: number
}

export interface BoscoJumpData {
  squatJumpCm?: number
  squatJumpFlightMs?: number
  cmjCm?: number
  cmjFlightMs?: number
  abalakovCm?: number
  abalakovFlightMs?: number
  dropJumpCm?: number
  dropJumpContactMs?: number
}

export interface Evaluation {
  id: string
  studentId: string
  studentName: string
  date: string
  type: EvaluationType
  status: EvaluationStatus
  protocol?: EvaluationProtocol
  vitalSigns?: VitalSignsData
  anthropometry?: AnthropometryData
  dynamometry?: DynamometryData
  romPatterns?: RomPatternsData
  isometricStrength?: IsometricStrengthData
  boscoJump?: BoscoJumpData
  seniorFitness?: SeniorFitnessData
  velocityEncoder?: VelocityEncoderData
  clinicalAlert?: string
  metrics: EvaluationMetric[]
  observations: string
  conductedBy?: string
  createdAt: string
  workbookData?: Record<string, any>
}

export interface CreateEvaluationDTO {
  studentId: string
  studentName?: string
  date: string
  type: EvaluationType
  status?: EvaluationStatus
  protocol?: EvaluationProtocol
  vitalSigns?: VitalSignsData
  anthropometry?: AnthropometryData
  dynamometry?: DynamometryData
  romPatterns?: RomPatternsData
  isometricStrength?: IsometricStrengthData
  boscoJump?: BoscoJumpData
  seniorFitness?: SeniorFitnessData
  velocityEncoder?: VelocityEncoderData
  clinicalAlert?: string
  metrics: EvaluationMetric[]
  observations: string
  conductedBy?: string
  workbookData?: Record<string, any>
}
