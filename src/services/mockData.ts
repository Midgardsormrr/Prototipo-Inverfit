import type { Student, Payment, ScheduleBlock, Evaluation, ContractConsent, TermsVersion, DayAvailability, ScheduleException } from '@/types'

export const initialStudents: Student[] = [
  {
    id: 'alu-01',
    fullName: 'Alumno de ejemplo 01',
    email: 'alumno01@ejemplo.com',
    phone: '+56 9 1234 5671',
    planName: 'Plan estándar',
    frequency: '2 veces/sem.',
    membershipStatus: 'Por revisar',
    fixedSchedule: 'Lun y Mié 06:30',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Requiere revisión de horario y fecha de inicio de membresía.',
    healthConditions: 'Sin lesiones reportadas. Apto para entrenamiento.',
    emergencyContact: 'Contacto familiar (+56 9 9999 1111)',
    createdAt: '2026-10-01'
  },
  {
    id: 'alu-02',
    fullName: 'Alumno de ejemplo 02',
    email: 'alumno02@ejemplo.com',
    phone: '+56 9 1234 5672',
    planName: 'Plan estándar',
    frequency: '1 vez/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Sáb 08:00',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Asiste a bloque regular de sábado por la mañana.',
    healthConditions: 'Molestia leve en rodilla izquierda (precaución en sentadillas profundas).',
    emergencyContact: 'Contacto familiar (+56 9 9999 2222)',
    createdAt: '2026-10-05'
  },
  {
    id: 'alu-03',
    fullName: 'Alumno de ejemplo 03',
    email: 'alumno03@ejemplo.com',
    phone: '+56 9 1234 5673',
    planName: 'Plan estándar',
    frequency: '3 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Lun, Mié y Vie 09:30',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Objetivo acondicionamiento general y fuerza.',
    healthConditions: 'Ninguna contraindicación.',
    emergencyContact: 'Contacto familiar (+56 9 9999 3333)',
    createdAt: '2026-10-08'
  },
  {
    id: 'alu-04',
    fullName: 'Alumno de ejemplo 04',
    email: 'alumno04@ejemplo.com',
    phone: '+56 9 1234 5674',
    planName: 'Plan estándar',
    frequency: '4 veces/sem.',
    membershipStatus: 'Por revisar',
    fixedSchedule: 'Lun, Mar, Jue y Sáb 11:00',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Plan intensivo. Mantener monitoreo de carga.',
    healthConditions: 'Ninguna observada.',
    emergencyContact: 'Contacto familiar (+56 9 9999 4444)',
    createdAt: '2026-10-12'
  },
  {
    id: 'alu-05',
    fullName: 'Alumno de ejemplo 05',
    email: 'alumno05@ejemplo.com',
    phone: '+56 9 1234 5675',
    planName: 'Plan estándar',
    frequency: '2 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Mar y Jue 08:00',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Buen compromiso y regularidad.',
    healthConditions: 'Apto para toda actividad física.',
    emergencyContact: 'Contacto familiar (+56 9 9999 5555)',
    createdAt: '2026-10-15'
  },
  {
    id: 'alu-06',
    fullName: 'Alumno de ejemplo 06',
    email: 'alumno06@ejemplo.com',
    phone: '+56 9 1234 5676',
    planName: 'Plan estándar',
    frequency: '1 vez/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Mié 12:30',
    membershipStartDate: 'Por registrar',
    membershipDueDate: 'Por definir',
    notes: 'Sesión complementaria semanal.',
    healthConditions: 'Ninguna.',
    emergencyContact: 'Contacto familiar (+56 9 9999 6666)',
    createdAt: '2026-10-20'
  }
]

