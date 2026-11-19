<script setup lang="ts">
import { ref, watch } from 'vue'
import { CreditCard, Check, ExternalLink, ShieldCheck } from 'lucide-vue-next'
import { useGymStore } from '@/stores/gym'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const gymStore = useGymStore()
const linkInput = ref('')

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    linkInput.value = gymStore.sumupPermanentLink
  }
})

function handleSave() {
  const trimmed = linkInput.value.trim()
  if (!trimmed) {
    gymStore.showToast('El enlace de pago SumUp no puede estar vacío.', 'warning')
    return
  }
  gymStore.updateSumupPermanentLink(trimmed)
  emit('close')
}

function testLink() {
  if (linkInput.value) {
    window.open(linkInput.value, '_blank')
  }
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-xs"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-md w-full border border-[#E2DFD8] overflow-hidden animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-[#CEA549] bg-[#191A1D] px-2 py-0.5 rounded">
            Configuración de Cobro
          </span>
          <h2 class="text-base font-bold text-[#191A1D] mt-1 flex items-center gap-2">
            <CreditCard class="w-4 h-4 text-[#CEA549]" />
            Link Permanente de SumUp
          </h2>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <form @submit.prevent="handleSave" class="p-6 space-y-4 text-xs">
        <p class="text-[#4A4C52] leading-relaxed">
          Pega aquí tu enlace permanente de cobro de <strong>SumUp</strong>. Se usará automáticamente en todos los avisos por correo del 1 al 5, mensajes de WhatsApp y botones de pago del sistema.
        </p>

        <div>
          <label class="block text-[11px] font-semibold text-[#6B6D74] uppercase tracking-wider mb-1.5">
            URL Permanente SumUp *
          </label>
          <input
            v-model="linkInput"
            type="url"
            required
            placeholder="https://link.sumup.com/pay/inverfit"
            class="w-full px-3 py-2.5 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] font-mono bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all"
          />
        </div>

        <div class="p-3 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg space-y-1.5 text-[11px] text-[#6B6D74]">
          <div class="flex items-center gap-1.5 font-bold text-[#191A1D]">
            <ShieldCheck class="w-4 h-4 text-[#CEA549]" /> Beneficio del Link Permanente
          </div>
          <p>
            Tus clientes accederán directamente a la pasarela segura de SumUp para pagar con tarjeta de crédito o débito, sin necesidad de crear links distintos por cada alumno.
          </p>
        </div>

        <!-- Acciones -->
        <div class="pt-2 flex items-center justify-between gap-2">
          <button
            type="button"
            @click="testLink"
            class="px-3 py-2 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-md text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer"
          >
            <ExternalLink class="w-3.5 h-3.5 text-[#8C6B1B]" />
            Probar Link
          </button>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="emit('close')"
              class="px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#6B6D74] hover:bg-[#F5F3EE] font-medium transition-colors cursor-pointer"
            >
              Cancelar
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold transition-colors shadow-xs cursor-pointer flex items-center gap-1"
            >
              <Check class="w-3.5 h-3.5" />
              Guardar Link
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>
