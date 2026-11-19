import type { Student, Payment, ScheduleBlock, Evaluation, ContractConsent, TermsVersion, DayAvailability, ScheduleException } from '@/types'

export const initialStudents: Student[] = [
  {
    id: 'alu-01',
    fullName: 'Carlos Valenzuela',
    email: 'carlos@ejemplo.com',
    phone: '+56 9 8765 4321',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Lun y Mié 06:30',
    membershipStartDate: '2026-11-01',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Acondicionamiento general y entrenamiento funcional.',
    healthConditions: 'Sin lesiones reportadas. Apto para entrenamiento.',
    emergencyContact: 'Andrea Muñoz (Cónyuge) +56 9 8888 1111',
    createdAt: '2026-10-01'
  },
  {
    id: 'alu-02',
    fullName: 'Nelson Test Onboarding',
    email: 'nelson.test@inverfit.cl',
    phone: '+56 9 9876 5432',
    planName: 'Membresía Mensual',
    frequency: '3 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Lun, Mié y Vie 09:30',
    membershipStartDate: '2026-11-01',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Objetivo: Salud, funcionalidad y longevidad. Ficha completada vía onboarding.',
    healthConditions: 'Cirugía menisco rodilla derecha en 2020. Precaución en rangos profundos.',
    emergencyContact: 'María Valenzuela (Madre) +56 9 1122 3344',
    onboarding: {
      rut: '18.765.432-1',
      birthDate: '1994-05-15',
      age: 32,
      gender: 'Masculino',
      comuna: 'Concepción',
      occupation: 'Ingeniero',
      hasChildren: 'No',
      primaryGoal: 'Salud, funcionalidad y longevidad',
      goalDetails: 'Mejorar postura y ganar fuerza lumbar sin molestias',
      hasTrainedBefore: 'Si',
      trainingExperience: '1-6 meses',
      sportsPracticed: 'Fútbol fin de semana',
      availableDaysWeekly: '2-4 veces a la semana',
      hadCovid: 'No',
      chronicDiseases: 'Ninguna',
      surgeriesInjuries: 'Cirugía menisco rodilla derecha en 2020',
      injuriesPreventExercise: 'Si',
      consentAccepted: true,
      consentTimestamp: 'miércoles, 18 de noviembre de 2026, 5:35:39 p. m.',
      digitalSignature: 'data:image/png;base64,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',
      termsVersion: 'v1.0 (Noviembre 2026 - Oficial Pablo)'
    },
    createdAt: '2026-11-18'
  },
  {
    id: 'alu-03',
    fullName: 'Camila Morales',
    email: 'camila.morales@ejemplo.cl',
    phone: '+56 9 7654 3210',
    planName: 'Membresía Mensual',
    frequency: '3 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Lun, Mié y Vie 11:00',
    membershipStartDate: '2026-11-01',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Acondicionamiento físico general y tonificación.',
    healthConditions: 'Ninguna contraindicación médica.',
    emergencyContact: 'Patricio Morales (Padre) +56 9 9999 3333',
    createdAt: '2026-10-08'
  },
  {
    id: 'alu-04',
    fullName: 'Felipe Arancibia',
    email: 'felipe.arancibia@ejemplo.cl',
    phone: '+56 9 6543 2109',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Mar y Jue 18:30',
    membershipStartDate: '2026-11-01',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Entrenamiento de fuerza e hipertrofia.',
    healthConditions: 'Molestia lumbar leve postural.',
    emergencyContact: 'Francisca Pérez (Pareja) +56 9 9999 4444',
    createdAt: '2026-10-12'
  },
  {
    id: 'alu-05',
    fullName: 'Valeria Soto',
    email: 'valeria.soto@ejemplo.cl',
    phone: '+56 9 5432 1098',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    membershipStatus: 'Por revisar',
    fixedSchedule: 'Mar y Jue 08:00',
    membershipStartDate: 'Por registrar',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Inicio de ciclo. Requiere coordinación de horarios.',
    healthConditions: 'Apta para toda actividad física.',
    emergencyContact: 'Marcela Soto (Hermana) +56 9 9999 5555',
    createdAt: '2026-10-15'
  },
  {
    id: 'alu-06',
    fullName: 'Ignacio Peña',
    email: 'ignacio.pena@ejemplo.cl',
    phone: '+56 9 4321 0987',
    planName: 'Membresía Mensual',
    frequency: '1 vez/sem.',
    membershipStatus: 'Activo',
    fixedSchedule: 'Sáb 09:00',
    membershipStartDate: '2026-11-01',
    membershipDueDate: '1 al 5 de cada mes',
    notes: 'Sesión complementaria de fuerza los sábados.',
    healthConditions: 'Apto sin limitaciones.',
    emergencyContact: 'Rodrigo Peña (Hermano) +56 9 9999 6666',
    createdAt: '2026-10-20'
  }
]

