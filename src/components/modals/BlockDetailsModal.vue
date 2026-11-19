<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { DAY_NAMES, timeToMinutes, addMinutesToTime } from '@/services/schedule'
import type { ScheduleBlock, SessionOccurrence, DayOfWeek } from '@/types'
import { Users, Clock, Trash2, UserPlus, ArrowRight, Calendar, RotateCcw } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  block: ScheduleBlock | null
  occurrence?: SessionOccurrence | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'add-student', block: ScheduleBlock): void
  (e: 'edit-block', block: ScheduleBlock): void
  (e: 'deleted', blockId: string): void
  (e: 'cancelled', occurrence: SessionOccurrence): void
  (e: 'restored'): void
  (e: 'updated'): void
}>()

const router = useRouter()
const gymStore = useGymStore()

// Live reactive block resolved from store so removals and additions reflect immediately
const currentBlock = computed<ScheduleBlock | null>(() => {
  if (!props.block) return null
  return gymStore.schedule.find(b => b.id === props.block!.id) || props.block
})

const assignedStudents = computed(() => {
  if (!currentBlock.value) return []
  return gymStore.students.filter(s => (currentBlock.value!.studentIds || []).includes(s.id))
})

// Exception / reschedule state
const showReschedulePanel = ref(false)
const rescheduleDay = ref<DayOfWeek>('Lun')
const rescheduleDate = ref('')
const rescheduleStartTime = ref('')
const rescheduleEndTime = ref('')
const isSubmittingException = ref(false)
const exceptionError = ref('')

function openRescheduleForm() {
  if (!props.occurrence) return
  rescheduleDate.value = props.occurrence.date
  rescheduleDay.value = props.occurrence.day
  rescheduleStartTime.value = props.occurrence.startTime
  rescheduleEndTime.value = props.occurrence.endTime
  showReschedulePanel.value = true
  exceptionError.value = ''
}

function handleRescheduleStartTimeChange() {
  if (rescheduleStartTime.value) {
    rescheduleEndTime.value = addMinutesToTime(rescheduleStartTime.value, 75)
  }
}

async function handleSaveException() {
  if (!currentBlock.value || !props.occurrence) return

  const s = timeToMinutes(rescheduleStartTime.value)
  const e = timeToMinutes(rescheduleEndTime.value)
  if (e <= s) {
    exceptionError.value = 'La hora de término debe ser posterior a la de inicio.'
    return
  }

  isSubmittingException.value = true
  exceptionError.value = ''

  try {
    await gymStore.createException({
      blockId: currentBlock.value.id,
      date: props.occurrence.date,
      type: 'rescheduled',
      newDate: rescheduleDate.value,
      newDay: rescheduleDay.value,
      newStartTime: rescheduleStartTime.value,
      newEndTime: rescheduleEndTime.value,
      notes: `Excepción semanal: movido a ${rescheduleDay.value} ${rescheduleStartTime.value}.`
    })

    showReschedulePanel.value = false
    emit('updated')
    emit('close')
  } catch (err: unknown) {
    exceptionError.value = err instanceof Error ? err.message : 'Error al reagendar excepción.'
  } finally {
    isSubmittingException.value = false
  }
}

async function handleCancelThisOccurrence() {
  if (!currentBlock.value || !props.occurrence) return
  if (confirm(`¿Cancelar la sesión sólo para esta semana (${props.occurrence.date})? El horario recurrente se mantendrá intacto para las demás semanas.`)) {
    try {
      const occ = props.occurrence
      await gymStore.createException({
        blockId: currentBlock.value.id,
        date: props.occurrence.date,
        type: 'cancelled',
        notes: 'Cancelado para esta fecha específica.'
      })
      emit('cancelled', occ)
      emit('updated')
      emit('close')
    } catch (err: unknown) {
      gymStore.showToast(err instanceof Error ? err.message : 'Error al cancelar ocurrencia', 'error')
    }
  }
}

async function handleRestoreRecurring() {
  if (!props.occurrence?.exceptionId) return
  if (confirm('¿Restablecer esta sesión a su horario recurrente habitual?')) {
    try {
      await gymStore.deleteException(props.occurrence.exceptionId, props.occurrence.blockId)
      emit('restored')
      emit('updated')
      emit('close')
    } catch (err: unknown) {
      gymStore.showToast(err instanceof Error ? err.message : 'Error al restaurar sesión habitual', 'error')
    }
  }
}

async function handleRemoveStudent(studentId: string, studentName: string) {
  if (!currentBlock.value) return
  if (confirm(`¿Quitar a ${studentName} de esta sesión?`)) {
    try {
      await gymStore.removeStudentFromSchedule(currentBlock.value.id, studentId)
      emit('updated')
    } catch {
      // Toast gestionado en el store
    }
  }
}

async function handleDeleteBlock() {
  if (!currentBlock.value) return
  if (confirm(`¿Eliminar definitivamente esta sesión recurrente de ${currentBlock.value.day} ${currentBlock.value.startTime}?`)) {
    const id = currentBlock.value.id
    try {
      await gymStore.deleteBlock(id)
      emit('deleted', id)
      emit('close')
    } catch {
      // Si falla, el store muestra toast y el bloque NO se elimina falsamente de la vista
    }
  }
}

