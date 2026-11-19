<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { FrequencyType, StudentOnboarding } from '@/types'
import {
  User,
  Heart,
  Activity,
  FileCheck2,
  CheckCircle2,
  AlertTriangle,
  RotateCcw,
  Send,
  ArrowRight,
  ArrowLeft,
  Phone,
  ShieldCheck,
  Check
} from 'lucide-vue-next'

const gymStore = useGymStore()

// Wizard state
const currentStep = ref(1)
const totalSteps = 4
const isSubmitting = ref(false)
const submitSuccess = ref(false)
const errorMessage = ref('')

// Step 1: Identificación y Contacto
const form = ref({
  // Identificación
  fullName: '',
  rut: '',
  birthDate: '',
  age: '',
  gender: 'Femenino' as 'Femenino' | 'Masculino' | 'Otro',
  phone: '',
  email: '',
  comuna: 'Concepción',
  occupation: '',
  hasChildren: 'No' as 'Si' | 'No',
  childrenCount: '',
  emergencyContact: '',

  // Perfil Deportivo y Objetivos
  primaryGoal: 'Salud, funcionalidad y longevidad' as 'Salud, funcionalidad y longevidad' | 'Estética' | 'Rendimiento deportivo' | 'Alto Rendimiento' | 'Otros',
  goalDetails: '',
  hasTrainedBefore: 'Si' as 'Si' | 'No',
  trainingExperience: '1-6 meses' as 'Sin experiencia previa' | '1-6 meses' | '6-12 meses' | 'Más de 1 año',
  sportsPracticed: '',
  frequency: '2 veces/sem.' as FrequencyType,

  // Ficha Clínica (PAR-Q+)
  chronicDiseases: '',
  surgeriesInjuries: '',
  injuriesPreventExercise: 'No' as 'Si' | 'No',
  currentDifficulties: '',
  hadCovid: 'No' as 'Si' | 'No',
  covidSequelae: '',
  frequentConsumption: '',

  // Consentimiento
  consentTermsRead: false,
  consentTruthDeclared: false
})

// Auto-calculate age from birthDate
function handleBirthDateChange() {
  if (!form.value.birthDate) return
  const birth = new Date(form.value.birthDate)
  const now = new Date()
  let age = now.getFullYear() - birth.getFullYear()
  const m = now.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
    age--
  }
  if (age > 0 && age < 120) {
    form.value.age = String(age)
  }
}

// Digital Signature Canvas
const signatureCanvas = ref<HTMLCanvasElement | null>(null)
const isDrawing = ref(false)
const hasSignature = ref(false)
const signatureDataUrl = ref('')
const signatureTimestamp = ref('')

function setupCanvas() {
  const canvas = signatureCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Set physical vs display resolution for retina/high DPI
  const rect = canvas.getBoundingClientRect()
  const dpr = window.devicePixelRatio || 1
  canvas.width = rect.width * dpr
  canvas.height = rect.height * dpr
  ctx.scale(dpr, dpr)

  ctx.lineWidth = 2.5
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.strokeStyle = '#191A1D'

  // If there was a previous signature saved, redraw
  if (signatureDataUrl.value) {
    const img = new Image()
    img.onload = () => {
      ctx.drawImage(img, 0, 0, rect.width, rect.height)
    }
    img.src = signatureDataUrl.value
  }
}

function getPointerPos(e: MouseEvent | TouchEvent) {
  const canvas = signatureCanvas.value
  if (!canvas) return { x: 0, y: 0 }
  const rect = canvas.getBoundingClientRect()
  if ('touches' in e && e.touches.length > 0) {
    return {
      x: e.touches[0].clientX - rect.left,
      y: e.touches[0].clientY - rect.top
    }
  } else if ('clientX' in e) {
    return {
      x: (e as MouseEvent).clientX - rect.left,
      y: (e as MouseEvent).clientY - rect.top
    }
  }
  return { x: 0, y: 0 }
}

function startDrawing(e: MouseEvent | TouchEvent) {
  const canvas = signatureCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  isDrawing.value = true
  const pos = getPointerPos(e)
  ctx.beginPath()
  ctx.moveTo(pos.x, pos.y)
}

function draw(e: MouseEvent | TouchEvent) {
  if (!isDrawing.value) return
  const canvas = signatureCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const pos = getPointerPos(e)
  ctx.lineTo(pos.x, pos.y)
  ctx.stroke()
  hasSignature.value = true
}

function stopDrawing() {
  if (!isDrawing.value) return
  isDrawing.value = false
  const canvas = signatureCanvas.value
  if (!canvas) return

  signatureDataUrl.value = canvas.toDataURL('image/png')
  signatureTimestamp.value = new Date().toLocaleString('es-CL', {
    dateStyle: 'full',
    timeStyle: 'medium'
  })
}

function clearSignature() {
  const canvas = signatureCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.width, canvas.height)
  hasSignature.value = false
  signatureDataUrl.value = ''
  signatureTimestamp.value = ''
}

