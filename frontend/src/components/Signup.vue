<template>
  <section class="flex items-center justify-center min-h-screen bg-gradient-to-r from-teal-100 to-slate-100">
    <div class="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
      <h2 class="text-3xl font-bold text-gray-800 mb-8">Sign Up</h2>
      <form @submit.prevent="submitForm" class="space-y-4">
        <div>
          <label class="block text-gray-600 text-left">Full Name</label>
          <input data-test="fullname-input" type="text" v-model="fullName"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" required />
        </div>
        <div>
          <label class="block text-gray-600 text-left">Username</label>
          <input data-test="username-input" type="text" v-model="username"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" required />
        </div>
        <div>
          <label class="block text-gray-600 text-left">Bio (Optional)</label>
          <textarea data-test="bio-input" v-model="bio"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows="3"></textarea>
        </div>
        <div>
          <label class="block text-gray-600 text-left">Password</label>
          <input data-test="password-input" type="password" v-model="password"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" required />
        </div>
        <button data-test="signup-button" type="submit"
          class="w-full bg-primary text-white font-semibold px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-200">
          Sign Up
        </button>
      </form>
      <p class="text-gray-600 mt-6">
        Already have an account?
        <router-link to="/login" class="text-blue-500 hover:underline">Login</router-link>
      </p>
    </div>

    <MessageBox v-if="messageBoxVisible" :title="messageTitle" :message="messageText" :type="messageType" />
  </section>
</template>


<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import MessageBox from '@/components/MessageBox.vue'

const router = useRouter()
const authStore = useAuthStore()
const fullName = ref('')
const username = ref('')
const bio = ref('')
const password = ref('')
const messageBoxVisible = ref(false)
const messageTitle = ref('')
const messageText = ref('')
const messageType = ref('info')

const submitForm = async () => {
  try {
    await authStore.signup({
      full_name: fullName.value,
      username: username.value,
      bio: bio.value,
      password: password.value
    })
    router.push('/login')
  } catch (error) {
    showMessageBox('Signup Failed', `There was an error creating your account. Please try again. ${error}`, 'error');
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
