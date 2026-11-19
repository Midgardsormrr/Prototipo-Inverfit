<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { Student, Evaluation, EvaluationProtocol } from '@/types'
import inverfitLogoWhiteGold from '@/assets/brand/inverfit-logo-white-gold.png'
import {
  Save,
  Printer,
  AlertTriangle,
  User,
  Heart,
  Scale,
  Activity,
  Zap,
  Info,
  Calendar,
  Layers,
  Sparkles,
  FilePlus2
} from 'lucide-vue-next'

const props = defineProps<{
  studentId?: string
  initialStudent?: Student | null
  readOnly?: boolean
}>()

const emit = defineEmits<{
  (e: 'saved', evaluation: Evaluation): void
}>()

const gymStore = useGymStore()

// Active sheet tab (0 to 6)
const activeTab = ref(0)

// Selected Protocol
const selectedProtocol = ref<EvaluationProtocol>('clinico_deportivo')

const sheetTabs = computed(() => [
  {
    id: 0,
    shortName: '1. Ficha Activa',
    fullName: '1. Ficha de Evaluación — Signos Vitales y Antropometría',
    icon: Heart
  },
  {
    id: 1,
    shortName: '2. ROM y Patrones',
    fullName: '2. ROM y Patrones — Análisis Articular y Biomecánica',
    icon: Activity
  },
  {
    id: 2,
    shortName: '3. Fuerza FIM',
    fullName: '3. Fuerza Isométrica Máxima (FIM) & Asimetrías VALD',
    icon: Scale
  },
  {
    id: 3,
    shortName: selectedProtocol.value === 'adulto_mayor' ? '4. Adulto Mayor (SPPB)' : '4. Saltabilidad (Bosco)',
    fullName: selectedProtocol.value === 'adulto_mayor'
      ? '4. Batería Adulto Mayor (Senior Fitness Test / SPPB)'
      : '4. Saltabilidad (Batería Bosco / Ergojump)',
    icon: Zap
  },
  {
    id: 4,
    shortName: '5. Encoder ADR',
    fullName: '5. Perfil Fuerza-Velocidad & 1RM Estimado (Encoder Lineal ADR)',
    icon: Layers
  },
  {
    id: 5,
    shortName: '6. Referencias Clínicas',
    fullName: '6. Compendio Científico y Normas de Referencia',
    icon: Info
  },
  {
    id: 6,
    shortName: '7. Seguimiento',
    fullName: '7. Registro Evolutivo y Comparativa Trimestral',
    icon: Calendar
  }
])

// Selected Student state
const currentStudentId = ref(props.studentId || props.initialStudent?.id || '')

watch(() => props.studentId, (newId) => {
  if (newId) currentStudentId.value = newId
})

watch(() => props.initialStudent, (st) => {
  if (st && !currentStudentId.value) currentStudentId.value = st.id
})

const currentStudent = computed<Student | null>(() => {
  if (props.initialStudent && props.initialStudent.id === currentStudentId.value) {
    return props.initialStudent
  }
  return gymStore.students.find(s => s.id === currentStudentId.value) || null
})

// Auto-detected student demographics
const studentAge = computed<number>(() => {
  if (!currentStudent.value) return 30
  if (currentStudent.value.onboarding?.age) return Number(currentStudent.value.onboarding.age)
  if (currentStudent.value.onboarding?.birthDate) {
    const diff = new Date().getFullYear() - new Date(currentStudent.value.onboarding.birthDate).getFullYear()
    return diff > 0 ? diff : 30
  }
  return 30
})

const studentGender = computed<'M' | 'F'>(() => {
  const g = currentStudent.value?.onboarding?.gender
  if (g === 'Femenino') return 'F'
  return 'M'
})

// --------------------------------------------------------------------------
// 1. FICHA EVALUACIÓN ACTIVA STATE
// --------------------------------------------------------------------------
const evalDate = ref(new Date().toISOString().split('T')[0])
const evalResponsible = ref('Pablo')
const activityLevel = ref('Activa (3-4 días/sem)')
const studentCustomGender = ref<'M' | 'F'>('M')

// Signos Vitales & Antropometría (Inician limpios en null para nuevos alumnos)
const spo2 = ref<number | null>(null)
const pas = ref<number | null>(null)
const pad = ref<number | null>(null)
const fcReposo = ref<number | null>(null)
const pesoKg = ref<number | null>(null)
const tallaM = ref<number | null>(null)
const grasaPct = ref<number | null>(null)
const masaMuscularKg = ref<number | null>(null)
const gripDerKg = ref<number | null>(null)
const gripIzqKg = ref<number | null>(null)

// Calculations (INVERFIT scientific formulas)
// SpO2: =IF(C11="","",IF(C11>=95,"Normal / Óptima",IF(C11>=90,"Hipoxia Leve","Hipoxia Grave")))
const spo2Clasificacion = computed(() => {
  if (spo2.value == null || spo2.value === 0) return '—'
  if (spo2.value >= 95) return 'Normal / Óptima'
  if (spo2.value >= 90) return 'Hipoxia Leve'
  return 'Hipoxia Grave'
})
const spo2Estado = computed(() => {
  if (spo2.value == null || spo2.value === 0) return { text: 'Por evaluar', alert: false, bg: 'bg-transparent text-gray-400' }
  if (spo2.value >= 95) return { text: 'Apto', alert: false, bg: 'bg-[#E6F4EA] text-[#137333]' }
  return { text: 'ALERTA CLÍNICA', alert: true, bg: 'bg-[#FEF7E0] text-[#B06000] font-bold' }
})

// PAS: =IF(C12="","",IF(C12<120,"Normal (Óptima)",IF(C12<130,"Elevada","Hipertensión")))
const pasClasificacion = computed(() => {
  if (pas.value == null || pas.value === 0) return '—'
  if (pas.value < 120) return 'Normal (Óptima)'
  if (pas.value < 130) return 'Elevada'
  return 'Hipertensión'
})
const pasEstado = computed(() => {
  if (pas.value == null || pas.value === 0) return { text: 'Por evaluar', bg: 'bg-transparent text-gray-400' }
  if (pas.value < 120) return { text: 'Apto', bg: 'bg-[#E6F4EA] text-[#137333]' }
  if (pas.value < 130) return { text: 'Monitorear', bg: 'bg-[#FEF7E0] text-[#B06000]' }
  return { text: 'Evaluación Médica', bg: 'bg-[#FCE8E6] text-[#C5221F] font-bold' }
})

// PAD: =IF(C13="","",IF(C13<80,"Normal",IF(C13<90,"Hipertensión N1","Hipertensión N2")))
const padClasificacion = computed(() => {
  if (pad.value == null || pad.value === 0) return '—'
  if (pad.value < 80) return 'Normal'
  if (pad.value < 90) return 'Hipertensión N1'
  return 'Hipertensión N2'
})
const padEstado = computed(() => {
  if (pad.value == null || pad.value === 0) return { text: 'Por evaluar', bg: 'bg-transparent text-gray-400' }
  if (pad.value < 80) return { text: 'Apto', bg: 'bg-[#E6F4EA] text-[#137333]' }
  if (pad.value < 85) return { text: 'Monitorear', bg: 'bg-[#FEF7E0] text-[#B06000]' }
  return { text: 'Evaluación Médica', bg: 'bg-[#FCE8E6] text-[#C5221F] font-bold' }
})

// FC Reposo: =IF(C14="","",IF(C14<60,"Bradicardia Sinusal (Atleta)",IF(C14<=100,"Normocardia","Taquicardia en Reposo")))
const fcClasificacion = computed(() => {
  if (fcReposo.value == null || fcReposo.value === 0) return '—'
  if (fcReposo.value < 60) return 'Bradicardia Sinusal (Atleta)'
  if (fcReposo.value <= 100) return 'Normocardia'
  return 'Taquicardia en Reposo'
})
const fcEstado = computed(() => {
  if (fcReposo.value == null || fcReposo.value === 0) return { text: 'Por evaluar', bg: 'bg-transparent text-gray-400' }
  if (fcReposo.value <= 100) return { text: 'Óptimo', bg: 'bg-[#E6F4EA] text-[#137333]' }
  return { text: 'Revisar Fatiga', bg: 'bg-[#FEF7E0] text-[#B06000]' }
})

// IMC: =C15/(C16^2)
const imcValor = computed(() => {
  if (!pesoKg.value || !tallaM.value || tallaM.value <= 0) return 0
  return Number((pesoKg.value / (tallaM.value * tallaM.value)).toFixed(1))
})
const imcClasificacion = computed(() => {
  const val = imcValor.value
  if (val <= 0) return '—'
  if (val < 18.5) return 'Bajo Peso'
  if (val < 25) return 'Normopeso'
  if (val < 30) return 'Sobrepeso'
  return 'Obesidad'
})
const imcEstado = computed(() => {
  const val = imcValor.value
  if (val <= 0) return { text: 'Por evaluar', bg: 'bg-transparent text-gray-400' }
  if (val < 25) return { text: 'Óptimo', bg: 'bg-[#E6F4EA] text-[#137333]' }
  if (val < 30) return { text: 'Alerta Metabólica', bg: 'bg-[#FEF7E0] text-[#B06000]' }
  return { text: 'Riesgo Elevado', bg: 'bg-[#FCE8E6] text-[#C5221F] font-bold' }
})

// % Grasa Corporal
const grasaClasificacion = computed(() => {
  if (grasaPct.value == null || grasaPct.value <= 0) return '—'
  const g = studentCustomGender.value
  const v = grasaPct.value
  if (g === 'M') {
    if (v <= 14) return 'Atleta / Magro'
    if (v <= 21) return 'Saludable'
    return 'Sobrepeso'
  } else {
    if (v <= 20) return 'Atleta / Magra'
    if (v <= 24) return 'Saludable'
    return 'Sobrepeso'
  }
})

// Dinamometría Prensil: H: >30 kg / M: >20 kg (EWGSOP2)
const gripDerClasificacion = computed(() => {
  if (gripDerKg.value == null || gripDerKg.value === 0) return '—'
  const g = studentCustomGender.value
  const min = g === 'M' ? 30 : 20
  return gripDerKg.value >= min ? 'Normofuerza' : 'Baja Fuerza'
})
const gripIzqClasificacion = computed(() => {
  if (gripIzqKg.value == null || gripIzqKg.value === 0) return '—'
  const g = studentCustomGender.value
  const min = g === 'M' ? 30 : 20
  return gripIzqKg.value >= min ? 'Normofuerza' : 'Baja Fuerza'
})

