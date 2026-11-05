<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import ManageTermsModal from '@/components/modals/ManageTermsModal.vue'
import type { ContractConsent, ContractStatus } from '@/types'

const router = useRouter()
const gymStore = useGymStore()

const showManageTermsModal = ref(false)
const showDocViewerModal = ref(false)
const selectedContract = ref<ContractConsent | null>(null)

const searchName = ref('')
const selectedStatus = ref('')
const selectedDocType = ref('')

const filteredContracts = computed(() => {
  return gymStore.contracts.filter(c => {
    const q = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (q && !c.studentName.toLowerCase().includes(q)) return false
    if (selectedStatus.value && c.status !== selectedStatus.value) return false
    if (selectedDocType.value && c.documentType !== selectedDocType.value) return false
    return true
  })
})

const currentVersion = computed(() => {
  const current = gymStore.terms.find(t => t.isCurrent)
  return current ? current.version : 'v1.0 (2026)'
})

function viewContract(contract: ContractConsent) {
  selectedContract.value = contract
  showDocViewerModal.value = true
}

async function handleToggleStatus(contract: ContractConsent) {
  const nextStatus: ContractStatus = contract.status === 'Aceptado' ? 'Pendiente' : 'Aceptado'
  await gymStore.updateContractStatus(contract.id, nextStatus)
  if (selectedContract.value && selectedContract.value.id === contract.id) {
    selectedContract.value.status = nextStatus
  }
}

function goToStudent(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header matching wireframe 07 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Contratos y consentimientos</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Controla documentos vigentes y aceptación de términos.</p>
      </div>

      <button
        @click="showManageTermsModal = true"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        Gestionar términos
      </button>
    </div>

    <!-- 3 KPI Summary Cards matching wireframe 07: [x=278, y=193], [x=623, y=193], [x=968, y=193] -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pendientes de aceptación</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingContractsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Solicitudes por completar</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Aceptados</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.acceptedContractsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Consentimientos registrados</span>
      </div>

      <!-- Card 3 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Versión vigente</span>
        <span class="text-[22px] font-bold text-[#191A1D] leading-tight truncate my-1">
          {{ currentVersion }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Por confirmar con el cliente</span>
      </div>
    </div>

    <!-- Toolbar: Search + Estado de aceptación + Tipo de documento matching wireframe -->
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
        <select
          v-model="selectedStatus"
          class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Estado de aceptación  ▾</option>
          <option value="Aceptado">Aceptado</option>
          <option value="Pendiente">Pendiente</option>
        </select>

        <select
          v-model="selectedDocType"
          class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Tipo de documento  ▾</option>
          <option value="Términos y condiciones">Términos y condiciones</option>
          <option value="Consentimiento informado">Consentimiento informado</option>
          <option value="Contrato de servicios">Contrato de servicios</option>
        </select>
      </div>
    </div>

    <!-- Table matching wireframe 07: Alumno | Documento | Estado | Fecha | Acción -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[750px]">
          <thead>
            <tr class="border-b border-[#E2DFD8] bg-[#FAF9F6]">
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Alumno
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Documento
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Estado
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Fecha
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider text-right">
                Acción
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-[#ECE8E1]">
            <tr v-if="filteredContracts.length === 0">
              <td colspan="5" class="py-12 text-center text-sm text-[#8E9097]">
                No se encontraron contratos registrados.
              </td>
            </tr>

            <tr
              v-for="contract in filteredContracts"
              :key="contract.id"
              class="hover:bg-[#FAF9F6] transition-colors group cursor-pointer"
              @click="viewContract(contract)"
            >
              <!-- Alumno -->
              <td class="py-4 px-6 text-[13px] font-medium text-[#191A1D]">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[12px] font-bold">
                    {{ contract.studentName.charAt(0) }}
                  </div>
                  <div>
                    <span class="group-hover:text-black font-bold text-[#191A1D]">{{ contract.studentName }}</span>
                    <p class="text-[11px] text-[#8E9097] leading-none mt-0.5">
                      Versión: {{ contract.termsVersion }}
                    </p>
                  </div>
                </div>
              </td>

              <!-- Documento -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ contract.documentType }}
              </td>

              <!-- Estado -->
              <td class="py-4 px-6 text-[13px]">
                <span
                  class="px-2.5 py-1 rounded-[4px] text-[12px] font-semibold inline-block"
                  :class="contract.status === 'Aceptado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                >
                  {{ contract.status }}
                </span>
              </td>

              <!-- Fecha -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ contract.date }}
              </td>

              <!-- Acción -->
              <td class="py-4 px-6 text-[13px] font-semibold text-right whitespace-nowrap" @click.stop>
                <div class="flex items-center justify-end gap-3">
                  <button
                    @click="handleToggleStatus(contract)"
                    class="text-xs px-2.5 py-1 rounded border border-[#DCD8CF] hover:bg-white font-semibold text-[#191A1D] transition-colors cursor-pointer"
                    title="Simular cambio de estado de aceptación"
                  >
                    {{ contract.status === 'Aceptado' ? 'Marcar pendiente' : 'Aceptar (Demo)' }}
                  </button>

                  <button
                    @click="viewContract(contract)"
                    class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer font-bold"
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

    <!-- Informational Card matching wireframe 07: [x=280, y=718] and [x=280, y=742] -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-1">
      <h3 class="text-[14px] font-bold text-[#191A1D]">Documento y respaldo</h3>
      <p class="text-[12px] text-[#6B6D74]">
        El mecanismo de firma, versionado y resguardo legal debe validarse antes de implementar. Esta sección es una simulación visual y documental de demostración.
      </p>
    </div>

    <!-- Modals -->
    <ManageTermsModal
      :open="showManageTermsModal"
      @close="showManageTermsModal = false"
    />

    <!-- Document Viewer Modal -->
    <div
      v-if="showDocViewerModal && selectedContract"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs overflow-y-auto"
      @click.self="showDocViewerModal = false"
    >
      <div class="bg-white rounded-xl shadow-2xl max-w-xl w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
          <div>
            <div class="flex items-center gap-2">
              <span
                class="px-2 py-0.5 rounded text-[11px] font-bold"
                :class="selectedContract.status === 'Aceptado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
              >
                {{ selectedContract.status }}
              </span>
              <h2 class="text-base font-bold text-[#191A1D]">{{ selectedContract.documentType }}</h2>
            </div>
            <p class="text-xs text-[#6B6D74] mt-0.5">Alumno: {{ selectedContract.studentName }} · Versión: {{ selectedContract.termsVersion }}</p>
          </div>
          <button
            @click="showDocViewerModal = false"
            class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
          >
            ✕
          </button>
        </div>

        <div class="p-6 space-y-4 max-h-[65vh] overflow-y-auto text-xs text-[#191A1D] leading-relaxed font-mono bg-[#FAF9F6] border-y border-[#E2DFD8] whitespace-pre-wrap">
