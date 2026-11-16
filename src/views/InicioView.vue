<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { Calendar, Users, CreditCard, ChevronRight } from 'lucide-vue-next'
import RegisterStudentModal from '@/components/modals/RegisterStudentModal.vue'
import BlockDetailsModal from '@/components/modals/BlockDetailsModal.vue'
import AssignScheduleModal from '@/components/modals/AssignScheduleModal.vue'
import type { ScheduleBlock } from '@/types'

const router = useRouter()
const gymStore = useGymStore()

onMounted(() => {
  if (gymStore.isScheduleApiMode()) {
    gymStore.loadSchedule().catch(() => {})
  }
  if (gymStore.isPaymentsApiMode()) {
    gymStore.loadPayments().catch(() => {})
  }
})

const showRegisterModal = ref(false)
const showBlockModal = ref(false)
const showAssignModal = ref(false)
const activeBlockId = ref<string | null>(null)
const activeBlock = computed(() => {
  if (!activeBlockId.value) return null
  return gymStore.schedule.find(b => b.id === activeBlockId.value) || null
})

// Calculate today blocks (or Lun if Dom)
const todayDayName = computed(() => {
  const map = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']
  const d = new Date().getDay()
  return map[d] === 'Dom' ? 'Lun' : map[d]
})

const todayBlocks = computed(() => {
  const list = gymStore.schedule.filter(b => b.day === todayDayName.value)
  if (list.length > 0) return list
  // fallback to Lun blocks for rich preview
  return gymStore.schedule.filter(b => b.day === 'Lun')
})

function openBlockDetail(block: ScheduleBlock) {
  activeBlockId.value = block.id
  showBlockModal.value = true
}

function handleAddStudentToBlock(block: ScheduleBlock) {
  showBlockModal.value = false
  activeBlockId.value = block.id
  showAssignModal.value = true
}

