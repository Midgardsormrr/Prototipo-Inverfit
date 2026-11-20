import type { Student, Payment, ScheduleBlock, Evaluation, ContractConsent, TermsVersion, DayAvailability, ScheduleException } from '@/types'

// Colecciones iniciales limpias para entorno productivo
export const initialStudents: Student[] = []

export const initialPayments: Payment[] = []

export const initialScheduleBlocks: ScheduleBlock[] = []

export const initialScheduleExceptions: ScheduleException[] = []

export const initialEvaluations: Evaluation[] = []

export const initialContracts: ContractConsent[] = []

// Disponibilidad horaria preestablecida oficial de Pablo (Lunes a Sábado, 75 min)
export const initialCoachAvailability: DayAvailability[] = [
  {
    day: 'Lun',
    enabled: true,
    intervals: [
      { id: 'int-lun-1', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Mar',
    enabled: true,
    intervals: [
      { id: 'int-mar-1', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Mié',
    enabled: true,
    intervals: [
      { id: 'int-mie-1', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Jue',
    enabled: true,
    intervals: [
      { id: 'int-jue-1', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Vie',
    enabled: true,
    intervals: [
      { id: 'int-vie-1', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Sáb',
    enabled: true,
    intervals: [
      { id: 'int-sab-eval', startTime: '06:30', endTime: '22:00' }
    ]
  },
  {
    day: 'Dom',
    enabled: false,
    intervals: []
  }
]

// Términos y condiciones y reglamento oficial INVERFIT
export const initialTermsVersions: TermsVersion[] = [
  {
    version: 'v1.0 (2026)',
    title: 'Términos y Condiciones del Servicio de Entrenamiento — INVERFIT',
    effectiveDate: 'Octubre 2026',
    isCurrent: true,
    content: `CONTRATO DE PRESTACIÓN DE SERVICIOS Y CONSENTIMIENTO INFORMADO

1. NATURALEZA DEL SERVICIO
INVERFIT ofrece sesiones de entrenamiento personalizado y grupal orientadas a la mejora de la condición física, guiadas directamente por el entrenador Pablo. Las sesiones tienen una duración planificada estándar de 75 minutos (adaptable según el plan acordado), organizadas según la disponibilidad horaria configurada por el entrenador.

2. ASIGNACIÓN Y PUNTUALIDAD DE HORARIOS
Los cupos se coordinan en bloques fijos semanales (1, 2, 3 o 4 veces por semana). La administración de la agenda recae exclusivamente en el entrenador. El alumno se compromete a asistir puntualmente a sus bloques fijados.

3. POLÍTICA DE AUSENCIAS Y REAGENDAMIENTO
Dado que los cupos son fijos y limitados, las sesiones perdidas por inasistencia no se acumulan para meses posteriores ni aplazan la fecha de término de la membresía. Todo cambio excepcional requiere aviso previo al entrenador y está estrictamente sujeto a disponibilidad en la misma semana operativa. De no existir bloque disponible, la clase se da por perdida.

4. CONDICIÓN FÍSICA Y SALUD
El alumno declara encontrarse en condiciones físicas aptas para realizar ejercicio físico de intensidad variable y se compromete a informar oportunamente al entrenador cualquier dolor, lesión o prescripción médica preexistente o sobreviniente.

5. PAGOS Y MEMBRESÍAS
Las membresías se pagan por adelantado de forma mensual mediante los canales autorizados (transferencia bancaria, enlace de pago simulado o cobro presencial). La vigencia de los planes es mensual e improrrogable, salvo acuerdo formal previo por causa justificada.

6. RESGUARDO DE ANTECEDENTES Y REGLAMENTO INTERNO
Este reglamento rige las condiciones de entrenamiento personalizado y consentimiento informado para todos los alumnos activos de INVERFIT.`
  }
]
