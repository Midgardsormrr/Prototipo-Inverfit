<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useGymStore } from '@/stores/gym'
import { Send, Edit, AlertCircle } from 'lucide-vue-next'
import SimulatedActionModal from '@/components/modals/SimulatedActionModal.vue'
import RegisterPaymentModal from '@/components/modals/RegisterPaymentModal.vue'
import CreateEvaluationModal from '@/components/modals/CreateEvaluationModal.vue'
import type { FrequencyType, StudentStatus } from '@/types'

const props = defineProps<{
  id: string
}>()

const gymStore = useGymStore()

const currentTab = ref<'resumen' | 'datos' | 'pagos' | 'evaluaciones' | 'contrato'>('resumen')
const isEditing = ref(false)

// Modals
const showLinkModal = ref(false)
const showPaymentModal = ref(false)
const showEvalModal = ref(false)
const generatedLink = ref('')

const student = computed(() => {
  const s = gymStore.students.find(item => item.id === props.id)
  return s || gymStore.students[0] || null
})

// Edit form model
const editForm = ref({
  fullName: '',
  planName: '',
  frequency: '2 veces/sem.' as FrequencyType,
  membershipStatus: 'Activo' as StudentStatus,
  fixedSchedule: '',
  membershipStartDate: '',
  membershipDueDate: '',
  phone: '',
  email: '',
  emergencyContact: '',
  healthConditions: '',
  notes: ''
})

watch(student, (s) => {
  if (s) {
    editForm.value = {
      fullName: s.fullName,
      planName: s.planName,
      frequency: s.frequency,
      membershipStatus: s.membershipStatus,
      fixedSchedule: s.fixedSchedule,
      membershipStartDate: s.membershipStartDate,
      membershipDueDate: s.membershipDueDate,
      phone: s.phone,
      email: s.email,
      emergencyContact: s.emergencyContact || '',
      healthConditions: s.healthConditions || '',
      notes: s.notes || ''
    }
  }
}, { immediate: true })

const studentPayments = computed(() => {
  if (!student.value) return []
  return gymStore.payments.filter(p => p.studentId === student.value!.id)
})

const studentEvaluations = computed(() => {
  if (!student.value) return []
  return gymStore.evaluations.filter(e => e.studentId === student.value!.id)
})

const studentContract = computed(() => {
  if (!student.value) return null
  return gymStore.contracts.find(c => c.studentId === student.value!.id) || null
})

async function handleSaveEdit() {
  if (!student.value) return
  await gymStore.updateStudent(student.value.id, {
    fullName: editForm.value.fullName,
    planName: editForm.value.planName,
    frequency: editForm.value.frequency,
    membershipStatus: editForm.value.membershipStatus,
    fixedSchedule: editForm.value.fixedSchedule,
    membershipStartDate: editForm.value.membershipStartDate,
    membershipDueDate: editForm.value.membershipDueDate,
    phone: editForm.value.phone,
    email: editForm.value.email,
    emergencyContact: editForm.value.emergencyContact,
    healthConditions: editForm.value.healthConditions,
    notes: editForm.value.notes
  })
  isEditing.value = false
}

async function handleSendPaymentLink() {
  if (!student.value) return
  const res = await gymStore.simulatePaymentLink(student.value.id)
  generatedLink.value = res.link
  showLinkModal.value = true
}

async function toggleContractStatus() {
  if (!studentContract.value) return
  const nextStatus = studentContract.value.status === 'Aceptado' ? 'Pendiente' : 'Aceptado'
  await gymStore.updateContractStatus(studentContract.value.id, nextStatus)
}
</script>

