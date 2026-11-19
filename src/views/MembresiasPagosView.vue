<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { paymentService } from '@/services'
import { Download, RefreshCw, AlertCircle, CreditCard, Mail, FileCheck, Calendar, Settings, Paperclip } from 'lucide-vue-next'
import RegisterPaymentModal from '@/components/modals/RegisterPaymentModal.vue'
import SimulatedActionModal from '@/components/modals/SimulatedActionModal.vue'
import BatchRemindersModal from '@/components/modals/BatchRemindersModal.vue'
import ConfigureSumupLinkModal from '@/components/modals/ConfigureSumupLinkModal.vue'
import type { Payment } from '@/types'

const router = useRouter()
const gymStore = useGymStore()

const showRegisterModal = ref(false)
const showSimModal = ref(false)
const showBatchModal = ref(false)
const showSumupConfigModal = ref(false)

const simMode = ref<'link' | 'receipt' | 'reminder'>('link')
const selectedPayment = ref<Payment | null>(null)
const generatedLink = ref('')

const searchName = ref('')
const filterStatus = ref<string>('')

onMounted(() => {
  if (gymStore.isPaymentsApiMode()) {
    gymStore.loadPayments()
  }
})

const filteredPayments = computed(() => {
  return gymStore.payments.filter(p => {
    const q = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (q && !p.studentName.toLowerCase().includes(q)) return false
    if (filterStatus.value && p.status !== filterStatus.value) return false
    return true
  })
})

const pendingPaymentsList = computed(() => {
  return gymStore.payments.filter(p => p.status === 'Pendiente')
})

const selectedStudent = computed(() => {
  if (!selectedPayment.value) return null
  return gymStore.students.find(s => s.id === selectedPayment.value?.studentId) || null
})

function handleSendLink(p: Payment) {
  selectedPayment.value = p
  simMode.value = 'link'
  // Use Coach Pablo's permanent SumUp link
  generatedLink.value = gymStore.sumupPermanentLink
  showSimModal.value = true
}

function handleSendReminder(p: Payment) {
  selectedPayment.value = p
  simMode.value = 'reminder'
  // Use Coach Pablo's permanent SumUp link
  generatedLink.value = gymStore.sumupPermanentLink
  showSimModal.value = true
}

function handleSendReceipt(p: Payment) {
  selectedPayment.value = p
  simMode.value = 'receipt'
  generatedLink.value = gymStore.sumupPermanentLink
  showSimModal.value = true
}

function handleSendRemindersBatch() {
  if (pendingPaymentsList.value.length === 0) {
    gymStore.showToast('No hay pagos pendientes para enviar avisos en este ciclo.', 'info')
    return
  }
  showBatchModal.value = true
}

