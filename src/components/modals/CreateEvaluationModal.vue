<script setup lang="ts">
import { ref, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { EvaluationType, EvaluationStatus, EvaluationMetric } from '@/types'
import { Plus, Trash2 } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  initialStudentId?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created'): void
}>()

const gymStore = useGymStore()

const selectedStudentId = ref(props.initialStudentId || '')
const evaluationType = ref<EvaluationType>('Evaluación inicial')
const status = ref<EvaluationStatus>('Registrada')
const evaluationDate = ref(new Date().toISOString().split('T')[0])
const observations = ref('')

const metrics = ref<EvaluationMetric[]>([
  { name: 'Peso corporal', value: '72.5', unit: 'kg' },
  { name: 'Estatura', value: '174', unit: 'cm' },
  { name: 'Movilidad articular', value: 'Buena movilidad de cadera y tobillo' },
  { name: 'Test fuerza básica', value: 'Sentadilla 60kg x 8 reps' }
])

const isSubmitting = ref(false)
const errorMessage = ref('')

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    if (props.initialStudentId) {
      selectedStudentId.value = props.initialStudentId
    } else if (gymStore.students.length > 0 && !selectedStudentId.value) {
      selectedStudentId.value = gymStore.students[0].id
    }
    errorMessage.value = ''
  }
})

function addMetric() {
  metrics.value.push({ name: '', value: '', unit: '' })
}

function removeMetric(index: number) {
  metrics.value.splice(index, 1)
}

async function handleSubmit() {
  if (!selectedStudentId.value) {
    errorMessage.value = 'Debes seleccionar un alumno.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const validMetrics = metrics.value.filter(m => m.name.trim() !== '')
    await gymStore.createEvaluation({
      studentId: selectedStudentId.value,
      date: evaluationDate.value,
      type: evaluationType.value,
      status: status.value,
      metrics: validMetrics,
      observations: observations.value.trim()
    })

    emit('created')
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al crear evaluación.'
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
          <h2 class="text-lg font-bold text-[#191A1D]">Nueva Evaluación</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Registro de control y progreso físico del alumno.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4 max-h-[75vh] overflow-y-auto">
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-sm">
          {{ errorMessage }}
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Alumno a Evaluar *
          </label>
          <select
            v-model="selectedStudentId"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            required
          >
            <option value="" disabled>Selecciona un alumno...</option>
            <option
              v-for="student in gymStore.students"
              :key="student.id"
              :value="student.id"
            >
              {{ student.fullName }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Tipo de Evaluación
            </label>
            <select
              v-model="evaluationType"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            >
              <option value="Evaluación inicial">Evaluación inicial</option>
              <option value="Control mensual">Control mensual</option>
              <option value="Control trimestral">Control trimestral</option>
              <option value="Re-evaluación">Re-evaluación</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Fecha
            </label>
            <input
              v-model="evaluationDate"
              type="date"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Estado de la Ficha
          </label>
          <select
            v-model="status"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          >
            <option value="Registrada">Registrada</option>
            <option value="Completada">Completada</option>
            <option value="Pendiente revisión">Pendiente revisión</option>
          </select>
        </div>

        <!-- Dynamic Metrics List -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider">
              Métricas y Mediciones (Personalizables)
            </label>
            <button
              type="button"
              @click="addMetric"
              class="text-xs text-[#8C6B1B] hover:text-[#614910] font-bold flex items-center gap-1 cursor-pointer"
            >
              <Plus class="w-3.5 h-3.5" />
              + Añadir otra métrica
            </button>
          </div>

          <div class="space-y-2">
            <div
              v-for="(metric, idx) in metrics"
              :key="idx"
              class="flex items-center gap-2 p-2 bg-[#FAF9F6] border border-[#E2DFD8] rounded-md"
            >
              <input
                v-model="metric.name"
                type="text"
                placeholder="Nombre métrica"
                class="flex-1 px-2.5 py-1.5 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
              />
              <input
                v-model="metric.value"
                type="text"
                placeholder="Valor"
                class="w-32 px-2.5 py-1.5 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
              />
              <input
                v-model="metric.unit"
                type="text"
                placeholder="Unidad"
                class="w-16 px-2.5 py-1.5 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549]"
              />
              <button
                type="button"
                @click="removeMetric(idx)"
                class="p-1.5 text-[#8E9097] hover:text-rose-600 rounded cursor-pointer"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
          <p class="text-[11px] text-[#8E9097] mt-1.5">
            Estructura flexible para registrar métricas sin imponer medidas fijas previo acuerdo con el cliente.
          </p>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Observaciones del Entrenador
          </label>
          <textarea
            v-model="observations"
            rows="3"
            placeholder="Progresión, técnica, sensaciones o recomendaciones para la planificación..."
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
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
            :disabled="isSubmitting"
            class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold transition-colors shadow-xs cursor-pointer disabled:opacity-50"
          >
            {{ isSubmitting ? 'Guardando...' : 'Registrar evaluación' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