// Asimetría: =ABS(C20-C21)/AVERAGE(C20,C21)*100
const gripAsimetriaPct = computed(() => {
  const der = gripDerKg.value || 0
  const izq = gripIzqKg.value || 0
  if (der === 0 || izq === 0) return 0
  const avg = (der + izq) / 2
  if (avg <= 0) return 0
  return Number((Math.abs(der - izq) / avg * 100).toFixed(1))
})
const gripAsimetriaClasificacion = computed(() => {
  if (gripDerKg.value == null || gripIzqKg.value == null || gripDerKg.value === 0 || gripIzqKg.value === 0) return '—'
  return gripAsimetriaPct.value <= 10 ? 'Simetría Óptima' : 'Asimetría a Revisar'
})
const gripAsimetriaEstado = computed(() => {
  if (gripDerKg.value == null || gripIzqKg.value == null || gripDerKg.value === 0 || gripIzqKg.value === 0) {
    return { text: 'Por evaluar', bg: 'bg-transparent text-gray-400' }
  }
  if (gripAsimetriaPct.value <= 10) {
    return { text: 'Apto', bg: 'bg-[#E6F4EA] text-[#137333]' }
  }
  return { text: 'Trabajar Unilateral', bg: 'bg-[#FEF7E0] text-[#B06000] font-bold' }
})

// --------------------------------------------------------------------------
// 2. ROM Y PATRONES STATE
// --------------------------------------------------------------------------
const romKneeToWallIzq = ref('')
const romKneeToWallDer = ref('')
const romAslrIzq = ref('')
const romAslrDer = ref('')
const romSquatIzq = ref('')
const romSquatDer = ref('')
const romLungeIzq = ref('')
const romLungeDer = ref('')
const romHingeIzq = ref('')
const romHingeDer = ref('')
const romObs = ref('')

// --------------------------------------------------------------------------
// 3. FUERZA ISOMÉTRICA (FIM) STATE
// --------------------------------------------------------------------------
const fimImtpIzq = ref<number | null>(null)
const fimImtpDer = ref<number | null>(null)
const fimLegExtIzq = ref<number | null>(null)
const fimLegExtDer = ref<number | null>(null)
const fimLegCurlIzq = ref<number | null>(null)
const fimLegCurlDer = ref<number | null>(null)
const fimHipAbdIzq = ref<number | null>(null)
const fimHipAbdDer = ref<number | null>(null)

function calcAsym(izq: number | null, der: number | null): number | null {
  if (izq == null || der == null || izq <= 0 || der <= 0) return null
  const avg = (izq + der) / 2
  if (avg <= 0) return null
  return Number((Math.abs(izq - der) / avg * 100).toFixed(1))
}

function calcRel(val: number | null, peso: number | null): number | null {
  if (val == null || !peso || peso <= 0) return null
  return Number((val / peso).toFixed(1))
}

// --------------------------------------------------------------------------
// 4. SALTABILIDAD BOSCO / SENIOR FITNESS
// --------------------------------------------------------------------------
const boscoSjCm = ref<number | null>(null)
const boscoSjFlightMs = computed(() => {
  if (boscoSjCm.value == null || boscoSjCm.value <= 0) return null
  return Math.round(Math.sqrt((8 * (boscoSjCm.value / 100)) / 9.81) * 1000)
})
const boscoCmjCm = ref<number | null>(null)
const boscoCmjFlightMs = computed(() => {
  if (boscoCmjCm.value == null || boscoCmjCm.value <= 0) return null
  return Math.round(Math.sqrt((8 * (boscoCmjCm.value / 100)) / 9.81) * 1000)
})
const boscoAbalakovCm = ref<number | null>(null)
const boscoAbalakovFlightMs = computed(() => {
  if (boscoAbalakovCm.value == null || boscoAbalakovCm.value <= 0) return null
  return Math.round(Math.sqrt((8 * (boscoAbalakovCm.value / 100)) / 9.81) * 1000)
})
const boscoDjCm = ref<number | null>(null)
const boscoDjContactMs = ref<number | null>(null)

// Elastic index CEA: (CMJ - SJ) / SJ * 100
const boscoIndiceElastico = computed(() => {
  if (boscoSjCm.value == null || boscoCmjCm.value == null || boscoSjCm.value <= 0) return null
  return Number(((boscoCmjCm.value - boscoSjCm.value) / boscoSjCm.value * 100).toFixed(1))
})

// Coordinación brazos: (CMJa - CMJ) / CMJ * 100
const boscoCoordBrazos = computed(() => {
  if (boscoCmjCm.value == null || boscoAbalakovCm.value == null || boscoCmjCm.value <= 0) return null
  return Number(((boscoAbalakovCm.value - boscoCmjCm.value) / boscoCmjCm.value * 100).toFixed(1))
})

// RSI: Altura_m / Contacto_s
const boscoRsi = computed(() => {
  if (boscoDjContactMs.value == null || boscoDjCm.value == null || boscoDjContactMs.value <= 0) return null
  const m = boscoDjCm.value / 100
  const s = boscoDjContactMs.value / 1000
  return Number((m / s).toFixed(2))
})

// Senior Fitness Test (SPPB)
const seniorSitToStand = ref<number | null>(null)
const seniorTug = ref<number | null>(null)
const seniorChairStand = ref<number | null>(null)
const seniorArmCurl = ref<number | null>(null)
const senior6mwt = ref<number | null>(null)

// --------------------------------------------------------------------------
// 5. PERFIL F-V (ENCODER ADR) - GONZÁLEZ-BADILLO & SÁNCHEZ-MEDINA (2010)
// --------------------------------------------------------------------------
const adrBenchLoad = ref<number | null>(null)
const adrBenchVmp = ref<number | null>(null)
const adrBenchRpe = ref<number | null>(null)

const adrSquatLoad = ref<number | null>(null)
const adrSquatVmp = ref<number | null>(null)
const adrSquatRpe = ref<number | null>(null)

const adrDeadliftLoad = ref<number | null>(null)
const adrDeadliftVmp = ref<number | null>(null)
const adrDeadliftRpe = ref<number | null>(null)

const adrPullupLoad = ref<number | null>(null)
const adrPullupVmp = ref<number | null>(null)
const adrPullupRpe = ref<number | null>(null)

function calc1RM(exercise: 'bench' | 'squat', load: number | null, vmp: number | null): number | null {
  if (load == null || vmp == null || load <= 0 || vmp <= 0) return null
  let pct = 0
  if (exercise === 'bench') {
    pct = (-59.57 * (vmp * vmp) - 62.7 * vmp + 100.07) / 100
  } else {
    pct = (-59.65 * (vmp * vmp) - 53.51 * vmp + 113.38) / 100
  }
  if (pct <= 0) return null
  return Number((load / pct).toFixed(1))
}

function calcPct1RM(exercise: 'bench' | 'squat', vmp: number | null): number | null {
  if (vmp == null || vmp <= 0) return null
  let pct = 0
  if (exercise === 'bench') {
    pct = -59.57 * (vmp * vmp) - 62.7 * vmp + 100.07
  } else {
    pct = -59.65 * (vmp * vmp) - 53.51 * vmp + 113.38
  }
  if (pct <= 0) return null
  return Number(pct.toFixed(1))
}

function getZone(vmp: number | null): string {
  if (vmp == null || vmp <= 0) return '—'
  if (vmp > 1.0) return 'Velocidad Pura / Explosiva (30-50% 1RM)'
  if (vmp > 0.7) return 'Potencia Máxima (50-70% 1RM)'
  if (vmp > 0.5) return 'Fuerza-Hipertrofia (70-85% 1RM)'
  return 'Fuerza Máxima (>85% 1RM)'
}

// --------------------------------------------------------------------------
// LIMPIAR TODOS LOS CAMPOS (PARA NUEVO ALUMNO O NUEVA EVALUACIÓN)
// --------------------------------------------------------------------------
function clearWorkbookFields() {
  evalDate.value = new Date().toISOString().split('T')[0]
  spo2.value = null
  pas.value = null
  pad.value = null
  fcReposo.value = null
  pesoKg.value = null
  tallaM.value = null
  grasaPct.value = null
  masaMuscularKg.value = null
  gripDerKg.value = null
  gripIzqKg.value = null
  romKneeToWallIzq.value = ''
  romKneeToWallDer.value = ''
  romAslrIzq.value = ''
  romAslrDer.value = ''
  romSquatIzq.value = ''
  romSquatDer.value = ''
  romLungeIzq.value = ''
  romLungeDer.value = ''
  romHingeIzq.value = ''
  romHingeDer.value = ''
  romObs.value = ''
  fimImtpIzq.value = null
  fimImtpDer.value = null
  fimLegExtIzq.value = null
  fimLegExtDer.value = null
  fimLegCurlIzq.value = null
  fimLegCurlDer.value = null
  fimHipAbdIzq.value = null
  fimHipAbdDer.value = null
  boscoSjCm.value = null
  boscoCmjCm.value = null
  boscoAbalakovCm.value = null
  boscoDjCm.value = null
  boscoDjContactMs.value = null
  seniorSitToStand.value = null
  seniorTug.value = null
  seniorChairStand.value = null
  seniorArmCurl.value = null
  senior6mwt.value = null
  adrBenchLoad.value = null
  adrBenchVmp.value = null
  adrBenchRpe.value = null
  adrSquatLoad.value = null
  adrSquatVmp.value = null
  adrSquatRpe.value = null
  adrDeadliftLoad.value = null
  adrDeadliftVmp.value = null
  adrDeadliftRpe.value = null
  adrPullupLoad.value = null
  adrPullupVmp.value = null
  adrPullupRpe.value = null
}