function goToStudent(studentId: string) {
  emit('close')
  router.push(`/alumnos/${studentId}`)
}
</script>

<template>
  <div
    v-if="open && currentBlock"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2 py-0.5 bg-[#191A1D] text-[#CEA549] text-xs font-bold rounded">
              {{ occurrence?.day || currentBlock.day }}
            </span>
            <span v-if="occurrence?.isException" class="px-2 py-0.5 bg-amber-100 text-amber-900 border border-amber-300 text-[11px] font-bold rounded">
              Excepción de esta semana
            </span>
            <span v-if="currentBlock.type === 'Evaluación' || occurrence?.type === 'Evaluación'" class="px-2 py-0.5 bg-[#CEA549]/15 text-[#8C6B1B] border border-[#CEA549]/30 text-[11px] font-bold rounded">
              Evaluación
            </span>
            <h2 class="text-lg font-bold text-[#191A1D]">
              {{ (currentBlock.type === 'Evaluación' || occurrence?.type === 'Evaluación') ? 'Bloque de Evaluación' : 'Sesión de Entrenamiento' }}
            </h2>
          </div>
          <p class="text-xs text-[#6B6D74] mt-1 flex items-center gap-2">
            <Clock class="w-3.5 h-3.5 text-[#CEA549]" />
            <strong class="text-[#191A1D]">{{ occurrence?.startTime || currentBlock.startTime }} — {{ occurrence?.endTime || currentBlock.endTime }}</strong>
            <span>({{ occurrence?.durationMinutes || currentBlock.durationMinutes }} min)</span>
            <span>· {{ currentBlock.type }}</span>
            <span v-if="occurrence?.date" class="text-[#8E9097]">· {{ occurrence.date }}</span>
          </p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-4 max-h-[75vh] overflow-y-auto">
        <!-- Assigned Students List -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <h3 class="text-xs font-bold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
                <Users class="w-4 h-4 text-[#CEA549]" />
                Alumnos Asignados
              </h3>
              <span
                class="px-2 py-0.5 rounded-full text-[11px] font-bold"
                :class="assignedStudents.length >= 3
                  ? 'bg-amber-100 text-amber-900 border border-amber-300'
                  : 'bg-emerald-50 text-emerald-800 border border-emerald-200'"
              >
                {{ assignedStudents.length }} / 3 cupos
                <span v-if="assignedStudents.length >= 3">· Completo</span>
              </span>
            </div>

            <button
              v-if="assignedStudents.length < 3"
              @click="emit('add-student', currentBlock)"
              class="text-xs text-[#8C6B1B] hover:text-[#614910] font-bold flex items-center gap-1 hover:underline cursor-pointer"
            >
              <UserPlus class="w-3.5 h-3.5" />
              + Asignar alumno
            </button>
            <span v-else class="text-[11px] text-[#8E9097] font-medium">
              Cupo máximo alcanzado
            </span>
          </div>

          <div v-if="assignedStudents.length === 0" class="p-6 text-center border border-dashed border-[#DCD8CF] rounded-lg text-sm text-[#6B6D74] bg-[#FAF9F6]">
            No hay alumnos asignados a esta sesión actualmente.
            <div class="mt-2.5">
              <button
                @click="emit('add-student', currentBlock)"
                class="px-3.5 py-1.5 bg-[#CEA549] text-[#0B0B0D] text-xs font-bold rounded hover:bg-[#BA923B] transition-colors shadow-xs cursor-pointer"
              >
                + Asignar alumno
              </button>
            </div>
          </div>

          <div v-else class="divide-y divide-[#ECE8E1] border border-[#E2DFD8] rounded-lg overflow-hidden bg-white">
            <div
              v-for="student in assignedStudents"
              :key="student.id"
              class="p-3.5 flex items-center justify-between hover:bg-[#FAF9F6] transition-colors"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/30 flex items-center justify-center text-xs font-bold">
                  {{ student.fullName.charAt(0) }}
                </div>
                <div>
                  <p class="text-sm font-bold text-[#191A1D]">{{ student.fullName }}</p>
                  <p class="text-xs text-[#6B6D74]">
                    {{ student.planName }} · {{ student.frequency }}
                    <span
                      class="ml-1.5 px-1.5 py-0.2 rounded text-[10px] font-semibold"
                      :class="student.membershipStatus === 'Activo' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                    >
                      {{ student.membershipStatus }}
                    </span>
                  </p>
                </div>
              </div>

              <div class="flex items-center gap-2">
                <button
                  @click="goToStudent(student.id)"
                  class="px-2.5 py-1 text-xs border border-[#DCD8CF] rounded text-[#191A1D] hover:bg-[#F5F3EE] font-medium flex items-center gap-1 transition-colors cursor-pointer"
                >
                  Ficha <ArrowRight class="w-3 h-3 text-[#CEA549]" />
                </button>
                <button
                  @click="handleRemoveStudent(student.id, student.fullName)"
                  title="Quitar de este bloque"
                  class="p-1.5 text-[#8E9097] hover:text-rose-600 hover:bg-rose-50 rounded transition-colors cursor-pointer"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Exception Panel (Section 2.E: Cambios Excepcionales) -->
        <div v-if="occurrence && !showReschedulePanel" class="p-3.5 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg space-y-2">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-xs font-bold text-[#191A1D]">Cambio Excepcional para esta Semana</span>
              <p class="text-[11px] text-[#8E9097]">
                Mueve o cancela la sesión de la fecha {{ occurrence.date }} sin modificar las demás semanas.
              </p>
            </div>
            <div class="flex items-center gap-2">
              <button
                v-if="occurrence.isException"
                type="button"
                @click="handleRestoreRecurring"
                class="px-2.5 py-1 text-xs font-bold text-[#8C6B1B] hover:bg-[#ECE8E1] rounded transition-colors cursor-pointer flex items-center gap-1"
              >
                <RotateCcw class="w-3 h-3" />
                Restablecer horario habitual
              </button>
              <button
                v-else
                type="button"
                @click="openRescheduleForm"
                class="px-3 py-1.5 bg-white border border-[#DCD8CF] hover:border-[#CEA549] text-xs font-bold text-[#191A1D] rounded transition-colors cursor-pointer shadow-2xs"
              >
                Reagendar sólo esta semana
              </button>
              <button
                type="button"
                @click="handleCancelThisOccurrence"
                class="px-2.5 py-1.5 text-xs text-rose-700 hover:bg-rose-50 rounded transition-colors cursor-pointer"
                title="Cancelar sólo esta ocurrencia"
              >
                Cancelar ocurrencia
              </button>
            </div>
          </div>
        </div>

        <!-- Reschedule Form active -->
        <div v-if="showReschedulePanel" class="p-4 bg-amber-50/60 border border-amber-200 rounded-lg space-y-3 animate-in fade-in">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-amber-950 flex items-center gap-1.5">
              <Calendar class="w-3.5 h-3.5 text-amber-700" />
              Reagendar sesión para la semana del {{ occurrence?.date }}
            </span>
            <button @click="showReschedulePanel = false" class="text-xs text-[#8E9097] hover:text-[#191A1D]">✕ Cancelar</button>
          </div>

          <div v-if="exceptionError" class="p-2 bg-rose-50 border border-rose-200 text-rose-800 rounded text-xs">
            {{ exceptionError }}
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-bold text-[#6B6D74] uppercase mb-1">Nuevo Día</label>
              <select
                v-model="rescheduleDay"
                class="w-full px-2.5 py-1.5 border border-[#DCD8CF] bg-white rounded text-xs font-bold text-[#191A1D]"
              >
                <option v-for="d in DAY_NAMES" :key="d" :value="d">{{ d }}</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-bold text-[#6B6D74] uppercase mb-1">Fecha Específica</label>
              <input
                v-model="rescheduleDate"
                type="date"
                class="w-full px-2.5 py-1.5 border border-[#DCD8CF] bg-white rounded text-xs font-bold text-[#191A1D]"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-[#6B6D74] uppercase mb-1">Nuevo Inicio</label>
              <input
                v-model="rescheduleStartTime"
                type="time"
                @change="handleRescheduleStartTimeChange"
                class="w-full px-2.5 py-1.5 border border-[#DCD8CF] bg-white rounded text-xs font-bold text-[#191A1D]"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-[#6B6D74] uppercase mb-1">Nuevo Término</label>
              <input
                v-model="rescheduleEndTime"
                type="time"
                class="w-full px-2.5 py-1.5 border border-[#DCD8CF] bg-white rounded text-xs font-bold text-[#191A1D]"
              />
            </div>
          </div>

          <div class="pt-2 flex justify-end gap-2">
            <button
              type="button"
              @click="showReschedulePanel = false"
              class="px-3 py-1 bg-white border border-[#DCD8CF] rounded text-xs font-medium text-[#6B6D74]"
            >
              Volver
            </button>
            <button
              type="button"
              @click="handleSaveException"
              :disabled="isSubmittingException"
              class="px-4 py-1.5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded text-xs font-bold shadow-xs cursor-pointer"
            >
              {{ isSubmittingException ? 'Reagendando...' : 'Aplicar cambio excepcional' }}
            </button>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="pt-3 border-t border-[#E2DFD8] flex items-center justify-between">
          <button
            @click="handleDeleteBlock"
            class="text-xs text-rose-600 hover:text-rose-800 font-medium flex items-center gap-1 cursor-pointer"
          >
            <Trash2 class="w-3.5 h-3.5" />
            Eliminar sesión recurrente
          </button>

          <button
            @click="emit('close')"
            class="px-4 py-2 bg-[#191A1D] text-white rounded-md text-sm font-semibold hover:bg-[#24262B] transition-colors cursor-pointer"
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
