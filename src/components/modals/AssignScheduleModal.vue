<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import { scheduleService, addMinutesToTime, timeToMinutes, DAY_NAMES } from '@/services/schedule'
import type { DayOfWeek } from '@/types'
import { Clock, Users, Calendar, AlertTriangle, CheckCircle2 } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  initialDay?: DayOfWeek
  initialTime?: string
  initialEndTime?: string
  initialStudentId?: string
  blockId?: string // If adding/editing an existing block
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'assigned'): void
}>()

const gymStore = useGymStore()

// State
const selectedStudentIds = ref<string[]>([])
const selectedDay = ref<DayOfWeek>(props.initialDay || 'Lun')
const startTime = ref(props.initialTime || '08:00')
const endTime = ref(props.initialEndTime || '09:15')
const isRecurring = ref(true)
const notes = ref('')

const isSubmitting = ref(false)
const errorMessage = ref('')

// Initialize / sync when modal opens
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    if (props.blockId) {
      const existingBlock = gymStore.schedule.find(b => b.id === props.blockId)
      if (existingBlock) {
        selectedStudentIds.value = [...existingBlock.studentIds]
        selectedDay.value = existingBlock.day
        startTime.value = existingBlock.startTime
        endTime.value = existingBlock.endTime
        isRecurring.value = existingBlock.isRecurring !== false
        notes.value = existingBlock.notes || ''
      }
    } else {
      if (props.initialStudentId) {
        selectedStudentIds.value = [props.initialStudentId]
      } else if (selectedStudentIds.value.length === 0 && gymStore.students.length > 0) {
        selectedStudentIds.value = [gymStore.students[0].id]
      }

      if (props.initialDay) {
        selectedDay.value = props.initialDay
      }

      if (props.initialTime) {
        startTime.value = props.initialTime
      }

      if (props.initialEndTime) {
        endTime.value = props.initialEndTime
      } else {
        endTime.value = addMinutesToTime(startTime.value, 75)
      }
    }

    errorMessage.value = ''
  }
})

// Auto-suggest endTime when startTime changes
function handleStartTimeChange() {
  if (startTime.value) {
    endTime.value = addMinutesToTime(startTime.value, 75)
  }
}

// Calculate duration in minutes live
const durationMinutes = computed(() => {
  if (!startTime.value || !endTime.value) return 0
  const s = timeToMinutes(startTime.value)
  const e = timeToMinutes(endTime.value)
  return Math.max(0, e - s)
})

// Live validation check
const validationCheck = computed(() => {
  if (!selectedDay.value || !startTime.value || !endTime.value) {
    return { valid: true }
  }
  return scheduleService.validateSession(selectedDay.value, startTime.value, endTime.value, props.blockId)
})

// Availability info for the selected day
const dayAvailabilityInfo = computed(() => {
  const avail = gymStore.availability.find(a => a.day === selectedDay.value)
  if (!avail || !avail.enabled || avail.intervals.length === 0) {
    return 'Sin disponibilidad configurada para este día'
  }
  return avail.intervals.map(i => `${i.startTime} a ${i.endTime}`).join(' y ')
})

function toggleStudent(studentId: string) {
  const index = selectedStudentIds.value.indexOf(studentId)
  if (index >= 0) {
    if (selectedStudentIds.value.length > 1) {
      selectedStudentIds.value.splice(index, 1)
    }
  } else {
    selectedStudentIds.value.push(studentId)
  }
}

