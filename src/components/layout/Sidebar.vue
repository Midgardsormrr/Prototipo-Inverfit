<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import {
  LayoutDashboard,
  Calendar,
  Users,
  CreditCard,
  ClipboardCheck,
  FileText,
  LogOut,
  X
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import BrandLogo from '@/components/brand/BrandLogo.vue'

defineProps<{
  mobileOpen?: boolean
}>()

const emit = defineEmits<{
  (e: 'close-mobile'): void
}>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const navItems = [
  { name: 'Inicio', path: '/inicio', icon: LayoutDashboard },
  { name: 'Agenda', path: '/agenda', icon: Calendar },
  { name: 'Alumnos', path: '/alumnos', icon: Users, matchPrefix: '/alumnos' },
  { name: 'Membresías y pagos', path: '/membresias-pagos', icon: CreditCard },
  { name: 'Evaluaciones', path: '/evaluaciones', icon: ClipboardCheck },
  { name: 'Contratos', path: '/contratos', icon: FileText }
]

function isActive(item: typeof navItems[0]) {
  if (item.matchPrefix) {
    return route.path.startsWith(item.matchPrefix)
  }
  return route.path === item.path
}

function navigate(path: string) {
  router.push(path)
  emit('close-mobile')
}

function handleLogout() {
  if (confirm('¿Cerrar sesión de Coach Pablo?')) {
    authStore.logout()
    router.push('/login')
  }
}
</script>

<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 w-[232px] bg-[#191A1D] border-r border-[#282A30] flex flex-col justify-between transition-transform duration-200 ease-in-out md:translate-x-0 shadow-xl md:shadow-none"
    :class="mobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
  >
    <!-- Top section -->
    <div>
      <!-- Header with INVERFIT Logo -->
      <div class="px-5 pt-6 pb-4 relative">
        <div class="flex items-center justify-between">
          <div class="cursor-pointer" @click="navigate('/inicio')">
            <BrandLogo variant="dark" :size="24" />
            <p class="text-[11px] font-medium text-[#8E9097] tracking-wider uppercase mt-1">
              Panel del entrenador
            </p>
          </div>
          <button
            v-if="mobileOpen"
            @click="emit('close-mobile')"
            class="md:hidden p-1.5 text-[#8E9097] hover:text-white hover:bg-[#24262B] rounded transition-colors"
            aria-label="Cerrar menú"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="h-[1px] bg-[#282A30] mt-4 w-full"></div>
      </div>

      <!-- Navigation items -->
      <nav class="px-3 space-y-1.5 mt-2">
        <button
          v-for="item in navItems"
          :key="item.path"
          @click="navigate(item.path)"
          class="w-full h-[42px] rounded-[6px] px-3 flex items-center gap-3 relative text-[14px] transition-all cursor-pointer text-left group"
          :class="[
            isActive(item)
              ? 'bg-[#24262B] text-[#CEA549] font-bold shadow-xs'
              : 'text-[#A0A2A8] font-normal hover:bg-[#24262B]/60 hover:text-white'
          ]"
        >
          <!-- Active vertical bar on the left matching Figma & brand gold -->
          <span
            v-if="isActive(item)"
            class="absolute left-0 top-1 bottom-1 w-[3.5px] bg-[#CEA549] rounded-r-[2px]"
          ></span>

          <!-- Circle indicator / icon matching Figma with gold accent -->
          <div
            class="w-[9px] h-[9px] rounded-full transition-colors flex-shrink-0"
            :class="isActive(item)
              ? 'bg-[#CEA549] ring-2 ring-[#CEA549]/30'
              : 'bg-[#3A3C43] group-hover:bg-[#5A5D66]'"
          ></div>

          <span class="truncate">{{ item.name }}</span>
        </button>
      </nav>
    </div>

    <!-- Bottom Coach Profile & Logout -->
    <div class="border-t border-[#282A30] p-4 bg-[#141518]">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-[36px] h-[36px] rounded-full bg-[#24262B] border border-[#CEA549]/50 flex items-center justify-center text-[#CEA549] font-bold text-[14px] shadow-xs">
            P
          </div>
          <div>
            <p class="text-[13px] font-semibold text-white leading-none">Pablo</p>
            <p class="text-[11px] text-[#8E9097] mt-1 leading-none">Entrenador</p>
          </div>
        </div>

        <button
          @click="handleLogout"
          title="Cerrar sesión"
          class="p-2 text-[#8E9097] hover:text-[#D93025] hover:bg-[#24262B] rounded transition-colors cursor-pointer"
          aria-label="Cerrar sesión"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>

      <div class="mt-3 pt-2 text-[10px] text-[#6B6D74] text-center border-t border-[#282A30] flex items-center justify-center gap-1.5 font-medium">
        <span class="tracking-wider text-[#A0A2A8] font-bold">INVERFIT</span>
        <span class="text-[#CEA549] font-black">❯❯</span>
        <span class="text-[#CEA549] font-bold">v1.1</span>
      </div>
    </div>
  </aside>
</template>
