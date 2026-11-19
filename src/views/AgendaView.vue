<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import {
  getWeekDates,
  timeToMinutes,
  PABLO_FIXED_PRESET_BLOCKS,
  PABLO_SATURDAY_EVALUATION_BLOCKS,
  INVERFIT_MAX_STUDENTS_PER_BLOCK,
  type WeekDayInfo,
  type FixedPresetBlock
} from '@/services/schedule'
import type { DayOfWeek, ScheduleBlock, SessionOccurrence } from '@/types'
import {
  Clock,
  AlertCircle,
  RefreshCw,
  CheckCircle2,
  Plus
} from 'lucide-vue-next'
import AssignScheduleModal from '@/components/modals/AssignScheduleModal.vue'
import BlockDetailsModal from '@/components/modals/BlockDetailsModal.vue'

const gymStore = useGymStore()

// State
const weekOffset = ref(0)
const selectedDay = ref<DayOfWeek>('Lun')
const showAssignModal = ref(false)
const showBlockModal = ref(false)
const isRefreshing = ref(false)
const weekError = ref<string | null>(null)

const activeBlockId = ref<string | null>(null)
const activeBlock = computed(() => {
  if (!activeBlockId.value) return null
  return gymStore.schedule.find(b => b.id === activeBlockId.value) || null
})
const activeOccurrence = ref<SessionOccurrence | null>(null)
const prefillBlockId = ref<string | undefined>(undefined)
const prefillDay = ref<DayOfWeek>('Lun')
const prefillTime = ref('08:00')
const prefillEndTime = ref<string | undefined>(undefined)
const prefillDate = ref<string | undefined>(undefined)

// Week dates calculation
const weekDays = computed<WeekDayInfo[]>(() => {
  return getWeekDates(weekOffset.value)
})

// Current week label
const weekRangeLabel = computed(() => {
  if (weekDays.value.length === 0) return ''
  const first = weekDays.value[0]
  const last = weekDays.value[6]
  return `${first.dayNumber} ${first.monthLabel} — ${last.dayNumber} ${last.monthLabel}`
})

const weekOffsetLabel = computed(() => {
  if (weekOffset.value === 0) return 'Semana actual'
  if (weekOffset.value > 0) return `Semana +${weekOffset.value}`
  return `Semana ${weekOffset.value}`
})

async function refreshAgenda() {
  if (weekDays.value.length === 0) return
  const from = weekDays.value[0].dateStr
  const to = weekDays.value[6].dateStr
  isRefreshing.value = true
  weekError.value = null
  try {
    if (gymStore.isScheduleApiMode()) {
      await Promise.all([
        gymStore.fetchWeekOccurrences(from, to),
        gymStore.loadSchedule(from, to)
      ])
    } else {
      await gymStore.loadSchedule(from, to)
    }
  } catch (err: unknown) {
    weekError.value = err instanceof Error ? err.message : 'Error al cargar la agenda semanal'
  } finally {
    isRefreshing.value = false
  }
}

onMounted(() => {
  refreshAgenda()
})

watch(weekOffset, () => {
  refreshAgenda()
})

// Calculate occurrences for the current week taking exceptions into account
const weekOccurrences = computed<SessionOccurrence[]>(() => {
  return gymStore.getWeekOccurrences(weekDays.value)
})

// Occurrences grouped by day
function getOccurrencesForDay(day: DayOfWeek): SessionOccurrence[] {
  return weekOccurrences.value.filter(o => o.day === day)
}

// Student names helper
function getStudentNames(studentIds: string[]): string {
  const students = gymStore.students.filter(s => studentIds.includes(s.id))
  if (students.length === 0) return 'Sin alumnos asignados'
  if (students.length <= 2) return students.map(s => s.fullName).join(', ')
  return `${students[0].fullName} +${students.length - 1}`
}

// ==============================================================================
// ESTRUCTURA DE BLOQUES FIJOS INVERFIT (06:30 a 22:00)
// ==============================================================================

export interface DayScheduleSlot {
  id: string
  presetId?: string
  startTime: string
  endTime: string
  durationMinutes: number
  shift: 'mañana' | 'tarde'
  isEvaluation: boolean
  isOccupied: boolean
  occurrence?: SessionOccurrence
  studentIds: string[]
  studentNames: string
  cuposCount: number
  isFull: boolean
  isException?: boolean
  notes?: string
}

/**
 * Obtiene los bloques fijos oficiales de Pablo para un día determinado,
 * vinculando las sesiones efectivamente asignadas o marcándolos como disponibles.
 */