// Validation per step
function validateStep(step: number): boolean {
  errorMessage.value = ''
  if (step === 1) {
    if (!form.value.fullName.trim()) {
      errorMessage.value = 'Por favor ingresa tu nombre completo.'
      return false
    }
    if (!form.value.rut.trim()) {
      errorMessage.value = 'Por favor ingresa tu RUT.'
      return false
    }
    if (!form.value.phone.trim()) {
      errorMessage.value = 'Por favor ingresa tu teléfono o WhatsApp de contacto.'
      return false
    }
    if (!form.value.email.trim() || !form.value.email.includes('@')) {
      errorMessage.value = 'Por favor ingresa un correo electrónico válido.'
      return false
    }
    if (!form.value.emergencyContact.trim()) {
      errorMessage.value = 'Por favor indica un contacto de emergencia (Nombre y teléfono).'
      return false
    }
    return true
  }

  if (step === 2) {
    if (!form.value.primaryGoal) {
      errorMessage.value = 'Por favor selecciona tu objetivo principal.'
      return false
    }
    return true
  }

  if (step === 3) {
    // Par-q fields are informative, require at least confirmation or details
    return true
  }

  if (step === 4) {
    if (!form.value.consentTermsRead || !form.value.consentTruthDeclared) {
      errorMessage.value = 'Debes aceptar las cláusulas de consentimiento y declaración jurada.'
      return false
    }
    if (!hasSignature.value || !signatureDataUrl.value) {
      errorMessage.value = 'Por favor estampa tu firma digital en el recuadro antes de enviar.'
      return false
    }
    return true
  }

  return true
}

