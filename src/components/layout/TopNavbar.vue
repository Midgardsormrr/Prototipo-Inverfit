<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Menu, HelpCircle, LogOut, ExternalLink } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import BrandLogo from '@/components/brand/BrandLogo.vue'

defineProps<{
  pageTitle?: string
}>()

const emit = defineEmits<{
  (e: 'toggle-mobile'): void
}>()

const router = useRouter()
const authStore = useAuthStore()

const showHelpModal = ref(false)
const showUserMenu = ref(false)

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="h-[70px] bg-white border-b border-[#E2DFD8] px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30 shadow-2xs">
    <!-- Left: Hamburger on mobile + Mobile Logo + Clean Workspace Title -->
    <div class="flex items-center gap-3 sm:gap-4">
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

      <!-- Desktop Status / Context Indicator -->
      <div class="hidden sm:flex items-center gap-2.5 px-3 py-1.5 rounded-full bg-[#FAF9F6] border border-[#E2DFD8] text-xs">
        <span class="w-2 h-2 rounded-full bg-[#137333] shadow-xs"></span>
        <span class="font-bold text-[#191A1D]">INVERFIT</span>
        <span class="text-[#DCD8CF]">|</span>
        <span class="text-[#6B6D74]">Sistema de Gestión Operativa</span>
      </div>
    </div>

    <!-- Right: Portal Alumnos, Ayuda, Coach Profile & Logout -->
    <div class="flex items-center gap-3 sm:gap-4 relative">
      <!-- Direct Link to Public Onboarding Portal -->
      <router-link
        to="/registro-alumno"
        target="_blank"
        class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-[#8C6B1B] hover:text-[#614910] hover:bg-[#FAF9F6] rounded-md border border-[#E2DFD8] transition-colors"
        title="Abrir formulario de onboarding público para alumnos"
      >
        <span>Portal Alumnos</span>
        <ExternalLink class="w-3.5 h-3.5" />
      </router-link>

      <!-- Help Button -->
      <button
        @click="showHelpModal = true"
        class="text-[13px] text-[#6B6D74] hover:text-[#191A1D] font-medium transition-colors flex items-center gap-1.5 cursor-pointer py-1.5 px-2 rounded hover:bg-[#FAF9F6]"
      >
        <HelpCircle class="w-4 h-4 text-[#8E9097]" />
        <span>Ayuda</span>
      </button>

      <!-- Coach Avatar & Dropdown Trigger -->
      <div class="relative">
        <button
          @click="showUserMenu = !showUserMenu"
          class="flex items-center gap-2 pl-2 pr-1 py-1 rounded-full hover:bg-[#FAF9F6] transition-colors cursor-pointer border border-transparent hover:border-[#E2DFD8]"
        >
          <div
            class="w-[36px] h-[36px] rounded-full bg-[#191A1D] border border-[#CEA549]/60 flex items-center justify-center text-[14px] font-bold text-[#CEA549] select-none shadow-xs"
            title="Coach Pablo"
          >
            P
          </div>
          <span class="text-xs font-bold text-[#191A1D] hidden md:inline">Pablo</span>
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="showUserMenu"
          class="absolute right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-[#E2DFD8] py-2 z-50 animate-in fade-in zoom-in-95 duration-100"
          @click="showUserMenu = false"
        >
          <div class="px-4 py-2 border-b border-[#ECE8E1]">
            <p class="text-xs font-bold text-[#191A1D]">Coach Pablo</p>
            <p class="text-[11px] text-[#8E9097]">Entrenador Principal · INVERFIT</p>
          </div>

          <div class="py-1">
            <button
              @click="handleLogout"
              class="w-full px-4 py-2 text-left text-xs text-[#D93025] hover:bg-[#FDF2F2] flex items-center gap-2 font-semibold transition-colors cursor-pointer"
            >
              <LogOut class="w-4 h-4" />
              <span>Cerrar sesión</span>
            </button>
          </div>
        </div>
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
            <span class="text-xs font-semibold text-[#6B6D74] uppercase tracking-wider">· Guía Operativa INVERFIT</span>
          </div>
          <button @click="showHelpModal = false" class="text-[#8E9097] hover:text-[#191A1D] text-lg font-bold">✕</button>
        </div>

        <div class="space-y-3 text-sm text-[#4A4C52] leading-relaxed max-h-[60vh] overflow-y-auto pr-1">
          <p>
            Bienvenido al panel centralizado de administración de <strong>INVERFIT</strong> (v1.1). Este sistema permite al entrenador gestionar de forma integral toda la operación del centro deportivo.
          </p>

          <div class="bg-[#FAF9F6] p-3.5 rounded-lg border border-[#E2DFD8] space-y-2">
            <p class="font-semibold text-[#191A1D] flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-[#CEA549]"></span>
              Reglas operativas oficiales:
            </p>
            <ul class="list-disc pl-5 space-y-1 text-xs text-[#555861]">
              <li><strong>Agenda y Horarios:</strong> Bloques oficiales de 75 minutos (06:30 a 20:00 lun-vie) y evaluaciones sabatinas, con capacidad técnica máxima de 4 alumnos por turno.</li>
              <li><strong>Membresías y Cobranza:</strong> Membresía mensual única con ciclo oficial del 1 al 5 de cada mes. Integración con enlace permanente de SumUp y registro de comprobantes con N° de operación.</li>
              <li><strong>Evaluaciones Físicas:</strong> Sistema nativo de baremos M-F, monitoreo de signos vitales, dinamometría y percentiles.</li>
              <li><strong>Contratos y Firmas:</strong> Registro de términos v1.0 y respaldo de firmas digitales táctiles capturadas desde el portal de onboarding.</li>
            </ul>
          </div>

          <p class="text-xs text-[#8E9097]">
            * Para compartir el enlace de ingreso a nuevos alumnos, utiliza la opción <strong>Portal Alumnos</strong> en la barra superior.
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