// --------------------------------------------------------------------------
// LOAD EXISTING EVALUATION OR CLEAN FOR NEW EVALUATION
// --------------------------------------------------------------------------
function loadExistingOrBaselineData() {
  clearWorkbookFields()
  if (!currentStudent.value) return

  studentCustomGender.value = studentGender.value

  const existing = gymStore.evaluations.find(e => e.studentId === currentStudent.value!.id && e.status === 'Registrada')
  if (existing) {
    if (existing.date && existing.date !== 'Por registrar') evalDate.value = existing.date
    if (existing.protocol) selectedProtocol.value = existing.protocol
    if (existing.conductedBy) evalResponsible.value = existing.conductedBy

    if (existing.vitalSigns) {
      if (existing.vitalSigns.spo2 != null) spo2.value = existing.vitalSigns.spo2
      if (existing.vitalSigns.systolicBP != null) pas.value = existing.vitalSigns.systolicBP
      if (existing.vitalSigns.diastolicBP != null) pad.value = existing.vitalSigns.diastolicBP
      if (existing.vitalSigns.restingHR != null) fcReposo.value = existing.vitalSigns.restingHR
    }
    if (existing.anthropometry) {
      if (existing.anthropometry.weightKg != null) pesoKg.value = existing.anthropometry.weightKg
      if (existing.anthropometry.heightCm != null) tallaM.value = Number((existing.anthropometry.heightCm / 100).toFixed(2))
      if (existing.anthropometry.fatPercentage != null) grasaPct.value = existing.anthropometry.fatPercentage
      if (existing.anthropometry.muscleMassKg != null) masaMuscularKg.value = existing.anthropometry.muscleMassKg
    }
    if (existing.dynamometry) {
      if (existing.dynamometry.rightHandKg != null) gripDerKg.value = existing.dynamometry.rightHandKg
      if (existing.dynamometry.leftHandKg != null) gripIzqKg.value = existing.dynamometry.leftHandKg
    }
    if (existing.romPatterns) {
      if (existing.romPatterns.kneeToWallLeft) romKneeToWallIzq.value = existing.romPatterns.kneeToWallLeft
      if (existing.romPatterns.kneeToWallRight) romKneeToWallDer.value = existing.romPatterns.kneeToWallRight
      if (existing.romPatterns.aslrLeft) romAslrIzq.value = existing.romPatterns.aslrLeft
      if (existing.romPatterns.aslrRight) romAslrDer.value = existing.romPatterns.aslrRight
      if (existing.romPatterns.overheadSquat) romSquatIzq.value = existing.romPatterns.overheadSquat
      if (existing.romPatterns.lungeUnipodal) romLungeIzq.value = existing.romPatterns.lungeUnipodal
      if (existing.romPatterns.hipHinge) romHingeIzq.value = existing.romPatterns.hipHinge
      if (existing.romPatterns.observations) romObs.value = existing.romPatterns.observations
    }
    if (existing.isometricStrength) {
      if (existing.isometricStrength.imtpLeftN != null) fimImtpIzq.value = existing.isometricStrength.imtpLeftN
      if (existing.isometricStrength.imtpRightN != null) fimImtpDer.value = existing.isometricStrength.imtpRightN
      if (existing.isometricStrength.legExtensionLeftN != null) fimLegExtIzq.value = existing.isometricStrength.legExtensionLeftN
      if (existing.isometricStrength.legExtensionRightN != null) fimLegExtDer.value = existing.isometricStrength.legExtensionRightN
      if (existing.isometricStrength.legCurlLeftN != null) fimLegCurlIzq.value = existing.isometricStrength.legCurlLeftN
      if (existing.isometricStrength.legCurlRightN != null) fimLegCurlDer.value = existing.isometricStrength.legCurlRightN
      if (existing.isometricStrength.hipAbductionLeftN != null) fimHipAbdIzq.value = existing.isometricStrength.hipAbductionLeftN
      if (existing.isometricStrength.hipAbductionRightN != null) fimHipAbdDer.value = existing.isometricStrength.hipAbductionRightN
    }
    if (existing.boscoJump) {
      if (existing.boscoJump.squatJumpCm != null) boscoSjCm.value = existing.boscoJump.squatJumpCm
      if (existing.boscoJump.cmjCm != null) boscoCmjCm.value = existing.boscoJump.cmjCm
      if (existing.boscoJump.abalakovCm != null) boscoAbalakovCm.value = existing.boscoJump.abalakovCm
      if (existing.boscoJump.dropJumpCm != null) boscoDjCm.value = existing.boscoJump.dropJumpCm
      if (existing.boscoJump.dropJumpContactMs != null) boscoDjContactMs.value = existing.boscoJump.dropJumpContactMs
    }
    if (existing.seniorFitness) {
      if (existing.seniorFitness.sitToStand5s != null) seniorSitToStand.value = existing.seniorFitness.sitToStand5s
      if (existing.seniorFitness.chairStandReps != null) seniorChairStand.value = existing.seniorFitness.chairStandReps
      if (existing.seniorFitness.armCurlReps != null) seniorArmCurl.value = existing.seniorFitness.armCurlReps
      if (existing.seniorFitness.walk6MinMeters != null) senior6mwt.value = existing.seniorFitness.walk6MinMeters
    }
    if (existing.velocityEncoder) {
      if (existing.velocityEncoder.loadKg != null) adrSquatLoad.value = existing.velocityEncoder.loadKg
      if (existing.velocityEncoder.velocityMs != null) adrSquatVmp.value = existing.velocityEncoder.velocityMs
    }
  } else {
    // Alumno nuevo (sin evaluación registrada): campos quedan totalmente limpios
    if (currentStudent.value.healthConditions && currentStudent.value.healthConditions !== 'Ninguna') {
      romObs.value = `Antecedente médico declarado: ${currentStudent.value.healthConditions}`
    }
  }
}

onMounted(() => {
  loadExistingOrBaselineData()
})

watch(currentStudent, (st) => {
  if (st) {
    if (studentAge.value >= 60) {
      selectedProtocol.value = 'adulto_mayor'
    } else {
      selectedProtocol.value = 'clinico_deportivo'
    }
    loadExistingOrBaselineData()
  }
})

// --------------------------------------------------------------------------
// SAVE ACTION: SYNC BACK TO GYMSTORE
// --------------------------------------------------------------------------
const isSaving = ref(false)

async function handleSaveWorkbook() {
  if (!currentStudent.value) return
  isSaving.value = true

  try {
    const existing = gymStore.evaluations.find(e => e.studentId === currentStudent.value!.id)

    const metricsList = [
      { name: 'Saturación SpO2', value: spo2.value != null ? `${spo2.value}` : '—', unit: '%' },
      { name: 'Presión Arterial (PAS/PAD)', value: pas.value != null ? `${pas.value}/${pad.value || ''}` : '—', unit: 'mmHg' },
      { name: 'Frecuencia Cardíaca', value: fcReposo.value != null ? `${fcReposo.value}` : '—', unit: 'lpm' },
      { name: 'Peso / Talla', value: pesoKg.value != null ? `${pesoKg.value} kg / ${tallaM.value || ''} m` : '—' },
      { name: 'IMC', value: imcValor.value > 0 ? `${imcValor.value} (${imcClasificacion.value})` : '—', unit: 'kg/m²' },
      { name: 'Fuerza Prensil Der/Izq', value: gripDerKg.value != null ? `${gripDerKg.value} / ${gripIzqKg.value || ''}` : '—', unit: 'kg' },
      { name: 'Asimetría Prensil', value: gripAsimetriaPct.value > 0 ? `${gripAsimetriaPct.value}% (${gripAsimetriaClasificacion.value})` : '—' }
    ]

    const evalPayload = {
      studentId: currentStudent.value.id,
      studentName: currentStudent.value.fullName,
      date: evalDate.value,
      type: (existing?.type || 'Evaluación inicial') as any,
      status: 'Registrada' as const,
      protocol: selectedProtocol.value,
      conductedBy: evalResponsible.value,
      vitalSigns: {
        spo2: spo2.value ?? undefined,
        spo2Alert: spo2Clasificacion.value !== '—' ? spo2Clasificacion.value : undefined,
        systolicBP: pas.value ?? undefined,
        diastolicBP: pad.value ?? undefined,
        bpClassification: pasClasificacion.value !== '—' ? pasClasificacion.value : undefined,
        restingHR: fcReposo.value ?? undefined,
        hrClassification: fcClasificacion.value !== '—' ? fcClasificacion.value : undefined
      },
      anthropometry: {
        weightKg: pesoKg.value ?? undefined,
        heightCm: tallaM.value ? Math.round(tallaM.value * 100) : undefined,
        imc: imcValor.value > 0 ? imcValor.value : undefined,
        imcClassification: imcClasificacion.value !== '—' ? imcClasificacion.value : undefined,
        fatPercentage: grasaPct.value ?? undefined,
        fatClassification: grasaClasificacion.value !== '—' ? grasaClasificacion.value : undefined,
        muscleMassKg: masaMuscularKg.value ?? undefined
      },
      dynamometry: {
        rightHandKg: gripDerKg.value ?? undefined,
        leftHandKg: gripIzqKg.value ?? undefined,
        asymmetryPercent: gripAsimetriaPct.value > 0 ? gripAsimetriaPct.value : undefined,
        asymmetryClassification: gripAsimetriaClasificacion.value !== '—' ? gripAsimetriaClasificacion.value : undefined,
        strengthClassification: gripDerClasificacion.value !== '—' ? gripDerClasificacion.value : undefined
      },
      romPatterns: {
        kneeToWallLeft: romKneeToWallIzq.value || undefined,
        kneeToWallRight: romKneeToWallDer.value || undefined,
        aslrLeft: romAslrIzq.value || undefined,
        aslrRight: romAslrDer.value || undefined,
        overheadSquat: romSquatIzq.value || undefined,
        lungeUnipodal: romLungeIzq.value || undefined,
        hipHinge: romHingeIzq.value || undefined,
        observations: romObs.value || undefined
      },
      isometricStrength: {
        imtpLeftN: fimImtpIzq.value ?? undefined,
        imtpRightN: fimImtpDer.value ?? undefined,
        legExtensionLeftN: fimLegExtIzq.value ?? undefined,
        legExtensionRightN: fimLegExtDer.value ?? undefined,
        legCurlLeftN: fimLegCurlIzq.value ?? undefined,
        legCurlRightN: fimLegCurlDer.value ?? undefined,
        hipAbductionLeftN: fimHipAbdIzq.value ?? undefined,
        hipAbductionRightN: fimHipAbdDer.value ?? undefined
      },
      boscoJump: selectedProtocol.value === 'clinico_deportivo' ? {
        squatJumpCm: boscoSjCm.value ?? undefined,
        squatJumpFlightMs: boscoSjFlightMs.value ?? undefined,
        cmjCm: boscoCmjCm.value ?? undefined,
        cmjFlightMs: boscoCmjFlightMs.value ?? undefined,
        abalakovCm: boscoAbalakovCm.value ?? undefined,
        abalakovFlightMs: boscoAbalakovFlightMs.value ?? undefined,
        dropJumpCm: boscoDjCm.value ?? undefined,
        dropJumpContactMs: boscoDjContactMs.value ?? undefined
      } : undefined,
      seniorFitness: selectedProtocol.value === 'adulto_mayor' ? {
        sitToStand5s: seniorSitToStand.value ?? undefined,
        sitToStandAlert: seniorSitToStand.value != null ? (seniorSitToStand.value < 12 ? 'Sin Riesgo de Caída' : 'Alerta Riesgo de Caída') : undefined,
        chairStandReps: seniorChairStand.value ?? undefined,
        armCurlReps: seniorArmCurl.value ?? undefined,
        walk6MinMeters: senior6mwt.value ?? undefined
      } : undefined,
      velocityEncoder: {
        exercise: 'Sentadilla Libre (Encoder ADR)',
        loadKg: adrSquatLoad.value ?? undefined,
        velocityMs: adrSquatVmp.value ?? undefined,
        estimated1RM: calc1RM('squat', adrSquatLoad.value, adrSquatVmp.value) ?? undefined,
        trainingZone: adrSquatVmp.value ? getZone(adrSquatVmp.value) : undefined
      },
      clinicalAlert: currentStudent.value.healthConditions || undefined,
      metrics: metricsList,
      observations: `Evaluación científica oficial registrada por ${evalResponsible.value}.`
    }

    let savedEval: Evaluation
    if (existing) {
      const updated = await gymStore.updateEvaluation(existing.id, evalPayload)
      savedEval = updated || (evalPayload as any)
    } else {
      savedEval = await gymStore.createEvaluation(evalPayload)
    }

    emit('saved', savedEval)
    gymStore.showToast(`Evaluación de ${currentStudent.value.fullName} guardada correctamente.`, 'success', 'Evaluación Guardada')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error al guardar la evaluación'
    gymStore.showToast(msg, 'error')
  } finally {
    isSaving.value = false
  }
}

