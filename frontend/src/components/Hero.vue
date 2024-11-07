<template>
  <section class="bg-gradient-to-r from-teal-100 to-slate-100 min-h-screen flex items-center">
    <div class="container mx-auto px-6 text-center">
      <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-4 leading-tight">
        Manage All Your Links in One Place
      </h1>
      <p class="text-gray-800 text-lg md:text-xl mb-8 max-w-2xl mx-auto">
        Linktree allows you to share multiple links with your followers from one single URL.
        It's easy, customizable, and perfect for creators, influencers, and businesses.
      </p>
      <router-link v-if="!authStore.isAuthenticated" to="/signup"
        class="inline-block bg-blue-500 text-white font-semibold px-8 py-3 rounded-lg hover:bg-blue-600 transition duration-300 ease-in-out shadow-lg transform hover:scale-105">
        Get Started
      </router-link>
      <router-link v-else to="/profile"
        class="inline-block bg-green-500 text-white font-semibold px-8 py-3 rounded-lg hover:bg-green-600 transition duration-300 ease-in-out shadow-lg transform hover:scale-105">
        Go to Profile
      </router-link>
    </div>
  </section>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  setup() {
    const authStore = useAuthStore();

    return {
      authStore,
    };
  },
}
</script>

<style scoped>
h1 {
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
}

section {
  min-height: 100vh;
}
</style>
