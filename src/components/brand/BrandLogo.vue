<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'dark' | 'light' | 'gold'
    compact?: boolean
    size?: 'sm' | 'md' | 'lg' | 'xl' | number
    subtitle?: string
  }>(),
  {
    variant: 'light',
    compact: false,
    size: 'md'
  }
)

const height = computed(() => {
  if (typeof props.size === 'number') return props.size
  switch (props.size) {
    case 'sm': return 20
    case 'lg': return 34
    case 'xl': return 44
    case 'md':
    default: return 26
  }
})

const textColor = computed(() => {
  if (props.variant === 'dark') return '#FFFFFF'
  if (props.variant === 'gold') return '#CEA549'
  return '#191A1D' // graphite for light surfaces
})

const goldColor = computed(() => {
  return '#CEA549'
})
</script>

<template>
  <div class="inline-flex items-center gap-2 select-none">
    <!-- COMPACT SYMBOL: ONLY DOUBLE CHEVRON -->
    <svg
      v-if="compact"
      :height="height"
      viewBox="0 0 46 36"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="flex-shrink-0 transition-transform duration-200 hover:scale-105"
      aria-label="INVERFIT"
    >
      <polygon points="6,2 17,2 30,18 17,34 6,34 19,18" :fill="goldColor" />
      <polygon points="21,2 32,2 45,18 32,34 21,34 34,18" :fill="goldColor" />
    </svg>

    <!-- FULL LOGO: INVERFIT WORDMARK + DOUBLE CHEVRON -->
    <svg
      v-else
      :height="height"
      viewBox="0 0 205 36"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="flex-shrink-0 transition-all duration-200"
      aria-label="INVERFIT"
    >
      <!-- Stylized Italic Wordmark -->
      <text
        x="0"
        y="27"
        font-family="'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
        font-style="italic"
        font-weight="900"
        font-size="26"
        letter-spacing="0.5"
        :fill="textColor"
      >
        INVERFIT
      </text>

      <!-- Iconic Slanted Golden Chevrons -->
      <polygon points="160,5 170,5 183,19 170,33 160,33 173,19" :fill="goldColor" />
      <polygon points="175,5 185,5 198,19 185,33 175,33 188,19" :fill="goldColor" />
    </svg>

    <!-- Optional small brand subtitle/tagline if provided -->
    <span
      v-if="subtitle"
      class="text-[10px] font-semibold uppercase tracking-wider px-1.5 py-0.5 rounded"
      :class="variant === 'dark' ? 'bg-[#24262B] text-[#CEA549]' : 'bg-[#CEA549]/10 text-[#8C6B1B]'"
    >
      {{ subtitle }}
    </span>
  </div>
</template>