async function handleAssign() {
  if (selectedStudentIds.value.length === 0) {
    errorMessage.value = 'Debes seleccionar al menos un alumno participante.'
    return
  }

  if (!validationCheck.value.valid) {
    errorMessage.value = validationCheck.value.reason || 'Conflicto de horario.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const res = await gymStore.assignSession({
      blockId: props.blockId,
      studentIds: selectedStudentIds.value,
      day: selectedDay.value,
      startTime: startTime.value,
      endTime: endTime.value,
      isRecurring: isRecurring.value,
      notes: notes.value.trim()
    })

    if (!res.success) {
      errorMessage.value = res.message
      return
    }

    emit('assigned')
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al programar sesión.'
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
    <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <h2 class="text-lg font-bold text-[#191A1D]">
            {{ blockId ? 'Editar Sesión de Entrenamiento' : 'Programar Sesión de Entrenamiento' }}
          </h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Asigna alumnos y define libremente el horario de inicio y término.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleAssign" class="p-6 space-y-4 max-h-[75vh] overflow-y-auto">
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-xs">
          {{ errorMessage }}
        </div>

        <!-- Student selector -->
        <div>
          <div class="flex items-center justify-between mb-1.5">
            <label class="block text-xs font-bold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
              <Users class="w-3.5 h-3.5 text-[#CEA549]" />
              Alumno(s) Participante(s) *
            </label>
            <span class="text-[11px] text-[#8E9097]">
              {{ selectedStudentIds.length }} seleccionado(s)
            </span>
          </div>

          <!-- Multi-select pill container -->
          <div class="p-2 border border-[#DCD8CF] rounded-lg bg-[#FAF9F6] max-h-36 overflow-y-auto space-y-1.5">
            <div
              v-for="student in gymStore.students"
              :key="student.id"
              @click="toggleStudent(student.id)"
              class="px-2.5 py-1.5 rounded-md text-xs font-medium flex items-center justify-between cursor-pointer transition-colors"
              :class="selectedStudentIds.includes(student.id)
                ? 'bg-[#191A1D] text-[#CEA549] font-bold shadow-2xs'
                : 'bg-white hover:bg-neutral-100 text-[#191A1D] border border-[#ECE8E1]'"
            >
              <span>{{ student.fullName }}</span>
              <span class="text-[10px] opacity-75">
                {{ student.frequency }} · {{ student.membershipStatus }}
              </span>
            </div>
          </div>
          <p class="text-[11px] text-[#8E9097] mt-1">
            Puedes seleccionar varios alumnos para formar una sesión grupal en el mismo bloque.
          </p>
        </div>

        <!-- Day selector -->
        <div>
          <label class="block text-xs font-bold text-[#6B6D74] uppercase tracking-wider mb-1 flex items-center gap-1.5">
            <Calendar class="w-3.5 h-3.5 text-[#CEA549]" />
            Día de la Semana *
          </label>
          <div class="grid grid-cols-7 gap-1">
            <button
              v-for="d in DAY_NAMES"
              :key="d"
              type="button"
              @click="selectedDay = d"
              class="py-2 text-xs font-bold rounded-md border transition-all cursor-pointer text-center"
              :class="selectedDay === d
                ? 'bg-[#191A1D] border-[#191A1D] text-[#CEA549] shadow-xs'
                : 'bg-white border-[#DCD8CF] text-[#6B6D74] hover:bg-[#FAF9F6]'"
            >
              {{ d }}
            </button>
          </div>

          <!-- Day Availability indicator -->
          <div class="mt-2 text-[11px] text-[#6B6D74] flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full" :class="dayAvailabilityInfo.includes('Sin') ? 'bg-amber-500' : 'bg-emerald-500'"></span>
            <span>Disponibilidad {{ selectedDay }}: <strong>{{ dayAvailabilityInfo }}</strong></span>
          </div>
        </div>

        <!-- Start & End Time Inputs (Customizable & Flexible) -->
        <div class="grid grid-cols-2 gap-4 p-3.5 bg-[#FAF9F6] border border-[#E2DFD8] rounded-xl">
          <div>
            <label class="block text-xs font-bold text-[#6B6D74] uppercase tracking-wider mb-1 flex items-center gap-1">
              <Clock class="w-3 h-3 text-[#CEA549]" />
              Hora de Inicio *
            </label>
            <input
              v-model="startTime"
              type="time"
              @change="handleStartTimeChange"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm font-bold text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#6B6D74] uppercase tracking-wider mb-1 flex items-center gap-1">
              <Clock class="w-3 h-3 text-[#CEA549]" />
              Hora de Término *
            </label>
            <input
              v-model="endTime"
              type="time"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm font-bold text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>

          <!-- Duration live calculation -->
          <div class="col-span-2 pt-1 border-t border-[#ECE8E1] flex items-center justify-between text-xs">
            <span class="text-[#6B6D74]">Duración calculada:</span>
            <span class="font-bold text-[#191A1D] px-2 py-0.5 rounded bg-white border border-[#DCD8CF]">
              {{ durationMinutes }} minutos
              <span v-if="durationMinutes === 75" class="text-[#8C6B1B] font-semibold text-[11px] ml-1">(Predeterminado)</span>
            </span>
          </div>
        </div>

        <!-- Recurrence mode -->
        <div class="p-3 bg-white border border-[#DCD8CF] rounded-lg flex items-center justify-between">
          <div>
            <span class="text-xs font-bold text-[#191A1D] block">Modalidad recurrente</span>
            <p class="text-[11px] text-[#8E9097]">Se programará automáticamente todas las semanas en este día y horario.</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input
              v-model="isRecurring"
              type="checkbox"
              class="sr-only peer"
            />
            <div class="w-10 h-5 bg-[#DCD8CF] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#CEA549]"></div>
          </label>
        </div>

        <!-- Validation status banner -->
        <div
          v-if="!validationCheck.valid"
          class="p-3 bg-amber-50 border border-amber-300 rounded-lg text-xs text-amber-900 flex items-start gap-2"
        >
          <AlertTriangle class="w-4 h-4 text-amber-700 flex-shrink-0 mt-0.5" />
          <span>{{ validationCheck.reason }}</span>
        </div>

        <div
          v-else-if="startTime && endTime"
          class="p-2.5 bg-emerald-50 border border-emerald-200 rounded-lg text-xs text-emerald-800 flex items-center gap-2"
        >
          <CheckCircle2 class="w-4 h-4 text-emerald-600 flex-shrink-0" />
          <span>Horario válido y dentro de la disponibilidad del entrenador.</span>
        </div>

        <!-- Notes (optional) -->
        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Notas u Observaciones (Opcional)
          </label>
          <textarea
            v-model="notes"
            rows="2"
            placeholder="Enfoque de la sesión, adaptaciones específicas o recordatorios..."
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          ></textarea>
        </div>

        <!-- Footer -->
        <div class="pt-4 border-t border-[#E2DFD8] flex items-center justify-end gap-3">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#6B6D74] hover:bg-[#F5F3EE] hover:text-[#191A1D] font-medium transition-colors cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !validationCheck.valid"
            class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold transition-all shadow-xs cursor-pointer disabled:opacity-50"
          >
            {{ isSubmitting ? 'Guardando...' : (blockId ? 'Actualizar sesión' : '+ Programar sesión') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
