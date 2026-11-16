<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { PaymentMethod } from '@/types'

const props = defineProps<{
  open: boolean
  initialStudentId?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const gymStore = useGymStore()

const selectedStudentId = ref(props.initialStudentId || '')
const amount = ref<number | null>(null)
const paymentMethod = ref<PaymentMethod>('Transferencia')
const paymentDate = ref(new Date().toISOString().split('T')[0])
const notes = ref('')

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

const selectedStudent = computed(() => {
  return gymStore.students.find(s => s.id === selectedStudentId.value)
})

const paymentMethods: PaymentMethod[] = ['Transferencia', 'Webpay', 'Servipag', 'Crédito', 'Débito']

async function handleSubmit() {
  if (!selectedStudentId.value) {
    errorMessage.value = 'Debes seleccionar un alumno.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    await gymStore.registerPayment({
      studentId: selectedStudentId.value,
      amount: (typeof amount.value === 'number' && !isNaN(amount.value) && amount.value >= 0) ? amount.value : null,
      paymentMethod: paymentMethod.value,
      paymentDate: paymentDate.value,
      notes: notes.value.trim()
    })

    emit('saved')
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al registrar pago.'
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
          <h2 class="text-lg font-bold text-[#191A1D]">Registrar Pago</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Control de ingresos y activación de membresía.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-sm">
          {{ errorMessage }}
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Alumno *
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
              {{ student.fullName }} ({{ student.membershipStatus }})
            </option>
          </select>
        </div>

        <div v-if="selectedStudent" class="p-3 bg-[#FAF9F6] border border-[#E2DFD8] rounded-md text-xs space-y-1">
          <div class="flex justify-between">
            <span class="text-[#6B6D74]">Plan actual:</span>
            <span class="font-semibold text-[#191A1D]">{{ selectedStudent.planName }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-[#6B6D74]">Frecuencia:</span>
            <span class="font-semibold text-[#191A1D]">{{ selectedStudent.frequency }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-[#6B6D74]">Estado actual:</span>
            <span
              class="font-semibold"
              :class="selectedStudent.membershipStatus === 'Activo' ? 'text-emerald-700' : 'text-amber-700'"
            >
              {{ selectedStudent.membershipStatus }}
            </span>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Método de Pago
            </label>
            <select
              v-model="paymentMethod"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            >
              <option v-for="m in paymentMethods" :key="m" :value="m">{{ m }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Fecha de Pago
            </label>
            <input
              v-model="paymentDate"
              type="date"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Monto (CLP, opcional)
          </label>
          <input
            v-model.number="amount"
            type="number"
            placeholder="Ej. 45000 (o dejar vacío para Por definir)"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Detalle / Referencia de Transferencia
          </label>
          <input
            v-model="notes"
            type="text"
            placeholder="N° operación banco o comprobante SumUp..."
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <div class="p-3 bg-[#FAF9F6] rounded-md text-xs text-[#6B6D74] border border-[#E2DFD8]">
          <strong class="text-[#191A1D]">Aviso:</strong> Esta acción registra el pago en el historial local y actualiza la membresía a "Activa". No se conecta con pasarelas externas reales.
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
            {{ isSubmitting ? 'Guardando...' : 'Registrar pago' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
