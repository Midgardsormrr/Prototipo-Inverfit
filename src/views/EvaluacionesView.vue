<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { onMounted } from 'vue'
import { useGymStore } from '@/stores/gym'
import { Trash2, AlertTriangle, Activity, List, Plus, ArrowLeft } from 'lucide-vue-next'
import PabloEvaluationWorkbook from '@/components/evaluations/PabloEvaluationWorkbook.vue'

const router = useRouter()
const route = useRoute()
const gymStore = useGymStore()

// View Mode: 'list' (general overview of all students) as default, 'workbook' (official evaluation workbook)
const viewMode = ref<'list' | 'workbook'>('list')
const selectedWorkbookStudentId = ref('')

const searchName = ref('')
const selectedStatus = ref('')
const selectedProtocol = ref('')

onMounted(() => {
  if (route.query.studentId) {
    selectedWorkbookStudentId.value = String(route.query.studentId)
    viewMode.value = 'workbook'
  } else if (route.query.view === 'workbook') {
    viewMode.value = 'workbook'
  }
})

// Unified list of all students combined with evaluation status
const allEvaluationRows = computed(() => {
  const rows: Array<{
    id: string
    studentId: string
    studentName: string
    frequency?: string
    date: string
    type: string
    status: 'Registrada' | 'Pendiente'
    protocol?: string
    clinicalAlert?: string
    observations?: string
    isRegistered: boolean
  }> = []

  const evaluatedStudentIds = new Set<string>()

  // 1. Add evaluations recorded in gymStore
  for (const ev of gymStore.evaluations) {
    evaluatedStudentIds.add(ev.studentId)
    const st = gymStore.students.find(s => s.id === ev.studentId)
    rows.push({
      id: ev.id,
      studentId: ev.studentId,
      studentName: ev.studentName,
      frequency: st?.frequency,
      date: ev.date,
      type: ev.type,
      status: ev.status === 'Registrada' ? 'Registrada' : 'Pendiente',
      protocol: ev.protocol,
      clinicalAlert: ev.clinicalAlert,
      observations: ev.observations,
      isRegistered: ev.status === 'Registrada'
    })
  }

  // 2. Add registered gym students who do NOT have an evaluation yet (e.g. newly created students)
  for (const st of gymStore.students) {
    if (!evaluatedStudentIds.has(st.id)) {
      rows.push({
        id: `pending-${st.id}`,
        studentId: st.id,
        studentName: st.fullName,
        frequency: st.frequency,
        date: 'Por registrar',
        type: 'Evaluación inicial',
        status: 'Pendiente',
        protocol: st.onboarding?.age && Number(st.onboarding.age) >= 60 ? 'adulto_mayor' : 'clinico_deportivo',
        clinicalAlert: st.healthConditions && st.healthConditions !== 'Ninguna' ? st.healthConditions : undefined,
        observations: 'Pendiente de evaluación física y funcional.',
        isRegistered: false
      })
    }
  }

  return rows
})

const filteredEvaluations = computed(() => {
  return allEvaluationRows.value.filter(ev => {
    const q = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (q && !ev.studentName.toLowerCase().includes(q)) return false
    if (selectedStatus.value && ev.status !== selectedStatus.value) return false
    if (selectedProtocol.value && ev.protocol !== selectedProtocol.value) return false
    return true
  })
})

const registeredEvaluationsCount = computed(() => {
  return allEvaluationRows.value.filter(r => r.isRegistered).length
})

const pendingEvaluationsCount = computed(() => {
  return allEvaluationRows.value.filter(r => !r.isRegistered).length
})

const lastUpdateDate = computed(() => {
  const list = allEvaluationRows.value.filter(e => e.isRegistered && e.date !== 'Por registrar')
  if (list.length === 0) return '—'
  return list[0].date
})

function openWorkbookForStudent(studentId: string) {
  selectedWorkbookStudentId.value = studentId
  viewMode.value = 'workbook'
}

function startNewEvaluation(studentId?: string) {
  if (studentId) selectedWorkbookStudentId.value = studentId
  viewMode.value = 'workbook'
}

