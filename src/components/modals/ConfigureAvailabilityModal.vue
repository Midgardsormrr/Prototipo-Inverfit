<script setup lang="ts">
import { ref, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import { timeToMinutes, minutesToTime, DAY_NAMES } from '@/services/schedule'
import type { DayOfWeek, DayAvailability, ScheduleBlock } from '@/types'
import { Plus, Trash2, Clock, AlertTriangle, Check, Copy, Sparkles } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  initialDay?: DayOfWeek
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const gymStore = useGymStore()

const activeDay = ref<DayOfWeek>('Lun')
const draftAvailability = ref<DayAvailability[]>([])
const conflicts = ref<ScheduleBlock[]>([])
const errorMessage = ref('')
const isSubmitting = ref(false)

function initDraft() {
  conflicts.value = []
  errorMessage.value = ''
  // Deep clone current availability from store
  draftAvailability.value = DAY_NAMES.map(day => {
    const existing = gymStore.availability.find(a => a.day === day)
    if (existing) {
      return {
        day,
        enabled: existing.enabled,
        intervals: (existing.intervals || []).map(i => ({ ...i }))
      }
    }
    return {
      day,
      enabled: false,
      intervals: []
    }
  })

  if (props.initialDay) {
    activeDay.value = props.initialDay
  }
}

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    initDraft()
  }
})

const currentDayAvail = ref<DayAvailability>({
  day: 'Lun',
  enabled: true,
  intervals: []
})

watch([draftAvailability, activeDay], () => {
  const found = draftAvailability.value.find(a => a.day === activeDay.value)
  if (found) {
    currentDayAvail.value = found
  }
}, { immediate: true })

function toggleDayEnabled() {
  currentDayAvail.value.enabled = !currentDayAvail.value.enabled
  if (currentDayAvail.value.enabled && currentDayAvail.value.intervals.length === 0) {
    // Propose default interval
    currentDayAvail.value.intervals.push({
      id: `int-${activeDay.value.toLowerCase()}-${Date.now()}`,
      startTime: '08:00',
      endTime: '13:00'
    })
  }
}

function addInterval() {
  const intervals = currentDayAvail.value.intervals
  let nextStart = '15:00'
  let nextEnd = '19:00'

  if (intervals.length > 0) {
    const lastEnd = intervals[intervals.length - 1].endTime
    const lastEndM = timeToMinutes(lastEnd)
    if (lastEndM + 60 < 22 * 60) {
      nextStart = minutesToTime(lastEndM + 60)
      nextEnd = minutesToTime(Math.min(lastEndM + 240, 22 * 60))
    }
  }

  intervals.push({
    id: `int-${activeDay.value.toLowerCase()}-${Date.now()}-${Math.random().toString(36).substring(2, 5)}`,
    startTime: nextStart,
    endTime: nextEnd
  })
}

function removeInterval(index: number) {
  currentDayAvail.value.intervals.splice(index, 1)
  if (currentDayAvail.value.intervals.length === 0) {
    currentDayAvail.value.enabled = false
  }
}

function copyToWeekdays() {
  const sourceIntervals = (currentDayAvail.value.intervals || []).map(i => ({
    ...i,
    id: `int-copy-${Date.now()}-${Math.random().toString(36).substring(2, 5)}`
  }))
  const isEnabled = currentDayAvail.value.enabled

  const weekdays: DayOfWeek[] = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie']
  weekdays.forEach(day => {
    const target = draftAvailability.value.find(a => a.day === day)
    if (target) {
      target.enabled = isEnabled
      target.intervals = sourceIntervals.map(i => ({
        ...i,
        id: `int-${day.toLowerCase()}-${Math.random().toString(36).substring(2, 6)}`
      }))
    }
  })

  gymStore.showToast(`Horario copiado a todos los días hábiles (Lun a Vie).`, 'info')
}

function applyPabloTemplate() {
  const weekdays: DayOfWeek[] = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie']
  draftAvailability.value.forEach(day => {
    if (weekdays.includes(day.day)) {
      day.enabled = true
      day.intervals = [
        { id: `int-${day.day.toLowerCase()}-1`, startTime: '06:30', endTime: '14:00' },
        { id: `int-${day.day.toLowerCase()}-2`, startTime: '16:00', endTime: '22:00' }
      ]
    } else if (day.day === 'Sáb') {
      day.enabled = true
      day.intervals = [
        { id: `int-sab-eval`, startTime: '06:30', endTime: '13:00' }
      ]
    } else {
      // Domingo descanso
      day.enabled = false
      day.intervals = []
    }
  })
  gymStore.showToast('Plantilla oficial de Pablo aplicada (Lun-Vie 06:30-14:00 y 16:00-22:00, Sáb 06:30-13:00).', 'info', 'Plantilla cargada')
}

