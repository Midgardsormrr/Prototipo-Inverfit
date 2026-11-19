<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { Trash2, RefreshCw, AlertCircle, Share2, Check, Activity, Clock, Phone } from 'lucide-vue-next'
import RegisterStudentModal from '@/components/modals/RegisterStudentModal.vue'
import type { Student } from '@/types'

const router = useRouter()
const gymStore = useGymStore()

const showRegisterModal = ref(false)
const searchName = ref('')
const selectedStatus = ref('')
const selectedFrequency = ref('')
const copiedLink = ref(false)

function getStudentSchedule(student: Student) {
  const blocks = gymStore.schedule.filter(b => b.studentIds.includes(student.id))
  if (blocks.length > 0) {
    const dayOrder = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
    return blocks
      .slice()
      .sort((a, b) => dayOrder.indexOf(a.day) - dayOrder.indexOf(b.day))
      .map(b => `${b.day} ${b.startTime}`)
      .join(', ')
  }
  return student.fixedSchedule || 'Por coordinar'
}

function getStudentEvaluation(studentId: string) {
  return gymStore.evaluations.find(e => e.studentId === studentId && e.status === 'Registrada')
}

async function copyPublicLink() {
  const url = `${window.location.origin}/registro-alumno`
  try {
    await navigator.clipboard.writeText(url)
    copiedLink.value = true
    gymStore.showToast('Enlace copiado al portapapeles. ¡Pégalo en WhatsApp para que el alumno complete su ficha y firma!', 'success', 'Enlace copiado')
    setTimeout(() => {
      copiedLink.value = false
    }, 3000)
  } catch {
    gymStore.showToast(`Enlace: ${url}`, 'info')
  }
}

onMounted(async () => {
  await gymStore.loadStudents()
})

const filteredStudents = computed(() => {
  return gymStore.students.filter(student => {
    // Global or local search
    const query = (searchName.value || gymStore.globalSearch).toLowerCase().trim()
    if (query && !student.fullName.toLowerCase().includes(query)) {
      return false
    }

    // Status filter
    if (selectedStatus.value && student.membershipStatus !== selectedStatus.value) {
      return false
    }

    // Frequency filter
    if (selectedFrequency.value && student.frequency !== selectedFrequency.value) {
      return false
    }

    return true
  })
})

function goToProfile(studentId: string) {
  router.push(`/alumnos/${studentId}`)
}

async function handleDelete(studentId: string, name: string) {
  const isApi = gymStore.dataSource === 'api'
  const msg = isApi
    ? `¿Archivar al alumno "${name}"? El alumno será marcado como archivado en MongoDB y sus antecedentes se conservarán intactos.`
    : `¿Eliminar al alumno de demostración "${name}"? Esta acción removerá sus pagos, evaluaciones y asignaciones de bloques en localStorage.`

  if (confirm(msg)) {
    await gymStore.deleteStudent(studentId)
  }
}
</script>

