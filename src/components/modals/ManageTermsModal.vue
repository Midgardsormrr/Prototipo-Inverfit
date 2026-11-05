<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useGymStore } from '@/stores/gym'
import { contractService } from '@/services'
import { AlertTriangle } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const gymStore = useGymStore()

const version = ref('v1.0 (2026)')
const title = ref('Términos y Condiciones del Servicio de Entrenamiento — INVERFIT')
const content = ref('')
const isSubmitting = ref(false)

async function loadCurrentTerms() {
  const current = await contractService.getCurrentTerms()
  if (current) {
    version.value = current.version
    title.value = current.title
    content.value = current.content
  }
}

onMounted(() => {
  loadCurrentTerms()
})

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    loadCurrentTerms()
  }
})

async function handleSave() {
  if (!version.value.trim() || !content.value.trim()) {
    gymStore.showToast('Versión y contenido son requeridos.', 'warning')
    return
  }

  isSubmitting.value = true
  try {
    await gymStore.saveTerms(version.value.trim(), title.value.trim(), content.value.trim())
    emit('saved')
    emit('close')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error al guardar versión de términos.'
    gymStore.showToast(msg, 'error')
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
    <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full border border-[#E2DFD8] overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <h2 class="text-lg font-bold text-[#191A1D]">Gestionar Términos y Condiciones</h2>
          <p class="text-xs text-[#6B6D74] mt-0.5">Control de versiones del acuerdo de servicio de INVERFIT.</p>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSave" class="p-6 space-y-4 max-h-[75vh] overflow-y-auto">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Versión
            </label>
            <input
              v-model="version"
              type="text"
              placeholder="v1.1 (2026)"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
              required
            />
          </div>

          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Título del Documento
            </label>
            <input
              v-model="title"
              type="text"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-sm text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
            Contenido y Cláusulas
          </label>
          <textarea
            v-model="content"
            rows="12"
            class="w-full px-3 py-2 font-mono text-xs text-[#191A1D] border border-[#DCD8CF] rounded-md leading-relaxed focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
            required
          ></textarea>
        </div>

        <div class="p-3 bg-[#FEF7E0] border border-[#F0D878] rounded-md text-xs text-[#8C6B1B] flex items-start gap-2">
          <AlertTriangle class="w-4 h-4 flex-shrink-0 mt-0.5 text-[#B06000]" />
          <span>
            <strong>Aviso legal:</strong> Mecanismo de demostración. No constituye validez jurídica ni firma electrónica avanzada conforme a la Ley chilena Nº 19.799. Requiere validación posterior con el cliente.
          </span>
        </div>

        <!-- Footer -->
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
            {{ isSubmitting ? 'Guardando...' : 'Guardar versión vigente' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
