<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import TopNavbar from '@/components/layout/TopNavbar.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const mobileOpen = ref(false)

function toggleMobile() {
  mobileOpen.value = !mobileOpen.value
}

function closeMobile() {
  mobileOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-[#F5F3EE] text-[#191A1D] flex flex-col font-sans">
    <!-- Sidebar -->
    <Sidebar :mobile-open="mobileOpen" @close-mobile="closeMobile" />

    <!-- Backdrop for mobile drawer -->
    <div
      v-if="mobileOpen"
      @click="closeMobile"
      class="fixed inset-0 bg-black/40 z-30 md:hidden backdrop-blur-xs transition-opacity"
    ></div>

    <!-- Main Container -->
    <div class="flex-1 flex flex-col md:pl-[232px] transition-all">
      <TopNavbar @toggle-mobile="toggleMobile" />

      <main class="flex-1 p-4 sm:p-6 lg:p-8">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Global Toast Notifications -->
    <ToastContainer />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
