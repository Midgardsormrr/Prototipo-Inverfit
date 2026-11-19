<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import BrandLogo from '@/components/brand/BrandLogo.vue'
import { Lock, Mail, Eye, EyeOff, ShieldCheck, ArrowRight, AlertCircle } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('pablo@inverfit.cl')
const password = ref('inverfit2026')
const showPassword = ref(false)
const rememberMe = ref(true)
const isLoading = ref(false)
const errorMsg = ref<string | null>(null)

function fillPabloCredentials() {
  email.value = 'pablo@inverfit.cl'
  password.value = 'inverfit2026'
  errorMsg.value = null
}

async function handleLogin() {
  errorMsg.value = null
  if (!email.value || !password.value) {
    errorMsg.value = 'Por favor ingresa tu correo y contraseña.'
    return
  }

  isLoading.value = true
  // Small simulated latency for smooth UX
  await new Promise(r => setTimeout(r, 300))

  const success = authStore.login(email.value, password.value)
  isLoading.value = false

  if (success) {
    const redirect = (route.query.redirect as string) || '/inicio'
    router.push(redirect)
  } else {
    errorMsg.value = authStore.authError || 'Credenciales inválidas.'
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#111215] text-[#FAF9F6] flex flex-col justify-between items-center p-4 sm:p-6 relative overflow-hidden font-sans">
    <!-- Background glow / subtle branding accents -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-[#CEA549]/10 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-[#CEA549]/10 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Header: Simple logo center -->
    <div class="pt-6 sm:pt-10 z-10">
      <BrandLogo variant="dark" :size="32" />
    </div>

    <!-- Main Card -->
    <div class="w-full max-w-[440px] z-10 my-auto">
      <div class="bg-[#191A1D] border border-[#2D3038] rounded-2xl p-6 sm:p-8 shadow-2xl relative">
        <div class="absolute top-0 left-8 right-8 h-[2px] bg-gradient-to-r from-transparent via-[#CEA549] to-transparent"></div>

        <!-- Title -->
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-[#24262B] border border-[#CEA549]/40 mb-3 shadow-inner">
            <ShieldCheck class="w-6 h-6 text-[#CEA549]" />
          </div>
          <h1 class="text-2xl font-bold tracking-tight text-white">Panel del Entrenador</h1>
          <p class="text-xs text-[#8E9097] mt-1.5">
            Ingreso exclusivo para gestión operativa de INVERFIT
          </p>
        </div>

        <!-- Quick Access Helper Chip -->
        <div class="mb-5 p-3 rounded-lg bg-[#24262B]/80 border border-[#3A3D47] flex items-center justify-between text-xs">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-[#137333] animate-pulse"></span>
            <span class="text-[#A0A2A8]">Acceso rápido: <strong>Pablo</strong></span>
          </div>
          <button
            type="button"
            @click="fillPabloCredentials"
            class="text-[#CEA549] hover:underline font-bold text-[11px] cursor-pointer"
          >
            Autocompletar
          </button>
        </div>

        <!-- Error Alert -->
        <div
          v-if="errorMsg"
          class="mb-4 p-3 rounded-lg bg-[#FEF7E0]/10 border border-[#CEA549]/40 text-[#CEA549] text-xs flex items-center gap-2"
        >
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <!-- Email Field -->
          <div>
            <label class="block text-xs font-semibold text-[#A0A2A8] mb-1.5 uppercase tracking-wider">
              Correo Electrónico
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-[#6B6D74]">
                <Mail class="w-4 h-4" />
              </span>
              <input
                v-model="email"
                type="text"
                autocomplete="email"
                required
                placeholder="pablo@inverfit.cl"
                class="w-full h-11 pl-10 pr-3.5 bg-[#141518] border border-[#2D3038] rounded-lg text-sm text-white placeholder-[#5A5D66] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549] transition-all"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-xs font-semibold text-[#A0A2A8] uppercase tracking-wider">
                Contraseña
              </label>
            </div>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-[#6B6D74]">
                <Lock class="w-4 h-4" />
              </span>
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                placeholder="••••••••••••"
                class="w-full h-11 pl-10 pr-10 bg-[#141518] border border-[#2D3038] rounded-lg text-sm text-white placeholder-[#5A5D66] focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549] transition-all font-mono"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-[#6B6D74] hover:text-[#A0A2A8] cursor-pointer"
                aria-label="Alternar visibilidad"
              >
                <component :is="showPassword ? EyeOff : Eye" class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Remember me & Hint -->
          <div class="flex items-center justify-between text-xs pt-1">
            <label class="flex items-center gap-2 cursor-pointer text-[#8E9097] hover:text-[#C5C7CE]">
              <input
                v-model="rememberMe"
                type="checkbox"
                class="rounded border-[#2D3038] bg-[#141518] text-[#CEA549] focus:ring-[#CEA549]"
              />
              <span>Recordar sesión</span>
            </label>
            <span class="text-[#6B6D74] text-[11px]">Coach Pablo · INVERFIT</span>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full h-11 mt-2 bg-[#CEA549] hover:bg-[#BA923B] disabled:opacity-50 text-[#0B0B0D] font-bold rounded-lg text-sm transition-all flex items-center justify-center gap-2 cursor-pointer shadow-md hover:shadow-lg active:scale-[0.99]"
          >
            <span v-if="isLoading" class="w-4 h-4 border-2 border-black/40 border-t-black rounded-full animate-spin"></span>
            <span v-else>Iniciar Sesión</span>
            <ArrowRight v-if="!isLoading" class="w-4 h-4" />
          </button>
        </form>

        <!-- Divider -->
        <div class="mt-6 pt-5 border-t border-[#2D3038] text-center">
          <router-link
            to="/registro-alumno"
            class="text-xs text-[#8E9097] hover:text-[#CEA549] transition-colors inline-flex items-center gap-1.5"
          >
            <span>¿Eres alumno? Ir al registro y consentimiento público</span>
            <ArrowRight class="w-3 h-3" />
          </router-link>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="pb-6 text-center text-xs text-[#5A5D66] z-10 space-y-1">
      <p class="font-medium">INVERFIT · Centro de Entrenamiento Personalizado</p>
      <p class="text-[11px] text-[#44464E]">Versión 1.1 · Despliegue de Producción Oficial</p>
    </div>
  </div>
</template>