export const initialPayments: Payment[] = [
  {
    id: 'pay-01',
    studentId: 'alu-01',
    studentName: 'Alumno de ejemplo 01',
    planName: 'Plan estándar',
    frequency: '2 veces/sem.',
    dueDate: 'Por definir',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'Por definir',
    reminderSentCount: 1,
    lastReminderAt: '2026-11-02 10:15',
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Enlace enviado, esperando confirmación de transferencia o pago.'
  },
  {
    id: 'pay-02',
    studentId: 'alu-02',
    studentName: 'Alumno de ejemplo 02',
    planName: 'Plan estándar',
    frequency: '1 vez/sem.',
    dueDate: 'Por definir',
    amount: null,
    status: 'Pagado',
    paymentDate: 'Por registrar',
    paymentMethod: 'Transferencia',
    transactionReference: 'TRF-882914',
    reminderSentCount: 0,
    lastReminderAt: null,
    receiptSent: true,
    receiptSentAt: '2026-10-29 14:20',
    notes: 'Pago registrado por transferencia electrónica.'
  },
  {
    id: 'pay-03',
    studentId: 'alu-03',
    studentName: 'Alumno de ejemplo 03',
    planName: 'Plan estándar',
    frequency: '3 veces/sem.',
    dueDate: 'Por definir',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'Por definir',
    reminderSentCount: 0,
    lastReminderAt: null,
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Pendiente de emisión de enlace de pago.'
  },
  {
    id: 'pay-04',
    studentId: 'alu-04',
    studentName: 'Alumno de ejemplo 04',
    planName: 'Plan estándar',
    frequency: '4 veces/sem.',
    dueDate: 'Por definir',
    amount: null,
    status: 'Pagado',
    paymentDate: 'Por registrar',
    paymentMethod: 'Débito',
    transactionReference: 'SUM-102941',
    reminderSentCount: 1,
    lastReminderAt: '2026-10-25 09:00',
    receiptSent: true,
    receiptSentAt: '2026-10-28 11:30',
    notes: 'Pago realizado vía tarjeta presencial SumUp.'
  },
  {
    id: 'pay-05',
    studentId: 'alu-05',
    studentName: 'Alumno de ejemplo 05',
    planName: 'Plan estándar',
    frequency: '2 veces/sem.',
    dueDate: 'Por definir',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'Por definir',
    reminderSentCount: 2,
    lastReminderAt: '2026-11-03 16:40',
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Recordatorio enviado vía WhatsApp / link simulado.'
  }
]

export const initialScheduleBlocks: ScheduleBlock[] = [
  // 06:30 (duración 75 min: 06:30 a 07:45, descanso hasta 08:00)
  {
    id: 'blk-lun-0630',
    day: 'Lun',
    startTime: '06:30',
    endTime: '07:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-01'],
    notes: 'Bloque matutino inicial. Grupo de fuerza técnica.'
  },
  {
    id: 'blk-mie-0630',
    day: 'Mié',
    startTime: '06:30',
    endTime: '07:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-01'],
    notes: 'Seguimiento de sesión técnica.'
  },
  {
    id: 'blk-jue-0630',
    day: 'Jue',
    startTime: '06:30',
    endTime: '07:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Sesión grupal programada.'
  },
  {
    id: 'blk-sab-0630',
    day: 'Sáb',
    startTime: '06:30',
    endTime: '07:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Bloque especial sábado madrugada.'
  },

  // 08:00 (duración 75 min: 08:00 a 09:15, descanso hasta 09:30)
  {
    id: 'blk-lun-0800',
    day: 'Lun',
    startTime: '08:00',
    endTime: '09:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Bloque con alumnos confirmados.'
  },
  {
    id: 'blk-mar-0800',
    day: 'Mar',
    startTime: '08:00',
    endTime: '09:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-05'],
    notes: 'Grupo enfocado en acondicionamiento metabólico.'
  },
  {
    id: 'blk-jue-0800',
    day: 'Jue',
    startTime: '08:00',
    endTime: '09:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-05'],
    notes: 'Continuación sesión martes.'
  },
  {
    id: 'blk-sab-0800',
    day: 'Sáb',
    startTime: '08:00',
    endTime: '09:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-02'],
    notes: 'Grupo sábado mañana.'
  },

  // 09:30 (duración 75 min: 09:30 a 10:45, descanso hasta 11:00)
  {
    id: 'blk-lun-0930',
    day: 'Lun',
    startTime: '09:30',
    endTime: '10:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-03'],
    notes: 'Bloque lunes media mañana.'
  },
  {
    id: 'blk-mie-0930',
    day: 'Mié',
    startTime: '09:30',
    endTime: '10:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-03'],
    notes: 'Media mañana miércoles.'
  },
  {
    id: 'blk-jue-0930',
    day: 'Jue',
    startTime: '09:30',
    endTime: '10:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Sesión grupal abierta.'
  },
  {
    id: 'blk-sab-0930',
    day: 'Sáb',
    startTime: '09:30',
    endTime: '10:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Sesión grupal sábado.'
  },

  // 11:00 (duración 75 min: 11:00 a 12:15, descanso hasta 12:30)
  {
    id: 'blk-lun-1100',
    day: 'Lun',
    startTime: '11:00',
    endTime: '12:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-04'],
    notes: 'Fuerza hipertrofia grupo mediodía.'
  },
  {
    id: 'blk-mie-1100',
    day: 'Mié',
    startTime: '11:00',
    endTime: '12:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Sesión grupal.'
  },
  {
    id: 'blk-jue-1100',
    day: 'Jue',
    startTime: '11:00',
    endTime: '12:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-04'],
    notes: 'Seguimiento grupo 11:00.'
  },
  {
    id: 'blk-sab-1100',
    day: 'Sáb',
    startTime: '11:00',
    endTime: '12:15',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-04'],
    notes: 'Cierre bloque sábado mediodía.'
  },

  // 12:30 (duración 75 min: 12:30 a 13:45, descanso hasta 14:00)
  {
    id: 'blk-mar-1230',
    day: 'Mar',
    startTime: '12:30',
    endTime: '13:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-04'],
    notes: 'Sesión post-mañana.'
  },
  {
    id: 'blk-mie-1230',
    day: 'Mié',
    startTime: '12:30',
    endTime: '13:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-06'],
    notes: 'Sesión grupal miércoles tarde.'
  },
  {
    id: 'blk-vie-1230',
    day: 'Vie',
    startTime: '12:30',
    endTime: '13:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: ['alu-03'],
    notes: 'Sesión viernes cierre semana.'
  },
  {
    id: 'blk-sab-1230',
    day: 'Sáb',
    startTime: '12:30',
    endTime: '13:45',
    durationMinutes: 75,
    type: 'Sesión grupal',
    studentIds: [],
    notes: 'Sesión sábado mediodía.'
  }
]

