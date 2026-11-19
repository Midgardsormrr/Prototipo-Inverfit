<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Mail, Check, Send, Users, ShieldCheck, Copy } from 'lucide-vue-next'
import { useGymStore } from '@/stores/gym'
import type { Payment } from '@/types'

const props = defineProps<{
  open: boolean
  pendingPayments: Payment[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'dispatched'): void
}>()

const gymStore = useGymStore()

const selectedPaymentIds = ref<string[]>([])
const isProcessing = ref(false)
const dispatchComplete = ref(false)
const copiedTemplate = ref(false)

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    selectedPaymentIds.value = props.pendingPayments.map(p => p.id)
    dispatchComplete.value = false
    copiedTemplate.value = false
    isProcessing.value = false
  }
})

const allSelected = computed({
  get: () => selectedPaymentIds.value.length === props.pendingPayments.length && props.pendingPayments.length > 0,
  set: (val: boolean) => {
    selectedPaymentIds.value = val ? props.pendingPayments.map(p => p.id) : []
  }
})

const selectedPayments = computed(() => {
  return props.pendingPayments.filter(p => selectedPaymentIds.value.includes(p.id))
})

// Current billing cycle info
const today = new Date()
const currentDay = today.getDate()
const isWithinWindow = currentDay >= 1 && currentDay <= 5

const cycleNoticeSubject = `INVERFIT — Recordatorio de Renovación de Membresía Mensual (Período 1 al 5)`

const cycleNoticeBody = computed(() => {
  return `Estimados alumnos de INVERFIT,

Les recordamos que nos encontramos en el período oficial de recaudación y renovación de membresías mensuales, correspondiente a los días 1 al 5 de cada mes.

Para mantener activos sus bloques fijos de entrenamiento y asegurar la continuidad en la agenda, les solicitamos realizar su pago dentro de esta ventana.

CANALES OFICIALES DE PAGO:

1. Tarjetas de Débito y Crédito (Enlace permanente SumUp):
${gymStore.sumupPermanentLink}

2. Transferencia Bancaria Directa:
Banco: Banco de Chile
Tipo de Cuenta: Cuenta Corriente
N° de Cuenta: 00-12345678-01
Nombre: INVERFIT - Pablo
RUT: 12.345.678-9
Correo para comprobantes: pagos@inverfit.cl

3. Efectivo en Box de Entrenamiento

Una vez realizado su pago, por favor remitan el comprobante para validar su renovación en el sistema.

Atentamente,
Coach Pablo — INVERFIT
Centro de Entrenamiento de Fuerza y Rendimiento`
})

function getStudentEmail(studentId: string): string {
  const student = gymStore.students.find(s => s.id === studentId)
  return student?.email || 'correo@ejemplo.com'
}

async function handleDispatchBatch() {
  if (selectedPaymentIds.value.length === 0) {
    gymStore.showToast('Selecciona al menos un alumno para enviar los avisos.', 'warning')
    return
  }

  isProcessing.value = true

  try {
    for (const id of selectedPaymentIds.value) {
      await gymStore.simulateReminder(id)
    }

    dispatchComplete.value = true
    gymStore.showToast(
      `Avisos registrados exitosamente para ${selectedPaymentIds.value.length} alumno(s).`,
      'success',
      'Despacho Completado'
    )
    emit('dispatched')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error al despachar avisos'
    gymStore.showToast(msg, 'error')
  } finally {
    isProcessing.value = false
  }
}

function openBatchEmailClient() {
  const emails = selectedPayments.value.map(p => getStudentEmail(p.studentId)).filter(Boolean)
  if (emails.length === 0) return

  // Use BCC (Copia Oculta) so recipient emails are never exposed to each other
  const bccList = encodeURIComponent(emails.join(','))
  const mailto = `mailto:pagos@inverfit.cl?bcc=${bccList}&subject=${encodeURIComponent(cycleNoticeSubject)}&body=${encodeURIComponent(cycleNoticeBody.value)}`
  window.location.href = mailto

  handleDispatchBatch()
}

