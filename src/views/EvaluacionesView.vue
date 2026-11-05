<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { Trash2 } from 'lucide-vue-next'
import CreateEvaluationModal from '@/components/modals/CreateEvaluationModal.vue'

const router = useRouter()
const gymStore = useGymStore()

const showCreateModal = ref(false)
const targetStudentId = ref('')
const searchName = ref('')
const selectedType = ref('')

const filteredEvaluations = computed(() => {
  return gymStore.evaluations.filter(ev => {
    const q = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (q && !ev.studentName.toLowerCase().includes(q)) return false
    if (selectedType.value && ev.type !== selectedType.value) return false
    return true
  })
})

const lastUpdateDate = computed(() => {
  const list = gymStore.evaluations
  if (list.length === 0) return '—'
  return list[0].date
})

function openRegisterForStudent(studentId: string) {
  targetStudentId.value = studentId
  showCreateModal.value = true
}

function openCreateModal() {
  targetStudentId.value = ''
  showCreateModal.value = true
}

function goToStudent(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}

async function handleDelete(id: string) {
  if (confirm('¿Eliminar este registro de evaluación de demostración?')) {
    await gymStore.deleteEvaluation(id)
  }
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header matching wireframe 06 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Evaluaciones</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Registra y consulta resultados de seguimiento de los alumnos.</p>
      </div>

      <button
        @click="openCreateModal"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        + Nueva evaluación
      </button>
    </div>

    <!-- 3 KPI Summary Cards matching wireframe 06: [x=278, y=193], [x=623, y=193], [x=968, y=193] -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <!-- Card 1 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Evaluaciones registradas</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.registeredEvaluationsCount }}
        </span>
        <span class="text-[12px] text-[#8E9097]">Historial del gimnasio</span>
      </div>

      <!-- Card 2 -->
      <div class="relative bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col justify-between shadow-2xs hover:border-[#CEA549]/60 transition-all group overflow-hidden">
        <div class="absolute top-0 left-0 right-0 h-[3px] bg-transparent group-hover:bg-[#CEA549] transition-colors"></div>
        <span class="text-[13px] font-semibold text-[#6B6D74]">Alumnos por evaluar</span>
        <span class="text-[30px] font-bold text-[#191A1D] leading-tight my-1">
          {{ gymStore.pendingEvaluationsCount }}
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
        <span class="text-[12px] text-[#8E9097]">Fecha por confirmar</span>
      </div>
    </div>

    <!-- Toolbar: Search + Tipo de evaluación + Período matching wireframe -->
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
          v-model="selectedType"
          class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Tipo de evaluación  ▾</option>
          <option value="Evaluación inicial">Evaluación inicial</option>
          <option value="Seguimiento">Seguimiento</option>
        </select>

        <select
          class="h-[40px] px-3 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] hover:bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Período  ▾</option>
          <option value="mes">Mes actual</option>
          <option value="trimestre">Últimos 3 meses</option>
          <option value="todos">Todos los períodos</option>
        </select>
      </div>
    </div>

    <!-- Table matching wireframe 06: Alumno | Fecha | Tipo | Estado | Acción -->
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
                Tipo
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
              @click="goToStudent(ev.studentId)"
            >
              <!-- Alumno -->
              <td class="py-4 px-6 text-[13px] font-medium text-[#191A1D]">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[12px] font-bold">
                    {{ ev.studentName.charAt(0) }}
                  </div>
                  <div>
                    <span class="group-hover:text-black font-bold text-[#191A1D]">{{ ev.studentName }}</span>
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

              <!-- Tipo -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ ev.type }}
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
                    @click="openRegisterForStudent(ev.studentId)"
                    class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer font-bold"
                  >
                    Registrar →
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

    <!-- Informational Card matching wireframe 06: [x=280, y=718] -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs space-y-1">
      <h3 class="text-[14px] font-bold text-[#191A1D]">Al registrar una evaluación</h3>
      <p class="text-[12px] text-[#6B6D74]">
        Se debe definir qué métricas, resultados y observaciones se almacenarán en la ficha del alumno. La plataforma ofrece campos configurables y flexibles.
      </p>
    </div>

    <!-- Modals -->
    <CreateEvaluationModal
      :open="showCreateModal"
      :initial-student-id="targetStudentId"
      @close="showCreateModal = false"
    />
  </div>
</template>