async function handleSave() {
  errorMessage.value = ''
  conflicts.value = []

  // Validate intervals
  for (const day of draftAvailability.value) {
    if (!day.enabled) continue
    for (const interval of day.intervals) {
      const s = timeToMinutes(interval.startTime)
      const e = timeToMinutes(interval.endTime)
      if (e <= s) {
        errorMessage.value = `En ${day.day}, la hora de término (${interval.endTime}) debe ser posterior al inicio (${interval.startTime}).`
        activeDay.value = day.day
        return
      }
    }

    // Check overlap between intervals on same day
    for (let i = 0; i < day.intervals.length; i++) {
      for (let j = i + 1; j < day.intervals.length; j++) {
        const i1 = day.intervals[i]
        const i2 = day.intervals[j]
        const s1 = timeToMinutes(i1.startTime)
        const e1 = timeToMinutes(i1.endTime)
        const s2 = timeToMinutes(i2.startTime)
        const e2 = timeToMinutes(i2.endTime)
        if (s1 < e2 && e1 > s2) {
          errorMessage.value = `En ${day.day}, los intervalos de disponibilidad no pueden solaparse (${i1.startTime}-${i1.endTime} con ${i2.startTime}-${i2.endTime}).`
          activeDay.value = day.day
          return
        }
      }
    }
  }

  isSubmitting.value = true
  try {
    const res = await gymStore.saveAvailability(draftAvailability.value)
    if (!res.success) {
      if (res.conflicts && res.conflicts.length > 0) {
        conflicts.value = res.conflicts
      } else {
        errorMessage.value = res.message
      }
      return
    }

    emit('saved')
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al guardar disponibilidad.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <h2 class="text-lg font-bold text-[#191A1D]">Configurar Disponibilidad del Entrenador</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Establece manualmente los días e intervalos horarios en que trabajas.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 space-y-5 max-h-[75vh] overflow-y-auto">
        <!-- Error alert -->
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-xs">
          {{ errorMessage }}
        </div>

        <!-- Conflict alert when reducing availability leaves sessions outside -->
        <div v-if="conflicts.length > 0" class="p-4 bg-amber-50 border border-amber-300 rounded-lg space-y-2">
          <div class="flex items-center gap-2 text-amber-900 font-bold text-sm">
            <AlertTriangle class="w-4 h-4 text-amber-700" />
            Conflicto de disponibilidad con sesiones existentes
          </div>
          <p class="text-xs text-amber-800 leading-relaxed">
            Las siguientes sesiones ya están programadas y quedarían fuera del nuevo horario configurado. Debes reubicar o cancelar estas sesiones antes de reducir tu disponibilidad:
          </p>
          <div class="bg-white/80 p-2.5 rounded border border-amber-200 divide-y divide-amber-100 text-xs text-amber-950">
            <div v-for="c in conflicts" :key="c.id" class="py-1.5 flex items-center justify-between">
              <span class="font-bold">{{ c.day }} · {{ c.startTime }} a {{ c.endTime }}</span>
              <span class="text-amber-800">{{ c.studentIds.length }} alumno(s)</span>
            </div>
          </div>
        </div>

        <!-- Pablo's Standard Template Banner & Quick Action -->
        <div class="p-3.5 bg-[#FAF6EC] border border-[#E8DEC7] rounded-xl flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            <span class="text-xs font-bold text-[#191A1D] flex items-center gap-1.5">
              <Sparkles class="w-3.5 h-3.5 text-[#CEA549]" />
              Plantilla Semanal Oficial INVERFIT (Pablo)
            </span>
            <p class="text-[11px] text-[#6B6D74] mt-0.5">
              Lun a Vie (06:30–14:00 y 16:00–22:00) · Sábados: Evaluaciones (06:30–13:00) · Dom: Cerrado
            </p>
          </div>
          <button
            type="button"
            @click="applyPabloTemplate"
            class="px-3 py-1.5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] text-xs font-bold rounded-md transition-colors cursor-pointer shadow-xs self-start sm:self-auto flex items-center gap-1.5"
          >
            <Sparkles class="w-3.5 h-3.5" />
            Aplicar plantilla
          </button>
        </div>

        <!-- Day Selector Tabs -->
        <div>
          <label class="block text-xs font-bold text-[#6B6D74] uppercase tracking-wider mb-2">
            Selecciona el Día de la Semana
          </label>
          <div class="grid grid-cols-7 gap-1.5 p-1 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg">
            <button
              v-for="day in DAY_NAMES"
              :key="day"
              type="button"
              @click="activeDay = day"
              class="py-2.5 rounded-md text-xs font-bold transition-all flex flex-col items-center gap-1 cursor-pointer"
              :class="[
                activeDay === day
                  ? 'bg-[#191A1D] text-[#CEA549] shadow-xs'
                  : 'text-[#6B6D74] hover:bg-white hover:text-[#191A1D]'
              ]"
            >
              <span>{{ day }}</span>
              <!-- Availability dot indicator -->
              <span
                class="w-1.5 h-1.5 rounded-full"
                :class="draftAvailability.find(a => a.day === day)?.enabled ? 'bg-[#137333]' : 'bg-[#DCD8CF]'"
              ></span>
            </button>
          </div>
        </div>

        <!-- Active Day Configuration Box -->
        <div class="p-5 bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl space-y-4">
          <div class="flex items-center justify-between border-b pb-3 border-[#ECE8E1]">
            <div class="flex items-center gap-3">
              <span class="px-2.5 py-1 bg-[#191A1D] text-[#CEA549] text-xs font-bold rounded">
                {{ activeDay }}
              </span>
              <div>
                <h3 class="text-sm font-bold text-[#191A1D]">
                  {{ currentDayAvail.enabled ? 'Día con disponibilidad activa' : 'Sin disponibilidad para este día' }}
                </h3>
                <p class="text-[11px] text-[#8E9097]">
                  {{ currentDayAvail.enabled ? 'El entrenador puede recibir sesiones en los intervalos indicados.' : 'No se podrán programar sesiones en este día.' }}
                </p>
              </div>
            </div>

            <!-- Toggle Day Button -->
            <button
              type="button"
              @click="toggleDayEnabled"
              class="px-3 py-1.5 text-xs font-bold rounded transition-colors cursor-pointer"
              :class="currentDayAvail.enabled
                ? 'bg-[#E6F4EA] text-[#137333] hover:bg-rose-100 hover:text-rose-800'
                : 'bg-[#191A1D] text-[#CEA549] hover:bg-[#24262B]'"
            >
              {{ currentDayAvail.enabled ? '✓ Disponible (clic para desactivar)' : '+ Habilitar día' }}
            </button>
          </div>

          <!-- Interval List if Enabled -->
          <div v-if="currentDayAvail.enabled" class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
                <Clock class="w-3.5 h-3.5 text-[#CEA549]" />
                Intervalos Horarios de Atención
              </label>

              <button
                type="button"
                @click="addInterval"
                class="text-xs font-bold text-[#8C6B1B] hover:text-[#614910] flex items-center gap-1 cursor-pointer"
              >
                <Plus class="w-3.5 h-3.5" />
                + Agregar otro intervalo
              </button>
            </div>

            <div class="space-y-2">
              <div
                v-for="(interval, idx) in currentDayAvail.intervals"
                :key="interval.id"
                class="p-3 bg-white border border-[#DCD8CF] rounded-lg flex items-center justify-between gap-3 shadow-2xs"
              >
                <div class="flex items-center gap-2 flex-1">
                  <span class="text-xs font-bold text-[#6B6D74] w-20">Turno {{ idx + 1 }}:</span>

                  <div class="flex items-center gap-2">
                    <label class="text-xs text-[#8E9097]">De:</label>
                    <input
                      v-model="interval.startTime"
                      type="time"
                      class="px-2.5 py-1.5 border border-[#DCD8CF] rounded text-xs font-bold text-[#191A1D] bg-[#FAF9F6] focus:outline-none focus:border-[#CEA549] focus:bg-white"
                      required
                    />
                  </div>

                  <div class="flex items-center gap-2">
                    <label class="text-xs text-[#8E9097]">Hasta:</label>
                    <input
                      v-model="interval.endTime"
                      type="time"
                      class="px-2.5 py-1.5 border border-[#DCD8CF] rounded text-xs font-bold text-[#191A1D] bg-[#FAF9F6] focus:outline-none focus:border-[#CEA549] focus:bg-white"
                      required
                    />
                  </div>
                </div>

                <button
                  type="button"
                  @click="removeInterval(idx)"
                  class="p-1.5 text-[#8E9097] hover:text-rose-600 rounded transition-colors cursor-pointer"
                  title="Eliminar este intervalo"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>

            <!-- Quick copy action -->
            <div class="pt-2 flex items-center justify-between text-xs text-[#6B6D74]">
              <span class="text-[11px] text-[#8E9097]">¿Mismo horario para la semana?</span>
              <button
                type="button"
                @click="copyToWeekdays"
                class="text-xs font-bold text-[#191A1D] hover:text-[#CEA549] flex items-center gap-1.5 transition-colors cursor-pointer underline"
              >
                <Copy class="w-3.5 h-3.5" />
                Copiar a todos los días hábiles (Lun a Vie)
              </button>
            </div>
          </div>

          <div v-else class="py-6 text-center text-xs text-[#8E9097]">
            Este día está marcado como no disponible (cerrado). Ninguna sesión podrá programarse para este día.
          </div>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="px-6 py-4 border-t border-[#E2DFD8] flex items-center justify-end gap-3 bg-[#FAF9F6]">
        <button
          type="button"
          @click="emit('close')"
          class="px-4 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#6B6D74] hover:bg-[#F5F3EE] hover:text-[#191A1D] font-medium transition-colors cursor-pointer"
        >
          Cancelar
        </button>
        <button
          type="button"
          @click="handleSave"
          :disabled="isSubmitting"
          class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold transition-all shadow-xs cursor-pointer disabled:opacity-50 flex items-center gap-2"
        >
          <Check class="w-4 h-4" />
          {{ isSubmitting ? 'Guardando...' : 'Guardar Disponibilidad' }}
        </button>
      </div>
    </div>
  </div>
</template>