function getSlotsForDay(day: DayOfWeek): DayScheduleSlot[] {
  if (day === 'Dom') {
    return []
  }

  const isSat = day === 'Sáb'
  const presets: FixedPresetBlock[] = isSat
    ? PABLO_SATURDAY_EVALUATION_BLOCKS
    : PABLO_FIXED_PRESET_BLOCKS

  const dayOccurrences = getOccurrencesForDay(day)
  const matchedOccIds = new Set<string>()

  // 1. Mapear cada bloque preestablecido oficial
  const slots: DayScheduleSlot[] = presets.map((preset) => {
    // Buscar si hay sesión asignada en este horario de inicio
    const occ = dayOccurrences.find(o => o.startTime === preset.startTime)
    if (occ) {
      matchedOccIds.add(occ.instanceId)
      return {
        id: occ.instanceId,
        presetId: preset.id,
        startTime: occ.startTime,
        endTime: occ.endTime,
        durationMinutes: occ.durationMinutes || 75,
        shift: preset.shift,
        isEvaluation: occ.type === 'Evaluación' || isSat,
        isOccupied: occ.studentIds && occ.studentIds.length > 0,
        occurrence: occ,
        studentIds: occ.studentIds || [],
        studentNames: getStudentNames(occ.studentIds || []),
        cuposCount: (occ.studentIds || []).length,
        isFull: (occ.studentIds || []).length >= INVERFIT_MAX_STUDENTS_PER_BLOCK,
        isException: occ.isException,
        notes: occ.notes
      }
    }

    // Bloque libre / disponible
    return {
      id: `slot-${day}-${preset.id}`,
      presetId: preset.id,
      startTime: preset.startTime,
      endTime: preset.endTime,
      durationMinutes: 75,
      shift: preset.shift,
      isEvaluation: isSat,
      isOccupied: false,
      occurrence: undefined,
      studentIds: [],
      studentNames: '',
      cuposCount: 0,
      isFull: false,
      isException: false
    }
  })

  // 2. Agregar cualquier sesión creada fuera de presets estándar
  const extraOccs = dayOccurrences.filter(o => !matchedOccIds.has(o.instanceId))
  for (const occ of extraOccs) {
    slots.push({
      id: occ.instanceId,
      startTime: occ.startTime,
      endTime: occ.endTime,
      durationMinutes: occ.durationMinutes || 75,
      shift: timeToMinutes(occ.startTime) < 14 * 60 ? 'mañana' : 'tarde',
      isEvaluation: occ.type === 'Evaluación' || isSat,
      isOccupied: occ.studentIds && occ.studentIds.length > 0,
      occurrence: occ,
      studentIds: occ.studentIds || [],
      studentNames: getStudentNames(occ.studentIds || []),
      cuposCount: (occ.studentIds || []).length,
      isFull: (occ.studentIds || []).length >= INVERFIT_MAX_STUDENTS_PER_BLOCK,
      isException: occ.isException,
      notes: occ.notes
    })
  }

  return slots.sort((a, b) => timeToMinutes(a.startTime) - timeToMinutes(b.startTime))
}

/**
 * Estadísticas de ocupación para un día
 */
function getDayStats(day: DayOfWeek) {
  const slots = getSlotsForDay(day)
  const totalSlots = slots.length
  const occupiedSlots = slots.filter(s => s.isOccupied).length
  const totalOccupiedCupos = slots.reduce((acc, s) => acc + s.cuposCount, 0)
  const maxCupos = totalSlots * INVERFIT_MAX_STUDENTS_PER_BLOCK
  const uniqueStudents = new Set(slots.flatMap(s => s.studentIds)).size
  const occupancyPercentage = maxCupos > 0 ? Math.round((totalOccupiedCupos / maxCupos) * 100) : 0

  return {
    totalSlots,
    occupiedSlots,
    freeSlots: totalSlots - occupiedSlots,
    totalOccupiedCupos,
    maxCupos,
    uniqueStudents,
    occupancyPercentage
  }
}

// Estadísticas de la semana completa
const weekTotalStats = computed(() => {
  let occupiedCupos = 0
  let totalCupos = 0
  let activeSessions = 0

  for (const dayInfo of weekDays.value) {
    if (dayInfo.day !== 'Dom') {
      const stats = getDayStats(dayInfo.day)
      occupiedCupos += stats.totalOccupiedCupos
      totalCupos += stats.maxCupos
      activeSessions += stats.occupiedSlots
    }
  }

  const percentage = totalCupos > 0 ? Math.round((occupiedCupos / totalCupos) * 100) : 0
  return { occupiedCupos, totalCupos, activeSessions, percentage }
})

// Navigation methods
function prevWeek() {
  weekOffset.value--
}