export const initialCoachAvailability: DayAvailability[] = [
  {
    day: 'Lun',
    enabled: true,
    intervals: [
      { id: 'int-lun-1', startTime: '06:30', endTime: '13:30' }
    ]
  },
  {
    day: 'Mar',
    enabled: true,
    intervals: [
      { id: 'int-mar-1', startTime: '08:00', endTime: '12:00' },
      { id: 'int-mar-2', startTime: '16:00', endTime: '21:00' }
    ]
  },
  {
    day: 'Mié',
    enabled: true,
    intervals: [
      { id: 'int-mie-1', startTime: '06:30', endTime: '14:30' }
    ]
  },
  {
    day: 'Jue',
    enabled: true,
    intervals: [
      { id: 'int-jue-1', startTime: '06:30', endTime: '12:30' },
      { id: 'int-jue-2', startTime: '16:00', endTime: '20:00' }
    ]
  },
  {
    day: 'Vie',
    enabled: true,
    intervals: [
      { id: 'int-vie-1', startTime: '07:30', endTime: '13:30' }
    ]
  },
  {
    day: 'Sáb',
    enabled: true,
    intervals: [
      { id: 'int-sab-1', startTime: '06:30', endTime: '14:00' }
    ]
  },
  {
    day: 'Dom',
    enabled: false,
    intervals: []
  }
]

export const initialScheduleExceptions: ScheduleException[] = []

