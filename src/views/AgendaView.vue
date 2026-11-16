<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import {
  getWeekDates,
  timeToMinutes,
  type WeekDayInfo
} from '@/services/schedule'
import type { DayOfWeek, ScheduleBlock, SessionOccurrence } from '@/types'
import {
  Plus,
  Clock,
  Users,
  Settings,
  ChevronRight,
  AlertCircle,
  RefreshCw
} from 'lucide-vue-next'
import AssignScheduleModal from '@/components/modals/AssignScheduleModal.vue'
import BlockDetailsModal from '@/components/modals/BlockDetailsModal.vue'
import ConfigureAvailabilityModal from '@/components/modals/ConfigureAvailabilityModal.vue'

const gymStore = useGymStore()

// State
const weekOffset = ref(0)
const selectedDay = ref<DayOfWeek>('Lun')
const showAssignModal = ref(false)
const showBlockModal = ref(false)
const showAvailabilityModal = ref(false)
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

// Availability for a day
function getAvailabilityForDay(day: DayOfWeek) {
  const found = gymStore.availability.find(a => a.day === day)
  if (found) {
    return {
      ...found,
      intervals: Array.isArray(found.intervals) ? found.intervals : []
    }
  }
  return { day, enabled: false, intervals: [] }
}

// Selected day info
const selectedDayInfo = computed(() => {
  return weekDays.value.find(w => w.day === selectedDay.value) || weekDays.value[0]
})

const selectedDayAvailability = computed(() => {
  return getAvailabilityForDay(selectedDay.value)
})