function nextWeek() {
  weekOffset.value++
}

function resetWeek() {
  weekOffset.value = 0
}

function selectDay(day: DayOfWeek) {
  selectedDay.value = day
}

// Modal open handlers: directly called from the slots!
function openNewSessionModal(day?: DayOfWeek, time?: string, dateStr?: string) {
  prefillBlockId.value = undefined
  prefillEndTime.value = undefined
  prefillDay.value = day || selectedDay.value
  prefillTime.value = time || '08:00'
  prefillDate.value = dateStr || weekDays.value.find(w => w.day === (day || selectedDay.value))?.dateStr
  showAssignModal.value = true
}

function openSessionDetails(occurrence: SessionOccurrence) {
  activeOccurrence.value = occurrence
  activeBlockId.value = occurrence.blockId
  showBlockModal.value = true
}

function handleAddStudentFromBlock(block: ScheduleBlock) {
  showBlockModal.value = false
  prefillBlockId.value = block.id
  prefillDay.value = block.day
  prefillTime.value = block.startTime
  prefillEndTime.value = block.endTime
  showAssignModal.value = true
}
</script>

<template>
  <div class="space-y-5 max-w-[1550px] mx-auto pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Horario Semanal INVERFIT</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">
          Lunes a Viernes 06:30 a 22:00 hrs (9 bloques de 75 min) · Sábados Evaluaciones 06:30 a 13:00 hrs · Haz clic en cualquier bloque para agendar o ver alumnos
        </p>
      </div>

      <div class="flex items-center gap-3">
        <!-- Button: Refresh -->
        <button
          @click="refreshAgenda"
          :disabled="isRefreshing"
          class="h-[38px] px-3.5 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-[6px] text-[13px] font-semibold transition-all flex items-center gap-2 cursor-pointer shadow-2xs disabled:opacity-50"
          title="Recargar sesiones de la semana"
        >
          <RefreshCw class="w-4 h-4 text-[#8C6B1B]" :class="{ 'animate-spin': isRefreshing }" />
          <span class="hidden sm:inline">Recargar</span>
        </button>
      </div>
    </div>

    <!-- Error Alert Banner -->
    <div
      v-if="weekError"
      class="p-4 bg-rose-50 border border-rose-200 text-rose-800 rounded-lg text-sm flex items-center justify-between"
    >
      <div class="flex items-center gap-2">
        <AlertCircle class="w-5 h-5 flex-shrink-0 text-rose-600" />
        <span>{{ weekError }}</span>
      </div>
      <button
        @click="refreshAgenda"
        class="px-3 py-1 bg-white border border-rose-300 text-rose-700 hover:bg-rose-100 rounded text-xs font-bold cursor-pointer transition-colors"
      >
        Reintentar conexión
      </button>
    </div>

    <!-- Toolbar: Week Switcher & Clean Legend -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] px-4 py-3 flex flex-wrap items-center justify-between gap-4 shadow-2xs">
      <div class="flex items-center gap-2.5">
        <span class="px-3 py-1 bg-[#191A1D] border border-[#282A30] rounded-[6px] text-[13px] font-bold text-[#CEA549]">
          {{ weekOffsetLabel }}
        </span>

        <span class="text-sm font-bold text-[#191A1D] px-1">
          {{ weekRangeLabel }}
        </span>

        <button
          @click="resetWeek"
          class="h-[32px] px-3 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] rounded-[6px] text-[12px] font-semibold text-[#191A1D] transition-colors cursor-pointer"
        >
          Hoy
        </button>

        <div class="flex items-center border border-[#DCD8CF] rounded-[6px] overflow-hidden bg-white">
          <button
            @click="prevWeek"
            class="h-[32px] px-3 hover:bg-[#FAF9F6] text-[#191A1D] text-sm border-r border-[#DCD8CF] transition-colors cursor-pointer"
            title="Semana anterior"
          >
            ‹
          </button>
          <button
            @click="nextWeek"
            class="h-[32px] px-3 hover:bg-[#FAF9F6] text-[#191A1D] text-sm transition-colors cursor-pointer"
            title="Semana siguiente"
          >
            ›
          </button>
        </div>

        <span class="text-xs text-[#8E9097] ml-2 hidden md:inline">
          Ocupación semanal: <strong>{{ weekTotalStats.occupiedCupos }} de {{ weekTotalStats.totalCupos }} cupos</strong> ({{ weekTotalStats.percentage }}%)
        </span>
      </div>

      <!-- Clean Legend -->
      <div class="flex items-center gap-5 text-xs text-[#6B6D74]">
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-[#191A1D] border border-[#CEA549]"></span>
          <span>Bloque con alumnos (1 a 3)</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
          <span>Bloque disponible (0/3)</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-[#CEA549]"></span>
          <span>Evaluaciones (Sábados)</span>
        </div>
      </div>
    </div>

    <!-- MAIN FULL-WIDTH TIMETABLE GRID (7 COLUMNS) -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 sm:p-5 shadow-2xs space-y-4">
      <!-- 7 Days Grid Matrix -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
        <div
          v-for="dayInfo in weekDays"
          :key="dayInfo.day"
          @click="selectDay(dayInfo.day)"
          class="flex flex-col rounded-lg border transition-all overflow-hidden min-h-[460px]"
          :class="[
            selectedDay === dayInfo.day
              ? 'border-[#CEA549] ring-2 ring-[#CEA549]/20 bg-[#FAF9F6]/40'
              : 'border-[#E2DFD8] hover:border-[#DCD8CF] bg-white'
          ]"
        >
          <!-- Day Column Header -->
          <div
            class="p-3 border-b transition-colors flex flex-col gap-1"
            :class="[
              selectedDay === dayInfo.day ? 'bg-[#191A1D] text-white' : 'bg-[#FAF9F6] text-[#191A1D]'
            ]"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold uppercase tracking-wider">
                {{ dayInfo.day }}
              </span>
              <span
                v-if="dayInfo.day === 'Sáb'"
                class="px-1.5 py-0.2 bg-[#CEA549]/20 text-[#8C6B1B] text-[9px] font-bold rounded border border-[#CEA549]/30"
              >
                Evaluaciones
              </span>
              <span
                v-else-if="dayInfo.isToday"
                class="px-1.5 py-0.2 bg-[#CEA549] text-[#0B0B0D] text-[10px] font-bold rounded"
              >
                Hoy
              </span>
            </div>

            <div class="flex items-baseline justify-between mt-0.5">
              <div class="flex items-baseline gap-1">
                <span class="text-xl font-bold leading-none">
                  {{ dayInfo.dayNumber }}
                </span>
                <span class="text-xs opacity-75">
                  {{ dayInfo.monthLabel }}
                </span>
              </div>

              <!-- Cupos summary badge -->
              <span
                v-if="dayInfo.day !== 'Dom'"
                class="text-[10px] px-1.5 py-0.2 rounded font-semibold"
                :class="selectedDay === dayInfo.day
                  ? 'bg-[#CEA549]/20 text-[#CEA549]'
                  : 'bg-white border border-[#DCD8CF] text-[#6B6D74]'"
              >
                {{ getDayStats(dayInfo.day).totalOccupiedCupos }}/{{ getDayStats(dayInfo.day).maxCupos }}
              </span>
            </div>

            <!-- Operating hours subtitle -->
            <div class="text-[10px] leading-tight opacity-75 mt-0.5">
              <span v-if="dayInfo.day === 'Dom'">Cerrado</span>
              <span v-else-if="dayInfo.day === 'Sáb'">06:30 – 13:00 hrs</span>
              <span v-else>06:30 – 22:00 hrs</span>
            </div>
          </div>

          <!-- Day Column Body: Full Timeline of Blocks -->
          <div class="p-2 flex-1 flex flex-col space-y-2 overflow-y-auto max-h-[720px]">
            <!-- DOMINGO: CERRADO -->
            <div
              v-if="dayInfo.day === 'Dom'"
              class="flex-1 flex flex-col items-center justify-center p-4 text-center text-xs text-[#8E9097] border border-dashed border-[#ECE8E1] rounded bg-[#FAF9F6]/50 min-h-[300px]"
            >
              <Clock class="w-6 h-6 text-[#8E9097] mb-2 opacity-40" />
              <span class="font-bold text-[#191A1D]">Gimnasio Cerrado</span>
              <span class="text-[11px] text-[#8E9097] mt-1">Día de descanso semanal</span>
            </div>

            <!-- LUNES A SÁBADO: LISTA DE BLOQUES OFICIALES DE PABLO -->
            <template v-else>
              <div
                v-for="slot in getSlotsForDay(dayInfo.day)"
                :key="slot.id"
                class="transition-all rounded-md text-left"
              >
                <!-- 1. BLOQUE OCUPADO (CON ALUMNOS) -->
                <div
                  v-if="slot.isOccupied && slot.occurrence"
                  @click.stop="openSessionDetails(slot.occurrence)"
                  class="p-2.5 rounded-lg border text-left cursor-pointer shadow-2xs group transition-all"
                  :class="[
                    slot.isException
                      ? 'bg-amber-50/80 border-amber-300 hover:border-amber-500'
                      : 'bg-white border-[#DCD8CF] hover:border-[#CEA549] hover:shadow-xs'
                  ]"
                >
                  <!-- Time & Duration -->
                  <div class="flex items-center justify-between">
                    <span class="text-[11px] font-bold text-[#191A1D] flex items-center gap-1">
                      <Clock class="w-3 h-3 text-[#CEA549]" />
                      {{ slot.startTime }} – {{ slot.endTime }}
                    </span>
                    <span class="text-[10px] px-1 rounded font-semibold bg-[#FAF9F6] border border-[#ECE8E1] text-[#6B6D74]">
                      {{ slot.durationMinutes }}m
                    </span>
                  </div>

                  <!-- Type & Cupos Badges -->
                  <div class="flex items-center justify-between mt-1.5 gap-1">
                    <span
                      v-if="slot.isEvaluation"
                      class="text-[9px] font-bold px-1.5 py-0.2 rounded bg-[#CEA549]/15 text-[#8C6B1B] border border-[#CEA549]/30"
                    >
                      Eval
                    </span>
                    <span
                      v-else
                      class="text-[9px] font-semibold px-1.5 py-0.2 rounded bg-neutral-100 text-[#6B6D74]"
                    >
                      {{ slot.studentIds.length > 1 ? 'Grupal' : 'Individual' }}
                    </span>

                    <span
                      class="text-[9px] font-bold px-1.5 py-0.2 rounded ml-auto"
                      :class="slot.isFull ? 'bg-amber-100 text-amber-900 border border-amber-300' : 'bg-emerald-50 text-emerald-800 border border-emerald-200'"
                    >
                      {{ slot.cuposCount }}/3
                      <span v-if="slot.isFull">· Lleno</span>
                    </span>
                  </div>

                  <!-- Student names -->
                  <p class="text-[11px] font-bold text-[#191A1D] truncate leading-tight mt-1.5 group-hover:text-black">
                    {{ slot.studentNames }}
                  </p>
                </div>

                <!-- 2. BLOQUE LIBRE / DISPONIBLE (CLIC PARA AGENDAR) -->
                <div
                  v-else
                  @click.stop="openNewSessionModal(dayInfo.day, slot.startTime, dayInfo.dateStr)"
                  class="p-2.5 rounded-lg border border-dashed border-[#DCD8CF] bg-[#FAF9F6]/60 hover:bg-white hover:border-[#CEA549] text-left cursor-pointer group transition-all"
                >
                  <div class="flex items-center justify-between">
                    <span class="text-[11px] font-medium text-[#6B6D74] group-hover:text-[#191A1D] flex items-center gap-1">
                      <Clock class="w-3 h-3 text-[#8E9097] group-hover:text-[#CEA549]" />
                      {{ slot.startTime }} – {{ slot.endTime }}
                    </span>
                    <span class="text-[9px] px-1.5 py-0.2 rounded font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                      0/3 cupos
                    </span>
                  </div>

                  <div class="flex items-center justify-between mt-1 text-[11px] text-[#8E9097] group-hover:text-[#8C6B1B]">
                    <span>Disponible</span>
                    <span class="font-bold flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity text-[#CEA549]">
                      <Plus class="w-3 h-3" /> Agendar
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- Footnote matching Pablo's philosophy -->
      <div class="p-3.5 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] text-xs text-[#6B6D74] flex items-start gap-2">
        <CheckCircle2 class="w-4 h-4 text-[#8C6B1B] flex-shrink-0 mt-0.5" />
        <p class="leading-relaxed">
          <strong>Protocolo Oficial INVERFIT:</strong> Cada bloque de entrenamiento dura 75 minutos exactos y cuenta con un máximo de 3 cupos simultáneos para garantizar técnica y carga controlada. Haz clic directamente sobre cualquier bloque libre para agendar alumnos o sobre un bloque ocupado para ver y editar su detalle.
        </p>
      </div>
    </div>

    <!-- Modals -->
    <AssignScheduleModal
      :open="showAssignModal"
      :block-id="prefillBlockId"
      :initial-day="prefillDay"
      :initial-time="prefillTime"
      :initial-end-time="prefillEndTime"
      :initial-date="prefillDate"
      @close="showAssignModal = false"
      @assigned="refreshAgenda"
    />

    <BlockDetailsModal
      :open="showBlockModal"
      :block="activeBlock"
      :occurrence="activeOccurrence"
      @close="showBlockModal = false"
      @deleted="refreshAgenda"
      @cancelled="refreshAgenda"
      @restored="refreshAgenda"
      @updated="refreshAgenda"
      @add-student="handleAddStudentFromBlock"
    />
  </div>
</template>
