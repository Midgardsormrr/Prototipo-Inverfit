<script setup lang="ts">
import { ref } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { FrequencyType, StudentStatus, StudentOnboarding } from '@/types'
import {
  Share2,
  Check,
  User,
  Heart,
  Activity,
  Zap,
  ClipboardList,
  ShieldCheck
} from 'lucide-vue-next'

defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', studentId: string): void
}>()

const gymStore = useGymStore()

// Mode: 'quick' (10 seconds) or 'full' (20 questions official in-person onboarding)
const registerMode = ref<'quick' | 'full'>('full')
const activeFullTab = ref<'demographics' | 'sports' | 'clinical'>('demographics')

// Common Fields
const fullName = ref('')
const frequency = ref<FrequencyType>('2 veces/sem.')
const fixedSchedule = ref('')
const email = ref('')
const phone = ref('')
const membershipStatus = ref<StudentStatus>('Activo')
const emergencyContact = ref('')
const healthConditions = ref('')
const notes = ref('')

// Full Onboarding 20 Questions Fields
const rut = ref('')
const birthDate = ref('')
const age = ref('')
const gender = ref<'Femenino' | 'Masculino' | 'Otro'>('Femenino')
const comuna = ref('Concepción')
const occupation = ref('')
const hasChildren = ref<'Si' | 'No'>('No')
const childrenCount = ref('')

const primaryGoal = ref<'Salud, funcionalidad y longevidad' | 'Estética' | 'Rendimiento deportivo' | 'Alto Rendimiento' | 'Otros'>('Salud, funcionalidad y longevidad')
const goalDetails = ref('')
const hasTrainedBefore = ref<'Si' | 'No'>('Si')
const trainingExperience = ref<'Sin experiencia previa' | '1-6 meses' | '6-12 meses' | 'Más de 1 año'>('1-6 meses')
const sportsPracticed = ref('')

const chronicDiseases = ref('')
const surgeriesInjuries = ref('')
const injuriesPreventExercise = ref<'Si' | 'No'>('No')
const currentDifficulties = ref('')
const hadCovid = ref<'Si' | 'No'>('No')
const covidSequelae = ref('')
const frequentConsumption = ref('')
const consentAcceptedInPerson = ref(true)

const isSubmitting = ref(false)
const errorMessage = ref('')
const copiedOnboardingLink = ref(false)

function handleBirthDateChange() {
  if (!birthDate.value) return
  const birth = new Date(birthDate.value)
  const now = new Date()
  let calculatedAge = now.getFullYear() - birth.getFullYear()
  const m = now.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
    calculatedAge--
  }
  if (calculatedAge > 0 && calculatedAge < 120) {
    age.value = String(calculatedAge)
  }
}

async function copyOnboardingLink() {
  const url = `${window.location.origin}/registro-alumno`
  try {
    await navigator.clipboard.writeText(url)
    copiedOnboardingLink.value = true
    gymStore.showToast('Enlace copiado al portapapeles. ¡Pégalo en WhatsApp para que el alumno se auto-registre!', 'success', 'Enlace copiado')
    setTimeout(() => {
      copiedOnboardingLink.value = false
    }, 3000)
  } catch {
    gymStore.showToast(`Enlace: ${url}`, 'info')
  }
}

function resetForm() {
  fullName.value = ''
  frequency.value = '2 veces/sem.'
  fixedSchedule.value = ''
  email.value = ''
  phone.value = ''
  membershipStatus.value = 'Activo'
  emergencyContact.value = ''
  healthConditions.value = ''
  notes.value = ''

  rut.value = ''
  birthDate.value = ''
  age.value = ''
  gender.value = 'Femenino'
  comuna.value = 'Concepción'
  occupation.value = ''
  hasChildren.value = 'No'
  childrenCount.value = ''

  primaryGoal.value = 'Salud, funcionalidad y longevidad'
  goalDetails.value = ''
  hasTrainedBefore.value = 'Si'
  trainingExperience.value = '1-6 meses'
  sportsPracticed.value = ''

  chronicDiseases.value = ''
  surgeriesInjuries.value = ''
  injuriesPreventExercise.value = 'No'
  currentDifficulties.value = ''
  hadCovid.value = 'No'
  covidSequelae.value = ''
  frequentConsumption.value = ''
  consentAcceptedInPerson.value = true

  errorMessage.value = ''
}

