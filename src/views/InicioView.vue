<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import {
  Calendar,
  Users,
  CreditCard,
  ClipboardCheck,
  FileText,
  ChevronRight,
  Share2,
  Copy,
  Check,
  ShieldCheck,
  Send,
  Plus,
  Clock
} from 'lucide-vue-next'
import RegisterStudentModal from '@/components/modals/RegisterStudentModal.vue'
import BlockDetailsModal from '@/components/modals/BlockDetailsModal.vue'
import AssignScheduleModal from '@/components/modals/AssignScheduleModal.vue'
import BatchRemindersModal from '@/components/modals/BatchRemindersModal.vue'
import ConfigureSumupLinkModal from '@/components/modals/ConfigureSumupLinkModal.vue'
import type { ScheduleBlock, DayOfWeek, SessionOccurrence } from '@/types'
import {
  getWeekDates,
  timeToMinutes,
  PABLO_FIXED_PRESET_BLOCKS,
  PABLO_SATURDAY_EVALUATION_BLOCKS,
  INVERFIT_MAX_STUDENTS_PER_BLOCK,
  type WeekDayInfo,
  type FixedPresetBlock
} from '@/services/schedule'

const router = useRouter()
const gymStore = useGymStore()

onMounted(async () => {
  const dates = getWeekDates(0)
  const from = dates[0].dateStr
  const to = dates[6].dateStr

  await Promise.all([
    gymStore.loadStudents().catch(() => {}),
    gymStore.loadPayments().catch(() => {}),
    gymStore.isScheduleApiMode()
      ? Promise.all([
          gymStore.fetchWeekOccurrences(from, to).catch(() => {}),
          gymStore.loadSchedule(from, to).catch(() => {})
        ]).catch(() => {})
      : gymStore.loadSchedule(from, to).catch(() => {})
  ])
})

// Modals
const showRegisterModal = ref(false)
const showBlockModal = ref(false)
const showAssignModal = ref(false)
const showBatchRemindersModal = ref(false)
const showConfigureSumupModal = ref(false)

// Active block and occurrence selection
const activeBlockId = ref<string | null>(null)
const activeOccurrence = ref<SessionOccurrence | null>(null)

const activeBlock = computed<ScheduleBlock | null>(() => {
  if (activeBlockId.value) {
    const found = gymStore.schedule.find(b => b.id === activeBlockId.value)
    if (found) return found
  }
  if (activeOccurrence.value) {
    const found = gymStore.schedule.find(b => b.id === activeOccurrence.value!.blockId)
    if (found) return found
    return {
      id: activeOccurrence.value.blockId,
      day: activeOccurrence.value.day,
      startTime: activeOccurrence.value.startTime,
      endTime: activeOccurrence.value.endTime,
      durationMinutes: activeOccurrence.value.durationMinutes || 75,
      maxCapacity: INVERFIT_MAX_STUDENTS_PER_BLOCK,
      studentIds: activeOccurrence.value.studentIds || [],
      type: activeOccurrence.value.type || 'Sesión grupal'
    }
  }
  return null
})

// Current week info & today identification
const weekDays = computed<WeekDayInfo[]>(() => getWeekDates(0))

const todayInfo = computed<WeekDayInfo>(() => {
  const found = weekDays.value.find(w => w.isToday)
  return found || weekDays.value[0]
})

const todayDayName = computed<DayOfWeek>(() => todayInfo.value.day)

// Occurrences for this week
const weekOccurrences = computed<SessionOccurrence[]>(() => {
  return gymStore.getWeekOccurrences(weekDays.value)
})

function getOccurrencesForToday(): SessionOccurrence[] {
  return weekOccurrences.value.filter(o => o.day === todayDayName.value)
}

export interface TodaySlotItem {
  id: string
  presetId?: string
  startTime: string
  endTime: string
  durationMinutes: number
  shift: 'mañana' | 'tarde'
  isEvaluation: boolean
  isOccupied: boolean
  occurrence?: SessionOccurrence
  blockId?: string
  studentIds: string[]
  studentNames: string[]
  cuposCount: number
  isFull: boolean
  isException?: boolean
  notes?: string
}

/**
 * Genera los bloques oficiales del día actual vinculando las sesiones agendadas reales.
 * Si hoy no hay alumnos programados, los bloques permanecen disponibles (NUNCA usa fallback a Lunes).
 */
