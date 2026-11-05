<script setup lang="ts">
import { useGymStore } from '@/stores/gym'
import { CheckCircle2, AlertCircle, Info, AlertTriangle, X } from 'lucide-vue-next'

const gymStore = useGymStore()

function iconFor(type: string) {
  switch (type) {
    case 'success': return CheckCircle2
    case 'error': return AlertCircle
    case 'warning': return AlertTriangle
    default: return Info
  }
}

function colorClasses(type: string) {
  switch (type) {
    case 'success': return 'bg-emerald-50 text-emerald-900 border-emerald-200'
    case 'error': return 'bg-rose-50 text-rose-900 border-rose-200'
    case 'warning': return 'bg-amber-50 text-amber-900 border-amber-200'
    default: return 'bg-blue-50 text-blue-900 border-blue-200'
  }
}

function iconColor(type: string) {
  switch (type) {
    case 'success': return 'text-emerald-600'
    case 'error': return 'text-rose-600'
    case 'warning': return 'text-amber-600'
    default: return 'text-blue-600'
  }
}
</script>

<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
    <TransitionGroup
      enter-active-class="transform ease-out duration-200 transition"
      enter-from-class="translate-y-2 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-for="toast in gymStore.toasts"
        :key="toast.id"
        class="pointer-events-auto p-4 rounded-lg border shadow-lg flex items-start gap-3 transition-all"
        :class="colorClasses(toast.type)"
      >
        <component :is="iconFor(toast.type)" class="w-5 h-5 flex-shrink-0 mt-0.5" :class="iconColor(toast.type)" />
        <div class="flex-1 text-sm">
          <p v-if="toast.title" class="font-semibold mb-0.5">{{ toast.title }}</p>
          <p class="leading-snug">{{ toast.message }}</p>
        </div>
        <button
          @click="gymStore.dismissToast(toast.id)"
          class="text-neutral-400 hover:text-neutral-700 p-0.5 rounded transition-colors"
          aria-label="Cerrar notificación"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
