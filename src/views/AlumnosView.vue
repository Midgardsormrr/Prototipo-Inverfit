<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGymStore } from '@/stores/gym'
import { Trash2 } from 'lucide-vue-next'
import RegisterStudentModal from '@/components/modals/RegisterStudentModal.vue'

const router = useRouter()
const gymStore = useGymStore()

const showRegisterModal = ref(false)
const searchName = ref('')
const selectedStatus = ref('')
const selectedFrequency = ref('')

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
  if (confirm(`¿Eliminar al alumno de demostración "${name}"? Esta acción removerá sus pagos, evaluaciones y asignaciones de bloques.`)) {
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

      <button
        @click="showRegisterModal = true"
        class="h-[40px] px-5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-[6px] text-[14px] font-bold transition-all flex items-center justify-center gap-1.5 self-start sm:self-auto cursor-pointer shadow-xs hover:shadow"
      >
        + Nuevo alumno
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
          <option value="">Todos los planes  ▾</option>
          <option value="1 vez/sem.">1 vez/sem.</option>
          <option value="2 veces/sem.">2 veces/sem.</option>
          <option value="3 veces/sem.">3 veces/sem.</option>
          <option value="4 veces/sem.">4 veces/sem.</option>
        </select>
      </div>
    </div>

    <!-- Table matching wireframe 03: Alumno | Plan | Frecuencia | Próximo pago | Estado | Acción -->
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
                Frecuencia
              </th>
              <th class="py-3 px-6 text-[12px] font-bold text-[#6B6D74] uppercase tracking-wider">
                Próximo pago
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
            <tr v-if="filteredStudents.length === 0">
              <td colspan="6" class="py-12 text-center text-sm text-[#8E9097]">
                No se encontraron alumnos que coincidan con la búsqueda o filtro.
              </td>
            </tr>

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
                    <span class="group-hover:text-black font-bold text-[#191A1D]">{{ student.fullName }}</span>
                    <p class="text-[11px] text-[#8E9097] leading-none mt-0.5">{{ student.fixedSchedule }}</p>
                  </div>
                </div>
              </td>

              <!-- Plan -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ student.planName }}
              </td>

              <!-- Frecuencia -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ student.frequency }}
              </td>

              <!-- Próximo pago -->
              <td class="py-4 px-6 text-[13px] text-[#4A4C52]">
                {{ student.membershipDueDate }}
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
                    Ver →
                  </button>
                  <button
                    @click="handleDelete(student.id, student.fullName)"
                    title="Eliminar de demostración"
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

    <!-- Bottom Informational card matching wireframe 03 -->
    <div class="bg-white rounded-[10px] border border-[#E2DFD8] p-6 shadow-2xs">
      <h3 class="text-[15px] font-bold text-[#191A1D]">Ficha del alumno</h3>
      <p class="text-[13px] text-[#6B6D74] mt-1">
        Desde el perfil se podrán consultar datos, horario fijo, pagos, evaluaciones y contrato.
      </p>
      <p class="text-[12px] text-[#8E9097] mt-3">
        Datos de ejemplo. No utilizar información real de alumnos en el prototipo.
      </p>
    </div>

    <!-- Modals -->
    <RegisterStudentModal
      :open="showRegisterModal"
      @close="showRegisterModal = false"
    />
  </div>
</template>