<template>
  <div v-if="student" class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header matching wireframe 04 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Ficha del alumno</h1>
        <div class="flex items-center gap-2 text-[14px] text-[#6B6D74] mt-0.5">
          <router-link to="/alumnos" class="hover:underline text-[#8C6B1B] font-medium">Alumnos</router-link>
          <span>/</span>
          <span class="text-[#191A1D] font-bold">Perfil</span>
        </div>
      </div>

      <button
        @click="isEditing = !isEditing"
        class="h-[38px] px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-[6px] text-[13px] font-semibold transition-colors flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-2xs"
      >
        <Edit class="w-4 h-4 text-[#8C6B1B]" />
        {{ isEditing ? 'Ver perfil' : 'Editar información' }}
      </button>
    </div>

    <!-- Student Hero Profile Card matching wireframe 04 -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div class="flex items-center gap-4">
        <!-- Avatar circle matching wireframe: [x=300, y=214] 'A' with brand styling -->
        <div class="w-[56px] h-[56px] rounded-full bg-[#191A1D] border border-[#CEA549]/50 flex items-center justify-center text-[22px] font-bold text-[#CEA549] shadow-xs">
          {{ student.fullName.charAt(0) }}
        </div>

        <div>
          <h2 class="text-[19px] font-bold text-[#191A1D] leading-tight">
            {{ student.fullName }}
          </h2>
          <p class="text-[13px] text-[#6B6D74] mt-1">
            Membresía: <strong class="text-[#191A1D]">{{ student.planName }}</strong> · {{ student.frequency }}
          </p>
          <div class="mt-2 flex items-center gap-2">
            <span
              class="px-2.5 py-0.5 rounded-[4px] text-[12px] font-semibold"
              :class="student.membershipStatus === 'Activo' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
            >
              {{ student.membershipStatus }}
            </span>
            <span class="text-[12px] text-[#8E9097]">
              · Horario: {{ student.fixedSchedule }}
            </span>
          </div>
        </div>
      </div>

      <!-- Action: Enviar enlace de pago in Brand Gold -->
      <div class="flex flex-wrap items-center gap-3">
        <button
          @click="handleSendPaymentLink"
          class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center gap-2 cursor-pointer shadow-xs hover:shadow"
        >
          <Send class="w-4 h-4" />
          Enviar enlace de pago
        </button>
      </div>
    </div>

    <!-- Navigation Tabs matching wireframe: Resumen | Datos personales | Pagos | Evaluaciones | Contrato -->
    <div class="border-b border-[#E2DFD8] flex items-center gap-2 overflow-x-auto bg-white px-4 rounded-t-[10px]">
      <button
        @click="currentTab = 'resumen'"
        class="py-3 px-4 text-[13px] font-bold border-b-2 transition-colors cursor-pointer whitespace-nowrap"
        :class="currentTab === 'resumen' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
      >
        Resumen
      </button>

      <button
        @click="currentTab = 'datos'"
        class="py-3 px-4 text-[13px] font-bold border-b-2 transition-colors cursor-pointer whitespace-nowrap"
        :class="currentTab === 'datos' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
      >
        Datos personales
      </button>

      <button
        @click="currentTab = 'pagos'"
        class="py-3 px-4 text-[13px] font-bold border-b-2 transition-colors cursor-pointer whitespace-nowrap flex items-center gap-1.5"
        :class="currentTab === 'pagos' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
      >
        Pagos
        <span class="px-1.5 py-0.2 rounded-full bg-[#FAF9F6] border border-[#E2DFD8] text-[#191A1D] text-[10px] font-bold">
          {{ studentPayments.length }}
        </span>
      </button>

      <button
        @click="currentTab = 'evaluaciones'"
        class="py-3 px-4 text-[13px] font-bold border-b-2 transition-colors cursor-pointer whitespace-nowrap flex items-center gap-1.5"
        :class="currentTab === 'evaluaciones' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
      >
        Evaluaciones
        <span class="px-1.5 py-0.2 rounded-full bg-[#FAF9F6] border border-[#E2DFD8] text-[#191A1D] text-[10px] font-bold">
          {{ studentEvaluations.length }}
        </span>
      </button>

      <button
        @click="currentTab = 'contrato'"
        class="py-3 px-4 text-[13px] font-bold border-b-2 transition-colors cursor-pointer whitespace-nowrap"
        :class="currentTab === 'contrato' ? 'border-[#CEA549] text-[#191A1D]' : 'border-transparent text-[#6B6D74] hover:text-[#191A1D]'"
      >
        Contrato
      </button>
    </div>

    <!-- TAB 1: RESUMEN (Matches wireframe 04 layout with 4 cards) -->
    <div v-if="currentTab === 'resumen'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Card 1: Datos y plan -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4">
        <h3 class="text-[16px] font-bold text-[#191A1D]">Datos y plan</h3>

        <div class="divide-y divide-[#ECE8E1] text-sm">
          <div class="py-2.5 flex items-center justify-between">
            <span class="text-[12px] text-[#6B6D74]">Frecuencia</span>
            <span class="text-[13px] text-[#191A1D] font-semibold">{{ student.frequency }}</span>
          </div>

          <div class="py-2.5 flex items-center justify-between">
            <span class="text-[12px] text-[#6B6D74]">Horario fijo</span>
            <span class="text-[13px] text-[#191A1D] font-semibold">{{ student.fixedSchedule }}</span>
          </div>

          <div class="py-2.5 flex items-center justify-between">
            <span class="text-[12px] text-[#6B6D74]">Inicio de membresía</span>
            <span class="text-[13px] text-[#191A1D] font-semibold">{{ student.membershipStartDate }}</span>
          </div>

          <div class="py-2.5 flex items-center justify-between">
            <span class="text-[12px] text-[#6B6D74]">Vencimiento</span>
            <span class="text-[13px] text-[#191A1D] font-semibold">{{ student.membershipDueDate }}</span>
          </div>
        </div>

        <div class="pt-2">
          <button
            @click="currentTab = 'datos'"
            class="text-[13px] font-bold text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer"
          >
            Editar datos y plan →
          </button>
        </div>
      </div>

      <!-- Card 2: Contrato y consentimiento -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4 flex flex-col justify-between">
        <div>
          <h3 class="text-[16px] font-bold text-[#191A1D]">Contrato y consentimiento</h3>

          <div class="divide-y divide-[#ECE8E1] text-sm mt-4">
            <div class="py-2.5 flex items-center justify-between">
              <span class="text-[12px] text-[#6B6D74]">Estado de aceptación</span>
              <span
                class="px-2 py-0.5 rounded text-[12px] font-semibold"
                :class="studentContract?.status === 'Aceptado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
              >
                {{ studentContract?.status || 'Pendiente' }}
              </span>
            </div>

            <div class="py-2.5 flex items-center justify-between">
              <span class="text-[12px] text-[#6B6D74]">Versión de términos</span>
              <span class="text-[13px] text-[#191A1D] font-semibold">{{ studentContract?.termsVersion || 'v1.0 (2026)' }}</span>
            </div>

            <div class="py-2.5 flex items-center justify-between">
              <span class="text-[12px] text-[#6B6D74]">Documento firmado</span>
              <span class="text-[13px] text-[#191A1D] font-semibold">
                {{ studentContract?.signedCopyAvailable ? 'Disponible (Demostración)' : 'No disponible' }}
              </span>
            </div>
          </div>
        </div>

        <div class="pt-2">
          <button
            @click="currentTab = 'contrato'"
            class="h-[36px] px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] rounded-[6px] text-[13px] font-semibold text-[#191A1D] transition-colors cursor-pointer w-full text-center"
          >
            Ver / gestionar contrato
          </button>
        </div>
      </div>

      <!-- Card 3: Últimos pagos -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4 flex flex-col justify-between">
        <div>
          <h3 class="text-[16px] font-bold text-[#191A1D]">Últimos pagos</h3>

          <div v-if="studentPayments.length === 0" class="py-6 text-[13px] text-[#6B6D74]">
            Todavía no hay pagos cargados en este wireframe.
          </div>

          <div v-else class="divide-y divide-[#ECE8E1] mt-3">
            <div
              v-for="p in studentPayments.slice(0, 3)"
              :key="p.id"
              class="py-2.5 flex items-center justify-between text-xs"
            >
              <div>
                <p class="font-bold text-[#191A1D]">{{ p.paymentMethod }} · {{ p.dueDate }}</p>
                <p class="text-[#6B6D74] text-[11px]">{{ p.notes || 'Pago de membresía' }}</p>
              </div>
              <span
                class="px-2 py-0.5 rounded text-[11px] font-semibold"
                :class="p.status === 'Pagado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
              >
                {{ p.status }}
              </span>
            </div>
          </div>
        </div>

        <div class="pt-2 flex items-center gap-3">
          <button
            @click="currentTab = 'pagos'"
            class="h-[36px] px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] rounded-[6px] text-[13px] font-semibold text-[#191A1D] transition-colors cursor-pointer"
          >
            Ir a pagos
          </button>
          <button
            @click="showPaymentModal = true"
            class="h-[36px] px-3.5 bg-[#CEA549] text-[#0B0B0D] rounded-[6px] text-[12px] font-bold hover:bg-[#BA923B] transition-colors cursor-pointer shadow-xs"
          >
            + Registrar pago
          </button>
        </div>
      </div>

      <!-- Card 4: Evaluaciones y resultados -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4 flex flex-col justify-between">
        <div>
          <h3 class="text-[16px] font-bold text-[#191A1D]">Evaluaciones y resultados</h3>

          <div v-if="studentEvaluations.length === 0" class="py-6 text-[13px] text-[#6B6D74]">
            Registrar y consultar el seguimiento del alumno.
          </div>

          <div v-else class="divide-y divide-[#ECE8E1] mt-3">
            <div
              v-for="ev in studentEvaluations.slice(0, 3)"
              :key="ev.id"
              class="py-2.5 flex items-center justify-between text-xs"
            >
              <div>
                <p class="font-bold text-[#191A1D]">{{ ev.type }} · {{ ev.date }}</p>
                <p class="text-[#6B6D74] text-[11px] truncate max-w-[200px]">{{ ev.observations || 'Sin observaciones' }}</p>
              </div>
              <span class="px-2 py-0.5 rounded bg-[#FAF9F6] border border-[#E2DFD8] text-[#191A1D] font-semibold text-[11px]">
                {{ ev.status }}
              </span>
            </div>
          </div>
        </div>

        <div class="pt-2 flex items-center gap-3">
          <button
            @click="currentTab = 'evaluaciones'"
            class="h-[36px] px-4 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] rounded-[6px] text-[13px] font-semibold text-[#191A1D] transition-colors cursor-pointer"
          >
            Ver evaluaciones
          </button>
          <button
            @click="showEvalModal = true"
            class="h-[36px] px-3.5 bg-[#CEA549] text-[#0B0B0D] rounded-[6px] text-[12px] font-bold hover:bg-[#BA923B] transition-colors cursor-pointer shadow-xs"
          >
            + Nueva evaluación
          </button>
        </div>
      </div>
    </div>

    <!-- TAB 2: DATOS PERSONALES -->
    <div v-else-if="currentTab === 'datos'" class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
      <div class="flex items-center justify-between mb-4 border-b pb-3 border-[#ECE8E1]">
        <div>
          <h3 class="text-[16px] font-bold text-[#191A1D]">Datos Personales y Membresía</h3>
          <p class="text-xs text-[#6B6D74]">Consulta y modifica los datos de contacto y condiciones del alumno.</p>
        </div>
        <button
          v-if="!isEditing"
          @click="isEditing = true"
          class="px-3.5 py-1.5 border border-[#DCD8CF] rounded text-xs font-bold text-[#191A1D] hover:bg-[#FAF9F6] cursor-pointer"
        >
          Editar
        </button>
      </div>

      <form @submit.prevent="handleSaveEdit" class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Nombre Completo</label>
            <input
              v-model="editForm.fullName"
              type="text"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Frecuencia</label>
            <select
              v-model="editForm.frequency"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            >
              <option value="1 vez/sem.">1 vez/sem.</option>
              <option value="2 veces/sem.">2 veces/sem.</option>
              <option value="3 veces/sem.">3 veces/sem.</option>
              <option value="4 veces/sem.">4 veces/sem.</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Horario Fijo Semanal</label>
            <input
              v-model="editForm.fixedSchedule"
              type="text"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Estado de Membresía</label>
            <select
              v-model="editForm.membershipStatus"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            >
              <option value="Activo">Activo</option>
              <option value="Por revisar">Por revisar</option>
              <option value="Inactivo">Inactivo</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Teléfono</label>
            <input
              v-model="editForm.phone"
              type="text"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Correo Electrónico</label>
            <input
              v-model="editForm.email"
              type="email"
              :disabled="!isEditing"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Contacto de Emergencia</label>
          <input
            v-model="editForm.emergencyContact"
            type="text"
            :disabled="!isEditing"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Condiciones Médicas / Lesiones</label>
          <textarea
            v-model="editForm.healthConditions"
            rows="2"
            :disabled="!isEditing"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
          ></textarea>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">Notas del Entrenador</label>
          <textarea
            v-model="editForm.notes"
            rows="2"
            :disabled="!isEditing"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 disabled:bg-[#FAF9F6] disabled:border-[#ECE8E1]"
          ></textarea>
        </div>

        <div v-if="isEditing" class="pt-4 flex justify-end gap-3 border-t border-[#ECE8E1]">
          <button
            type="button"
            @click="isEditing = false"
            class="px-4 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#6B6D74] hover:bg-[#FAF9F6] cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="submit"
            class="px-5 py-2 bg-[#CEA549] text-[#0B0B0D] rounded-md text-sm font-bold hover:bg-[#BA923B] transition-colors cursor-pointer shadow-xs"
          >
            Guardar cambios
          </button>
        </div>
      </form>
    </div>

    <!-- TAB 3: PAGOS -->
    <div v-else-if="currentTab === 'pagos'" class="space-y-4">
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-[16px] font-bold text-[#191A1D]">Historial de Pagos</h3>
            <p class="text-xs text-[#6B6D74]">Control de membresías y emisión simulada de comprobantes.</p>
          </div>
          <button
            @click="showPaymentModal = true"
            class="px-4 py-2 bg-[#CEA549] text-[#0B0B0D] text-xs font-bold rounded hover:bg-[#BA923B] transition-colors cursor-pointer shadow-xs"
          >
            + Registrar pago
          </button>
        </div>

        <div v-if="studentPayments.length === 0" class="text-center py-10 text-sm text-[#8E9097]">
          No hay registros de pago para este alumno.
        </div>

        <div v-else class="divide-y divide-[#ECE8E1] border border-[#E2DFD8] rounded-lg overflow-hidden">
          <div
            v-for="p in studentPayments"
            :key="p.id"
            class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-[#FAF9F6]"
          >
            <div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-sm text-[#191A1D]">{{ p.planName }} ({{ p.frequency }})</span>
                <span
                  class="px-2 py-0.5 rounded text-[11px] font-semibold"
                  :class="p.status === 'Pagado' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                >
                  {{ p.status }}
                </span>
              </div>
              <p class="text-xs text-[#6B6D74] mt-1">
                Medio: <strong>{{ p.paymentMethod }}</strong> · Vencimiento: {{ p.dueDate }} · Fecha de pago: {{ p.paymentDate || 'Pendiente' }}
              </p>
              <p v-if="p.transactionReference" class="text-[11px] text-[#8E9097] font-mono mt-0.5">
                Ref: {{ p.transactionReference }}
              </p>
            </div>

            <div class="flex items-center gap-2">
              <button
                @click="handleSendPaymentLink"
                class="px-3 py-1.5 text-xs border border-[#DCD8CF] rounded font-semibold text-[#191A1D] hover:bg-white cursor-pointer"
              >
                Enviar enlace
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 4: EVALUACIONES -->
    <div v-else-if="currentTab === 'evaluaciones'" class="space-y-4">
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-[16px] font-bold text-[#191A1D]">Evaluaciones Físicas y Seguimiento</h3>
            <p class="text-xs text-[#6B6D74]">Resultados de test, mediciones y observaciones periódicas.</p>
          </div>
          <button
            @click="showEvalModal = true"
            class="px-4 py-2 bg-[#CEA549] text-[#0B0B0D] text-xs font-bold rounded hover:bg-[#BA923B] transition-colors cursor-pointer shadow-xs"
          >
            + Nueva evaluación
          </button>
        </div>

        <div v-if="studentEvaluations.length === 0" class="text-center py-10 text-sm text-[#8E9097]">
          No hay evaluaciones registradas para este alumno.
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="ev in studentEvaluations"
            :key="ev.id"
            class="p-4 border border-[#E2DFD8] rounded-lg hover:border-[#CEA549]/60 transition-colors bg-white space-y-3"
          >
            <div class="flex items-center justify-between border-b pb-2 border-[#ECE8E1]">
              <div>
                <span class="font-bold text-sm text-[#191A1D]">{{ ev.type }}</span>
                <span class="text-xs text-[#6B6D74] ml-2">· {{ ev.date }}</span>
              </div>
              <span class="px-2 py-0.5 text-xs rounded bg-[#FAF9F6] border border-[#E2DFD8] text-[#191A1D] font-semibold">
                {{ ev.status }}
              </span>
            </div>

            <div v-if="ev.metrics && ev.metrics.length > 0" class="grid grid-cols-2 sm:grid-cols-4 gap-2">
              <div
                v-for="(m, mi) in ev.metrics"
                :key="mi"
                class="p-2 bg-[#FAF9F6] rounded border border-[#ECE8E1] text-xs"
              >
                <p class="text-[#8E9097] text-[10px] uppercase font-bold">{{ m.name }}</p>
                <p class="text-[#191A1D] font-bold mt-0.5">{{ m.value }} {{ m.unit || '' }}</p>
              </div>
            </div>

            <p class="text-xs text-[#6B6D74] bg-[#FAF9F6] p-2.5 rounded border border-[#ECE8E1] italic">
              "{{ ev.observations || 'Sin observaciones registradas' }}"
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 5: CONTRATO -->
    <div v-else-if="currentTab === 'contrato'" class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-4">
      <div class="flex items-center justify-between border-b pb-3 border-[#ECE8E1]">
        <div>
          <h3 class="text-[16px] font-bold text-[#191A1D]">Contrato y Términos y Condiciones</h3>
          <p class="text-xs text-[#6B6D74]">Estado de consentimiento informado y aceptación de servicios.</p>
        </div>
        <button
          @click="toggleContractStatus"
          class="px-3.5 py-1.5 rounded text-xs font-bold transition-colors cursor-pointer"
          :class="studentContract?.status === 'Aceptado' ? 'bg-[#FEF7E0] text-[#B06000] hover:bg-[#FEEFC3]' : 'bg-[#E6F4EA] text-[#137333] hover:bg-[#CEEAD6]'"
        >
          {{ studentContract?.status === 'Aceptado' ? 'Marcar como pendiente' : 'Marcar como aceptado (Demo)' }}
        </button>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 p-4 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] text-xs">
        <div>
          <span class="text-[#6B6D74]">Estado de Aceptación:</span>
          <p class="font-bold text-[#191A1D] text-sm mt-0.5">{{ studentContract?.status || 'Pendiente' }}</p>
        </div>
        <div>
          <span class="text-[#6B6D74]">Versión de Términos:</span>
          <p class="font-bold text-[#191A1D] text-sm mt-0.5">{{ studentContract?.termsVersion || 'v1.0 (2026)' }}</p>
        </div>
        <div>
          <span class="text-[#6B6D74]">Copia de Respaldo:</span>
          <p class="font-bold text-[#191A1D] text-sm mt-0.5">
            {{ studentContract?.signedCopyAvailable ? 'Disponible (Digital)' : 'Pendiente' }}
          </p>
        </div>
      </div>

      <div class="p-4 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] text-xs text-[#191A1D] font-mono leading-relaxed max-h-60 overflow-y-auto whitespace-pre-wrap">
CONTRATO DE ADHESIÓN Y CONSENTIMIENTO DE ENTRENAMIENTO — INVERFIT
Alumno: {{ student.fullName }}
Plan: {{ student.planName }} ({{ student.frequency }})
Horario fijo acordado: {{ student.fixedSchedule }}

El alumno declara encontrarse en condiciones físicas para participar en sesiones de fuerza y acondicionamiento. Acepta que los cambios de horario dependen de la disponibilidad del coach dentro de la misma semana y que las inasistencias sin aviso previo no aplazan la mensualidad.
      </div>

      <div class="p-3 bg-[#FEF7E0] border border-[#F0D878] rounded text-xs text-[#8C6B1B] flex items-start gap-2">
        <AlertCircle class="w-4 h-4 flex-shrink-0 mt-0.5 text-[#B06000]" />
        <span>
          Aviso: Este mecanismo de aceptación es una demostración técnica para validar la visualización del estado documental. No constituye firma electrónica legal.
        </span>
      </div>
    </div>

    <!-- Modals -->
    <SimulatedActionModal
      :open="showLinkModal"
      mode="link"
      :student-name="student.fullName"
      :link-url="generatedLink"
      @close="showLinkModal = false"
    />

    <RegisterPaymentModal
      :open="showPaymentModal"
      :initial-student-id="student.id"
      @close="showPaymentModal = false"
    />

    <CreateEvaluationModal
      :open="showEvalModal"
      :initial-student-id="student.id"
      @close="showEvalModal = false"
    />
  </div>
</template>
