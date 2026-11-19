<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Upload, FileText, Check, Trash2 } from 'lucide-vue-next'
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
const amount = ref<number | null>(45000)
const paymentMethod = ref<PaymentMethod>('SumUp (Link de pago)')
const paymentDate = ref(new Date().toISOString().split('T')[0])
const transactionReference = ref('')
const notes = ref('')

// Voucher attachment
const fileInputRef = ref<HTMLInputElement | null>(null)
const voucherUrl = ref<string>('')
const voucherFileName = ref<string>('')

const isSubmitting = ref(false)
const errorMessage = ref('')

// Membresía única: pago mensual
const defaultMonthlyFee = 45000

const selectedStudent = computed(() => {
  return gymStore.students.find(s => s.id === selectedStudentId.value)
})

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    if (props.initialStudentId) {
      selectedStudentId.value = props.initialStudentId
    } else if (gymStore.students.length > 0 && !selectedStudentId.value) {
      selectedStudentId.value = gymStore.students[0].id
    }
    voucherUrl.value = ''
    voucherFileName.value = ''
    transactionReference.value = ''
    notes.value = ''
    errorMessage.value = ''
    amount.value = defaultMonthlyFee
  }
})

watch(selectedStudent, (student) => {
  if (student && !amount.value) {
    amount.value = defaultMonthlyFee
  }
})

const paymentMethods: PaymentMethod[] = [
  'SumUp (Link de pago)',
  'SumUp (POS Tarjeta)',
  'Transferencia Bancaria',
  'Efectivo'
]

function handleFileUpload(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    voucherFileName.value = file.name

    const reader = new FileReader()
    reader.onload = (uploadEvent) => {
      voucherUrl.value = (uploadEvent.target?.result as string) || ''
    }
    reader.readAsDataURL(file)
  }
}

function removeVoucher() {
  voucherUrl.value = ''
  voucherFileName.value = ''
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

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
      transactionReference: transactionReference.value.trim() || undefined,
      voucherUrl: voucherUrl.value || undefined,
      voucherFileName: voucherFileName.value || undefined,
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
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] overflow-hidden my-6 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-[#CEA549] bg-[#191A1D] px-2 py-0.5 rounded">
            Recaudación INVERFIT
          </span>
          <h2 class="text-lg font-bold text-[#191A1D] mt-1">Registrar Pago & Comprobante</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Control de ingresos, adjunto de voucher SumUp y activación de membresía.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4 max-h-[calc(85vh-120px)] overflow-y-auto">
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-sm">
          {{ errorMessage }}
        </div>

        <!-- Alumno -->
        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Alumno a Registrar *
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
              {{ student.fullName }} ({{ student.frequency }} · {{ student.membershipStatus }})
            </option>
          </select>
        </div>

        <!-- Resumen de Membresía Única -->
        <div v-if="selectedStudent" class="p-3 bg-[#FAF9F6] border border-[#E2DFD8] rounded-md text-xs space-y-1">
          <div class="flex justify-between">
            <span class="text-[#6B6D74]">Membresía Única:</span>
            <span class="font-bold text-[#191A1D]">Membresía Mensual</span>
          </div>
          <div class="flex justify-between">
            <span class="text-[#6B6D74]">Modalidad de Cobro:</span>
            <span class="font-semibold text-[#191A1D]">Pago Mensual (del 1 al 5)</span>
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

        <!-- Método de Pago y Fecha -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Método de Pago *
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
              Fecha de Recepción *
            </label>
            <input
              v-model="paymentDate"
              type="date"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>
        </div>

        <!-- Monto de Mensualidad Única -->
        <div>
          <div class="flex items-center justify-between mb-1">
            <label class="text-xs font-semibold text-[#6B6D74] uppercase tracking-wider">
              Monto Mensualidad (CLP) *
            </label>
            <span class="text-[11px] text-[#8E9097]">
              Pago Mensual
            </span>
          </div>
          <input
            v-model.number="amount"
            type="number"
            placeholder="45000"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] font-bold focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <!-- N° Operación SumUp / Código Transferencia -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              N° Operación SumUp / Banco
            </label>
            <input
              v-model="transactionReference"
              type="text"
              placeholder="Ej. SUM-849102 o TRF-8829"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] font-mono focus:outline-none focus:border-[#CEA549]"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Notas u Observaciones
            </label>
            <input
              v-model="notes"
              type="text"
              placeholder="Ej. Comprobante recibido por WhatsApp..."
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] focus:outline-none focus:border-[#CEA549]"
            />
          </div>
        </div>

        <!-- Adjuntar Comprobante / Voucher de SumUp -->
        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1.5">
            Adjuntar Voucher / Comprobante de SumUp (Opcional)
          </label>

          <!-- Dropzone para adjuntar archivo -->
          <div
            v-if="!voucherUrl"
            class="border-2 border-dashed border-[#DCD8CF] hover:border-[#CEA549] rounded-lg p-4 text-center bg-[#FAF9F6] transition-colors cursor-pointer"
            @click="fileInputRef?.click()"
          >
            <Upload class="w-5 h-5 mx-auto text-[#8C6B1B] mb-1" />
            <p class="text-xs font-bold text-[#191A1D]">Subir captura o foto del voucher de SumUp</p>
            <p class="text-[11px] text-[#8E9097] mt-0.5">Formatos: PNG, JPG o PDF de SumUp o transferencia</p>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*,application/pdf"
              class="hidden"
              @change="handleFileUpload"
            />
          </div>

          <!-- Preview de Comprobante Adjuntado -->
          <div v-else class="p-3 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg flex items-center justify-between gap-3">
            <div class="flex items-center gap-2.5 truncate">
              <img
                v-if="voucherUrl.startsWith('data:image')"
                :src="voucherUrl"
                class="w-12 h-12 object-cover rounded border border-[#DCD8CF] shrink-0"
                alt="Vista previa voucher"
              />
              <FileText v-else class="w-8 h-8 text-[#CEA549] shrink-0" />
              <div class="truncate">
                <span class="text-xs font-bold text-[#191A1D] block truncate">{{ voucherFileName }}</span>
                <span class="text-[11px] text-emerald-700 font-semibold flex items-center gap-1">
                  <Check class="w-3 h-3" /> Voucher adjuntado con éxito
                </span>
              </div>
            </div>
            <button
              type="button"
              @click="removeVoucher"
              class="px-2.5 py-1.5 bg-white hover:bg-rose-50 text-rose-700 border border-rose-200 rounded text-xs font-bold transition-colors cursor-pointer flex items-center gap-1 shrink-0"
            >
              <Trash2 class="w-3.5 h-3.5" /> Quitar
            </button>
          </div>
        </div>

        <div class="p-3 bg-[#FAF9F6] rounded-md text-xs text-[#6B6D74] border border-[#E2DFD8]">
          <strong class="text-[#191A1D]">Registro Contable:</strong> Al guardar el pago, se confirma la membresía del alumno para el ciclo activo, se guarda el voucher de SumUp y queda disponible el comprobante oficial.
        </div>

        <!-- Footer -->
        <div class="pt-3 border-t border-[#E2DFD8] flex items-center justify-end gap-3">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#6B6D74] hover:bg-[#F5F3EE] hover:text-[#191A1D] font-medium transition-colors cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold transition-colors shadow-xs cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
          >
            <Check class="w-4 h-4" />
            {{ isSubmitting ? 'Guardando...' : 'Registrar Pago & Voucher' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
