<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import ManageTermsModal from '@/components/modals/ManageTermsModal.vue'
import type { ContractConsent, ContractStatus, Student } from '@/types'
import {
  FileText,
  CheckCircle2,
  Clock,
  Printer,
  Share2,
  Check,
  ShieldCheck,
  AlertCircle
} from 'lucide-vue-next'

const router = useRouter()
const gymStore = useGymStore()

const showManageTermsModal = ref(false)
const showDocViewerModal = ref(false)
const selectedContract = ref<ContractConsent | null>(null)
const copiedContractId = ref<string | null>(null)

const searchName = ref('')
const selectedStatus = ref('')
const selectedDocType = ref('')

onMounted(async () => {
  await gymStore.loadStudents()
})

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

function getStudent(studentId: string): Student | undefined {
  return gymStore.students.find(s => s.id === studentId)
}

const selectedStudent = computed(() => {
  if (!selectedContract.value) return null
  return gymStore.students.find(s => s.id === selectedContract.value!.studentId) || null
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

async function copyConsentLink(contract: ContractConsent) {
  const url = `${window.location.origin}/registro-alumno`
  const text = `Hola ${contract.studentName}! Te comparto el enlace para que firmes digitalmente tu Consentimiento Informado y Reglamento de INVERFIT desde tu celular:\n${url}`
  try {
    await navigator.clipboard.writeText(text)
    copiedContractId.value = contract.id
    gymStore.showToast(`Enlace para ${contract.studentName} copiado. ¡Pégalo en WhatsApp!`, 'success', 'Enlace copiado')
    setTimeout(() => {
      copiedContractId.value = null
    }, 3000)
  } catch {
    gymStore.showToast(`Enlace: ${url}`, 'info')
  }
}

function printContract() {
  window.print()
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
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Contratos y consentimientos</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Controla documentos vigentes, firmas digitales y aceptación de términos de los alumnos.</p>
      </div>

      <button
        @click="showManageTermsModal = true"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        <FileText class="w-4 h-4" />
        Gestionar términos
      </button>
    </div>

    <!-- 3 KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Pendientes de aceptación</span>
          <Clock class="w-4 h-4 text-[#B06000]" />
        </div>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingContractsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Alumnos por firmar consentimiento</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Aceptados y vigentes</span>
          <CheckCircle2 class="w-4 h-4 text-[#137333]" />
        </div>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.acceptedContractsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Consentimientos con respaldo activo</span>
      </div>

      <!-- Card 3 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-[#6B6D74]">Versión vigente</span>
          <ShieldCheck class="w-4 h-4 text-[#8C6B1B]" />
        </div>
        <span class="text-[20px] font-bold text-[#191A1D] leading-tight truncate my-1">
          {{ currentVersion }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Reglamento oficial INVERFIT</span>
      </div>
    </div>

    <!-- Toolbar: Search + Filtro Estado + Filtro Documento -->
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

      <div class="flex items-center gap-3 w-full md:w-auto justify-end flex-wrap">
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
          <option value="Consentimiento informado y reglamento">Consentimiento informado y reglamento</option>
          <option value="Términos y condiciones">Términos y condiciones</option>
          <option value="Consentimiento informado">Consentimiento informado</option>
        </select>
      </div>
    </div>

    <!-- Table: Alumno | Documento | Estado | Fecha | Acción -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[780px]">
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
                Acciones
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-[#ECE8E1]">
            <tr v-if="filteredContracts.length === 0">
              <td colspan="5" class="py-12 text-center text-sm text-[#8E9097]">
                No se encontraron consentimientos o contratos registrados.
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
                  <div class="w-9 h-9 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[13px] font-bold shrink-0">
                    {{ contract.studentName.charAt(0) }}
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <span class="group-hover:text-black font-bold text-[#191A1D]">{{ contract.studentName }}</span>
                      <span
                        v-if="getStudent(contract.studentId)?.onboarding?.digitalSignature"
                        class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-[#E6F4EA] text-[#137333] border border-[#137333]/20"
                        title="Firma digital táctil registrada"
                      >
                        ✍️ Firma digital
                      </span>
                    </div>
                    <p class="text-[11px] text-[#8E9097] leading-none mt-1">
                      {{ getStudent(contract.studentId)?.email || 'Sin correo registrado' }}
                    </p>
                  </div>
                </div>
              </td>

              <!-- Documento -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                <p class="font-semibold text-[#191A1D]">{{ contract.documentType }}</p>
                <p class="text-[11px] text-[#8E9097] mt-0.5">Versión: {{ contract.termsVersion }}</p>
              </td>

              <!-- Estado -->
              <td class="py-4 px-6 text-[13px]">
                <span
                  class="px-2.5 py-1 rounded-[4px] text-[12px] font-semibold inline-flex items-center gap-1.5"
                  :class="contract.status === 'Aceptado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="contract.status === 'Aceptado' ? 'bg-[#137333]' : 'bg-[#B06000]'"></span>
                  {{ contract.status }}
                </span>
              </td>

              <!-- Fecha -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52] whitespace-nowrap">
                {{ contract.date }}
              </td>

              <!-- Acciones -->
              <td class="py-4 px-6 text-[13px] font-semibold text-right whitespace-nowrap" @click.stop>
                <div class="flex items-center justify-end gap-2">
                  <!-- Botón WhatsApp para alumnos pendientes -->
                  <button
                    v-if="contract.status === 'Pendiente'"
                    @click="copyConsentLink(contract)"
                    class="h-[30px] px-2.5 rounded border border-[#DCD8CF] hover:bg-white text-xs font-semibold text-[#191A1D] flex items-center gap-1.5 transition-colors cursor-pointer"
                    title="Copiar enlace para WhatsApp del alumno"
                  >
                    <component :is="copiedContractId === contract.id ? Check : Share2" class="w-3.5 h-3.5" :class="copiedContractId === contract.id ? 'text-[#137333]' : 'text-[#8C6B1B]'" />
                    <span>{{ copiedContractId === contract.id ? 'Copiado!' : 'Link WhatsApp' }}</span>
                  </button>

                  <!-- Toggle Status Button -->
                  <button
                    @click="handleToggleStatus(contract)"
                    class="h-[30px] px-2.5 rounded border border-[#DCD8CF] hover:bg-white text-xs font-semibold text-[#191A1D] transition-colors cursor-pointer"
                    :title="contract.status === 'Aceptado' ? 'Marcar como pendiente' : 'Registrar aceptación presencial'"
                  >
                    {{ contract.status === 'Aceptado' ? 'Marcar pendiente' : 'Marcar aceptado' }}
                  </button>

                  <!-- Ver Documento -->
                  <button
                    @click="viewContract(contract)"
                    class="h-[30px] px-3 bg-[#191A1D] hover:bg-[#2A2B2F] text-white rounded text-xs font-bold transition-colors cursor-pointer flex items-center gap-1"
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

    <!-- Informational Card: Respaldo legal y consentimientos -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-2">
      <div class="flex items-center gap-2">
        <ShieldCheck class="w-5 h-5 text-[#8C6B1B]" />
        <h3 class="text-[14px] font-bold text-[#191A1D]">Respaldo legal y consentimientos informados</h3>
      </div>
      <p class="text-[12px] text-[#6B6D74] leading-relaxed">
        Todos los términos, políticas de clase y declaraciones de salud son aceptados y firmados digitalmente por los alumnos durante su registro de onboarding, quedando debidamente archivados bajo la versión vigente del reglamento interno de INVERFIT. Si un alumno aún no firma, puedes compartirle el enlace de WhatsApp para que estampe su firma táctil desde su teléfono móvil.
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
      <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
        <!-- Modal Top Bar -->
        <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-[#191A1D] border border-[#CEA549]/40 flex items-center justify-center text-[#CEA549] font-black text-xs">
              IN
            </div>
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
              <p class="text-xs text-[#6B6D74] mt-0.5">Versión oficial: {{ selectedContract.termsVersion }} · Fecha: {{ selectedContract.date }}</p>
            </div>
          </div>

          <button
            @click="showDocViewerModal = false"
            class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
          >
            ✕
          </button>
        </div>

        <!-- Student Data Strip -->
        <div class="px-6 py-3 bg-[#FAF9F6] border-b border-[#E2DFD8] grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
          <div>
            <span class="text-[#8E9097] block text-[11px]">Alumno:</span>
            <span class="font-bold text-[#191A1D]">{{ selectedContract.studentName }}</span>
          </div>
          <div>
            <span class="text-[#8E9097] block text-[11px]">RUT / Identificación:</span>
            <span class="font-semibold text-[#191A1D]">{{ selectedStudent?.onboarding?.rut || 'Sin registrar' }}</span>
          </div>
          <div>
            <span class="text-[#8E9097] block text-[11px]">Teléfono:</span>
            <span class="font-semibold text-[#191A1D]">{{ selectedStudent?.phone || 'Sin registrar' }}</span>
          </div>
          <div>
            <span class="text-[#8E9097] block text-[11px]">Contacto de Emergencia:</span>
            <span class="font-semibold text-[#191A1D] truncate block">{{ selectedStudent?.emergencyContact || 'Sin registrar' }}</span>
          </div>
        </div>

        <!-- Document Body -->
        <div class="p-6 space-y-4 max-h-[55vh] overflow-y-auto text-xs text-[#191A1D] leading-relaxed bg-white font-sans">
          <div class="border-b border-[#ECE8E1] pb-3 text-center">
            <h3 class="text-sm font-bold uppercase tracking-wider text-[#191A1D]">
              Consentimiento Informado, Declaración de Salud y Reglamento del Servicio
            </h3>
            <p class="text-[11px] text-[#6B6D74] mt-0.5">Centro de Entrenamiento Personalizado INVERFIT · Coach Pablo</p>
          </div>

          <div class="space-y-3 text-[12px] text-[#3A3B40]">
            <div>
              <p class="font-bold text-[#191A1D]">1. DECLARACIÓN DE SALUD Y APTITUD FÍSICA (PAR-Q+)</p>
              <p class="mt-0.5 text-justify">
                El alumno declara bajo fe de juramento encontrarse en condiciones médicas adecuadas para realizar actividad física progresiva y entrenamiento de fuerza. Declara haber informado con total veracidad antecedentes quirúrgicos, dolencias articulares o condiciones crónicas en su ficha de ingreso, eximiendo a INVERFIT y a su entrenador principal de cualquier responsabilidad derivada de la omisión o falsedad de antecedentes clínicos.
              </p>
            </div>

            <div>
              <p class="font-bold text-[#191A1D]">2. ORGANIZACIÓN EN BLOQUES FIJOS DE 75 MINUTOS</p>
              <p class="mt-0.5 text-justify">
                Las sesiones de entrenamiento se organizan en bloques fijos de 75 minutos con un aforo técnico riguroso de máximo 4 alumnos por bloque para garantizar atención personalizada. El alumno se compromete a asistir puntualmente. Las inasistencias no son acumulables ni aplazan el valor de la mensualidad, salvo aviso previo con al menos 24 horas y reasignación excepcional sujeta a cupos disponibles en la misma semana.
              </p>
            </div>

            <div>
              <p class="font-bold text-[#191A1D]">3. MEMBRESÍA MENSUAL Y CICLO DE RECAUDACIÓN (DEL 1 AL 5)</p>
              <p class="mt-0.5 text-justify">
                El servicio opera bajo una membresía mensual única con ciclo oficial de pago del 1 al 5 de cada mes calendario. Los pagos deben regularizarse a través del enlace permanente oficial de SumUp o transferencia bancaria acreditada con comprobante. El no pago dentro del plazo faculta a INVERFIT a liberar el cupo en el bloque horario.
              </p>
            </div>

            <div>
              <p class="font-bold text-[#191A1D]">4. RESGUARDO DE DATOS Y SEGUIMIENTO</p>
              <p class="mt-0.5 text-justify">
                El alumno autoriza el registro de sus evaluaciones físicas, percentiles y métricas corporales para fines exclusivos de monitoreo técnico, planificación de sobrecarga y resguardo deportivo confidencial.
              </p>
            </div>
          </div>

          <!-- Digital Signature Box -->
          <div
            v-if="selectedStudent?.onboarding?.digitalSignature"
            class="mt-4 p-4 rounded-lg bg-[#FAF9F6] border-2 border-[#CEA549]/50 shadow-2xs space-y-2"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2 text-[#137333] font-bold text-xs">
                <ShieldCheck class="w-4 h-4 text-[#137333]" />
                <span>Firmado digitalmente vía Web Onboarding INVERFIT</span>
              </div>
              <span class="text-[11px] text-[#6B6D74]">{{ selectedStudent.onboarding.consentTimestamp }}</span>
            </div>

            <div class="p-2 bg-white rounded border border-[#E2DFD8] flex items-center justify-center">
              <img
                :src="selectedStudent.onboarding.digitalSignature"
                alt="Firma digital del alumno"
                class="max-h-24 object-contain"
              />
            </div>
            <p class="text-[10px] text-[#8E9097] text-center">
              Firma electrónica estampada en pantalla táctil con hash de registro y confirmación biométrica.
            </p>
          </div>

          <!-- Unsigned Info Box -->
          <div
            v-else
            class="mt-4 p-4 rounded-lg bg-[#FEF7E0]/60 border border-[#B06000]/30 space-y-2 text-xs"
          >
            <div class="flex items-center gap-2 text-[#B06000] font-bold">
              <AlertCircle class="w-4 h-4" />
              <span>Pendiente de firma digital táctil</span>
            </div>
            <p class="text-[#6B6D74]">
              Este alumno no tiene una firma digital estampada en el portal. Puedes enviarle el enlace para que firme en su teléfono o marcar la aceptación si ya firmó presencialmente.
            </p>
            <button
              @click="copyConsentLink(selectedContract)"
              class="px-3 py-1.5 bg-[#191A1D] hover:bg-[#2A2B2F] text-[#CEA549] rounded text-xs font-bold transition-colors cursor-pointer flex items-center gap-1.5"
            >
              <Share2 class="w-3.5 h-3.5" />
              <span>Copiar mensaje y enlace para WhatsApp</span>
            </button>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-4 bg-[#FAF9F6] flex flex-col sm:flex-row items-center justify-between gap-3 border-t border-[#ECE8E1]">
          <div class="flex items-center gap-3">
            <button
              @click="goToStudent(selectedContract.studentId); showDocViewerModal = false"
              class="text-xs text-[#8C6B1B] hover:text-[#614910] font-bold underline cursor-pointer"
            >
              Ir a la ficha del alumno →
            </button>

            <button
              @click="printContract"
              class="text-xs text-[#191A1D] hover:text-black font-semibold flex items-center gap-1 cursor-pointer"
              title="Imprimir contrato"
            >
              <Printer class="w-3.5 h-3.5 text-[#8C6B1B]" />
              <span>Imprimir / PDF</span>
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="handleToggleStatus(selectedContract)"
              class="px-3.5 py-2 rounded text-xs font-bold transition-colors cursor-pointer"
              :class="selectedContract.status === 'Aceptado' ? 'bg-[#FEF7E0] text-[#B06000] hover:bg-[#FEEFC3]' : 'bg-[#E6F4EA] text-[#137333] hover:bg-[#CEEAD6]'"
            >
              {{ selectedContract.status === 'Aceptado' ? 'Cambiar a Pendiente' : 'Marcar como Aceptado' }}
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