function nextStep() {
  if (validateStep(currentStep.value)) {
    if (currentStep.value < totalSteps) {
      currentStep.value++
      window.scrollTo({ top: 0, behavior: 'smooth' })
      if (currentStep.value === 4) {
        nextTick(() => {
          setupCanvas()
        })
      }
    }
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

// Final Submit
async function handleSubmit() {
  if (!validateStep(4)) return

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const onboardingPayload: StudentOnboarding = {
      rut: form.value.rut.trim(),
      birthDate: form.value.birthDate,
      age: form.value.age ? Number(form.value.age) : undefined,
      gender: form.value.gender,
      comuna: form.value.comuna.trim(),
      occupation: form.value.occupation.trim(),
      hasChildren: form.value.hasChildren,
      childrenCount: form.value.childrenCount.trim(),

      primaryGoal: form.value.primaryGoal,
      goalDetails: form.value.goalDetails.trim(),
      hasTrainedBefore: form.value.hasTrainedBefore,
      trainingExperience: form.value.trainingExperience,
      sportsPracticed: form.value.sportsPracticed.trim(),
      availableDaysWeekly: form.value.frequency === '4 veces/sem.' ? '4-6 veces a la semana' : '2-4 veces a la semana',

      chronicDiseases: form.value.chronicDiseases.trim(),
      surgeriesInjuries: form.value.surgeriesInjuries.trim(),
      injuriesPreventExercise: form.value.injuriesPreventExercise,
      currentDifficulties: form.value.currentDifficulties.trim(),
      hadCovid: form.value.hadCovid,
      covidSequelae: form.value.covidSequelae.trim(),
      frequentConsumption: form.value.frequentConsumption.trim(),

      consentAccepted: true,
      consentTimestamp: signatureTimestamp.value || new Date().toISOString(),
      digitalSignature: signatureDataUrl.value,
      termsVersion: 'v1.0 (Noviembre 2026 - Oficial Pablo)'
    }

    // Build synthesized health conditions
    const healthParts: string[] = []
    if (form.value.chronicDiseases.trim()) healthParts.push(`Enfermedades: ${form.value.chronicDiseases.trim()}`)
    if (form.value.surgeriesInjuries.trim()) healthParts.push(`Cirugías/Lesiones: ${form.value.surgeriesInjuries.trim()}`)
    if (form.value.injuriesPreventExercise === 'Si') healthParts.push('Presenta dolor/limitación física')
    if (form.value.hadCovid === 'Si' && form.value.covidSequelae.trim()) healthParts.push(`Secuelas COVID: ${form.value.covidSequelae.trim()}`)

    await gymStore.createStudent({
      fullName: form.value.fullName.trim(),
      email: form.value.email.trim(),
      phone: form.value.phone.trim(),
      frequency: form.value.frequency,
      membershipStatus: 'Por revisar', // Mark as pending coach approval
      fixedSchedule: 'Por coordinar con Coach Pablo',
      emergencyContact: form.value.emergencyContact.trim(),
      healthConditions: healthParts.join(' | ') || 'Sin patologías ni lesiones declaradas',
      notes: `Objetivo: ${form.value.primaryGoal}. ${form.value.goalDetails ? 'Detalle: ' + form.value.goalDetails : ''}`,
      onboarding: onboardingPayload
    })

    submitSuccess.value = true
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Ocurrió un error al enviar el formulario. Por favor intenta nuevamente.'
  } finally {
    isSubmitting.value = false
  }
}

// Resize canvas on orientation change or window resize
function handleResize() {
  if (currentStep.value === 4 && signatureCanvas.value) {
    setupCanvas()
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

const whatsappMessageUrl = computed(() => {
  const phone = '56987654321' // Official Coach phone placeholder or direct
  const text = encodeURIComponent(
    `Hola Coach Pablo! Acabo de completar mi Ficha de Ingreso y Consentimiento en INVERFIT para ${form.value.fullName}. Quedo atento a tus indicaciones y a la confirmación de mis horarios.`
  )
  return `https://wa.me/${phone}?text=${text}`
})
</script>

<template>
  <div class="min-h-screen bg-[#F5F3EE] flex flex-col justify-between">
    <!-- Top Branded Header -->
    <header class="bg-[#191A1D] border-b border-[#2C2E33] text-white py-5 px-4 sm:px-8 shadow-md">
      <div class="max-w-3xl mx-auto flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img src="/brand/Logos_INVERFIT-02.png" alt="INVERFIT" class="h-9 sm:h-11 object-contain" />
          <div class="h-7 w-[1px] bg-white/20 hidden sm:block"></div>
          <div class="hidden sm:block">
            <p class="text-xs font-bold text-[#CEA549] uppercase tracking-wider">Entrenamiento Personalizado</p>
            <p class="text-[11px] text-[#A6A8B0]">Ficha de Ingreso y Consentimiento</p>
          </div>
        </div>

        <div class="flex items-center gap-2 text-xs text-[#A6A8B0] bg-[#222429] px-3 py-1.5 rounded-full border border-white/10">
          <ShieldCheck class="w-4 h-4 text-[#CEA549]" />
          <span class="font-medium text-white">Registro Oficial</span>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <main class="flex-1 max-w-3xl w-full mx-auto p-4 sm:p-6 lg:p-8">
      <!-- SUCCESS SCREEN -->
      <div v-if="submitSuccess" class="bg-white rounded-2xl border border-[#E2DFD8] p-6 sm:p-10 shadow-lg text-center space-y-6 animate-in fade-in zoom-in-95 duration-200">
        <div class="w-16 h-16 sm:w-20 sm:h-20 mx-auto rounded-full bg-[#E6F4EA] border-2 border-[#137333]/30 flex items-center justify-center text-[#137333] shadow-inner">
          <CheckCircle2 class="w-10 h-10 sm:w-12 sm:h-12" />
        </div>

        <div class="space-y-2">
          <span class="inline-block px-3 py-1 bg-[#CEA549]/15 text-[#8C6B1B] border border-[#CEA549]/30 rounded-full text-xs font-bold uppercase tracking-wider">
            ¡Ficha Registrada con Éxito!
          </span>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-[#191A1D]">
            ¡Bienvenido(a) a INVERFIT, {{ form.fullName }}!
          </h1>
          <p class="text-sm sm:text-base text-[#6B6D74] max-w-xl mx-auto leading-relaxed">
            Tus datos de salud, objetivos deportivos y consentimiento informado firmado digitalmente han quedado registrados en nuestro sistema.
          </p>
        </div>

        <!-- Summary box -->
        <div class="bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl p-5 text-left max-w-lg mx-auto text-xs sm:text-sm space-y-3">
          <div class="flex justify-between items-center border-b border-[#ECE8E1] pb-2">
            <span class="text-[#6B6D74]">RUT Alumno:</span>
            <span class="font-bold text-[#191A1D]">{{ form.rut }}</span>
          </div>
          <div class="flex justify-between items-center border-b border-[#ECE8E1] pb-2">
            <span class="text-[#6B6D74]">Objetivo Principal:</span>
            <span class="font-bold text-[#191A1D]">{{ form.primaryGoal }}</span>
          </div>
          <div class="flex justify-between items-center border-b border-[#ECE8E1] pb-2">
            <span class="text-[#6B6D74]">Frecuencia Solicitada:</span>
            <span class="font-bold text-[#191A1D]">{{ form.frequency }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-[#6B6D74]">Consentimiento & Firma:</span>
            <span class="font-bold text-[#137333] flex items-center gap-1">
              <Check class="w-4 h-4" /> Firmado digitalmente
            </span>
          </div>
        </div>

        <div class="p-4 bg-[#FEF7E0]/60 border border-[#F0D878] rounded-xl text-left max-w-lg mx-auto text-xs text-[#8C6B1B] space-y-1">
          <p class="font-bold text-[#70520F]">¿Cuál es el siguiente paso?</p>
          <p class="text-[#8C6B1B]">
            El Coach Pablo revisará tus antecedentes para afinar tu planificación y se comunicará contigo vía WhatsApp para confirmar tu bloque de entrenamiento (máximo 3 alumnos por turno).
          </p>
        </div>

        <div class="pt-4 flex flex-col sm:flex-row items-center justify-center gap-3 max-w-md mx-auto">
          <a
            :href="whatsappMessageUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="w-full sm:w-auto flex-1 h-12 px-6 bg-[#25D366] hover:bg-[#20ba59] text-white rounded-xl text-sm font-bold flex items-center justify-center gap-2 shadow-md hover:shadow-lg transition-all"
          >
            <Phone class="w-4 h-4" />
            <span>Avisar a Pablo por WhatsApp</span>
          </a>
        </div>
      </div>

      <!-- ONBOARDING FORM WIZARD -->
      <div v-else class="bg-white rounded-2xl border border-[#E2DFD8] shadow-sm overflow-hidden">
        <!-- Progress Stepper Header -->
        <div class="bg-[#FAF9F6] border-b border-[#E2DFD8] p-4 sm:p-6">
          <div class="flex items-center justify-between mb-3 text-xs sm:text-sm font-bold text-[#191A1D]">
            <span>Paso {{ currentStep }} de {{ totalSteps }}</span>
            <span class="text-[#8C6B1B] font-semibold">{{ Math.round((currentStep / totalSteps) * 100) }}% completado</span>
          </div>

          <!-- Progress Bar -->
          <div class="w-full bg-[#E2DFD8] h-2 rounded-full overflow-hidden">
            <div
              class="bg-[#CEA549] h-full transition-all duration-300 ease-out"
              :style="{ width: `${(currentStep / totalSteps) * 100}%` }"
            ></div>
          </div>

          <!-- Step Indicators -->
          <div class="grid grid-cols-4 gap-2 mt-4 text-center">
            <button
              type="button"
              @click="currentStep > 1 && (currentStep = 1)"
              :class="currentStep === 1 ? 'text-[#8C6B1B] font-bold' : currentStep > 1 ? 'text-[#191A1D] font-medium' : 'text-[#A6A8B0]'"
              class="text-[11px] sm:text-xs flex flex-col items-center gap-1 cursor-pointer"
            >
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center text-xs"
                :class="currentStep === 1 ? 'bg-[#CEA549] text-black font-bold' : currentStep > 1 ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#E2DFD8] text-[#6B6D74]'"
              >
                1
              </div>
              <span class="hidden sm:inline">Datos Personales</span>
              <span class="sm:hidden">Datos</span>
            </button>

            <button
              type="button"
              @click="currentStep > 2 && (currentStep = 2)"
              :class="currentStep === 2 ? 'text-[#8C6B1B] font-bold' : currentStep > 2 ? 'text-[#191A1D] font-medium' : 'text-[#A6A8B0]'"
              class="text-[11px] sm:text-xs flex flex-col items-center gap-1 cursor-pointer"
            >
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center text-xs"
                :class="currentStep === 2 ? 'bg-[#CEA549] text-black font-bold' : currentStep > 2 ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#E2DFD8] text-[#6B6D74]'"
              >
                2
              </div>
              <span class="hidden sm:inline">Objetivos</span>
              <span class="sm:hidden">Metas</span>
            </button>

            <button
              type="button"
              @click="currentStep > 3 && (currentStep = 3)"
              :class="currentStep === 3 ? 'text-[#8C6B1B] font-bold' : currentStep > 3 ? 'text-[#191A1D] font-medium' : 'text-[#A6A8B0]'"
              class="text-[11px] sm:text-xs flex flex-col items-center gap-1 cursor-pointer"
            >
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center text-xs"
                :class="currentStep === 3 ? 'bg-[#CEA549] text-black font-bold' : currentStep > 3 ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#E2DFD8] text-[#6B6D74]'"
              >
                3
              </div>
              <span class="hidden sm:inline">Salud (PAR-Q+)</span>
              <span class="sm:hidden">Salud</span>
            </button>

            <button
              type="button"
              :class="currentStep === 4 ? 'text-[#8C6B1B] font-bold' : 'text-[#A6A8B0]'"
              class="text-[11px] sm:text-xs flex flex-col items-center gap-1"
            >
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center text-xs"
                :class="currentStep === 4 ? 'bg-[#CEA549] text-black font-bold' : 'bg-[#E2DFD8] text-[#6B6D74]'"
              >
                4
              </div>
              <span class="hidden sm:inline">Firma & Legal</span>
              <span class="sm:hidden">Firma</span>
            </button>
          </div>
        </div>

        <!-- Form Body -->
        <div class="p-6 sm:p-8 space-y-6">
          <!-- Error alert -->
          <div v-if="errorMessage" class="p-3.5 bg-rose-50 border border-rose-200 text-rose-800 rounded-xl text-xs sm:text-sm flex items-start gap-2.5 animate-shake">
            <AlertTriangle class="w-4 h-4 flex-shrink-0 mt-0.5 text-rose-600" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- STEP 1: DATOS PERSONALES -->
          <div v-if="currentStep === 1" class="space-y-5 animate-in fade-in duration-150">
            <div>
              <h2 class="text-lg sm:text-xl font-bold text-[#191A1D] flex items-center gap-2">
                <User class="w-5 h-5 text-[#CEA549]" />
                Paso 1: Identificación y Contacto
              </h2>
              <p class="text-xs sm:text-sm text-[#6B6D74] mt-1">
                Tus datos básicos para abrir tu expediente deportivo oficial en INVERFIT.
              </p>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- Nombre -->
              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Nombre Completo <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.fullName"
                  type="text"
                  placeholder="Ej: Nelson López Valenzuela"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-2 focus:ring-[#CEA549]/20 transition-all bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <!-- RUT -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  RUT <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.rut"
                  type="text"
                  placeholder="Ej: 19.876.543-2"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-2 focus:ring-[#CEA549]/20 transition-all bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <!-- Fecha Nacimiento / Edad -->
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                    Nacimiento
                  </label>
                  <input
                    v-model="form.birthDate"
                    type="date"
                    @change="handleBirthDateChange"
                    class="w-full px-3 py-2.5 rounded-lg border border-[#DCD8CF] text-xs text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  />
                </div>
                <div>
                  <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                    Edad
                  </label>
                  <input
                    v-model="form.age"
                    type="number"
                    placeholder="Años"
                    class="w-full px-3 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  />
                </div>
              </div>

              <!-- Teléfono / WhatsApp -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Teléfono / WhatsApp <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.phone"
                  type="tel"
                  placeholder="+56 9 1234 5678"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-2 focus:ring-[#CEA549]/20 transition-all bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <!-- Correo -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Correo Electrónico <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.email"
                  type="email"
                  placeholder="ejemplo@correo.cl"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-2 focus:ring-[#CEA549]/20 transition-all bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <!-- Género -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Género
                </label>
                <select
                  v-model="form.gender"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-white"
                >
                  <option value="Femenino">Femenino</option>
                  <option value="Masculino">Masculino</option>
                  <option value="Otro">Otro / Prefiero no decir</option>
                </select>
              </div>

              <!-- Comuna -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Comuna de Residencia
                </label>
                <input
                  v-model="form.comuna"
                  type="text"
                  placeholder="Ej: Concepción, San Pedro de la Paz, Talcahuano"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>

              <!-- Ocupación -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Ocupación o Actividad
                </label>
                <input
                  v-model="form.occupation"
                  type="text"
                  placeholder="Ej: Profesional, Estudiante, Independiente..."
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>

              <!-- Hijos -->
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Tienes hijos?
                </label>
                <div class="flex items-center gap-4 py-2">
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="No" v-model="form.hasChildren" class="accent-[#CEA549]" /> No
                  </label>
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="Si" v-model="form.hasChildren" class="accent-[#CEA549]" /> Sí
                  </label>
                  <input
                    v-if="form.hasChildren === 'Si'"
                    v-model="form.childrenCount"
                    type="text"
                    placeholder="¿Cuántos?"
                    class="w-24 px-2 py-1 border border-[#DCD8CF] rounded text-xs"
                  />
                </div>
              </div>

              <!-- Contacto de Emergencia -->
              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Contacto de Emergencia (Nombre, Parentesco y Teléfono) <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.emergencyContact"
                  type="text"
                  placeholder="Ej: Carmen Valenzuela (Madre) +56 9 8877 6655"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-2 focus:ring-[#CEA549]/20 transition-all bg-[#FAF9F6]/50"
                  required
                />
              </div>
            </div>
          </div>

          <!-- STEP 2: PERFIL DEPORTIVO Y OBJETIVOS -->
          <div v-else-if="currentStep === 2" class="space-y-5 animate-in fade-in duration-150">
            <div>
              <h2 class="text-lg sm:text-xl font-bold text-[#191A1D] flex items-center gap-2">
                <Activity class="w-5 h-5 text-[#CEA549]" />
                Paso 2: Perfil Deportivo y Objetivos
              </h2>
              <p class="text-xs sm:text-sm text-[#6B6D74] mt-1">
                Cuéntanos qué buscas lograr para estructurar las progresiones con el Coach Pablo.
              </p>
            </div>

            <!-- Objetivo Principal (Cards) -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-2">
                ¿Cuál es tu objetivo principal? <span class="text-rose-500">*</span>
              </label>

              <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <!-- Opción 1: Salud y Longevidad -->
                <div
                  @click="form.primaryGoal = 'Salud, funcionalidad y longevidad'"
                  :class="form.primaryGoal === 'Salud, funcionalidad y longevidad' ? 'border-[#CEA549] bg-[#CEA549]/10 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] hover:border-[#CEA549]/40 bg-white'"
                  class="p-4 rounded-xl border cursor-pointer transition-all flex flex-col justify-between"
                >
                  <div class="space-y-2">
                    <Heart class="w-5 h-5 text-[#CEA549]" />
                    <p class="font-bold text-xs sm:text-sm text-[#191A1D]">Salud, Funcionalidad & Longevidad</p>
                    <p class="text-[11px] text-[#6B6D74] leading-relaxed">
                      Mejorar postura, articulaciones, vitalidad diaria y masa muscular saludable.
                    </p>
                  </div>
                  <div class="mt-3 flex items-center justify-end">
                    <div
                      class="w-4 h-4 rounded-full border flex items-center justify-center"
                      :class="form.primaryGoal === 'Salud, funcionalidad y longevidad' ? 'border-[#CEA549] bg-[#CEA549]' : 'border-[#DCD8CF]'"
                    >
                      <div v-if="form.primaryGoal === 'Salud, funcionalidad y longevidad'" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                    </div>
                  </div>
                </div>

                <!-- Opción 2: Estética -->
                <div
                  @click="form.primaryGoal = 'Estética'"
                  :class="form.primaryGoal === 'Estética' ? 'border-[#CEA549] bg-[#CEA549]/10 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] hover:border-[#CEA549]/40 bg-white'"
                  class="p-4 rounded-xl border cursor-pointer transition-all flex flex-col justify-between"
                >
                  <div class="space-y-2">
                    <Activity class="w-5 h-5 text-[#CEA549]" />
                    <p class="font-bold text-xs sm:text-sm text-[#191A1D]">Estética & Composición</p>
                    <p class="text-[11px] text-[#6B6D74] leading-relaxed">
                      Hipertrofia muscular, recomposición corporal, tonificación y reducción de grasa.
                    </p>
                  </div>
                  <div class="mt-3 flex items-center justify-end">
                    <div
                      class="w-4 h-4 rounded-full border flex items-center justify-center"
                      :class="form.primaryGoal === 'Estética' ? 'border-[#CEA549] bg-[#CEA549]' : 'border-[#DCD8CF]'"
                    >
                      <div v-if="form.primaryGoal === 'Estética'" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                    </div>
                  </div>
                </div>

                <!-- Opción 3: Rendimiento -->
                <div
                  @click="form.primaryGoal = 'Rendimiento deportivo'"
                  :class="form.primaryGoal === 'Rendimiento deportivo' ? 'border-[#CEA549] bg-[#CEA549]/10 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] hover:border-[#CEA549]/40 bg-white'"
                  class="p-4 rounded-xl border cursor-pointer transition-all flex flex-col justify-between"
                >
                  <div class="space-y-2">
                    <User class="w-5 h-5 text-[#CEA549]" />
                    <p class="font-bold text-xs sm:text-sm text-[#191A1D]">Rendimiento Deportivo</p>
                    <p class="text-[11px] text-[#6B6D74] leading-relaxed">
                      Aumento de fuerza máxima, potencia atlética y complemento a otros deportes.
                    </p>
                  </div>
                  <div class="mt-3 flex items-center justify-end">
                    <div
                      class="w-4 h-4 rounded-full border flex items-center justify-center"
                      :class="form.primaryGoal === 'Rendimiento deportivo' ? 'border-[#CEA549] bg-[#CEA549]' : 'border-[#DCD8CF]'"
                    >
                      <div v-if="form.primaryGoal === 'Rendimiento deportivo'" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Detalle de objetivos -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                Cuéntanos con tus propias palabras qué te gustaría conseguir
              </label>
              <textarea
                v-model="form.goalDetails"
                rows="2"
                placeholder="Ej: Quiero volver a entrenar sin que me duela la zona lumbar, ganar fuerza en piernas y tener más energía durante el día..."
                class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              ></textarea>
            </div>

            <!-- Experiencia previa -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Has entrenado fuerza o pesas antes?
                </label>
                <div class="flex items-center gap-6 py-2">
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="Si" v-model="form.hasTrainedBefore" class="accent-[#CEA549]" /> Sí, tengo experiencia
                  </label>
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="No" v-model="form.hasTrainedBefore" class="accent-[#CEA549]" /> No, seré principiante
                  </label>
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  Tiempo de experiencia entrenando
                </label>
                <select
                  v-model="form.trainingExperience"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-white"
                >
                  <option value="Sin experiencia previa">Sin experiencia previa</option>
                  <option value="1-6 meses">1 a 6 meses</option>
                  <option value="6-12 meses">6 a 12 meses</option>
                  <option value="Más de 1 año">Más de 1 año</option>
                </select>
              </div>
            </div>

            <!-- Deportes practicados -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                ¿Qué deportes o disciplinas practicas o has practicado?
              </label>
              <input
                v-model="form.sportsPracticed"
                type="text"
                placeholder="Ej: Fútbol, running, natación, artes marciales, ninguno..."
                class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>

            <!-- Frecuencia deseada -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                ¿Cuántos días a la semana planeas entrenar en INVERFIT? <span class="text-rose-500">*</span>
              </label>
              <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
                <button
                  type="button"
                  v-for="freq in ['1 vez/sem.', '2 veces/sem.', '3 veces/sem.', '4 veces/sem.'] as FrequencyType[]"
                  :key="freq"
                  @click="form.frequency = freq"
                  class="py-2.5 px-3 rounded-lg border text-xs font-bold transition-all text-center cursor-pointer"
                  :class="form.frequency === freq ? 'bg-[#CEA549] text-black border-[#CEA549] shadow-xs' : 'bg-white text-[#191A1D] border-[#DCD8CF] hover:border-[#CEA549]/50'"
                >
                  {{ freq }}
                </button>
              </div>
            </div>
          </div>

          <!-- STEP 3: SALUD Y ANTECEDENTES MÉDICOS (PAR-Q+) -->
          <div v-else-if="currentStep === 3" class="space-y-5 animate-in fade-in duration-150">
            <div>
              <h2 class="text-lg sm:text-xl font-bold text-[#191A1D] flex items-center gap-2">
                <Heart class="w-5 h-5 text-rose-600" />
                Paso 3: Ficha Clínica y Antecedentes Médicos (PAR-Q+)
              </h2>
              <p class="text-xs sm:text-sm text-[#6B6D74] mt-1">
                La seguridad es lo primero. Esta información es estrictamente confidencial para dosificar cargas y evitar lesiones.
              </p>
            </div>

            <!-- Enfermedades Crónicas -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                ¿Padeces alguna enfermedad crónica diagnosticada?
              </label>
              <p class="text-[11px] text-[#6B6D74] mb-1.5">
                (Asma, hipertensión, diabetes, afecciones cardíacas, problemas tiroideos, etc. Si no tienes, escribe "Ninguna")
              </p>
              <input
                v-model="form.chronicDiseases"
                type="text"
                placeholder="Ej: Hipertensión controlada con medicamentos, o Ninguna"
                class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>

            <!-- Operaciones / Lesiones -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                ¿Has tenido operaciones, fracturas o lesiones crónicas?
              </label>
              <p class="text-[11px] text-[#6B6D74] mb-1.5">
                (Indica qué zona del cuerpo y fecha aproximada. Ej: fractura de clavícula 2021, meniscos rodilla derecha 2019)
              </p>
              <textarea
                v-model="form.surgeriesInjuries"
                rows="2"
                placeholder="Ej: Operación de menisco derecho en 2020. Esguince de tobillo recurrente..."
                class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              ></textarea>
            </div>

            <!-- Limitaciones o dolor actual -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Tienes algún dolor o limitación física actualmente?
                </label>
                <div class="flex items-center gap-4 py-2">
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="No" v-model="form.injuriesPreventExercise" class="accent-[#CEA549]" /> No, me siento bien
                  </label>
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="Si" v-model="form.injuriesPreventExercise" class="accent-[#CEA549]" /> Sí, tengo molestias
                  </label>
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Alguna dificultad para comenzar a entrenar?
                </label>
                <input
                  v-model="form.currentDifficulties"
                  type="text"
                  placeholder="Ej: Fatiga rápida, mareos con calor, o Ninguna"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>
            </div>

            <!-- Antecedente COVID-19 -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 border-t border-[#ECE8E1] pt-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Tuviste COVID-19 en el pasado?
                </label>
                <div class="flex items-center gap-4 py-2">
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="No" v-model="form.hadCovid" class="accent-[#CEA549]" /> No
                  </label>
                  <label class="flex items-center gap-1.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                    <input type="radio" value="Si" v-model="form.hadCovid" class="accent-[#CEA549]" /> Sí
                  </label>
                </div>
              </div>

              <div v-if="form.hadCovid === 'Si'">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                  ¿Quedaste con secuelas respiratorias o de fatiga?
                </label>
                <input
                  v-model="form.covidSequelae"
                  type="text"
                  placeholder="Ej: Ninguna, o ligera agitación al subir escaleras"
                  class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>
            </div>

            <!-- Hábitos de consumo -->
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                Hábitos y consumo frecuente
              </label>
              <input
                v-model="form.frequentConsumption"
                type="text"
                placeholder="Ej: Café frecuente, suplemento de creatina, fumador ocasional, o Ninguno..."
                class="w-full px-3.5 py-2.5 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>
          </div>

          <!-- STEP 4: CONSENTIMIENTO INFORMADO Y FIRMA DIGITAL TÁCTIL -->
          <div v-else-if="currentStep === 4" class="space-y-5 animate-in fade-in duration-150">
            <div>
              <h2 class="text-lg sm:text-xl font-bold text-[#191A1D] flex items-center gap-2">
                <FileCheck2 class="w-5 h-5 text-[#CEA549]" />
                Paso 4: Consentimiento Informado y Firma Digital
              </h2>
              <p class="text-xs sm:text-sm text-[#6B6D74] mt-1">
                Por favor lee las condiciones oficiales de INVERFIT y estampa tu firma digital con tu dedo o puntero.
              </p>
            </div>

            <!-- Terms Box -->
            <div class="p-4 bg-[#FAF9F6] rounded-xl border border-[#E2DFD8] text-xs text-[#4A4C52] space-y-3 max-h-52 overflow-y-auto leading-relaxed shadow-inner">
              <p class="font-bold text-[#191A1D] uppercase tracking-wider">
                TÉRMINOS Y CONDICIONES DEL SERVICIO DE ENTRENAMIENTO — INVERFIT
              </p>

              <div>
                <p class="font-bold text-[#191A1D]">1. Declaración Jurada de Salud y Veracidad:</p>
                <p>
                  El alumno declara bajo fe de juramento que la información médica consignada en este formulario es completa y verídica, eximiendo a INVERFIT de consecuencias derivadas de omisiones deliberadas sobre su salud.
                </p>
              </div>

              <div>
                <p class="font-bold text-[#191A1D]">2. Asunción Voluntaria de Riesgos:</p>
                <p>
                  El alumno comprende que la actividad física guiada de fuerza y acondicionamiento involucra exigencia corporal. Se compromete a acatar las instrucciones de técnica y seguridad impartidas por el Coach Pablo y a reportar inmediatamente cualquier molestia o dolor agudo durante las sesiones.
                </p>
              </div>

              <div>
                <p class="font-bold text-[#191A1D]">3. Política de Facturación y Pago Mensual (Días 1 al 5):</p>
                <p>
                  Para garantizar la reserva de su cupo en el bloque horario fijo acordado (máximo 3 alumnos por bloque de 75 minutos), <strong class="text-[#191A1D]">la mensualidad debe cancelarse de manera anticipada durante los primeros 5 días hábiles de cada mes (del 1 al 5)</strong>. Pagos posteriores al día 5 quedan sujetos a pérdida de reserva del bloque horario preferencial.
                </p>
              </div>

              <div>
                <p class="font-bold text-[#191A1D]">4. Asistencia y Recuperaciones:</p>
                <p>
                  Las inasistencias avisadas con al menos 24 horas de anticipación podrán ser reagendadas dentro de la misma semana calendario, según cupos disponibles en otros bloques y mutuo acuerdo con el entrenador. Inasistencias sin aviso no aplazan la mensualidad.
                </p>
              </div>
            </div>

            <!-- Checkboxes obligatorios -->
            <div class="space-y-2.5 pt-2">
              <label class="flex items-start gap-2.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                <input
                  type="checkbox"
                  v-model="form.consentTermsRead"
                  class="w-4 h-4 mt-0.5 rounded border-[#DCD8CF] text-[#CEA549] focus:ring-[#CEA549] accent-[#CEA549]"
                  required
                />
                <span>
                  He leído y acepto íntegramente las políticas del gimnasio, en especial el pago de la mensualidad entre los días 1 y 5 de cada mes y los bloques reducidos de 75 minutos. <span class="text-rose-500">*</span>
                </span>
              </label>

              <label class="flex items-start gap-2.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                <input
                  type="checkbox"
                  v-model="form.consentTruthDeclared"
                  class="w-4 h-4 mt-0.5 rounded border-[#DCD8CF] text-[#CEA549] focus:ring-[#CEA549] accent-[#CEA549]"
                  required
                />
                <span>
                  Doy fe de que los datos médicos y antecedentes personales declarados son verídicos y autorizo al Coach Pablo a custodiarlos en mi ficha deportiva confidencial. <span class="text-rose-500">*</span>
                </span>
              </label>
            </div>

            <!-- CANVAS DE FIRMA DIGITAL -->
            <div class="space-y-2 pt-2">
              <div class="flex items-center justify-between">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider">
                  Firma Digital Táctil <span class="text-rose-500">*</span>
                </label>
                <button
                  type="button"
                  @click="clearSignature"
                  class="text-xs text-[#8C6B1B] hover:text-[#614910] font-semibold flex items-center gap-1 cursor-pointer"
                >
                  <RotateCcw class="w-3.5 h-3.5" />
                  Borrar firma
                </button>
              </div>

              <div class="relative bg-white rounded-xl border-2 border-dashed border-[#DCD8CF] overflow-hidden group hover:border-[#CEA549]/70 transition-colors shadow-2xs">
                <canvas
                  ref="signatureCanvas"
                  class="w-full h-40 touch-none cursor-crosshair bg-white"
                  @mousedown="startDrawing"
                  @mousemove="draw"
                  @mouseup="stopDrawing"
                  @mouseleave="stopDrawing"
                  @touchstart.prevent="startDrawing"
                  @touchmove.prevent="draw"
                  @touchend.prevent="stopDrawing"
                ></canvas>

                <!-- Watermark placeholder when empty -->
                <div
                  v-if="!hasSignature"
                  class="absolute inset-0 pointer-events-none flex flex-col items-center justify-center text-[#A6A8B0] text-xs gap-1"
                >
                  <FileCheck2 class="w-6 h-6 opacity-40 text-[#CEA549]" />
                  <span>Dibuja tu firma aquí con tu dedo o puntero</span>
                </div>

                <!-- Verified stamp when signed -->
                <div
                  v-if="hasSignature"
                  class="absolute bottom-2 right-2 pointer-events-none bg-[#E6F4EA] border border-[#137333]/30 px-2 py-0.5 rounded text-[10px] font-bold text-[#137333] flex items-center gap-1 shadow-xs"
                >
                  <Check class="w-3 h-3" /> Firma capturada
                </div>
              </div>

              <p v-if="signatureTimestamp" class="text-[11px] text-[#6B6D74] italic">
                Fecha de firma digital: {{ signatureTimestamp }}
              </p>
            </div>
          </div>

          <!-- Navigation Footer Buttons -->
          <div class="pt-6 border-t border-[#ECE8E1] flex items-center justify-between gap-3">
            <button
              v-if="currentStep > 1"
              type="button"
              @click="prevStep"
              :disabled="isSubmitting"
              class="h-11 px-5 rounded-xl border border-[#DCD8CF] bg-white hover:bg-[#FAF9F6] text-xs sm:text-sm font-bold text-[#191A1D] flex items-center gap-2 cursor-pointer transition-colors"
            >
              <ArrowLeft class="w-4 h-4" />
              <span>Anterior</span>
            </button>
            <div v-else></div>

            <button
              v-if="currentStep < totalSteps"
              type="button"
              @click="nextStep"
              class="h-11 px-6 rounded-xl bg-[#CEA549] hover:bg-[#BA923B] text-xs sm:text-sm font-bold text-[#0B0B0D] flex items-center gap-2 cursor-pointer shadow-md transition-all ml-auto"
            >
              <span>Continuar</span>
              <ArrowRight class="w-4 h-4" />
            </button>

            <button
              v-else
              type="button"
              @click="handleSubmit"
              :disabled="isSubmitting"
              class="h-11 px-7 rounded-xl bg-[#CEA549] hover:bg-[#BA923B] text-xs sm:text-sm font-bold text-[#0B0B0D] flex items-center gap-2 cursor-pointer shadow-md transition-all ml-auto disabled:opacity-60 disabled:cursor-not-allowed"
            >
              <span v-if="isSubmitting">Enviando expediente...</span>
              <span v-else>Finalizar y Enviar Ficha</span>
              <Send class="w-4 h-4" :class="{ 'animate-pulse': isSubmitting }" />
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="py-6 text-center text-xs text-[#8E9097] border-t border-[#E2DFD8] mt-8 bg-white/60">
      <div class="max-w-3xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
        <p>© 2026 INVERFIT — Pablo Henríquez San Martín · Entrenador Personal</p>
        <p class="text-[11px] text-[#A6A8B0]">Formulario Oficial y Consentimiento Informado v1.0</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}
.animate-shake {
  animation: shake 0.3s ease-in-out;
}
</style>