function copyTemplate() {
  navigator.clipboard.writeText(cycleNoticeBody.value)
  copiedTemplate.value = true
  setTimeout(() => { copiedTemplate.value = false }, 2500)
  gymStore.showToast('Mensaje de aviso copiado al portapapeles.', 'success')
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full border border-[#E2DFD8] overflow-hidden my-6 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <div class="flex items-center gap-2">
            <span class="text-[10px] font-bold uppercase tracking-wider text-[#CEA549] bg-[#191A1D] px-2 py-0.5 rounded">
              Ciclo Oficial 1 al 5
            </span>
            <span class="text-xs text-[#8E9097] font-semibold">
              Coach Pablo · Automatización de Cobros
            </span>
          </div>
          <h2 class="text-lg font-bold text-[#191A1D] mt-1.5 flex items-center gap-2">
            <Mail class="w-5 h-5 text-[#CEA549]" />
            Despacho de Avisos Masivos de Cobro
          </h2>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 space-y-4 max-h-[calc(85vh-120px)] overflow-y-auto">
        <!-- Banner del Ciclo del 1 al 5 -->
        <div
          class="p-4 rounded-lg border flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs"
          :class="isWithinWindow ? 'bg-[#FEF7E0] border-[#F0D878] text-[#8C6B1B]' : 'bg-[#FAF9F6] border-[#E2DFD8] text-[#4A4C52]'"
        >
          <div class="space-y-1">
            <div class="flex items-center gap-2 font-bold text-[#191A1D] text-sm">
              <span class="w-2.5 h-2.5 rounded-full" :class="isWithinWindow ? 'bg-[#CEA549] animate-pulse' : 'bg-[#8E9097]'"></span>
              <span>
                {{ isWithinWindow ? `Ventana Oficial Activa (Día ${currentDay} del mes)` : `Próximo Ciclo Oficial: Días 1 al 5` }}
              </span>
            </div>
            <p class="text-[#6B6D74] leading-relaxed">
              Fecha límite oficial de pago: <strong>Día 5</strong>. Los avisos se envían a los alumnos activos que mantienen mensualidad pendiente.
            </p>
          </div>

          <div class="sm:text-right shrink-0">
            <span class="text-[11px] text-[#8E9097] block">Pendientes de cobro</span>
            <span class="text-2xl font-bold text-[#191A1D]">{{ pendingPayments.length }} alumnos</span>
          </div>
        </div>

        <!-- Enlace permanente SumUp activo -->
        <div class="p-3 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] flex items-center justify-between gap-2 text-xs">
          <div class="flex items-center gap-2 truncate">
            <ShieldCheck class="w-4 h-4 text-[#CEA549] shrink-0" />
            <span class="text-[#6B6D74]">Link Permanente SumUp:</span>
            <code class="font-mono text-[#191A1D] font-bold truncate">{{ gymStore.sumupPermanentLink }}</code>
          </div>
          <button
            @click="copyTemplate"
            type="button"
            class="px-2 py-1 bg-white border border-[#DCD8CF] hover:border-[#CEA549] rounded text-[11px] font-semibold text-[#191A1D] flex items-center gap-1 shrink-0 transition-colors cursor-pointer"
          >
            <Copy class="w-3 h-3 text-[#8C6B1B]" />
            {{ copiedTemplate ? 'Copiado' : 'Copiar Texto' }}
          </button>
        </div>

        <!-- Lista de Alumnos Pendientes con Checkboxes -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="text-xs font-semibold text-[#6B6D74] uppercase tracking-wider flex items-center gap-1.5">
              <Users class="w-4 h-4 text-[#CEA549]" /> Alumnos con Cobro Pendiente ({{ selectedPaymentIds.length }} de {{ pendingPayments.length }} seleccionados)
            </label>

            <button
              @click="allSelected = !allSelected"
              class="text-xs text-[#8C6B1B] hover:text-[#614910] font-bold hover:underline cursor-pointer"
            >
              {{ allSelected ? 'Desmarcar todos' : 'Seleccionar todos' }}
            </button>
          </div>

          <div class="border border-[#E2DFD8] rounded-lg divide-y divide-[#ECE8E1] max-h-48 overflow-y-auto bg-white">
            <div
              v-for="p in pendingPayments"
              :key="p.id"
              class="p-3 flex items-center justify-between gap-3 hover:bg-[#FAF9F6] transition-colors cursor-pointer"
              @click="selectedPaymentIds.includes(p.id) ? selectedPaymentIds = selectedPaymentIds.filter(id => id !== p.id) : selectedPaymentIds.push(p.id)"
            >
              <div class="flex items-center gap-3">
                <input
                  type="checkbox"
                  :value="p.id"
                  v-model="selectedPaymentIds"
                  @click.stop
                  class="rounded border-[#DCD8CF] text-[#CEA549] focus:ring-[#CEA549] cursor-pointer"
                />
                <div>
                  <span class="font-bold text-xs text-[#191A1D]">{{ p.studentName }}</span>
                  <p class="text-[11px] text-[#8E9097] leading-none mt-0.5">
                    {{ getStudentEmail(p.studentId) }} · {{ p.frequency }}
                  </p>
                </div>
              </div>

              <div class="text-right shrink-0">
                <span class="text-xs font-bold text-[#191A1D]">
                  {{ p.amount ? `$${p.amount.toLocaleString('es-CL')}` : 'Por definir' }}
                </span>
                <span class="text-[10px] text-[#8E9097] block">
                  {{ p.reminderSentCount > 0 ? `${p.reminderSentCount} aviso(s) previo(s)` : 'Sin avisos' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Previsualización del Correo de Cobro Masivo -->
        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Plantilla Oficial de Aviso (Ciclo 1 al 5)
          </label>
          <div class="p-3.5 bg-[#FAF9F6] rounded-lg text-xs text-[#4A4C52] border border-[#E2DFD8] font-mono whitespace-pre-wrap max-h-36 overflow-y-auto leading-relaxed select-all">
{{ cycleNoticeBody }}
          </div>
        </div>

        <!-- Banner de éxito si ya se despacharon -->
        <div v-if="dispatchComplete" class="p-3 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-md text-xs font-semibold flex items-center gap-2">
          <Check class="w-4 h-4 text-emerald-700" />
          <span>¡Avisos registrados exitosamente en el historial de cobros de los alumnos seleccionados!</span>
        </div>

        <!-- Botones de Acción Masiva -->
        <div class="space-y-2.5 pt-2">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <button
              @click="openBatchEmailClient"
              type="button"
              :disabled="selectedPaymentIds.length === 0 || isProcessing"
              class="w-full py-2.5 px-4 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
            >
              <Send class="w-4 h-4" />
              <span>Abrir en Correo Masivo (BCC Oculto)</span>
            </button>

            <button
              @click="handleDispatchBatch"
              type="button"
              :disabled="selectedPaymentIds.length === 0 || isProcessing"
              class="w-full py-2.5 px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
            >
              <Check class="w-4 h-4 text-[#8C6B1B]" />
              <span>{{ isProcessing ? 'Procesando...' : `Registrar Avisos (${selectedPaymentIds.length})` }}</span>
            </button>
          </div>

          <p class="text-[11px] text-[#8E9097] text-center">
            * Al abrir en correo, se genera un mensaje con <strong>Copia Oculta (BCC)</strong> para resguardar la privacidad de los alumnos.
          </p>
        </div>

        <!-- Footer -->
        <div class="pt-3 border-t border-[#E2DFD8] flex items-center justify-between">
          <span class="text-[11px] text-[#8E9097]">
            INVERFIT · Módulo de Recaudación Oficial
          </span>
          <button
            @click="emit('close')"
            class="px-4 py-2 bg-[#F5F3EE] hover:bg-[#EAE6DD] text-[#191A1D] rounded-md text-xs font-semibold transition-colors cursor-pointer border border-[#DDD8CE]"
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