<template>
  <div class="space-y-6 max-w-[1100px] mx-auto pb-10">
    <!-- Header matching wireframe 03 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-bold text-[#191A1D] tracking-tight">Alumnos</h1>
        <p class="text-[14px] text-[#6B6D74] mt-0.5">Consulta alumnos, planes y estado de sus membresías.</p>
      </div>

      <div class="flex items-center gap-2 self-start sm:self-auto flex-wrap">
        <button
          @click="copyPublicLink"
          title="Copiar enlace de ficha oficial y consentimiento para enviar por WhatsApp"
          class="h-[40px] px-3.5 bg-white hover:bg-[#FAF9F6] border border-[#DCD8CF] text-[#191A1D] rounded-[6px] text-[13px] font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-2xs hover:border-[#CEA549]"
        >
          <Check v-if="copiedLink" class="w-4 h-4 text-[#137333]" />
          <Share2 v-else class="w-4 h-4 text-[#8C6B1B]" />
          <span class="hidden sm:inline">{{ copiedLink ? '¡Enlace copiado!' : 'Copiar link auto-ingreso' }}</span>
          <span class="sm:hidden">{{ copiedLink ? '¡Copiado!' : 'Link WhatsApp' }}</span>
        </button>

        <button
          @click="gymStore.loadStudents()"
          title="Recargar lista de alumnos"
          :disabled="gymStore.loadingStudents"
          class="h-[40px] px-3 bg-white hover:bg-[#FAF9F6] border border-[#DCD8CF] text-[#4A4C52] rounded-[6px] text-[13px] font-medium transition-all flex items-center justify-center gap-1.5 cursor-pointer disabled:opacity-50"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': gymStore.loadingStudents }" />
          <span class="hidden sm:inline">Refrescar</span>
        </button>

        <button
          @click="showRegisterModal = true"
          class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs hover:shadow"
        >
          + Nuevo alumno
        </button>
      </div>
    </div>

    <!-- Error Alert if API failed -->
    <div
      v-if="gymStore.studentsError"
      class="bg-rose-50 border border-rose-200 rounded-[10px] p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-rose-800 shadow-2xs animate-in fade-in duration-200"
    >
      <div class="flex items-start gap-3">
        <AlertCircle class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" />
        <div>
          <p class="text-[14px] font-bold">No se pudieron cargar los alumnos desde el backend</p>
          <p class="text-[12px] text-rose-700 mt-0.5">{{ gymStore.studentsError }}</p>
        </div>
      </div>
      <button
        @click="gymStore.loadStudents()"
        class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-[6px] text-[13px] font-semibold transition-colors flex items-center justify-center gap-1.5 shrink-0 cursor-pointer self-start sm:self-auto"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': gymStore.loadingStudents }" />
        Reintentar conexión
      </button>
    </div>

    <!-- Filters Bar matching wireframe: Search + Estado filter + Plan filter -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-4 flex flex-col md:flex-row items-center justify-between gap-4 shadow-2xs">
      <div class="w-full md:w-[320px] relative">
        <input
          v-model="searchName"
          type="text"
          placeholder="⌕  Buscar por nombre..."
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

      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <select
          v-model="selectedStatus"
          class="h-[40px] px-3.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:bg-white focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Todos los estados  ▾</option>
          <option value="Activo">Activo</option>
          <option value="Por revisar">Por revisar</option>
          <option value="Inactivo">Inactivo</option>
        </select>

        <select
          v-model="selectedFrequency"
          class="h-[40px] px-3.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] font-semibold text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:bg-white focus:ring-1 focus:ring-[#CEA549]/30 transition-all cursor-pointer"
        >
          <option value="">Todas las frecuencias  ▾</option>
          <option value="1 vez/sem.">1 vez/sem.</option>
          <option value="2 veces/sem.">2 veces/sem.</option>
          <option value="3 veces/sem.">3 veces/sem.</option>
          <option value="4 veces/sem.">4 veces/sem.</option>
        </select>
      </div>
    </div>

    <!-- Table: Alumno | Plan (Frecuencia) | Horario Asignado | Evaluación | Estado | Acción -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] overflow-hidden shadow-2xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[750px]">
          <thead>
            <tr class="border-b border-[#E2DFD8] bg-[#FAF9F6]">
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Alumno
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Plan (Frecuencia)
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Horario Asignado
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Evaluación
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
            <!-- Loading State -->
            <tr v-if="gymStore.loadingStudents && gymStore.students.length === 0">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                <div class="flex items-center justify-center gap-2">
                  <RefreshCw class="w-4 h-4 animate-spin text-[#CEA549]" />
                  <span>Cargando alumnos desde {{ gymStore.dataSource === 'api' ? 'MongoDB...' : 'localStorage...' }}</span>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-else-if="filteredStudents.length === 0 && !gymStore.studentsError">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                <p class="font-medium text-[#4A4C52]">No se encontraron alumnos que coincidan con la búsqueda o filtro.</p>
                <p v-if="gymStore.dataSource === 'api'" class="text-xs text-[#8E9097] mt-1">
                  La colección en MongoDB está lista para recibir registros. Haz clic en "+ Nuevo alumno" para registrar el primero.
                </p>
              </td>
            </tr>

            <!-- Data Rows -->
            <tr
              v-for="student in filteredStudents"
              :key="student.id"
              class="hover:bg-[#FAF9F6] transition-colors group cursor-pointer"
              @click="goToProfile(student.id)"
            >
              <!-- Alumno -->
              <td class="py-4 px-6 text-[13px] font-medium text-[#191A1D]">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-[#191A1D] text-[#CEA549] border border-[#CEA549]/40 flex items-center justify-center text-[12px] font-bold">
                    {{ student.fullName.charAt(0) }}
                  </div>
                  <div>
                    <div class="flex items-center gap-1.5 flex-wrap">
                      <span class="group-hover:text-black font-bold text-[#191A1D]">{{ student.fullName }}</span>
                      <span
                        v-if="student.onboarding?.consentAccepted"
                        title="Ficha completada y consentimiento firmado digitalmente"
                        class="inline-flex items-center px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#E6F4EA] text-[#137333] border border-[#137333]/20"
                      >
                        ✓ Firmado
                      </span>
                      <span
                        v-if="student.healthConditions && student.healthConditions !== 'Sin patologías ni lesiones declaradas'"
                        title="Atención médica: reporta lesiones o enfermedades en su ficha"
                        class="inline-flex items-center px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#FEF7E0] text-[#B06000] border border-[#B06000]/20"
                      >
                        PAR-Q+ Alerta
                      </span>
                    </div>
                    <p v-if="student.phone" class="text-[11px] text-[#8E9097] flex items-center gap-1 mt-0.5">
                      <Phone class="w-3 h-3 text-[#CEA549]" />
                      {{ student.phone }}
                    </p>
                    <p v-else class="text-[11px] text-[#8E9097] mt-0.5">{{ student.email || 'Sin contacto directo' }}</p>
                  </div>
                </div>
              </td>

              <!-- Plan (Frecuencia) -->
              <td class="py-4 px-6 text-[13px] font-bold text-[#191A1D]">
                {{ student.frequency }}
              </td>

              <!-- Horario Asignado -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                <div class="flex items-center gap-1.5">
                  <Clock class="w-3.5 h-3.5 text-[#CEA549] shrink-0" />
                  <span :class="getStudentSchedule(student) === 'Por coordinar' ? 'text-[#8E9097] italic' : 'font-medium text-[#191A1D]'">
                    {{ getStudentSchedule(student) }}
                  </span>
                </div>
              </td>

              <!-- Evaluación -->
              <td class="py-4 px-6 text-[13px]" @click.stop>
                <router-link
                  v-if="getStudentEvaluation(student.id)"
                  :to="'/evaluaciones?studentId=' + student.id"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-[4px] text-[11px] font-bold bg-[#E6F4EA] text-[#137333] border border-[#137333]/20 hover:bg-[#CEEAD6] transition-colors"
                  title="Evaluación física registrada. Clic para ver evaluación"
                >
                  <Activity class="w-3.5 h-3.5 text-[#137333]" />
                  Evaluación Lista
                </router-link>
                <router-link
                  v-else
                  :to="'/evaluaciones?studentId=' + student.id"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-[4px] text-[11px] font-bold bg-[#FEF7E0] text-[#B06000] border border-[#F0D878] hover:bg-[#FEEFC3] transition-colors"
                  title="Falta evaluación física. Clic para evaluar alumno"
                >
                  Pendiente →
                </router-link>
              </td>

              <!-- Estado -->
              <td class="py-4 px-6 text-[13px]">
                <span
                  class="px-2.5 py-1 rounded-[4px] text-[12px] font-semibold inline-block"
                  :class="[
                    student.membershipStatus === 'Activo' && 'bg-[#E6F4EA] text-[#137333]',
                    student.membershipStatus === 'Por revisar' && 'bg-[#FEF7E0] text-[#B06000]',
                    student.membershipStatus === 'Inactivo' && 'bg-[#F1F3F4] text-[#5F6368]'
                  ]"
                >
                  {{ student.membershipStatus }}
                </span>
              </td>

              <!-- Acción -->
              <td class="py-4 px-6 text-[13px] font-semibold text-right whitespace-nowrap" @click.stop>
                <div class="flex items-center justify-end gap-3">
                  <button
                    @click="goToProfile(student.id)"
                    class="text-[#8C6B1B] hover:text-[#614910] hover:underline cursor-pointer flex items-center gap-1 font-bold"
                  >
                    Ficha →
                  </button>
                  <button
                    @click="handleDelete(student.id, student.fullName)"
                    :title="gymStore.dataSource === 'api' ? 'Archivar alumno' : 'Eliminar de demostración'"
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

    <!-- Modals -->
    <RegisterStudentModal
      :open="showRegisterModal"
      @close="showRegisterModal = false"
    />
  </div>
</template>