CONTRATO Y TÉRMINOS DEL SERVICIO — INVERFIT
Alumno: {{ selectedContract.studentName }}
Estado de aceptación: {{ selectedContract.status }}
Fecha de registro: {{ selectedContract.date }}

1. DECLARACIÓN DE SALUD Y APTITUD
El alumno declara encontrarse médicamente apto para el entrenamiento impartido en INVERFIT y exime al entrenador de responsabilidad por omisión de antecedentes médicos.

2. CUMPLIMIENTO DE HORARIOS Y POLÍTICA DE CLASES
Los cupos asignados en la agenda son fijos por semana. Las cancelaciones no aplazan la membresía salvo disponibilidad y coordinación directa con el entrenador en la misma semana.
        </div>

        <div class="p-4 bg-white flex items-center justify-between gap-3 border-t border-[#ECE8E1]">
          <button
            @click="goToStudent(selectedContract.studentId); showDocViewerModal = false"
            class="text-xs text-[#8C6B1B] hover:text-[#614910] font-bold underline cursor-pointer"
          >
            Ir a la ficha del alumno →
          </button>

          <div class="flex items-center gap-2">
            <button
              @click="handleToggleStatus(selectedContract)"
              class="px-3.5 py-2 rounded text-xs font-bold transition-colors cursor-pointer"
              :class="selectedContract.status === 'Aceptado' ? 'bg-[#FEF7E0] text-[#B06000] hover:bg-[#FEEFC3]' : 'bg-[#E6F4EA] text-[#137333] hover:bg-[#CEEAD6]'"
            >
              {{ selectedContract.status === 'Aceptado' ? 'Cambiar a Pendiente' : 'Aceptar en Demo' }}
            </button>
            <button
              @click="showDocViewerModal = false"
              class="px-4 py-2 bg-[#191A1D] text-white rounded text-xs font-semibold hover:bg-[#24262B] transition-colors cursor-pointer"
            >
              Cerrar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