const todaySlots = computed<TodaySlotItem[]>(() => {
  const day = todayDayName.value
  if (day === 'Dom') return []

  const isSat = day === 'Sáb'
  const presets: FixedPresetBlock[] = isSat
    ? PABLO_SATURDAY_EVALUATION_BLOCKS
    : PABLO_FIXED_PRESET_BLOCKS

  const dayOccurrences = getOccurrencesForToday()
  const matchedOccIds = new Set<string>()

  // 1. Mapear cada bloque preestablecido oficial de Pablo
  const slots: TodaySlotItem[] = presets.map((preset) => {
    const occ = dayOccurrences.find(o => o.startTime === preset.startTime)
    if (occ && occ.studentIds && occ.studentIds.length > 0) {
      matchedOccIds.add(occ.instanceId)
      const names = occ.studentIds.map(id => {
        const s = gymStore.students.find(st => st.id === id)
        return s ? s.fullName : 'Alumno'
      })
      return {
        id: occ.instanceId,
        presetId: preset.id,
        startTime: occ.startTime,
        endTime: occ.endTime,
        durationMinutes: occ.durationMinutes || 75,
        shift: preset.shift,
        isEvaluation: occ.type === 'Evaluación' || isSat,
        isOccupied: true,
        occurrence: occ,
        blockId: occ.blockId,
        studentIds: occ.studentIds,
        studentNames: names,
        cuposCount: occ.studentIds.length,
        isFull: occ.studentIds.length >= INVERFIT_MAX_STUDENTS_PER_BLOCK,
        isException: occ.isException,
        notes: occ.notes
      }
    }

    return {
      id: `slot-${day}-${preset.id}`,
      presetId: preset.id,
      startTime: preset.startTime,
      endTime: preset.endTime,
      durationMinutes: 75,
      shift: preset.shift,
      isEvaluation: isSat,
      isOccupied: false,
      occurrence: occ,
      blockId: occ?.blockId,
      studentIds: [],
      studentNames: [],
      cuposCount: 0,
      isFull: false,
      isException: false
    }
  })

  // 2. Sesiones fuera de presets estándar (si existen)
  const extraOccs = dayOccurrences.filter(o => !matchedOccIds.has(o.instanceId) && o.studentIds && o.studentIds.length > 0)
  for (const occ of extraOccs) {
    const names = (occ.studentIds || []).map(id => {
      const s = gymStore.students.find(st => st.id === id)
      return s ? s.fullName : 'Alumno'
    })
    slots.push({
      id: occ.instanceId,
      startTime: occ.startTime,
      endTime: occ.endTime,
      durationMinutes: occ.durationMinutes || 75,
      shift: timeToMinutes(occ.startTime) < 14 * 60 ? 'mañana' : 'tarde',
      isEvaluation: occ.type === 'Evaluación' || isSat,
      isOccupied: true,
      occurrence: occ,
      blockId: occ.blockId,
      studentIds: occ.studentIds || [],
      studentNames: names,
      cuposCount: (occ.studentIds || []).length,
      isFull: (occ.studentIds || []).length >= INVERFIT_MAX_STUDENTS_PER_BLOCK,
      isException: occ.isException,
      notes: occ.notes
    })
  }

  return slots.sort((a, b) => timeToMinutes(a.startTime) - timeToMinutes(b.startTime))
})

// Bloques ocupados y disponibles de hoy
const todayOccupiedSlots = computed(() => todaySlots.value.filter(s => s.isOccupied))
const todayAvailableSlots = computed(() => todaySlots.value.filter(s => !s.isOccupied))

// Métricas oficiales de ocupación de hoy (3 cupos máximos por bloque)
const todayTotalCupos = computed(() => todaySlots.value.length * INVERFIT_MAX_STUDENTS_PER_BLOCK)
const todayOccupiedCupos = computed(() => todaySlots.value.reduce((acc, s) => acc + s.cuposCount, 0))
const todayFreeCupos = computed(() => Math.max(0, todayTotalCupos.value - todayOccupiedCupos.value))

// Resumen financiero (Ciclo del 1 al 5)
const paidPayments = computed(() => gymStore.payments.filter(p => p.status === 'Pagado'))
const pendingPayments = computed(() => gymStore.payments.filter(p => p.status === 'Pendiente'))

const totalCollected = computed(() => {
  return paidPayments.value.reduce((acc, p) => acc + (p.amount || 45000), 0)
})
const totalPending = computed(() => {
  return pendingPayments.value.reduce((acc, p) => acc + (p.amount || 45000), 0)
})

// Evaluaciones y contratos recientes
const latestEvaluations = computed(() => {
  return gymStore.evaluations.slice(0, 3)
})

const latestContracts = computed(() => {
  return gymStore.contracts.slice(0, 3)
})

function getStudent(studentId: string) {
  return gymStore.students.find(s => s.id === studentId)
}

// Prefill values for AssignScheduleModal
const assignInitialDay = ref<DayOfWeek>('Lun')
const assignInitialTime = ref('06:30')
const assignInitialEndTime = ref('07:45')
const assignInitialDate = ref<string | undefined>(undefined)
const assignBlockId = ref<string | null>(null)