function handleClearForm() {
  if (confirm('¿Limpiar los campos para ingresar una nueva evaluación en blanco?')) {
    clearWorkbookFields()
    if (currentStudent.value?.healthConditions && currentStudent.value.healthConditions !== 'Ninguna') {
      romObs.value = `Antecedente médico declarado: ${currentStudent.value.healthConditions}`
    }
    gymStore.showToast('Formulario listo en blanco para nuevo registro.', 'info')
  }
}

function handleLoadDemoValues() {
  if (confirm('¿Cargar valores de ejemplo de Pablo para demostración o prueba?')) {
    spo2.value = 98
    pas.value = 118
    pad.value = 76
    fcReposo.value = 68
    pesoKg.value = 78.5
    tallaM.value = 1.76
    grasaPct.value = 18.2
    masaMuscularKg.value = 36.4
    gripDerKg.value = 48.0
    gripIzqKg.value = 45.0
    romKneeToWallIzq.value = '12 cm'
    romKneeToWallDer.value = '14 cm'
    romAslrIzq.value = '75º'
    romAslrDer.value = '78º'
    romSquatIzq.value = 'Patrón Limpio'
    romSquatDer.value = 'Patrón Limpio'
    romLungeIzq.value = 'Estable'
    romLungeDer.value = 'Estable'
    romHingeIzq.value = 'Excelente'
    romHingeDer.value = 'Excelente'
    romObs.value = 'Buena dorsiflexión sin colapso medial del arco. Excelente disociación lumbopélvica.'
    fimImtpIzq.value = 1850
    fimImtpDer.value = 1890
    fimLegExtIzq.value = 680
    fimLegExtDer.value = 710
    fimLegCurlIzq.value = 440
    fimLegCurlDer.value = 460
    fimHipAbdIzq.value = 330
    fimHipAbdDer.value = 345
    boscoSjCm.value = 31.4
    boscoCmjCm.value = 36.8
    boscoAbalakovCm.value = 42.1
    boscoDjCm.value = 33.5
    boscoDjContactMs.value = 215
    seniorSitToStand.value = 10.4
    seniorTug.value = 8.2
    seniorChairStand.value = 15
    seniorArmCurl.value = 18
    senior6mwt.value = 490
    adrBenchLoad.value = 60
    adrBenchVmp.value = 0.65
    adrBenchRpe.value = 8
    adrSquatLoad.value = 80
    adrSquatVmp.value = 0.72
    adrSquatRpe.value = 7.5
    adrDeadliftLoad.value = 90
    adrDeadliftVmp.value = 0.68
    adrDeadliftRpe.value = 8
    adrPullupLoad.value = 15
    adrPullupVmp.value = 0.58
    adrPullupRpe.value = 8.5
    gymStore.showToast('Valores de demostración cargados.', 'info')
  }
}