async function handleSubmit() {
  if (!fullName.value.trim()) {
    errorMessage.value = 'El nombre del alumno es obligatorio.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    let onboardingData: StudentOnboarding | undefined = undefined
    let finalHealth = healthConditions.value.trim()
    let finalNotes = notes.value.trim()

    if (registerMode.value === 'full') {
      onboardingData = {
        rut: rut.value.trim(),
        birthDate: birthDate.value,
        age: age.value ? Number(age.value) : undefined,
        gender: gender.value,
        comuna: comuna.value.trim(),
        occupation: occupation.value.trim(),
        hasChildren: hasChildren.value,
        childrenCount: childrenCount.value.trim(),

        primaryGoal: primaryGoal.value,
        goalDetails: goalDetails.value.trim(),
        hasTrainedBefore: hasTrainedBefore.value,
        trainingExperience: trainingExperience.value,
        sportsPracticed: sportsPracticed.value.trim(),
        availableDaysWeekly: frequency.value === '4 veces/sem.' ? '4-6 veces a la semana' : '2-4 veces a la semana',

        chronicDiseases: chronicDiseases.value.trim(),
        surgeriesInjuries: surgeriesInjuries.value.trim(),
        injuriesPreventExercise: injuriesPreventExercise.value,
        currentDifficulties: currentDifficulties.value.trim(),
        hadCovid: hadCovid.value,
        covidSequelae: covidSequelae.value.trim(),
        frequentConsumption: frequentConsumption.value.trim(),

        consentAccepted: consentAcceptedInPerson.value,
        consentTimestamp: new Date().toLocaleString('es-CL', { dateStyle: 'full', timeStyle: 'medium' }),
        termsVersion: 'v1.0 (Registro Presencial Oficial Coach Pablo)'
      }

      // Synthesize medical conditions
      const parts: string[] = []
      if (chronicDiseases.value.trim()) parts.push(`Enfermedades: ${chronicDiseases.value.trim()}`)
      if (surgeriesInjuries.value.trim()) parts.push(`Cirugías/Lesiones: ${surgeriesInjuries.value.trim()}`)
      if (injuriesPreventExercise.value === 'Si') parts.push('Presenta dolor/limitación física')
      if (hadCovid.value === 'Si' && covidSequelae.value.trim()) parts.push(`Secuelas COVID: ${covidSequelae.value.trim()}`)
      finalHealth = parts.join(' | ') || 'Sin patologías ni lesiones declaradas'

      if (goalDetails.value.trim()) {
        finalNotes = `Objetivo: ${primaryGoal.value}. Detalle: ${goalDetails.value.trim()}`
      }
    }

    const student = await gymStore.createStudent({
      fullName: fullName.value.trim(),
      frequency: frequency.value,
      fixedSchedule: fixedSchedule.value.trim() || 'Días y horas por confirmar',
      email: email.value.trim(),
      phone: phone.value.trim(),
      membershipStatus: membershipStatus.value,
      emergencyContact: emergencyContact.value.trim(),
      healthConditions: finalHealth,
      notes: finalNotes,
      onboarding: onboardingData
    })

    resetForm()
    emit('created', student.id)
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al registrar alumno.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-3 sm:p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full border border-[#E2DFD8] overflow-hidden my-6 animate-in fade-in zoom-in-95 duration-150 flex flex-col max-h-[90vh]">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] bg-[#FAF9F6] flex items-center justify-between shrink-0">
        <div>
          <h2 class="text-lg font-bold text-[#191A1D] flex items-center gap-2">
            <span>Registrar Alumno</span>
            <span class="text-xs px-2 py-0.5 rounded bg-[#CEA549]/20 text-[#8C6B1B] font-semibold border border-[#CEA549]/30">
              INVERFIT
            </span>
          </h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Añade un nuevo expediente deportivo a la plataforma.</p>
        </div>

        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <!-- Mode Selector Tabs: Rápido vs Ficha Completa -->
      <div class="px-6 pt-3 bg-white border-b border-[#ECE8E1] flex items-center justify-between gap-2 shrink-0">
        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="registerMode = 'full'"
            class="py-2.5 px-3.5 text-xs font-bold border-b-2 transition-all flex items-center gap-1.5 cursor-pointer"
            :class="registerMode === 'full' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
          >
            <ClipboardList class="w-4 h-4 text-[#8C6B1B]" />
            <span>Ficha Oficial Completa (20 preguntas)</span>
          </button>

          <button
            type="button"
            @click="registerMode = 'quick'"
            class="py-2.5 px-3.5 text-xs font-bold border-b-2 transition-all flex items-center gap-1.5 cursor-pointer"
            :class="registerMode === 'quick' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
          >
            <Zap class="w-4 h-4 text-[#8C6B1B]" />
            <span>Registro Rápido Express</span>
          </button>
        </div>

        <!-- WhatsApp Link Copy Helper Button -->
        <button
          type="button"
          @click="copyOnboardingLink"
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1 text-xs font-bold text-[#8C6B1B] hover:text-[#614910] bg-[#FAF9F6] border border-[#E2DFD8] hover:border-[#CEA549] rounded-lg transition-colors cursor-pointer"
          title="Copiar enlace para enviar por WhatsApp al alumno"
        >
          <Check v-if="copiedOnboardingLink" class="w-3.5 h-3.5 text-[#137333]" />
          <Share2 v-else class="w-3.5 h-3.5" />
          <span>{{ copiedOnboardingLink ? '¡Link copiado!' : 'Copiar link WhatsApp' }}</span>
        </button>
      </div>

      <!-- Scrollable Form Body -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4 overflow-y-auto flex-1">
        <!-- Error Banner -->
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-lg text-xs sm:text-sm">
          {{ errorMessage }}
        </div>

        <!-- Helper Callout if in Quick Mode -->
        <div v-if="registerMode === 'quick'" class="p-3 bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl flex items-center justify-between gap-3 text-xs">
          <div>
            <p class="font-bold text-[#191A1D]">¿El alumno prefiere completar sus datos desde su celular?</p>
            <p class="text-[#6B6D74]">Puedes enviarle el enlace oficial de auto-ingreso para que responda las 20 preguntas y firme.</p>
          </div>
          <button
            type="button"
            @click="copyOnboardingLink"
            class="px-3 py-1.5 bg-[#CEA549] hover:bg-[#BA923B] text-black font-bold rounded-lg transition-colors shrink-0 flex items-center gap-1 cursor-pointer shadow-xs"
          >
            <Check v-if="copiedOnboardingLink" class="w-3.5 h-3.5 text-[#137333]" />
            <Share2 v-else class="w-3.5 h-3.5 text-black" />
            <span>{{ copiedOnboardingLink ? 'Copiado' : 'Copiar link' }}</span>
          </button>
        </div>

        <!-- MODE 1: REGISTRO RÁPIDO (EXPRESS) -->
        <div v-if="registerMode === 'quick'" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
              Nombre Completo *
            </label>
            <input
              v-model="fullName"
              type="text"
              placeholder="Ej. Nelson López Valenzuela"
              class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all bg-[#FAF9F6]/50"
              required
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Frecuencia Semanal *
              </label>
              <select
                v-model="frequency"
                class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] transition-all"
              >
                <option value="1 vez/sem.">1 vez/sem.</option>
                <option value="2 veces/sem.">2 veces/sem.</option>
                <option value="3 veces/sem.">3 veces/sem.</option>
                <option value="4 veces/sem.">4 veces/sem.</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Estado Inicial
              </label>
              <select
                v-model="membershipStatus"
                class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] transition-all"
              >
                <option value="Activo">Activo</option>
                <option value="Por revisar">Por revisar</option>
                <option value="Inactivo">Inactivo</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
              Horario Fijo Sugerido
            </label>
            <input
              v-model="fixedSchedule"
              type="text"
              placeholder="Ej. Lun y Mié 08:00 (o por coordinar)"
              class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Teléfono de Contacto
              </label>
              <input
                v-model="phone"
                type="text"
                placeholder="+56 9 1234 5678"
                class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Correo Electrónico
              </label>
              <input
                v-model="email"
                type="email"
                placeholder="alumno@ejemplo.cl"
                class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
              Contacto de Emergencia
            </label>
            <input
              v-model="emergencyContact"
              type="text"
              placeholder="Nombre y teléfono de familiar o contacto"
              class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
              Condiciones de Salud o Lesiones
            </label>
            <textarea
              v-model="healthConditions"
              rows="2"
              placeholder="Antecedentes médicos relevantes, cirugías o dolores..."
              class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
              Observaciones Adicionales
            </label>
            <textarea
              v-model="notes"
              rows="2"
              placeholder="Notas generales del entrenador..."
              class="w-full px-3.5 py-2.5 border border-[#DCD8CF] rounded-lg text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
            ></textarea>
          </div>
        </div>

        <!-- MODE 2: FICHA COMPLETA OFICIAL (20 PREGUNTAS) -->
        <div v-else class="space-y-5">
          <!-- Subtabs: 1. Demográficos | 2. Objetivos | 3. PAR-Q+ Clínico -->
          <div class="grid grid-cols-3 gap-2 p-1 bg-[#FAF9F6] rounded-xl border border-[#ECE8E1] text-center text-xs">
            <button
              type="button"
              @click="activeFullTab = 'demographics'"
              class="py-2 px-2 rounded-lg font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer"
              :class="activeFullTab === 'demographics' ? 'bg-[#CEA549] text-black shadow-xs' : 'text-[#6B6D74] hover:text-[#191A1D]'"
            >
              <User class="w-3.5 h-3.5" />
              <span>1. Demografía</span>
            </button>

            <button
              type="button"
              @click="activeFullTab = 'sports'"
              class="py-2 px-2 rounded-lg font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer"
              :class="activeFullTab === 'sports' ? 'bg-[#CEA549] text-black shadow-xs' : 'text-[#6B6D74] hover:text-[#191A1D]'"
            >
              <Activity class="w-3.5 h-3.5" />
              <span>2. Metas</span>
            </button>

            <button
              type="button"
              @click="activeFullTab = 'clinical'"
              class="py-2 px-2 rounded-lg font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer"
              :class="activeFullTab === 'clinical' ? 'bg-[#CEA549] text-black shadow-xs' : 'text-[#6B6D74] hover:text-[#191A1D]'"
            >
              <Heart class="w-3.5 h-3.5" />
              <span>3. PAR-Q+</span>
            </button>
          </div>

          <!-- SUBTAB 1: DEMOGRAFÍA & CONTACTO -->
          <div v-if="activeFullTab === 'demographics'" class="space-y-4 animate-in fade-in duration-150">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Nombre Completo *
                </label>
                <input
                  v-model="fullName"
                  type="text"
                  placeholder="Ej: Nelson López Valenzuela"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  RUT *
                </label>
                <input
                  v-model="rut"
                  type="text"
                  placeholder="Ej: 19.876.543-2"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <div class="grid grid-cols-2 gap-2">
                <div>
                  <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                    Nacimiento
                  </label>
                  <input
                    v-model="birthDate"
                    type="date"
                    @change="handleBirthDateChange"
                    class="w-full px-2 py-2 rounded-lg border border-[#DCD8CF] text-xs text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  />
                </div>
                <div>
                  <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                    Edad
                  </label>
                  <input
                    v-model="age"
                    type="number"
                    placeholder="Años"
                    class="w-full px-2 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Teléfono / WhatsApp *
                </label>
                <input
                  v-model="phone"
                  type="tel"
                  placeholder="+56 9 1234 5678"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Correo Electrónico *
                </label>
                <input
                  v-model="email"
                  type="email"
                  placeholder="alumno@inverfit.cl"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Género
                </label>
                <select
                  v-model="gender"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
                >
                  <option value="Femenino">Femenino</option>
                  <option value="Masculino">Masculino</option>
                  <option value="Otro">Otro</option>
                </select>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Comuna
                </label>
                <input
                  v-model="comuna"
                  type="text"
                  placeholder="Ej: Concepción, San Pedro"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Ocupación
                </label>
                <input
                  v-model="occupation"
                  type="text"
                  placeholder="Ej: Ingeniero, Estudiante"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  ¿Tiene hijos?
                </label>
                <div class="flex items-center gap-3 py-1">
                  <label class="flex items-center gap-1 text-xs cursor-pointer">
                    <input type="radio" value="No" v-model="hasChildren" class="accent-[#CEA549]" /> No
                  </label>
                  <label class="flex items-center gap-1 text-xs cursor-pointer">
                    <input type="radio" value="Si" v-model="hasChildren" class="accent-[#CEA549]" /> Sí
                  </label>
                  <input
                    v-if="hasChildren === 'Si'"
                    v-model="childrenCount"
                    type="text"
                    placeholder="¿Cuántos?"
                    class="w-20 px-2 py-1 border border-[#DCD8CF] rounded text-xs"
                  />
                </div>
              </div>

              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Contacto de Emergencia (Nombre, Parentesco y Teléfono) *
                </label>
                <input
                  v-model="emergencyContact"
                  type="text"
                  placeholder="Ej: Carmen Valenzuela (Madre) +56 9 8877 6655"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                  required
                />
              </div>
            </div>

            <div class="pt-2 flex justify-end">
              <button
                type="button"
                @click="activeFullTab = 'sports'"
                class="px-4 py-2 bg-[#FAF9F6] border border-[#DCD8CF] hover:border-[#CEA549] rounded-lg text-xs font-bold text-[#191A1D] cursor-pointer"
              >
                Siguiente: Objetivos y Perfil →
              </button>
            </div>
          </div>

          <!-- SUBTAB 2: PERFIL DEPORTIVO Y OBJETIVOS -->
          <div v-else-if="activeFullTab === 'sports'" class="space-y-4 animate-in fade-in duration-150">
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1.5">
                Objetivo Principal (Categoría Oficial INVERFIT) *
              </label>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                <div
                  @click="primaryGoal = 'Salud, funcionalidad y longevidad'"
                  :class="primaryGoal === 'Salud, funcionalidad y longevidad' ? 'border-[#CEA549] bg-[#CEA549]/15 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] bg-white'"
                  class="p-3 rounded-lg border cursor-pointer text-xs space-y-1"
                >
                  <p class="font-bold text-[#191A1D]">Salud & Longevidad</p>
                  <p class="text-[11px] text-[#6B6D74]">Postura, articulaciones y vitalidad.</p>
                </div>

                <div
                  @click="primaryGoal = 'Estética'"
                  :class="primaryGoal === 'Estética' ? 'border-[#CEA549] bg-[#CEA549]/15 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] bg-white'"
                  class="p-3 rounded-lg border cursor-pointer text-xs space-y-1"
                >
                  <p class="font-bold text-[#191A1D]">Estética & Hipertrofia</p>
                  <p class="text-[11px] text-[#6B6D74]">Recomposición y tonificación.</p>
                </div>

                <div
                  @click="primaryGoal = 'Rendimiento deportivo'"
                  :class="primaryGoal === 'Rendimiento deportivo' ? 'border-[#CEA549] bg-[#CEA549]/15 ring-1 ring-[#CEA549]' : 'border-[#E2DFD8] bg-white'"
                  class="p-3 rounded-lg border cursor-pointer text-xs space-y-1"
                >
                  <p class="font-bold text-[#191A1D]">Rendimiento Deportivo</p>
                  <p class="text-[11px] text-[#6B6D74]">Fuerza máxima y potencia atlética.</p>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Detalle de Metas Declaradas
              </label>
              <textarea
                v-model="goalDetails"
                rows="2"
                placeholder="Metas particulares que el alumno busca alcanzar..."
                class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              ></textarea>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  ¿Experiencia entrenando pesas?
                </label>
                <div class="flex items-center gap-4 py-2">
                  <label class="flex items-center gap-1.5 text-xs cursor-pointer">
                    <input type="radio" value="Si" v-model="hasTrainedBefore" class="accent-[#CEA549]" /> Sí
                  </label>
                  <label class="flex items-center gap-1.5 text-xs cursor-pointer">
                    <input type="radio" value="No" v-model="hasTrainedBefore" class="accent-[#CEA549]" /> No (Principiante)
                  </label>
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Tiempo de experiencia
                </label>
                <select
                  v-model="trainingExperience"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
                >
                  <option value="Sin experiencia previa">Sin experiencia previa</option>
                  <option value="1-6 meses">1 a 6 meses</option>
                  <option value="6-12 meses">6 a 12 meses</option>
                  <option value="Más de 1 año">Más de 1 año</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Deportes o Actividad Física Practicada
              </label>
              <input
                v-model="sportsPracticed"
                type="text"
                placeholder="Ej: Fútbol fin de semana, running, natación, artes marciales..."
                class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Frecuencia Solicitada *
                </label>
                <select
                  v-model="frequency"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
                >
                  <option value="1 vez/sem.">1 vez/sem.</option>
                  <option value="2 veces/sem.">2 veces/sem.</option>
                  <option value="3 veces/sem.">3 veces/sem.</option>
                  <option value="4 veces/sem.">4 veces/sem.</option>
                </select>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Horario Fijo Acordado
                </label>
                <input
                  v-model="fixedSchedule"
                  type="text"
                  placeholder="Ej. Mar y Jue 08:00 (o por coordinar)"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>
            </div>

            <div class="pt-2 flex justify-between">
              <button
                type="button"
                @click="activeFullTab = 'demographics'"
                class="px-4 py-2 border border-[#DCD8CF] rounded-lg text-xs font-bold text-[#6B6D74] cursor-pointer"
              >
                ← Demografía
              </button>
              <button
                type="button"
                @click="activeFullTab = 'clinical'"
                class="px-4 py-2 bg-[#FAF9F6] border border-[#DCD8CF] hover:border-[#CEA549] rounded-lg text-xs font-bold text-[#191A1D] cursor-pointer"
              >
                Siguiente: PAR-Q+ Clínico →
              </button>
            </div>
          </div>

          <!-- SUBTAB 3: FICHA CLÍNICA PAR-Q+ Y CONSENTIMIENTO -->
          <div v-else-if="activeFullTab === 'clinical'" class="space-y-4 animate-in fade-in duration-150">
            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                ¿Enfermedades crónicas diagnosticadas?
              </label>
              <input
                v-model="chronicDiseases"
                type="text"
                placeholder="Ej: Asma, hipertensión, diabetes, o Ninguna"
                class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                Operaciones, Cirugías, Fracturas o Lesiones y Fechas
              </label>
              <textarea
                v-model="surgeriesInjuries"
                rows="2"
                placeholder="Ej: Cirugía meniscos rodilla derecha 2021, fractura clavícula 2019..."
                class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
              ></textarea>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  ¿Molestias o dolor físico actual?
                </label>
                <div class="flex items-center gap-4 py-1.5">
                  <label class="flex items-center gap-1.5 text-xs cursor-pointer">
                    <input type="radio" value="No" v-model="injuriesPreventExercise" class="accent-[#CEA549]" /> No
                  </label>
                  <label class="flex items-center gap-1.5 text-xs cursor-pointer">
                    <input type="radio" value="Si" v-model="injuriesPreventExercise" class="accent-[#CEA549]" /> Sí, presenta molestias
                  </label>
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Dificultades para comenzar
                </label>
                <input
                  v-model="currentDifficulties"
                  type="text"
                  placeholder="Ej: Fatiga rápida, mareos, o Ninguna"
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  ¿Antecedente COVID-19?
                </label>
                <div class="flex items-center gap-3 py-1.5">
                  <label class="flex items-center gap-1 text-xs cursor-pointer">
                    <input type="radio" value="No" v-model="hadCovid" class="accent-[#CEA549]" /> No
                  </label>
                  <label class="flex items-center gap-1 text-xs cursor-pointer">
                    <input type="radio" value="Si" v-model="hadCovid" class="accent-[#CEA549]" /> Sí
                  </label>
                  <input
                    v-if="hadCovid === 'Si'"
                    v-model="covidSequelae"
                    type="text"
                    placeholder="¿Secuelas?"
                    class="flex-1 px-2 py-1 border border-[#DCD8CF] rounded text-xs"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-[#191A1D] uppercase tracking-wider mb-1">
                  Hábitos de consumo frecuente
                </label>
                <input
                  v-model="frequentConsumption"
                  type="text"
                  placeholder="Ej: Café, tabaco ocasional, suplementos..."
                  class="w-full px-3.5 py-2 rounded-lg border border-[#DCD8CF] text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] bg-[#FAF9F6]/50"
                />
              </div>
            </div>

            <!-- Consentimiento Presencial Checkbox -->
            <div class="p-3.5 bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl space-y-2">
              <div class="flex items-center gap-2">
                <ShieldCheck class="w-4 h-4 text-[#CEA549]" />
                <span class="text-xs font-bold text-[#191A1D]">Consentimiento Informado Presencial</span>
              </div>
              <label class="flex items-start gap-2.5 text-xs text-[#191A1D] font-medium cursor-pointer">
                <input
                  type="checkbox"
                  v-model="consentAcceptedInPerson"
                  class="w-4 h-4 mt-0.5 rounded border-[#DCD8CF] text-[#CEA549] focus:ring-[#CEA549] accent-[#CEA549]"
                />
                <span>
                  El alumno revisó y aceptó presencialmente las normas de INVERFIT (pagos del 1 al 5 de cada mes, bloques de 75 min) y dio fe de la veracidad de su información médica.
                </span>
              </label>
            </div>
          </div>
        </div>

        <!-- Footer actions -->
        <div class="pt-4 border-t border-[#E2DFD8] flex items-center justify-between gap-3 shrink-0">
          <div class="text-xs text-[#8E9097]">
            <span v-if="registerMode === 'full'">Modo: Ficha Completa Oficial</span>
            <span v-else>Modo: Registro Rápido</span>
          </div>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="emit('close')"
              class="px-4 py-2 border border-[#DCD8CF] rounded-lg text-xs sm:text-sm text-[#6B6D74] hover:bg-[#F5F3EE] hover:text-[#191A1D] font-medium transition-colors cursor-pointer"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-lg text-xs sm:text-sm font-bold transition-all shadow-xs cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
            >
              <span v-if="isSubmitting">Guardando...</span>
              <span v-else>+ Registrar Alumno</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>