export const initialEvaluations: Evaluation[] = [
  {
    id: 'eval-01',
    studentId: 'alu-01',
    studentName: 'Alumno de ejemplo 01',
    date: 'Por registrar',
    type: 'Evaluación inicial',
    status: 'Pendiente',
    metrics: [
      { name: 'Peso corporal', value: 'Por definir', unit: 'kg' },
      { name: 'Estatura', value: 'Por definir', unit: 'cm' },
      { name: 'Movilidad articular', value: 'Por evaluar' },
      { name: 'Test fuerza básica', value: 'Por registrar' }
    ],
    observations: 'Evaluación de ingreso pendiente de coordinación en la primera sesión presencial.',
    conductedBy: 'Pablo',
    createdAt: '2026-10-01'
  },
  {
    id: 'eval-02',
    studentId: 'alu-02',
    studentName: 'Alumno de ejemplo 02',
    date: 'Por registrar',
    type: 'Seguimiento',
    status: 'Pendiente',
    metrics: [
      { name: 'Peso corporal', value: 'Por definir', unit: 'kg' },
      { name: 'Movilidad rodilla', value: 'En observación' },
      { name: 'Carga sentadilla', value: 'Por evaluar' }
    ],
    observations: 'Revisión periódica de molestia en rodilla y progreso de técnica.',
    conductedBy: 'Pablo',
    createdAt: '2026-10-05'
  },
  {
    id: 'eval-03',
    studentId: 'alu-03',
    studentName: 'Alumno de ejemplo 03',
    date: 'Por registrar',
    type: 'Seguimiento',
    status: 'Pendiente',
    metrics: [
      { name: 'Resistencia cardiovascular', value: 'Por evaluar' },
      { name: 'Fuerza tren superior', value: 'Por registrar' }
    ],
    observations: 'Control mensual de rendimiento y adaptación a 3 sesiones semanales.',
    conductedBy: 'Pablo',
    createdAt: '2026-10-08'
  },
  {
    id: 'eval-04',
    studentId: 'alu-04',
    studentName: 'Alumno de ejemplo 04',
    date: 'Por registrar',
    type: 'Evaluación inicial',
    status: 'Pendiente',
    metrics: [
      { name: 'Screening funcional', value: 'Pendiente' },
      { name: 'Composición corporal', value: 'Por definir' }
    ],
    observations: 'Ficha inicial para alumno de 4 frecuencias semanales.',
    conductedBy: 'Pablo',
    createdAt: '2026-10-12'
  },
  {
    id: 'eval-05',
    studentId: 'alu-05',
    studentName: 'Alumno de ejemplo 05',
    date: 'Por registrar',
    type: 'Seguimiento',
    status: 'Pendiente',
    metrics: [
      { name: 'Carga de trabajo', value: 'Por evaluar' },
      { name: 'Recuperación', value: 'Adecuada' }
    ],
    observations: 'Monitoreo de adaptación a bloques de martes y jueves.',
    conductedBy: 'Pablo',
    createdAt: '2026-10-15'
  }
]

export const initialContracts: ContractConsent[] = [
  {
    id: 'cnt-01',
    studentId: 'alu-01',
    studentName: 'Alumno de ejemplo 01',
    documentType: 'Términos y condiciones',
    status: 'Pendiente',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: false,
    notes: 'Documento pendiente de envío y confirmación digital del alumno.'
  },
  {
    id: 'cnt-02',
    studentId: 'alu-02',
    studentName: 'Alumno de ejemplo 02',
    documentType: 'Términos y condiciones',
    status: 'Aceptado',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: true,
    notes: 'Aceptación registrada en demostración. Consentimiento informado vigente.'
  },
  {
    id: 'cnt-03',
    studentId: 'alu-03',
    studentName: 'Alumno de ejemplo 03',
    documentType: 'Términos y condiciones',
    status: 'Pendiente',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: false,
    notes: 'Pendiente de aceptación de términos de entrenamiento.'
  },
  {
    id: 'cnt-04',
    studentId: 'alu-04',
    studentName: 'Alumno de ejemplo 04',
    documentType: 'Términos y condiciones',
    status: 'Aceptado',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: true,
    notes: 'Aceptación registrada para plan intensivo.'
  },
  {
    id: 'cnt-05',
    studentId: 'alu-05',
    studentName: 'Alumno de ejemplo 05',
    documentType: 'Términos y condiciones',
    status: 'Pendiente',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: false,
    notes: 'Falta validación de consentimiento informado.'
  }
]

export const initialTermsVersions: TermsVersion[] = [
  {
    version: 'v1.0 (2026)',
    title: 'Términos y Condiciones del Servicio de Entrenamiento — INVERFIT',
    effectiveDate: 'Octubre 2026',
    isCurrent: true,
    content: `CONTRATO DE PRESTACIÓN DE SERVICIOS Y CONSENTIMIENTO INFORMADO (DEMOSTRACIÓN)

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

6. AVISO LEGAL DE DEMOSTRACIÓN
Este documento es un modelo ilustrativo y demostrativo para el prototipo digital de gestión de INVERFIT. No constituye asesoría legal ni firma electrónica validada conforme a la Ley 19.799.`
  }
]