function openOccupiedSlot(slot: TodaySlotItem) {
  if (slot.occurrence) {
    activeOccurrence.value = slot.occurrence
    activeBlockId.value = slot.occurrence.blockId
    showBlockModal.value = true
  } else {
    openAssignForSlot(slot)
  }
}

function openAssignForSlot(slot: TodaySlotItem) {
  assignInitialDay.value = todayDayName.value
  assignInitialTime.value = slot.startTime
  assignInitialEndTime.value = slot.endTime
  assignInitialDate.value = todayInfo.value.dateStr
  assignBlockId.value = slot.blockId || null
  showAssignModal.value = true
}

function openAssignForToday() {
  assignInitialDay.value = todayDayName.value
  assignInitialTime.value = '06:30'
  assignInitialEndTime.value = '07:45'
  assignInitialDate.value = todayInfo.value.dateStr
  assignBlockId.value = null
  showAssignModal.value = true
}

function handleAddStudentToBlock(block: ScheduleBlock) {
  showBlockModal.value = false
  assignBlockId.value = block.id
  assignInitialDay.value = block.day
  assignInitialTime.value = block.startTime
  assignInitialEndTime.value = block.endTime
  assignInitialDate.value = activeOccurrence.value?.date || todayInfo.value.dateStr
  showAssignModal.value = true
}

async function handleScheduleChanged() {
  const dates = getWeekDates(0)
  const from = dates[0].dateStr
  const to = dates[6].dateStr
  await gymStore.refreshScheduleState()
  if (gymStore.isScheduleApiMode()) {
    await gymStore.fetchWeekOccurrences(from, to).catch(() => {})
  }
}

// Enlaces SumUp y Onboarding
const copiedSumupLink = ref(false)
const copiedOnboardingLink = ref(false)

async function copyPermanentSumupLink() {
  const url = gymStore.sumupPermanentLink || 'https://link.sumup.com/pay/inverfit'
  try {
    await navigator.clipboard.writeText(url)
    copiedSumupLink.value = true
    gymStore.showToast('Enlace de SumUp copiado al portapapeles.', 'success', 'Enlace SumUp')
    setTimeout(() => {
      copiedSumupLink.value = false
    }, 3000)
  } catch {
    gymStore.showToast(`Enlace: ${url}`, 'info')
  }
}

async function copyPublicOnboardingLink() {
  const url = `${window.location.origin}/registro-alumno`
  const text = `¡Hola! Te comparto tu enlace de ingreso oficial a INVERFIT para completar tu ficha de salud y firmar tu consentimiento informado:\n${url}`
  try {
    await navigator.clipboard.writeText(text)
    copiedOnboardingLink.value = true
    gymStore.showToast('Enlace de Onboarding copiado. ¡Listo para enviar por WhatsApp!', 'success', 'Enlace copiado')
    setTimeout(() => {
      copiedOnboardingLink.value = false
    }, 3000)
  } catch {
    gymStore.showToast(`Enlace: ${url}`, 'info')
  }
}
</script>