export const initialPayments: Payment[] = [
  {
    id: 'pay-01',
    studentId: 'alu-01',
    studentName: 'Carlos Valenzuela',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    dueDate: '1 al 5 de cada mes',
    amount: 45000,
    status: 'Pagado',
    paymentDate: '2026-11-02',
    paymentMethod: 'SumUp (POS Tarjeta)',
    transactionReference: 'SUM-102941',
    reminderSentCount: 0,
    lastReminderAt: null,
    receiptSent: true,
    receiptSentAt: '2026-11-02 11:30',
    notes: 'Pago con tarjeta presencial vía SumUp en recepción.'
  },
  {
    id: 'pay-02',
    studentId: 'alu-02',
    studentName: 'Nelson Test Onboarding',
    planName: 'Membresía Mensual',
    frequency: '3 veces/sem.',
    dueDate: '1 al 5 de cada mes',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'SumUp (Link de pago)',
    transactionReference: undefined,
    reminderSentCount: 1,
    lastReminderAt: '2026-11-02 10:15',
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Enlace de SumUp generado, pendiente confirmación.'
  },
  {
    id: 'pay-03',
    studentId: 'alu-03',
    studentName: 'Camila Morales',
    planName: 'Membresía Mensual',
    frequency: '3 veces/sem.',
    dueDate: '1 al 5 de cada mes',
    amount: 45000,
    status: 'Pagado',
    paymentDate: '2026-11-03',
    paymentMethod: 'Transferencia Bancaria',
    transactionReference: 'TRF-882914',
    reminderSentCount: 0,
    lastReminderAt: null,
    receiptSent: true,
    receiptSentAt: '2026-11-03 14:20',
    notes: 'Pago registrado por transferencia electrónica.'
  },
  {
    id: 'pay-04',
    studentId: 'alu-04',
    studentName: 'Felipe Arancibia',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    dueDate: '1 al 5 de cada mes',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'SumUp (Link de pago)',
    reminderSentCount: 0,
    lastReminderAt: null,
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Aviso por correo programado para el ciclo del 1 al 5.'
  },
  {
    id: 'pay-05',
    studentId: 'alu-05',
    studentName: 'Valeria Soto',
    planName: 'Membresía Mensual',
    frequency: '2 veces/sem.',
    dueDate: '1 al 5 de cada mes',
    amount: null,
    status: 'Pendiente',
    paymentDate: null,
    paymentMethod: 'SumUp (Link de pago)',
    reminderSentCount: 2,
    lastReminderAt: '2026-11-03 16:40',
    receiptSent: false,
    receiptSentAt: null,
    notes: 'Aviso enviado por correo y enlace SumUp compartido.'
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
      { id: 'int-sab-eval', startTime: '06:30', endTime: '13:00' }
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
    id: 'eval-nelson-official',
    studentId: '6afe0c22633d726169887111',
    studentName: 'Nelson Test Onboarding',
    date: '2026-11-18',
    type: 'Evaluación inicial',
    status: 'Registrada',
    protocol: 'clinico_deportivo',
    vitalSigns: {
      spo2: 98,
      spo2Alert: 'Normal / Óptima',
      systolicBP: 118,
      diastolicBP: 76,
      bpClassification: 'Normal (Óptima)',
      restingHR: 68,
      hrClassification: 'Normocardia (ACSM)'
    },
    anthropometry: {
      weightKg: 78.5,
      heightCm: 176,
      imc: 25.3,
      imcClassification: 'Sobrepeso Leve (OMS)',
      fatPercentage: 18.2,
      fatClassification: 'Normal / Saludable (ACSM)',
      muscleMassKg: 36.4
    },
    dynamometry: {
      rightHandKg: 48.0,
      leftHandKg: 45.0,
      asymmetryPercent: 6.5,
      asymmetryClassification: 'Simetría Óptima (< 10%)',
      strengthClassification: 'Normofuerza (EWGSOP2)'
    },
    romPatterns: {
      kneeToWallLeft: '12 cm',
      kneeToWallRight: '14 cm',
      aslrLeft: '75º',
      aslrRight: '78º',
      overheadSquat: 'Patrón Limpio',
      lungeUnipodal: 'Estable',
      hipHinge: 'Excelente',
      observations: 'Buena movilidad lumbopélvica. Precaución en zancada profunda por antecedente rodilla derecha.'
    },
    isometricStrength: {
      imtpLeftN: 1850,
      imtpRightN: 1890,
      legExtensionLeftN: 680,
      legExtensionRightN: 710,
      legCurlLeftN: 440,
      legCurlRightN: 460,
      hipAbductionLeftN: 330,
      hipAbductionRightN: 345
    },
    boscoJump: {
      squatJumpCm: 31.4,
      squatJumpFlightMs: 506,
      cmjCm: 36.8,
      cmjFlightMs: 548,
      abalakovCm: 42.1,
      abalakovFlightMs: 586,
      dropJumpCm: 33.5,
      dropJumpContactMs: 215
    },
    velocityEncoder: {
      exercise: 'Sentadilla Libre (Encoder ADR)',
      loadKg: 80,
      velocityMs: 0.72,
      estimated1RM: 115.0,
      trainingZone: 'Potencia Máxima (50-70% 1RM)'
    },
    clinicalAlert: 'Atención: Cirugía menisco rodilla derecha (2020) — Monitorear valgo dinámico en carga axial',
    metrics: [
      { name: 'Saturación SpO2', value: '98', unit: '%' },
      { name: 'Presión Arterial (PAS/PAD)', value: '118/76', unit: 'mmHg' },
      { name: 'Frecuencia Cardíaca', value: '68', unit: 'lpm' },
      { name: 'IMC', value: '25.3 (Sobrepeso Leve)', unit: 'kg/m²' },
      { name: 'Fuerza Prensil Der/Izq', value: '48.0 / 45.0', unit: 'kg' },
      { name: 'Asimetría Prensil', value: '6.5% (Simetría Óptima)' },
      { name: '1RM Sentadilla (Encoder)', value: '115.0 kg @ 0.72 m/s' }
    ],
    observations: 'Evaluación inicial completa ingresada con baremos de Pablo. Alumno apto con supervisión en sentadilla profunda.',
    conductedBy: 'Pablo',
    createdAt: '2026-11-18'
  },
  {
    id: 'eval-carlos-valenzuela',
    studentId: '6afde896b2776112af7b290f',
    studentName: 'Carlos Valenzuela',
    date: '2026-11-16',
    type: 'Evaluación inicial',
    status: 'Registrada',
    protocol: 'clinico_deportivo',
    vitalSigns: {
      spo2: 99,
      spo2Alert: 'Normal / Óptima',
      systolicBP: 115,
      diastolicBP: 74,
      bpClassification: 'Normal (Óptima)',
      restingHR: 64,
      hrClassification: 'Normocardia (ACSM)'
    },
    anthropometry: {
      weightKg: 82.0,
      heightCm: 178,
      imc: 25.9,
      imcClassification: 'Sobrepeso Leve (OMS)',
      fatPercentage: 17.5,
      fatClassification: 'Normal (ACSM)',
      muscleMassKg: 38.5
    },
    dynamometry: {
      rightHandKg: 52.0,
      leftHandKg: 50.0,
      asymmetryPercent: 3.9,
      asymmetryClassification: 'Simetría Óptima (< 10%)',
      strengthClassification: 'Normofuerza (EWGSOP2)'
    },
    velocityEncoder: {
      exercise: 'Press de Banca (Encoder ADR)',
      loadKg: 65,
      velocityMs: 0.68,
      estimated1RM: 92.5,
      trainingZone: 'Fuerza - Hipertrofia (70-85% 1RM)'
    },
    metrics: [
      { name: 'Saturación SpO2', value: '99', unit: '%' },
      { name: 'Presión Arterial', value: '115/74', unit: 'mmHg' },
      { name: 'Fuerza Prensil', value: '52.0 / 50.0', unit: 'kg' },
      { name: '1RM Press Banca', value: '92.5 kg' }
    ],
    observations: 'Excelente base muscular. Sin antecedentes clínicos limitantes.',
    conductedBy: 'Pablo',
    createdAt: '2026-11-16'
  },
  {
    id: 'eval-01',
    studentId: 'alu-01',
    studentName: 'Carlos Valenzuela',
    date: '2026-11-10',
    type: 'Evaluación inicial',
    status: 'Registrada',
    protocol: 'clinico_deportivo',
    vitalSigns: {
      spo2: 98,
      systolicBP: 120,
      diastolicBP: 78,
      restingHR: 70
    },
    anthropometry: {
      weightKg: 74.0,
      heightCm: 172,
      imc: 25.0,
      fatPercentage: 19.0,
      muscleMassKg: 33.5
    },
    dynamometry: {
      rightHandKg: 42.0,
      leftHandKg: 40.0,
      asymmetryPercent: 4.9
    },
    metrics: [
      { name: 'Peso corporal', value: '74.0', unit: 'kg' },
      { name: 'Estatura', value: '172', unit: 'cm' },
      { name: 'IMC', value: '25.0', unit: 'kg/m²' },
      { name: 'Fuerza Prensil', value: '42 kg' }
    ],
    observations: 'Evaluación de ingreso completada y apto para inicio de ciclo.',
    conductedBy: 'Pablo',
    createdAt: '2026-11-10'
  },
  {
    id: 'eval-02',
    studentId: 'alu-02',
    studentName: 'Nelson Test Onboarding',
    date: '2026-11-12',
    type: 'Seguimiento',
    status: 'Registrada',
    protocol: 'clinico_deportivo',
    vitalSigns: {
      spo2: 97,
      systolicBP: 122,
      diastolicBP: 80,
      restingHR: 72
    },
    anthropometry: {
      weightKg: 69.5,
      heightCm: 168,
      imc: 24.6,
      fatPercentage: 21.0,
      muscleMassKg: 31.0
    },
    dynamometry: {
      rightHandKg: 38.0,
      leftHandKg: 35.0,
      asymmetryPercent: 8.2
    },
    metrics: [
      { name: 'Peso corporal', value: '69.5', unit: 'kg' },
      { name: 'Estatura', value: '168', unit: 'cm' },
      { name: 'IMC', value: '24.6', unit: 'kg/m²' }
    ],
    observations: 'Revisión periódica de molestia en rodilla y progreso de técnica.',
    conductedBy: 'Pablo',
    createdAt: '2026-11-12'
  },
  {
    id: 'eval-03',
    studentId: 'alu-03',
    studentName: 'Camila Morales',
    date: '2026-11-14',
    type: 'Seguimiento',
    status: 'Pendiente',
    metrics: [
      { name: 'Resistencia cardiovascular', value: 'Por evaluar' },
      { name: 'Fuerza tren superior', value: 'Por registrar' }
    ],
    observations: 'Control mensual agendado para el próximo sábado.',
    conductedBy: 'Pablo',
    createdAt: '2026-11-14'
  }
]

