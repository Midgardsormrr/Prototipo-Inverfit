<script setup lang="ts">
import { ref } from 'vue'
import { Menu, HelpCircle, RotateCcw } from 'lucide-vue-next'
import { useGymStore } from '@/stores/gym'
import BrandLogo from '@/components/brand/BrandLogo.vue'

defineProps<{
  pageTitle?: string
}>()

const emit = defineEmits<{
  (e: 'toggle-mobile'): void
}>()

const gymStore = useGymStore()
const showHelpModal = ref(false)

function handleResetDemo() {
  if (confirm('¿Restablecer todos los datos de demostración a su estado inicial de Figma?')) {
    gymStore.resetDemoData()
  }
}
</script>

<template>
  <header class="h-[70px] bg-white border-b border-[#E2DFD8] px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30 shadow-2xs">
    <!-- Left: Hamburger on mobile + Mobile Logo + Search Box -->
    <div class="flex items-center gap-3 sm:gap-4 flex-1 max-w-[500px]">
      <button
        @click="emit('toggle-mobile')"
        class="md:hidden p-2 text-[#191A1D] hover:bg-[#F5F3EE] rounded-md transition-colors"
        aria-label="Abrir menú"
      >
        <Menu class="w-5 h-5" />
      </button>

      <!-- Mobile brand logo (visible only on small screens when sidebar is hidden) -->
      <div class="md:hidden flex items-center">
        <BrandLogo variant="light" :size="18" />
      </div>

      <!-- Search input matching Figma wireframe -->
      <div class="relative w-full max-w-[300px] hidden sm:block">
        <input
          v-model="gymStore.globalSearch"
          type="text"
          placeholder="⌕  Buscar alumno, bloque o pago..."
          class="w-full h-[38px] px-3.5 bg-[#FAF9F6] border border-[#DCD8CF] rounded-[6px] text-[13px] text-[#191A1D] placeholder-[#8E9097] focus:outline-none focus:border-[#CEA549] focus:bg-white focus:ring-2 focus:ring-[#CEA549]/20 transition-all"
        />
        <button
          v-if="gymStore.globalSearch"
          @click="gymStore.globalSearch = ''"
          class="absolute right-2.5 top-1/2 -translate-y-1/2 text-xs text-[#8E9097] hover:text-[#191A1D]"
        >
          ✕
        </button>
      </div>
    </div>

    <!-- Right: Reiniciar demo, Ayuda, Avatar -->
    <div class="flex items-center gap-3 sm:gap-4">
      <button
        @click="handleResetDemo"
        class="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-1 text-[12px] font-medium text-[#6B6D74] bg-[#F5F3EE] hover:bg-[#EAE6DD] hover:text-[#191A1D] rounded-[6px] border border-[#DDD8CE] transition-colors cursor-pointer"
        title="Restaurar datos de ejemplo del prototipo"
      >
        <RotateCcw class="w-3.5 h-3.5 text-[#8E9097]" />
        Reiniciar demo
      </button>

      <button
        @click="showHelpModal = true"
        class="text-[13px] text-[#6B6D74] hover:text-[#191A1D] font-medium transition-colors flex items-center gap-1.5 cursor-pointer"
      >
        <HelpCircle class="w-4 h-4 text-[#8E9097]" />
        Ayuda
      </button>

      <!-- Coach Avatar Circle matching Figma cx=1310, cy=35, r=18 with brand styling -->
      <div
        class="w-[36px] h-[36px] rounded-full bg-[#191A1D] border border-[#CEA549]/50 flex items-center justify-center text-[14px] font-bold text-[#CEA549] select-none cursor-pointer hover:border-[#CEA549] shadow-xs transition-colors"
        title="Pablo — Entrenador (INVERFIT)"
      >
        P
      </div>
    </div>

    <!-- Help Modal -->
    <div
      v-if="showHelpModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
      @click.self="showHelpModal = false"
    >
      <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] p-6 space-y-4 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between border-b pb-3 border-[#E2DFD8]">
          <div class="flex items-center gap-2">
            <BrandLogo variant="light" :size="20" />
            <span class="text-xs font-semibold text-[#6B6D74] uppercase tracking-wider">· Guía</span>
          </div>
          <button @click="showHelpModal = false" class="text-[#8E9097] hover:text-[#191A1D] text-lg font-bold">✕</button>
        </div>

        <div class="space-y-3 text-sm text-[#4A4C52] leading-relaxed max-h-[60vh] overflow-y-auto pr-1">
          <p>
            Bienvenido al panel administrativo de gestión de gimnasio para el entrenador de <strong>INVERFIT</strong>. Esta aplicación replica con fidelidad visual y funcional los 7 wireframes de Figma adaptados a la identidad corporativa.
          </p>

          <div class="bg-[#FAF9F6] p-3.5 rounded-lg border border-[#E2DFD8] space-y-2">
            <p class="font-semibold text-[#191A1D] flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-[#CEA549]"></span>
              Reglas operativas del negocio:
            </p>
            <ul class="list-disc pl-5 space-y-1 text-xs text-[#555861]">
              <li><strong>Agenda del coach:</strong> El agendamiento es administrado 100% por el entrenador en bloques fijos semanales (1 a 4 sesiones semanales).</li>
              <li><strong>Disponibilidad y duración:</strong> Disponibilidad configurable por día. Sesiones de 75 minutos predeterminadas (libres y editables), sin descansos obligatorios.</li>
              <li><strong>Horario operativo:</strong> 06:30 a 22:00.</li>
              <li><strong>Simulaciones:</strong> El envío de enlaces de pago (<code class="bg-[#ECE8E1] px-1 py-0.5 rounded text-[11px]">https://demo.inverfit.invalid/checkout/...</code>), recordatorios por WhatsApp/correo y comprobantes son simulaciones funcionales con confirmaciones visuales.</li>
              <li><strong>Persistencia:</strong> Todos los cambios se guardan localmente en localStorage. Puedes reiniciar la demostración en cualquier momento con el botón "Reiniciar demo".</li>
            </ul>
          </div>

          <p class="text-xs text-[#8E9097]">
            * Se resguardan los datos mediante registros ficticios (Alumno de ejemplo 01 a 06).
          </p>
        </div>

        <div class="pt-2 flex justify-end">
          <button
            @click="showHelpModal = false"
            class="px-5 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-sm font-bold transition-colors shadow-xs cursor-pointer"
          >
            Entendido
          </button>
        </div>
      </div>
    </div>
  </header>
</template>