<template>
  <div class="space-y-6 max-w-[1140px] mx-auto pb-12">
    <!-- Top Executive Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <h1 class="text-[28px] font-bold text-[#191A1D] tracking-tight">Panel de Control</h1>
          <span class="px-2.5 py-0.5 rounded-full bg-[#FAF9F6] border border-[#E2DFD8] text-[11px] font-bold text-[#8C6B1B]">
            Coach Pablo · v1.1
          </span>
        </div>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">
          Resumen ejecutivo de operaciones, agenda del día, ciclo de pagos y estado de alumnos.
        </p>
      </div>

      <!-- Quick Action Buttons -->
      <div class="flex items-center gap-2.5 flex-wrap">
        <button
          @click="copyPublicOnboardingLink"
          class="h-[40px] px-3.5 bg-white border border-[#DCD8CF] hover:border-[#CEA549] text-[#191A1D] rounded-[6px] text-[13px] font-semibold transition-all flex items-center gap-2 cursor-pointer shadow-2xs hover:bg-[#FAF9F6]"
          title="Copiar mensaje con link de onboarding y firma digital para WhatsApp"
        >
          <component :is="copiedOnboardingLink ? Check : Share2" class="w-4 h-4" :class="copiedOnboardingLink ? 'text-[#137333]' : 'text-[#8C6B1B]'" />
          <span>{{ copiedOnboardingLink ? '¡Enlace copiado!' : 'WhatsApp Alumno' }}</span>
        </button>

        <button
          @click="showRegisterModal = true"
          class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center gap-1.5 cursor-pointer shadow-xs hover:shadow"
        >
          <Plus class="w-4 h-4 stroke-[3]" />
          <span>Registrar alumno</span>
        </button>
      </div>
    </div>

    <!-- 4 High-Impact KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1: Alumnos Activos -->
      <div
        @click="router.push('/alumnos')"
        class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-5 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 hover:shadow-sm transition-all group overflow-hidden cursor-pointer"
      >
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Alumnos activos</span>
          <div class="w-8 h-8 rounded-lg bg-[#FAF9F6] border border-[#E2DFD8] flex items-center justify-center text-[#8C6B1B] group-hover:bg-[#CEA549] group-hover:text-black transition-colors">
            <Users class="w-4 h-4" />
          </div>
        </div>
        <div class="my-2">
          <span class="text-[32px] font-bold text-[#191A1D] leading-tight">
            {{ gymStore.activeStudentsCount }}
          </span>
          <span class="text-xs text-[#8E9097] ml-2 font-medium">alumnos registrados</span>
        </div>
        <span class="text-[12px] text-[#8C6B1B] font-semibold flex items-center gap-1 group-hover:underline">
          Ver directorio completo →
        </span>
      </div>

      <!-- Card 2: Recaudación Ciclo 1 al 5 -->
      <div
        @click="router.push('/membresias-pagos')"
        class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-5 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 hover:shadow-sm transition-all group overflow-hidden cursor-pointer"
      >
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Recaudado ciclo 1 al 5</span>
          <div class="w-8 h-8 rounded-lg bg-[#FAF9F6] border border-[#E2DFD8] flex items-center justify-center text-[#137333] group-hover:bg-[#137333] group-hover:text-white transition-colors">
            <CreditCard class="w-4 h-4" />
          </div>
        </div>
        <div class="my-2">
          <span class="text-[30px] font-bold text-[#137333] leading-tight">
            ${{ totalCollected.toLocaleString('es-CL') }}
          </span>
          <p class="text-[11px] text-[#8E9097] mt-0.5 font-medium">
            {{ paidPayments.length }} pagadas · {{ pendingPayments.length }} por cobrar (${{ totalPending.toLocaleString('es-CL') }})
          </p>
        </div>
        <span class="text-[12px] text-[#8C6B1B] font-semibold flex items-center gap-1 group-hover:underline">
          Gestionar pagos SumUp →
        </span>
      </div>

      <!-- Card 3: Sesiones y Ocupación Hoy -->
      <div
        @click="router.push('/agenda')"
        class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-5 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 hover:shadow-sm transition-all group overflow-hidden cursor-pointer"
      >
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Ocupación de hoy ({{ todayDayName }})</span>
          <div class="w-8 h-8 rounded-lg bg-[#FAF9F6] border border-[#E2DFD8] flex items-center justify-center text-[#8C6B1B] group-hover:bg-[#CEA549] group-hover:text-black transition-colors">
            <Calendar class="w-4 h-4" />
          </div>
        </div>
        <div class="my-2">
          <span class="text-[32px] font-bold text-[#191A1D] leading-tight">
            {{ todayOccupiedCupos }} / {{ todayTotalCupos }}
          </span>
          <span class="text-xs text-[#8E9097] ml-2 font-medium">cupos ocupados</span>
          <p class="text-[11px] text-[#137333] font-semibold mt-0.5">
            {{ todayFreeCupos }} cupos libres disponibles
          </p>
        </div>
        <span class="text-[12px] text-[#8C6B1B] font-semibold flex items-center gap-1 group-hover:underline">
          Ver agenda de 75 min →
        </span>
      </div>

      <!-- Card 4: Consentimientos & Firmas -->
      <div
        @click="router.push('/contratos')"
        class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-5 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 hover:shadow-sm transition-all group overflow-hidden cursor-pointer"
      >
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Consentimientos legales</span>
          <div class="w-8 h-8 rounded-lg bg-[#FAF9F6] border border-[#E2DFD8] flex items-center justify-center text-[#8C6B1B] group-hover:bg-[#CEA549] group-hover:text-black transition-colors">
            <ShieldCheck class="w-4 h-4" />
          </div>
        </div>
        <div class="my-2">
          <span class="text-[32px] font-bold text-[#191A1D] leading-tight">
            {{ gymStore.acceptedContractsCount }}
          </span>
          <span class="text-xs text-[#137333] font-bold ml-2">vigentes</span>
          <p class="text-[11px] text-[#B06000] font-semibold mt-0.5">
            {{ gymStore.pendingContractsCount }} pendientes de firma digital
          </p>
        </div>
        <span class="text-[12px] text-[#8C6B1B] font-semibold flex items-center gap-1 group-hover:underline">
          Ver consentimientos y firmas →
        </span>
      </div>
    </div>

    <!-- Main Operational Grid: Agenda del Día (Left 7 cols) & Cobranza SumUp (Right 5 cols) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left Column: Agenda Oficial de Hoy (7 cols) -->
      <div class="lg:col-span-7 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <div>
              <div class="flex items-center gap-2">
                <h2 class="text-[18px] font-bold text-[#191A1D]">
                  Agenda de hoy · {{ todayDayName }} {{ todayInfo.dayNumber }} {{ todayInfo.monthLabel }}
                </h2>
                <span class="px-2 py-0.5 rounded text-[11px] font-bold bg-[#FAF9F6] border border-[#E2DFD8] text-[#8C6B1B]">
                  75 min · Máx 3 alumnos
                </span>
              </div>
              <p class="text-[13px] text-[#6B6D74] mt-0.5">
                Turnos preestablecidos y asignación de cupos
              </p>
            </div>

            <button
              @click="router.push('/agenda')"
              class="text-[13px] text-[#8C6B1B] hover:text-[#614910] font-bold flex items-center gap-1 hover:underline cursor-pointer"
            >
              Ver agenda completa →
            </button>
          </div>

          <!-- Sunday (Gym Closed) -->
          <div v-if="todayDayName === 'Dom'" class="p-6 rounded-[8px] border border-dashed border-[#E2DFD8] bg-[#FAF9F6] text-center my-3">
            <div class="w-10 h-10 rounded-full bg-[#E2DFD8]/60 flex items-center justify-center mx-auto mb-2 text-[#6B6D74]">
              <Clock class="w-5 h-5 text-[#8C6B1B]" />
            </div>
            <h3 class="text-[14px] font-bold text-[#191A1D]">Gimnasio Cerrado (Domingo)</h3>
            <p class="text-[12px] text-[#6B6D74] mt-0.5">Día oficial de descanso semanal del Coach Pablo. Los entrenamientos se reanudan el Lunes a las 06:30 hrs.</p>
          </div>

          <!-- Case 1: When there ARE occupied sessions for today -->
          <div v-else-if="todayOccupiedSlots.length > 0" class="space-y-3">
            <div
              v-for="slot in todayOccupiedSlots"
              :key="slot.id"
              class="p-4 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 transition-all group bg-white cursor-pointer"
              @click="openOccupiedSlot(slot)"
            >
              <div class="flex items-center gap-4">
                <!-- Time Badge -->
                <div class="w-[72px] text-center py-2 px-1 bg-[#191A1D] rounded-[6px] border border-[#CEA549]/40 text-[12px] font-bold text-[#CEA549] shrink-0">
                  {{ slot.startTime }}
                </div>

                <div>
                  <div class="flex items-center gap-2 flex-wrap">
                    <p class="text-[14px] font-bold text-[#191A1D] group-hover:text-black">
                      {{ slot.isEvaluation ? 'Evaluación física' : 'Bloque de entrenamiento' }}
                    </p>

                    <!-- Cupo Badge -->
                    <span
                      class="px-2 py-0.5 rounded-[4px] text-[11px] font-bold"
                      :class="slot.isFull ? 'bg-[#FEF7E0] text-[#B06000]' : 'bg-[#E6F4EA] text-[#137333]'"
                    >
                      {{ slot.cuposCount }}/3 cupos {{ slot.isFull ? '· Lleno' : '' }}
                    </span>
                  </div>

                  <!-- Student Names Pills -->
                  <div class="flex items-center gap-1.5 mt-1.5 flex-wrap">
                    <span
                      v-for="name in slot.studentNames"
                      :key="name"
                      class="px-2 py-0.5 rounded text-[11px] font-semibold bg-[#F5F3EE] border border-[#E2DFD8] text-[#191A1D]"
                    >
                      {{ name }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Right: Cupo Dots (3 max) + Chevron -->
              <div class="flex items-center gap-3 self-end sm:self-center">
                <div class="flex items-center gap-1">
                  <span
                    v-for="i in 3"
                    :key="i"
                    class="w-2.5 h-2.5 rounded-full transition-colors"
                    :class="i <= slot.cuposCount ? 'bg-[#CEA549]' : 'bg-[#E2DFD8]'"
                  ></span>
                </div>
                <ChevronRight class="w-4 h-4 text-[#8E9097] group-hover:text-[#CEA549] transition-colors" />
              </div>
            </div>

            <!-- Summary of available slots below occupied ones -->
            <div v-if="todayAvailableSlots.length > 0" class="pt-2">
              <p class="text-[11px] text-[#6B6D74] font-medium flex items-center justify-between">
                <span>{{ todayAvailableSlots.length }} bloques disponibles hoy para agendar</span>
                <button
                  @click="router.push('/agenda')"
                  class="text-[#8C6B1B] hover:underline font-bold cursor-pointer"
                >
                  Ver en agenda semanal →
                </button>
              </p>
            </div>
          </div>

          <!-- Case 2: When there are ZERO sessions scheduled for today (e.g. Jueves) -->
          <div v-else class="space-y-4 my-1">
            <!-- Empty state banner -->
            <div class="p-5 rounded-[8px] border border-dashed border-[#DCD8CF] bg-[#FAF9F6] text-center">
              <div class="w-10 h-10 rounded-full bg-[#E2DFD8]/50 flex items-center justify-center mx-auto mb-2 text-[#8C6B1B]">
                <Calendar class="w-5 h-5" />
              </div>
              <h3 class="text-[15px] font-bold text-[#191A1D]">
                Sin sesiones agendadas para hoy ({{ todayDayName }} {{ todayInfo.dayNumber }} {{ todayInfo.monthLabel }})
              </h3>
              <p class="text-[12px] text-[#6B6D74] mt-1 max-w-md mx-auto leading-relaxed">
                No hay alumnos programados en este día. Todos los {{ todaySlots.length }} bloques oficiales de 75 min del Coach Pablo están 100% libres y disponibles (0 / {{ todayTotalCupos }} cupos ocupados).
              </p>
              <div class="mt-3.5 flex items-center justify-center gap-2.5 flex-wrap">
                <button
                  @click="openAssignForToday"
                  class="h-[34px] px-3.5 rounded-[6px] bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] font-bold text-[12px] flex items-center gap-1.5 transition-colors cursor-pointer shadow-xs"
                >
                  <Plus class="w-3.5 h-3.5 stroke-[3]" />
                  <span>Asignar alumno a bloque hoy</span>
                </button>
                <button
                  @click="router.push('/agenda')"
                  class="h-[34px] px-3 rounded-[6px] bg-white border border-[#DCD8CF] hover:border-[#CEA549] text-[#191A1D] font-semibold text-[12px] transition-colors cursor-pointer"
                >
                  Ver agenda semanal completa
                </button>
              </div>
            </div>

            <!-- List of today's available blocks with quick 1-click assign -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-[11px] font-bold uppercase tracking-wider text-[#6B6D74]">
                  Bloques oficiales disponibles hoy ({{ todaySlots.length }} horarios)
                </span>
                <span class="text-[11px] text-[#137333] font-semibold">
                  0/3 cupos · 100% disponible
                </span>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-[220px] overflow-y-auto pr-1">
                <div
                  v-for="slot in todaySlots"
                  :key="slot.id"
                  @click="openAssignForSlot(slot)"
                  class="p-2.5 rounded-[6px] border border-[#E2DFD8] hover:border-[#CEA549] hover:bg-[#FAF9F6] bg-white flex items-center justify-between transition-colors cursor-pointer group"
                >
                  <div class="flex items-center gap-2">
                    <span class="px-2 py-0.5 rounded bg-[#191A1D] text-[#CEA549] font-mono text-[11px] font-bold">
                      {{ slot.startTime }}
                    </span>
                    <span class="text-[12px] text-[#6B6D74] group-hover:text-[#191A1D] font-medium">
                      {{ slot.startTime }} – {{ slot.endTime }}
                    </span>
                  </div>
                  <span class="text-[11px] text-[#8C6B1B] font-bold group-hover:underline flex items-center gap-1">
                    + Asignar
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] flex items-center justify-between text-[12px] text-[#6B6D74]">
          <span>¿Necesitas asignar un alumno a un bloque libre?</span>
          <button
            @click="openAssignForToday"
            class="text-[#8C6B1B] font-bold hover:underline cursor-pointer flex items-center gap-1"
          >
            + Asignar a bloque
          </button>
        </div>
      </div>

      <!-- Right Column: Recaudación & Pasarela SumUp (5 cols) -->
      <div class="lg:col-span-5 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-[18px] font-bold text-[#191A1D]">Recaudación SumUp</h2>
              <p class="text-[13px] text-[#6B6D74] mt-0.5">Ciclo oficial del 1 al 5 · $45.000 mensual</p>
            </div>
            <button
              @click="router.push('/membresias-pagos')"
              class="text-[13px] text-[#8C6B1B] hover:text-[#614910] font-bold flex items-center gap-1 hover:underline cursor-pointer"
            >
              Ver detalle →
            </button>
          </div>

          <!-- Permanent SumUp Link Card -->
          <div class="p-4 rounded-[8px] bg-[#191A1D] text-white border border-[#CEA549]/40 mb-4 space-y-3 shadow-xs">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <CreditCard class="w-4 h-4 text-[#CEA549]" />
                <span class="text-xs font-bold uppercase tracking-wider text-[#CEA549]">Link Permanente SumUp</span>
              </div>
              <button
                @click="showConfigureSumupModal = true"
                class="text-[11px] text-[#A0A2A8] hover:text-white underline cursor-pointer"
              >
                Configurar
              </button>
            </div>

            <div class="p-2.5 rounded bg-[#24262B] border border-[#3A3D47] flex items-center justify-between gap-2 font-mono text-xs">
              <span class="truncate text-[#E2DFD8]">
                {{ gymStore.sumupPermanentLink || 'https://link.sumup.com/pay/inverfit' }}
              </span>
              <button
                @click="copyPermanentSumupLink"
                class="px-2.5 py-1 rounded bg-[#CEA549] hover:bg-[#BA923B] text-black font-bold text-[11px] transition-colors shrink-0 flex items-center gap-1 cursor-pointer"
              >
                <component :is="copiedSumupLink ? Check : Copy" class="w-3.5 h-3.5" />
                <span>{{ copiedSumupLink ? '¡Listo!' : 'Copiar' }}</span>
              </button>
            </div>

            <p class="text-[11px] text-[#8E9097] leading-relaxed">
              Enlace único para que los alumnos paguen con tarjeta de débito o crédito sin generar links individuales.
            </p>
          </div>

          <!-- Status of Pending Notices (Ciclo 1 al 5) -->
          <div class="space-y-3">
            <div class="p-3.5 rounded-[8px] border border-[#E2DFD8] bg-[#FAF9F6] flex items-center justify-between">
              <div>
                <p class="text-xs font-semibold text-[#6B6D74]">Membresías al día este mes</p>
                <p class="text-lg font-bold text-[#137333]">{{ paidPayments.length }} alumnos</p>
              </div>
              <span class="px-2.5 py-1 rounded text-xs font-bold bg-[#E6F4EA] text-[#137333]">
                Al día ✓
              </span>
            </div>

            <div class="p-3.5 rounded-[8px] border border-[#E2DFD8] bg-[#FEF7E0]/40 flex items-center justify-between">
              <div>
                <p class="text-xs font-semibold text-[#6B6D74]">Pagos pendientes ciclo 1 al 5</p>
                <p class="text-lg font-bold text-[#B06000]">{{ pendingPayments.length }} alumnos</p>
              </div>
              <button
                @click="showBatchRemindersModal = true"
                class="px-3 py-1.5 bg-[#B06000] hover:bg-[#924F00] text-white rounded text-xs font-bold transition-colors cursor-pointer shadow-xs flex items-center gap-1.5"
              >
                <Send class="w-3 h-3" />
                <span>Despacho 1 al 5</span>
              </button>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] text-[11px] text-[#8E9097]">
          Los pagos se pueden acreditar con el N° de operación y comprobante SumUp / transferencia.
        </div>
      </div>
    </div>

    <!-- Secondary Row: Evaluaciones Recientes (6 cols) & Consentimientos/Firmas (6 cols) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left: Evaluaciones Físicas Recientes (6 cols) -->
      <div class="lg:col-span-6 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <ClipboardCheck class="w-5 h-5 text-[#8C6B1B]" />
              <div>
                <h2 class="text-[17px] font-bold text-[#191A1D]">Evaluaciones & Baremos M-F</h2>
                <p class="text-[12px] text-[#6B6D74]">Últimos registros de dinamometría y signos vitales</p>
              </div>
            </div>

            <button
              @click="router.push('/evaluaciones')"
              class="text-[13px] text-[#8C6B1B] hover:text-[#614910] font-bold hover:underline cursor-pointer"
            >
              Ver todas →
            </button>
          </div>

          <div class="space-y-3">
            <div
              v-for="ev in latestEvaluations"
              :key="ev.id"
              class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] transition-all flex items-center justify-between cursor-pointer group"
              @click="router.push('/evaluaciones')"
            >
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-xs font-bold shrink-0">
                  {{ ev.studentName.charAt(0) }}
                </div>
                <div>
                  <p class="text-[13px] font-bold text-[#191A1D] group-hover:text-black">{{ ev.studentName }}</p>
                  <p class="text-[11px] text-[#8E9097]">
                    Fecha: {{ ev.date }} · {{ ev.type }}
                  </p>
                </div>
              </div>

              <!-- Metric Preview Chip -->
              <div class="text-right">
                <span class="px-2 py-0.5 rounded text-[11px] font-bold bg-[#E6F4EA] text-[#137333]">
                  {{ ev.status }}
                </span>
                <p class="text-[11px] text-[#6B6D74] mt-0.5 font-mono">
                  {{ ev.dynamometry ? `${ev.dynamometry.rightHandKg} kg prensil` : 'Registrada' }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] flex items-center justify-between text-[12px] text-[#6B6D74]">
          <span>Protocolo clínico y baremos nativos diferenciados</span>
          <button
            @click="router.push('/evaluaciones')"
            class="text-[#8C6B1B] font-bold hover:underline cursor-pointer"
          >
            Nueva evaluación +
          </button>
        </div>
      </div>

      <!-- Right: Contratos & Firmas Digitales (6 cols) -->
      <div class="lg:col-span-6 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <ShieldCheck class="w-5 h-5 text-[#8C6B1B]" />
              <div>
                <h2 class="text-[17px] font-bold text-[#191A1D]">Consentimientos & Firmas</h2>
                <p class="text-[12px] text-[#6B6D74]">Respaldo legal y reglamento interno v1.0</p>
              </div>
            </div>

            <button
              @click="router.push('/contratos')"
              class="text-[13px] text-[#8C6B1B] hover:text-[#614910] font-bold hover:underline cursor-pointer"
            >
              Ver todos →
            </button>
          </div>

          <div class="space-y-3">
            <div
              v-for="contract in latestContracts"
              :key="contract.id"
              class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] transition-all flex items-center justify-between cursor-pointer group"
              @click="router.push('/contratos')"
            >
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-xs font-bold shrink-0">
                  {{ contract.studentName.charAt(0) }}
                </div>
                <div>
                  <div class="flex items-center gap-1.5">
                    <p class="text-[13px] font-bold text-[#191A1D] group-hover:text-black">{{ contract.studentName }}</p>
                    <span
                      v-if="getStudent(contract.studentId)?.onboarding?.digitalSignature"
                      class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#E6F4EA] text-[#137333]"
                    >
                      ✍️
                    </span>
                  </div>
                  <p class="text-[11px] text-[#8E9097] leading-none mt-1">
                    {{ contract.documentType }}
                  </p>
                </div>
              </div>

              <div class="text-right">
                <span
                  class="px-2 py-0.5 rounded text-[11px] font-semibold"
                  :class="contract.status === 'Aceptado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                >
                  {{ contract.status }}
                </span>
                <p class="text-[11px] text-[#8E9097] mt-0.5">
                  {{ contract.date }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] flex items-center justify-between text-[12px] text-[#6B6D74]">
          <span>Firmas capturadas vía pantalla táctil en onboarding</span>
          <button
            @click="copyPublicOnboardingLink"
            class="text-[#8C6B1B] font-bold hover:underline cursor-pointer flex items-center gap-1"
          >
            <Share2 class="w-3.5 h-3.5" />
            <span>Compartir enlace</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Bottom Dock: Accesos Rápidos -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-5 shadow-2xs">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h3 class="text-[14px] font-bold text-[#191A1D]">Accesos directos del sistema</h3>
          <p class="text-[12px] text-[#6B6D74]">Navegación rápida a los módulos de Coach Pablo</p>
        </div>

        <div class="flex flex-wrap items-center gap-2.5">
          <button
            @click="router.push('/agenda')"
            class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549] transition-all flex items-center gap-2 cursor-pointer shadow-2xs"
          >
            <Calendar class="w-4 h-4 text-[#8C6B1B]" />
            <span>Agenda 75 min</span>
          </button>

          <button
            @click="router.push('/alumnos')"
            class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549] transition-all flex items-center gap-2 cursor-pointer shadow-2xs"
          >
            <Users class="w-4 h-4 text-[#8C6B1B]" />
            <span>Alumnos</span>
          </button>

          <button
            @click="router.push('/membresias-pagos')"
            class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549] transition-all flex items-center gap-2 cursor-pointer shadow-2xs"
          >
            <CreditCard class="w-4 h-4 text-[#8C6B1B]" />
            <span>Cobranza SumUp</span>
          </button>

          <button
            @click="router.push('/evaluaciones')"
            class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549] transition-all flex items-center gap-2 cursor-pointer shadow-2xs"
          >
            <ClipboardCheck class="w-4 h-4 text-[#8C6B1B]" />
            <span>Evaluaciones</span>
          </button>

          <button
            @click="router.push('/contratos')"
            class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549] transition-all flex items-center gap-2 cursor-pointer shadow-2xs"
          >
            <FileText class="w-4 h-4 text-[#8C6B1B]" />
            <span>Contratos</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <RegisterStudentModal
      :open="showRegisterModal"
      @close="showRegisterModal = false"
    />

    <BlockDetailsModal
      :open="showBlockModal"
      :block="activeBlock"
      :occurrence="activeOccurrence"
      @close="showBlockModal = false"
      @deleted="handleScheduleChanged"
      @updated="handleScheduleChanged"
      @cancelled="handleScheduleChanged"
      @add-student="handleAddStudentToBlock"
    />

    <AssignScheduleModal
      :open="showAssignModal"
      :block-id="assignBlockId || undefined"
      :initial-day="assignInitialDay"
      :initial-time="assignInitialTime"
      :initial-end-time="assignInitialEndTime"
      :initial-date="assignInitialDate"
      @close="showAssignModal = false"
      @assigned="handleScheduleChanged"
    />

    <BatchRemindersModal
      :open="showBatchRemindersModal"
      :pending-payments="pendingPayments"
      @close="showBatchRemindersModal = false"
    />

    <ConfigureSumupLinkModal
      :open="showConfigureSumupModal"
      @close="showConfigureSumupModal = false"
    />
  </div>
</template>