export const initialContracts: ContractConsent[] = [
  {
    id: 'cnt-01',
    studentId: 'alu-01',
    studentName: 'Carlos Valenzuela',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Aceptado',
    date: '2026-11-15',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: true,
    notes: 'Consentimiento informado y reglamento interno aceptado formalmente.'
  },
  {
    id: 'cnt-02',
    studentId: 'alu-02',
    studentName: 'Nelson Test Onboarding',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Aceptado',
    date: '2026-11-18',
    termsVersion: 'v1.0 (Noviembre 2026 - Oficial Pablo)',
    signedCopyAvailable: true,
    notes: 'Consentimiento y reglamento firmado digitalmente en pantalla táctil durante onboarding.'
  },
  {
    id: 'cnt-03',
    studentId: 'alu-03',
    studentName: 'Camila Morales',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Aceptado',
    date: '2026-11-10',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: true,
    notes: 'Consentimiento informado aceptado al ingresar.'
  },
  {
    id: 'cnt-04',
    studentId: 'alu-04',
    studentName: 'Felipe Arancibia',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Pendiente',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: false,
    notes: 'Pendiente de firma o aceptación de términos de entrenamiento.'
  },
  {
    id: 'cnt-05',
    studentId: 'alu-05',
    studentName: 'Valeria Soto',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Pendiente',
    date: 'Por registrar',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: false,
    notes: 'Pendiente de validación de consentimiento informado.'
  },
  {
    id: 'cnt-06',
    studentId: 'alu-06',
    studentName: 'Ignacio Peña',
    documentType: 'Consentimiento informado y reglamento',
    status: 'Aceptado',
    date: '2026-11-05',
    termsVersion: 'v1.0 (2026)',
    signedCopyAvailable: true,
    notes: 'Aceptación registrada para bloques de sábado.'
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

6. RESGUARDO DE ANTECEDENTES Y REGLAMENTO INTERNO
Este reglamento rige las condiciones de entrenamiento personalizado y consentimiento informado para todos los alumnos activos de INVERFIT.`
  }
]