const selectedDayOccurrences = computed(() => {
  return getOccurrencesForDay(selectedDay.value)
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

// Modal open handlers
function openNewSessionModal(day?: DayOfWeek, time?: string, dateStr?: string) {
  prefillBlockId.value = undefined
  prefillEndTime.value = undefined
  prefillDay.value = day || selectedDay.value
  prefillTime.value = time || '08:00'
  prefillDate.value = dateStr || weekDays.value.find(w => w.day === (day || selectedDay.value))?.dateStr
  showAssignModal.value = true
}

function openConfigureAvailability(day?: DayOfWeek) {
  if (day) selectedDay.value = day
  showAvailabilityModal.value = true
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

// Student names helper
function getStudentNames(studentIds: string[]): string {
  const students = gymStore.students.filter(s => studentIds.includes(s.id))
  if (students.length === 0) return 'Sin alumnos'
  if (students.length <= 2) return students.map(s => s.fullName).join(', ')
  return `${students[0].fullName} +${students.length - 1}`
}

// Computed hours summary for contextual panel
const dayTotalAvailableHours = computed(() => {
  const intervals = selectedDayAvailability.value.intervals || []
  if (!selectedDayAvailability.value.enabled || intervals.length === 0) return 0
  const totalM = intervals.reduce((acc, i) => acc + (timeToMinutes(i.endTime) - timeToMinutes(i.startTime)), 0)
  return (totalM / 60).toFixed(1)
})

const dayScheduledHours = computed(() => {
  const totalM = selectedDayOccurrences.value.reduce((acc, o) => acc + o.durationMinutes, 0)
  return (totalM / 60).toFixed(1)
})
</script>

<template>
  <div class="space-y-6 max-w-[1300px] mx-auto pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Agenda Semanal</h1>
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold tracking-wide uppercase"
            :class="gymStore.scheduleDataSource === 'api' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-amber-50 text-amber-800 border border-amber-200'"
          >
            <span class="w-1.5 h-1.5 rounded-full" :class="gymStore.scheduleDataSource === 'api' ? 'bg-emerald-500' : 'bg-amber-500'"></span>
            Fuente: {{ gymStore.scheduleDataSource === 'api' ? 'API Go + MongoDB' : 'Demostración local (Mock)' }}
          </span>
        </div>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">
          Gestión flexible de sesiones y configuración manual de disponibilidad por día.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <!-- Button: Refresh -->
        <button
          @click="refreshAgenda"
          :disabled="isRefreshing"
          class="h-[40px] px-3.5 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-[6px] text-[13px] font-semibold transition-all flex items-center gap-2 cursor-pointer shadow-2xs disabled:opacity-50"
          title="Recargar agenda semanal"
        >
          <RefreshCw class="w-4 h-4 text-[#8C6B1B]" :class="{ 'animate-spin': isRefreshing }" />
          <span class="hidden sm:inline">Recargar</span>
        </button>

        <!-- Button: Configure Availability -->
        <button
          @click="openConfigureAvailability()"
          class="h-[40px] px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-[6px] text-[13px] font-bold transition-all flex items-center gap-2 cursor-pointer shadow-2xs hover:border-[#CEA549]/60"
        >
          <Settings class="w-4 h-4 text-[#8C6B1B]" />
          Configurar disponibilidad
        </button>

        <!-- Button: Program Session -->
        <button
          @click="openNewSessionModal(selectedDay, undefined, selectedDayInfo.dateStr)"
          class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs hover:shadow"
        >
          <Plus class="w-4 h-4" />
          + Programar sesión
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

    <!-- Toolbar: Week Switcher & Legend -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-wrap items-center justify-between gap-4 shadow-2xs">
      <div class="flex items-center gap-2.5">
        <span class="px-3 py-1.5 bg-[#191A1D] border border-[#282A30] rounded-[6px] text-[13px] font-bold text-[#CEA549]">
          {{ weekOffsetLabel }}
        </span>

        <span class="text-sm font-bold text-[#191A1D] px-2">
          {{ weekRangeLabel }}
        </span>

        <button
          @click="resetWeek"
          class="h-[34px] px-3 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] rounded-[6px] text-[13px] font-semibold text-[#191A1D] transition-colors cursor-pointer"
        >
          Hoy
        </button>

        <div class="flex items-center border border-[#DCD8CF] rounded-[6px] overflow-hidden bg-white">
          <button
            @click="prevWeek"
            class="h-[34px] px-3 hover:bg-[#FAF9F6] text-[#191A1D] text-sm border-r border-[#DCD8CF] transition-colors cursor-pointer"
            title="Semana anterior"
          >
            ‹
          </button>
          <button
            @click="nextWeek"
            class="h-[34px] px-3 hover:bg-[#FAF9F6] text-[#191A1D] text-sm transition-colors cursor-pointer"
            title="Semana siguiente"
          >
            ›
          </button>
        </div>
      </div>

      <!-- Legend -->
      <div class="flex items-center gap-4 text-xs text-[#6B6D74]">
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-[#191A1D] border border-[#CEA549]"></span>
          <span>Sesión programada</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-amber-400"></span>
          <span>Excepción semanal</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-[#137333]"></span>
          <span>Disponibilidad activa</span>
        </div>
      </div>
    </div>

    <!-- MAIN TWO-COLUMN LAYOUT: WEEKLY COLUMNS (LEFT 8 COLS) + CONTEXTUAL PANEL (RIGHT 4 COLS) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <!-- LEFT: 7 DAY COLUMNS (8 cols on lg) -->
      <div class="lg:col-span-8 bg-white rounded-[10px] border border-[#E2DFD8] p-4 sm:p-5 shadow-2xs space-y-4">
        <div class="flex items-center justify-between border-b pb-3 border-[#ECE8E1]">
          <div>
            <h2 class="text-[17px] font-bold text-[#191A1D]">Visor Semanal de Sesiones</h2>
            <p class="text-xs text-[#6B6D74] mt-0.5">
              Consulta las sesiones efectivamente asignadas y los espacios libres de cada día.
            </p>
          </div>
          <span class="text-xs text-[#8E9097] hidden sm:inline">
            Haz clic en un día para ver su detalle en el panel derecho.
          </span>
        </div>

        <!-- 7 Days Grid Matrix -->
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-7 gap-3">
          <div
            v-for="dayInfo in weekDays"
            :key="dayInfo.day"
            @click="selectDay(dayInfo.day)"
            class="flex flex-col rounded-lg border transition-all cursor-pointer overflow-hidden min-h-[340px]"
            :class="[
              selectedDay === dayInfo.day
                ? 'border-[#CEA549] ring-2 ring-[#CEA549]/25 bg-[#FAF9F6]'
                : 'border-[#E2DFD8] hover:border-[#DCD8CF] bg-white'
            ]"
          >
            <!-- Day Column Header -->
            <div
              class="p-2.5 border-b transition-colors flex flex-col gap-1"
              :class="[
                selectedDay === dayInfo.day ? 'bg-[#191A1D] text-white' : 'bg-[#FAF9F6] text-[#191A1D]'
              ]"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider">
                  {{ dayInfo.day }}
                </span>
                <span
                  v-if="dayInfo.isToday"
                  class="px-1.5 py-0.2 bg-[#CEA549] text-[#0B0B0D] text-[10px] font-bold rounded"
                >
                  Hoy
                </span>
              </div>

              <div class="flex items-baseline gap-1">
                <span class="text-lg font-bold leading-none">
                  {{ dayInfo.dayNumber }}
                </span>
                <span class="text-[11px] opacity-75">
                  {{ dayInfo.monthLabel }}
                </span>
              </div>

              <!-- Day availability summary -->
              <div class="pt-1 text-[10px] leading-tight truncate">
                <span
                  v-if="getAvailabilityForDay(dayInfo.day).enabled && (getAvailabilityForDay(dayInfo.day).intervals || []).length > 0"
                  :class="selectedDay === dayInfo.day ? 'text-[#CEA549]' : 'text-[#137333] font-semibold'"
                >
                  {{ (getAvailabilityForDay(dayInfo.day).intervals || []).map(i => `${i.startTime}–${i.endTime}`).join(', ') }}
                </span>
                <span v-else class="text-[#8E9097] italic">
                  Sin disponibilidad
                </span>
              </div>
            </div>

            <!-- Day Column Sessions Body -->
            <div class="p-2 flex-1 flex flex-col justify-between space-y-2">
              <!-- If Day is NOT enabled / No intervals -->
              <div
                v-if="!getAvailabilityForDay(dayInfo.day).enabled || (getAvailabilityForDay(dayInfo.day).intervals || []).length === 0"
                class="flex-1 flex flex-col items-center justify-center p-3 text-center text-[11px] text-[#8E9097] border border-dashed border-[#ECE8E1] rounded bg-[#FAF9F6]/50"
              >
                <span>No disponible</span>
                <button
                  type="button"
                  @click.stop="openConfigureAvailability(dayInfo.day)"
                  class="mt-2 text-[10px] font-bold text-[#8C6B1B] hover:underline"
                >
                  + Habilitar
                </button>
              </div>

              <!-- If Day IS enabled -->
              <div v-else class="space-y-2 flex-1">
                <!-- If No sessions yet -->
                <div
                  v-if="getOccurrencesForDay(dayInfo.day).length === 0"
                  class="p-3 text-center border border-dashed border-[#DCD8CF] rounded-md bg-white text-[11px] text-[#8E9097]"
                >
                  <span>Espacio disponible</span>
                  <div class="mt-2">
                    <button
                      type="button"
                      @click.stop="openNewSessionModal(dayInfo.day, undefined, dayInfo.dateStr)"
                      class="px-2 py-1 bg-white border border-[#DCD8CF] hover:border-[#CEA549] text-[10px] font-bold text-[#191A1D] rounded transition-colors shadow-2xs"
                    >
                      + Sesión
                    </button>
                  </div>
                </div>

                <!-- Session Cards List -->
                <div
                  v-for="occ in getOccurrencesForDay(dayInfo.day)"
                  :key="occ.instanceId"
                  @click.stop="openSessionDetails(occ)"
                  class="p-2.5 rounded-lg border transition-all text-left group cursor-pointer shadow-2xs"
                  :class="[
                    occ.isException
                      ? 'bg-amber-50/70 border-amber-300 hover:border-amber-500'
                      : 'bg-white border-[#DCD8CF] hover:border-[#CEA549]'
                  ]"
                >
                  <!-- Time & Duration -->
                  <div class="flex items-center justify-between">
                    <span class="text-[11px] font-bold text-[#191A1D] flex items-center gap-1">
                      <Clock class="w-3 h-3 text-[#CEA549]" />
                      {{ occ.startTime }} – {{ occ.endTime }}
                    </span>
                    <span class="text-[10px] px-1 py-0.2 rounded font-semibold bg-[#FAF9F6] border border-[#ECE8E1] text-[#6B6D74]">
                      {{ occ.durationMinutes }}m
                    </span>
                  </div>

                  <!-- Students info -->
                  <div class="mt-1.5">
                    <p class="text-xs font-bold text-[#191A1D] truncate leading-tight group-hover:text-black">
                      {{ getStudentNames(occ.studentIds) }}
                    </p>
                    <div class="flex items-center gap-1.5 mt-1 text-[10px] text-[#6B6D74]">
                      <span class="w-1.5 h-1.5 rounded-full" :class="occ.studentIds.length > 1 ? 'bg-[#CEA549]' : 'bg-[#191A1D]'"></span>
                      <span>{{ occ.studentIds.length > 1 ? `Grupal (${occ.studentIds.length})` : 'Individual' }}</span>
                      <span v-if="occ.isException" class="text-amber-800 font-bold ml-auto">Reagendado</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Quick add button at the bottom of enabled day -->
              <button
                v-if="getAvailabilityForDay(dayInfo.day).enabled"
                type="button"
                @click.stop="openNewSessionModal(dayInfo.day, undefined, dayInfo.dateStr)"
                class="w-full py-1.5 border border-dashed border-[#DCD8CF] hover:border-[#CEA549] text-[11px] font-bold text-[#8E9097] hover:text-[#191A1D] rounded transition-colors flex items-center justify-center gap-1 cursor-pointer"
              >
                <Plus class="w-3 h-3" />
                Agregar
              </button>
            </div>
          </div>
        </div>

        <!-- Footnote matching requirements: NO forced breaks -->
        <div class="p-3.5 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] text-xs text-[#6B6D74] flex items-start gap-2">
          <AlertCircle class="w-4 h-4 text-[#8C6B1B] flex-shrink-0 mt-0.5" />
          <p class="leading-relaxed">
            <strong>Agenda flexible del entrenador:</strong> Tú defines tus horas de atención por día y los descansos entre clases según tu criterio. Las sesiones se programan libremente sin intervalos de espera obligatorios.
          </p>
        </div>
      </div>

      <!-- RIGHT: CONTEXTUAL DAY PANEL (4 cols on lg) -->
      <div class="lg:col-span-4 bg-white rounded-[10px] border border-[#E2DFD8] p-5 shadow-2xs space-y-5 sticky top-24">
        <!-- Panel Title -->
        <div class="border-b pb-3 border-[#ECE8E1] flex items-center justify-between">
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2.5 py-0.5 bg-[#191A1D] text-[#CEA549] text-xs font-bold rounded">
                {{ selectedDayInfo.day }}
              </span>
              <span v-if="selectedDayInfo.isToday" class="text-xs font-bold text-[#137333] bg-[#E6F4EA] px-2 py-0.5 rounded">
                Hoy
              </span>
            </div>
            <h2 class="text-base font-bold text-[#191A1D] mt-1">
              {{ selectedDayInfo.dayNumber }} de {{ selectedDayInfo.monthLabel }}, 2026
            </h2>
          </div>

          <button
            @click="openNewSessionModal(selectedDay)"
            class="h-[34px] px-3 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded text-xs font-bold transition-all shadow-xs flex items-center gap-1 cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            Nueva sesión
          </button>
        </div>

        <!-- 1. Configured Availability for this day -->
        <div class="p-3.5 bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
              <Clock class="w-3.5 h-3.5 text-[#CEA549]" />
              Disponibilidad Configurada
            </span>
            <button
              @click="openConfigureAvailability(selectedDay)"
              class="text-xs font-bold text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer"
            >
              Configurar →
            </button>
          </div>

          <div v-if="selectedDayAvailability.enabled && selectedDayAvailability.intervals.length > 0" class="space-y-1.5">
            <div
              v-for="(int, idx) in selectedDayAvailability.intervals"
              :key="int.id"
              class="px-3 py-1.5 bg-white border border-[#DCD8CF] rounded text-xs font-bold text-[#191A1D] flex items-center justify-between shadow-2xs"
            >
              <span>Turno {{ idx + 1 }}:</span>
              <span class="text-[#8C6B1B]">{{ int.startTime }} a {{ int.endTime }}</span>
            </div>
          </div>
          <div v-else class="text-xs text-[#8E9097] py-1">
            Sin disponibilidad activa para este día. El entrenador no atiende los {{ selectedDay }}.
          </div>
        </div>

        <!-- 2. Metrics / Working hours summary -->
        <div class="grid grid-cols-2 gap-3 text-center">
          <div class="p-3 bg-white border border-[#E2DFD8] rounded-lg">
            <span class="text-[11px] font-semibold text-[#8E9097] block">Disponibilidad total</span>
            <span class="text-lg font-bold text-[#191A1D]">{{ dayTotalAvailableHours }} hrs</span>
          </div>
          <div class="p-3 bg-white border border-[#E2DFD8] rounded-lg">
            <span class="text-[11px] font-semibold text-[#8E9097] block">Sesiones asignadas</span>
            <span class="text-lg font-bold text-[#191A1D]">{{ selectedDayOccurrences.length }} ({{ dayScheduledHours }} hrs)</span>
          </div>
        </div>

        <!-- 3. Scheduled Sessions List for this day -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-xs font-bold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
              <Users class="w-3.5 h-3.5 text-[#CEA549]" />
              Sesiones del Día ({{ selectedDayOccurrences.length }})
            </h3>
          </div>

          <div v-if="selectedDayOccurrences.length === 0" class="p-5 text-center border border-dashed border-[#DCD8CF] rounded-lg text-xs text-[#8E9097]">
            No hay sesiones asignadas para este día.
            <div class="mt-2.5">
              <button
                @click="openNewSessionModal(selectedDay)"
                class="px-3.5 py-1.5 bg-[#FAF9F6] border border-[#DCD8CF] hover:border-[#CEA549] text-xs font-bold text-[#191A1D] rounded transition-colors cursor-pointer"
              >
                + Programar para {{ selectedDay }}
              </button>
            </div>
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="occ in selectedDayOccurrences"
              :key="occ.instanceId"
              @click="openSessionDetails(occ)"
              class="p-3 rounded-lg border transition-all cursor-pointer group flex items-center justify-between"
              :class="occ.isException
                ? 'bg-amber-50/70 border-amber-300 hover:border-amber-400'
                : 'bg-[#FAF9F6] border-[#E2DFD8] hover:border-[#CEA549] hover:bg-white'"
            >
              <div>
                <div class="flex items-center gap-2">
                  <span class="text-xs font-bold text-[#191A1D]">
                    {{ occ.startTime }} — {{ occ.endTime }}
                  </span>
                  <span class="text-[10px] px-1.5 py-0.2 rounded font-semibold bg-white border border-[#DCD8CF] text-[#6B6D74]">
                    {{ occ.durationMinutes }}m
                  </span>
                  <span v-if="occ.isException" class="text-[10px] px-1.5 py-0.2 rounded font-bold bg-amber-100 text-amber-900 border border-amber-300">
                    Reagendado
                  </span>
                </div>
                <p class="text-xs font-bold text-[#191A1D] mt-1 group-hover:text-black">
                  {{ getStudentNames(occ.studentIds) }}
                </p>
                <p class="text-[11px] text-[#8E9097]">
                  {{ occ.type }} · {{ occ.studentIds.length }} alumno(s)
                </p>
              </div>

              <ChevronRight class="w-4 h-4 text-[#8E9097] group-hover:text-[#CEA549] transition-colors" />
            </div>
          </div>
        </div>
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

    <ConfigureAvailabilityModal
      :open="showAvailabilityModal"
      :initial-day="selectedDay"
      @close="showAvailabilityModal = false"
      @saved="refreshAgenda"
    />
  </div>
</template>
