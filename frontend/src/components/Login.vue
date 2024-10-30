<template>
  <section class="flex items-center justify-center min-h-screen bg-gradient-to-r from-teal-100 to-slate-100">
    <div class="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
      <h2 class="text-3xl font-bold text-gray-800 mb-8">Login</h2>
      <form @submit.prevent="submitForm" class="space-y-4">
        <div>
          <label class="block text-gray-600 text-left">Username</label>
          <input type="text" v-model="username"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" required />
        </div>
        <div>
          <label class="block text-gray-600 text-left">Password</label>
          <input type="password" v-model="password"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" required />
        </div>
        <button type="submit"
          class="w-full bg-primary text-white font-semibold px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-200">
          Login
        </button>
      </form>

      <MessageBox v-if="messageBoxVisible" :title="messageTitle" :message="messageText" :type="messageType" />

      <p class="text-gray-600 mt-6">
        Don’t have an account?
        <router-link to="/signup" class="text-blue-500 hover:underline">Sign Up</router-link>
      </p>
    </div>
  </section>
</template>


<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import MessageBox from '@/components/MessageBox.vue'

const router = useRouter()
const authStore = useAuthStore()
const username = ref('')
const password = ref('')
const messageBoxVisible = ref(false)
const messageTitle = ref('')
const messageText = ref('')
const messageType = ref('info')

const submitForm = async () => {
  try {
    await authStore.login(username.value, password.value)
    router.push('/profile')
  } catch (error) {
    showMessageBox('Login Failed', error.message || 'Invalid username or password. Please try again.', 'error');
  }
}

const showMessageBox = (title, message, type) => {
  messageTitle.value = title;
  messageText.value = message;
  messageType.value = type;
  messageBoxVisible.value = true;

  setTimeout(() => {
    messageBoxVisible.value = false;
  }, 5000);
}
</script>
