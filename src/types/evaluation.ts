export type EvaluationType = 'Evaluación inicial' | 'Seguimiento'
export type EvaluationStatus = 'Registrada' | 'Pendiente'

export interface EvaluationMetric {
  name: string
  value: string
  unit?: string
}

export interface Evaluation {
  id: string
  studentId: string
  studentName: string
  date: string
  type: EvaluationType
  status: EvaluationStatus
  metrics: EvaluationMetric[]
  observations: string
  conductedBy?: string
  createdAt: string
}

export interface CreateEvaluationDTO {
  studentId: string
  date: string
  type: EvaluationType
  status?: EvaluationStatus
  metrics: EvaluationMetric[]
  observations: string
}