function getAssignedCount(block: ScheduleBlock) {
  // Read live from gymStore so it reacts immediately
  const live = gymStore.schedule.find(b => b.id === block.id)
  return live ? live.studentIds.length : block.studentIds.length
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Page Header matching wireframe -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Inicio</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Resumen general de la actividad del gimnasio.</p>
      </div>

      <button
        @click="showRegisterModal = true"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        + Registrar alumno
      </button>
    </div>

    <!-- 4 KPI Summary Cards matching wireframe -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Alumnos activos</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.activeStudentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Total de alumnos vigentes</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Pagos pendientes</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingPaymentsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Membresías por revisar</span>
      </div>

      <!-- Card 3 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Sesiones de hoy</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ todayBlocks.length }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Bloques programados ({{ todayDayName }})</span>
      </div>

      <!-- Card 4 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Contratos pendientes</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingContractsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Aceptaciones por completar</span>
      </div>
    </div>

    <!-- Main Content Grid: Agenda de hoy (Left) & Pendientes de atención (Right) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left Column: Agenda de hoy (width ~60% -> 7 cols) -->
      <div class="lg:col-span-7 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-[17px] font-bold text-[#191A1D]">Agenda de hoy</h2>
            <p class="text-[13px] text-[#6B6D74] mt-0.5">Próximos bloques de entrenamiento · {{ todayDayName }}</p>
          </div>
          <button
            @click="router.push('/agenda')"
            class="text-[13px] text-[#8C6B1B] hover:text-[#614910] font-bold flex items-center gap-1 hover:underline cursor-pointer"
          >
            Ver completa →
          </button>
        </div>

        <!-- Blocks list matching wireframe 01 -->
        <div class="space-y-3">
          <div
            v-for="block in todayBlocks"
            :key="block.id"
            class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] flex items-center justify-between transition-all group bg-white cursor-pointer"
            @click="openBlockDetail(block)"
          >
            <div class="flex items-center gap-4">
              <!-- Time Badge -->
              <div class="w-[58px] text-center py-1.5 px-2 bg-[#191A1D] rounded-[6px] border border-[#CEA549]/40 text-[12px] font-bold text-[#CEA549]">
                {{ block.startTime }}
              </div>

              <div>
                <p class="text-[14px] font-semibold text-[#191A1D] group-hover:text-black">
                  Bloque de entrenamiento
                </p>
                <p class="text-[12px] text-[#8E9097] flex items-center gap-1.5 mt-0.5">
                  <span>Alumnos asignados: {{ getAssignedCount(block) }}</span>
                  <span class="text-[#8C6B1B] font-medium group-hover:underline">· ver detalle</span>
                </p>
              </div>
            </div>

            <span class="text-xs text-[#8E9097] group-hover:text-[#CEA549] transition-colors">
              <ChevronRight class="w-4 h-4" />
            </span>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] flex items-center justify-between text-[11px] text-[#8E9097]">
          <span>Sesiones de duración flexible · El entrenador administra libremente sus tiempos</span>
          <button
            @click="showAssignModal = true"
            class="text-[#8C6B1B] font-bold hover:underline cursor-pointer"
          >
            + Asignar a bloque
          </button>
        </div>
      </div>

      <!-- Right Column: Pendientes de atención (width ~40% -> 5 cols) -->
      <div class="lg:col-span-5 bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs flex flex-col justify-between">
        <div>
          <h2 class="text-[17px] font-bold text-[#191A1D] mb-4">Pendientes de atención</h2>

          <div class="space-y-4">
            <!-- Item 1: Membresías por vencer -->
            <div
              @click="router.push('/membresias-pagos')"
              class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] cursor-pointer transition-all"
            >
              <div class="flex items-center justify-between">
                <p class="text-[14px] font-semibold text-[#191A1D]">Membresías por vencer</p>
                <span class="text-[11px] font-semibold px-2 py-0.5 rounded bg-[#FEF7E0] text-[#B06000]">
                  {{ gymStore.expiringMembershipsCount }} pendientes
                </span>
              </div>
              <p class="text-[12px] text-[#6B6D74] mt-1">Revisar fechas y enviar enlace de pago</p>
            </div>

            <!-- Item 2: Contratos sin aceptar -->
            <div
              @click="router.push('/contratos')"
              class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] cursor-pointer transition-all"
            >
              <div class="flex items-center justify-between">
                <p class="text-[14px] font-semibold text-[#191A1D]">Contratos sin aceptar</p>
                <span class="text-[11px] font-semibold px-2 py-0.5 rounded bg-[#FEF7E0] text-[#B06000]">
                  {{ gymStore.pendingContractsCount }} sin aceptar
                </span>
              </div>
              <p class="text-[12px] text-[#6B6D74] mt-1">Solicitar aceptación al alumno</p>
            </div>

            <!-- Item 3: Evaluaciones por registrar -->
            <div
              @click="router.push('/evaluaciones')"
              class="p-3.5 rounded-[8px] border border-[#E2DFD8] hover:border-[#CEA549]/60 hover:bg-[#FAF9F6] cursor-pointer transition-all"
            >
              <div class="flex items-center justify-between">
                <p class="text-[14px] font-semibold text-[#191A1D]">Evaluaciones por registrar</p>
                <span class="text-[11px] font-semibold px-2 py-0.5 rounded bg-[#E8F0FE] text-[#1967D2]">
                  {{ gymStore.pendingEvaluationsCount }} pendientes
                </span>
              </div>
              <p class="text-[12px] text-[#6B6D74] mt-1">Revisar seguimiento de alumnos</p>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-[#ECE8E1] text-[11px] text-[#8E9097]">
          Alertas calculadas automáticamente desde los datos sincronizados.
        </div>
      </div>
    </div>

    <!-- Bottom Section: Accesos rápidos matching wireframe 01 -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
      <h3 class="text-[14px] font-bold text-[#191A1D] mb-3">Accesos rápidos</h3>

      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex flex-wrap items-center gap-3">
          <button
            @click="router.push('/agenda')"
            class="h-[38px] px-4 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549]/60 transition-colors flex items-center gap-2 cursor-pointer"
          >
            <Calendar class="w-4 h-4 text-[#8C6B1B]" />
            Ver agenda
          </button>

          <button
            @click="router.push('/alumnos')"
            class="h-[38px] px-4 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549]/60 transition-colors flex items-center gap-2 cursor-pointer"
          >
            <Users class="w-4 h-4 text-[#8C6B1B]" />
            Buscar alumno
          </button>

          <button
            @click="router.push('/membresias-pagos')"
            class="h-[38px] px-4 bg-white border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-[#FAF9F6] hover:border-[#CEA549]/60 transition-colors flex items-center gap-2 cursor-pointer"
          >
            <CreditCard class="w-4 h-4 text-[#8C6B1B]" />
            Gestionar pagos
          </button>
        </div>

        <span class="text-[12px] text-[#8E9097] italic">
          Wireframe con datos ficticios para validación.
        </span>
      </div>
    </div>

    <!-- Modals -->
    <RegisterStudentModal
      :open="showRegisterModal"
      @close="showRegisterModal = false"
    />

    <BlockDetailsModal
      :open="showBlockModal"
      :block="activeBlock"
      @close="showBlockModal = false"
      @deleted="gymStore.refreshScheduleState"
      @updated="gymStore.refreshScheduleState"
      @add-student="handleAddStudentToBlock"
    />

    <AssignScheduleModal
      :open="showAssignModal"
      :block-id="activeBlockId || undefined"
      :initial-day="activeBlock?.day"
      :initial-time="activeBlock?.startTime"
      :initial-end-time="activeBlock?.endTime"
      @close="showAssignModal = false"
      @assigned="gymStore.refreshScheduleState"
    />
  </div>
</template>
