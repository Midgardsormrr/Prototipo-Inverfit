<script setup lang="ts">
import { ref } from 'vue'
import { useGymStore } from '@/stores/gym'
import type { FrequencyType, StudentStatus } from '@/types'

defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', studentId: string): void
}>()

const gymStore = useGymStore()

const fullName = ref('')
const frequency = ref<FrequencyType>('2 veces/sem.')
const fixedSchedule = ref('')
const email = ref('')
const phone = ref('')
const membershipStatus = ref<StudentStatus>('Activo')
const emergencyContact = ref('')
const healthConditions = ref('')
const notes = ref('')

const isSubmitting = ref(false)
const errorMessage = ref('')

function resetForm() {
  fullName.value = ''
  frequency.value = '2 veces/sem.'
  fixedSchedule.value = ''
  email.value = ''
  phone.value = ''
  membershipStatus.value = 'Activo'
  emergencyContact.value = ''
  healthConditions.value = ''
  notes.value = ''
  errorMessage.value = ''
}

async function handleSubmit() {
  if (!fullName.value.trim()) {
    errorMessage.value = 'El nombre del alumno es obligatorio.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const student = await gymStore.createStudent({
      fullName: fullName.value.trim(),
      frequency: frequency.value,
      fixedSchedule: fixedSchedule.value.trim() || 'Días y horas por confirmar',
      email: email.value.trim(),
      phone: phone.value.trim(),
      membershipStatus: membershipStatus.value,
      emergencyContact: emergencyContact.value.trim(),
      healthConditions: healthConditions.value.trim(),
      notes: notes.value.trim()
    })

    resetForm()
    emit('created', student.id)
    emit('close')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al registrar alumno.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <h2 class="text-lg font-bold text-[#191A1D]">Registrar Alumno</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Ingresa los datos para su ficha y asignación de bloques.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4 max-h-[75vh] overflow-y-auto">
        <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-md text-sm">
          {{ errorMessage }}
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Nombre del Alumno *
          </label>
          <input
            v-model="fullName"
            type="text"
            placeholder="Ej. Alumno de ejemplo 07"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            required
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Frecuencia Semanal *
            </label>
            <select
              v-model="frequency"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            >
              <option value="1 vez/sem.">1 vez/sem.</option>
              <option value="2 veces/sem.">2 veces/sem.</option>
              <option value="3 veces/sem.">3 veces/sem.</option>
              <option value="4 veces/sem.">4 veces/sem.</option>
            </select>
            <p class="text-[11px] text-[#8E9097] mt-1">Planes operativos acordados.</p>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Estado Inicial
            </label>
            <select
              v-model="membershipStatus"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            >
              <option value="Activo">Activo</option>
              <option value="Por revisar">Por revisar</option>
              <option value="Inactivo">Inactivo</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Horario Fijo Sugerido
          </label>
          <input
            v-model="fixedSchedule"
            type="text"
            placeholder="Ej. Lun y Mié 08:00 (o por confirmar)"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Teléfono de Contacto
            </label>
            <input
              v-model="phone"
              type="text"
              placeholder="+56 9 1234 5678"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Correo Electrónico
            </label>
            <input
              v-model="email"
              type="email"
              placeholder="alumno@ejemplo.com"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Contacto de Emergencia
          </label>
          <input
            v-model="emergencyContact"
            type="text"
            placeholder="Nombre y teléfono de familiar o contacto"
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Condiciones de Salud o Lesiones
          </label>
          <textarea
            v-model="healthConditions"
            rows="2"
            placeholder="Antecedentes médicos relevantes, cirugías previas o dolores articulares..."
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          ></textarea>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Observaciones Adicionales
          </label>
          <textarea
            v-model="notes"
            rows="2"
            placeholder="Objetivos de entrenamiento, disponibilidad u observaciones del entrenador..."
            class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          ></textarea>
        </div>

        <!-- Footer actions -->
        <div class="pt-4 border-t border-[#E2DFD8] flex items-center justify-end gap-3">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#6B6D74] hover:bg-[#F5F3EE] hover:text-[#191A1D] font-medium transition-colors cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold transition-colors shadow-xs cursor-pointer disabled:opacity-50"
          >
            {{ isSubmitting ? 'Registrando...' : '+ Registrar alumno' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