async function exportCSV() {
  try {
    const csv = await paymentService.exportToCSV()
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.setAttribute('href', url)
    link.setAttribute('download', `inverfit_pagos_${new Date().toISOString().split('T')[0]}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    gymStore.showToast('Archivo CSV exportado exitosamente.', 'success')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error al exportar CSV'
    gymStore.showToast(msg, 'error')
  }
}

function goToStudent(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Membresías y Pagos</h1>
          <span
            class="px-2.5 py-0.5 rounded-full text-[11px] font-bold border tracking-wide"
            :class="gymStore.paymentsDataSource === 'api' ? 'bg-[#E6F4EA] text-[#137333] border-[#137333]/30' : 'bg-[#F1EFEA] text-[#6B6D74] border-[#DCD8CF]'"
          >
            {{ gymStore.paymentsDataSource === 'api' ? 'API REST (Go)' : 'Mock (localStorage)' }}
          </span>
        </div>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">
          Ciclo oficial de cobro del 1 al 5 · Link permanente SumUp y avisos por correo.
        </p>
      </div>

      <div class="flex items-center gap-2.5 self-start sm:self-auto flex-wrap">
        <button
          @click="showSumupConfigModal = true"
          class="h-[40px] px-3.5 bg-white border border-[#DCD8CF] hover:border-[#CEA549] rounded-[6px] text-[13px] font-bold text-[#191A1D] flex items-center gap-1.5 transition-colors cursor-pointer shadow-2xs"
          title="Configurar Link Permanente SumUp"
        >
          <CreditCard class="w-3.5 h-3.5 text-[#8C6B1B]" />
          Link SumUp
        </button>

        <button
          @click="gymStore.loadPayments()"
          class="h-[40px] px-3.5 bg-white border border-[#DCD8CF] hover:border-[#CEA549] rounded-[6px] text-[13px] font-bold text-[#191A1D] flex items-center gap-1.5 transition-colors cursor-pointer shadow-2xs"
          title="Recargar pagos"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': gymStore.loadingPayments }" />
          Recargar
        </button>

        <button
          @click="showRegisterModal = true"
          class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs hover:shadow"
        >
          + Registrar pago
        </button>
      </div>
    </div>

    <!-- Error Alert Banner -->
    <div
      v-if="gymStore.paymentsError"
      class="p-4 bg-rose-50 border border-rose-200 rounded-[10px] flex items-center justify-between shadow-2xs"
    >
      <div class="flex items-center gap-2.5 text-rose-800 text-sm">
        <AlertCircle class="w-5 h-5 shrink-0" />
        <span>{{ gymStore.paymentsError }}</span>
      </div>
      <button
        @click="gymStore.loadPayments()"
        class="px-3 py-1.5 bg-rose-100 hover:bg-rose-200 text-rose-900 rounded-[6px] text-xs font-bold transition-colors cursor-pointer"
      >
        Reintentar
      </button>
    </div>

    <!-- 4 KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Período 1 al 5</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ pendingPaymentsList.length }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Ventana oficial de cobro</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pagos pendientes</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingPaymentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Sin registrar este ciclo</span>
      </div>

      <!-- Card 3 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pagos al día</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.paidPaymentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Membresías activas</span>
      </div>

      <!-- Card 4: Acciones (Despacho Masivo) -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs">
        <span class="text-[13px] font-semibold text-[#6B6D74]">Avisos Masivos</span>
        <button
          @click="handleSendRemindersBatch"
          class="h-[36px] w-full px-3 bg-[#FAF9F6] border border-[#DCD8CF] hover:bg-[#F5F3EE] hover:border-[#CEA549]/60 rounded-[6px] text-[13px] font-bold text-[#191A1D] transition-colors flex items-center justify-center gap-1.5 cursor-pointer shadow-2xs"
        >
          <Mail class="w-3.5 h-3.5 text-[#8C6B1B]" />
          Avisos del 1 al 5
        </button>
        <span class="text-[11px] text-[#8E9097] truncate">{{ pendingPaymentsList.length }} pendientes este mes</span>
      </div>
    </div>

    <!-- Toolbar: Search + Filtrar + Exportar -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col md:flex-row items-center justify-between gap-4 shadow-2xs">
      <div class="w-full md:w-[320px] relative">
        <input
          v-model="searchName"
          type="text"
          placeholder="⌕  Buscar alumno..."
          class="w-full h-[40px] px-3.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] text-[#191A1D] placeholder-[#8E9097] focus:outline-none focus:border-[#CEA549] focus:bg-white focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
        />
        <button
          v-if="searchName"
          @click="searchName = ''"
          class="absolute right-2.5 top-1/2 -translate-y-1/2 text-xs text-[#8E9097] hover:text-[#191A1D]"
        >
          ✕
        </button>
      </div>

      <div class="flex items-center gap-3 w-full md:w-auto justify-end">
        <!-- Filtrar dropdown -->
        <select
          v-model="filterStatus"
          class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Filtrar: Todos los estados</option>
          <option value="Pagado">Pagado</option>
          <option value="Pendiente">Pendiente</option>
        </select>

        <!-- Exportar button -->
        <button
          @click="exportCSV"
          class="h-[40px] px-4 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-bold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549]/60 transition-colors flex items-center gap-1.5 cursor-pointer shadow-2xs"
        >
          <Download class="w-4 h-4 text-[#8C6B1B]" />
          Exportar
        </button>
      </div>
    </div>

    <!-- Table: Alumno | Membresía | Período de Cobro | Monto | Estado | Acción -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[780px]">
          <thead>
            <tr class="border-b border-[#E2DFD8] bg-[#FAF9F6]">
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Alumno
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Membresía
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider whitespace-nowrap">
                Período de Cobro
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Monto
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Estado
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider text-right">
                Acción
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-[#ECE8E1]">
            <tr v-if="gymStore.loadingPayments">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                <div class="flex items-center justify-center gap-2">
                  <RefreshCw class="w-4 h-4 animate-spin text-[#CEA549]" />
                  <span>Cargando registros de pagos...</span>
                </div>
              </td>
            </tr>

            <tr v-else-if="filteredPayments.length === 0">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                No se encontraron registros de pago.
              </td>
            </tr>

            <tr
              v-else
              v-for="payment in filteredPayments"
              :key="payment.id"
              class="hover:bg-[#FAF9F6] transition-colors group cursor-pointer"
              @click="goToStudent(payment.studentId)"
            >
              <!-- Alumno -->
              <td class="py-4 px-6 text-[13px] font-medium text-[#191A1D]">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[12px] font-bold">
                    {{ payment.studentName.charAt(0) }}
                  </div>
                  <div>
                    <span class="group-hover:text-black font-bold text-[#191A1D]">{{ payment.studentName }}</span>
                    <p class="text-[11px] text-[#8E9097] leading-none mt-0.5 flex items-center gap-1.5 flex-wrap">
                      <span>{{ payment.paymentMethod }}</span>
                      <span>·</span>
                      <span>{{ payment.reminderSentCount > 0 ? `${payment.reminderSentCount} aviso(s) enviado(s)` : 'Sin avisos' }}</span>
                      <span v-if="payment.voucherUrl" class="text-[10px] bg-emerald-50 text-emerald-800 border border-emerald-200 px-1.5 py-0.5 rounded font-semibold flex items-center gap-0.5">
                        <Paperclip class="w-3 h-3" /> Voucher
                      </span>
                    </p>
                  </div>
                </div>
              </td>

              <!-- Membresía (Única mensual + frecuencia) -->
              <td class="py-4 px-6 text-[13px]">
                <div class="flex flex-col">
                  <span class="font-bold text-[#191A1D] text-[13px]">Membresía Mensual</span>
                  <span class="text-[11px] text-[#8E9097]">{{ payment.frequency }}</span>
                </div>
              </td>

              <!-- Período de Cobro (1 al 5) con icono estilizado y whitespace-nowrap -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                <div class="flex items-center gap-1.5 whitespace-nowrap">
                  <Calendar class="w-4 h-4 text-[#CEA549] shrink-0" />
                  <span class="font-medium text-[#191A1D]">{{ payment.dueDate || '1 al 5 de cada mes' }}</span>
                </div>
              </td>

              <!-- Monto -->
              <td class="py-4 px-6 text-[13px] font-semibold text-[#191A1D]">
                {{ payment.amount ? `$${payment.amount.toLocaleString('es-CL')}` : 'Por definir' }}
              </td>

              <!-- Estado -->
              <td class="py-4 px-6 text-[13px]">
                <span
                  class="px-2.5 py-1 rounded-[4px] text-[12px] font-semibold inline-block"
                  :class="payment.status === 'Pagado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                >
                  {{ payment.status }}
                </span>
              </td>

              <!-- Acción -->
              <td class="py-4 px-6 text-[13px] font-semibold text-right whitespace-nowrap" @click.stop>
                <div class="flex items-center justify-end gap-1.5">
                  <!-- Acciones para estado Pendiente -->
                  <template v-if="payment.status === 'Pendiente'">
                    <button
                      @click="handleSendLink(payment)"
                      class="px-2.5 py-1 text-xs border border-[#CEA549] bg-[#CEA549]/15 hover:bg-[#CEA549]/25 text-[#191A1D] rounded font-bold transition-colors cursor-pointer flex items-center gap-1 shadow-2xs"
                      title="Copiar o enviar el enlace permanente de SumUp"
                    >
                      <CreditCard class="w-3.5 h-3.5 text-[#8C6B1B]" />
                      Link SumUp
                    </button>

                    <button
                      @click="handleSendReminder(payment)"
                      class="px-2.5 py-1 text-xs border border-[#DCD8CF] hover:border-[#CEA549] bg-white hover:bg-[#FAF9F6] text-[#191A1D] rounded font-bold transition-colors cursor-pointer flex items-center gap-1 shadow-2xs"
                      title="Enviar aviso individual de cobro al correo del alumno"
                    >
                      <Mail class="w-3.5 h-3.5 text-[#8C6B1B]" />
                      Aviso Correo
                    </button>
                  </template>

                  <!-- Acciones para estado Pagado -->
                  <template v-if="payment.status === 'Pagado'">
                    <button
                      @click="handleSendReceipt(payment)"
                      class="px-2.5 py-1 text-xs border border-[#DCD8CF] hover:border-emerald-600/50 bg-white hover:bg-[#FAF9F6] rounded font-bold text-[#191A1D] transition-colors cursor-pointer flex items-center gap-1 shadow-2xs"
                      title="Ver y enviar comprobante de pago"
                    >
                      <FileCheck class="w-3.5 h-3.5 text-emerald-700" />
                      Comprobante
                    </button>
                  </template>

                  <button
                    @click="goToStudent(payment.studentId)"
                    class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer font-bold ml-1 text-xs"
                  >
                    Ver →
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Política Oficial INVERFIT de Membresías y Cobranza -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4">
      <div class="flex items-center justify-between border-b border-[#ECE8E1] pb-3">
        <div class="flex items-center gap-2">
          <span class="w-2.5 h-2.5 rounded-full bg-[#CEA549]"></span>
          <h3 class="text-sm font-bold text-[#191A1D] uppercase tracking-wide">
            Política Oficial de Membresía & Recaudación INVERFIT
          </h3>
        </div>
        <span class="text-xs font-bold text-[#CEA549] bg-[#191A1D] px-2.5 py-0.5 rounded">
          Coach Pablo
        </span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
        <div class="p-3 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] space-y-1">
          <span class="font-bold text-[#191A1D] block">1. Membresía Única Mensual</span>
          <p class="text-[#6B6D74] leading-relaxed">
            En INVERFIT existe una sola membresía mensual enfocada en entrenamiento de fuerza y longevidad, tarificada exclusivamente según la frecuencia semanal acordada (1, 2, 3 o 4 veces por semana).
          </p>
        </div>

        <div class="p-3 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] space-y-1">
          <span class="font-bold text-[#191A1D] block">2. Ventana de Pago del 1 al 5</span>
          <p class="text-[#6B6D74] leading-relaxed">
            La recaudación oficial se realiza puntualmente entre los días 1 y 5 de cada mes calendario para confirmar y reservar los bloques fijos de entrenamiento del alumno en la agenda.
          </p>
        </div>

        <!-- Columna 3 con acceso directo al Link Permanente -->
        <div class="p-3 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] space-y-1.5">
          <div class="flex items-center justify-between">
            <span class="font-bold text-[#191A1D] block">3. Pasarela SumUp & Avisos</span>
            <button
              @click="showSumupConfigModal = true"
              type="button"
              class="text-[11px] text-[#8C6B1B] font-bold hover:underline cursor-pointer flex items-center gap-1"
            >
              <Settings class="w-3 h-3" />
              Editar Link
            </button>
          </div>
          <p class="text-[#6B6D74] leading-relaxed">
            Link permanente activo:
            <code class="font-mono text-[#191A1D] font-bold block truncate mt-0.5">{{ gymStore.sumupPermanentLink }}</code>
          </p>
          <p class="text-[11px] text-[#8E9097] leading-relaxed">
            Este enlace oficial se inserta automáticamente en todos los correos de cobro y recordatorios de WhatsApp.
          </p>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <RegisterPaymentModal
      :open="showRegisterModal"
      @close="showRegisterModal = false"
    />

    <SimulatedActionModal
      :open="showSimModal"
      :mode="simMode"
      :student-name="selectedPayment?.studentName || ''"
      :student-email="selectedStudent?.email || ''"
      :student-phone="selectedStudent?.phone || ''"
      :link-url="generatedLink"
      :payment-id="selectedPayment?.id"
      :amount="selectedPayment?.amount"
      :frequency="selectedPayment?.frequency"
      :due-date="selectedPayment?.dueDate"
      :method="selectedPayment?.paymentMethod"
      :date="selectedPayment?.paymentDate || undefined"
      :transaction-reference="selectedPayment?.transactionReference"
      :voucher-url="selectedPayment?.voucherUrl"
      :voucher-file-name="selectedPayment?.voucherFileName"
      @close="showSimModal = false"
    />

    <BatchRemindersModal
      :open="showBatchModal"
      :pending-payments="pendingPaymentsList"
      @close="showBatchModal = false"
      @dispatched="gymStore.loadPayments()"
    />

    <ConfigureSumupLinkModal
      :open="showSumupConfigModal"
      @close="showSumupConfigModal = false"
    />
  </div>
</template>