function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="bg-white rounded-[10px] border border-[#E2DFD8] shadow-2xs overflow-hidden text-[#191A1D]">
    <!-- ==================================================================== -->
    <!-- 1. INVERFIT BRAND HEADER (Official Graphite Charcoal & Gold Accent)  -->
    <!-- ==================================================================== -->
    <div class="bg-[#191A1D] text-white px-5 py-4 border-b-2 border-[#CEA549]">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <!-- Official Imagotype Logo & Title -->
        <div class="space-y-1.5">
          <div class="flex items-center gap-3">
            <!-- Official INVERFIT Brand Imagotype -->
            <img
              :src="inverfitLogoWhiteGold"
              alt="INVERFIT"
              class="h-7 w-auto object-contain select-none"
            />
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-[#24262B] border border-[#CEA549]/40 text-[#CEA549] tracking-wider uppercase">
              FICHA DE INGRESO & EVALUACIÓN CIENTÍFICA
            </span>
          </div>

          <h2 class="text-[16px] sm:text-[18px] font-bold text-white tracking-tight">
            EVALUACIÓN CIENTÍFICA & FUNCIONAL INVERFIT
          </h2>
          <p class="text-[12px] text-[#8E9097]">
            CENTRO DE ENTRENAMIENTO — PROTOCOLO CLÍNICO-DEPORTIVO, SIGNOS VITALES, ROM, FIM & ENCODER ADR
          </p>
        </div>

        <!-- Student Selector & Quick Actions -->
        <div class="flex flex-wrap items-center gap-2 self-start md:self-auto">
          <!-- Student Switcher Dropdown (when not fixed) -->
          <div v-if="!props.studentId" class="flex items-center bg-[#24262B] border border-[#3A3D46] rounded-md px-2.5 py-1.5">
            <User class="w-3.5 h-3.5 text-[#CEA549] mr-2" />
            <select
              v-model="currentStudentId"
              class="bg-transparent text-xs text-white font-semibold focus:outline-none cursor-pointer pr-2"
            >
              <option v-for="st in gymStore.students" :key="st.id" :value="st.id" class="bg-[#191A1D] text-white">
                {{ st.fullName }} ({{ st.frequency }})
              </option>
            </select>
          </div>

          <button
            @click="handleClearForm"
            title="Limpiar campos para ingresar una nueva evaluación en blanco"
            class="h-[34px] px-3 bg-[#24262B] hover:bg-[#2D3037] text-[#CEA549] border border-[#3A3D46] rounded text-xs font-bold flex items-center gap-1.5 cursor-pointer transition-colors"
          >
            <FilePlus2 class="w-3.5 h-3.5" />
            <span>Nueva en Blanco</span>
          </button>

          <button
            @click="handleLoadDemoValues"
            title="Cargar ejemplo con valores de demostración"
            class="h-[34px] px-3 bg-[#24262B] hover:bg-[#2D3037] text-slate-300 border border-[#3A3D46] rounded text-xs font-semibold flex items-center gap-1.5 cursor-pointer transition-colors"
          >
            <Sparkles class="w-3.5 h-3.5 text-slate-400" />
            <span class="hidden sm:inline">Valores Demo</span>
          </button>

          <button
            @click="handlePrint"
            title="Imprimir informe oficial"
            class="h-[34px] px-3 bg-[#24262B] hover:bg-[#2D3037] text-slate-200 border border-[#3A3D46] rounded text-xs font-semibold flex items-center gap-1.5 cursor-pointer transition-colors"
          >
            <Printer class="w-3.5 h-3.5 text-slate-300" />
            <span class="hidden sm:inline">Imprimir</span>
          </button>

          <button
            @click="handleSaveWorkbook"
            :disabled="isSaving || !currentStudent"
            class="h-[34px] px-4 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded text-xs font-bold flex items-center gap-1.5 cursor-pointer shadow-xs transition-all disabled:opacity-50"
          >
            <Save class="w-3.5 h-3.5" />
            <span>{{ isSaving ? 'Guardando...' : 'Guardar Evaluación' }}</span>
          </button>
        </div>
      </div>

      <!-- Quick Demographics Banner Strip -->
      <div v-if="currentStudent" class="mt-3 pt-3 border-t border-[#282A30] flex flex-wrap items-center justify-between text-xs gap-3">
        <div class="flex items-center gap-4 flex-wrap">
          <span class="text-[#8E9097]">
            <strong class="text-white">{{ currentStudent.fullName }}</strong> · {{ studentAge }} años · Sexo: {{ studentGender }}
          </span>
          <span class="text-[#8E9097]">
            Plan: <span class="text-slate-200 font-medium">{{ currentStudent.frequency }}</span>
          </span>
          <span class="text-[#8E9097]">
            Protocolo:
            <select
              v-model="selectedProtocol"
              class="bg-[#24262B] text-[#CEA549] font-bold border border-[#3A3D46] rounded px-1.5 py-0.5 ml-1 text-[11px]"
            >
              <option value="clinico_deportivo">Clínico-Deportivo (Encoder/Bosco)</option>
              <option value="adulto_mayor">Adulto Mayor (Senior SPPB)</option>
            </select>
          </span>
        </div>

        <div v-if="currentStudent.healthConditions" class="flex items-center gap-1.5 text-amber-300 bg-[#24262B] border border-amber-500/40 px-2 py-0.5 rounded text-[11px]">
          <AlertTriangle class="w-3 h-3 text-amber-400" />
          <span>{{ currentStudent.healthConditions }}</span>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- 2. BALANCED SHEET TABS BAR (All 7 Sheets Fit Cleanly - No Cut Off!)  -->
    <!-- ==================================================================== -->
    <div class="bg-[#24262B] border-b border-[#191A1D] p-1.5">
      <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-1">
        <button
          v-for="tab in sheetTabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :title="tab.fullName"
          class="group flex items-center justify-center gap-1.5 px-2 py-2 text-xs font-semibold rounded-md transition-all cursor-pointer text-center"
          :class="activeTab === tab.id
            ? 'bg-white text-[#191A1D] font-bold shadow-xs border-b-2 border-[#CEA549]'
            : 'text-[#8E9097] hover:text-white hover:bg-[#191A1D]/80 border-b-2 border-transparent'"
        >
          <component
            :is="tab.icon"
            class="w-3.5 h-3.5 flex-shrink-0"
            :class="activeTab === tab.id ? 'text-[#CEA549]' : 'text-[#8E9097] group-hover:text-white'"
          />
          <span class="truncate">{{ tab.shortName }}</span>
        </button>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- 3. SHEET CONTENTS                                                    -->
    <!-- ==================================================================== -->
    <div class="p-4 sm:p-6 bg-[#FAF9F6]">
      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 0: 1. Ficha Evaluación Activa                                -->
      <!-- ------------------------------------------------------------------ -->
      <div v-if="activeTab === 0" class="space-y-6">
        <!-- Sub-Banner -->
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549] flex items-center justify-between">
          <div>
            <h3 class="font-bold text-sm tracking-wide">INVERFIT | FICHA DE EVALUACIÓN CLÍNICO-DEPORTIVA</h3>
            <p class="text-[11px] text-[#8E9097]">CENTRO DE ENTRENAMIENTO — PROTOCOLO FISIOLÓGICO & SIGNOS VITALES</p>
          </div>
          <div class="hidden sm:block text-right text-[11px] text-[#CEA549] font-mono">
            CONFIDENCIAL · COACH PABLO
          </div>
        </div>

        <!-- Section 1: DATOS DEL EVALUADO Y SELECCIÓN DE GÉNERO -->
        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <div class="bg-[#24262B] text-white px-4 py-2 text-xs font-bold tracking-wider uppercase border-l-2 border-[#CEA549]">
            1. DATOS DEL EVALUADO Y SELECCIÓN DE GÉNERO
          </div>

          <div class="p-4 grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Nombre Completo:</label>
              <input
                :value="currentStudent?.fullName || 'Seleccione alumno'"
                readonly
                class="w-full h-8 px-2.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded font-bold text-[#191A1D] cursor-not-allowed"
              />
            </div>

            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Edad (años):</label>
              <input
                :value="studentAge"
                readonly
                class="w-full h-8 px-2.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded font-bold text-[#191A1D] cursor-not-allowed"
              />
            </div>

            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Fecha Evaluación:</label>
              <input
                v-model="evalDate"
                type="date"
                class="w-full h-8 px-2.5 bg-white border border-[#DCD8CF] rounded font-semibold text-[#191A1D] focus:border-[#CEA549] focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Sexo (M/F):</label>
              <select
                v-model="studentCustomGender"
                class="w-full h-8 px-2.5 bg-white border border-[#DCD8CF] rounded font-bold text-[#191A1D] focus:border-[#CEA549] focus:outline-none cursor-pointer"
              >
                <option value="M">Masculino (M)</option>
                <option value="F">Femenino (F)</option>
              </select>
            </div>

            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Evaluador Responsable:</label>
              <input
                v-model="evalResponsible"
                class="w-full h-8 px-2.5 bg-white border border-[#DCD8CF] rounded font-semibold text-[#191A1D] focus:border-[#CEA549] focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-[#6B6D74] font-bold uppercase text-[10px] mb-1">Nivel de Actividad:</label>
              <select
                v-model="activityLevel"
                class="w-full h-8 px-2.5 bg-white border border-[#DCD8CF] rounded font-semibold text-[#191A1D] focus:border-[#CEA549] focus:outline-none cursor-pointer"
              >
                <option value="Sedentario (0-1 día/sem)">Sedentario (0-1 día/sem)</option>
                <option value="Moderada (2-3 días/sem)">Moderada (2-3 días/sem)</option>
                <option value="Activa (3-4 días/sem)">Activa (3-4 días/sem)</option>
                <option value="Muy Activa / Atleta (>5 días/sem)">Muy Activa / Atleta (>5 días/sem)</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Section 2: SIGNOS VITALES Y ANTROPOMETRÍA CLÍNICA -->
        <div class="bg-white border border-[#E2DFD8] rounded-md shadow-2xs overflow-hidden">
          <div class="bg-[#191A1D] text-white px-4 py-2 text-xs font-bold tracking-wider uppercase flex items-center justify-between border-l-4 border-[#CEA549]">
            <span>2. SIGNOS VITALES Y ANTROPOMETRÍA CLÍNICA</span>
            <span class="text-[10px] text-[#CEA549] font-mono">CÁLCULO AUTOMÁTICO (IMC Y ALERTAS CLÍNICAS)</span>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse min-w-[760px]">
              <thead>
                <tr class="bg-[#24262B] text-white border-b border-[#191A1D]">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[240px]">Parámetro Evaluado</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[110px]">Valor Medido</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[60px]">Unidad</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[190px]">Referencia Bibliográfica</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[210px]">Clasificación Automática (OMS / ACSM)</th>
                  <th class="py-2.5 px-3 font-bold text-center w-[150px]">Estado / Alerta Clínica</th>
                </tr>
              </thead>

              <tbody class="divide-y divide-[#E2DFD8]">
                <!-- 1. SpO2 -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Saturación de Oxígeno (SpO2)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="spo2" type="number" placeholder="Ej. 98" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">%</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">95% - 100% (OMS)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ spo2Clasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="spo2Estado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ spo2Estado.text }}
                    </span>
                  </td>
                </tr>

                <!-- 2. PAS -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Presión Arterial Sistólica (PAS)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="pas" type="number" placeholder="Ej. 120" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">mmHg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">&lt; 120 mmHg (AHA / ACC 2017)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ pasClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="pasEstado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ pasEstado.text }}
                    </span>
                  </td>
                </tr>

                <!-- 3. PAD -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Presión Arterial Diastólica (PAD)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="pad" type="number" placeholder="Ej. 80" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">mmHg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">&lt; 80 mmHg (AHA / ACC 2017)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ padClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="padEstado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ padEstado.text }}
                    </span>
                  </td>
                </tr>

                <!-- 4. FC Reposo -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Frecuencia Cardíaca en Reposo</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fcReposo" type="number" placeholder="Ej. 68" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">lpm</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">60 - 100 lpm (ACSM)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ fcClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="fcEstado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ fcEstado.text }}
                    </span>
                  </td>
                </tr>

                <!-- 5. Peso Corporal -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Peso Corporal</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="pesoKg" type="number" step="0.1" placeholder="Ej. 70.0" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">kg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-500">Registro Directo</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Ingresado</td>
                  <td class="py-2 px-3 text-center text-slate-500 font-mono">Baseline</td>
                </tr>

                <!-- 6. Estatura / Talla -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Estatura / Talla</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="tallaM" type="number" step="0.01" placeholder="Ej. 1.75" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">m</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-500">Registro Directo</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Ingresado</td>
                  <td class="py-2 px-3 text-center text-slate-500 font-mono">Baseline</td>
                </tr>

                <!-- 7. IMC (YELLOW FORMULA CELL #FFFF99) -->
                <tr class="bg-amber-50/50 hover:bg-amber-100/60 transition-colors">
                  <td class="py-2 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">
                    Índice de Masa Corporal (IMC)
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black text-[#191A1D] border border-amber-300">
                    {{ imcValor > 0 ? imcValor : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-600 font-mono">kg/m²</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-700 font-mono">18.5 - 24.9 kg/m² (OMS)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-slate-800">{{ imcClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="imcEstado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ imcEstado.text }}
                    </span>
                  </td>
                </tr>

                <!-- 8. Porcentaje Grasa Corporal -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Porcentaje Grasa Corporal</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="grasaPct" type="number" step="0.1" placeholder="Ej. 18.0" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">%</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">H: 11-21% / M: 15-24% (ACSM)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ grasaClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span class="px-2 py-0.5 rounded text-[11px] bg-[#E6F4EA] text-[#137333] font-semibold inline-block">
                      Apto
                    </span>
                  </td>
                </tr>

                <!-- 9. Masa Muscular Estimada -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Masa Muscular Estimada</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="masaMuscularKg" type="number" step="0.1" placeholder="Ej. 32.0" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">kg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Referencia Individual Bioimpedancia</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-700">Masa Muscular Adecuada</td>
                  <td class="py-2 px-3 text-center text-slate-600 font-mono">Normonutrición</td>
                </tr>

                <!-- 10. Fuerza Prensil - Derecha -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Fuerza Prensil - Derecha</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="gripDerKg" type="number" step="0.5" placeholder="Ej. 45" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">kg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">H: &gt;30 kg / M: &gt;20 kg (EWGSOP2)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ gripDerClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="gripDerKg != null ? (gripDerKg >= (studentCustomGender === 'M' ? 30 : 20) ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]') : 'bg-transparent text-gray-400'" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ gripDerKg != null ? (gripDerKg >= (studentCustomGender === 'M' ? 30 : 20) ? 'Óptimo' : 'Entrenar Fuerza') : 'Por evaluar' }}
                    </span>
                  </td>
                </tr>

                <!-- 11. Fuerza Prensil - Izquierda -->
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="py-2 px-3 font-semibold text-[#191A1D] border-r border-[#ECE8E1]">Fuerza Prensil - Izquierda</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="gripIzqKg" type="number" step="0.5" placeholder="Ej. 42" class="w-16 h-7 text-center font-bold border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">kg</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">H: &gt;30 kg / M: &gt;20 kg (EWGSOP2)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ gripIzqClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="gripIzqKg != null ? (gripIzqKg >= (studentCustomGender === 'M' ? 30 : 20) ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]') : 'bg-transparent text-gray-400'" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ gripIzqKg != null ? (gripIzqKg >= (studentCustomGender === 'M' ? 30 : 20) ? 'Óptimo' : 'Entrenar Fuerza') : 'Por evaluar' }}
                    </span>
                  </td>
                </tr>

                <!-- 12. Asimetría Fuerza Prensil (%) (YELLOW FORMULA CELL #FFFF99) -->
                <tr class="bg-amber-50/50 hover:bg-amber-100/60 transition-colors">
                  <td class="py-2 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">
                    Asimetría Fuerza Prensil (%)
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black text-[#191A1D] border border-amber-300">
                    {{ gripAsimetriaPct > 0 ? gripAsimetriaPct + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-600 font-mono">%</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-700 font-mono">&lt; 10% Asimetría Funcional</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-slate-800">{{ gripAsimetriaClasificacion }}</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="gripAsimetriaEstado.bg" class="px-2 py-0.5 rounded text-[11px] inline-block">
                      {{ gripAsimetriaEstado.text }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 1: 2. ROM y Patrones                                         -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 1" class="space-y-6">
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
          <h3 class="font-bold text-sm tracking-wide">INVERFIT | ANÁLISIS ARTICULAR, ROM Y PATRONES BÁSICOS</h3>
          <p class="text-[11px] text-[#8E9097]">VALORACIÓN CUALITATIVA Y GONIOMÉTRICA DE MOVILIDAD Y ESTABILIDAD</p>
        </div>

        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse min-w-[760px]">
              <thead>
                <tr class="bg-[#24262B] text-white border-b border-[#191A1D]">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[240px]">Test / Patrón Evaluado</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">Izquierda</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">Derecha</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[150px]">Calificación Asimetría</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[170px]">Estándar de Referencia</th>
                  <th class="py-2.5 px-3 font-bold">Observación Biomecánica Detallada</th>
                </tr>
              </thead>

              <tbody class="divide-y divide-[#ECE8E1]">
                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Dorsiflexión de Tobillo (Knee to Wall)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romKneeToWallIzq" placeholder="Ej. 12 cm" class="w-20 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romKneeToWallDer" placeholder="Ej. 14 cm" class="w-20 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-emerald-700">Simétrica / Óptima</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">&gt; 10 cm / &gt; 35º</td>
                  <td class="py-2 px-3 text-slate-700">Buena dorsiflexión sin colapso medial del arco</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Flexión de Cadera (ASLR)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romAslrIzq" placeholder="Ej. 75º" class="w-20 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romAslrDer" placeholder="Ej. 78º" class="w-20 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-emerald-700">Simétrica</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">70º - 90º</td>
                  <td class="py-2 px-3 text-slate-700">Excelente disociación lumbopélvica</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Overhead Squat (Sentadilla)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romSquatIzq" placeholder="Limpio" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romSquatDer" placeholder="Limpio" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-emerald-700">Óptimo</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Control lumbopélvico y torácico</td>
                  <td class="py-2 px-3 text-slate-700">Sin pérdida de alineación en flexión profunda</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Lunge Unipodal (Zancada)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romLungeIzq" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romLungeDer" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-emerald-700">Simétrico</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Estabilidad frontal nivel 3</td>
                  <td class="py-2 px-3 text-slate-700">Control de valgo dinámico; proteger menisco derecho</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Hip Hinge (Bisagra de Cadera)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romHingeIzq" placeholder="Excelente" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model="romHingeDer" placeholder="Excelente" class="w-24 h-7 text-center font-semibold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-emerald-700">Óptimo</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Disociación coxofemoral</td>
                  <td class="py-2 px-3 text-slate-700">Patrón biomecánico consolidado para carga axial</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="p-3 bg-[#FAF9F6] border-t border-[#ECE8E1]">
            <label class="block text-[11px] font-bold text-[#6B6D74] uppercase mb-1">Observaciones Biomecánicas Globales:</label>
            <input v-model="romObs" class="w-full h-8 px-3 text-xs bg-white border border-[#CBD5E1] rounded focus:border-[#CEA549] focus:outline-none" />
          </div>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 2: 3. Fuerza Isométrica (FIM)                                -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 2" class="space-y-6">
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
          <h3 class="font-bold text-sm tracking-wide">INVERFIT | FUERZA ISOMÉTRICA MÁXIMA (FIM) & ASIMETRÍAS NEUROMUSCULARES</h3>
          <p class="text-[11px] text-[#8E9097]">CENTRO DE ENTRENAMIENTO — IMTP, LEG EXTENSION, LEG CURL Y ABDUCCIÓN DE CADERA</p>
          <p class="text-[10px] text-[#CEA549] mt-0.5 font-medium">REFERENCIA CIENTÍFICA: VALD Performance, Stone et al., McMaster University (Baremos M-F)</p>
        </div>

        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse min-w-[760px]">
              <thead>
                <tr class="bg-[#24262B] text-white border-b border-[#191A1D]">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[240px]">Test Neuromuscular Isométrico</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[110px]">Izquierda (N)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[110px]">Derecha (N)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[110px]">Asimetría (%)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[130px]">Fuerza Relativa (N/kg)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[210px]">Evaluación & Baremación M-F</th>
                  <th class="py-2.5 px-3 font-bold text-center w-[140px]">Estado Norma M-F</th>
                </tr>
              </thead>

              <tbody class="divide-y divide-[#ECE8E1]">
                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">IMTP (Isometric Mid-Thigh Pull)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimImtpIzq" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimImtpDer" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcAsym(fimImtpIzq, fimImtpDer) != null ? calcAsym(fimImtpIzq, fimImtpDer) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcRel(Math.max(fimImtpIzq || 0, fimImtpDer || 0), pesoKg) != null ? calcRel(Math.max(fimImtpIzq || 0, fimImtpDer || 0), pesoKg) + ' N/kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700 bg-[#FFFF99] border border-amber-300">
                    Fuerza Pico Relativa Óptima (&gt;22 N/kg)
                  </td>
                  <td class="py-2 px-3 text-center bg-[#FFFF99] font-bold text-emerald-800 border border-amber-300">
                    DENTRO DE LA NORMA
                  </td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Leg Extension (Cuádriceps)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimLegExtIzq" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimLegExtDer" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcAsym(fimLegExtIzq, fimLegExtDer) != null ? calcAsym(fimLegExtIzq, fimLegExtDer) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcRel(Math.max(fimLegExtIzq || 0, fimLegExtDer || 0), pesoKg) != null ? calcRel(Math.max(fimLegExtIzq || 0, fimLegExtDer || 0), pesoKg) + ' N/kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700 bg-[#FFFF99] border border-amber-300">
                    Capacidad Extensora de Rodilla Apta
                  </td>
                  <td class="py-2 px-3 text-center bg-[#FFFF99] font-bold text-emerald-800 border border-amber-300">
                    DENTRO DE LA NORMA
                  </td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Leg Curl (Isquiosurales)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimLegCurlIzq" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimLegCurlDer" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcAsym(fimLegCurlIzq, fimLegCurlDer) != null ? calcAsym(fimLegCurlIzq, fimLegCurlDer) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcRel(Math.max(fimLegCurlIzq || 0, fimLegCurlDer || 0), pesoKg) != null ? calcRel(Math.max(fimLegCurlIzq || 0, fimLegCurlDer || 0), pesoKg) + ' N/kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700 bg-[#FFFF99] border border-amber-300">
                    Ratio H:Q Funcional (0.65) Óptimo
                  </td>
                  <td class="py-2 px-3 text-center bg-[#FFFF99] font-bold text-emerald-800 border border-amber-300">
                    DENTRO DE LA NORMA
                  </td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Abducción de Cadera (Glúteo Medio)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimHipAbdIzq" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="fimHipAbdDer" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcAsym(fimHipAbdIzq, fimHipAbdDer) != null ? calcAsym(fimHipAbdIzq, fimHipAbdDer) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcRel(Math.max(fimHipAbdIzq || 0, fimHipAbdDer || 0), pesoKg) != null ? calcRel(Math.max(fimHipAbdIzq || 0, fimHipAbdDer || 0), pesoKg) + ' N/kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700 bg-[#FFFF99] border border-amber-300">
                    Estabilidad Frontal Unipodal Apta
                  </td>
                  <td class="py-2 px-3 text-center bg-[#FFFF99] font-bold text-emerald-800 border border-amber-300">
                    DENTRO DE LA NORMA
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Reference Table Below -->
        <div class="bg-white border border-[#E2DFD8] rounded-md shadow-2xs overflow-hidden">
          <div class="bg-[#191A1D] text-white px-4 py-2 text-xs font-bold tracking-wider border-l-4 border-[#CEA549]">
            BAREMOS DE REFERENCIA CIENTÍFICA — FUERZA ISOMÉTRICA (VALD PERFORMANCE)
          </div>
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-100 text-[#191A1D] border-b border-slate-200 font-bold">
                <th class="py-2 px-3 border-r border-slate-200">Test / Parámetro FIM</th>
                <th class="py-2 px-3 border-r border-slate-200">Bajo Promedio</th>
                <th class="py-2 px-3 border-r border-slate-200">Norma / Saludable</th>
                <th class="py-2 px-3 border-r border-slate-200">Excelente / Atleta</th>
                <th class="py-2 px-3">Referencia</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 text-slate-700">
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">IMTP Pico Fuerza Relativa</td>
                <td class="py-2 px-3 border-r border-slate-200">&lt;30 N/kg (M) / &lt;25 N/kg (F)</td>
                <td class="py-2 px-3 border-r border-slate-200 font-bold text-emerald-700">32 - 42 N/kg (M) / 28 - 36 N/kg (F)</td>
                <td class="py-2 px-3 border-r border-slate-200">&gt;45 N/kg (M) / &gt;38 N/kg (F)</td>
                <td class="py-2 px-3 font-mono text-[11px]">Stone et al. / VALD ForceDecks</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Asimetría Bilateral Inter-Línea</td>
                <td class="py-2 px-3 border-r border-slate-200">&gt; 15.0%</td>
                <td class="py-2 px-3 border-r border-slate-200 font-bold text-emerald-700">5.0% - 10.0%</td>
                <td class="py-2 px-3 border-r border-slate-200">&lt; 5.0% (Simetría Absoluta)</td>
                <td class="py-2 px-3 font-mono text-[11px]">Bishop et al. / Impellizzeri</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Ratio Isquiosural / Cuádriceps (H:Q)</td>
                <td class="py-2 px-3 border-r border-slate-200">&lt; 0.50</td>
                <td class="py-2 px-3 border-r border-slate-200 font-bold text-emerald-700">0.60 - 0.70 (Óptimo)</td>
                <td class="py-2 px-3 border-r border-slate-200">&gt; 0.75</td>
                <td class="py-2 px-3 font-mono text-[11px]">Croisier et al. / Prevención LCA</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Abducción de Cadera Relativa</td>
                <td class="py-2 px-3 border-r border-slate-200">&lt; 3.0 N/kg</td>
                <td class="py-2 px-3 border-r border-slate-200 font-bold text-emerald-700">3.5 - 5.0 N/kg</td>
                <td class="py-2 px-3 border-r border-slate-200">&gt; 5.5 N/kg</td>
                <td class="py-2 px-3 font-mono text-[11px]">Stastny et al. / Estabilidad Valgo</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 3: 4. Saltabilidad Bosco / Batería Adulto Mayor              -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 3" class="space-y-6">
        <!-- Older Adults SPPB -->
        <div v-if="selectedProtocol === 'adulto_mayor'" class="space-y-4">
          <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
            <h3 class="font-bold text-sm tracking-wide">INVERFIT | BATERÍA ADULTO MAYOR (SENIOR FITNESS TEST / SPPB)</h3>
            <p class="text-[11px] text-[#8E9097]">VALORACIÓN FUNCIONAL, PREVENCIÓN DE CAÍDAS Y RIESGO DE SARCOPENIA (RIKLI & JONES / GURALNIK)</p>
          </div>

          <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
            <table class="w-full text-left text-xs border-collapse min-w-[700px]">
              <thead>
                <tr class="bg-[#24262B] text-white">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Test Funcional Senior</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">Resultado</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[90px]">Unidad</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[180px]">Criterio Clínico</th>
                  <th class="py-2.5 px-3 font-bold text-center w-[180px]">Clasificación / Riesgo</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#ECE8E1]">
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">5-STS (Sentarse y Levantarse 5x de la Silla)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="seniorSitToStand" type="number" step="0.1" class="w-20 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">segundos</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">&lt; 12.0 s (Sin riesgo)</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="seniorSitToStand != null ? (seniorSitToStand < 12 ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FCE8E6] text-[#C5221F] font-bold') : 'bg-transparent text-gray-400'" class="px-2 py-0.5 rounded text-[11px] inline-block">
                      {{ seniorSitToStand != null ? (seniorSitToStand < 12 ? 'Normal (Sin Riesgo)' : 'Alerta Riesgo Caída') : 'Por evaluar' }}
                    </span>
                  </td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Timed Up & Go (TUG 3 metros)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="seniorTug" type="number" step="0.1" class="w-20 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">segundos</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600 font-mono">&lt; 10.0 s (Independiente)</td>
                  <td class="py-2 px-3 text-center">
                    <span :class="seniorTug != null ? (seniorTug < 10 ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]') : 'bg-transparent text-gray-400'" class="px-2 py-0.5 rounded text-[11px] inline-block font-semibold">
                      {{ seniorTug != null ? (seniorTug < 10 ? 'Movilidad Óptima' : 'Movilidad Reducida') : 'Por evaluar' }}
                    </span>
                  </td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Chair Stand 30 seg (Reps completas)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="seniorChairStand" type="number" class="w-20 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">repeticiones</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Norma: 12 - 17 reps</td>
                  <td class="py-2 px-3 text-center text-emerald-700 font-bold">Saludable / Funcional</td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Arm Curl 30 seg (Mancuerna 3.6kg H / 2.3kg M)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="seniorArmCurl" type="number" class="w-20 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">repeticiones</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Norma: 14 - 19 reps</td>
                  <td class="py-2 px-3 text-center text-emerald-700 font-bold">Fuerza Brazo Adecuada</td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Marcha 6 Minutos (6MWT)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="senior6mwt" type="number" class="w-20 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center text-slate-500 font-mono">metros</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-slate-600">Referencia: 400 - 550 m</td>
                  <td class="py-2 px-3 text-center text-emerald-700 font-bold">Capacidad Aeróbica Normal</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Athletes & General: Batería Bosco -->
        <div v-else class="space-y-4">
          <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
            <h3 class="font-bold text-sm tracking-wide">INVERFIT | BATERÍA DE SALTABILIDAD BOSCO (PLATAFORMA / ERGOJUMP)</h3>
            <p class="text-[11px] text-[#8E9097]">VALORACIÓN DE POTENCIA MECÁNICA, ÍNDICE ELÁSTICO (IE) Y COORDINACIÓN DE BRAZOS (CARMELO BOSCO)</p>
          </div>

          <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
            <table class="w-full text-left text-xs border-collapse min-w-[700px]">
              <thead>
                <tr class="bg-[#24262B] text-white">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Test de Salto</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">Altura (cm)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[130px]">Tiempo Vuelo</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[180px]">Potencia / Índice Elástico</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[160px]">Clasificación</th>
                  <th class="py-2.5 px-3 font-bold">Baremos Bosco</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#ECE8E1]">
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Squat Jump (SJ - Fuerza Contráctil Pura)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="boscoSjCm" type="number" step="0.1" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ boscoSjFlightMs != null ? boscoSjFlightMs + ' ms' : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-semibold text-slate-700">Potencia Media</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-bold text-emerald-700">Rango Normal</td>
                  <td class="py-2 px-3 text-slate-600 font-mono">28 - 38 cm</td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Countermovement Jump (CMJ - Ciclo CEA)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="boscoCmjCm" type="number" step="0.1" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ boscoCmjFlightMs != null ? boscoCmjFlightMs + ' ms' : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ boscoIndiceElastico != null ? 'IE: ' + (boscoIndiceElastico >= 0 ? '+' : '') + boscoIndiceElastico + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-bold text-emerald-700">Reutilización Óptima</td>
                  <td class="py-2 px-3 text-slate-600 font-mono">34 - 44 cm (&gt;10% sobre SJ)</td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Abalakov (CMJa - Con Brazos)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="boscoAbalakovCm" type="number" step="0.1" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ boscoAbalakovFlightMs != null ? boscoAbalakovFlightMs + ' ms' : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ boscoCoordBrazos != null ? 'Aporte: ' + (boscoCoordBrazos >= 0 ? '+' : '') + boscoCoordBrazos + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-bold text-emerald-700">Segmentario Normal</td>
                  <td class="py-2 px-3 text-slate-600 font-mono">+10% a +15% sobre CMJ</td>
                </tr>
                <tr>
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Drop Jump (DJ - Reactividad / RSI)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="boscoDjCm" type="number" step="0.1" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">Contacto: {{ boscoDjContactMs != null ? boscoDjContactMs + ' ms' : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ boscoRsi != null ? 'RSI: ' + boscoRsi : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-bold text-emerald-700">Reactividad Alta</td>
                  <td class="py-2 px-3 text-slate-600 font-mono">RSI &gt; 1.50 (Flanagan)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 4: 5. Perfil F-V (Encoder ADR)                               -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 4" class="space-y-6">
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
          <h3 class="font-bold text-sm tracking-wide">INVERFIT | PERFIL FUERZA-VELOCIDAD & 1RM ESTIMADO (ENCODER LINEAL ADR)</h3>
          <p class="text-[11px] text-[#8E9097]">CENTRO DE ENTRENAMIENTO — VELOCIDAD CONCÉNTRICA, CARGA Y ZONAS DE ENTRENAMIENTO</p>
          <p class="text-[10px] text-[#CEA549] mt-0.5 font-medium">REFERENCIAS CIENTÍFICAS: González-Badillo & Sánchez-Medina (2010), Amador García-Ramos, Verkhoshansky</p>
        </div>

        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse min-w-[800px]">
              <thead>
                <tr class="bg-[#24262B] text-white border-b border-[#191A1D]">
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Ejercicio Evaluado</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">Carga (kg)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[120px]">VMP (m/s)</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[80px]">RPE</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[130px]">1RM Estimado</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center w-[100px]">% 1RM</th>
                  <th class="py-2.5 px-3 font-bold border-r border-[#34373F] w-[200px]">Zona F-V Automática</th>
                  <th class="py-2.5 px-3 font-bold">Referencia Científica</th>
                </tr>
              </thead>

              <tbody class="divide-y divide-[#ECE8E1]">
                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">SENTADILLA LIBRE (Encoder ADR)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrSquatLoad" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrSquatVmp" type="number" step="0.01" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ adrSquatRpe != null ? adrSquatRpe : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calc1RM('squat', adrSquatLoad, adrSquatVmp) != null ? calc1RM('squat', adrSquatLoad, adrSquatVmp) + ' kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcPct1RM('squat', adrSquatVmp) != null ? calcPct1RM('squat', adrSquatVmp) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">
                    {{ getZone(adrSquatVmp) }}
                  </td>
                  <td class="py-2 px-3 text-[11px] text-slate-600 font-mono">González-Badillo (2010)</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">PRESS BANCA (Encoder ADR)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrBenchLoad" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrBenchVmp" type="number" step="0.01" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ adrBenchRpe != null ? adrBenchRpe : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calc1RM('bench', adrBenchLoad, adrBenchVmp) != null ? calc1RM('bench', adrBenchLoad, adrBenchVmp) + ' kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcPct1RM('bench', adrBenchVmp) != null ? calcPct1RM('bench', adrBenchVmp) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">
                    {{ getZone(adrBenchVmp) }}
                  </td>
                  <td class="py-2 px-3 text-[11px] text-slate-600 font-mono">González-Badillo (2010)</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">PESO MUERTO / HINGE (ADR)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrDeadliftLoad" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrDeadliftVmp" type="number" step="0.01" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ adrDeadliftRpe != null ? adrDeadliftRpe : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calc1RM('squat', adrDeadliftLoad, adrDeadliftVmp) != null ? calc1RM('squat', adrDeadliftLoad, adrDeadliftVmp) + ' kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcPct1RM('squat', adrDeadliftVmp) != null ? calcPct1RM('squat', adrDeadliftVmp) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">
                    {{ getZone(adrDeadliftVmp) }}
                  </td>
                  <td class="py-2 px-3 text-[11px] text-slate-600 font-mono">González-Badillo (2010)</td>
                </tr>

                <tr class="hover:bg-slate-50">
                  <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">DOMINADAS LASTRADAS (ADR)</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrPullupLoad" type="number" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">
                    <input v-model.number="adrPullupVmp" type="number" step="0.01" class="w-18 h-7 text-center font-bold border border-[#CBD5E1] rounded text-xs" />
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">{{ adrPullupRpe != null ? adrPullupRpe : '—' }}</td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calc1RM('bench', adrPullupLoad, adrPullupVmp) != null ? calc1RM('bench', adrPullupLoad, adrPullupVmp) + ' kg' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                    {{ calcPct1RM('bench', adrPullupVmp) != null ? calcPct1RM('bench', adrPullupVmp) + '%' : '—' }}
                  </td>
                  <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">
                    {{ getZone(adrPullupVmp) }}
                  </td>
                  <td class="py-2 px-3 text-[11px] text-slate-600 font-mono">González-Badillo (2010)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Reference Table Below -->
        <div class="bg-white border border-[#E2DFD8] rounded-md shadow-2xs overflow-hidden">
          <div class="bg-[#191A1D] text-white px-4 py-2 text-xs font-bold tracking-wider border-l-4 border-[#CEA549]">
            TABLA DE REFERENCIA CIENTÍFICA — ZONAS DE ENTRENAMIENTO SEGÚN VELOCIDAD CONCÉNTRICA
          </div>
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-100 text-[#191A1D] border-b border-slate-200 font-bold">
                <th class="py-2 px-3 border-r border-slate-200">Zona de Entrenamiento</th>
                <th class="py-2 px-3 border-r border-slate-200">Velocidad Conc. (m/s)</th>
                <th class="py-2 px-3 border-r border-slate-200">% 1RM Estimado</th>
                <th class="py-2 px-3 border-r border-slate-200">Objetivo Neuromuscular</th>
                <th class="py-2 px-3">Autor / Referencia</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 text-slate-700">
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Velocidad / Explosivo</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">&gt; 1.00 m/s</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">30% - 50%</td>
                <td class="py-2 px-3 border-r border-slate-200">Desarrollo de tasa de aplicación de fuerza (RFD)</td>
                <td class="py-2 px-3 font-mono text-[11px]">Verkhoshansky / García-Ramos</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Potencia Máxima</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">0.75 - 1.00 m/s</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">50% - 70%</td>
                <td class="py-2 px-3 border-r border-slate-200">Optimización del pico de potencia mecánica (Watt/kg)</td>
                <td class="py-2 px-3 font-mono text-[11px]">Bosco / González-Badillo</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Fuerza - Hipertrofia</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">0.50 - 0.75 m/s</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">70% - 85%</td>
                <td class="py-2 px-3 border-r border-slate-200">Señalización hipertrófica y fuerza submáxima</td>
                <td class="py-2 px-3 font-mono text-[11px]">Schoenfeld / Badillo</td>
              </tr>
              <tr>
                <td class="py-2 px-3 font-semibold border-r border-slate-200">Fuerza Máxima</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">&lt; 0.50 m/s</td>
                <td class="py-2 px-3 border-r border-slate-200 font-mono">&gt; 85% - 100%</td>
                <td class="py-2 px-3 border-r border-slate-200">Reclutamiento de unidades motoras de alto umbral</td>
                <td class="py-2 px-3 font-mono text-[11px]">Zatsiorsky / Kraemer</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 5: Referencias y Baremos M-F                                 -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 5" class="space-y-6">
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
          <h3 class="font-bold text-sm tracking-wide">INVERFIT | COMPENDIO CIENTÍFICO Y NORMAS DE REFERENCIA</h3>
          <p class="text-[11px] text-[#8E9097]">MARCO BIBLIOGRÁFICO DE REFERENCIA (OMS, ACSM, EWGSOP2, VALD, BOSCO, BADILLO)</p>
        </div>

        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-[#24262B] text-white">
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Componente Evaluado</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Organización / Autor</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Criterios Hombres (M)</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Criterios Mujeres (F)</th>
                <th class="py-2.5 px-3 font-bold">Aplicación Práctica en INVERFIT</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#ECE8E1] text-slate-700">
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Índice de Masa Corporal (IMC)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">OMS (2024)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Normopeso: 18.5 - 24.9 kg/m²</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Normopeso: 18.5 - 24.9 kg/m²</td>
                <td class="py-2 px-3">Clasificación metabólica trimestral</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Presión Arterial en Reposo</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">AHA / ACC (2017)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Normal: &lt; 120 / 80 mmHg</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Normal: &lt; 120 / 80 mmHg</td>
                <td class="py-2 px-3">Filtro de seguridad previo al esfuerzo</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Porcentaje de Grasa Corporal</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">ACSM (11th Ed.)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Saludable (20-39a): 11.0% - 21.0%</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Saludable (20-39a): 15.0% - 24.0%</td>
                <td class="py-2 px-3">Composición corporal y salud metabólica</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Fuerza Prensil (Dinamometría)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">EWGSOP2 / Bohannon</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">Normal / Saludable: &gt; 30 kg</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-emerald-800">Normal / Saludable: &gt; 20 kg</td>
                <td class="py-2 px-3">Marcador de fuerza global y sarcopenia</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Perfil Fuerza-Velocidad (1RM)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">González-Badillo (2010)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Velocidad teórica 1RM: ~0.30 m/s</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">Velocidad teórica 1RM: ~0.30 m/s</td>
                <td class="py-2 px-3">Control de cargas y zonas de potencia</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Batería de Saltabilidad (SJ/CMJ/DJ)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">Carmelo Bosco (1983)</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">CMJ Normal: 29 - 38 cm</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">CMJ Normal: 24 - 31 cm</td>
                <td class="py-2 px-3">Cuantificación de capacidad reactiva</td>
              </tr>
              <tr>
                <td class="py-2.5 px-3 font-bold text-[#191A1D] border-r border-[#ECE8E1]">Fuerza Isométrica IMTP & FIM</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-600">VALD Performance / Stone</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">IMTP Relativo: 32 - 42 N/kg</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1]">IMTP Relativo: 28 - 36 N/kg</td>
                <td class="py-2 px-3">Evaluación neuromuscular y asimetrías</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ------------------------------------------------------------------ -->
      <!-- SHEET 6: Seguimiento Histórico                                     -->
      <!-- ------------------------------------------------------------------ -->
      <div v-else-if="activeTab === 6" class="space-y-6">
        <div class="bg-[#191A1D] text-white p-3 rounded-t-md border-l-4 border-[#CEA549]">
          <h3 class="font-bold text-sm tracking-wide">INVERFIT | REGISTRO EVOLUTIVO TRIMESTRAL</h3>
          <p class="text-[11px] text-[#8E9097]">CENTRO DE ENTRENAMIENTO — SEGUIMIENTO TRIMESTRAL & PROTOCOLO CLÍNICO-DEPORTIVO (OMS / ACSM)</p>
          <p class="text-[10px] text-[#CEA549] mt-0.5 font-medium">CONTROL HISTÓRICO DE EVALUACIONES (COMPARATIVA CADA 3 MESES)</p>
        </div>

        <div class="bg-white border border-[#E2DFD8] rounded-b-md shadow-2xs overflow-hidden">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-[#24262B] text-white">
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">ID Alumno</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Nombre Completo</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center">Sexo (M/F)</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center">Fecha Eval</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center">Peso (kg)</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center">Talla (m)</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F] text-center">IMC (kg/m²)</th>
                <th class="py-2.5 px-3 font-bold border-r border-[#34373F]">Clasificación IMC</th>
                <th class="py-2.5 px-3 font-bold text-center">Estado / Tendencia</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#ECE8E1]">
              <tr v-for="ev in gymStore.evaluations.filter(e => e.status === 'Registrada')" :key="ev.id" class="hover:bg-slate-50">
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-mono text-slate-500">{{ ev.id }}</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-bold text-[#191A1D]">{{ ev.studentName }}</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-mono">M</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">{{ ev.date }}</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] text-center font-bold">{{ ev.anthropometry?.weightKg || '78.5' }}</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] text-center">{{ ev.anthropometry?.heightCm ? (ev.anthropometry.heightCm/100).toFixed(2) : '1.76' }}</td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] text-center bg-[#FFFF99] font-black border border-amber-300">
                  {{ ev.anthropometry?.imc || '—' }}
                </td>
                <td class="py-2 px-3 border-r border-[#ECE8E1] font-semibold text-slate-700">{{ ev.anthropometry?.imcClassification || '—' }}</td>
                <td class="py-2 px-3 text-center">
                  <span class="px-2 py-0.5 rounded text-[11px] bg-[#E6F4EA] text-[#137333] font-semibold">
                    Óptimo / Saludable
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- 4. BOTTOM STATUS BAR                                                 -->
    <!-- ==================================================================== -->
    <div class="bg-[#191A1D] text-[#8E9097] px-4 py-2.5 border-t border-[#282A30] text-xs flex flex-col sm:flex-row sm:items-center justify-between gap-2">
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-[#CEA549]"></span>
        <span class="text-white font-semibold">Evaluación Oficial Activa</span>
        <span class="text-[#3A3D46]">|</span>
        <span class="text-[#8E9097]">Autoguardado en Almacén Local / API Go REST</span>
      </div>

      <div class="flex items-center gap-3">
        <span class="text-[11px] text-[#CEA549] font-mono font-bold">
          {{ sheetTabs[activeTab].shortName }}
        </span>
        <button
          @click="handleSaveWorkbook"
          :disabled="isSaving"
          class="px-3 py-1 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded text-xs font-bold transition-colors cursor-pointer shadow-xs"
        >
          {{ isSaving ? 'Guardando...' : 'Guardar Ficha' }}
        </button>
      </div>
    </div>
  </div>
</template>
