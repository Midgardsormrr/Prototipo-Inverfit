<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { paymentService } from '@/services'
import { Download, Send } from 'lucide-vue-next'
import RegisterPaymentModal from '@/components/modals/RegisterPaymentModal.vue'
import SimulatedActionModal from '@/components/modals/SimulatedActionModal.vue'
import type { Payment } from '@/types'

const router = useRouter()
const gymStore = useGymStore()

const showRegisterModal = ref(false)
const showSimModal = ref(false)
const simMode = ref<'link' | 'receipt' | 'reminder'>('link')
const selectedPayment = ref<Payment | null>(null)
const generatedLink = ref('')

const searchName = ref('')
const filterStatus = ref<string>('')

const filteredPayments = computed(() => {
  return gymStore.payments.filter(p => {
    const q = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (q && !p.studentName.toLowerCase().includes(q)) return false
    if (filterStatus.value && p.status !== filterStatus.value) return false
    return true
  })
})

function handleSendLink(p: Payment) {
  selectedPayment.value = p
  simMode.value = 'link'
  const token = Math.random().toString(36).substring(2, 9).toUpperCase()
  generatedLink.value = `https://demo.inverfit.invalid/checkout/${token}`
  showSimModal.value = true
}

function handleSendReceipt(p: Payment) {
  selectedPayment.value = p
  simMode.value = 'receipt'
  showSimModal.value = true
}

function handleSendRemindersBatch() {
  const pendingCount = gymStore.pendingPaymentsCount
  if (pendingCount === 0) {
    gymStore.showToast('No hay pagos pendientes para enviar recordatorios.', 'info')
    return
  }
  selectedPayment.value = gymStore.payments.find(p => p.status === 'Pendiente') || null
  simMode.value = 'reminder'
  showSimModal.value = true
}

function exportCSV() {
  const csv = paymentService.exportToCSV()
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `inverfit_pagos_${new Date().toISOString().split('T')[0]}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  gymStore.showToast('Archivo CSV exportado exitosamente.', 'success')
}

function goToStudent(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header matching wireframe 05 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Membresías y pagos</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Controla vencimientos, estados y enlaces de pago.</p>
      </div>

      <button
        @click="showRegisterModal = true"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        + Registrar pago
      </button>
    </div>

    <!-- 4 KPI Summary Cards matching wireframe 05 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Por vencer</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.expiringMembershipsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Revisar fechas de membresía</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pagos pendientes</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingPaymentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Enlace no pagado / por revisar</span>
      </div>

      <!-- Card 3 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pagos registrados</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.paidPaymentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Historial del período</span>
      </div>

      <!-- Card 4: Acciones with "Enviar recordatorios" button inside matching wireframe -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs">
        <span class="text-[13px] font-semibold text-[#6B6D74]">Acciones</span>
        <button
          @click="handleSendRemindersBatch"
          class="h-[36px] w-full px-3 bg-[#FAF9F6] border border-[#DCD8CF] hover:bg-[#F5F3EE] hover:border-[#CEA549]/60 rounded-[6px] text-[13px] font-bold text-[#191A1D] transition-colors flex items-center justify-center gap-1.5 cursor-pointer shadow-2xs"
        >
          <Send class="w-3.5 h-3.5 text-[#8C6B1B]" />
          Enviar recordatorios
        </button>
        <span class="text-[11px] text-[#8E9097] truncate">Simulación de avisos</span>
      </div>
    </div>

    <!-- Toolbar: Search + Filtrar + Exportar matching wireframe 05 -->
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

        <!-- Exportar button matching wireframe -->
        <button
          @click="exportCSV"
          class="h-[40px] px-4 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-bold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549]/60 transition-colors flex items-center gap-1.5 cursor-pointer shadow-2xs"
        >
          <Download class="w-4 h-4 text-[#8C6B1B]" />
          Exportar
        </button>
      </div>
    </div>

    <!-- Table matching wireframe 05: Alumno | Plan | Vencimiento | Monto | Estado | Acción -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[750px]">
          <thead>
            <tr class="border-b border-[#E2DFD8] bg-[#FAF9F6]">
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Alumno
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Plan
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Vencimiento
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
            <tr v-if="filteredPayments.length === 0">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                No se encontraron registros de pago.
              </td>
            </tr>

            <tr
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
                    <p class="text-[11px] text-[#8E9097] leading-none mt-0.5">
                      {{ payment.paymentMethod }} · {{ payment.reminderSentCount > 0 ? `${payment.reminderSentCount} recordatorio(s)` : 'Sin recordatorios' }}
                    </p>
                  </div>
                </div>
              </td>

              <!-- Plan -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ payment.frequency }}
              </td>

              <!-- Vencimiento -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ payment.dueDate }}
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
                <div class="flex items-center justify-end gap-2">
                  <button
                    v-if="payment.status === 'Pendiente'"
                    @click="handleSendLink(payment)"
                    class="px-2.5 py-1 text-xs border border-[#CEA549] bg-[#CEA549]/10 rounded hover:bg-[#CEA549]/20 font-bold text-[#191A1D] transition-colors cursor-pointer"
                    title="Simular envío de enlace de pago"
                  >
                    Enviar link
                  </button>

                  <button
                    v-if="payment.status === 'Pagado'"
                    @click="handleSendReceipt(payment)"
                    class="px-2.5 py-1 text-xs border border-[#DCD8CF] rounded hover:bg-white font-semibold text-[#191A1D] transition-colors cursor-pointer"
                    title="Simular envío de comprobante"
                  >
                    Comprobante
                  </button>

                  <button
                    @click="goToStudent(payment.studentId)"
                    class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer font-bold ml-1"
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

    <!-- Informational Notes matching wireframe 05: [x=260, y=690] and [x=260, y=718] -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-2">
      <p class="text-[13px] text-[#4A4C52]">
        <strong class="text-[#191A1D]">Medios mencionados por el cliente:</strong> Webpay, Servipag, crédito, débito y transferencia.
      </p>
      <p class="text-[12px] text-[#8E9097]">
        El proveedor y la confirmación automática del pago están pendientes de definición. Las acciones de enlace y comprobante operan como simulaciones funcionales.
      </p>
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
      :link-url="generatedLink"
      :payment-id="selectedPayment?.id"
      :amount="selectedPayment?.amount"
      :method="selectedPayment?.paymentMethod"
      :date="selectedPayment?.paymentDate || undefined"
      @close="showSimModal = false"
    />
  </div>
</template>