function goToStudent(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}

async function handleDelete(id: string) {
  if (confirm('¿Eliminar este registro de evaluación?')) {
    await gymStore.deleteEvaluation(id)
  }
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header with View Switcher -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Evaluaciones</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Protocolos clínicos, antropometría y seguimiento científico del alumno.</p>
      </div>

      <div class="flex items-center gap-2">
        <!-- View Toggle Pills -->
        <div class="bg-white p-1 rounded-lg border border-[#E2DFD8] flex items-center gap-1 shadow-2xs">
          <button
            @click="viewMode = 'list'"
            class="px-3 py-1.5 rounded-md text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
            :class="viewMode === 'list'
              ? 'bg-[#191A1D] text-[#CEA549] shadow-xs'
              : 'text-[#6B6D74] hover:text-[#191A1D]'"
          >
            <List class="w-3.5 h-3.5" />
            <span>Vista General (Alumnos)</span>
          </button>

          <button
            @click="viewMode = 'workbook'"
            class="px-3 py-1.5 rounded-md text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
            :class="viewMode === 'workbook'
              ? 'bg-[#191A1D] text-[#CEA549] shadow-xs'
              : 'text-[#6B6D74] hover:text-[#191A1D]'"
          >
            <Activity class="w-3.5 h-3.5" />
            <span>Ficha de Evaluación</span>
          </button>
        </div>

        <button
          v-if="viewMode === 'list'"
          @click="startNewEvaluation()"
          class="h-[38px] px-4 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[13px] font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs hover:shadow"
        >
          <Plus class="w-4 h-4" />
          <span>Evaluar Alumno</span>
        </button>
      </div>
    </div>

    <!-- VIEW 1: OFFICIAL EVALUATION WORKBOOK -->
    <div v-if="viewMode === 'workbook'" class="space-y-4">
      <div class="flex items-center justify-between pb-1">
        <button
          @click="viewMode = 'list'"
          class="inline-flex items-center gap-2 text-xs font-bold text-[#191A1D] hover:text-black px-3.5 py-2 rounded-md bg-white border border-[#E2DFD8] hover:border-[#CEA549] transition-all cursor-pointer shadow-2xs group"
        >
          <ArrowLeft class="w-4 h-4 text-[#CEA549] group-hover:-translate-x-0.5 transition-transform" />
          <span>← Volver al listado general de alumnos</span>
        </button>

        <span class="text-xs text-[#8E9097] hidden sm:inline">
          Modo ficha individual interactiva · <strong>Protocolo Oficial INVERFIT</strong>
        </span>
      </div>

      <PabloEvaluationWorkbook :student-id="selectedWorkbookStudentId" />
    </div>

    <!-- VIEW 2: GENERAL LIST & KPIS -->
    <div v-else class="space-y-6">
      <!-- 3 KPI Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <!-- Card 1 -->
        <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
          <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
          <span class="text-[13px] font-semibold text-[#6B6D74]">Evaluaciones registradas</span>
          <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
            {{ registeredEvaluationsCount }}
          </span>
          <span class="text-[12px] text-[#8E9097]">Historial del gimnasio</span>
        </div>

        <!-- Card 2 -->
        <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
          <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
          <span class="text-[13px] font-semibold text-[#6B6D74]">Alumnos por evaluar</span>
          <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
            {{ pendingEvaluationsCount }}
          </span>
          <span class="text-[12px] text-[#8E9097]">Seguimiento pendiente</span>
        </div>

        <!-- Card 3 -->
        <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
          <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
          <span class="text-[13px] font-semibold text-[#6B6D74]">Última actualización</span>
          <span class="text-[22px] font-bold text-[#191A1D] leading-tight truncate my-1">
            {{ lastUpdateDate }}
          </span>
          <span class="text-[12px] text-[#8E9097]">Evaluaciones oficiales</span>
        </div>
      </div>

      <!-- Toolbar -->
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
            <option value="">Todos los estados  ▾</option>
            <option value="Registrada">Evaluación Registrada</option>
            <option value="Pendiente">Pendiente de evaluar</option>
          </select>

          <select
            v-model="selectedProtocol"
            class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
          >
            <option value="">Todos los protocolos  ▾</option>
            <option value="clinico_deportivo">Clínico-Deportivo (Encoder ADR)</option>
            <option value="adulto_mayor">Adulto Mayor (Senior SPPB)</option>
          </select>
        </div>
      </div>

      <!-- Table -->
      <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse min-w-[750px]">
            <thead>
              <tr class="border-b border-[#E2DFD8] bg-[#FAF9F6]">
                <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                  Alumno
                </th>
                <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                  Fecha
                </th>
                <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                  Protocolo
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
              <tr v-if="filteredEvaluations.length === 0">
                <td colspan="5" class="py-12 text-center text-sm text-[#8E9097]">
                  No se encontraron evaluaciones registradas.
                </td>
              </tr>

              <tr
                v-for="ev in filteredEvaluations"
                :key="ev.id"
                class="hover:bg-[#FAF9F6] transition-colors group cursor-pointer"
                @click="openWorkbookForStudent(ev.studentId)"
              >
                <!-- Alumno -->
                <td class="py-4 px-6 text-[13px] font-medium text-[#191A1D]">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[12px] font-bold">
                      {{ ev.studentName.charAt(0) }}
                    </div>
                    <div>
                      <div class="flex items-center gap-1.5 flex-wrap">
                        <span class="group-hover:text-black font-bold text-[#191A1D]">{{ ev.studentName }}</span>
                        <span
                          v-if="ev.clinicalAlert"
                          class="inline-flex items-center gap-0.5 px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#FEF7E0] text-[#B06000] border border-[#F0D878]"
                          :title="ev.clinicalAlert"
                        >
                          <AlertTriangle class="w-3 h-3 text-[#B06000]" />
                          Alerta Médica
                        </span>
                      </div>
                      <p class="text-[11px] text-[#8E9097] leading-none mt-0.5 truncate max-w-[240px]">
                        {{ ev.observations || 'Evaluación física' }}
                      </p>
                    </div>
                  </div>
                </td>

                <!-- Fecha -->
                <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                  {{ ev.date }}
                </td>

                <!-- Protocolo -->
                <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                  <span
                    v-if="ev.protocol === 'adulto_mayor'"
                    class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-bold bg-[#FAF9F6] text-[#8C6B1B] border border-[#E2DFD8]"
                  >
                    Senior SPPB (Adulto Mayor)
                  </span>
                  <span
                    v-else
                    class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-bold bg-[#FAF9F6] text-[#CEA549] border border-[#E2DFD8]"
                  >
                    Clínico-Deportivo (Encoder ADR)
                  </span>
                </td>

                <!-- Estado -->
                <td class="py-4 px-6 text-[13px]">
                  <span
                    class="px-2.5 py-1 rounded-[4px] text-[12px] font-semibold inline-block"
                    :class="ev.status === 'Registrada' ? 'bg-[#E6F4EA] text-[#137333]' : 'bg-[#FEF7E0] text-[#B06000]'"
                  >
                    {{ ev.status }}
                  </span>
                </td>

                <!-- Acción -->
                <td class="py-4 px-6 text-[13px] font-semibold text-right whitespace-nowrap" @click.stop>
                  <div class="flex items-center justify-end gap-3">
                    <button
                      @click.stop="openWorkbookForStudent(ev.studentId)"
                      class="text-[#CEA549] hover:text-[#BA923B] hover:underline cursor-pointer font-bold flex items-center gap-1"
                    >
                      <component :is="ev.isRegistered ? Activity : Plus" class="w-3.5 h-3.5" />
                      <span>{{ ev.isRegistered ? 'Ver Evaluación →' : 'Evaluar →' }}</span>
                    </button>
                    <button
                      @click="goToStudent(ev.studentId)"
                      class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer font-bold"
                    >
                      Ficha Alumno
                    </button>
                    <button
                      @click="handleDelete(ev.id)"
                      title="Eliminar registro"
                      class="p-1 text-[#8E9097] hover:text-rose-600 transition-colors cursor-pointer"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

  </div>
</template>
