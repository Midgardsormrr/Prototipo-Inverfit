<script setup lang="ts">
import { ref } from 'vue'
import { Copy, Check, Send, AlertCircle, FileCheck } from 'lucide-vue-next'
import { useGymStore } from '@/stores/gym'

const props = defineProps<{
  open: boolean
  mode: 'link' | 'receipt' | 'reminder'
  studentName: string
  linkUrl?: string
  paymentId?: string
  amount?: number | null
  method?: string
  date?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const gymStore = useGymStore()
const copied = ref(false)
const actionDone = ref(false)

function copyLink() {
  if (props.linkUrl) {
    navigator.clipboard.writeText(props.linkUrl)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2500)
  }
}

async function simulateSendEmail() {
  if (props.paymentId) {
    await gymStore.simulateReceipt(props.paymentId)
  } else {
    gymStore.showToast(`Comprobante simulado enviado al correo de ${props.studentName}.`, 'success')
  }
  actionDone.value = true
}

async function simulateSendReminder() {
  if (props.paymentId) {
    await gymStore.simulateReminder(props.paymentId)
  } else {
    gymStore.showToast(`Recordatorio simulado enviado a ${props.studentName}.`, 'info')
  }
  actionDone.value = true
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-md w-full border border-[#E2DFD8] overflow-hidden animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-[#CEA549] bg-[#191A1D] px-2 py-0.5 rounded">
            Simulación Demostrativa
          </span>
          <h2 class="text-lg font-bold text-[#191A1D] mt-1.5">
            <span v-if="mode === 'link'">Enviar Enlace de Pago</span>
            <span v-else-if="mode === 'receipt'">Comprobante de Pago</span>
            <span v-else>Enviar Recordatorio</span>
          </h2>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 space-y-4">
        <!-- MODE: PAYMENT LINK -->
        <div v-if="mode === 'link'" class="space-y-3">
          <p class="text-sm text-[#4A4C52]">
            Se ha generado un enlace de pago seguro para <strong>{{ studentName }}</strong>.
          </p>

          <div class="p-3 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] flex items-center justify-between gap-2">
            <code class="text-xs text-[#191A1D] truncate font-mono select-all">
              {{ linkUrl || 'https://demo.inverfit.invalid/checkout/9A4DF7' }}
            </code>
            <button
              @click="copyLink"
              class="px-2.5 py-1.5 bg-white border border-[#CEA549] rounded text-xs font-bold text-[#191A1D] hover:bg-[#CEA549]/10 flex items-center gap-1 transition-colors flex-shrink-0 cursor-pointer shadow-2xs"
            >
              <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-600" />
              <Copy v-else class="w-3.5 h-3.5 text-[#8C6B1B]" />
              {{ copied ? 'Copiado' : 'Copiar' }}
            </button>
          </div>

          <div class="p-3 bg-emerald-50 text-emerald-900 border border-emerald-200 rounded-lg text-xs space-y-1">
            <p class="font-semibold flex items-center gap-1.5">
              <FileCheck class="w-4 h-4 text-emerald-700" /> Mensaje para enviar al alumno:
            </p>
            <p class="italic text-emerald-800">
              "Hola {{ studentName }}, te comparto tu enlace para renovar tu membresía de INVERFIT: {{ linkUrl }}. ¡Cualquier consulta me avisas!"
            </p>
          </div>
        </div>

        <!-- MODE: RECEIPT -->
        <div v-else-if="mode === 'receipt'" class="space-y-3">
          <div class="p-4 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg space-y-2 text-xs">
            <div class="flex justify-between border-b pb-2 border-[#E2DFD8] font-bold text-[#191A1D]">
              <span class="flex items-center gap-1">INVERFIT <span class="text-[#CEA549]">❯❯</span></span>
              <span>BOLETA / COMPROBANTE</span>
            </div>
            <div class="grid grid-cols-2 gap-1 text-[#6B6D74]">
              <span>Alumno:</span>
              <span class="font-semibold text-[#191A1D]">{{ studentName }}</span>
              <span>Medio:</span>
              <span class="font-semibold text-[#191A1D]">{{ method || 'Transferencia' }}</span>
              <span>Fecha:</span>
              <span class="font-semibold text-[#191A1D]">{{ date || '2026-11-04' }}</span>
              <span>Monto:</span>
              <span class="font-semibold text-[#191A1D]">{{ amount ? `$${amount.toLocaleString('es-CL')}` : 'Por definir' }}</span>
            </div>
            <div class="border-t pt-2 border-[#E2DFD8] text-[11px] text-[#8E9097] text-center">
              Emisión centralizada de comprobantes (simulación para evitar Gmail manual).
            </div>
          </div>

          <div v-if="!actionDone" class="pt-2">
            <button
              @click="simulateSendEmail"
              class="w-full py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <Send class="w-4 h-4" />
              Simular envío de comprobante al alumno
            </button>
          </div>
          <div v-else class="p-3 bg-emerald-50 text-emerald-800 rounded-md text-xs font-semibold text-center border border-emerald-200">
            ✓ Comprobante enviado en simulación exitosamente.
          </div>
        </div>

        <!-- MODE: REMINDER -->
        <div v-else class="space-y-3">
          <p class="text-sm text-[#4A4C52]">
            Envío de recordatorio de vencimiento de membresía para <strong>{{ studentName }}</strong>.
          </p>

          <div class="p-3 bg-[#FAF9F6] rounded-lg text-xs text-[#555861] italic border border-[#E2DFD8]">
            "Estimado(a) {{ studentName }}, te recordamos que tu membresía en INVERFIT está próxima a vencer. Por favor comunícate con el coach Pablo para confirmar tus bloques."
          </div>

          <div v-if="!actionDone" class="pt-2">
            <button
              @click="simulateSendReminder"
              class="w-full py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <Send class="w-4 h-4" />
              Simular envío de recordatorio
            </button>
          </div>
          <div v-else class="p-3 bg-emerald-50 text-emerald-800 rounded-md text-xs font-semibold text-center border border-emerald-200">
            ✓ Recordatorio enviado en simulación.
          </div>
        </div>

        <!-- Prototype boundary alert -->
        <div class="p-3 bg-[#FEF7E0] border border-[#F0D878] rounded-md flex items-start gap-2 text-xs text-[#8C6B1B]">
          <AlertCircle class="w-4 h-4 flex-shrink-0 mt-0.5 text-[#B06000]" />
          <span>
            Esta acción es una simulación visual del prototipo. No se envían correos reales ni se realizan transacciones bancarias en esta etapa.
          </span>
        </div>

        <!-- Footer -->
        <div class="pt-2 border-t border-[#E2DFD8] flex justify-end">
          <button
            @click="emit('close')"
            class="px-4 py-2 bg-[#F5F3EE] hover:bg-[#EAE6DD] text-[#191A1D] rounded-md text-sm font-semibold transition-colors cursor-pointer border border-[#DDD8CE]"
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
