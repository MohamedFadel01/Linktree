<template>
  <div id="app">
    <Navbar />
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'
import Navbar from './components/Navbar.vue'
import axios from 'axios';

const authStore = useAuthStore()

onMounted(() => {
  if (authStore.token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${authStore.token}`
  }
})
</script>
